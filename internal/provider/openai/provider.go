// Package openai implements the OpenAI (vendor) provider. It is OpenAI-API
// compatible: the downstream path (/v1/...) is forwarded upstream unchanged.
package openai

import (
	"net/http"
	"strings"

	"github.com/insurgence-ai/llm-gateway/internal/loadbalancing"
	"github.com/insurgence-ai/llm-gateway/internal/model"
	"github.com/insurgence-ai/llm-gateway/internal/provider"
)

type Adapter struct {
	BaseURL    string
	Keys       provider.KeySource
	Instances  map[string][]string // model -> []deployment IDs
	Selector   loadbalancing.InstanceSelector
	ModelAlias map[string]string
	OrgFor     func(tenant string) string
}

func New(selector loadbalancing.InstanceSelector) *Adapter {
	return &Adapter{
		BaseURL:    "api.openai.com",
		Keys:       provider.KeySource{EnvVar: "OPENAI_API_KEY"},
		Instances:  map[string][]string{},
		Selector:   selector,
		ModelAlias: map[string]string{},
	}
}

func (a *Adapter) Prefix() string { return "/openai" }

func (a *Adapter) Rewrite(req *http.Request, suffix string, info provider.ReqInfo) error {
	modelKey := strings.ToLower(strings.TrimSpace(info.Model))
	instances := a.Instances[modelKey]
	if len(instances) == 0 {
		return nil // fallback, no deployments
	}
	// Use chosen deployment ID (string) to set header or param as you need, for now just a placeholder
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
		_ = provider.RewriteJSONModel(req, alias)
	}
	return nil
}

func BuildProvider(deployments []model.ModelDeployment, selector loadbalancing.InstanceSelector) *Adapter {
	if selector == nil {
		selector = loadbalancing.NewRoundRobinSelector()
	}
	adapter := New(selector)
	for _, md := range deployments {
		if md.Provider != "openai" {
			continue
		}
		adapter.Instances[md.Model] = append(adapter.Instances[md.Model], md.Deployment)
		adapter.ModelAlias[md.Model] = md.Meta["Alias"]
	}
	if len(adapter.Instances) == 0 {
		return nil
	}
	return adapter
}
