package policies

import (
	"context"
	"fmt"

	"github.com/WebDeveloperBen/ai-gateway/internal/logger"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

// RequestSizePolicy enforces maximum request body size
type RequestSizePolicy struct {
	config model.RequestSizeConfig
}

// NewRequestSizePolicy creates a new request size policy
func NewRequestSizePolicy(config model.RequestSizeConfig) *RequestSizePolicy {
	return &RequestSizePolicy{config: config}
}

// Type returns the policy type
func (p *RequestSizePolicy) Type() model.PolicyType {
	return model.PolicyTypeRequestSize
}

// PreCheck checks if the request size exceeds the limit
func (p *RequestSizePolicy) PreCheck(ctx context.Context, req *PreRequestContext) error {
	if p.config.MaxRequestBytes > 0 && req.RequestSizeBytes > p.config.MaxRequestBytes {
		logger.GetLogger(ctx).Warn().
			Int("request_size", req.RequestSizeBytes).
			Int("limit", p.config.MaxRequestBytes).
			Str("app_id", req.AppID).
			Msg("Request size limit exceeded")
		return fmt.Errorf("request size limit exceeded")
	}

	return nil
}

// PostCheck is a no-op for request size
func (p *RequestSizePolicy) PostCheck(ctx context.Context, req *PostRequestContext) {
	// No post-processing needed for request size
}
