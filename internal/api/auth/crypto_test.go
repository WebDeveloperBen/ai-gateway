package auth_test

import (
	"net/http"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/api/auth"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/require"
)

func TestGenerateState_Indirect(t *testing.T) {
	// Test the generateState function indirectly through the login endpoint
	// Since generateState is not exported, we verify it produces different states
	oidcSvc := auth.NewMockOIDCService()

	router := auth.NewRouter(oidcSvc, nil)
	api := testkit.SetupAuthTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	// Make multiple requests and verify they produce different Location headers
	// (which contain the state parameter)
	var locations []string
	for i := 0; i < 5; i++ {
		resp := api.Get("/auth/login")
		require.Equal(t, http.StatusFound, resp.Code)
		location := resp.Header().Get("Location")
		require.NotEmpty(t, location)
		locations = append(locations, location)
	}

	// Verify all locations are different (different states)
	for i := 0; i < len(locations); i++ {
		for j := i + 1; j < len(locations); j++ {
			require.NotEqual(t, locations[i], locations[j], "login states should be unique")
		}
	}
}
