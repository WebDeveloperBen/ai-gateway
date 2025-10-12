package auth

import (
	"context"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

// MockOIDCService implements OIDCServiceInterface for testing and development
type MockOIDCService struct {
	OAuth2Config *oauth2.Config
}

func NewMockOIDCService() *MockOIDCService {
	oauth2Config := &oauth2.Config{
		ClientID:     "mock-client-id",
		ClientSecret: "mock-client-secret",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://mock-auth-url.com",
			TokenURL: "https://mock-token-url.com",
		},
		RedirectURL: "http://localhost:8080/auth/callback",
		Scopes:      []string{"openid", "profile", "email"},
	}

	return &MockOIDCService{
		OAuth2Config: oauth2Config,
	}
}

func (m *MockOIDCService) GetOAuth2Config() *oauth2.Config {
	return m.OAuth2Config
}

func (m *MockOIDCService) GetVerifier() *oidc.IDTokenVerifier {
	// For mock, return nil since it's not used in the current routes
	return nil
}

func (m *MockOIDCService) VerifyIDToken(ctx context.Context, tok *oauth2.Token) (*oidc.IDToken, map[string]any, error) {
	// Return mock ID token and claims
	idToken := &oidc.IDToken{
		Issuer:          "https://mock-issuer.com",
		Audience:        []string{"mock-client-id"},
		Subject:         "mock-user-id",
		Expiry:          time.Now().Add(time.Hour),
		IssuedAt:        time.Now(),
		Nonce:           "mock-nonce",
		AccessTokenHash: "",
	}

	claims := map[string]any{
		"sub":                "mock-user-id",
		"email":              "test@example.com",
		"name":               "Test User",
		"given_name":         "Test",
		"family_name":        "User",
		"preferred_username": "testuser",
		"roles":              []string{"admin"},
		"groups":             []string{"developers"},
	}

	return idToken, claims, nil
}

func (m *MockOIDCService) ClaimsToScopedToken(claims map[string]any, idToken *oidc.IDToken) model.ScopedToken {
	getStr := func(k string) string {
		if v, ok := claims[k].(string); ok {
			return v
		}
		return ""
	}

	getStrSlice := func(k string) []string {
		out := []string{}
		if arr, ok := claims[k].([]any); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok {
					out = append(out, s)
				}
			}
		}
		return out
	}

	return model.ScopedToken{
		Email:             getStr("email"),
		Name:              getStr("name"),
		GivenName:         getStr("given_name"),
		FamilyName:        getStr("family_name"),
		PreferredUsername: getStr("preferred_username"),
		Roles:             getStrSlice("roles"),
		Groups:            getStrSlice("groups"),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   getStr("sub"),
			Issuer:    idToken.Issuer,
			ExpiresAt: jwt.NewNumericDate(idToken.Expiry),
		},
	}
}
