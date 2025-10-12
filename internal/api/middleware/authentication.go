package middleware

import (
	"net/http"

	"github.com/WebDeveloperBen/ai-gateway/internal/config"
	"github.com/WebDeveloperBen/ai-gateway/internal/logger"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthCookieMiddleware authenticates requests by verifying a JWT stored in an HttpOnly, Secure cookie.
// Sets user claims into the context on success and denies requests with missing or invalid cookies.
// Use for browser-based flows where the authentication token is managed via cookies.
func AuthCookieMiddleware(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		req, err := huma.ReadCookie(ctx, AuthCookieName)

		if err != nil || req == nil {
			logger.Logger.Info().Msgf("[AuthCookieMiddleware Error]: %+v", err)
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
