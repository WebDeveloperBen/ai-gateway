package auth_test

import (
	"context"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/api/auth"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestMockOIDCService(t *testing.T) {
	svc := auth.NewMockOIDCService()

	t.Run("GetOAuth2Config", func(t *testing.T) {
		config := svc.GetOAuth2Config()
		require.NotNil(t, config)
		require.Equal(t, "mock-client-id", config.ClientID)
		require.Equal(t, "http://localhost:8080/auth/callback", config.RedirectURL)
	})

	t.Run("VerifyIDToken", func(t *testing.T) {
		token := &oauth2.Token{}
		token = token.WithExtra(map[string]interface{}{
			"id_token": "mock.jwt.token",
		})

		idToken, claims, err := svc.VerifyIDToken(context.Background(), token)
		require.NoError(t, err)
		require.NotNil(t, idToken)
		require.NotNil(t, claims)
		require.Equal(t, "mock-user-id", claims["sub"])
		require.Equal(t, "test@example.com", claims["email"])
	})

	t.Run("ClaimsToScopedToken", func(t *testing.T) {
		claims := map[string]any{
			"sub":                "user-123",
			"email":              "test@example.com",
			"name":               "Test User",
			"given_name":         "Test",
			"family_name":        "User",
			"preferred_username": "testuser",
			"roles":              []any{"admin"},
			"groups":             []any{"developers"},
		}

		idToken := &oidc.IDToken{
			Issuer: "https://mock-issuer.com",
		}

		scoped := svc.ClaimsToScopedToken(claims, idToken)

		require.Equal(t, "test@example.com", scoped.Email)
		require.Equal(t, "Test User", scoped.Name)
		require.Equal(t, "user-123", scoped.Subject)
		require.Equal(t, []string{"admin"}, scoped.Roles)
		require.Equal(t, []string{"developers"}, scoped.Groups)
		require.Equal(t, "https://mock-issuer.com", scoped.Issuer)
	})
}
