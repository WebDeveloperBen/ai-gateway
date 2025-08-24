// Package azureopenai implements the azure open ai provider requirements and handles all necessary mapping logic
package azureopenai

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/insurgence-ai/llm-gateway/internal/provider"
)

type Entry struct {
	BaseURL    string // resource host or absolute URL
	Deployment string // AOAI deployment name
	APIVer     string // e.g., "2024-07-01-preview"
}

type Adapter struct {
	Global   map[string]Entry            // model (lowercase) -> entry
	ByTenant map[string]map[string]Entry // tenant -> model -> entry
	Default  *Entry

	APIKeyEnv string                     // e.g., "AOAI_API_KEY"
	APIKeyFor func(tenant string) string // optional injector for tests
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
	model := strings.ToLower(strings.TrimSpace(info.Model))

	ent, ok := a.resolve(info.Tenant, model)
	if !ok {
		// NEW: if exactly one global entry exists, use it as an implicit default
		if e, ok1 := a.singleGlobal(); ok1 {
			ent, ok = e, true
		}
	}
	if !ok {
		return fmt.Errorf("unknown model %q and no default route", info.Model)
	}
	if ent.BaseURL == "" || ent.Deployment == "" || ent.APIVer == "" {
		return fmt.Errorf("aoai route incomplete")
	}

	base, err := normalizeBase(ent.BaseURL)
	if err != nil {
		return err
	}

	// Build AOAI URL: /openai/deployments/{deployment} + /v1/... tail (without /v1)
	u, _ := url.Parse(base)
	u.Path = path.Join(u.Path, "/openai/deployments", url.PathEscape(ent.Deployment), strings.TrimPrefix(suffix, "/v1"))
	q := req.URL.Query() // preserve client query if present
	q.Set("api-version", ent.APIVer)
	u.RawQuery = q.Encode()

	req.URL.Scheme, req.URL.Host, req.URL.Path, req.URL.RawQuery = u.Scheme, u.Host, u.Path, u.RawQuery

	// Auth: api-key (Managed Identity can be added later by a RoundTripper)
	req.Header.Del("Authorization")
	key := ""
	if a.APIKeyFor != nil {
		key = a.APIKeyFor(info.Tenant)
	}
	if key == "" {
		key = os.Getenv(a.APIKeyEnv)
	}
	if key != "" {
		req.Header.Set("api-key", key)
	}
	return nil
}

func (a *Adapter) resolve(tenant, model string) (Entry, bool) {
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
	if a.Default != nil {
		return *a.Default, true
	}
	return Entry{}, false
}

func normalizeBase(base string) (string, error) {
	b := strings.TrimRight(strings.TrimSpace(base), "/")
	if b == "" {
		return "", fmt.Errorf("empty BaseURL")
	}
	if strings.Contains(b, "://") {
		return b, nil
	}
	if strings.Contains(b, ".") {
		return "https://" + b, nil
	}
	return "https://" + b + ".openai.azure.com", nil
}

func (a *Adapter) singleGlobal() (Entry, bool) {
	if len(a.Global) != 1 {
		return Entry{}, false
	}
	for _, e := range a.Global {
		return e, true
	}
	return Entry{}, false
}
