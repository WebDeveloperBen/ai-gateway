package usage

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	dbdriver "github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	usagerepo "github.com/WebDeveloperBen/ai-gateway/internal/repository/usage"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUsageMetrics_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())
	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	keyID := fixtures.CreateTestAPIKey(t, orgID, appID)

	repo := usagerepo.NewPostgresRepo(pg.Queries)
	svc := NewService(repo)
	router := NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	now := time.Now()
	metric := &model.UsageMetric{
		ID:                uuid.New(),
		OrgID:             orgID,
		AppID:             appID,
		APIKeyID:          keyID,
		Provider:          "openai",
		ModelName:         "gpt-4",
		PromptTokens:      100,
		CompletionTokens:  50,
		TotalTokens:       150,
		RequestSizeBytes:  1000,
		ResponseSizeBytes: 500,
		Timestamp:         now,
	}

	err = repo.Create(context.Background(), metric)
	require.NoError(t, err)

	start := now.Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	end := now.Add(1 * time.Hour).UTC().Format(time.RFC3339)

	ctx := context.WithValue(context.Background(), "org_id", orgID)
	resp := api.GetCtx(ctx, "/api/usage/metrics?app_id="+appID.String()+"&start="+start+"&end="+end)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())

	var result []*UsageMetric
	err = json.Unmarshal(resp.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.NotEmpty(t, result[0].ID)
	assert.Equal(t, "openai", result[0].Provider)
	assert.Equal(t, "gpt-4", result[0].ModelName)
	assert.Equal(t, 100, result[0].PromptTokens)
	assert.Equal(t, 50, result[0].CompletionTokens)
	assert.Equal(t, 150, result[0].TotalTokens)
}

func TestGetUsageMetrics_InvalidAppID(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())
	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, _ := fixtures.CreateTestOrgAndApp(t)

	repo := usagerepo.NewPostgresRepo(pg.Queries)
	svc := NewService(repo)
	router := NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	now := time.Now()
	start := now.Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	end := now.Add(1 * time.Hour).UTC().Format(time.RFC3339)

	ctx := context.WithValue(context.Background(), "org_id", orgID)
	resp := api.GetCtx(ctx, "/api/usage/metrics?app_id=invalid-uuid&start="+start+"&end="+end)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetUsageMetrics_InvalidStartTime(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())
	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	repo := usagerepo.NewPostgresRepo(pg.Queries)
	svc := NewService(repo)
	router := NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	now := time.Now()
	end := now.Add(1 * time.Hour).Format(time.RFC3339)

	ctx := context.WithValue(context.Background(), "org_id", orgID)
	resp := api.GetCtx(ctx, "/api/usage/metrics?app_id="+appID.String()+"&start=invalid-time&end="+end)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetUsageMetrics_InvalidEndTime(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())
	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	repo := usagerepo.NewPostgresRepo(pg.Queries)
	svc := NewService(repo)
	router := NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	now := time.Now()
	start := now.Add(-1 * time.Hour).Format(time.RFC3339)

	ctx := context.WithValue(context.Background(), "org_id", orgID)
	resp := api.GetCtx(ctx, "/api/usage/metrics?app_id="+appID.String()+"&start="+start+"&end=invalid-time")

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetUsageSummaryByApp_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())
	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	keyID := fixtures.CreateTestAPIKey(t, orgID, appID)

	repo := usagerepo.NewPostgresRepo(pg.Queries)
	svc := NewService(repo)
	router := NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	now := time.Now()

	metric1 := &model.UsageMetric{
		ID:               uuid.New(),
		OrgID:            orgID,
		AppID:            appID,
		APIKeyID:         keyID,
		Provider:         "openai",
		ModelName:        "gpt-4",
		PromptTokens:     100,
		CompletionTokens: 50,
		TotalTokens:      150,
		Timestamp:        now,
	}
	metric2 := &model.UsageMetric{
		ID:               uuid.New(),
		OrgID:            orgID,
		AppID:            appID,
		APIKeyID:         keyID,
		Provider:         "openai",
		ModelName:        "gpt-3.5-turbo",
		PromptTokens:     200,
		CompletionTokens: 100,
		TotalTokens:      300,
		Timestamp:        now,
	}

	err = repo.Create(context.Background(), metric1)
	require.NoError(t, err)
	err = repo.Create(context.Background(), metric2)
	require.NoError(t, err)

	start := now.Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	end := now.Add(1 * time.Hour).UTC().Format(time.RFC3339)

	ctx := context.WithValue(context.Background(), "org_id", orgID)
	resp := api.GetCtx(ctx, "/api/usage/summary/app/"+appID.String()+"?start="+start+"&end="+end)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())

	var result *TokenSummary
	err = json.Unmarshal(resp.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, 300, result.TotalPromptTokens)
	assert.Equal(t, 150, result.TotalCompletionTokens)
	assert.Equal(t, 450, result.TotalTokens)
	assert.Equal(t, 2, result.RequestCount)
}

