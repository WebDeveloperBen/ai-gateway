// middleware.go
package middleware

import (
	"github.com/danielgtaylor/huma/v2"
)

const orgQueriesKey ctxKey = "orgQueries"

type (
	ctxKey                string
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
