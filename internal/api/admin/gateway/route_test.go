package gateway_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apigw "github.com/WebDeveloperBen/ai-gateway/internal/api/admin/gateway"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/require"
)

func TestProxyRoutes(t *testing.T) {
	fx := testkit.NewAOAIUnit(t,
		testkit.AOAIUnitWithMapping("gpt-4o", "https://example.openai.azure.com", "any-deploy", "2024-07-01-preview"),
		testkit.AOAIUnitWithKey("sekret-key"),
	)

	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		w.WriteHeader(http.StatusOK)
		w.WriteString(`{"echo":"ok"}`)
		return w.Result(), nil
	})

	core := gateway.NewCoreWithAdapters(transport, fx.Authenticator, fx.Adapter)

	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		apigw.RegisterProvider(grp, fx.BasePath, core)
	})

	body := map[string]any{"model": fx.Model, "messages": []string{"hi"}}
	b, _ := json.Marshal(body)
	resp := api.Post("/api"+fx.BasePath+"/v1/chat/completions", "Content-Type: application/json", bytes.NewReader(b))
	require.Equal(t, http.StatusOK, resp.Code)
	require.JSONEq(t, `{"echo":"ok"}`, resp.Body.String())
}

func TestUnitProxy_AzureOpenAI(t *testing.T) {
	fx := testkit.NewAOAIUnit(t,
		testkit.AOAIUnitWithMapping("gpt-4o", "https://dev-insurgence-openai.openai.azure.com", "dev-openai-gpt4-1", "2024-07-01-preview"),
		testkit.AOAIUnitWithKey("sekret-key"),
	)

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

	core := gateway.NewCoreWithAdapters(transport, fx.Authenticator, fx.Adapter)
	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		apigw.RegisterProvider(grp, fx.BasePath, core)
	})

	body := map[string]any{
		"model":    fx.Model,
		"messages": []map[string]string{{"role": "user", "content": "Say hello test!"}},
	}
	b, _ := json.Marshal(body)

	resp := api.Post(
		"/api"+fx.BasePath+"/v1/chat/completions",
		"Content-Type: application/json",
		"Authorization: Bearer client-token",
		bytes.NewReader(b),
	)

	require.True(t, called, "upstream transport was not invoked")
	require.Equal(t, http.StatusOK, resp.Code)
	require.JSONEq(t, `{"echo":"ok"}`, resp.Body.String())
}

func TestUnitProxy_AzureOpenAI_EnvFallback(t *testing.T) {
	t.Setenv("AOAI_TEST_KEY", "env-secret")
	fx := testkit.NewAOAIUnit(t,
		testkit.AOAIUnitWithMapping("gpt-4o", "https://dev-insurgence-openai.openai.azure.com", "dev-openai-gpt4-1", "2024-07-01-preview"),
		testkit.AOAIUnitWithKeyEnv("AOAI_TEST_KEY"),
	)

	called := false
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		called = true
		require.Equal(t, "env-secret", req.Header.Get("api-key"))
		w := httptest.NewRecorder()
		w.WriteHeader(http.StatusOK)
		w.WriteString(`{"ok":true}`)
		return w.Result(), nil
	})

	core := gateway.NewCoreWithAdapters(transport, fx.Authenticator, fx.Adapter)
	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		apigw.RegisterProvider(grp, fx.BasePath, core)
	})

	resp := api.Post("/api"+fx.BasePath+"/v1/chat/completions", "Content-Type: application/json",
		bytes.NewReader([]byte(`{"model":"`+fx.Model+`","messages":[]}`)),
	)
	require.True(t, called)
	require.Equal(t, http.StatusOK, resp.Code)
}

func TestE2EProxy_AzureOpenAI(t *testing.T) {
	fx := testkit.NewAOAIE2E(t)
	core := gateway.NewCoreWithAdapters(http.DefaultTransport, fx.Authenticator, fx.Adapter)
	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) { apigw.RegisterProvider(grp, fx.BasePath, core) })
	body := map[string]any{"model": fx.Model, "messages": []map[string]string{{"role": "user", "content": "Say hello test!"}}}
	b, _ := json.Marshal(body)
	resp := api.Post("/api"+fx.BasePath+"/v1/chat/completions", "Content-Type: application/json", bytes.NewReader(b))
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
}

// --- test helpers ---

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
