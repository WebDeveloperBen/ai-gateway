package policies_test

import (
	"context"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/policies"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/stretchr/testify/require"
)

func TestCELPolicyIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	t.Run("PreCheckExpressionAllowsValidRequests", func(t *testing.T) {
		config := `{
			"pre_check_expression": "estimated_tokens < 1000 && model == 'gpt-4'"
		}`

		policy, err := policies.NewCELPolicy(model.PolicyTypeCustomCEL, []byte(config))
		require.NoError(t, err)

		req := &policies.PreRequestContext{
			AppID:            "test-app",
			OrgID:            "test-org",
			Model:            "gpt-4",
			EstimatedTokens:  500,
			RequestSizeBytes: 1024,
		}

		err = policy.PreCheck(ctx, req)
		require.NoError(t, err, "Valid request should pass CEL expression")
	})

	t.Run("PreCheckExpressionBlocksInvalidRequests", func(t *testing.T) {
		config := `{
			"pre_check_expression": "estimated_tokens < 1000 && model == 'gpt-4'"
		}`

		policy, err := policies.NewCELPolicy(model.PolicyTypeCustomCEL, []byte(config))
		require.NoError(t, err)

		// Test cases that should fail
		testCases := []struct {
			name string
			req  *policies.PreRequestContext
		}{
			{
				name: "too many tokens",
				req: &policies.PreRequestContext{
					AppID:            "test-app",
					OrgID:            "test-org",
					Model:            "gpt-4",
					EstimatedTokens:  1500, // Exceeds 1000
					RequestSizeBytes: 1024,
				},
			},
			{
				name: "wrong model",
				req: &policies.PreRequestContext{
					AppID:            "test-app",
					OrgID:            "test-org",
					Model:            "gpt-3.5-turbo", // Not gpt-4
					EstimatedTokens:  500,
					RequestSizeBytes: 1024,
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := policy.PreCheck(ctx, tc.req)
				require.Error(t, err, "Invalid request should fail CEL expression")
				require.Contains(t, err.Error(), "CEL expression returned false")
			})
		}
	})

	t.Run("ComplexPreCheckExpression", func(t *testing.T) {
		config := `{
			"pre_check_expression": "estimated_tokens <= 2000 && (model == 'gpt-4' || model == 'claude-3') && request_size_bytes < 50000"
		}`

		policy, err := policies.NewCELPolicy(model.PolicyTypeCustomCEL, []byte(config))
		require.NoError(t, err)

		// Valid request
		req := &policies.PreRequestContext{
			AppID:            "test-app",
			OrgID:            "test-org",
			Model:            "gpt-4",
			EstimatedTokens:  1000,
			RequestSizeBytes: 10000,
		}

		err = policy.PreCheck(ctx, req)
		require.NoError(t, err, "Valid complex request should pass")

		// Invalid request - too large
		req.RequestSizeBytes = 60000
		err = policy.PreCheck(ctx, req)
		require.Error(t, err, "Oversized request should fail")
	})

	t.Run("PostCheckExpressionEvaluation", func(t *testing.T) {
		config := `{
			"pre_check_expression": "estimated_tokens < 1000",
			"post_check_expression": "total_tokens < 2000 && latency_ms < 30000"
		}`

		policy, err := policies.NewCELPolicy(model.PolicyTypeCustomCEL, []byte(config))
		require.NoError(t, err)

		// Pre-check should pass
		preReq := &policies.PreRequestContext{
			AppID:           "test-app",
			OrgID:           "test-org",
			Model:           "gpt-4",
			EstimatedTokens: 500,
		}

		err = policy.PreCheck(ctx, preReq)
		require.NoError(t, err, "Pre-check should pass")

		// Post-check should not error (just logs)
		postReq := &policies.PostRequestContext{
			AppID:     "test-app",
			OrgID:     "test-org",
			ModelName: "gpt-4",
			ActualTokens: model.TokenUsage{
				PromptTokens:     400,
				CompletionTokens: 300,
				TotalTokens:      700,
			},
			LatencyMs:         5000,
			ResponseSizeBytes: 2048,
		}

		// Should not panic
		policy.PostCheck(ctx, postReq)
	})

	t.Run("InvalidCELExpressionFailsCompilation", func(t *testing.T) {
		config := `{
			"pre_check_expression": "invalid_syntax +++"
		}`

		_, err := policies.NewCELPolicy(model.PolicyTypeCustomCEL, []byte(config))
		require.Error(t, err, "Invalid CEL expression should fail compilation")
		require.Contains(t, err.Error(), "failed to compile")
	})

	t.Run("NonBooleanExpressionFails", func(t *testing.T) {
		config := `{
			"pre_check_expression": "estimated_tokens + 100"
		}`

		policy, err := policies.NewCELPolicy(model.PolicyTypeCustomCEL, []byte(config))
		require.NoError(t, err)

		req := &policies.PreRequestContext{
			AppID:           "test-app",
			EstimatedTokens: 500,
		}

		err = policy.PreCheck(ctx, req)
		require.Error(t, err, "Non-boolean expression should fail")
		require.Contains(t, err.Error(), "must return boolean")
	})

	t.Run("MissingPreCheckExpression", func(t *testing.T) {
		config := `{
			"post_check_expression": "total_tokens < 1000"
		}`

		policy, err := policies.NewCELPolicy(model.PolicyTypeCustomCEL, []byte(config))
		require.NoError(t, err)

		req := &policies.PreRequestContext{
			AppID:           "test-app",
			EstimatedTokens: 5000, // Would fail if expression existed
		}

		// Should pass since no pre-check expression
		err = policy.PreCheck(ctx, req)
		require.NoError(t, err, "Missing pre-check expression should allow all requests")
	})

	t.Run("PolicyTypeIsPreserved", func(t *testing.T) {
		config := `{
			"pre_check_expression": "true"
		}`

		policyType := model.PolicyType("custom_test_policy")
		policy, err := policies.NewCELPolicy(policyType, []byte(config))
		require.NoError(t, err)

		require.Equal(t, policyType, policy.Type())
	})
}
