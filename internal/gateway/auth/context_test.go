package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextKeyID(t *testing.T) {
	ctx := context.Background()
	keyID := "test-key-id"

	t.Run("WithKeyID and GetKeyID", func(t *testing.T) {
		ctx = WithKeyID(ctx, keyID)
		assert.Equal(t, keyID, GetKeyID(ctx))
	})

	t.Run("GetKeyID returns empty string when not set", func(t *testing.T) {
		ctx := context.Background()
		assert.Equal(t, "", GetKeyID(ctx))
	})

	t.Run("GetKeyID returns empty string for wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), contextKeyKeyID, 123)
		assert.Equal(t, "", GetKeyID(ctx))
	})
}

func TestContextOrgID(t *testing.T) {
	ctx := context.Background()
	orgID := "test-org-id"

	t.Run("WithOrgID and GetOrgID", func(t *testing.T) {
		ctx = WithOrgID(ctx, orgID)
		assert.Equal(t, orgID, GetOrgID(ctx))
	})

	t.Run("GetOrgID returns empty string when not set", func(t *testing.T) {
		ctx := context.Background()
		assert.Equal(t, "", GetOrgID(ctx))
	})

	t.Run("GetOrgID returns empty string for wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), contextKeyOrgID, 123)
		assert.Equal(t, "", GetOrgID(ctx))
	})
}

func TestContextAppID(t *testing.T) {
	ctx := context.Background()
	appID := "test-app-id"

	t.Run("WithAppID and GetAppID", func(t *testing.T) {
		ctx = WithAppID(ctx, appID)
		assert.Equal(t, appID, GetAppID(ctx))
	})

	t.Run("GetAppID returns empty string when not set", func(t *testing.T) {
		ctx := context.Background()
		assert.Equal(t, "", GetAppID(ctx))
	})

	t.Run("GetAppID returns empty string for wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), contextKeyAppID, 123)
		assert.Equal(t, "", GetAppID(ctx))
	})
}

func TestContextUserID(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-id"

	t.Run("WithUserID and GetUserID", func(t *testing.T) {
		ctx = WithUserID(ctx, userID)
		assert.Equal(t, userID, GetUserID(ctx))
	})

	t.Run("GetUserID returns empty string when not set", func(t *testing.T) {
		ctx := context.Background()
		assert.Equal(t, "", GetUserID(ctx))
	})

	t.Run("GetUserID returns empty string for wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), contextKeyUserID, 123)
		assert.Equal(t, "", GetUserID(ctx))
	})
}

func TestContextProvider(t *testing.T) {
	ctx := context.Background()
	provider := "openai"

	t.Run("WithProvider and GetProvider", func(t *testing.T) {
		ctx = WithProvider(ctx, provider)
		assert.Equal(t, provider, GetProvider(ctx))
	})

	t.Run("GetProvider returns empty string when not set", func(t *testing.T) {
		ctx := context.Background()
		assert.Equal(t, "", GetProvider(ctx))
	})

	t.Run("GetProvider returns empty string for wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), contextKeyProvider, 123)
		assert.Equal(t, "", GetProvider(ctx))
	})
}

func TestContextModelName(t *testing.T) {
	ctx := context.Background()
	modelName := "gpt-4"

	t.Run("WithModelName and GetModelName", func(t *testing.T) {
		ctx = WithModelName(ctx, modelName)
		assert.Equal(t, modelName, GetModelName(ctx))
	})

	t.Run("GetModelName returns empty string when not set", func(t *testing.T) {
		ctx := context.Background()
		assert.Equal(t, "", GetModelName(ctx))
	})

	t.Run("GetModelName returns empty string for wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), contextKeyModel, 123)
		assert.Equal(t, "", GetModelName(ctx))
	})
}

func TestContextPolicies(t *testing.T) {
	ctx := context.Background()
	policies := []string{"policy1", "policy2"}

	t.Run("WithPolicies and GetPolicies", func(t *testing.T) {
		ctx = WithPolicies(ctx, policies)
		result := GetPolicies(ctx)
		assert.NotNil(t, result)
		assert.Equal(t, policies, result)
	})

	t.Run("GetPolicies returns nil when not set", func(t *testing.T) {
		ctx := context.Background()
		assert.Nil(t, GetPolicies(ctx))
	})

	t.Run("GetPolicies preserves type", func(t *testing.T) {
		type CustomPolicy struct {
			Name string
		}
		policy := CustomPolicy{Name: "test"}
		ctx = WithPolicies(ctx, policy)
		result := GetPolicies(ctx)
		assert.Equal(t, policy, result)
	})
}

func TestContextChaining(t *testing.T) {
	t.Run("Multiple context values can be chained", func(t *testing.T) {
		ctx := context.Background()
		ctx = WithKeyID(ctx, "key-123")
		ctx = WithOrgID(ctx, "org-456")
		ctx = WithAppID(ctx, "app-789")
		ctx = WithUserID(ctx, "user-101")
		ctx = WithProvider(ctx, "openai")
		ctx = WithModelName(ctx, "gpt-4")
		ctx = WithPolicies(ctx, []string{"policy1"})

		assert.Equal(t, "key-123", GetKeyID(ctx))
		assert.Equal(t, "org-456", GetOrgID(ctx))
		assert.Equal(t, "app-789", GetAppID(ctx))
		assert.Equal(t, "user-101", GetUserID(ctx))
		assert.Equal(t, "openai", GetProvider(ctx))
		assert.Equal(t, "gpt-4", GetModelName(ctx))
		assert.Equal(t, []string{"policy1"}, GetPolicies(ctx))
	})
}
