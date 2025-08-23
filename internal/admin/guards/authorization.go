package guards

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type RequireFunc func(ctx huma.Context) error

func RequireMiddleware(api huma.API, checks ...RequireFunc) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		for _, check := range checks {
			if err := check(ctx); err != nil {
				huma.WriteErr(api, ctx, http.StatusForbidden, "Forbidden", err)
				return
			}
		}
		next(ctx)
	}
}
