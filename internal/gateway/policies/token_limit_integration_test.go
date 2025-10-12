package policies_test

import (
	"context"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/policies"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/stretchr/testify/require"
)

func TestTokenLimitPolicyIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	t.Run("PromptTokensWithinLimitPasses", func(t *testing.T) {
		config := model.TokenLimitConfig{
			MaxPromptTokens: 1000,
		}

		policy := policies.NewTokenLimitPolicy(config)

		req := &policies.PreRequestContext{
			AppID:           "test-app",
			EstimatedTokens: 500, // Within limit
			Model:           "gpt-4",
		}

		err := policy.PreCheck(ctx, req)
		require.NoError(t, err, "Request within prompt token limit should pass")
	})

	t.Run("PromptTokensExceedingLimitFails", func(t *testing.T) {
		config := model.TokenLimitConfig{
			MaxPromptTokens: 1000,
		}

		policy := policies.NewTokenLimitPolicy(config)

		req := &policies.PreRequestContext{
			AppID:           "test-app",
			EstimatedTokens: 1500, // Exceeds limit
			Model:           "gpt-4",
		}

		err := policy.PreCheck(ctx, req)
		require.Error(t, err, "Request exceeding prompt token limit should fail")
		require.Contains(t, err.Error(), "token limit exceeded")
	})

	t.Run("TotalTokensWithinLimitPasses", func(t *testing.T) {
		config := model.TokenLimitConfig{
			MaxTotalTokens: 2000,
		}

		policy := policies.NewTokenLimitPolicy(config)

		req := &policies.PreRequestContext{
			AppID:           "test-app",
			EstimatedTokens: 800, // Estimated total will be 1600 (800 * 2)
			Model:           "gpt-4",
		}

		err := policy.PreCheck(ctx, req)
		require.NoError(t, err, "Request within total token limit should pass")
	})

	t.Run("TotalTokensExceedingLimitFails", func(t *testing.T) {
		config := model.TokenLimitConfig{
			MaxTotalTokens: 1000,
		}

		policy := policies.NewTokenLimitPolicy(config)

		req := &policies.PreRequestContext{
			AppID:           "test-app",
			EstimatedTokens: 600, // Estimated total will be 1200 (600 * 2), exceeds 1000
			Model:           "gpt-4",
		}

		err := policy.PreCheck(ctx, req)
		require.Error(t, err, "Request exceeding total token limit should fail")
		require.Contains(t, err.Error(), "token limit exceeded")
	})

	t.Run("ZeroLimitsAllowAllRequests", func(t *testing.T) {
		config := model.TokenLimitConfig{
			MaxPromptTokens: 0, // No limit
			MaxTotalTokens:  0, // No limit
		}

		policy := policies.NewTokenLimitPolicy(config)

		req := &policies.PreRequestContext{
			AppID:           "test-app",
			EstimatedTokens: 100000, // Very large request
			Model:           "gpt-4",
		}

		err := policy.PreCheck(ctx, req)
		require.NoError(t, err, "Zero limits should allow all token counts")
	})

	t.Run("PostCheckLogsCompletionTokenExceedance", func(t *testing.T) {
		config := model.TokenLimitConfig{
			MaxCompletionTokens: 500,
		}

		policy := policies.NewTokenLimitPolicy(config)

		req := &policies.PostRequestContext{
			AppID:     "test-app",
			ModelName: "gpt-4",
			ActualTokens: model.TokenUsage{
				PromptTokens:     200,
				CompletionTokens: 600, // Exceeds limit
				TotalTokens:      800,
			},
		}

		// PostCheck should log but not error
		policy.PostCheck(ctx, req)
	})

	t.Run("PostCheckLogsTotalTokenExceedance", func(t *testing.T) {
		config := model.TokenLimitConfig{
			MaxTotalTokens: 1000,
		}

		policy := policies.NewTokenLimitPolicy(config)

		req := &policies.PostRequestContext{
			AppID:     "test-app",
			ModelName: "gpt-4",
			ActualTokens: model.TokenUsage{
				PromptTokens:     400,
				CompletionTokens: 700,
				TotalTokens:      1100, // Exceeds limit
			},
		}

		// PostCheck should log but not error
		policy.PostCheck(ctx, req)
	})

	t.Run("PostCheckWithinLimits", func(t *testing.T) {
		config := model.TokenLimitConfig{
			MaxCompletionTokens: 500,
			MaxTotalTokens:      1000,
		}

		policy := policies.NewTokenLimitPolicy(config)

		req := &policies.PostRequestContext{
			AppID:     "test-app",
			ModelName: "gpt-4",
			ActualTokens: model.TokenUsage{
				PromptTokens:     300,
				CompletionTokens: 400, // Within limit
				TotalTokens:      700, // Within limit
			},
		}

		// PostCheck should not log errors
		policy.PostCheck(ctx, req)
	})

	t.Run("PolicyTypeIsCorrect", func(t *testing.T) {
		config := model.TokenLimitConfig{
			MaxPromptTokens: 1000,
		}

		policy := policies.NewTokenLimitPolicy(config)

		require.Equal(t, model.PolicyTypeTokenLimit, policy.Type())
	})
}
