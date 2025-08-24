package testkit

import (
	"os"
	"testing"

	"github.com/insurgence-ai/llm-gateway/internal/config"
	"github.com/insurgence-ai/llm-gateway/internal/provider"
	"github.com/insurgence-ai/llm-gateway/internal/provider/azureopenai"
)

type AOAITest struct {
	Adapter  *azureopenai.Adapter
	Model    string
	BasePath string // where you mount it in your API (used by tests)
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
//	AOAI_BASE_URL      = https://dev-insurgence-openai.openai.azure.com
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
		BaseURL:    os.Getenv("AOAI_BASE_URL"),
		Deployment: os.Getenv("AOAI_DEPLOYMENT"),
		APIVer:     getenvDefault("AOAI_API_VERSION", "2024-07-01-preview"),
		KeyEnv:     getenvDefault("AOAI_KEY_ENV", "AZURE_OPENAI_API_KEY"),
	}
	for _, f := range opts {
		f(&o)
	}

	// Fail/skip loudly with actionable messages.
	if o.BaseURL == "" || o.Deployment == "" {
		t.Skipf("AOAI_BASE_URL and/or AOAI_DEPLOYMENT missing; set them in your .env or pass WithBaseURL/WithDeployment")
	}
	key := os.Getenv(o.KeyEnv)
	if key == "" && config.Envs.AzureOpenAiAPIKey == "" {
		t.Skipf("%s missing; set it or populate config.Envs.AzureOpenAiAPIKey via .env", o.KeyEnv)
	}

	ad := azureopenai.New()
	ad.Global[o.Model] = azureopenai.Entry{
		BaseURL:    o.BaseURL,
		Deployment: o.Deployment,
		APIVer:     o.APIVer,
	}
	// Prefer explicit env var; otherwise fall back to config.Envs
	if key != "" {
		ad.Keys = provider.KeySource{EnvVar: o.KeyEnv}
	} else {
		ad.Keys = provider.KeySource{ForTenant: func(string) string { return config.Envs.AzureOpenAiAPIKey }}
	}

	return &AOAITest{
		Adapter:  ad,
		Model:    o.Model,
		BasePath: "/azure/openai",
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

	ad := azureopenai.New()
	ad.Global[o.Model] = azureopenai.Entry{
		BaseURL:    o.BaseURL,
		Deployment: o.Deployment,
		APIVer:     o.APIVer,
	}

	if o.Key != "" {
		ad.Keys = provider.KeySource{ForTenant: func(string) string { return o.Key }}
	} else if o.KeyEnv != "" {
		ad.Keys = provider.KeySource{EnvVar: o.KeyEnv}
	} // else leave default EnvVar="AOAI_API_KEY" (rare in unit tests)

	return &AOAITest{Adapter: ad, Model: o.Model, BasePath: "/azure/openai"}
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
