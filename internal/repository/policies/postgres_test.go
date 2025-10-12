package policies

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	dbdriver "github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	testkit.CleanupSharedContainers()
	os.Exit(code)
}

func setupTestDB(t *testing.T) (*dbdriver.Postgres, *testkit.DBFixtures) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set up shared test containers
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	ctx := context.Background()

	// Set up database connection
	pg, err := dbdriver.NewPostgresDriver(ctx, pgConnStr)
	require.NoError(t, err)
	t.Cleanup(func() { pg.Pool.Close() })

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	return pg, fixtures
}

func TestPostgresRepo_GetByID(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test data
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	policyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 100}`)

	// Test GetByID
	policy, err := repo.GetByID(ctx, policyID)
	require.NoError(t, err)
	assert.Equal(t, policyID, policy.ID)
	assert.Equal(t, orgID, policy.OrgID)
	assert.Equal(t, appID, policy.AppID)
	assert.Equal(t, model.PolicyTypeRateLimit, policy.PolicyType)
	assert.True(t, policy.Enabled)
	assert.NotZero(t, policy.CreatedAt)
	assert.NotZero(t, policy.UpdatedAt)
}

func TestPostgresRepo_ListByAppID(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test data
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 100}`)
	fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeTokenLimit, `{"max_prompt_tokens": 1000}`)

	// Test ListByAppID
	policies, err := repo.ListByAppID(ctx, appID, 100, 0)
	require.NoError(t, err)
	assert.Len(t, policies, 2)

	// Check that policies are returned
	foundRateLimit := false
	foundTokenLimit := false
	for _, policy := range policies {
		if policy.PolicyType == model.PolicyTypeRateLimit {
			foundRateLimit = true
			assert.Equal(t, 100.0, policy.Config["requests_per_minute"])
		} else if policy.PolicyType == model.PolicyTypeTokenLimit {
			foundTokenLimit = true
			assert.Equal(t, 1000.0, policy.Config["max_prompt_tokens"])
		}
	}
	assert.True(t, foundRateLimit)
	assert.True(t, foundTokenLimit)
}

func TestPostgresRepo_ListEnabledByAppID(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test data
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	enabledPolicyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 100}`)
	disabledPolicyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeTokenLimit, `{"max_prompt_tokens": 1000}`)

	// Disable one policy
	err := repo.Disable(ctx, disabledPolicyID)
	require.NoError(t, err)

	// Test ListEnabledByAppID
	policies, err := repo.ListEnabledByAppID(ctx, appID, 100, 0)
	require.NoError(t, err)
	assert.Len(t, policies, 1)
	assert.Equal(t, enabledPolicyID, policies[0].ID)
	assert.Equal(t, model.PolicyTypeRateLimit, policies[0].PolicyType)
}

func TestPostgresRepo_GetByType(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test data
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 100}`)
	fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 200}`)

	// Test GetByType
	policies, err := repo.GetByType(ctx, appID, model.PolicyTypeRateLimit)
	require.NoError(t, err)
	assert.Len(t, policies, 2)

	for _, policy := range policies {
		assert.Equal(t, model.PolicyTypeRateLimit, policy.PolicyType)
		assert.Contains(t, []float64{100.0, 200.0}, policy.Config["requests_per_minute"])
	}
}

func TestPostgresRepo_Create(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test data
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	// Test Create
	config := map[string]any{"requests_per_minute": 150}
	policy, err := repo.Create(ctx, orgID, appID, model.PolicyTypeRateLimit, config, true)
	require.NoError(t, err)
	assert.Equal(t, orgID, policy.OrgID)
	assert.Equal(t, appID, policy.AppID)
	assert.Equal(t, model.PolicyTypeRateLimit, policy.PolicyType)
	assert.Equal(t, config, policy.Config)
	assert.True(t, policy.Enabled)
	assert.NotZero(t, policy.CreatedAt)
	assert.NotZero(t, policy.UpdatedAt)
}

func TestPostgresRepo_Update(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test data
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	policyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 100}`)

	// Test Update
	newConfig := map[string]any{"requests_per_minute": 200}
	updatedPolicy, err := repo.Update(ctx, policyID, model.PolicyTypeRateLimit, newConfig, false)
	require.NoError(t, err)
	assert.Equal(t, policyID, updatedPolicy.ID)
	assert.Equal(t, model.PolicyTypeRateLimit, updatedPolicy.PolicyType)
	assert.Equal(t, newConfig, updatedPolicy.Config)
	assert.False(t, updatedPolicy.Enabled)
}

func TestPostgresRepo_Enable_Disable(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test data
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	policyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 100}`)

	// Test Disable
	err := repo.Disable(ctx, policyID)
	require.NoError(t, err)

	policy, err := repo.GetByID(ctx, policyID)
	require.NoError(t, err)
	assert.False(t, policy.Enabled)

	// Test Enable
	err = repo.Enable(ctx, policyID)
	require.NoError(t, err)

	policy, err = repo.GetByID(ctx, policyID)
	require.NoError(t, err)
	assert.True(t, policy.Enabled)
}

func TestPostgresRepo_Delete(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test data
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	policyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 100}`)

	// Test Delete
	err := repo.Delete(ctx, policyID)
	require.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(ctx, policyID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows), "expected sql.ErrNoRows, got %v", err)
}

func TestPostgresRepo_GetByID_NotFound(t *testing.T) {
	pg, _ := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Test GetByID with non-existent ID
	_, err := repo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows), "expected sql.ErrNoRows, got %v", err)
}

func TestPostgresRepo_Create_Validation(t *testing.T) {
	pg, _ := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	config := map[string]any{"test": "config"}

	// Test Create with nil orgID
	_, err := repo.Create(ctx, uuid.Nil, uuid.New(), model.PolicyTypeRateLimit, config, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "orgID cannot be nil")

	// Test Create with nil appID
	_, err = repo.Create(ctx, uuid.New(), uuid.Nil, model.PolicyTypeRateLimit, config, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "appID cannot be nil")

	// Test Create with empty policyType
	_, err = repo.Create(ctx, uuid.New(), uuid.New(), "", config, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "policyType cannot be empty")
}

func TestPostgresRepo_Update_Validation(t *testing.T) {
	pg, _ := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	config := map[string]any{"test": "config"}

	// Test Update with nil ID
	_, err := repo.Update(ctx, uuid.Nil, model.PolicyTypeRateLimit, config, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id cannot be nil")

	// Test Update with empty policyType
	_, err = repo.Update(ctx, uuid.New(), "", config, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "policyType cannot be empty")
}

func TestPostgresRepo_Enable_Disable_Validation(t *testing.T) {
	pg, _ := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Test Enable with nil ID
	err := repo.Enable(ctx, uuid.Nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id cannot be nil")

	// Test Disable with nil ID
	err = repo.Disable(ctx, uuid.Nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id cannot be nil")
}

func TestPostgresRepo_Delete_Validation(t *testing.T) {
	pg, _ := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Test Delete with nil ID
	err := repo.Delete(ctx, uuid.Nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id cannot be nil")
}
