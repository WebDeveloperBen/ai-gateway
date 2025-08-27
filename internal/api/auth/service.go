package auth

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OIDCConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	TenantID     string
}

type OIDCService struct {
	Provider     *oidc.Provider
	OAuth2Config *oauth2.Config
	Verifier     *oidc.IDTokenVerifier
}

type OIDCServiceInterface interface {
	GetOAuth2Config() *oauth2.Config
	GetVerifier() *oidc.IDTokenVerifier
}

func (s *OIDCService) GetOAuth2Config() *oauth2.Config {
	return s.OAuth2Config
}

func (s *OIDCService) GetVerifier() *oidc.IDTokenVerifier {
	return s.Verifier
}

// TODO: change this to be config driven where the user is able to register authenticaton providers and we use the database to lookup the oidc service config options

func NewOIDCService(ctx context.Context, cfg OIDCConfig) (*OIDCService, error) {
	issuer := "https://login.microsoftonline.com/" + cfg.TenantID + "/v2.0"

	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	oauth2Config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  cfg.RedirectURL,
		Scopes:       []string{"openid", "profile", "email"},
	}

	return &OIDCService{Provider: provider, OAuth2Config: oauth2Config, Verifier: verifier}, nil
}
