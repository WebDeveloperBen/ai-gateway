package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/api/middleware"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/require"
)

func TestGetScopedToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		scopedToken := model.ScopedToken{
			Email: "test@example.com",
			Name:  "Test User",
		}

		ctx := context.WithValue(context.Background(), middleware.ScopedTokenKey, scopedToken)
		token, ok := middleware.GetScopedToken(ctx)

		require.True(t, ok)
		require.Equal(t, "test@example.com", token.Email)
		require.Equal(t, "Test User", token.Name)
	})

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()
		_, ok := middleware.GetScopedToken(ctx)

		require.False(t, ok)
	})
}

func TestGetOrgIDFromSession(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		scopedToken := model.ScopedToken{
			OrgID: "org-123",
		}

		ctx := context.WithValue(context.Background(), middleware.ScopedTokenKey, scopedToken)
		orgID, ok := middleware.GetOrgIDFromSession(ctx)

		require.True(t, ok)
		require.Equal(t, "org-123", orgID)
	})

	t.Run("empty org id", func(t *testing.T) {
		scopedToken := model.ScopedToken{
			OrgID: "",
		}

		ctx := context.WithValue(context.Background(), middleware.ScopedTokenKey, scopedToken)
		_, ok := middleware.GetOrgIDFromSession(ctx)

		require.False(t, ok)
	})

	t.Run("no token", func(t *testing.T) {
		ctx := context.Background()
		_, ok := middleware.GetOrgIDFromSession(ctx)

		require.False(t, ok)
	})
}

func TestGetOrgQueries(t *testing.T) {
	t.Run("nil context", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when orgQueriesKey not set in context")
			}
		}()
		middleware.GetOrgQueries(context.Background())
	})
}

func TestRequireCookieAuth(t *testing.T) {
	t.Run("success with valid token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		ctx := humatest.NewContext(&huma.Operation{}, req, w)
		ctx = huma.WithValue(ctx, middleware.ScopedTokenKey, model.ScopedToken{
			Email: "test@example.com",
		})

		requireFunc := middleware.RequireCookieAuth()
		err := requireFunc(ctx)

		require.NoError(t, err)
	})

	t.Run("failure without token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		ctx := humatest.NewContext(&huma.Operation{}, req, w)

		requireFunc := middleware.RequireCookieAuth()
		err := requireFunc(ctx)

		require.Error(t, err)
		require.Contains(t, err.Error(), "authentication required")
	})
}

func TestRequireOrganisationID(t *testing.T) {
	t.Run("success with valid org ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		ctx := humatest.NewContext(&huma.Operation{}, req, w)
		ctx = huma.WithValue(ctx, middleware.ScopedTokenKey, model.ScopedToken{
			OrgID: "org-123",
		})

		requireFunc := middleware.RequireOrganisationID()
		err := requireFunc(ctx)

		require.NoError(t, err)
	})

	t.Run("failure without org ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		ctx := humatest.NewContext(&huma.Operation{}, req, w)
		ctx = huma.WithValue(ctx, middleware.ScopedTokenKey, model.ScopedToken{
			Email: "test@example.com",
		})

		requireFunc := middleware.RequireOrganisationID()
		err := requireFunc(ctx)

		require.Error(t, err)
		require.Contains(t, err.Error(), "authentication required")
	})
}

func TestRequireMiddleware(t *testing.T) {
	t.Run("success with all checks passing", func(t *testing.T) {
		_, api := humatest.New(t)
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		ctx := humatest.NewContext(&huma.Operation{}, req, w)
		ctx = huma.WithValue(ctx, middleware.ScopedTokenKey, model.ScopedToken{
			Email: "test@example.com",
			OrgID: "org-123",
		})

		mw := middleware.RequireMiddleware(api, middleware.RequireCookieAuth(), middleware.RequireOrganisationID())

		called := false
		mw(ctx, func(c huma.Context) {
			called = true
		})

		require.True(t, called)
	})

	t.Run("failure when check fails", func(t *testing.T) {
		_, api := humatest.New(t)
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		ctx := humatest.NewContext(&huma.Operation{}, req, w)
		// No scoped token set

		mw := middleware.RequireMiddleware(api, middleware.RequireCookieAuth())

		called := false
		mw(ctx, func(c huma.Context) {
			called = true
		})

		require.False(t, called)
		// Check that error was written to response
		require.Equal(t, 401, w.Code)
	})
}

func TestAuthCookieMiddleware(t *testing.T) {
	t.Run("no cookie present", func(t *testing.T) {
		_, api := humatest.New(t)
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		ctx := humatest.NewContext(&huma.Operation{}, req, w)

		mw := middleware.AuthCookieMiddleware(api)

		called := false
		mw(ctx, func(c huma.Context) {
			called = true
		})

		// Should continue to next middleware even without cookie
		require.True(t, called)
	})

	t.Run("invalid cookie causes unauthorized", func(t *testing.T) {
		_, api := humatest.New(t)
		req := httptest.NewRequest("GET", "/test", nil)
		req.AddCookie(&http.Cookie{Name: middleware.AuthCookieName, Value: "invalid"})
		w := httptest.NewRecorder()
		ctx := humatest.NewContext(&huma.Operation{}, req, w)

		mw := middleware.AuthCookieMiddleware(api)

		called := false
		mw(ctx, func(c huma.Context) {
			called = true
		})

		// Should not continue to next middleware
		require.False(t, called)
		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestUse(t *testing.T) {
	t.Run("wraps middleware correctly", func(t *testing.T) {
		_, api := humatest.New(t)
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		ctx := humatest.NewContext(&huma.Operation{}, req, w)

		// Create a test middleware that adds a value to context
		testMiddleware := func(api huma.API) func(huma.Context, func(huma.Context)) {
			return func(ctx huma.Context, next func(huma.Context)) {
				ctx = huma.WithValue(ctx, "test", "value")
				next(ctx)
			}
		}

		wrapped := middleware.Use(api, testMiddleware)

		called := false
		wrapped(ctx, func(c huma.Context) {
			called = true
			val := c.Context().Value("test")
			require.Equal(t, "value", val)
		})

		require.True(t, called)
	})
}
