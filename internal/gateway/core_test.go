package gateway_test

import (
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway"
	"github.com/WebDeveloperBen/ai-gateway/internal/provider"
	"github.com/stretchr/testify/require"
)

func TestIndexOfSegment(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		needle   string
		expected int
	}{
		{
			name:     "exact match at start",
			path:     "/azure/openai/v1/chat",
			needle:   "/azure/openai",
			expected: 0,
		},
		{
			name:     "match at start of segment",
			path:     "/azure/openai/v1/chat",
			needle:   "/azure/openai",
			expected: 0, // "/azure/openai" starts at position 0, at segment boundary
		},
		{
			name:     "no match - not on segment boundary",
			path:     "/azureopenai/v1/chat",
			needle:   "/azure/openai",
			expected: -1,
		},
		{
			name:     "no match - partial segment",
			path:     "/azure/openai-extra/v1/chat",
			needle:   "/azure/openai",
			expected: -1,
		},
		{
			name:     "empty needle",
			path:     "/azure/openai/v1/chat",
			needle:   "",
			expected: -1,
		},
		{
			name:     "match at end",
			path:     "/azure/openai",
			needle:   "/azure/openai",
			expected: 0,
		},
		{
			name:     "single segment match",
			path:     "/openai/v1/chat",
			needle:   "/openai",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gateway.IndexOfSegment(tt.path, tt.needle)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestListPrefixes(t *testing.T) {
	t.Run("empty adapters list", func(t *testing.T) {
		var adapters []provider.Adapter
		result := gateway.ListPrefixes(adapters)
		require.Empty(t, result)
	})

	t.Run("single adapter", func(t *testing.T) {
		adapters := []provider.Adapter{
			&mockAdapter{prefix: "/openai"},
		}
		result := gateway.ListPrefixes(adapters)
		require.Equal(t, []string{"/openai"}, result)
	})

	t.Run("multiple adapters", func(t *testing.T) {
		adapters := []provider.Adapter{
			&mockAdapter{prefix: "/openai"},
			&mockAdapter{prefix: "/azure/openai"},
			&mockAdapter{prefix: "/anthropic"},
		}
		result := gateway.ListPrefixes(adapters)
		require.Equal(t, []string{"/openai", "/azure/openai", "/anthropic"}, result)
	})
}

func TestExtractModel(t *testing.T) {
	tests := []struct {
		name     string
		jsonBody string
		expected string
	}{
		{
			name:     "empty body",
			jsonBody: "",
			expected: "",
		},
		{
			name:     "standard model field",
			jsonBody: `{"model": "gpt-4", "messages": []}`,
			expected: "gpt-4",
		},
		{
			name:     "deployment field (AOAI)",
			jsonBody: `{"deployment": "gpt-4-deployment", "messages": []}`,
			expected: "gpt-4-deployment",
		},
		{
			name:     "engine field (legacy)",
			jsonBody: `{"engine": "text-davinci-003", "prompt": "test"}`,
			expected: "text-davinci-003",
		},
		{
			name:     "model takes precedence over deployment",
			jsonBody: `{"model": "gpt-4", "deployment": "ignored", "messages": []}`,
			expected: "gpt-4",
		},
		{
			name:     "whitespace trimmed",
			jsonBody: `{"model": "  gpt-4  ", "messages": []}`,
			expected: "gpt-4",
		},
		{
			name:     "empty model field",
			jsonBody: `{"model": "", "messages": []}`,
			expected: "",
		},
		{
			name:     "invalid json",
			jsonBody: `{"model": "gpt-4"`,
			expected: "",
		},
		{
			name:     "fallback to map lookup",
			jsonBody: `{"custom_model": "gpt-4", "model": "gpt-3.5"}`,
			expected: "gpt-3.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := []byte(tt.jsonBody)
			result := gateway.ExtractModel(raw)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple string",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "string with quotes",
			input:    `error: "invalid request"`,
			expected: `error: \"invalid request\"`,
		},
		{
			name:     "string with newlines",
			input:    "line1\nline2",
			expected: "line1\\nline2",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "string with backslashes",
			input:    `path\to\file`,
			expected: `path\\to\\file`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gateway.Escape(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestGetProviderName(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		expected string
	}{
		{
			name:     "azure openai prefix",
			prefix:   "/azure/openai",
			expected: "azureopenai",
		},
		{
			name:     "openai prefix",
			prefix:   "/openai",
			expected: "openai",
		},
		{
			name:     "anthropic prefix",
			prefix:   "/anthropic",
			expected: "anthropic",
		},
		{
			name:     "empty prefix",
			prefix:   "",
			expected: "unknown",
		},
		{
			name:     "root prefix",
			prefix:   "/",
			expected: "unknown",
		},
		{
			name:     "deep nested prefix",
			prefix:   "/api/v1/azure/openai",
			expected: "apiv1azureopenai",
		},
		{
			name:     "single segment",
			prefix:   "/test",
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &mockAdapter{prefix: tt.prefix}
			result := gateway.GetProviderName(adapter)
			require.Equal(t, tt.expected, result)
		})
	}
}
