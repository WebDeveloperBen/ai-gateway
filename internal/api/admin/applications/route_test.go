package applications_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/api/admin/applications"
	"github.com/WebDeveloperBen/ai-gateway/internal/api/middleware"
	dbdriver "github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	apprepo "github.com/WebDeveloperBen/ai-gateway/internal/repository/applications"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateApplication_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, _ := fixtures.CreateTestOrgAndApp(t)

	appRepo, err := apprepo.NewRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	svc := applications.NewService(appRepo)
	router := applications.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	desc := "Test Application Description"
	reqBody := applications.CreateApplicationBody{
		Name:        "Test Application",
		Description: &desc,
	}
	body, _ := json.Marshal(reqBody)

	ctx := context.WithValue(context.Background(), middleware.ScopedTokenKey, model.ScopedToken{
		OrgID: orgID.String(),
	})
	resp := api.PostCtx(ctx, "/api/applications", "Content-Type: application/json", bytes.NewReader(body))

	require.Equal(t, http.StatusCreated, resp.Code, resp.Body.String())

	var responseBody applications.Application
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.NotEmpty(t, responseBody.ID)
	require.Equal(t, orgID.String(), responseBody.OrgID)
	require.Equal(t, "Test Application", responseBody.Name)
	require.NotNil(t, responseBody.Description)
	require.Equal(t, desc, *responseBody.Description)
	require.NotZero(t, responseBody.CreatedAt)
}

func TestCreateApplication_NoOrgIDInContext(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	appRepo, err := apprepo.NewRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	svc := applications.NewService(appRepo)
	router := applications.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	reqBody := applications.CreateApplicationBody{
		Name: "Test Application",
	}
	body, _ := json.Marshal(reqBody)

	resp := api.Post("/api/applications", "Content-Type: application/json", bytes.NewReader(body))

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestCreateApplication_ValidationError(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, _ := fixtures.CreateTestOrgAndApp(t)

	appRepo, err := apprepo.NewRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	svc := applications.NewService(appRepo)
	router := applications.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	reqBody := applications.CreateApplicationBody{
		Name: "",
	}
	body, _ := json.Marshal(reqBody)

	ctx := context.WithValue(context.Background(), middleware.ScopedTokenKey, model.ScopedToken{
		OrgID: orgID.String(),
	})
	resp := api.PostCtx(ctx, "/api/applications", "Content-Type: application/json", bytes.NewReader(body))

	require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}

func TestListApplications_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, _ := fixtures.CreateTestOrgAndApp(t)

	appRepo, err := apprepo.NewRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	_, err = appRepo.Create(context.Background(), orgID, "App 1", nil)
	require.NoError(t, err)
	_, err = appRepo.Create(context.Background(), orgID, "App 2", nil)
	require.NoError(t, err)

	svc := applications.NewService(appRepo)
	router := applications.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	ctx := context.WithValue(context.Background(), middleware.ScopedTokenKey, model.ScopedToken{
		OrgID: orgID.String(),
	})
	resp := api.GetCtx(ctx, "/api/applications")

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())

	var apps []*applications.Application
	err = json.Unmarshal(resp.Body.Bytes(), &apps)
	require.NoError(t, err)

	require.Len(t, apps, 3)
}

func TestListApplications_NoOrgIDInContext(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	appRepo, err := apprepo.NewRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	svc := applications.NewService(appRepo)
	router := applications.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Get("/api/applications")

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestGetApplication_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	appRepo, err := apprepo.NewRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	svc := applications.NewService(appRepo)
	router := applications.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Get("/api/applications/" + appID.String())

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())

	var app applications.Application
	err = json.Unmarshal(resp.Body.Bytes(), &app)
	require.NoError(t, err)

	require.Equal(t, appID.String(), app.ID)
	require.Equal(t, orgID.String(), app.OrgID)
}

func TestGetApplication_InvalidID(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	appRepo, err := apprepo.NewRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	svc := applications.NewService(appRepo)
	router := applications.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Get("/api/applications/invalid-uuid")

	require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}

func TestGetApplication_NotFound(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	appRepo, err := apprepo.NewRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	svc := applications.NewService(appRepo)
	router := applications.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	nonExistentID := uuid.New()
	resp := api.Get("/api/applications/" + nonExistentID.String())

	require.Equal(t, http.StatusNotFound, resp.Code)
}

func TestUpdateApplication_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	appRepo, err := apprepo.NewRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	svc := applications.NewService(appRepo)
	router := applications.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	newDesc := "Updated Description"
	reqBody := applications.UpdateApplicationBody{
		Name:        "Updated Name",
		Description: &newDesc,
	}
	body, _ := json.Marshal(reqBody)

	ctx := context.WithValue(context.Background(), middleware.ScopedTokenKey, model.ScopedToken{
		OrgID: orgID.String(),
	})
	resp := api.PutCtx(ctx, "/api/applications/"+appID.String(), "Content-Type: application/json", bytes.NewReader(body))

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())

	var app applications.Application
	err = json.Unmarshal(resp.Body.Bytes(), &app)
	require.NoError(t, err)

	require.Equal(t, appID.String(), app.ID)
	require.Equal(t, orgID.String(), app.OrgID)
	require.Equal(t, "Updated Name", app.Name)
	require.NotNil(t, app.Description)
	require.Equal(t, newDesc, *app.Description)
}

func TestUpdateApplication_InvalidID(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	appRepo, err := apprepo.NewRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	svc := applications.NewService(appRepo)
	router := applications.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	reqBody := applications.UpdateApplicationBody{
		Name: "Updated Name",
	}
	body, _ := json.Marshal(reqBody)

	resp := api.Put("/api/applications/invalid-uuid", "Content-Type: application/json", bytes.NewReader(body))

	require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}

func TestUpdateApplication_ValidationError(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	appRepo, err := apprepo.NewRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	svc := applications.NewService(appRepo)
	router := applications.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	reqBody := applications.UpdateApplicationBody{
		Name: "",
	}
	body, _ := json.Marshal(reqBody)

	ctx := context.WithValue(context.Background(), middleware.ScopedTokenKey, model.ScopedToken{
		OrgID: orgID.String(),
	})
	resp := api.PutCtx(ctx, "/api/applications/"+appID.String(), "Content-Type: application/json", bytes.NewReader(body))

	require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}

func TestDeleteApplication_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, _ := fixtures.CreateTestOrgAndApp(t)

	appRepo, err := apprepo.NewRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	app, err := appRepo.Create(context.Background(), orgID, "To Delete", nil)
	require.NoError(t, err)

	svc := applications.NewService(appRepo)
	router := applications.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Delete("/api/applications/" + app.ID.String())

	require.Equal(t, http.StatusNoContent, resp.Code)

	_, err = appRepo.GetByID(context.Background(), app.ID)
	require.Error(t, err)
}

func TestDeleteApplication_InvalidID(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	appRepo, err := apprepo.NewRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	svc := applications.NewService(appRepo)
	router := applications.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Delete("/api/applications/invalid-uuid")

	require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}

func TestDeleteApplication_NotFound(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	appRepo, err := apprepo.NewRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	svc := applications.NewService(appRepo)
	router := applications.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	nonExistentID := uuid.New()
	resp := api.Delete("/api/applications/" + nonExistentID.String())

	require.Equal(t, http.StatusNoContent, resp.Code)
}
