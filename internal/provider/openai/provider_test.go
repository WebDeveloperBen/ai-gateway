package openai_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insurgence-ai/llm-gateway/internal/gateway/loadbalancing"
	"github.com/insurgence-ai/llm-gateway/internal/provider"
	openai "github.com/insurgence-ai/llm-gateway/internal/provider/openai"
)

func TestRewrite_ForwardsPathQuery_AndSetsBearerAndOrg(t *testing.T) {
	ad := openai.New(loadbalancing.NewRoundRobinSelector())
	ad.Keys = provider.KeySource{ForTenant: func(string) string { return "k123" }}
	ad.OrgFor = func(string) string { return "org_abc" }
	ad.Instances["gpt-4o"] = []string{"gpt-4o"}

	body := `{"model":"gpt-4o","messages":[]}`
	req := httptest.NewRequest("POST", "/v1/chat/completions?stream=true", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer client-token") // should be stripped

	err := ad.Rewrite(req, "/v1/chat/completions", provider.ReqInfo{Tenant: "t1", Model: "gpt-4o"})
	require.NoError(t, err)

	require.Equal(t, "https", req.URL.Scheme)
	require.Equal(t, "api.openai.com", req.URL.Host)
	require.Equal(t, "/v1/chat/completions", req.URL.Path)
	require.Equal(t, "stream=true", req.URL.RawQuery)

	require.Equal(t, "Bearer k123", req.Header.Get("Authorization"))
	require.Equal(t, "org_abc", req.Header.Get("OpenAI-Organization"))
}

func TestRewrite_ModelAlias_RewritesJSONModel_AndContentLength(t *testing.T) {
	ad := openai.New(loadbalancing.NewRoundRobinSelector())
	ad.Keys = provider.KeySource{ForTenant: func(string) string { return "k" }}
	ad.ModelAlias = map[string]string{
		"gpt-4o": "gpt-4o-2024-08-06",
	}
	ad.Instances["gpt-4o"] = []string{"gpt-4o"}

	orig := `{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}]}`
	req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewBufferString(orig))
	req.Header.Set("Content-Type", "application/json")

	err := ad.Rewrite(req, "/v1/chat/completions", provider.ReqInfo{Tenant: "t1", Model: "gpt-4o"})
	require.NoError(t, err)

	b, _ := io.ReadAll(req.Body)
	var got map[string]any
	_ = json.Unmarshal(b, &got)

	require.Equal(t, "gpt-4o-2024-08-06", got["model"])
	require.Equal(t, int64(len(b)), req.ContentLength)
}

func TestRewrite_NoAlias_PreservesBody_AndForwardsEmbeddings(t *testing.T) {
	ad := openai.New(loadbalancing.NewRoundRobinSelector())
	ad.Keys = provider.KeySource{ForTenant: func(string) string { return "k" }}
	ad.Instances["gpt-4o"] = []string{"gpt-4o"}

	orig := `{"model":"gpt-4o","input":"hello"}`
	req := httptest.NewRequest("POST", "/v1/embeddings?foo=1", bytes.NewBufferString(orig))

	err := ad.Rewrite(req, "/v1/embeddings", provider.ReqInfo{Model: "gpt-4o"})
	require.NoError(t, err)

	require.Equal(t, "https://api.openai.com/v1/embeddings?foo=1", req.URL.String())

	b, _ := io.ReadAll(req.Body)
	var got map[string]any
	_ = json.Unmarshal(b, &got)
	require.Equal(t, "gpt-4o", got["model"])
}

func TestRewrite_UsesEnvWhenNoTenantKey(t *testing.T) {
	ad := openai.New(loadbalancing.NewRoundRobinSelector())
	ad.Instances["gpt-4o"] = []string{"gpt-4o"}
	ad.Keys = provider.KeySource{EnvVar: "OPENAI_API_KEY"} // default env name
	t.Setenv("OPENAI_API_KEY", "envk")

	req := httptest.NewRequest("POST", "/v1/embeddings", bytes.NewBufferString(`{}`))
	err := ad.Rewrite(req, "/v1/embeddings", provider.ReqInfo{Model: "gpt-4o"})
	require.NoError(t, err)

	require.Equal(t, "Bearer envk", req.Header.Get("Authorization"))
}

func TestRewrite_ModelAlias_IsCaseInsensitive(t *testing.T) {
	ad := openai.New(loadbalancing.NewRoundRobinSelector())
	ad.Keys = provider.KeySource{ForTenant: func(string) string { return "k" }}
	ad.ModelAlias = map[string]string{
		"gpt-4o": "gpt-4o-2024-08-06",
	}
	ad.Instances["gpt-4o"] = []string{"gpt-4o"}

	req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewBufferString(`{"model":"GPT-4O"}`))

	err := ad.Rewrite(req, "/v1/chat/completions", provider.ReqInfo{Model: "  GPT-4O  "})
	require.NoError(t, err)

	b, _ := io.ReadAll(req.Body)
	var got map[string]any
	_ = json.Unmarshal(b, &got)
	require.Equal(t, "gpt-4o-2024-08-06", got["model"])
}
