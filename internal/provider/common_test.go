package provider

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureAbsoluteBase(t *testing.T) {
	tests := []struct {
		name              string
		base              string
		defaultHostSuffix string
		expected          string
		wantErr           bool
	}{
		{
			name:              "Empty base",
			base:              "",
			defaultHostSuffix: "",
			wantErr:           true,
		},
		{
			name:              "Full URL",
			base:              "https://api.example.com",
			defaultHostSuffix: "",
			expected:          "https://api.example.com",
		},
		{
			name:              "Hostname with protocol",
			base:              "http://example.com",
			defaultHostSuffix: "",
			expected:          "http://example.com",
		},
		{
			name:              "Hostname with dot",
			base:              "api.example.com",
			defaultHostSuffix: "",
			expected:          "https://api.example.com",
		},
		{
			name:              "Simple name with default suffix",
			base:              "my-resource",
			defaultHostSuffix: "openai.azure.com",
			expected:          "https://my-resource.openai.azure.com",
		},
		{
			name:              "Simple name without suffix",
			base:              "my-resource",
			defaultHostSuffix: "",
			expected:          "https://my-resource",
		},
		{
			name:              "URL with trailing slash",
			base:              "https://api.example.com/",
			defaultHostSuffix: "",
			expected:          "https://api.example.com",
		},
		{
			name:              "Hostname with spaces",
			base:              "  api.example.com  ",
			defaultHostSuffix: "",
			expected:          "https://api.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := EnsureAbsoluteBase(tt.base, tt.defaultHostSuffix)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJoinURL(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		segs     []string
		query    url.Values
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple join",
			base:     "https://api.example.com",
			segs:     []string{"v1", "chat", "completions"},
			expected: "https://api.example.com/v1/chat/completions",
		},
		{
			name:     "With query params",
			base:     "https://api.example.com",
			segs:     []string{"v1", "models"},
			query:    url.Values{"api-version": []string{"2023-05-15"}},
			expected: "https://api.example.com/v1/models?api-version=2023-05-15",
		},
		{
			name:     "Empty segments",
			base:     "https://api.example.com",
			segs:     []string{"", "v1", "", "models"},
			expected: "https://api.example.com/v1/models",
		},
		{
			name:     "Base with existing path",
			base:     "https://api.example.com/base",
			segs:     []string{"v1"},
			expected: "https://api.example.com/base/v1",
		},
		{
			name:     "Merge queries",
			base:     "https://api.example.com?existing=param",
			segs:     []string{"v1"},
			query:    url.Values{"new": []string{"value"}},
			expected: "https://api.example.com/v1?existing=param&new=value",
		},
		{
			name:    "Invalid base URL",
			base:    "://invalid",
			segs:    []string{"v1"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := JoinURL(tt.base, tt.segs, tt.query)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result.String())
		})
	}
}

func TestSetUpstreamURL(t *testing.T) {
	req, err := http.NewRequest("GET", "https://original.com/path?query=1", nil)
	require.NoError(t, err)

	upstreamURL, _ := url.Parse("https://upstream.com/new/path?newquery=2")
	SetUpstreamURL(req, upstreamURL)

	assert.Equal(t, "https", req.URL.Scheme)
	assert.Equal(t, "upstream.com", req.URL.Host)
	assert.Equal(t, "/new/path", req.URL.Path)
	assert.Equal(t, "newquery=2", req.URL.RawQuery)
	assert.Equal(t, "upstream.com", req.Host)
}

