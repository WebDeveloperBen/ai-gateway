package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/auth"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/tokens"
)

// RequestBuffer buffers and parses the request body once
type RequestBuffer struct {
	estimator *tokens.Estimator
}

// NewRequestBuffer creates a new request buffer middleware
func NewRequestBuffer() *RequestBuffer {
	return &RequestBuffer{
		estimator: tokens.NewEstimator(),
	}
}

// Middleware returns a RoundTripper middleware that buffers and parses the request body
func (rb *RequestBuffer) Middleware(next http.RoundTripper) http.RoundTripper {
	return roundTripFunc(func(r *http.Request) (*http.Response, error) {
		// Read request body once
		var bodyBytes []byte
		if r.Body != nil {
			var err error
			bodyBytes, err = io.ReadAll(r.Body)
			if err != nil {
				return deny(400, "failed to read request body"), nil
			}
			r.Body.Close()
		}

		// Replace body with a fresh reader for upstream
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		// Parse request ONCE and extract all needed data
		parsed := rb.parseRequest(bodyBytes)

		// Store parsed request in context (small struct, not full body)
		ctx := auth.WithParsedRequest(r.Context(), parsed)
		r = r.WithContext(ctx)

		return next.RoundTrip(r)
	})
}

// parseRequest parses the LLM request body and extracts all needed information
func (rb *RequestBuffer) parseRequest(body []byte) *auth.ParsedRequest {
	parsed := &auth.ParsedRequest{
		RequestSize: len(body),
	}

	if len(body) == 0 {
		return parsed
	}

	// Parse JSON once into full structure
	var req struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages,omitempty"`
		Prompt string `json:"prompt,omitempty"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		// If parsing fails, return partial data
		return parsed
	}

	// Extract model
	parsed.Model = req.Model

	// Convert messages
	if len(req.Messages) > 0 {
		parsed.Messages = make([]auth.Message, len(req.Messages))
		for i, msg := range req.Messages {
			parsed.Messages[i] = auth.Message{
				Role:    msg.Role,
				Content: msg.Content,
			}
		}
	}

	// Extract prompt (for completion endpoints)
	parsed.Prompt = req.Prompt

	// Estimate tokens using the parsed data
	estimatedTokens, err := rb.estimator.EstimateRequest(parsed.Model, body)
	if err == nil {
		parsed.EstimatedTokens = estimatedTokens
	}

	return parsed
}
