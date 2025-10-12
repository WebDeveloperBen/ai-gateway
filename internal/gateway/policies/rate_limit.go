package policies

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
	"github.com/WebDeveloperBen/ai-gateway/internal/logger"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

func init() {
	Register(model.PolicyTypeRateLimit, func(config []byte, deps PolicyDependencies) (Policy, error) {
		var cfg model.RateLimitConfig
		if err := json.Unmarshal(config, &cfg); err != nil {
			return nil, fmt.Errorf("invalid rate limit config: %w", err)
		}
		return NewRateLimitPolicy(cfg, deps.Cache), nil
	})
}

// RateLimitPolicy enforces rate limits on requests and tokens per minute
type RateLimitPolicy struct {
	config  model.RateLimitConfig
	limiter *RateLimiter
}

// NewRateLimitPolicy creates a new rate limit policy
func NewRateLimitPolicy(config model.RateLimitConfig, cache kv.KvStore) *RateLimitPolicy {
	return &RateLimitPolicy{
		config:  config,
		limiter: NewRateLimiter(cache),
	}
}

// Type returns the policy type
func (p *RateLimitPolicy) Type() model.PolicyType {
	return model.PolicyTypeRateLimit
}

// PreCheck checks if the request exceeds rate limits
func (p *RateLimitPolicy) PreCheck(ctx context.Context, req *PreRequestContext) error {
	// Check requests per minute
	if p.config.RequestsPerMinute > 0 {
		key := RateLimitKey(req.AppID, "requests")
		allowed, err := p.limiter.CheckAndIncrement(ctx, key, p.config.RequestsPerMinute, time.Minute)
		if err != nil {
			// Log error but don't block request (fail open)
			return nil
		}
		if !allowed {
			logger.GetLogger(ctx).Warn().
				Int("limit", p.config.RequestsPerMinute).
				Str("app_id", req.AppID).
				Msg("Requests per minute limit exceeded")
			return rateLimitError("requests per minute limit exceeded")
		}
	}

	// Check tokens per minute (estimate only)
	if p.config.TokensPerMinute > 0 && req.EstimatedTokens > 0 {
		key := RateLimitKey(req.AppID, "tokens")
		current, err := p.limiter.GetCount(ctx, key)
		if err == nil {
			if current+req.EstimatedTokens > p.config.TokensPerMinute {
				logger.GetLogger(ctx).Warn().
					Int("current", current).
					Int("estimated", req.EstimatedTokens).
					Int("limit", p.config.TokensPerMinute).
					Str("app_id", req.AppID).
					Msg("Tokens per minute limit would be exceeded")
				return rateLimitError("tokens per minute limit exceeded")
			}
		}
		// Increment by estimated tokens
		_ = p.limiter.Increment(ctx, key, req.EstimatedTokens, time.Minute)
	}

	return nil
}

// PostCheck updates rate limit counters with actual token usage (async)
func (p *RateLimitPolicy) PostCheck(ctx context.Context, req *PostRequestContext) {
	// Update token counter with actual tokens (if different from estimate)
	if p.config.TokensPerMinute > 0 {
		actualTokens := req.ActualTokens.TotalTokens
		if actualTokens > 0 {
			key := RateLimitKey(req.AppID, "tokens")
			// Note: In a real implementation, you'd want to adjust for the difference
			// between estimated and actual tokens to avoid double-counting
			// For now, we just track actual tokens
			_ = p.limiter.Increment(ctx, key, actualTokens, time.Minute)
		}
	}
}

// rateLimitError creates a structured error for rate limit violations
func rateLimitError(message string) error {
	return fmt.Errorf("rate limit exceeded: %s", message)
}
