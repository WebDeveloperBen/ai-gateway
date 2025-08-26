// Package auth provides authentication and authorization primitives,
// OpenID Connect (OIDC) integration with Azure Entra ID, and request/response
// DTOs for the API interface. It handles OAuth login/callback flows, secure
// JWT session encoding/validation, and related helpers for user context.
package auth

import (
	"context"
	"fmt"
	"net/http"

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
		DefaultStatus: http.StatusFound,
		Tags:          []string{"Auth"},
	}, exceptions.Handle(func(ctx context.Context, in *CallbackRequest) (*CallbackRedirect, error) {
		if in.Error != "" {
			return nil, fmt.Errorf("OIDC error: %s", in.Error)
		}

		if in.Code == "" {
			return nil, exceptions.Unauthorized("code not found")
		}

		tok, err := svc.OAuth2Config.Exchange(ctx, in.Code)
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

		scoped := model.ScopedTokenClaims{
			Email: email,
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
			Expires:  idToken.Expiry,
		}

		return &CallbackRedirect{
			SetCookies: []*http.Cookie{&cookie},
			Location:   "http://localhost:3000",
		}, nil
	}))

	huma.Register(grp, huma.Operation{
		OperationID:   "auth-me",
		Method:        http.MethodGet,
		Path:          "/me",
		Summary:       "Get user info",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Auth"},
	}, exceptions.Handle(func(ctx context.Context, _ *struct{}) (*MeResponse, error) {
		claims, ok := ctx.Value(model.UserClaimsKey).(model.ScopedTokenClaims)
		if !ok {
			return nil, fmt.Errorf("no session")
		}
		return &MeResponse{User: map[string]any{
			"email": claims.Email,
			"sub":   claims.Subject,
		}}, nil
	}))
}
