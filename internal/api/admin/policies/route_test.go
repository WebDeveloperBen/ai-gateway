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

func TestUpdatePolicy_Integration(t *testing.T) {
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

	// Update the policy
	updateBody := policies.UpdatePolicyBody{
		PolicyType: model.PolicyTypeRateLimit,
		Config:     map[string]interface{}{"requests_per_minute": 120},
		Enabled:    false,
	}
	body, _ := json.Marshal(updateBody)

	updateResp := api.Put("/api/policies/"+policyID.String(), "Content-Type: application/json", bytes.NewReader(body))
	require.Equal(t, http.StatusOK, updateResp.Code)

	var updatedPolicy policies.Policy
	err = json.Unmarshal(updateResp.Body.Bytes(), &updatedPolicy)
	require.NoError(t, err)

	require.Equal(t, policyID.String(), updatedPolicy.ID)
	require.Equal(t, model.PolicyTypeRateLimit, updatedPolicy.PolicyType)
	require.Equal(t, 120.0, updatedPolicy.Config["requests_per_minute"])
	require.False(t, updatedPolicy.Enabled)
}

func TestListEnabledPolicies_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	policiesRepo := policiesrepo.NewPostgresRepo(pg.Queries)

	// Create two policies - one enabled, one disabled
	enabledPolicyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 60}`)
	disabledPolicyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 30}`)

	// Disable the second policy
	err = policiesRepo.Disable(context.Background(), disabledPolicyID)
	require.NoError(t, err)

	svc := policies.NewService(policiesRepo)
	router := policies.NewRouter(svc)

	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	// List only enabled policies
	resp := api.Get("/api/policies/enabled?app_id=" + appID.String())
	require.Equal(t, http.StatusOK, resp.Code)

	var enabledPolicies []*policies.Policy
	err = json.Unmarshal(resp.Body.Bytes(), &enabledPolicies)
	require.NoError(t, err)

	require.Len(t, enabledPolicies, 1)
	require.Equal(t, enabledPolicyID.String(), enabledPolicies[0].ID)
	require.True(t, enabledPolicies[0].Enabled)
}

