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
	"github.com/insurgence-ai/llm-gateway/internal/api/middleware"
	"github.com/insurgence-ai/llm-gateway/internal/api/organisations"
	"github.com/insurgence-ai/llm-gateway/internal/config"
	"github.com/insurgence-ai/llm-gateway/internal/exceptions"
)

type AuthService struct {
	oidcService OIDCServiceInterface
	orgService  organisations.OrganisationServiceInterface
}

func NewRouter(oidc OIDCServiceInterface, orgSvc organisations.OrganisationServiceInterface) *AuthService {
	return &AuthService{oidcService: oidc, orgService: orgSvc}
}

func (svc *AuthService) RegisterRoutes(grp *huma.Group) {
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
		url := svc.oidcService.GetOAuth2Config().AuthCodeURL(state)
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
			return nil, exceptions.Unauthorized("oidc error")
		}

		if req.Code == "" {
			return nil, exceptions.Unauthorized("response code not found")
		}

		tok, err := svc.oidcService.GetOAuth2Config().Exchange(ctx, req.Code)
		if err != nil {
			return nil, exceptions.Unauthorized("token exchange failed")
		}

		idToken, claims, err := svc.oidcService.VerifyIDToken(ctx, tok)
		if err != nil {
			return nil, exceptions.Unauthorized(err.Error())
		}

		scoped := svc.oidcService.ClaimsToScopedToken(claims, idToken)
		if scoped.Subject == "" {
			return nil, fmt.Errorf("OIDC claims or ID token missing required fields: %+v", claims)
		}

		user, org, err := svc.orgService.EnsureUserAndOrg(ctx, scoped)
		if err != nil || user == nil || org == nil {
			fmt.Printf("EnsureUserAndOrg failed: user=%+v org=%+v err=%v scoped=%+v\n", user, org, err, scoped)
			return nil, fmt.Errorf("failed to persist user/org: %w", err)
		}

		// Hydrate ScopedToken with org info
		scoped.OrgID = org.ID
		scoped.UserID = user.ID
		scoped.Roles = user.Roles

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
