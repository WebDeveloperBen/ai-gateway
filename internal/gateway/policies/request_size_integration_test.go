package policies_test

import (
	"context"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/policies"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/stretchr/testify/require"
)

func TestRequestSizePolicyIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	t.Run("RequestWithinLimitPasses", func(t *testing.T) {
		config := model.RequestSizeConfig{
			MaxRequestBytes: 1024,
		}

		policy := policies.NewRequestSizePolicy(config)

		req := &policies.PreRequestContext{
			AppID:            "test-app",
			RequestSizeBytes: 512, // Within limit
		}

		err := policy.PreCheck(ctx, req)
		require.NoError(t, err, "Request within size limit should pass")
	})

	t.Run("RequestExceedingLimitFails", func(t *testing.T) {
		config := model.RequestSizeConfig{
			MaxRequestBytes: 1024,
		}

		policy := policies.NewRequestSizePolicy(config)

		req := &policies.PreRequestContext{
			AppID:            "test-app",
			RequestSizeBytes: 2048, // Exceeds limit
		}

		err := policy.PreCheck(ctx, req)
		require.Error(t, err, "Request exceeding size limit should fail")
		require.Contains(t, err.Error(), "request size limit exceeded")
	})

	t.Run("ZeroLimitAllowsAllRequests", func(t *testing.T) {
		config := model.RequestSizeConfig{
			MaxRequestBytes: 0, // No limit
		}

		policy := policies.NewRequestSizePolicy(config)

		req := &policies.PreRequestContext{
			AppID:            "test-app",
			RequestSizeBytes: 1000000, // Very large request
		}

		err := policy.PreCheck(ctx, req)
		require.NoError(t, err, "Zero limit should allow all request sizes")
	})

	t.Run("ExactLimitPasses", func(t *testing.T) {
		config := model.RequestSizeConfig{
			MaxRequestBytes: 1024,
		}

		policy := policies.NewRequestSizePolicy(config)

		req := &policies.PreRequestContext{
			AppID:            "test-app",
			RequestSizeBytes: 1024, // Exactly at limit
		}

		err := policy.PreCheck(ctx, req)
		require.NoError(t, err, "Request at exact size limit should pass")
	})

	t.Run("PostCheckIsNoOp", func(t *testing.T) {
		config := model.RequestSizeConfig{
			MaxRequestBytes: 1024,
		}

		policy := policies.NewRequestSizePolicy(config)

		req := &policies.PostRequestContext{
			AppID: "test-app",
		}

		// PostCheck should not panic or error
		policy.PostCheck(ctx, req)
	})

	t.Run("PolicyTypeIsCorrect", func(t *testing.T) {
		config := model.RequestSizeConfig{
			MaxRequestBytes: 1024,
		}

		policy := policies.NewRequestSizePolicy(config)

		require.Equal(t, model.PolicyTypeRequestSize, policy.Type())
	})
}
