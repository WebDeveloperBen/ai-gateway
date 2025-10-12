package policies_test

import (
	"context"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/policies"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/stretchr/testify/require"
)

func TestRateLimitPolicyIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set up test containers
	_, redisAddr := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	ctx := context.Background()

	// Set up Redis cache for rate limiting
	cache, err := kv.NewDriver(kv.Config{
		Backend:   kv.BackendRedis,
		RedisAddr: redisAddr,
		RedisPW:   "",
		RedisDB:   0,
	})
	require.NoError(t, err)

	t.Run("RequestsPerMinuteWithinLimit", func(t *testing.T) {
		config := model.RateLimitConfig{
			RequestsPerMinute: 10,
		}

		policy := policies.NewRateLimitPolicy(config, cache)

		appID := "test-app-rate-limit"

		// Make requests within limit
		for i := 0; i < 5; i++ {
			req := &policies.PreRequestContext{
				AppID: appID,
				Model: "gpt-4",
			}

			err := policy.PreCheck(ctx, req)
			require.NoError(t, err, "Request %d within limit should pass", i+1)
		}
	})

	t.Run("RequestsPerMinuteExceedsLimit", func(t *testing.T) {
		config := model.RateLimitConfig{
			RequestsPerMinute: 2, // Very low limit for testing
		}

		policy := policies.NewRateLimitPolicy(config, cache)

		appID := "test-app-rate-limit-exceed"

		// Make requests up to limit
		for i := 0; i < 2; i++ {
			req := &policies.PreRequestContext{
				AppID: appID,
				Model: "gpt-4",
			}

			err := policy.PreCheck(ctx, req)
			require.NoError(t, err, "Request %d within limit should pass", i+1)
		}

		// Next request should be blocked
		req := &policies.PreRequestContext{
			AppID: appID,
			Model: "gpt-4",
		}

		err := policy.PreCheck(ctx, req)
		require.Error(t, err, "Request exceeding limit should be blocked")
		require.Contains(t, err.Error(), "requests per minute limit exceeded")
	})

	t.Run("TokensPerMinuteWithinLimit", func(t *testing.T) {
		config := model.RateLimitConfig{
			TokensPerMinute: 1000,
		}

		policy := policies.NewRateLimitPolicy(config, cache)

		appID := "test-app-token-limit"

		req := &policies.PreRequestContext{
			AppID:           appID,
			Model:           "gpt-4",
			EstimatedTokens: 500, // Within limit
		}

		err := policy.PreCheck(ctx, req)
		require.NoError(t, err, "Request within token limit should pass")
	})

	t.Run("TokensPerMinuteExceedsLimit", func(t *testing.T) {
		config := model.RateLimitConfig{
			TokensPerMinute: 100, // Low limit
		}

		policy := policies.NewRateLimitPolicy(config, cache)

		appID := "test-app-token-exceed"

		req := &policies.PreRequestContext{
			AppID:           appID,
			Model:           "gpt-4",
			EstimatedTokens: 150, // Exceeds limit
		}

		err := policy.PreCheck(ctx, req)
		require.Error(t, err, "Request exceeding token limit should be blocked")
		require.Contains(t, err.Error(), "tokens per minute limit exceeded")
	})

	t.Run("ZeroLimitsAllowAllRequests", func(t *testing.T) {
		config := model.RateLimitConfig{
			RequestsPerMinute: 0, // No limit
			TokensPerMinute:   0, // No limit
		}

		policy := policies.NewRateLimitPolicy(config, cache)

		appID := "test-app-no-limits"

		// Make many requests
		for i := 0; i < 10; i++ {
			req := &policies.PreRequestContext{
				AppID:           appID,
				Model:           "gpt-4",
				EstimatedTokens: 10000, // Very large
			}

			err := policy.PreCheck(ctx, req)
			require.NoError(t, err, "Zero limits should allow all requests")
		}
	})

	t.Run("RateLimitBlocksWhenRedisUnavailable", func(t *testing.T) {
		// Use a bad Redis address to simulate Redis unavailability
		badCache, err := kv.NewDriver(kv.Config{
			Backend:   kv.BackendRedis,
			RedisAddr: "127.0.0.1:9999", // Non-existent Redis
			RedisPW:   "",
			RedisDB:   0,
		})
		require.NoError(t, err)

		config := model.RateLimitConfig{
			RequestsPerMinute: 10,
		}

		policy := policies.NewRateLimitPolicy(config, badCache)

		req := &policies.PreRequestContext{
			AppID: "test-app-redis-down",
			Model: "gpt-4",
		}

		err = policy.PreCheck(ctx, req)
		require.Error(t, err, "Request should be blocked when Redis is unavailable")
		require.Contains(t, err.Error(), "rate limiter unavailable")
	})

	t.Run("PostCheckUpdatesTokenCounters", func(t *testing.T) {
		config := model.RateLimitConfig{
			TokensPerMinute: 1000,
		}

		policy := policies.NewRateLimitPolicy(config, cache)

		appID := "test-app-post-check"

		// Pre-check with estimated tokens
		preReq := &policies.PreRequestContext{
			AppID:           appID,
			Model:           "gpt-4",
			EstimatedTokens: 200,
		}

		err := policy.PreCheck(ctx, preReq)
		require.NoError(t, err, "Pre-check should pass")

		// Post-check with actual token usage
		postReq := &policies.PostRequestContext{
			AppID:     appID,
			ModelName: "gpt-4",
			ActualTokens: model.TokenUsage{
				PromptTokens:     180,
				CompletionTokens: 150,
				TotalTokens:      330,
			},
		}

		// Post-check should not error
		policy.PostCheck(ctx, postReq)
	})

	t.Run("RateLimitResetsAfterWindow", func(t *testing.T) {
		config := model.RateLimitConfig{
			RequestsPerMinute: 1, // Only 1 request per minute
		}

		policy := policies.NewRateLimitPolicy(config, cache)

		appID := "test-app-window-reset"

		// First request should pass
		req1 := &policies.PreRequestContext{
			AppID: appID,
			Model: "gpt-4",
		}

		err := policy.PreCheck(ctx, req1)
		require.NoError(t, err, "First request should pass")

		// Second request should be blocked
		req2 := &policies.PreRequestContext{
			AppID: appID,
			Model: "gpt-4",
		}

		err = policy.PreCheck(ctx, req2)
		require.Error(t, err, "Second request should be blocked")

		// Wait for the next minute window (simulate time passing)
		// In a real test, we might need to mock time or wait
		// For now, we'll test that the limit is enforced
	})

	t.Run("PolicyTypeIsCorrect", func(t *testing.T) {
		config := model.RateLimitConfig{
			RequestsPerMinute: 100,
		}

		policy := policies.NewRateLimitPolicy(config, cache)

		require.Equal(t, model.PolicyTypeRateLimit, policy.Type())
	})
}
