package azureopenai

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/insurgence-ai/llm-gateway/internal/gateway"
)

type Entry struct {
	BaseURL    string
	Deployment string
	APIVer     string
}

type Router struct {
	Global    map[string]Entry            // model (lowercase) -> entry
	ByTenant  map[string]map[string]Entry // tenant -> model -> entry
	Default   *Entry
	APIKeyEnv string // e.g. "AOAI_API_KEY" for POC
}

func New() *Router {
	return &Router{
		Global:    map[string]Entry{},
		ByTenant:  map[string]map[string]Entry{},
		APIKeyEnv: "AOAI_API_KEY",
	}
}

func (r *Router) Route(info gateway.RequestInfo) (string, gateway.MutateFunc, error) {
	model := strings.ToLower(strings.TrimSpace(info.Model))

	var ent Entry
	ok := false
	if info.Tenant != "" {
		if tmap, okT := r.ByTenant[info.Tenant]; okT {
			if v, okM := tmap[model]; okM {
				ent, ok = v, true
			}
		}
	}
	if !ok {
		if v, okG := r.Global[model]; okG {
			ent, ok = v, true
		}
	}
	if !ok {
		if r.Default == nil {
			return "", nil, fmt.Errorf("unknown model %q and no default route", info.Model)
		}
		ent, ok = *r.Default, true
	}

	if ent.BaseURL == "" || ent.Deployment == "" || ent.APIVer == "" {
		return "", nil, fmt.Errorf("aoai route incomplete")
	}
	if !strings.HasPrefix(info.Path, "/v1/") {
		return "", nil, fmt.Errorf("unsupported path %q for AOAI", info.Path)
	}

	rest := strings.TrimPrefix(info.Path, "/v1") // keep /<resource> suffix
	path := "/openai/deployments/" + url.PathEscape(ent.Deployment) + rest
	target := ent.BaseURL + path
	if info.Query != "" {
		target += "?" + info.Query + "&api-version=" + url.QueryEscape(ent.APIVer)
	} else {
		target += "?api-version=" + url.QueryEscape(ent.APIVer)
	}

	mutate := func(hctx huma.Context, req *http.Request, _ []byte) error {
		req.Header.Del("Authorization") // strip caller auth
		if key := os.Getenv(r.APIKeyEnv); key != "" {
			req.Header.Set("api-key", key) // POC: key auth
		}
		return nil
	}
	return target, mutate, nil
}
