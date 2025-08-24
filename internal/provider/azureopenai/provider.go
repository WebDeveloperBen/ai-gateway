// Package azureopenai implements the azure open ai provider requirements and handles all necessary mapping logic
package azureopenai

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/insurgence-ai/llm-gateway/internal/model/models"
	"github.com/insurgence-ai/llm-gateway/internal/provider"
)

type Entry struct {
	BaseURL    string // resource host or absolute URL
	Deployment string // AOAI deployment name
	APIVer     string // e.g., "2024-07-01-preview"
	SecretRef  string
}

type Adapter struct {
	Global   map[string]Entry            // model -> entry
	ByTenant map[string]map[string]Entry // tenant -> model -> entry
	Default  *Entry                      // optional
	Keys     provider.KeySource
}

func New() *Adapter {
	return &Adapter{
		Global:   map[string]Entry{},
		ByTenant: map[string]map[string]Entry{},
		Keys:     provider.KeySource{EnvVar: "AOAI_API_KEY"},
	}
}

func (a *Adapter) Prefix() string { return "/azure/openai" }

func (a *Adapter) Rewrite(req *http.Request, suffix string, info provider.ReqInfo) error {
	// Resolve routing entry
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
		fmt.Printf("[AzureOpenAI] model lookup failed. Requested: %q; known: %v; byTenant: %v; tenant: %q\n", info.Model, a.Global, a.ByTenant, info.Tenant)
		return fmt.Errorf("unknown model %q and no default route", info.Model)
	}

	// sanity
	if ent.BaseURL == "" || ent.Deployment == "" || ent.APIVer == "" {
		return fmt.Errorf("aoai route incomplete")
	}

	// Base normalization (AOAI well-known suffix)
	base, err := provider.EnsureAbsoluteBase(ent.BaseURL, "openai.azure.com")
	if err != nil {
		return err
	}

	// Build upstream URL: /openai/deployments/{deployment} + trimmed OpenAI path
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

	// strip caller, then set AOAI key if present
	provider.StripCallerAuth(req.Header)

	key := a.Keys.Resolve(info.Tenant, "AOAI_API_KEY")
	// TODO: preference managed identity lookup over keys
	// TODO: replace this with a kv lookup instead in future
	// prefer per-resource secret if it's present
	if ent.SecretRef != "" {
		if v := os.Getenv(ent.SecretRef); v != "" {
			key = v
		}
	}

	provider.SetAPIKey(req.Header, "api-key", key)

	return nil
}

// BuildProvider builds and returns a provider.Adapter configured with all models/deployments.
// Used to dynamically instantiate tenant/model/provider adapters from registry/model config at runtime.
// Accepts all deployments, must filter & populate its own adapter-specific mapping.
func BuildProvider(deployments []models.ModelDeployment) *Adapter {
	has := false
	adapter := New()
	for _, md := range deployments {
		if md.Provider != "azure" {
			continue
		}
		ent := Entry{
			BaseURL:    md.Meta["BaseURL"],
			Deployment: md.Deployment,
			APIVer:     md.Meta["APIVer"],
			SecretRef:  md.Meta["SecretRef"],
		}
		if md.Tenant != "" {
			if adapter.ByTenant[md.Tenant] == nil {
				adapter.ByTenant[md.Tenant] = map[string]Entry{}
			}
			adapter.ByTenant[md.Tenant][md.Model] = ent
		} else {
			adapter.Global[md.Model] = ent
		}
		has = true
	}
	if !has {
		return nil
	}
	return adapter
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
