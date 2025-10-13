package organisations

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

func TestPostgresRepo_Create(t *testing.T) {
	pg := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)

	t.Run("create organisation", func(t *testing.T) {
		org, err := repo.Create(context.Background(), "Test Organisation")
		assert.NoError(t, err)
		assert.NotNil(t, org)
		assert.Equal(t, "Test Organisation", org.Name)
		assert.NotEmpty(t, org.ID)
		assert.NotZero(t, org.CreatedAt)
		assert.NotZero(t, org.UpdatedAt)
	})
}

func TestPostgresRepo_FindByID(t *testing.T) {
	pg := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)

	// Create an organisation first
	createdOrg, err := repo.Create(context.Background(), "Test Organisation")
	require.NoError(t, err)

	t.Run("find existing organisation", func(t *testing.T) {
		foundOrg, err := repo.FindByID(context.Background(), createdOrg.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundOrg)
		assert.Equal(t, createdOrg.ID, foundOrg.ID)
		assert.Equal(t, createdOrg.Name, foundOrg.Name)
	})

	t.Run("find non-existing organisation", func(t *testing.T) {
		foundOrg, err := repo.FindByID(context.Background(), uuid.New().String())
		assert.Error(t, err)
		assert.Nil(t, foundOrg)
		assert.Contains(t, err.Error(), "organisation not found")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		foundOrg, err := repo.FindByID(context.Background(), "invalid-uuid")
		assert.Error(t, err)
		assert.Nil(t, foundOrg)
		assert.Contains(t, err.Error(), "invalid uuid")
	})
}

func TestPostgresRepo_FindRoleByName(t *testing.T) {
	pg := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)

	// Create a role first
	role, err := repo.CreateRole(context.Background(), "find-test-role", "Test role description")
	require.NoError(t, err)

	t.Run("find existing role", func(t *testing.T) {
		foundRole, err := repo.FindRoleByName(context.Background(), "find-test-role")
		assert.NoError(t, err)
		assert.NotNil(t, foundRole)
		assert.Equal(t, role.ID, foundRole.ID)
		assert.Equal(t, role.Name, foundRole.Name)
		assert.Equal(t, role.Description, foundRole.Description)
	})

	t.Run("find non-existing role", func(t *testing.T) {
		foundRole, err := repo.FindRoleByName(context.Background(), "non-existing-role")
		assert.Error(t, err)
		assert.Nil(t, foundRole)
		assert.Contains(t, err.Error(), "role not found")
	})
}

func TestPostgresRepo_CreateRole(t *testing.T) {
	pg := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)

	t.Run("create role", func(t *testing.T) {
		role, err := repo.CreateRole(context.Background(), "unique-test-role", "Test role description")
		assert.NoError(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, "unique-test-role", role.Name)
		assert.Equal(t, "Test role description", role.Description)
		assert.NotEmpty(t, role.ID)
		assert.NotZero(t, role.CreatedAt)
	})

	t.Run("create duplicate role", func(t *testing.T) {
		// Create first role
		_, err := repo.CreateRole(context.Background(), "duplicate-test-role", "Description")
		require.NoError(t, err)

		// Try to create duplicate
		_, err = repo.CreateRole(context.Background(), "duplicate-test-role", "Description")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "role already exists")
	})
}

func TestPostgresRepo_AssignRole(t *testing.T) {
	pg := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)

	// Create org and role
	org, err := repo.Create(context.Background(), "Test Org")
	require.NoError(t, err)

	role, err := repo.CreateRole(context.Background(), "assign-test-role", "Test role")
	require.NoError(t, err)

	t.Run("assign role to organisation", func(t *testing.T) {
		err := repo.AssignRole(context.Background(), org.ID, role.ID)
		assert.NoError(t, err)
	})

	t.Run("assign same role again", func(t *testing.T) {
		// Should succeed due to ON CONFLICT DO NOTHING
		err := repo.AssignRole(context.Background(), org.ID, role.ID)
		assert.NoError(t, err)
	})

	t.Run("assign to non-existing organisation", func(t *testing.T) {
		err := repo.AssignRole(context.Background(), uuid.New().String(), role.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reference")
	})

	t.Run("assign non-existing role", func(t *testing.T) {
		err := repo.AssignRole(context.Background(), org.ID, uuid.New().String())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reference")
	})

	t.Run("invalid org uuid", func(t *testing.T) {
		err := repo.AssignRole(context.Background(), "invalid-uuid", role.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid uuid")
	})

	t.Run("invalid role uuid", func(t *testing.T) {
		err := repo.AssignRole(context.Background(), org.ID, "invalid-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid uuid")
	})
}

func TestPostgresRepo_EnsureMembership(t *testing.T) {
	pg := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)

	// Create org
	org, err := repo.Create(context.Background(), "Test Org")
	require.NoError(t, err)

	orgUUID, err := uuid.Parse(org.ID)
	require.NoError(t, err)

	// Create a user first
	sub := "test-user-sub"
	name := "Test User"
	user, err := pg.Queries.CreateUser(context.Background(), db.CreateUserParams{
		OrgID: orgUUID,
		Sub:   &sub,
		Email: "test@example.com",
		Name:  &name,
	})
	require.NoError(t, err)

	t.Run("ensure membership", func(t *testing.T) {
		err := repo.EnsureMembership(context.Background(), orgUUID, user.ID)
		assert.NoError(t, err)
	})

	t.Run("ensure membership again", func(t *testing.T) {
		// Should succeed due to ON CONFLICT DO NOTHING
		err := repo.EnsureMembership(context.Background(), orgUUID, user.ID)
		assert.NoError(t, err)
	})
}
