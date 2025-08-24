// Package provider defines the abstraction and implementations for
// upstream LLM providers, including routing, retries, and usage reporting.
package provider

import "net/http"

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
	// Prefix is the façade base mounted in your API, e.g. "/azure/openai" or "/openai".
	// Requests must start with this prefix (or this prefix + "/").
	Prefix() string

	// Rewrite mutates req in-place to point at the upstream provider.
	// suffix is the "/v1/..." part of the path.
	Rewrite(req *http.Request, suffix string, info ReqInfo) error
}
