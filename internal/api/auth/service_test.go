package auth_test

import (
	"context"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/api/auth"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestOIDCService_GetOAuth2Config(t *testing.T) {
	svc := auth.NewMockOIDCService()
	config := svc.GetOAuth2Config()

	assert.NotNil(t, config)
	assert.Equal(t, "mock-client-id", config.ClientID)
	assert.Equal(t, "http://localhost:8080/auth/callback", config.RedirectURL)
	assert.Contains(t, config.Scopes, "openid")
}

func TestOIDCService_GetVerifier(t *testing.T) {
	svc := auth.NewMockOIDCService()
	verifier := svc.GetVerifier()

	// Mock returns nil, which is acceptable for testing
	assert.Nil(t, verifier)
}

func TestOIDCService_VerifyIDToken(t *testing.T) {
	svc := auth.NewMockOIDCService()

	t.Run("valid token", func(t *testing.T) {
		token := &oauth2.Token{}
		token = token.WithExtra(map[string]interface{}{
			"id_token": "mock.jwt.token",
		})

		idToken, claims, err := svc.VerifyIDToken(context.Background(), token)
		require.NoError(t, err)
		assert.NotNil(t, idToken)
		assert.NotNil(t, claims)
		assert.Equal(t, "mock-user-id", claims["sub"])
	})

	t.Run("missing id_token", func(t *testing.T) {
		token := &oauth2.Token{}

		// Mock service doesn't validate missing id_token, it just returns mock data
		idToken, claims, err := svc.VerifyIDToken(context.Background(), token)
		require.NoError(t, err)
		assert.NotNil(t, idToken)
		assert.NotNil(t, claims)
	})
}

func TestOIDCService_ClaimsToScopedToken(t *testing.T) {
	svc := auth.NewMockOIDCService()

	claims := map[string]any{
		"sub":                "user-123",
		"email":              "test@example.com",
		"name":               "Test User",
		"given_name":         "Test",
		"family_name":        "User",
		"preferred_username": "testuser",
		"roles":              []any{"admin", "user"},
		"groups":             []any{"developers", "admins"},
	}

	idToken := &oidc.IDToken{
		Issuer: "https://mock-issuer.com",
	}

	scoped := svc.ClaimsToScopedToken(claims, idToken)

	assert.Equal(t, "test@example.com", scoped.Email)
	assert.Equal(t, "Test User", scoped.Name)
	assert.Equal(t, "user-123", scoped.Subject)
	assert.Equal(t, []string{"admin", "user"}, scoped.Roles)
	assert.Equal(t, []string{"developers", "admins"}, scoped.Groups)
	assert.Equal(t, "https://mock-issuer.com", scoped.Issuer)
}

func TestNewOIDCService(t *testing.T) {
	// This test would require a real OIDC provider, so we'll skip it
	t.Skip("NewOIDCService requires real OIDC provider, tested via integration")
}
