// Package provider defines the abstraction and implementations for
// upstream LLM providers, including routing, retries, and usage reporting.
package provider

import (
	"net/http"
)

const (
	AzureOpenAIPrefix = "/azure/openai"
	OpenAIPrefix      = "/openai"
)

// ReqInfo is what adapters need to decide routing/rewrite.
// Keep this in provider pkg to avoid gateway<->provider cycles.
type ReqInfo struct {
	Method string
	Path   string // normalized "/v1/..."
	Query  string
	Model  string
	Tenant string
	App    string
}

// Adapter is implemented by each provider package (azureopenai, openai, anthropic, ...).
type Adapter interface {
	// Prefix returns the public API prefix (e.g. "/azure/openai" or "/openai").
	// Requests handled by this adapter must start with this prefix or this prefix + "/".
	Prefix() string

	// Rewrite mutates req in-place to target the upstream provider.
	// suffix is the "/v1/..." portion of the path.
	// info contains normalized, auth-aware request context.
	Rewrite(req *http.Request, suffix string, info ReqInfo) error
}
