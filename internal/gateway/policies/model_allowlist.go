package policies

import (
	"context"
	"fmt"
	"slices"

	"github.com/WebDeveloperBen/ai-gateway/internal/logger"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

// ModelAllowlistPolicy restricts which models an application can use
type ModelAllowlistPolicy struct {
	config model.ModelAllowlistConfig
}

// NewModelAllowlistPolicy creates a new model allowlist policy
func NewModelAllowlistPolicy(config model.ModelAllowlistConfig) *ModelAllowlistPolicy {
	return &ModelAllowlistPolicy{config: config}
}

// Type returns the policy type
func (p *ModelAllowlistPolicy) Type() model.PolicyType {
	return model.PolicyTypeModelAllowlist
}

// PreCheck checks if the requested model is in the allowlist
func (p *ModelAllowlistPolicy) PreCheck(ctx context.Context, req *PreRequestContext) error {
	if len(p.config.AllowedModelIDs) == 0 {
		// Empty allowlist means all models are allowed
		return nil
	}

	// Check if the requested model is in the allowlist
	if slices.Contains(p.config.AllowedModelIDs, req.Model) {
		return nil
	}

	logger.GetLogger(ctx).Warn().
		Str("model", req.Model).
		Str("app_id", req.AppID).
		Msg("Model not in allowlist")
	return fmt.Errorf("model not allowed")
}

// PostCheck is a no-op for model allowlist
func (p *ModelAllowlistPolicy) PostCheck(ctx context.Context, req *PostRequestContext) {
	// No post-processing needed for model allowlist
}
