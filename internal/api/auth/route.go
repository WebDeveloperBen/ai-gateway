// Package auth provides authentication and authorization primitives,
// OpenID Connect (OIDC) integration with Azure Entra ID, and request/response
// DTOs for the API interface. It handles OAuth login/callback flows, secure
// JWT session encoding/validation, and related helpers for user context.
package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	middleware "github.com/insurgence-ai/llm-gateway/internal/api"
	"github.com/insurgence-ai/llm-gateway/internal/config"
	"github.com/insurgence-ai/llm-gateway/internal/exceptions"
	"github.com/insurgence-ai/llm-gateway/internal/model"
)

func RegisterAuthRoutes(grp *huma.Group, svc *OIDCService) {
	huma.Register(grp, huma.Operation{
		OperationID:   "auth-login",
		Method:        http.MethodGet,
		Path:          "/login",
		Summary:       "Get Azure login URL",
		DefaultStatus: http.StatusFound,
		Tags:          []string{"Auth"},
	}, exceptions.Handle(func(ctx context.Context, _ *struct{}) (*LoginRedirect, error) {
		state, err := generateState(32)
		if err != nil {
			return nil, fmt.Errorf("error generating oauth state")
		}
		url := svc.OAuth2Config.AuthCodeURL(state)
		return &LoginRedirect{
			Location: url,
		}, nil
	}))

	huma.Register(grp, huma.Operation{
		OperationID:   "auth-callback",
		Method:        http.MethodGet,
		Path:          "/callback",
		Summary:       "OIDC callback handler",
		DefaultStatus: http.StatusSeeOther,
		Tags:          []string{"Auth"},
	}, exceptions.Handle(func(ctx context.Context, req *CallbackRequest) (*CallbackRedirect, error) {
		if req.Error != "" {
			return nil, fmt.Errorf("OIDC error: %s", req.Error)
		}

		if req.Code == "" {
			return nil, exceptions.Unauthorized("code not found")
		}

		tok, err := svc.OAuth2Config.Exchange(ctx, req.Code)
		if err != nil {
			return nil, exceptions.Unauthorized("token exchange failed")
		}

		rawIDToken, ok := tok.Extra("id_token").(string)
		if !ok {
			return nil, exceptions.Unauthorized("id_token missing")
		}

		idToken, err := svc.Verifier.Verify(ctx, rawIDToken)
		if err != nil {
			return nil, exceptions.Unauthorized("id_token verification failed")
		}

		var rawClaims map[string]any
		if err := idToken.Claims(&rawClaims); err != nil {
			return nil, exceptions.Unauthorized(fmt.Sprintf("claim parse failed: %+v", err))
		}

		email, _ := rawClaims["email"].(string)

		sub, _ := rawClaims["sub"].(string)

		name, _ := rawClaims["name"].(string)
		givenName, _ := rawClaims["given_name"].(string)
		familyName, _ := rawClaims["family_name"].(string)
		preferredUsername, _ := rawClaims["preferred_username"].(string)
		// Cast slices carefully (OIDC may have roles/groups as []any)
		var roles, groups []string
		if r, ok := rawClaims["roles"].([]any); ok {
			for _, v := range r {
				if s, ok := v.(string); ok {
					roles = append(roles, s)
				}
			}
		}
		if g, ok := rawClaims["groups"].([]any); ok {
			for _, v := range g {
				if s, ok := v.(string); ok {
					groups = append(groups, s)
				}
			}
		}

		scoped := model.ScopedToken{
			Email:             email,
			Name:              name,
			GivenName:         givenName,
			FamilyName:        familyName,
			PreferredUsername: preferredUsername,
			Roles:             roles,
			Groups:            groups,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   sub,
				Issuer:    idToken.Issuer,
				ExpiresAt: jwt.NewNumericDate(idToken.Expiry),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, scoped)
		cfg := config.Envs

		signed, err := token.SignedString([]byte(cfg.AuthSecret))
		if err != nil {
			return nil, fmt.Errorf("JWT signing error: %w", err)
		}

		cookie := http.Cookie{
			Name:     middleware.AuthCookieName,
			Value:    signed,
			Path:     "/",
			Secure:   cfg.IsProd,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(time.Until(idToken.Expiry).Seconds()),
		}

		return &CallbackRedirect{
			SetCookie: cookie,
			Location:  "/",
		}, nil
	}))

	huma.Register(grp, huma.Operation{
		OperationID:   "auth-me",
		Method:        http.MethodGet,
		Path:          "/me",
		Summary:       "Get user info",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Auth"},
		Middlewares:   huma.Middlewares{middleware.RequireMiddleware(grp.API, middleware.RequireCookieAuth())},
	}, exceptions.Handle(func(ctx context.Context, _ *struct{}) (*MeResponse, error) {
		claims, ok := middleware.GetScopedToken(ctx)
		if !ok {
			return nil, fmt.Errorf("no session")
		}
		return &MeResponse{Body: MeResponseBody{
			Email:             claims.Email,
			Sub:               claims.Subject,
			Name:              claims.Name,
			GivenName:         claims.GivenName,
			FamilyName:        claims.FamilyName,
			PreferredUsername: claims.PreferredUsername,
			Roles:             claims.Roles,
			Groups:            claims.Groups,
		}}, nil
	}))
}
