package testkit

import (
	"os"
	"strings"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/config"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/auth"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/loadbalancing"
	"github.com/WebDeveloperBen/ai-gateway/internal/provider"
	"github.com/WebDeveloperBen/ai-gateway/internal/provider/azureopenai"
)

type AOAITest struct {
	Adapter       *azureopenai.Adapter
	Model         string
	BasePath      string // where you mount it in your API (used by tests)
	Authenticator auth.KeyAuthenticator
}

// AOAIOption Functional options for easy overrides in specific tests
type AOAIOption func(*aoaiOpts)

type aoaiOpts struct {
	Model      string
	BaseURL    string
	Deployment string
	APIVer     string
	KeyEnv     string // env var name storing the key
}

type AOAIUnitOption func(*aoaiUnitOpts)

type aoaiUnitOpts struct {
	Model      string
	BaseURL    string
	Deployment string
	APIVer     string
	Key        string // if set, use direct key
	KeyEnv     string // else, set EnvVar to read from env in test
}

// NewAOAIE2E loads env, refreshes config, validates required vars, and returns a ready adapter.
// Defaults come from environment so you can configure once in .env:
//
//	AZURE_OPENAI_ENDPOINT      = https://dev-insurgence-openai.openai.azure.com
//	AOAI_DEPLOYMENT    = dev-openai-gpt4-1
//	AOAI_API_VERSION   = 2024-07-01-preview
//	AZURE_OPENAI_API_KEY = <secret>        (or choose your own via AOAI_KEY_ENV)
//	AOAI_MODEL         = gpt-4o            (optional)
//
// You can override any of these via AOAIOption.
func NewAOAIE2E(t *testing.T, opts ...AOAIOption) *AOAITest {
	t.Helper()

	LoadDotenvFromRepoRoot(t)

	// Make sure config.Envs is fresh in case tests rely on it elsewhere.
	config.Reload()

	o := aoaiOpts{
		Model:      getenvDefault("AOAI_MODEL", "dev-openai-gpt4-1"),
		BaseURL:    getenvDefault("AZURE_OPENAI_ENDPOINT", "https://dev-insurgence-openai.openai.azure.com"),
		Deployment: getenvDefault("AOAI_DEPLOYMENT", "dev-insurgence-openai"),
		APIVer:     getenvDefault("AOAI_API_VERSION", "2024-07-01-preview"),
		KeyEnv:     getenvDefault("AOAI_KEY_ENV", "AZURE_OPENAI_API_KEY"),
	}
	for _, f := range opts {
		f(&o)
	}

	// Fail/skip loudly with actionable messages.
	if o.BaseURL == "" || o.Deployment == "" {
		t.Skipf("AZURE_OPENAI_ENDPOINT and/or AOAI_DEPLOYMENT missing; set them in your .env or pass WithBaseURL/WithDeployment")
	}
	key := os.Getenv(o.KeyEnv)
	if key == "" && config.Envs.AzureOpenAiAPIKey == "" {
		t.Skipf("%s missing; set it or populate config.Envs.AzureOpenAiAPIKey via .env", o.KeyEnv)
	}

	ad := azureopenai.New(loadbalancing.NewRoundRobinSelector())
	// Prefer explicit env var; otherwise fall back to config.Envs
	if key != "" {
		ad.Keys = provider.KeySource{EnvVar: o.KeyEnv}
	} else {
		ad.Keys = provider.KeySource{ForTenant: func(string) string { return config.Envs.AzureOpenAiAPIKey }}
	}

	// Fix for E2E: set mapping for the tested model
	ad.Instances[strings.ToLower(o.Model)] = []azureopenai.Entry{
		{
			BaseURL:    o.BaseURL,
			Deployment: o.Deployment,
			APIVer:     o.APIVer,
			SecretRef:  "",
		},
	}

	return &AOAITest{
		Adapter:       ad,
		Model:         o.Model,
		BasePath:      "/azure/openai",
		Authenticator: &auth.NoopAuthenticator{},
	}
}

func NewAOAIUnit(t *testing.T, opts ...AOAIUnitOption) *AOAITest {
	t.Helper()
	// Safe defaults for unit tests
	o := aoaiUnitOpts{
		Model:      "gpt-4o",
		BaseURL:    "https://example.openai.azure.com",
		Deployment: "any-deploy",
		APIVer:     "2024-07-01-preview",
	}

	for _, f := range opts {
		f(&o)
	}

	ad := azureopenai.New(loadbalancing.NewRoundRobinSelector())

	if o.Key != "" {
		ad.Keys = provider.KeySource{ForTenant: func(string) string { return o.Key }}
	} else if o.KeyEnv != "" {
		ad.Keys = provider.KeySource{EnvVar: o.KeyEnv}
	}

	ad.Instances[strings.ToLower(o.Model)] = []azureopenai.Entry{
		{
			BaseURL:    o.BaseURL,
			Deployment: o.Deployment,
			APIVer:     o.APIVer,
			SecretRef:  "",
		},
	}

	return &AOAITest{Adapter: ad, Model: o.Model, BasePath: "/azure/openai", Authenticator: &auth.NoopAuthenticator{}}
}

// WithModel and the others are conviencance E2E test defaults
func WithModel(m string) AOAIOption      { return func(o *aoaiOpts) { o.Model = m } }
func WithBaseURL(u string) AOAIOption    { return func(o *aoaiOpts) { o.BaseURL = u } }
func WithDeployment(d string) AOAIOption { return func(o *aoaiOpts) { o.Deployment = d } }
func WithAPIVer(v string) AOAIOption     { return func(o *aoaiOpts) { o.APIVer = v } }
func WithKeyEnv(name string) AOAIOption  { return func(o *aoaiOpts) { o.KeyEnv = name } }

// AOAIUnitWithMapping and the others below are convienance unit test defaults
func AOAIUnitWithMapping(model, baseURL, deployment, apiVer string) AOAIUnitOption {
	return func(o *aoaiUnitOpts) {
		o.Model, o.BaseURL, o.Deployment, o.APIVer = model, baseURL, deployment, apiVer
	}
}
func AOAIUnitWithKey(key string) AOAIUnitOption     { return func(o *aoaiUnitOpts) { o.Key = key } }
func AOAIUnitWithKeyEnv(name string) AOAIUnitOption { return func(o *aoaiUnitOpts) { o.KeyEnv = name } }

func getenvDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
