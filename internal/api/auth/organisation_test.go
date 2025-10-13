package auth_test

import (
	"context"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/api/auth"
	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	dbdriver "github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/repository/organisations"
	"github.com/WebDeveloperBen/ai-gateway/internal/repository/users"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrganisationService_EnsureUserAndOrg(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())
	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)

	t.Run("existing user", func(t *testing.T) {
		orgID, _ := fixtures.CreateTestOrgAndApp(t)
		// Create a user manually for testing
		sub := "test-sub-123"
		name := "Test User"
		user, err := pg.Queries.CreateUser(context.Background(), db.CreateUserParams{
			OrgID: orgID,
			Sub:   &sub,
			Email: "existing@example.com",
			Name:  &name,
		})
		require.NoError(t, err)

		orgRepo := organisations.NewPostgresRepo(pg.Queries)
		userRepo := users.NewPostgresRepo(pg.Queries)
		svc := auth.NewOrganisationService(orgRepo, userRepo)

		scoped := model.ScopedToken{
			RegisteredClaims: jwt.RegisteredClaims{Subject: *user.Sub},
			Email:            user.Email,
			Name:             *user.Name,
		}

		resultUser, resultOrg, err := svc.EnsureUserAndOrg(context.Background(), scoped)
		require.NoError(t, err)
		require.NotNil(t, resultUser)
		require.NotNil(t, resultOrg)
		assert.Equal(t, user.ID.String(), resultUser.ID)
		assert.Equal(t, orgID.String(), resultOrg.ID)
	})

	t.Run("new user creates org and user", func(t *testing.T) {
		orgRepo := organisations.NewPostgresRepo(pg.Queries)
		userRepo := users.NewPostgresRepo(pg.Queries)
		svc := auth.NewOrganisationService(orgRepo, userRepo)

		scoped := model.ScopedToken{
			RegisteredClaims: jwt.RegisteredClaims{Subject: "new-user-123"},
			Email:            "newuser@example.com",
			Name:             "New User",
		}

		resultUser, resultOrg, err := svc.EnsureUserAndOrg(context.Background(), scoped)
		require.NoError(t, err)
		require.NotNil(t, resultUser)
		require.NotNil(t, resultOrg)

		// Verify user was created
		assert.Equal(t, "newuser@example.com", resultUser.Email)
		assert.Equal(t, "new-user-123", resultUser.Sub)
		assert.Equal(t, "New User", resultUser.Name)

		// Verify org was created
		assert.Contains(t, resultOrg.Name, "New User's Home")
	})

	t.Run("invalid claims", func(t *testing.T) {
		orgRepo := organisations.NewPostgresRepo(pg.Queries)
		userRepo := users.NewPostgresRepo(pg.Queries)
		svc := auth.NewOrganisationService(orgRepo, userRepo)

		// Missing subject
		scoped := model.ScopedToken{
			Email: "test@example.com",
			Name:  "Test User",
		}

		_, _, err := svc.EnsureUserAndOrg(context.Background(), scoped)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid claims")
	})
}

func TestOrganisationService_EnsureRole(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())
	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, _ := fixtures.CreateTestOrgAndApp(t)

	orgRepo := organisations.NewPostgresRepo(pg.Queries)
	userRepo := users.NewPostgresRepo(pg.Queries)
	svc := auth.NewOrganisationService(orgRepo, userRepo)

	t.Run("create new role", func(t *testing.T) {
		role, err := svc.EnsureRole(context.Background(), orgID.String(), "test-role", "Test role description")
		require.NoError(t, err)
		require.NotNil(t, role)
		assert.Equal(t, "test-role", role.Name)
		assert.Equal(t, "Test role description", role.Description)
	})

	t.Run("existing role", func(t *testing.T) {
		// Create role first
		role1, err := svc.EnsureRole(context.Background(), orgID.String(), "existing-role", "Existing role")
		require.NoError(t, err)

		// Try to ensure the same role again
		role2, err := svc.EnsureRole(context.Background(), orgID.String(), "existing-role", "Existing role")
		require.NoError(t, err)
		assert.Equal(t, role1.ID, role2.ID)
	})

	t.Run("invalid org ID", func(t *testing.T) {
		_, err := svc.EnsureRole(context.Background(), "invalid-uuid", "role", "desc")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid org uuid")
	})
}

func TestNewOrganisationService(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())
	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	orgRepo := organisations.NewPostgresRepo(pg.Queries)
	userRepo := users.NewPostgresRepo(pg.Queries)

	svc := auth.NewOrganisationService(orgRepo, userRepo)
	require.NotNil(t, svc)

	// Verify it's the right type
	_, ok := svc.(*auth.OrganisationService)
	assert.True(t, ok)
}
