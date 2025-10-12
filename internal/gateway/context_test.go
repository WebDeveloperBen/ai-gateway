package gateway_test

import (
	"context"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway"
	"github.com/stretchr/testify/require"
)

func TestTenantFrom(t *testing.T) {
	t.Run("no tenant in context returns empty string", func(t *testing.T) {
		ctx := context.Background()
		tenant := gateway.TenantFrom(ctx)
		require.Empty(t, tenant)
	})

	t.Run("nil context returns empty string", func(t *testing.T) {
		tenant := gateway.TenantFrom(nil)
		require.Empty(t, tenant)
	})
}

func TestAppFrom(t *testing.T) {
	t.Run("no app in context returns empty string", func(t *testing.T) {
		ctx := context.Background()
		app := gateway.AppFrom(ctx)
		require.Empty(t, app)
	})

	t.Run("nil context returns empty string", func(t *testing.T) {
		app := gateway.AppFrom(nil)
		require.Empty(t, app)
	})
}
