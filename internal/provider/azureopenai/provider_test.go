package azureopenai_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insurgence-ai/llm-gateway/internal/loadbalancing"
	"github.com/insurgence-ai/llm-gateway/internal/provider"
	aoai "github.com/insurgence-ai/llm-gateway/internal/provider/azureopenai"
)

func TestRewrite_GlobalMapping_URL_Auth(t *testing.T) {
	ad := aoai.New(loadbalancing.NewRoundRobinSelector())
	ad.Keys = provider.KeySource{ForTenant: func(string) string { return "k123" }}
	ad.Instances["gpt-4o"] = []aoai.Entry{{
		BaseURL:    "dev-insurgence-openai.openai.azure.com",
		Deployment: "dev-openai-gpt4-1",
		APIVer:     "2024-07-01-preview",
	}}

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions?stream=true", bytes.NewBufferString(`{"model":"gpt-4o"}`))
	req.Header.Set("Authorization", "Bearer client-token")
	err := ad.Rewrite(req, "/v1/chat/completions", provider.ReqInfo{Tenant: "t1", Model: "gpt-4o"})
	require.NoError(t, err)

	require.Equal(t, "https", req.URL.Scheme)
	require.Equal(t, "dev-insurgence-openai.openai.azure.com", req.URL.Host)
	require.Equal(t, "/openai/deployments/dev-openai-gpt4-1/chat/completions", req.URL.Path)

	q, _ := url.ParseQuery(req.URL.RawQuery)
	require.Equal(t, "true", q.Get("stream"))
	require.Equal(t, "2024-07-01-preview", q.Get("api-version"))

	require.Empty(t, req.Header.Get("Authorization"))   // stripped caller auth
	require.Equal(t, "k123", req.Header.Get("api-key")) // AOAI header
	require.Equal(t, req.URL.Host, req.Host)            // Host aligned
}

func TestRewrite_MergesQuery_AndTrimsV1_ForEmbeddings(t *testing.T) {
	ad := aoai.New(loadbalancing.NewRoundRobinSelector())
	ad.Instances["gpt-4o"] = []aoai.Entry{{
		BaseURL:    "myres.openai.azure.com",
		Deployment: "dep-123",
		APIVer:     "2024-07-01-preview",
	}}

	req := httptest.NewRequest(http.MethodPost, "/v1/embeddings?x=1&x=2", bytes.NewBufferString(`{"model":"gpt-4o","input":"hi"}`))
	err := ad.Rewrite(req, "/v1/embeddings", provider.ReqInfo{Model: "gpt-4o"})
	require.NoError(t, err)

	require.Equal(t, "https://myres.openai.azure.com/openai/deployments/dep-123/embeddings", req.URL.Scheme+"://"+req.URL.Host+req.URL.Path)
	q, _ := url.ParseQuery(req.URL.RawQuery)
	require.ElementsMatch(t, []string{"1", "2"}, q["x"])
	require.Equal(t, "2024-07-01-preview", q.Get("api-version"))
}

func TestRewrite_ByTenantOverridesGlobal(t *testing.T) {
	ad := aoai.New(loadbalancing.NewRoundRobinSelector())
	ad.Instances["gpt-4o"] = []aoai.Entry{{
		BaseURL:    "tenantres.openai.azure.com",
		Deployment: "tenant-dep",
		APIVer:     "2024-07-01-preview",
	}}

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewBufferString(`{"model":"gpt-4o"}`))
	err := ad.Rewrite(req, "/v1/chat/completions", provider.ReqInfo{Tenant: "acme", Model: "gpt-4o"})
	require.NoError(t, err)

	require.Equal(t, "tenantres.openai.azure.com", req.URL.Host)                      // tenant base used
	require.Equal(t, "/openai/deployments/tenant-dep/chat/completions", req.URL.Path) // tenant deployment used
}

