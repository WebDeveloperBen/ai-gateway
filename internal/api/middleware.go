package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/insurgence-ai/llm-gateway/internal/config"
	"github.com/insurgence-ai/llm-gateway/internal/exceptions"
	"github.com/insurgence-ai/llm-gateway/internal/logger"
	"github.com/insurgence-ai/llm-gateway/internal/model"
)

type (
	RequireFunc           func(ctx huma.Context) error
	scopedTokenContextKey struct{}
)

var ScopedTokenKey = scopedTokenContextKey{}

const AuthCookieName = "llm_gateway_auth_token"

// Use is a convienance function that wraps and allows the use of middleware
func Use(api huma.API, mw func(huma.API) func(huma.Context, func(huma.Context))) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		handler := mw(api) // create new instance on each request
		handler(ctx, next)
	}
}

// AuthCookieMiddleware authenticates requests by verifying a JWT stored in an HttpOnly, Secure cookie.
// Sets user claims into the context on success and denies requests with missing or invalid cookies.
// Use for browser-based flows where the authentication token is managed via cookies.
func AuthCookieMiddleware(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		req, err := huma.ReadCookie(ctx, AuthCookieName)

		logger.Logger.Info().Msgf("[AuthCookieMiddleware Cookie]: ", req.Value)
		if err != nil || req == nil {
			logger.Logger.Info().Msgf("[AuthCookieMiddleware Error]: ", err)
			// No cookie: just let pass through
			next(ctx)
			return
		}

		claims := &model.ScopedToken{}
		cfg := config.Envs

		logger.Logger.Info().Msgf("[Claims]: %+v", claims)
		token, err := jwt.ParseWithClaims(req.Value, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNoCookie
			}
			return []byte(cfg.AuthSecret), nil
		})
		if err != nil || !token.Valid {

			logger.Logger.Info().Msgf("[Token Error]: %+v %+v", err, token)
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid or expired auth cookie")
			return
		}

		ctx = huma.WithValue(ctx, ScopedTokenKey, *claims)

		logger.Logger.Info().Msgf("[Context Values ]: %+v", ctx)
		next(ctx)
	}
}

func GetScopedToken(ctx context.Context) (model.ScopedToken, bool) {
	claims, ok := ctx.Value(ScopedTokenKey).(model.ScopedToken)
	logger.Logger.Info().Msgf("[Get Scoped Token Claims]: %+v", claims)
	return claims, ok
}

func RequireMiddleware(api huma.API, checks ...RequireFunc) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		for _, check := range checks {
			if err := check(ctx); err != nil {

				logger.Logger.Info().Msgf("[Require Middleware Error]: %+v", err)
				if err.Error() == "authentication required" || strings.Contains(err.Error(), "401") {
					huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized", err)
				} else {
					huma.WriteErr(api, ctx, http.StatusForbidden, "Forbidden", err)
				}
				return
			}
		}
		next(ctx)
	}
}

func RequireCookieAuth() RequireFunc {
	return func(ctx huma.Context) error {
		err, ok := GetScopedToken(ctx.Context())

		if !ok {
			logger.Logger.Info().Msgf("[RequireCookieAuth Error]: %+v", err)
			return exceptions.Unauthorized("authentication required")
		}
		return nil
	}
}

// TODO: update this to be relevant to this app

// func GetScopedOrganisationID(ctx context.Context) (uuid.UUID, error) {
// 	session, ok := GetUserSession(ctx)
// 	if !ok {
// 		return uuid.Nil, exceptions.Unauthorized("unauthenticated or session missing")
// 	}
//
// 	if len(session.Organisations) == 0 {
// 		return uuid.Nil, exceptions.Forbidden("no organisation context in token")
// 	}
//
// 	if len(session.Organisations) > 1 {
// 		return uuid.Nil, exceptions.Forbidden("token not scoped to a single organisation")
// 	}
//
// 	return session.Organisations[0].OrganisationID, nil
// }

// func AuthenticationMiddleware(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
// 	return func(ctx huma.Context, next func(huma.Context)) {
// 		authHeader := ctx.Header("Authorization")
// 		if authHeader == "" {
// 			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Missing Authorization header")
// 			return
// 		}
//
// 		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
// 		claims := &model.ScopedToken{}
//
// 		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, errors.New("unexpected signing method")
// 			}
// 			return []byte(config.Envs.AuthSecret), nil
// 		})
// 		if err != nil || !token.Valid {
// 			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid token")
// 			return
// 		}
//
// 		userID, err := uuid.Parse(claims.Subject)
// 		if err != nil {
// 			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid user ID in subject")
// 			return
// 		}
//
// 		session := model.UserClaims{
// 			UserID: userID,
// 			Email:  claims.Email,
// 		}
//
// 		if claims.Organisation != nil && claims.Organisation.OrganisationID != uuid.Nil {
// 			session.Organisations = []model.OrganisationMembership{*claims.Organisation}
// 		} else if len(claims.Organisations) > 0 {
// 			session.Organisations = claims.Organisations
// 		}
//
// 		ctx = huma.WithValue(ctx, model.UserClaimsKey, session)
// 		next(ctx)
// 	}
// }
