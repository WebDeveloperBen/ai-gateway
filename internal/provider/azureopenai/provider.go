// Package azureopenai implements the azure open ai provider requirements and handles all necessary mapping logic
package azureopenai

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/insurgence-ai/llm-gateway/internal/gateway/loadbalancing"
	"github.com/insurgence-ai/llm-gateway/internal/model"
	"github.com/insurgence-ai/llm-gateway/internal/provider"
)

type Entry struct {
	BaseURL    string // resource host or absolute URL
	Deployment string // AOAI deployment name
	APIVer     string // e.g., "2024-07-01-preview"
	SecretRef  string
}

type Adapter struct {
	Instances map[string][]Entry // model -> []Entry
	Selector  loadbalancing.InstanceSelector
	Keys      provider.KeySource
}

func New(selector loadbalancing.InstanceSelector) *Adapter {
	return &Adapter{
		Instances: map[string][]Entry{},
		Selector:  selector,
		Keys:      provider.KeySource{EnvVar: "AZURE_OPENAI_API_KEY"},
	}
}

func (a *Adapter) Prefix() string { return provider.AzureOpenAIPrefix }

func (a *Adapter) Rewrite(req *http.Request, suffix string, info provider.ReqInfo) error {
	modelKey := strings.ToLower(strings.TrimSpace(info.Model))
	instances, ok := a.Instances[modelKey]
	if !ok || len(instances) == 0 {
		return fmt.Errorf("no deployments found for model %q", info.Model)
	}
	// Select instance
	var ids []string
	for _, ent := range instances {
		ids = append(ids, ent.Deployment)
	}
	chosen := a.Selector.Select(ids, info.Model)
	var ent Entry
	for _, inst := range instances {
		if inst.Deployment == chosen {
			ent = inst
			break
		}
	}
	if ent.BaseURL == "" || ent.Deployment == "" || ent.APIVer == "" {
		return fmt.Errorf("selected deployment incomplete")
	}

	base, err := provider.EnsureAbsoluteBase(ent.BaseURL, "openai.azure.com")
	if err != nil {
		return err
	}

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

	provider.StripCallerAuth(req.Header)
	key := a.Keys.Resolve(info.Tenant, "AZURE_OPENAI_API_KEY")
	if ent.SecretRef != "" {
		if v := os.Getenv(ent.SecretRef); v != "" {
			key = v
		}
	}
	provider.SetAPIKey(req.Header, "api-key", key)
	return nil
}

// BuildProvider builds and returns a provider.Adapter configured with all models/deployments.
// It accepts an instance selector for loadbalancing, defaulting to round robin if nil.
func BuildProvider(deployments []model.ModelDeployment, selector loadbalancing.InstanceSelector) *Adapter {
	if selector == nil {
		selector = loadbalancing.NewRoundRobinSelector()
	}
	adapter := New(selector)
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
		adapter.Instances[md.Model] = append(adapter.Instances[md.Model], ent)
	}
	if len(adapter.Instances) == 0 {
		return nil
	}
	return adapter
}
