package health_test

import (
	"net/http"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/api/public/health"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		health.RegisterPublicRoutes(grp)
	})

	resp := api.Get("/api/healthz")

	require.Equal(t, http.StatusOK, resp.Code)
	require.JSONEq(t, `{"message":"All systems online and healthy","status":200}`, resp.Body.String())
}

func TestHealthCheck_MultipleRequests(t *testing.T) {
	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		health.RegisterPublicRoutes(grp)
	})

	for i := 0; i < 10; i++ {
		resp := api.Get("/api/healthz")
		require.Equal(t, http.StatusOK, resp.Code)
		require.JSONEq(t, `{"message":"All systems online and healthy","status":200}`, resp.Body.String())
	}
}
