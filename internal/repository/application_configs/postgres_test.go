package application_configs

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	dbdriver "github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	if os.Getenv("CI") == "true" && os.Getenv("DOCKER_AVAILABLE") != "true" {
		os.Exit(0)
	}

	code := m.Run()

	testkit.CleanupSharedContainers()
	os.Exit(code)
}

func setupTestDB(t *testing.T) (*dbdriver.Postgres, *testkit.DBFixtures) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	ctx := context.Background()

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

	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	config := map[string]interface{}{
		"max_tokens": 1000,
		"enabled":    true,
	}
	created, err := repo.Create(ctx, appID, orgID, "production", config)
	require.NoError(t, err)

	cfg, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, cfg.ID)
	assert.Equal(t, appID, cfg.AppID)
	assert.Equal(t, orgID, cfg.OrgID)
	assert.Equal(t, "production", cfg.Environment)
	assert.Equal(t, float64(1000), cfg.Config["max_tokens"])
	assert.Equal(t, true, cfg.Config["enabled"])
	assert.NotZero(t, cfg.CreatedAt)
	assert.NotZero(t, cfg.UpdatedAt)
}

func TestPostgresRepo_GetByEnv(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	config := map[string]interface{}{
		"max_tokens": 2000,
		"model":      "gpt-4",
	}
	_, err := repo.Create(ctx, appID, orgID, "staging", config)
	require.NoError(t, err)

	cfg, err := repo.GetByEnv(ctx, appID, "staging")
	require.NoError(t, err)
	assert.Equal(t, appID, cfg.AppID)
	assert.Equal(t, "staging", cfg.Environment)
	assert.Equal(t, float64(2000), cfg.Config["max_tokens"])
	assert.Equal(t, "gpt-4", cfg.Config["model"])
}

func TestPostgresRepo_ListByAppID(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	prodConfig := map[string]interface{}{"env": "prod"}
	stagingConfig := map[string]interface{}{"env": "staging"}

	_, err := repo.Create(ctx, appID, orgID, "production", prodConfig)
	require.NoError(t, err)
	_, err = repo.Create(ctx, appID, orgID, "staging", stagingConfig)
	require.NoError(t, err)

	configs, err := repo.ListByAppID(ctx, appID)
	require.NoError(t, err)
	assert.Len(t, configs, 2)

	envs := make(map[string]bool)
	for _, cfg := range configs {
		assert.Equal(t, appID, cfg.AppID)
		envs[cfg.Environment] = true
	}
	assert.True(t, envs["production"])
	assert.True(t, envs["staging"])
}

func TestPostgresRepo_Create(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	config := map[string]interface{}{
		"rate_limit":   100,
		"timeout_ms":   5000,
		"allowed_ips":  []interface{}{"192.168.1.1", "10.0.0.1"},
		"feature_flag": true,
	}

	cfg, err := repo.Create(ctx, appID, orgID, "development", config)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, cfg.ID)
	assert.Equal(t, appID, cfg.AppID)
	assert.Equal(t, orgID, cfg.OrgID)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, float64(100), cfg.Config["rate_limit"])
	assert.Equal(t, float64(5000), cfg.Config["timeout_ms"])
	assert.Equal(t, true, cfg.Config["feature_flag"])
	assert.NotZero(t, cfg.CreatedAt)
	assert.NotZero(t, cfg.UpdatedAt)
}

func TestPostgresRepo_Update(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	originalConfig := map[string]interface{}{
		"version": 1,
		"enabled": true,
	}
	created, err := repo.Create(ctx, appID, orgID, "production", originalConfig)
	require.NoError(t, err)

	updatedConfig := map[string]interface{}{
		"version": 2,
		"enabled": false,
		"new_key": "new_value",
	}

	updated, err := repo.Update(ctx, created.ID, updatedConfig)
	require.NoError(t, err)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, float64(2), updated.Config["version"])
	assert.Equal(t, false, updated.Config["enabled"])
	assert.Equal(t, "new_value", updated.Config["new_key"])
	assert.NotZero(t, updated.UpdatedAt)
}

func TestPostgresRepo_Delete(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	config := map[string]interface{}{"test": "data"}
	created, err := repo.Create(ctx, appID, orgID, "production", config)
	require.NoError(t, err)

	err = repo.Delete(ctx, created.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, created.ID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows), "expected sql.ErrNoRows, got %v", err)
}

func TestPostgresRepo_GetByID_NotFound(t *testing.T) {
	pg, _ := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows), "expected sql.ErrNoRows, got %v", err)
}

func TestPostgresRepo_GetByEnv_NotFound(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	_, appID := fixtures.CreateTestOrgAndApp(t)

	_, err := repo.GetByEnv(ctx, appID, "non-existent-env")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows), "expected sql.ErrNoRows, got %v", err)
}

func TestPostgresRepo_Create_Validation(t *testing.T) {
	pg, _ := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	config := map[string]interface{}{"test": "data"}

	_, err := repo.Create(ctx, uuid.Nil, uuid.New(), "production", config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "appID cannot be nil")

	_, err = repo.Create(ctx, uuid.New(), uuid.Nil, "production", config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "orgID cannot be nil")

	_, err = repo.Create(ctx, uuid.New(), uuid.New(), "", config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "environment cannot be empty")
}

func TestPostgresRepo_Update_Validation(t *testing.T) {
	pg, _ := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	config := map[string]interface{}{"test": "data"}

	_, err := repo.Update(ctx, uuid.Nil, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id cannot be nil")
}

func TestPostgresRepo_Delete_Validation(t *testing.T) {
	pg, _ := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	err := repo.Delete(ctx, uuid.Nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id cannot be nil")
}

func TestPostgresRepo_Create_UniqueConstraint(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	config := map[string]interface{}{"test": "data"}
	_, err := repo.Create(ctx, appID, orgID, "production", config)
	require.NoError(t, err)

	_, err = repo.Create(ctx, appID, orgID, "production", config)
	assert.Error(t, err)
}

func TestPostgresRepo_EmptyConfig(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	emptyConfig := map[string]interface{}{}
	cfg, err := repo.Create(ctx, appID, orgID, "production", emptyConfig)
	require.NoError(t, err)
	assert.NotNil(t, cfg.Config)
	assert.Len(t, cfg.Config, 0)
}

func TestPostgresRepo_ComplexNestedConfig(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	complexConfig := map[string]interface{}{
		"database": map[string]interface{}{
			"host":     "localhost",
			"port":     5432,
			"settings": map[string]interface{}{"pool_size": 10, "timeout": 30},
		},
		"features": []interface{}{"auth", "logging", "metrics"},
		"metadata": map[string]interface{}{"version": "1.0.0", "tags": []interface{}{"prod", "us-east"}},
	}

	cfg, err := repo.Create(ctx, appID, orgID, "production", complexConfig)
	require.NoError(t, err)
	assert.NotNil(t, cfg.Config["database"])
	assert.NotNil(t, cfg.Config["features"])
	assert.NotNil(t, cfg.Config["metadata"])

	dbConfig := cfg.Config["database"].(map[string]interface{})
	assert.Equal(t, "localhost", dbConfig["host"])
	assert.Equal(t, float64(5432), dbConfig["port"])

	features := cfg.Config["features"].([]interface{})
	assert.Len(t, features, 3)
	assert.Contains(t, features, "auth")
}