func TestGetUsageSummaryByApp_InvalidAppID(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())
	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, _ := fixtures.CreateTestOrgAndApp(t)

	repo := usagerepo.NewPostgresRepo(pg.Queries)
	svc := NewService(repo)
	router := NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	now := time.Now()
	start := now.Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	end := now.Add(1 * time.Hour).UTC().Format(time.RFC3339)

	ctx := context.WithValue(context.Background(), "org_id", orgID)
	resp := api.GetCtx(ctx, "/api/usage/summary/app/invalid-uuid?start="+start+"&end="+end)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetUsageSummaryByOrg_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())
	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	keyID := fixtures.CreateTestAPIKey(t, orgID, appID)

	repo := usagerepo.NewPostgresRepo(pg.Queries)
	svc := NewService(repo)
	router := NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	now := time.Now()

	metric := &model.UsageMetric{
		ID:               uuid.New(),
		OrgID:            orgID,
		AppID:            appID,
		APIKeyID:         keyID,
		Provider:         "openai",
		ModelName:        "gpt-4",
		PromptTokens:     100,
		CompletionTokens: 50,
		TotalTokens:      150,
		Timestamp:        now,
	}

	err = repo.Create(context.Background(), metric)
	require.NoError(t, err)

	start := now.Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	end := now.Add(1 * time.Hour).UTC().Format(time.RFC3339)

	ctx := context.WithValue(context.Background(), "org_id", orgID)
	resp := api.GetCtx(ctx, "/api/usage/summary/org?start="+start+"&end="+end)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())

	var result *TokenSummary
	err = json.Unmarshal(resp.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, 100, result.TotalPromptTokens)
	assert.Equal(t, 50, result.TotalCompletionTokens)
	assert.Equal(t, 150, result.TotalTokens)
	assert.Equal(t, 1, result.RequestCount)
}

func TestGetUsageSummaryByOrg_MissingOrgID(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())
	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	repo := usagerepo.NewPostgresRepo(pg.Queries)
	svc := NewService(repo)
	router := NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	now := time.Now()
	start := now.Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	end := now.Add(1 * time.Hour).UTC().Format(time.RFC3339)

	ctx := context.Background()
	resp := api.GetCtx(ctx, "/api/usage/summary/org?start="+start+"&end="+end)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestGetUsageByModel_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())
	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	keyID := fixtures.CreateTestAPIKey(t, orgID, appID)

	repo := usagerepo.NewPostgresRepo(pg.Queries)
	svc := NewService(repo)
	router := NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	now := time.Now()

	metric1 := &model.UsageMetric{
		ID:               uuid.New(),
		OrgID:            orgID,
		AppID:            appID,
		APIKeyID:         keyID,
		Provider:         "openai",
		ModelName:        "gpt-4",
		PromptTokens:     100,
		CompletionTokens: 50,
		TotalTokens:      150,
		Timestamp:        now,
	}
	metric2 := &model.UsageMetric{
		ID:               uuid.New(),
		OrgID:            orgID,
		AppID:            appID,
		APIKeyID:         keyID,
		Provider:         "openai",
		ModelName:        "gpt-4",
		PromptTokens:     200,
		CompletionTokens: 100,
		TotalTokens:      300,
		Timestamp:        now,
	}

	err = repo.Create(context.Background(), metric1)
	require.NoError(t, err)
	err = repo.Create(context.Background(), metric2)
	require.NoError(t, err)

	start := now.Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	end := now.Add(1 * time.Hour).UTC().Format(time.RFC3339)

	ctx := context.WithValue(context.Background(), "org_id", orgID)
	resp := api.GetCtx(ctx, "/api/usage/by-model/"+appID.String()+"?start="+start+"&end="+end)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())

	var result []*ModelUsageSummary
	err = json.Unmarshal(resp.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "gpt-4", result[0].ModelName)
	assert.Equal(t, "openai", result[0].Provider)
	assert.Equal(t, 300, result[0].TotalPromptTokens)
	assert.Equal(t, 150, result[0].TotalCompletionTokens)
	assert.Equal(t, 450, result[0].TotalTokens)
	assert.Equal(t, 2, result[0].RequestCount)
}

func TestGetUsageByModel_InvalidAppID(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())
	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, _ := fixtures.CreateTestOrgAndApp(t)

	repo := usagerepo.NewPostgresRepo(pg.Queries)
	svc := NewService(repo)
	router := NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	now := time.Now()
	start := now.Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	end := now.Add(1 * time.Hour).UTC().Format(time.RFC3339)

	ctx := context.WithValue(context.Background(), "org_id", orgID)
	resp := api.GetCtx(ctx, "/api/usage/by-model/invalid-uuid?start="+start+"&end="+end)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetUsageMetrics_WithPagination(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())
	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	keyID := fixtures.CreateTestAPIKey(t, orgID, appID)

	repo := usagerepo.NewPostgresRepo(pg.Queries)
	svc := NewService(repo)
	router := NewRouter(svc)

	api := testkit.SetupAdminTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	now := time.Now()

	for i := 0; i < 5; i++ {
		metric := &model.UsageMetric{
			ID:               uuid.New(),
			OrgID:            orgID,
			AppID:            appID,
			APIKeyID:         keyID,
			Provider:         "openai",
			ModelName:        "gpt-4",
			PromptTokens:     100,
			CompletionTokens: 50,
			TotalTokens:      150,
			Timestamp:        now.Add(time.Duration(i) * time.Minute),
		}
		err = repo.Create(context.Background(), metric)
		require.NoError(t, err)
	}

	start := now.Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	end := now.Add(1 * time.Hour).UTC().Format(time.RFC3339)

	ctx := context.WithValue(context.Background(), "org_id", orgID)
	resp := api.GetCtx(ctx, "/api/usage/metrics?app_id="+appID.String()+"&start="+start+"&end="+end+"&limit=3")

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())

	var result []*UsageMetric
	err = json.Unmarshal(resp.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Len(t, result, 3)
}
