package tokens

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEstimator(t *testing.T) {
	estimator := NewEstimator()
	require.NotNil(t, estimator)
	assert.NotNil(t, estimator.encodings)
}

func TestGetEncodingForModel(t *testing.T) {
	tests := []struct {
		model    string
		expected string
	}{
		{"gpt-4", "cl100k_base"},
		{"gpt-4-32k", "cl100k_base"},
		{"gpt-4-turbo", "cl100k_base"},
		{"gpt-4-turbo-preview", "cl100k_base"},
		{"gpt-4o", "o200k_base"},
		{"gpt-4o-mini", "o200k_base"},
		{"gpt-3.5-turbo", "cl100k_base"},
		{"gpt-3.5-turbo-16k", "cl100k_base"},
		{"unknown-model", "cl100k_base"},
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			encoding := getEncodingForModel(tt.model)
			assert.Equal(t, tt.expected, encoding)
		})
	}
}

func TestEstimatorEstimateRequest(t *testing.T) {
	estimator := NewEstimator()
	ctx := context.Background()

	t.Run("Chat completion with messages", func(t *testing.T) {
		body := []byte(`{
			"model": "gpt-4",
			"messages": [
				{"role": "system", "content": "You are a helpful assistant."},
				{"role": "user", "content": "Hello, how are you?"}
			]
		}`)

		tokens, err := estimator.EstimateRequest(ctx, "gpt-4", body)
		require.NoError(t, err)
		assert.Greater(t, tokens, 0)
		assert.Greater(t, tokens, 10)
	})

	t.Run("Completion with prompt", func(t *testing.T) {
		body := []byte(`{
			"model": "gpt-3.5-turbo",
			"prompt": "Once upon a time in a land far away"
		}`)

		tokens, err := estimator.EstimateRequest(ctx, "gpt-3.5-turbo", body)
		require.NoError(t, err)
		assert.Greater(t, tokens, 0)
		assert.Greater(t, tokens, 5)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		body := []byte(`invalid json`)
		_, err := estimator.EstimateRequest(ctx, "gpt-4", body)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse request body")
	})

	t.Run("No messages or prompt", func(t *testing.T) {
		body := []byte(`{"model": "gpt-4"}`)
		_, err := estimator.EstimateRequest(ctx, "gpt-4", body)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no messages or prompt found")
	})

	t.Run("Empty messages", func(t *testing.T) {
		body := []byte(`{
			"model": "gpt-4",
			"messages": []
		}`)
		_, err := estimator.EstimateRequest(ctx, "gpt-4", body)
		assert.Error(t, err)
	})

	t.Run("Multiple messages increase token count", func(t *testing.T) {
		singleMessage := []byte(`{
			"model": "gpt-4",
			"messages": [
				{"role": "user", "content": "Hello"}
			]
		}`)

		multipleMessages := []byte(`{
			"model": "gpt-4",
			"messages": [
				{"role": "user", "content": "Hello"},
				{"role": "assistant", "content": "Hi there!"},
				{"role": "user", "content": "How are you?"}
			]
		}`)

		singleTokens, err := estimator.EstimateRequest(ctx, "gpt-4", singleMessage)
		require.NoError(t, err)

		multipleTokens, err := estimator.EstimateRequest(ctx, "gpt-4", multipleMessages)
		require.NoError(t, err)

		assert.Greater(t, multipleTokens, singleTokens)
	})

	t.Run("Overhead tokens included", func(t *testing.T) {
		body := []byte(`{
			"model": "gpt-4",
			"messages": [
				{"role": "user", "content": "Test"}
			]
		}`)

		tokens, err := estimator.EstimateRequest(ctx, "gpt-4", body)
		require.NoError(t, err)

		// Should include: actual tokens + message overhead + array overhead
		// Minimum: 1 token (Test) + 4 (message overhead) + 3 (array overhead) = 8
		assert.GreaterOrEqual(t, tokens, MessageOverheadTokens+ArrayOverheadTokens)
	})
}

func TestEstimatorCaching(t *testing.T) {
	estimator := NewEstimator()
	ctx := context.Background()

	body := []byte(`{
		"model": "gpt-4",
		"messages": [
			{"role": "user", "content": "Test message"}
		]
	}`)

	t.Run("Encoding is cached", func(t *testing.T) {
		// First call should create encoding
		tokens1, err := estimator.EstimateRequest(ctx, "gpt-4", body)
		require.NoError(t, err)

		// Second call should use cached encoding
		tokens2, err := estimator.EstimateRequest(ctx, "gpt-4", body)
		require.NoError(t, err)

		// Results should be the same
		assert.Equal(t, tokens1, tokens2)

		// Check cache contains the model
		_, ok := estimator.encodings.Get("gpt-4")
		assert.True(t, ok, "Encoding should be cached")
	})

	t.Run("Different models cache separately", func(t *testing.T) {
		// GPT-4 uses cl100k_base
		tokens4, err := estimator.EstimateRequest(ctx, "gpt-4", body)
		require.NoError(t, err)
		assert.Greater(t, tokens4, 0)

		// GPT-4o uses o200k_base (different encoding)
		tokens4o, err := estimator.EstimateRequest(ctx, "gpt-4o", body)
		require.NoError(t, err)
		assert.Greater(t, tokens4o, 0)

		// Both models should be cached
		_, ok4 := estimator.encodings.Get("gpt-4")
		assert.True(t, ok4, "gpt-4 encoding should be cached")

		_, ok4o := estimator.encodings.Get("gpt-4o")
		assert.True(t, ok4o, "gpt-4o encoding should be cached")
	})
}

func TestEstimatorConstants(t *testing.T) {
	t.Run("Constants have expected values", func(t *testing.T) {
		assert.Equal(t, 4, MessageOverheadTokens)
		assert.Equal(t, 3, ArrayOverheadTokens)
	})
}
