package applications

import (
	"context"
	"testing"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMemoryRepo(t *testing.T) {
	repo := NewMemoryRepo()
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.store)
}

func TestMemoryRepo_GetByID(t *testing.T) {
	repo := NewMemoryRepo()
	ctx := context.Background()

	// Test not found
	_, err := repo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "application not found")

	// Create an app
	orgID := uuid.New()
	app, err := repo.Create(ctx, orgID, "test-app", stringPtr("test description"))
	require.NoError(t, err)

	// Test found
	found, err := repo.GetByID(ctx, app.ID)
	assert.NoError(t, err)
	assert.Equal(t, app, found)
}

func TestMemoryRepo_GetByName(t *testing.T) {
	repo := NewMemoryRepo()
	ctx := context.Background()

	orgID := uuid.New()

	// Test not found
	_, err := repo.GetByName(ctx, orgID, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "application not found")

	// Create an app
	app, err := repo.Create(ctx, orgID, "test-app", stringPtr("test description"))
	require.NoError(t, err)

	// Test found
	found, err := repo.GetByName(ctx, orgID, "test-app")
	assert.NoError(t, err)
	assert.Equal(t, app, found)

	// Test wrong org
	wrongOrgID := uuid.New()
	_, err = repo.GetByName(ctx, wrongOrgID, "test-app")
	assert.Error(t, err)
}

func TestMemoryRepo_ListByOrgID(t *testing.T) {
	repo := NewMemoryRepo()
	ctx := context.Background()

	orgID1 := uuid.New()
	orgID2 := uuid.New()

	// Create apps in different orgs
	app1, err := repo.Create(ctx, orgID1, "app1", nil)
	require.NoError(t, err)
	app2, err := repo.Create(ctx, orgID1, "app2", nil)
	require.NoError(t, err)
	app3, err := repo.Create(ctx, orgID2, "app3", nil)
	require.NoError(t, err)

	// List org1 apps
	apps, err := repo.ListByOrgID(ctx, orgID1, 0, 0)
	assert.NoError(t, err)
	assert.Len(t, apps, 2)
	assert.Contains(t, apps, app1)
	assert.Contains(t, apps, app2)

	// List org2 apps
	apps, err = repo.ListByOrgID(ctx, orgID2, 0, 0)
	assert.NoError(t, err)
	assert.Len(t, apps, 1)
	assert.Contains(t, apps, app3)

	// Test limit and offset
	apps, err = repo.ListByOrgID(ctx, orgID1, 1, 0)
	assert.NoError(t, err)
	assert.Len(t, apps, 1)

	apps, err = repo.ListByOrgID(ctx, orgID1, 1, 1)
	assert.NoError(t, err)
	assert.Len(t, apps, 1)
}

func TestMemoryRepo_Create(t *testing.T) {
	repo := NewMemoryRepo()
	ctx := context.Background()

	orgID := uuid.New()

	// Test successful creation
	app, err := repo.Create(ctx, orgID, "test-app", stringPtr("description"))
	assert.NoError(t, err)
	assert.NotNil(t, app)
	assert.Equal(t, orgID, app.OrgID)
	assert.Equal(t, "test-app", app.Name)
	assert.Equal(t, stringPtr("description"), app.Description)
	assert.NotEqual(t, uuid.Nil, app.ID)
	assert.NotZero(t, app.CreatedAt)
	assert.NotZero(t, app.UpdatedAt)

	// Test duplicate name
	_, err = repo.Create(ctx, orgID, "test-app", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "application name already exists")

	// Test nil orgID
	_, err = repo.Create(ctx, uuid.Nil, "test-app2", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "orgID cannot be nil")

	// Test empty name
	_, err = repo.Create(ctx, orgID, "", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name cannot be empty")
}

func TestMemoryRepo_Update(t *testing.T) {
	repo := NewMemoryRepo()
	ctx := context.Background()

	orgID := uuid.New()

	// Create an app
	app, err := repo.Create(ctx, orgID, "test-app", stringPtr("old description"))
	require.NoError(t, err)
	originalUpdatedAt := app.UpdatedAt

	// Wait a bit to ensure updated time changes
	time.Sleep(1 * time.Millisecond)

	// Test successful update
	updated, err := repo.Update(ctx, app.ID, "new-name", stringPtr("new description"))
	assert.NoError(t, err)
	assert.Equal(t, "new-name", updated.Name)
	assert.Equal(t, stringPtr("new description"), updated.Description)
	assert.True(t, updated.UpdatedAt.After(originalUpdatedAt))

	// Test update to duplicate name
	app2, err := repo.Create(ctx, orgID, "another-app", nil)
	require.NoError(t, err)

	_, err = repo.Update(ctx, app2.ID, "new-name", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "application name already exists")

	// Test update non-existent app
	_, err = repo.Update(ctx, uuid.New(), "name", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "application not found")

	// Test nil ID
	_, err = repo.Update(ctx, uuid.Nil, "name", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id cannot be nil")

	// Test empty name
	_, err = repo.Update(ctx, app.ID, "", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name cannot be empty")
}

func TestMemoryRepo_Delete(t *testing.T) {
	repo := NewMemoryRepo()
	ctx := context.Background()

	orgID := uuid.New()

	// Create an app
	app, err := repo.Create(ctx, orgID, "test-app", nil)
	require.NoError(t, err)

	// Test successful delete
	err = repo.Delete(ctx, app.ID)
	assert.NoError(t, err)

	// Verify it's gone
	_, err = repo.GetByID(ctx, app.ID)
	assert.Error(t, err)

	// Test delete non-existent
	err = repo.Delete(ctx, uuid.New())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "application not found")

	// Test nil ID
	err = repo.Delete(ctx, uuid.Nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id cannot be nil")
}

func TestNewRepository(t *testing.T) {
	ctx := context.Background()

	// Test memory backend
	cfg := model.RepositoryConfig{Backend: model.RepositoryMemory}
	repo, err := NewRepository(ctx, cfg)
	assert.NoError(t, err)
	assert.IsType(t, &MemoryRepo{}, repo)

	// Test unsupported backend
	cfg = model.RepositoryConfig{Backend: "unsupported"}
	repo, err = NewRepository(ctx, cfg)
	assert.Error(t, err)
	assert.Nil(t, repo)
	assert.Contains(t, err.Error(), "unsupported applications repository backend")
}

func stringPtr(s string) *string {
	return &s
}
