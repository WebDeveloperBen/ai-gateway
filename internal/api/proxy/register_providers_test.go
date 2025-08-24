package proxy

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/insurgence-ai/llm-gateway/internal/config"
	"github.com/insurgence-ai/llm-gateway/internal/gateway"
	"github.com/insurgence-ai/llm-gateway/internal/provider/azureopenai"
	"github.com/insurgence-ai/llm-gateway/internal/testkit"
	"github.com/stretchr/testify/require"
)

// --- tests ---

func TestProxyRoutes(t *testing.T) {
	// Azure adapter with a simple global mapping
	aoai := azureopenai.New()
	aoai.Global["gpt-4o"] = azureopenai.Entry{
		BaseURL:    "https://example.openai.azure.com",
		Deployment: "any-deploy",
		APIVer:     "2024-07-01-preview",
	}
	aoai.APIKeyFor = func(_ string) string { return "sekret-key" }

	// Stub upstream: always 200 with {"echo":"ok"}
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		w.WriteHeader(http.StatusOK)
		w.WriteString(`{"echo":"ok"}`)
		return w.Result(), nil
	})

	core := gateway.NewCoreWithAdapters(transport, aoai)

	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		// Mount Azure facade under /azure/openai (and the API group is mounted at /api by testkit)
		RegisterProvider(grp, "/azure/openai", core)
	})

	t.Run("forwards chat completions and returns upstream body", func(t *testing.T) {
		body := map[string]any{"model": "gpt-4o", "messages": []string{"hi"}}
		b, _ := json.Marshal(body)
		resp := api.Post("/api/azure/openai/v1/chat/completions", "Content-Type: application/json", bytes.NewReader(b))
		require.Equal(t, http.StatusOK, resp.Code)
		require.JSONEq(t, `{"echo":"ok"}`, resp.Body.String())
	})
}

func TestUnitProxy_AzureOpenAI(t *testing.T) {
	aoai := azureopenai.New()
	aoai.Global["gpt-4o"] = azureopenai.Entry{
		BaseURL:    "https://dev-insurgence-openai.openai.azure.com",
		Deployment: "dev-openai-gpt4-1",
		APIVer:     "2024-07-01-preview",
	}
	aoai.APIKeyFor = func(_ string) string { return "sekret-key" } // avoid env in unit test

	// Capture the outgoing request the proxy makes.
	called := false

	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		called = true

		require.Equal(t, "dev-insurgence-openai.openai.azure.com", req.URL.Host)
		require.Equal(t, "/openai/deployments/dev-openai-gpt4-1/chat/completions", req.URL.Path)
		require.Equal(t, "api-version=2024-07-01-preview", req.URL.RawQuery)
		require.Equal(t, "sekret-key", req.Header.Get("api-key"))
		require.Empty(t, req.Header.Get("Authorization"))

		w := httptest.NewRecorder()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.WriteString(`{"echo":"ok"}`)
		return w.Result(), nil
	})

	core := gateway.NewCoreWithAdapters(transport, aoai)

	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		RegisterProvider(grp, "/azure/openai", core)
	})

	body := map[string]any{
		"model":    "gpt-4o",
		"messages": []map[string]string{{"role": "user", "content": "Say hello test!"}},
	}
	b, _ := json.Marshal(body)

	resp := api.Post(
		"/api/azure/openai/v1/chat/completions",
		"Content-Type: application/json",
		"Authorization: Bearer client-token", // should be stripped
		bytes.NewReader(b),
	)

	require.True(t, called, "upstream transport was not invoked")
	require.Equal(t, http.StatusOK, resp.Code)
	require.JSONEq(t, `{"echo":"ok"}`, resp.Body.String())
}

func TestE2EProxy_AzureOpenAI(t *testing.T) {
	// 2) Refresh captured config AFTER loading env
	config.Reload()
	cfg := config.Envs

	if cfg.AzureOpenAiAPIKey == "" {
		t.Skip("AZURE_OPENAI_API_KEY missing; skipping E2E (set it in your root .env)")
	}
	t.Logf("Loaded AZURE_OPENAI_API_KEY: **** (len=%d)", len(cfg.AzureOpenAiAPIKey))

	// 3) Wire real AOAI using env-sourced key
	aoai := azureopenai.New()
	aoai.Global["gpt-4o"] = azureopenai.Entry{
		BaseURL:    "https://dev-insurgence-openai.openai.azure.com",
		Deployment: "dev-openai-gpt4-1",
		APIVer:     "2024-07-01-preview",
	}
	aoai.APIKeyFor = func(_ string) string { return cfg.AzureOpenAiAPIKey }

	core := gateway.NewCoreWithAdapters(http.DefaultTransport, aoai)
	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		RegisterProvider(grp, "/azure/openai", core)
	})

	body := map[string]any{
		"model":    "gpt-4o",
		"messages": []map[string]string{{"role": "user", "content": "Say hello test!"}},
	}
	b, _ := json.Marshal(body)
	resp := api.Post("/api/azure/openai/v1/chat/completions", "Content-Type: application/json", bytes.NewReader(b))
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
}

// --- test helpers ---

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
