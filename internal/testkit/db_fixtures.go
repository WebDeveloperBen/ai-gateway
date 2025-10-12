package testkit

import (
	"context"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBFixtures provides database setup utilities for integration tests
type DBFixtures struct {
	Queries *db.Queries
	Pool    *pgxpool.Pool
}

// NewDBFixtures creates a new database fixtures helper
func NewDBFixtures(queries *db.Queries, pool *pgxpool.Pool) *DBFixtures {
	return &DBFixtures{Queries: queries, Pool: pool}
}

// CreateTestOrgAndApp creates a test organization and application for testing
func (df *DBFixtures) CreateTestOrgAndApp(t *testing.T) (uuid.UUID, uuid.UUID) {
	t.Helper()
	return df.CreateTestOrgAndAppWithSuffix(t, "")
}

// CreateTestOrgAndAppWithSuffix creates a test organization and application with a suffix
func (df *DBFixtures) CreateTestOrgAndAppWithSuffix(t *testing.T, suffix string) (uuid.UUID, uuid.UUID) {
	t.Helper()
	ctx := context.Background()

	orgName := "test-org"
	appName := "test-app"
	if suffix != "" {
		orgName += "-" + suffix
		appName += "-" + suffix
	}

	// Create test org
	org, err := df.Queries.CreateOrg(ctx, orgName)
	if err != nil {
		t.Fatalf("Failed to create test org: %v", err)
	}

	// Create test app
	app, err := df.Queries.CreateApplication(ctx, db.CreateApplicationParams{
		OrgID:       org.ID,
		Name:        appName,
		Description: stringPtr("Test application for integration tests"),
	})
	if err != nil {
		t.Fatalf("Failed to create test app: %v", err)
	}

	return org.ID, app.ID
}

// CreateTestPolicy creates a test policy for the given org and app
func (df *DBFixtures) CreateTestPolicy(t *testing.T, orgID, appID uuid.UUID, policyType model.PolicyType, config string) uuid.UUID {
	t.Helper()

	policy, err := df.Queries.CreatePolicy(context.Background(), db.CreatePolicyParams{
		OrgID:      orgID,
		AppID:      appID,
		PolicyType: string(policyType),
		Config:     []byte(config),
		Enabled:    true,
	})
	if err != nil {
		t.Fatalf("Failed to create test policy: %v", err)
	}

	return policy.ID
}

// CleanupTestPolicy removes a test policy
func (df *DBFixtures) CleanupTestPolicy(t *testing.T, policyID uuid.UUID) {
	t.Helper()

	err := df.Queries.DeletePolicy(context.Background(), policyID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup test policy %s: %v", policyID, err)
	}
}

// CleanupTestApp removes a test application
func (df *DBFixtures) CleanupTestApp(t *testing.T, appID uuid.UUID) {
	t.Helper()

	err := df.Queries.DeleteApplication(context.Background(), appID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup test app %s: %v", appID, err)
	}
}

// CreateTestAPIKey creates a test API key for the given org
func (df *DBFixtures) CreateTestAPIKey(t *testing.T, orgID uuid.UUID) uuid.UUID {
	t.Helper()

	// Create a user first (required by api_keys table)
	user, err := df.Queries.CreateUser(context.Background(), db.CreateUserParams{
		OrgID: orgID,
		Email: "test@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create API key directly with raw SQL since there's no SQLC query for it
	var apiKeyID uuid.UUID
	err = df.Pool.QueryRow(context.Background(), `
		INSERT INTO api_keys (org_id, user_id, key_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`, orgID, user.ID, "$argon2id$v=19$m=65536,t=3,p=2$dummy$salt$dummyhash").Scan(&apiKeyID)
	if err != nil {
		t.Fatalf("Failed to create test API key: %v", err)
	}

	return apiKeyID
}

// stringPtr returns a pointer to the given string
func stringPtr(s string) *string {
	return &s
}
