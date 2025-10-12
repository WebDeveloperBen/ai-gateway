// Package middleware provides HTTP middleware for the gateway
package middleware

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/auth"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/policies"
	"github.com/WebDeveloperBen/ai-gateway/internal/logger"
)

// PolicyEnforcer provides policy enforcement for requests
type PolicyEnforcer struct {
	engine *policies.Engine
}

// NewPolicyEnforcer creates a new policy enforcer
func NewPolicyEnforcer(engine *policies.Engine) *PolicyEnforcer {
	return &PolicyEnforcer{
		engine: engine,
	}
}

// Middleware returns a RoundTripper middleware that enforces policies
func (pe *PolicyEnforcer) Middleware(next http.RoundTripper) http.RoundTripper {
	return roundTripFunc(func(r *http.Request) (*http.Response, error) {
		ctx := r.Context()

		// Fast path: Skip policy enforcement if app context is missing
		// This indicates the request is not an LLM request (auth middleware not applied)
		appID := auth.GetAppID(ctx)
		if appID == "" {
			return next.RoundTrip(r)
		}

		// Load policies for this app
		policyList, err := pe.engine.LoadPolicies(ctx, appID)
		if err != nil {
			logger.GetLogger(ctx).Error().
				Err(err).
				Str("app_id", appID).
				Msg("Failed to load policies for enforcement")
			return deny(500, "failed to load policies"), nil
		}

		// Store policies in context for usage recording middleware
		ctx = auth.WithPolicies(ctx, policyList)
		r = r.WithContext(ctx)

		// Get pre-parsed request data from context
		// No more JSON unmarshaling, no more body access!
		parsedReq := auth.GetParsedRequest(ctx)
		if parsedReq == nil {
			// This shouldn't happen if RequestBuffer middleware ran
			logger.GetLogger(ctx).Error().
				Str("app_id", appID).
				Msg("Missing parsed request data")
			return deny(500, "internal error"), nil
		}

		// Build pre-request context using parsed data
		preCtx := &policies.PreRequestContext{
			Request:          r,
			OrgID:            auth.GetOrgID(ctx),
			AppID:            appID,
			APIKeyID:         auth.GetKeyID(ctx),
			Model:            parsedReq.Model,
			EstimatedTokens:  parsedReq.EstimatedTokens,
			RequestSizeBytes: parsedReq.RequestSize,
			Body:             nil, // No longer needed - policies use parsed data
		}

		// Run pre-checks (blocking)
		for _, policy := range policyList {
			if err := policy.PreCheck(ctx, preCtx); err != nil {
				logger.GetLogger(ctx).Warn().
					Err(err).
					Str("app_id", appID).
					Str("org_id", auth.GetOrgID(ctx)).
					Str("policy_type", string(policy.Type())).
					Str("model", parsedReq.Model).
					Int("estimated_tokens", parsedReq.EstimatedTokens).
					Msg("Policy check failed")
				return deny(429, "policy violation"), nil
			}
		}

		// Continue with request
		return next.RoundTrip(r)
	})
}

// roundTripFunc is a type adapter for http.RoundTripper
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// deny creates an error response
func deny(code int, msg string) *http.Response {
	return &http.Response{
		StatusCode: code,
		Status:     fmt.Sprintf("%d %s", code, http.StatusText(code)),
		Header:     http.Header{"Content-Type": []string{"application/problem+json"}},
		Body:       io.NopCloser(strings.NewReader(fmt.Sprintf(`{"title":"%s","status":%d}`, msg, code))),
	}
}
