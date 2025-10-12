package applications

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

	// Test GetByID
	app, err := repo.GetByID(ctx, appID)
	require.NoError(t, err)
	assert.Equal(t, appID, app.ID)
	assert.Equal(t, orgID, app.OrgID)
	assert.Equal(t, "test-app", app.Name)
	assert.NotNil(t, app.Description)
	assert.Equal(t, "Test application for integration tests", *app.Description)
	assert.NotZero(t, app.CreatedAt)
	assert.NotZero(t, app.UpdatedAt)
}

func TestPostgresRepo_GetByName(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test data
	orgID, _ := fixtures.CreateTestOrgAndApp(t)

	// Test GetByName
	app, err := repo.GetByName(ctx, orgID, "test-app")
	require.NoError(t, err)
	assert.Equal(t, orgID, app.OrgID)
	assert.Equal(t, "test-app", app.Name)
	assert.NotNil(t, app.Description)
	assert.Equal(t, "Test application for integration tests", *app.Description)
}

func TestPostgresRepo_ListByOrgID(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test data
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	// Test ListByOrgID
	apps, err := repo.ListByOrgID(ctx, orgID)
	require.NoError(t, err)
	assert.Len(t, apps, 1)
	assert.Equal(t, appID, apps[0].ID)
	assert.Equal(t, "test-app", apps[0].Name)
}

func TestPostgresRepo_Create(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test org
	orgID, _ := fixtures.CreateTestOrgAndApp(t)

	// Test Create
	description := "New test application"
	app, err := repo.Create(ctx, orgID, "new-app", &description)
	require.NoError(t, err)
	assert.Equal(t, orgID, app.OrgID)
	assert.Equal(t, "new-app", app.Name)
	assert.NotNil(t, app.Description)
	assert.Equal(t, description, *app.Description)
	assert.NotZero(t, app.CreatedAt)
	assert.NotZero(t, app.UpdatedAt)
}

func TestPostgresRepo_Update(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test data
	_, appID := fixtures.CreateTestOrgAndApp(t)

	// Test Update
	newDescription := "Updated description"
	updatedApp, err := repo.Update(ctx, appID, "updated-app", &newDescription)
	require.NoError(t, err)
	assert.Equal(t, appID, updatedApp.ID)
	assert.Equal(t, "updated-app", updatedApp.Name)
	assert.NotNil(t, updatedApp.Description)
	assert.Equal(t, newDescription, *updatedApp.Description)
	assert.NotZero(t, updatedApp.UpdatedAt)
}

func TestPostgresRepo_Delete(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test data
	_, appID := fixtures.CreateTestOrgAndApp(t)

	// Test Delete
	err := repo.Delete(ctx, appID)
	require.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(ctx, appID)
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

func TestPostgresRepo_GetByName_NotFound(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test org
	orgID, _ := fixtures.CreateTestOrgAndApp(t)

	// Test GetByName with non-existent name
	_, err := repo.GetByName(ctx, orgID, "non-existent-app")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows), "expected sql.ErrNoRows, got %v", err)
}

func TestPostgresRepo_Create_Validation(t *testing.T) {
	pg, _ := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Test Create with nil orgID
	_, err := repo.Create(ctx, uuid.Nil, "test-app", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "orgID cannot be nil")

	// Test Create with empty name
	_, err = repo.Create(ctx, uuid.New(), "", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name cannot be empty")
}

func TestPostgresRepo_Update_Validation(t *testing.T) {
	pg, _ := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Test Update with nil ID
	_, err := repo.Update(ctx, uuid.Nil, "test-app", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id cannot be nil")

	// Test Update with empty name
	_, err = repo.Update(ctx, uuid.New(), "", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name cannot be empty")
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
