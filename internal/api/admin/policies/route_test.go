package policies_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/api/admin/policies"
	dbdriver "github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	policiesrepo "github.com/WebDeveloperBen/ai-gateway/internal/repository/policies"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/require"
)

func TestCreatePolicy_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, _ := fixtures.CreateTestOrgAndApp(t)

	policiesRepo := policiesrepo.NewPostgresRepo(pg.Queries)

	svc := policies.NewService(policiesRepo)
	router := policies.NewRouter(svc)

	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	reqBody := policies.CreatePolicyBody{
		OrgID:      orgID.String(),
		PolicyType: model.PolicyTypeRateLimit,
		Config:     map[string]interface{}{"requests_per_minute": 100},
		Enabled:    true,
		// AppID is optional - can be attached later
	}
	body, _ := json.Marshal(reqBody)

	resp := api.Post("/api/policies", "Content-Type: application/json", bytes.NewReader(body))

	require.Equal(t, http.StatusCreated, resp.Code)

	var responseBody policies.Policy
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.Equal(t, orgID.String(), responseBody.OrgID)
	require.Equal(t, model.PolicyTypeRateLimit, responseBody.PolicyType)
	require.True(t, responseBody.Enabled)
}

func TestListPolicies_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	// Create a test policy and attach it to the app (both done by CreateTestPolicy)
	policyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 60}`)

	policiesRepo := policiesrepo.NewPostgresRepo(pg.Queries)
	svc := policies.NewService(policiesRepo)
	router := policies.NewRouter(svc)

	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Get("/api/policies?app_id=" + appID.String())

	require.Equal(t, http.StatusOK, resp.Code)

	var policiesList []*policies.Policy
	err = json.Unmarshal(resp.Body.Bytes(), &policiesList)
	require.NoError(t, err)

	require.Len(t, policiesList, 1)
	require.Equal(t, policyID.String(), policiesList[0].ID)
	require.Equal(t, model.PolicyTypeRateLimit, policiesList[0].PolicyType)
	require.True(t, policiesList[0].Enabled)
}

func TestGetPolicy_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	policiesRepo := policiesrepo.NewPostgresRepo(pg.Queries)

	// Create a test policy directly in the database
	policyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 60}`)

	svc := policies.NewService(policiesRepo)
	router := policies.NewRouter(svc)

	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Get("/api/policies/" + policyID.String())

	require.Equal(t, http.StatusOK, resp.Code)

	var policy policies.Policy
	err = json.Unmarshal(resp.Body.Bytes(), &policy)
	require.NoError(t, err)

	require.Equal(t, policyID.String(), policy.ID)
	require.Equal(t, model.PolicyTypeRateLimit, policy.PolicyType)
	require.True(t, policy.Enabled)
}

func TestEnableDisablePolicy_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	policiesRepo := policiesrepo.NewPostgresRepo(pg.Queries)

	// Create a test policy directly in the database (initially disabled)
	policyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 60}`)

	// Disable the policy first
	err = policiesRepo.Disable(context.Background(), policyID)
	require.NoError(t, err)

	svc := policies.NewService(policiesRepo)
	router := policies.NewRouter(svc)

	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	// Enable the policy
	enableResp := api.Post("/api/policies/" + policyID.String() + "/enable")
	require.Equal(t, http.StatusNoContent, enableResp.Code)

	// Verify it's enabled
	getResp := api.Get("/api/policies/" + policyID.String())
	require.Equal(t, http.StatusOK, getResp.Code)

	var policyResp policies.Policy
	err = json.Unmarshal(getResp.Body.Bytes(), &policyResp)
	require.NoError(t, err)
	require.True(t, policyResp.Enabled)

	// Disable the policy
	disableResp := api.Post("/api/policies/" + policyID.String() + "/disable")
	require.Equal(t, http.StatusNoContent, disableResp.Code)

	// Verify it's disabled
	getResp2 := api.Get("/api/policies/" + policyID.String())
	require.Equal(t, http.StatusOK, getResp2.Code)

	var policyResp2 policies.Policy
	err = json.Unmarshal(getResp2.Body.Bytes(), &policyResp2)
	require.NoError(t, err)
	require.False(t, policyResp2.Enabled)
}

func TestDeletePolicy_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	policiesRepo := policiesrepo.NewPostgresRepo(pg.Queries)

	// Create a test policy directly in the database
	policyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 60}`)

	svc := policies.NewService(policiesRepo)
	router := policies.NewRouter(svc)

	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	// Delete the policy
	deleteResp := api.Delete("/api/policies/" + policyID.String())
	require.Equal(t, http.StatusNoContent, deleteResp.Code)

	// Verify it's gone
	getResp := api.Get("/api/policies/" + policyID.String())
	require.Equal(t, http.StatusNotFound, getResp.Code)
}
