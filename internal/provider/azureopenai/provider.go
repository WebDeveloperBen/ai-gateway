// Package azureopenai implements the azure open ai provider requirements and handles all necessary mapping logic
package azureopenai

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/insurgence-ai/llm-gateway/internal/config"
	"github.com/insurgence-ai/llm-gateway/internal/provider"
)

type Entry struct {
	BaseURL    string // resource host or absolute URL
	Deployment string // AOAI deployment name
	APIVer     string // e.g., "2024-07-01-preview"
}

type Adapter struct {
	Global    map[string]Entry            // model -> entry
	ByTenant  map[string]map[string]Entry // tenant -> model -> entry
	Default   *Entry                      // optional
	APIKeyEnv string                      // e.g., "AOAI_API_KEY"
	APIKeyFor func(tenant string) string  // injector for tests/overrides
}

func New() *Adapter {
	return &Adapter{
		Global:    map[string]Entry{},
		ByTenant:  map[string]map[string]Entry{},
		APIKeyEnv: "AOAI_API_KEY",
	}
}

func (a *Adapter) Prefix() string { return "/azure/openai" }

func (a *Adapter) Rewrite(req *http.Request, suffix string, info provider.ReqInfo) error {
	// 1) Resolve routing entry
	modelKey, ok := provider.ModelOrDefault(
		info.Model,
		func(m string) bool { _, ok := a.peek(info.Tenant, m); return ok },
		a.singleGlobalKey,
		a.Default != nil,
		"__default__",
	)

	var ent Entry
	switch {
	case ok && modelKey == "__default__":
		ent = *a.Default
	case ok:
		ent, _ = a.peek(info.Tenant, modelKey)
	default:
		return fmt.Errorf("unknown model %q and no default route", info.Model)
	}

	// sanity
	if ent.BaseURL == "" || ent.Deployment == "" || ent.APIVer == "" {
		return fmt.Errorf("aoai route incomplete")
	}

	// 2) Base normalization (AOAI well-known suffix)
	base, err := provider.EnsureAbsoluteBase(ent.BaseURL, "openai.azure.com")
	if err != nil {
		return err
	}

	// 3) Build upstream URL: /openai/deployments/{deployment} + trimmed OpenAI path
	trimmed := strings.TrimPrefix(suffix, "/v1")
	q := provider.CopyQuery(req)
	if q == nil {
		q = url.Values{}
	}
	q.Set("api-version", ent.APIVer)

	u, err := provider.JoinURL(base, []string{"/openai/deployments", ent.Deployment, trimmed}, q)
	if err != nil {
		return err
	}
	provider.SetUpstreamURL(req, u)

	// 4) Auth: strip caller, then set AOAI key if present
	provider.StripCallerAuth(req.Header)
	key := ""
	if a.APIKeyFor != nil {
		key = a.APIKeyFor(info.Tenant)
	}
	if key == "" {
		key = config.Envs.AzureOpenAiAPIKey
	}
	provider.SetAPIKey(req.Header, "api-key", key)

	return nil
}

// ---- internals ----

func (a *Adapter) peek(tenant, model string) (Entry, bool) {
	model = strings.ToLower(strings.TrimSpace(model))
	if tenant != "" {
		if tmap, ok := a.ByTenant[tenant]; ok {
			if e, ok := tmap[model]; ok {
				return e, true
			}
		}
	}
	if e, ok := a.Global[model]; ok {
		return e, true
	}
	return Entry{}, false
}

func (a *Adapter) singleGlobalKey() (string, bool) {
	if len(a.Global) != 1 {
		return "", false
	}
	for k := range a.Global {
		return k, true
	}
	return "", false
}
