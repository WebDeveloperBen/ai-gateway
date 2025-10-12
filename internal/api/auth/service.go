package auth

import (
	"context"
	"fmt"

	"github.com/WebDeveloperBen/ai-gateway/internal/exceptions"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
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
	VerifyIDToken(ctx context.Context, tok *oauth2.Token) (*oidc.IDToken, map[string]any, error)
	ClaimsToScopedToken(claims map[string]any, idToken *oidc.IDToken) model.ScopedToken
}

func (s *OIDCService) GetOAuth2Config() *oauth2.Config {
	return s.OAuth2Config
}

func (s *OIDCService) GetVerifier() *oidc.IDTokenVerifier {
	return s.Verifier
}

func (s *OIDCService) VerifyIDToken(ctx context.Context, tok *oauth2.Token) (*oidc.IDToken, map[string]any, error) {
	raw, ok := tok.Extra("id_token").(string)
	if !ok {
		return nil, nil, exceptions.Unauthorized("id_token missing in token response")
	}

	idToken, err := s.GetVerifier().Verify(ctx, raw)
	if err != nil {
		return nil, nil, exceptions.Unauthorized(fmt.Sprintf("id_token verification failed: %v", err))
	}

	var claims map[string]any
	if err := idToken.Claims(&claims); err != nil {
		return nil, nil, exceptions.Unauthorized(fmt.Sprintf("id_token claim parse failed: %v", err))
	}

	return idToken, claims, nil
}

func (s *OIDCService) ClaimsToScopedToken(claims map[string]any, idToken *oidc.IDToken) model.ScopedToken {
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