func TestCopyQuery(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected url.Values
	}{
		{
			name:     "Single query param",
			url:      "https://example.com?key=value",
			expected: url.Values{"key": []string{"value"}},
		},
		{
			name:     "Multiple params",
			url:      "https://example.com?key1=value1&key2=value2",
			expected: url.Values{"key1": []string{"value1"}, "key2": []string{"value2"}},
		},
		{
			name:     "No query params",
			url:      "https://example.com",
			expected: url.Values{},
		},
		{
			name:     "Duplicate keys",
			url:      "https://example.com?key=value1&key=value2",
			expected: url.Values{"key": []string{"value1", "value2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			require.NoError(t, err)

			result := CopyQuery(req)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCopyQuery_NilURL(t *testing.T) {
	req := &http.Request{URL: nil}
	result := CopyQuery(req)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestStripCallerAuth(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Bearer token123")
	h.Set("Content-Type", "application/json")

	StripCallerAuth(h)

	assert.Empty(t, h.Get("Authorization"))
	assert.Equal(t, "application/json", h.Get("Content-Type"))
}

func TestSetAPIKey(t *testing.T) {
	tests := []struct {
		name       string
		headerName string
		key        string
		shouldSet  bool
	}{
		{
			name:       "Set key",
			headerName: "api-key",
			key:        "secret123",
			shouldSet:  true,
		},
		{
			name:       "Empty key does not set header",
			headerName: "api-key",
			key:        "",
			shouldSet:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.Header{}
			SetAPIKey(h, tt.headerName, tt.key)

			if tt.shouldSet {
				assert.Equal(t, tt.key, h.Get(tt.headerName))
			} else {
				assert.Empty(t, h.Get(tt.headerName))
			}
		})
	}
}

func TestForceContentLength(t *testing.T) {
	body := []byte(`{"test":"data"}`)
	req, err := http.NewRequest("POST", "https://example.com", bytes.NewReader(body))
	require.NoError(t, err)

	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("Content-Encoding", "gzip")

	ForceContentLength(req, len(body))

	assert.Equal(t, int64(len(body)), req.ContentLength)
	assert.Equal(t, "15", req.Header.Get("Content-Length"))
	assert.Empty(t, req.Header.Get("Transfer-Encoding"))
	assert.Empty(t, req.Header.Get("Content-Encoding"))
}

func TestModelOrDefault(t *testing.T) {
	hasExact := func(m string) bool {
		return m == "gpt-4"
	}
	single := func() (string, bool) {
		return "single-model", true
	}

	tests := []struct {
		name           string
		model          string
		hasExact       func(string) bool
		single         func() (string, bool)
		fallbackExists bool
		fallbackKey    string
		expectedKey    string
		expectedOK     bool
	}{
		{
			name:        "Exact match",
			model:       "gpt-4",
			hasExact:    hasExact,
			single:      single,
			expectedKey: "gpt-4",
			expectedOK:  true,
		},
		{
			name:        "Case insensitive match",
			model:       "GPT-4",
			hasExact:    hasExact,
			single:      single,
			expectedKey: "gpt-4",
			expectedOK:  true,
		},
		{
			name:     "Fall back to single",
			model:    "unknown-model",
			hasExact: func(string) bool { return false },
			single:   single,
			expectedKey: "single-model",
			expectedOK:  true,
		},
		{
			name:           "Fall back to default",
			model:          "unknown-model",
			hasExact:       func(string) bool { return false },
			single:         func() (string, bool) { return "", false },
			fallbackExists: true,
			fallbackKey:    "default-model",
			expectedKey:    "default-model",
			expectedOK:     true,
		},
		{
			name:        "No match",
			model:       "unknown-model",
			hasExact:    func(string) bool { return false },
			single:      func() (string, bool) { return "", false },
			expectedKey: "",
			expectedOK:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, ok := ModelOrDefault(tt.model, tt.hasExact, tt.single, tt.fallbackExists, tt.fallbackKey)
			assert.Equal(t, tt.expectedOK, ok)
			assert.Equal(t, tt.expectedKey, key)
		})
	}
}

func TestKeySource_Resolve(t *testing.T) {
	tests := []struct {
		name        string
		keySource   KeySource
		tenant      string
		defaultEnv  string
		envVars     map[string]string
		expected    string
	}{
		{
			name: "ForTenant returns value",
			keySource: KeySource{
				EnvVar: "API_KEY",
				ForTenant: func(t string) string {
					if t == "tenant1" {
						return "tenant1-key"
					}
					return ""
				},
			},
			tenant:   "tenant1",
			expected: "tenant1-key",
		},
		{
			name: "ForTenant empty, fallback to EnvVar",
			keySource: KeySource{
				EnvVar: "API_KEY",
				ForTenant: func(t string) string {
					return ""
				},
			},
			tenant:     "tenant1",
			defaultEnv: "DEFAULT_KEY",
			envVars:    map[string]string{"API_KEY": "env-key"},
			expected:   "env-key",
		},
		{
			name: "No ForTenant, use EnvVar",
			keySource: KeySource{
				EnvVar: "API_KEY",
			},
			envVars:  map[string]string{"API_KEY": "env-key"},
			expected: "env-key",
		},
		{
			name: "Empty EnvVar, use defaultEnv",
			keySource: KeySource{
				EnvVar: "",
			},
			defaultEnv: "DEFAULT_KEY",
			envVars:    map[string]string{"DEFAULT_KEY": "default-key"},
			expected:   "default-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			result := tt.keySource.Resolve(tt.tenant, tt.defaultEnv)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRewriteJSONField(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		field         string
		value         any
		expectedBody  string
		expectedError bool
	}{
		{
			name:         "Rewrite existing field",
			body:         `{"model":"gpt-3.5","temperature":0.5}`,
			field:        "model",
			value:        "gpt-4",
			expectedBody: `{"model":"gpt-4","temperature":0.5}`,
		},
		{
			name:         "Add new field",
			body:         `{"temperature":0.5}`,
			field:        "model",
			value:        "gpt-4",
			expectedBody: `{"model":"gpt-4","temperature":0.5}`,
		},
		{
			name:         "Rewrite with different type",
			body:         `{"max_tokens":100}`,
			field:        "max_tokens",
			value:        200,
			expectedBody: `{"max_tokens":200}`,
		},
		{
			name:         "Rewrite complex nested object",
			body:         `{"model":"gpt-3.5","options":{"temperature":0.5}}`,
			field:        "model",
			value:        "gpt-4",
			expectedBody: `{"model":"gpt-4","options":{"temperature":0.5}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "https://example.com", bytes.NewBufferString(tt.body))
			require.NoError(t, err)

			err = RewriteJSONField(req, tt.field, tt.value)
			require.NoError(t, err)

			resultBody, err := io.ReadAll(req.Body)
			require.NoError(t, err)

			assert.JSONEq(t, tt.expectedBody, string(resultBody))
		})
	}
}

func TestRewriteJSONField_NilBody(t *testing.T) {
	req, err := http.NewRequest("GET", "https://example.com", nil)
	require.NoError(t, err)

	err = RewriteJSONField(req, "model", "gpt-4")
	require.NoError(t, err)
}

func TestRewriteJSONModel(t *testing.T) {
	body := `{"model":"gpt-3.5","temperature":0.5}`
	req, err := http.NewRequest("POST", "https://example.com", bytes.NewBufferString(body))
	require.NoError(t, err)

	err = RewriteJSONModel(req, "gpt-4")
	require.NoError(t, err)

	resultBody, err := io.ReadAll(req.Body)
	require.NoError(t, err)

	expected := `{"model":"gpt-4","temperature":0.5}`
	assert.JSONEq(t, expected, string(resultBody))
}
