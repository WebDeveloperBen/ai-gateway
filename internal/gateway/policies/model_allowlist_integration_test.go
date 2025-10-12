package policies_test

import (
	"context"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/policies"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/stretchr/testify/require"
)

func TestModelAllowlistPolicyIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	t.Run("AllowlistWithAllowedModels", func(t *testing.T) {
		config := model.ModelAllowlistConfig{
			AllowedModelIDs: []string{"gpt-4", "gpt-3.5-turbo", "claude-3"},
		}

		policy := policies.NewModelAllowlistPolicy(config)

		// Test allowed models
		for _, modelID := range config.AllowedModelIDs {
			req := &policies.PreRequestContext{
				AppID: "test-app",
				Model: modelID,
			}

			err := policy.PreCheck(ctx, req)
			require.NoError(t, err, "Model %s should be allowed", modelID)
		}
	})

	t.Run("AllowlistBlocksDisallowedModels", func(t *testing.T) {
		config := model.ModelAllowlistConfig{
			AllowedModelIDs: []string{"gpt-4", "gpt-3.5-turbo"},
		}

		policy := policies.NewModelAllowlistPolicy(config)

		disallowedModels := []string{"claude-3", "gpt-4-turbo", "unknown-model"}

		for _, modelID := range disallowedModels {
			req := &policies.PreRequestContext{
				AppID: "test-app",
				Model: modelID,
			}

			err := policy.PreCheck(ctx, req)
			require.Error(t, err, "Model %s should be blocked", modelID)
			require.Contains(t, err.Error(), "model not allowed")
		}
	})

	t.Run("EmptyAllowlistAllowsAllModels", func(t *testing.T) {
		config := model.ModelAllowlistConfig{
			AllowedModelIDs: []string{}, // Empty allowlist
		}

		policy := policies.NewModelAllowlistPolicy(config)

		testModels := []string{"gpt-4", "claude-3", "any-model"}

		for _, modelID := range testModels {
			req := &policies.PreRequestContext{
				AppID: "test-app",
				Model: modelID,
			}

			err := policy.PreCheck(ctx, req)
			require.NoError(t, err, "Empty allowlist should allow all models including %s", modelID)
		}
	})

	t.Run("PostCheckIsNoOp", func(t *testing.T) {
		config := model.ModelAllowlistConfig{
			AllowedModelIDs: []string{"gpt-4"},
		}

		policy := policies.NewModelAllowlistPolicy(config)

		req := &policies.PostRequestContext{
			AppID: "test-app",
		}

		// PostCheck should not panic or error
		policy.PostCheck(ctx, req)
	})

	t.Run("PolicyTypeIsCorrect", func(t *testing.T) {
		config := model.ModelAllowlistConfig{
			AllowedModelIDs: []string{"gpt-4"},
		}

		policy := policies.NewModelAllowlistPolicy(config)

		require.Equal(t, model.PolicyTypeModelAllowlist, policy.Type())
	})
}
