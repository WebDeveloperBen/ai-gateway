// Package policies provides a flexible policy engine for request validation
// and usage tracking. It supports both predefined policy
// types and custom CEL (Common Expression Language) policies.
package policies

import (
	"context"
	"net/http"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

// Policy defines the interface that all policies must implement
type Policy interface {
	// Type returns the policy type
	Type() model.PolicyType

	// PreCheck is called before the request is sent upstream (blocking)
	// Returns an error if the policy check fails
	PreCheck(ctx context.Context, req *PreRequestContext) error

	// PostCheck is called after the response is received (async, non-blocking)
	PostCheck(ctx context.Context, req *PostRequestContext)
}

// PreRequestContext contains information available before sending the request
type PreRequestContext struct {
	Request          *http.Request
	OrgID            string
	AppID            string
	APIKeyID         string
	Model            string
	EstimatedTokens  int
	RequestSizeBytes int
	Body             []byte // Captured request body
}

// PostRequestContext contains information after receiving the response
type PostRequestContext struct {
	Request           *http.Request
	Response          *http.Response
	OrgID             string
	AppID             string
	APIKeyID          string
	Provider          string
	ModelName         string
	ActualTokens      model.TokenUsage
	RequestSizeBytes  int
	ResponseSizeBytes int
	LatencyMs         int64
}