func TestAttachDetachPolicy_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	_, appID2 := fixtures.CreateTestOrgAndApp(t)

	policiesRepo := policiesrepo.NewPostgresRepo(pg.Queries)

	// Create a policy not attached to any app
	policyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 60}`)

	svc := policies.NewService(policiesRepo)
	router := policies.NewRouter(svc)

	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	// Initially, policy should be attached to appID
	resp := api.Get("/api/policies?app_id=" + appID.String())
	require.Equal(t, http.StatusOK, resp.Code)

	var policiesList []*policies.Policy
	err = json.Unmarshal(resp.Body.Bytes(), &policiesList)
	require.NoError(t, err)
	require.Len(t, policiesList, 1)

	// Detach policy from appID
	detachResp := api.Post("/api/policies/" + policyID.String() + "/detach/" + appID.String())
	require.Equal(t, http.StatusNoContent, detachResp.Code)

	// Verify policy is no longer attached to appID
	resp = api.Get("/api/policies?app_id=" + appID.String())
	require.Equal(t, http.StatusOK, resp.Code)

	err = json.Unmarshal(resp.Body.Bytes(), &policiesList)
	require.NoError(t, err)
	require.Len(t, policiesList, 0)

	// Attach policy to appID2
	attachResp := api.Post("/api/policies/" + policyID.String() + "/attach/" + appID2.String())
	require.Equal(t, http.StatusNoContent, attachResp.Code)

	// Verify policy is now attached to appID2
	resp = api.Get("/api/policies?app_id=" + appID2.String())
	require.Equal(t, http.StatusOK, resp.Code)

	err = json.Unmarshal(resp.Body.Bytes(), &policiesList)
	require.NoError(t, err)
	require.Len(t, policiesList, 1)
	require.Equal(t, policyID.String(), policiesList[0].ID)
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

func TestPolicyErrorCases_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	policiesRepo := policiesrepo.NewPostgresRepo(pg.Queries)
	svc := policies.NewService(policiesRepo)
	router := policies.NewRouter(svc)

	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	t.Run("invalid_org_id_create", func(t *testing.T) {
		reqBody := policies.CreatePolicyBody{
			OrgID:      "invalid-uuid",
			PolicyType: model.PolicyTypeRateLimit,
			Config:     map[string]interface{}{"requests_per_minute": 100},
			Enabled:    true,
		}
		body, _ := json.Marshal(reqBody)

		resp := api.Post("/api/policies", "Content-Type: application/json", bytes.NewReader(body))
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("invalid_policy_id_get", func(t *testing.T) {
		resp := api.Get("/api/policies/invalid-uuid")
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("policy_not_found", func(t *testing.T) {
		validUUID := "550e8400-e29b-41d4-a716-446655440000"
		resp := api.Get("/api/policies/" + validUUID)
		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("invalid_app_id_list", func(t *testing.T) {
		resp := api.Get("/api/policies?app_id=invalid-uuid")
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("invalid_policy_id_update", func(t *testing.T) {
		updateBody := policies.UpdatePolicyBody{
			PolicyType: model.PolicyTypeRateLimit,
			Config:     map[string]interface{}{"requests_per_minute": 120},
			Enabled:    false,
		}
		body, _ := json.Marshal(updateBody)

		resp := api.Put("/api/policies/invalid-uuid", "Content-Type: application/json", bytes.NewReader(body))
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("invalid_policy_id_enable", func(t *testing.T) {
		resp := api.Post("/api/policies/invalid-uuid/enable")
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("invalid_policy_id_disable", func(t *testing.T) {
		resp := api.Post("/api/policies/invalid-uuid/disable")
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("invalid_policy_id_attach", func(t *testing.T) {
		resp := api.Post("/api/policies/invalid-uuid/attach/" + appID.String())
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("invalid_app_id_attach", func(t *testing.T) {
		policyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 60}`)
		resp := api.Post("/api/policies/" + policyID.String() + "/attach/invalid-uuid")
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("invalid_policy_id_detach", func(t *testing.T) {
		resp := api.Post("/api/policies/invalid-uuid/detach/" + appID.String())
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("invalid_app_id_detach", func(t *testing.T) {
		policyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 60}`)
		resp := api.Post("/api/policies/" + policyID.String() + "/detach/invalid-uuid")
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("create_policy_missing_required_fields", func(t *testing.T) {
		// Missing policy_type
		reqBody := policies.CreatePolicyBody{
			OrgID:   orgID.String(),
			Config:  map[string]interface{}{"requests_per_minute": 100},
			Enabled: true,
		}
		body, _ := json.Marshal(reqBody)

		resp := api.Post("/api/policies", "Content-Type: application/json", bytes.NewReader(body))
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("update_policy_missing_required_fields", func(t *testing.T) {
		policyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 60}`)

		// Missing policy_type
		updateBody := policies.UpdatePolicyBody{
			Config:  map[string]interface{}{"requests_per_minute": 120},
			Enabled: false,
		}
		body, _ := json.Marshal(updateBody)

		resp := api.Put("/api/policies/"+policyID.String(), "Content-Type: application/json", bytes.NewReader(body))
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestService_GetPoliciesByType(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	policiesRepo := policiesrepo.NewPostgresRepo(pg.Queries)
	svc := policies.NewService(policiesRepo)

	// Create multiple policies of different types
	rateLimitPolicyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeRateLimit, `{"requests_per_minute": 60}`)
	modelAllowlistPolicyID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeModelAllowlist, `{"allowed_models": ["gpt-4", "claude"]}`)

	// Test GetPoliciesByType for rate limit policies
	rateLimitPolicies, err := svc.GetPoliciesByType(context.Background(), appID, model.PolicyTypeRateLimit)
	require.NoError(t, err)
	require.Len(t, rateLimitPolicies, 1)
	require.Equal(t, rateLimitPolicyID.String(), rateLimitPolicies[0].ID)
	require.Equal(t, model.PolicyTypeRateLimit, rateLimitPolicies[0].PolicyType)

	// Test GetPoliciesByType for model allowlist policies
	modelAllowlistPolicies, err := svc.GetPoliciesByType(context.Background(), appID, model.PolicyTypeModelAllowlist)
	require.NoError(t, err)
	require.Len(t, modelAllowlistPolicies, 1)
	require.Equal(t, modelAllowlistPolicyID.String(), modelAllowlistPolicies[0].ID)
	require.Equal(t, model.PolicyTypeModelAllowlist, modelAllowlistPolicies[0].PolicyType)

	// Test GetPoliciesByType for non-existent type
	emptyPolicies, err := svc.GetPoliciesByType(context.Background(), appID, model.PolicyType("non-existent"))
	require.NoError(t, err)
	require.Len(t, emptyPolicies, 0)
}

func TestService_GetAppsForPolicy(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID1 := fixtures.CreateTestOrgAndApp(t)
	_, appID2 := fixtures.CreateTestOrgAndApp(t)

	policiesRepo := policiesrepo.NewPostgresRepo(pg.Queries)
	svc := policies.NewService(policiesRepo)

	// Create a policy and attach it to both apps
	policyID := fixtures.CreateTestPolicy(t, orgID, appID1, model.PolicyTypeRateLimit, `{"requests_per_minute": 60}`)

	// Attach to second app
	err = svc.AttachPolicyToApp(context.Background(), policyID, appID2)
	require.NoError(t, err)

	// Test GetAppsForPolicy
	apps, err := svc.GetAppsForPolicy(context.Background(), policyID)
	require.NoError(t, err)
	require.Len(t, apps, 2)

	// Verify both apps are returned (order may vary)
	appIDs := make(map[string]bool)
	for _, app := range apps {
		appIDs[app.ID.String()] = true
	}
	require.True(t, appIDs[appID1.String()])
	require.True(t, appIDs[appID2.String()])
}
