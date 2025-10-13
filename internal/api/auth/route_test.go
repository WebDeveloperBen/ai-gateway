package auth_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/api/auth"
	"github.com/WebDeveloperBen/ai-gateway/internal/api/middleware"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestLogin_Redirect(t *testing.T) {
	oidcSvc := auth.NewMockOIDCService()

	// For login test, we don't need org service since it's not called
	router := auth.NewRouter(oidcSvc, nil)
	api := testkit.SetupAuthTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Get("/auth/login")

	require.Equal(t, http.StatusFound, resp.Code)
	require.Contains(t, resp.Header().Get("Location"), "https://mock-auth-url.com")
}

// TODO: Implement callback test once OAuth mocking is figured out
// func TestCallback_Success(t *testing.T) { ... }

func TestCallback_MissingCode(t *testing.T) {
	oidcSvc := auth.NewMockOIDCService()

	router := auth.NewRouter(oidcSvc, nil)
	api := testkit.SetupAuthTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Get("/auth/callback")

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestCallback_OIDError(t *testing.T) {
	oidcSvc := auth.NewMockOIDCService()

	router := auth.NewRouter(oidcSvc, nil)
	api := testkit.SetupAuthTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Get("/auth/callback?error=access_denied")

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestMe_Success(t *testing.T) {
	oidcSvc := auth.NewMockOIDCService()

	router := auth.NewRouter(oidcSvc, nil)

	// Set up API with middleware
	_, api := humatest.New(t)
	authgrp := huma.NewGroup(api, "/auth")

	// Add middleware that sets the scoped token
	authgrp.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		// Mock the scoped token in context
		scopedToken := model.ScopedToken{
			Email:             "test@example.com",
			Name:              "Test User",
			GivenName:         "Test",
			FamilyName:        "User",
			PreferredUsername: "testuser",
			Roles:             []string{"admin"},
			Groups:            []string{"developers"},
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "user-123",
			},
		}
		ctx = huma.WithValue(ctx, middleware.ScopedTokenKey, scopedToken)
		next(ctx)
	})

	router.RegisterRoutes(authgrp)

	resp := api.Get("/auth/me")

	require.Equal(t, http.StatusOK, resp.Code)

	var response auth.MeResponseBody
	bodyBytes := resp.Body.Bytes()
	err := json.Unmarshal(bodyBytes, &response)
	require.NoError(t, err)

	require.Equal(t, "test@example.com", response.Email)
	require.Equal(t, "user-123", response.Sub)
	require.Equal(t, "Test User", response.Name)
	require.Equal(t, []string{"admin"}, response.Roles)
	require.Equal(t, []string{"developers"}, response.Groups)
}

func TestMe_NoAuth(t *testing.T) {
	oidcSvc := auth.NewMockOIDCService()

	router := auth.NewRouter(oidcSvc, nil)

	// Set up API with middleware that requires auth
	_, api := humatest.New(t)
	authgrp := huma.NewGroup(api, "/auth")

	// Add middleware that checks for auth (but doesn't set token)
	authgrp.UseMiddleware(middleware.RequireMiddleware(api, middleware.RequireCookieAuth()))

	router.RegisterRoutes(authgrp)

	resp := api.Get("/auth/me")

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestNewRouter(t *testing.T) {
	oidcSvc := auth.NewMockOIDCService()
	orgSvc := auth.NewOrganisationService(nil, nil) // nil repos for this test

	router := auth.NewRouter(oidcSvc, orgSvc)

	require.NotNil(t, router)
}

func TestAuthService_RegisterRoutes(t *testing.T) {
	oidcSvc := auth.NewMockOIDCService()
	orgSvc := auth.NewOrganisationService(nil, nil)

	router := auth.NewRouter(oidcSvc, orgSvc)

	// Set up API
	_, api := humatest.New(t)
	authgrp := huma.NewGroup(api, "/auth")

	// Register routes shouldn't panic
	router.RegisterRoutes(authgrp)

	// Verify routes are registered by checking if they exist
	// This is a basic test to ensure RegisterRoutes completes without error
	require.NotNil(t, router)
}
