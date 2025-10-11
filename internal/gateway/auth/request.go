package auth

import "context"

// ParsedRequest contains pre-parsed LLM request data to avoid multiple JSON unmarshals
type ParsedRequest struct {
	Model           string
	EstimatedTokens int
	RequestSize     int

	// Raw messages for advanced policies that need full content
	Messages []Message
	Prompt   string
}

// Message represents a chat message
type Message struct {
	Role    string
	Content string
}

// Context helpers for parsed request
type contextKeyParsedRequest struct{}

func WithParsedRequest(ctx context.Context, req *ParsedRequest) context.Context {
	return context.WithValue(ctx, contextKeyParsedRequest{}, req)
}

func GetParsedRequest(ctx context.Context) *ParsedRequest {
	if req, ok := ctx.Value(contextKeyParsedRequest{}).(*ParsedRequest); ok {
		return req
	}
	return nil
}

// Convenience getters
func GetModel(ctx context.Context) string {
	if req := GetParsedRequest(ctx); req != nil {
		return req.Model
	}
	return ""
}

func GetEstimatedTokens(ctx context.Context) int {
	if req := GetParsedRequest(ctx); req != nil {
		return req.EstimatedTokens
	}
	return 0
}

func GetRequestSize(ctx context.Context) int {
	if req := GetParsedRequest(ctx); req != nil {
		return req.RequestSize
	}
	return 0
}
