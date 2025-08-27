package middleware

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/insurgence-ai/llm-gateway/internal/exceptions"
	"github.com/insurgence-ai/llm-gateway/internal/logger"
)

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

func RequireOrganisationID() RequireFunc {
	return func(ctx huma.Context) error {
		err, ok := GetOrgIDFromSession(ctx.Context())

		if !ok {

			logger.Logger.Info().Msgf("[RequireCookieAuth Error]: %+v", err)
			return exceptions.Unauthorized("authentication required")
		}
		return nil
	}
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
