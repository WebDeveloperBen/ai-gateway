package users

import (
	"context"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	dbdriver "github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *dbdriver.Postgres {
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

	return pg
}

func TestNewPostgresRepo(t *testing.T) {
	pg := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	assert.NotNil(t, repo)
	assert.IsType(t, &postgresRepo{}, repo)
}

func TestPostgresRepo_FindBySubOrEmail(t *testing.T) {
	pg := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)

	// Create an organisation first
	org, err := pg.Queries.CreateOrg(context.Background(), "Test Org")
	require.NoError(t, err)

	// Create a user first
	sub := "test-sub"
	name := "Test User"
	user, err := pg.Queries.CreateUser(context.Background(), db.CreateUserParams{
		OrgID: org.ID,
		Sub:   &sub,
		Email: "test@example.com",
		Name:  &name,
	})
	require.NoError(t, err)

	t.Run("find by sub", func(t *testing.T) {
		found, err := repo.FindBySubOrEmail(context.Background(), "test-sub", "other@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, user.ID.String(), found.ID)
		assert.Equal(t, "test-sub", found.Sub)
	})

	t.Run("find by email", func(t *testing.T) {
		found, err := repo.FindBySubOrEmail(context.Background(), "other-sub", "test@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, user.ID.String(), found.ID)
		assert.Equal(t, "test@example.com", found.Email)
	})

	t.Run("not found", func(t *testing.T) {
		found, err := repo.FindBySubOrEmail(context.Background(), "non-existing-sub", "non-existing@example.com")
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestPostgresRepo_Create(t *testing.T) {
	pg := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)

	// Create an organisation first
	org, err := pg.Queries.CreateOrg(context.Background(), "Test Org")
	require.NoError(t, err)

	sub := "create-test-sub"
	name := "Create Test User"

	created, err := repo.Create(context.Background(), db.CreateUserParams{
		OrgID: org.ID,
		Sub:   &sub,
		Email: "create@example.com",
		Name:  &name,
	})

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, org.ID.String(), created.OrgID)
	assert.Equal(t, "create-test-sub", created.Sub)
	assert.Equal(t, "create@example.com", created.Email)
	assert.Equal(t, "Create Test User", created.Name)
	assert.NotEmpty(t, created.ID)
	assert.NotZero(t, created.CreatedAt)
	assert.NotZero(t, created.UpdatedAt)
}

func TestPostgresRepo_AssignRole(t *testing.T) {
	pg := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)

	// Create an organisation first
	org, err := pg.Queries.CreateOrg(context.Background(), "Test Org")
	require.NoError(t, err)

	// Create user
	sub := "assign-test-sub"
	name := "Assign Test User"
	user, err := pg.Queries.CreateUser(context.Background(), db.CreateUserParams{
		OrgID: org.ID,
		Sub:   &sub,
		Email: "assign@example.com",
		Name:  &name,
	})
	require.NoError(t, err)

	// Create role
	role, err := pg.Queries.CreateRole(context.Background(), db.CreateRoleParams{
		Name:        "assign-test-role",
		Description: &name,
	})
	require.NoError(t, err)

	t.Run("assign role to user", func(t *testing.T) {
		err := repo.AssignRole(context.Background(), user.ID.String(), role.ID.String(), org.ID)
		assert.NoError(t, err)
	})

	t.Run("assign same role again", func(t *testing.T) {
		// Should succeed due to ON CONFLICT DO NOTHING
		err := repo.AssignRole(context.Background(), user.ID.String(), role.ID.String(), org.ID)
		assert.NoError(t, err)
	})

	t.Run("invalid user uuid", func(t *testing.T) {
		err := repo.AssignRole(context.Background(), "invalid-uuid", role.ID.String(), org.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid uuid")
	})

	t.Run("invalid role uuid", func(t *testing.T) {
		err := repo.AssignRole(context.Background(), user.ID.String(), "invalid-uuid", org.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid uuid")
	})

	t.Run("nil org uuid", func(t *testing.T) {
		err := repo.AssignRole(context.Background(), user.ID.String(), role.ID.String(), uuid.Nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid uuid")
	})
}
