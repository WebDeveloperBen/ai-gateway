// Package openai implements the OpenAI (vendor) provider. It is OpenAI-API
// compatible: the downstream path (/v1/...) is forwarded upstream unchanged.
package openai

import (
	"net/http"
	"strings"

	"github.com/insurgence-ai/llm-gateway/internal/model"
	"github.com/insurgence-ai/llm-gateway/internal/provider"
)

type Adapter struct {
	BaseURL    string
	Keys       provider.KeySource
	ModelAlias map[string]string
	OrgFor     func(tenant string) string
}

func New() *Adapter {
	return &Adapter{
		BaseURL:    "api.openai.com",
		Keys:       provider.KeySource{EnvVar: "OPENAI_API_KEY"},
		ModelAlias: map[string]string{},
	}
}

// Prefix determines where you mount this provider under your API (for docs/tests).
func (a *Adapter) Prefix() string { return "/openai" }

// BuildProvider builds and returns a provider.Adapter configured with all models/deployments.
// Used to dynamically instantiate tenant/model/provider adapters from registry/model config at runtime.
// Accepts all deployments, must filter & populate its own adapter-specific mapping.
func BuildProvider(deployments []model.ModelDeployment) *Adapter {
	has := false
	adapter := New()
	for _, md := range deployments {
		if md.Provider != "openai" {
			continue
		}
		adapter.ModelAlias[md.Model] = md.Meta["Alias"]
		has = true
	}
	if !has {
		return nil
	}
	return adapter
}

// Rewrite turns an OpenAI-compatible downstream request into an OpenAI upstream request.
// Unlike Azure, we keep the /v1/... suffix and forward it. Optionally, we rewrite the
// "model" field if ModelAlias has an entry for info.Model.
func (a *Adapter) Rewrite(req *http.Request, suffix string, info provider.ReqInfo) error {
	base, _ := provider.EnsureAbsoluteBase(a.BaseURL, "api.openai.com")
	u, _ := provider.JoinURL(base, []string{suffix}, provider.CopyQuery(req))
	provider.SetUpstreamURL(req, u)

	provider.StripCallerAuth(req.Header)
	if key := a.Keys.Resolve(info.Tenant, "OPENAI_API_KEY"); key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}
	if a.OrgFor != nil {
		if org := strings.TrimSpace(a.OrgFor(info.Tenant)); org != "" {
			req.Header.Set("OpenAI-Organization", org)
		}
	}
	if alias, ok := a.ModelAlias[strings.ToLower(strings.TrimSpace(info.Model))]; ok && alias != "" && alias != info.Model {
		_ = provider.RewriteJSONModel(req, alias) // best-effort
	}
	return nil
}
