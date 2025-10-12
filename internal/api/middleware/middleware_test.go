package middleware_test

import (
	"context"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/api/middleware"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
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