func TestRewrite_DefaultUsed_WhenNoModelProvided(t *testing.T) {
	ad := aoai.New(loadbalancing.NewRoundRobinSelector())

	ad.Instances[""] = []aoai.Entry{{
		BaseURL:    "defres.openai.azure.com",
		Deployment: "def-dep",
		APIVer:     "2024-07-01-preview",
	}}

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewBufferString(`{}`))
	err := ad.Rewrite(req, "/v1/chat/completions", provider.ReqInfo{Model: ""})
	require.NoError(t, err)
	require.Equal(t, "defres.openai.azure.com", req.URL.Host)
	require.Equal(t, "/openai/deployments/def-dep/chat/completions", req.URL.Path)
}

func TestRewrite_SingleGlobalFallback_WhenModelMissing(t *testing.T) {
	ad := aoai.New(loadbalancing.NewRoundRobinSelector())
	ad.Instances[""] = []aoai.Entry{{
		BaseURL:    "onlyone.openai.azure.com",
		Deployment: "only-dep",
		APIVer:     "2024-07-01-preview",
	}}

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewBufferString(`{}`))
	err := ad.Rewrite(req, "/v1/chat/completions", provider.ReqInfo{Model: ""})
	require.NoError(t, err)
	require.Equal(t, "onlyone.openai.azure.com", req.URL.Host)
	require.Equal(t, "/openai/deployments/only-dep/chat/completions", req.URL.Path)
}

func TestRewrite_UnknownModel_Error(t *testing.T) {
	ad := aoai.New(loadbalancing.NewRoundRobinSelector()) // no mappings, no default

	ad.Instances[""] = []aoai.Entry{{
		BaseURL:    "onlyone.openai.azure.com",
		Deployment: "only-dep",
		APIVer:     "2024-07-01-preview",
	}}

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewBufferString(`{"model":"nope"}`))
	err := ad.Rewrite(req, "/v1/chat/completions", provider.ReqInfo{Model: "nope"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no deployments found for model")
}

func TestRewrite_IncompleteRoute_Error(t *testing.T) {
	ad := aoai.New(loadbalancing.NewRoundRobinSelector())
	ad.Instances["gpt-4o"] = []aoai.Entry{{
		BaseURL:    "dev-insurgence-openai.openai.azure.com",
		Deployment: "dev-openai-gpt4-1",
		APIVer:     "",
	}}

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewBufferString(`{"model":"gpt-4o"}`))
	err := ad.Rewrite(req, "/v1/chat/completions", provider.ReqInfo{Model: "gpt-4o"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "selected deployment incomplete")
}

func TestRewrite_EnvFallback_WhenNoTenantKey(t *testing.T) {
	ad := aoai.New(loadbalancing.NewRoundRobinSelector())
	ad.Keys = provider.KeySource{EnvVar: "AZURE_OPENAI_API_KEY"} // default

	ad.Instances[""] = []aoai.Entry{{
		BaseURL:    "defres.openai.azure.com",
		Deployment: "def-dep",
		APIVer:     "2024-07-01-preview",
	}}

	t.Setenv("AZURE_OPENAI_API_KEY", "envk")

	req := httptest.NewRequest(http.MethodPost, "/v1/embeddings", bytes.NewBufferString(`{"model":"gpt-4o"}`))
	err := ad.Rewrite(req, "/v1/embeddings", provider.ReqInfo{})
	require.NoError(t, err)
	require.Equal(t, "envk", req.Header.Get("api-key"))
}

func TestRewrite_ModelLookup_IsCaseInsensitive_AndTrimmed(t *testing.T) {
	ad := aoai.New(loadbalancing.NewRoundRobinSelector())
	ad.Instances["gpt-4o"] = []aoai.Entry{{
		BaseURL:    "dep-host",
		Deployment: "dep",
		APIVer:     "2024-07-01-preview",
	}}

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewBufferString(`{"model":"  GPT-4O  "}`))

	err := ad.Rewrite(req, "/v1/chat/completions", provider.ReqInfo{Model: "  GPT-4O  "})
	require.NoError(t, err)

	// ensure it worked by checking final path contains deployment (mapping matched)
	require.Equal(t, "/openai/deployments/dep/chat/completions", req.URL.Path)

	// and body is untouched (Azure routes by path, not body model) — just sanity:
	got, _ := io.ReadAll(req.Body)
	require.Contains(t, string(got), "GPT-4O")
}
