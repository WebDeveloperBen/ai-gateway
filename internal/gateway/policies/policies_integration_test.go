package policies_test

import (
	"context"
	"os"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	dbdriver "github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/policies"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// TestMain handles shared container setup and cleanup
func TestMain(m *testing.M) {
	// Skip integration tests if running in CI without Docker
	if os.Getenv("CI") == "true" && os.Getenv("DOCKER_AVAILABLE") != "true" {
		os.Exit(0)
	}

	// Run tests
	code := m.Run()

	// Force cleanup of shared containers at the very end
	// This ensures containers are cleaned up even if tests fail
	testkit.CleanupSharedContainers()
	os.Exit(code)
}

// TestEndToEndPolicyEnforcementWithDatabase tests complete policy enforcement flow with real database
func TestEndToEndPolicyEnforcementWithDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set up shared test containers
	pgConnStr, redisAddr := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	ctx := context.Background()

	// Set up database connection
	pg, err := dbdriver.NewPostgresDriver(ctx, pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	// Set up Redis cache
	cache, err := kv.NewDriver(kv.Config{
		Backend:   kv.BackendRedis,
		RedisAddr: redisAddr,
		RedisPW:   "",
		RedisDB:   0,
	})
	require.NoError(t, err)

	// Create policy engine
	engine := policies.NewEngine(pg.Queries, cache)

	// Create unique test data for this test
	testSuffix := uuid.New().String()[:8]
	orgID, appID := createTestOrgAndAppWithSuffix(t, pg.Queries, testSuffix)
	appIDStr := appID.String()

	// Create test policies in database
	rateLimitID := createTestPolicy(t, pg.Queries, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 1}`)
	tokenLimitID := createTestPolicy(t, pg.Queries, orgID, appID, model.PolicyTypeTokenLimit, `{"max_prompt_tokens": 100}`)

	t.Run("RateLimitBlocksExcessRequests", func(t *testing.T) {
		// Load policies from database
		loadedPolicies, err := engine.LoadPolicies(ctx, appIDStr)
		require.NoError(t, err)
		require.Len(t, loadedPolicies, 2, "Expected 2 policies to be loaded")

		// Debug: print policy types
		for _, p := range loadedPolicies {
			t.Logf("Loaded policy type: %s", p.Type())
		}

		// First request should pass
		preCtx1 := &policies.PreRequestContext{
			AppID: appIDStr,
			Model: "gpt-4",
		}
		err1 := engine.CheckPreRequest(ctx, loadedPolicies, preCtx1)
		t.Logf("First request result: %v", err1)
		require.NoError(t, err1, "First request should pass")

		// Second request should be blocked (rate limited)
		preCtx2 := &policies.PreRequestContext{
			AppID: appIDStr,
			Model: "gpt-4",
		}
		err2 := engine.CheckPreRequest(ctx, loadedPolicies, preCtx2)
		t.Logf("Second request result: %v", err2)
		require.Error(t, err2, "Second request should be rate limited")
	})

	t.Run("TokenLimitBlocksLargeRequests", func(t *testing.T) {
		// Load policies from database
		loadedPolicies, err := engine.LoadPolicies(ctx, appIDStr)
		require.NoError(t, err)
		require.Len(t, loadedPolicies, 2, "Expected 2 policies to be loaded")

		preCtx := &policies.PreRequestContext{
			AppID:           appIDStr,
			Model:           "gpt-4",
			EstimatedTokens: 200, // Exceeds 100 limit
		}

		err = engine.CheckPreRequest(ctx, loadedPolicies, preCtx)
		require.Error(t, err, "Large request should be blocked")
	})

	t.Run("PolicyCachingWorks", func(t *testing.T) {
		// Load policies (should hit DB first time)
		policies1, err := engine.LoadPolicies(ctx, appIDStr)
		require.NoError(t, err)
		require.Len(t, policies1, 2, "Expected 2 policies")

		// Load again (should hit cache)
		policies2, err := engine.LoadPolicies(ctx, appIDStr)
		require.NoError(t, err)
		require.Len(t, policies2, 2, "Expected 2 policies from cache")
	})

	// Clean up test data
	cleanupTestPolicy(t, pg.Queries, rateLimitID)
	cleanupTestPolicy(t, pg.Queries, tokenLimitID)
	cleanupTestApp(t, pg.Queries, appID)
}

// TestPolicyRegistryIntegrationWithDatabase tests registry with database-backed policies
func TestPolicyRegistryIntegrationWithDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set up shared test containers
	pgConnStr, redisAddr := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	ctx := context.Background()

	// Set up database connection
	pg, err := dbdriver.NewPostgresDriver(ctx, pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	// Set up Redis cache
	cache, err := kv.NewDriver(kv.Config{
		Backend:   kv.BackendRedis,
		RedisAddr: redisAddr,
		RedisPW:   "",
		RedisDB:   0,
	})
	require.NoError(t, err)

	// Create policy engine
	engine := policies.NewEngine(pg.Queries, cache)

	// Create unique test data for this test
	testSuffix := uuid.New().String()[:8]
	orgID, appID := createTestOrgAndAppWithSuffix(t, pg.Queries, testSuffix)
	appIDStr := appID.String()

	// Test all built-in policy types
	policyTypes := []model.PolicyType{
		model.PolicyTypeRateLimit,
		model.PolicyTypeTokenLimit,
		model.PolicyTypeModelAllowlist,
		model.PolicyTypeRequestSize,
	}

	for _, policyType := range policyTypes {
		t.Run(string(policyType), func(t *testing.T) {
			config := getTestConfigForType(policyType)
			policyID := createTestPolicy(t, pg.Queries, orgID, appID, policyType, config)

			t.Cleanup(func() {
				cleanupTestPolicy(t, pg.Queries, policyID)
				if err := engine.InvalidateCache(ctx, appIDStr); err != nil {
					t.Logf("Warning: failed to invalidate cache for app %s: %v", appIDStr, err)
				}
			})

			// Clear caches to ensure the latest DB state is fetched for each policy type
			require.NoError(t, engine.InvalidateCache(ctx, appIDStr))

			// Load policies from DB (after cache invalidation)
			policies, err := engine.LoadPolicies(ctx, appIDStr)
			require.NoError(t, err)

			// Should have at least one policy
			require.NotEmpty(t, policies, "Expected at least one policy to be loaded")

			// Verify the policy was created with correct type
			found := false
			for _, policy := range policies {
				if policy.Type() == policyType {
					found = true
					break
				}
			}
			require.True(t, found, "Policy of type %s not found in loaded policies", policyType)

			// Cleanup handled via t.Cleanup to ensure it always runs
		})
	}

	// Clean up
	cleanupTestApp(t, pg.Queries, appID)
}

// Helper functions

func createTestOrgAndApp(t *testing.T, queries *db.Queries) (uuid.UUID, uuid.UUID) {
	t.Helper()
	return createTestOrgAndAppWithSuffix(t, queries, "")
}

func createTestOrgAndAppWithSuffix(t *testing.T, queries *db.Queries, suffix string) (uuid.UUID, uuid.UUID) {
	t.Helper()
	ctx := context.Background()

	orgName := "test-org"
	appName := "test-app"
	if suffix != "" {
		orgName += "-" + suffix
		appName += "-" + suffix
	}

	// Create test org
	org, err := queries.CreateOrg(ctx, orgName)
	require.NoError(t, err)

	// Create test app
	app, err := queries.CreateApplication(ctx, db.CreateApplicationParams{
		OrgID:       org.ID,
		Name:        appName,
		Description: stringPtr("Test application for policy integration tests"),
	})
	require.NoError(t, err)

	return org.ID, app.ID
}

func createTestPolicy(t *testing.T, queries *db.Queries, orgID, appID uuid.UUID, policyType model.PolicyType, config string) uuid.UUID {
	t.Helper()

	policy, err := queries.CreatePolicy(context.Background(), db.CreatePolicyParams{
		OrgID:      orgID,
		PolicyType: string(policyType),
		Config:     []byte(config),
		Enabled:    true,
	})
	require.NoError(t, err)

	// Attach policy to the application
	err = queries.AttachPolicyToApp(context.Background(), db.AttachPolicyToAppParams{
		PolicyID: policy.ID,
		AppID:    appID,
	})
	require.NoError(t, err)

	return policy.ID
}

func cleanupTestPolicy(t *testing.T, queries *db.Queries, policyID uuid.UUID) {
	t.Helper()

	err := queries.DeletePolicy(context.Background(), policyID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup test policy %s: %v", policyID, err)
	}
}

func cleanupTestApp(t *testing.T, queries *db.Queries, appID uuid.UUID) {
	t.Helper()

	err := queries.DeleteApplication(context.Background(), appID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup test app %s: %v", appID, err)
	}
}

func stringPtr(s string) *string {
	return &s
}

func getTestConfigForType(policyType model.PolicyType) string {
	switch policyType {
	case model.PolicyTypeRateLimit:
		return `{"requests_per_minute": 1000}`
	case model.PolicyTypeTokenLimit:
		return `{"max_prompt_tokens": 4000}`
	case model.PolicyTypeModelAllowlist:
		return `{"allowed_model_ids": ["gpt-4", "gpt-3.5-turbo"]}`
	case model.PolicyTypeRequestSize:
		return `{"max_request_bytes": 51200}`
	default:
		return `{}`
	}
}
