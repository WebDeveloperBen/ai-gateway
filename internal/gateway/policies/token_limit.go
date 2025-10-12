package policies

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/WebDeveloperBen/ai-gateway/internal/logger"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

func init() {
	Register(model.PolicyTypeTokenLimit, func(config []byte, deps PolicyDependencies) (Policy, error) {
		var cfg model.TokenLimitConfig
		if err := json.Unmarshal(config, &cfg); err != nil {
			return nil, fmt.Errorf("invalid token limit config: %w", err)
		}
		return NewTokenLimitPolicy(cfg), nil
	})
}

// TokenLimitPolicy enforces maximum token limits per request
type TokenLimitPolicy struct {
	config model.TokenLimitConfig
}

// NewTokenLimitPolicy creates a new token limit policy
func NewTokenLimitPolicy(config model.TokenLimitConfig) *TokenLimitPolicy {
	return &TokenLimitPolicy{config: config}
}

// Type returns the policy type
func (p *TokenLimitPolicy) Type() model.PolicyType {
	return model.PolicyTypeTokenLimit
}

// PreCheck checks if the estimated tokens exceed the limit
func (p *TokenLimitPolicy) PreCheck(ctx context.Context, req *PreRequestContext) error {
	if p.config.MaxPromptTokens > 0 && req.EstimatedTokens > p.config.MaxPromptTokens {
		logger.GetLogger(ctx).Warn().
			Int("estimated_tokens", req.EstimatedTokens).
			Int("limit", p.config.MaxPromptTokens).
			Str("app_id", req.AppID).
			Str("model", req.Model).
			Msg("Prompt token limit exceeded")
		return fmt.Errorf("token limit exceeded")
	}

	if p.config.MaxTotalTokens > 0 {
		// Estimate total tokens (prompt + typical completion ratio)
		// For now, assume completion is roughly same size as prompt
		estimatedTotal := req.EstimatedTokens * 2
		if estimatedTotal > p.config.MaxTotalTokens {
			logger.GetLogger(ctx).Warn().
				Int("estimated_total", estimatedTotal).
				Int("limit", p.config.MaxTotalTokens).
				Str("app_id", req.AppID).
				Str("model", req.Model).
				Msg("Total token limit exceeded")
			return fmt.Errorf("token limit exceeded")
		}
	}

	return nil
}

// PostCheck validates actual token usage (async)
func (p *TokenLimitPolicy) PostCheck(ctx context.Context, req *PostRequestContext) {
	// Log if actual usage exceeded limits (for monitoring)
	if p.config.MaxCompletionTokens > 0 && req.ActualTokens.CompletionTokens > p.config.MaxCompletionTokens {
		logger.GetLogger(ctx).Warn().
			Int("completion_tokens", req.ActualTokens.CompletionTokens).
			Int("limit", p.config.MaxCompletionTokens).
			Str("app_id", req.AppID).
			Str("model", req.ModelName).
			Msg("Completion tokens exceeded limit (post-check)")
	}

	if p.config.MaxTotalTokens > 0 && req.ActualTokens.TotalTokens > p.config.MaxTotalTokens {
		logger.GetLogger(ctx).Warn().
			Int("total_tokens", req.ActualTokens.TotalTokens).
			Int("limit", p.config.MaxTotalTokens).
			Str("app_id", req.AppID).
			Str("model", req.ModelName).
			Msg("Total tokens exceeded limit (post-check)")
	}
}
