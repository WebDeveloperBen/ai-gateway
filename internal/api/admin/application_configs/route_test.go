package application_configs_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/api/admin/application_configs"
	dbdriver "github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	configrepo "github.com/WebDeveloperBen/ai-gateway/internal/repository/application_configs"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateApplicationConfig_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	reqBody := application_configs.CreateApplicationConfigBody{
		AppID:       appID.String(),
		Environment: "production",
		Config: map[string]any{
			"key1": "value1",
			"key2": 123,
		},
	}
	body, _ := json.Marshal(reqBody)

	ctx := context.WithValue(context.Background(), "org_id", orgID)
	resp := api.PostCtx(ctx, "/api/application-configs", "Content-Type: application/json", bytes.NewReader(body))

	require.Equal(t, http.StatusCreated, resp.Code, resp.Body.String())

	var responseBody application_configs.ApplicationConfig
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.NotEmpty(t, responseBody.ID)
	require.Equal(t, appID.String(), responseBody.AppID)
	require.Equal(t, orgID.String(), responseBody.OrgID)
	require.Equal(t, "production", responseBody.Environment)
	require.Equal(t, "value1", responseBody.Config["key1"])
	require.Equal(t, float64(123), responseBody.Config["key2"])
}

func TestCreateApplicationConfig_NoOrgIDInContext(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	reqBody := application_configs.CreateApplicationConfigBody{
		AppID:       uuid.New().String(),
		Environment: "production",
		Config:      map[string]any{},
	}
	body, _ := json.Marshal(reqBody)

	resp := api.Post("/api/application-configs", "Content-Type: application/json", bytes.NewReader(body))

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestCreateApplicationConfig_InvalidAppID(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, _ := fixtures.CreateTestOrgAndApp(t)

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	reqBody := application_configs.CreateApplicationConfigBody{
		AppID:       "invalid-uuid",
		Environment: "production",
		Config:      map[string]any{},
	}
	body, _ := json.Marshal(reqBody)

	ctx := context.WithValue(context.Background(), "org_id", orgID)
	resp := api.PostCtx(ctx, "/api/application-configs", "Content-Type: application/json", bytes.NewReader(body))

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCreateApplicationConfig_MissingEnvironment(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	reqBody := application_configs.CreateApplicationConfigBody{
		AppID:       appID.String(),
		Environment: "",
		Config:      map[string]any{},
	}
	body, _ := json.Marshal(reqBody)

	ctx := context.WithValue(context.Background(), "org_id", orgID)
	resp := api.PostCtx(ctx, "/api/application-configs", "Content-Type: application/json", bytes.NewReader(body))

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestListApplicationConfigs_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	_, err = configRepo.Create(context.Background(), appID, orgID, "development", map[string]any{"env": "dev"})
	require.NoError(t, err)
	_, err = configRepo.Create(context.Background(), appID, orgID, "production", map[string]any{"env": "prod"})
	require.NoError(t, err)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Get("/api/applications/" + appID.String() + "/configs")

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())

	var configs []*application_configs.ApplicationConfig
	err = json.Unmarshal(resp.Body.Bytes(), &configs)
	require.NoError(t, err)

	require.Len(t, configs, 2)
}

func TestListApplicationConfigs_InvalidAppID(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Get("/api/applications/invalid-uuid/configs")

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetApplicationConfig_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	cfg, err := configRepo.Create(context.Background(), appID, orgID, "staging", map[string]any{"key": "value"})
	require.NoError(t, err)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Get("/api/application-configs/" + cfg.ID.String())

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())

	var responseBody application_configs.ApplicationConfig
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.Equal(t, cfg.ID.String(), responseBody.ID)
	require.Equal(t, "staging", responseBody.Environment)
}

func TestGetApplicationConfig_InvalidID(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Get("/api/application-configs/invalid-uuid")

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetApplicationConfig_NotFound(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	nonExistentID := uuid.New()
	resp := api.Get("/api/application-configs/" + nonExistentID.String())

	require.Equal(t, http.StatusNotFound, resp.Code)
}

func TestGetApplicationConfigByEnv_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	_, err = configRepo.Create(context.Background(), appID, orgID, "test", map[string]any{"test_key": "test_value"})
	require.NoError(t, err)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Get("/api/applications/" + appID.String() + "/configs/test")

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())

	var responseBody application_configs.ApplicationConfig
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.Equal(t, "test", responseBody.Environment)
	require.Equal(t, "test_value", responseBody.Config["test_key"])
}

func TestGetApplicationConfigByEnv_InvalidAppID(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Get("/api/applications/invalid-uuid/configs/test")

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetApplicationConfigByEnv_NotFound(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	_, appID := fixtures.CreateTestOrgAndApp(t)

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Get("/api/applications/" + appID.String() + "/configs/nonexistent")

	require.Equal(t, http.StatusNotFound, resp.Code)
}

func TestUpdateApplicationConfig_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	cfg, err := configRepo.Create(context.Background(), appID, orgID, "update-test", map[string]any{"old": "value"})
	require.NoError(t, err)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	reqBody := application_configs.UpdateApplicationConfigBody{
		Config: map[string]any{"new": "updated_value"},
	}
	body, _ := json.Marshal(reqBody)

	resp := api.Put("/api/application-configs/"+cfg.ID.String(), "Content-Type: application/json", bytes.NewReader(body))

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())

	var responseBody application_configs.ApplicationConfig
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.Equal(t, cfg.ID.String(), responseBody.ID)
	require.Equal(t, "updated_value", responseBody.Config["new"])
}

func TestUpdateApplicationConfig_InvalidID(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	reqBody := application_configs.UpdateApplicationConfigBody{
		Config: map[string]any{"key": "value"},
	}
	body, _ := json.Marshal(reqBody)

	resp := api.Put("/api/application-configs/invalid-uuid", "Content-Type: application/json", bytes.NewReader(body))

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestDeleteApplicationConfig_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	cfg, err := configRepo.Create(context.Background(), appID, orgID, "delete-test", map[string]any{"key": "value"})
	require.NoError(t, err)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Delete("/api/application-configs/" + cfg.ID.String())

	require.Equal(t, http.StatusNoContent, resp.Code)

	_, err = configRepo.GetByID(context.Background(), cfg.ID)
	require.Error(t, err)
}

func TestDeleteApplicationConfig_InvalidID(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Delete("/api/application-configs/invalid-uuid")

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestDeleteApplicationConfig_NotFound(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	configRepo := configrepo.NewPostgresRepo(pg.Queries)

	svc := application_configs.NewService(configRepo)
	router := application_configs.NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	nonExistentID := uuid.New()
	resp := api.Delete("/api/application-configs/" + nonExistentID.String())

	require.Equal(t, http.StatusNoContent, resp.Code)
}
