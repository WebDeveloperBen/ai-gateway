package tokens

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenAIParser(t *testing.T) {
	parser := &OpenAIParser{}

	t.Run("Valid response", func(t *testing.T) {
		body := []byte(`{
			"id": "chatcmpl-123",
			"object": "chat.completion",
			"usage": {
				"prompt_tokens": 10,
				"completion_tokens": 20,
				"total_tokens": 30
			}
		}`)

		usage, err := parser.ParseResponse(body)
		require.NoError(t, err)
		assert.Equal(t, 10, usage.PromptTokens)
		assert.Equal(t, 20, usage.CompletionTokens)
		assert.Equal(t, 30, usage.TotalTokens)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		body := []byte(`invalid json`)
		_, err := parser.ParseResponse(body)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse OpenAI response")
	})

	t.Run("No usage data", func(t *testing.T) {
		body := []byte(`{
			"id": "chatcmpl-123",
			"object": "chat.completion"
		}`)
		_, err := parser.ParseResponse(body)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no usage data in response")
	})

	t.Run("Zero tokens", func(t *testing.T) {
		body := []byte(`{
			"usage": {
				"prompt_tokens": 0,
				"completion_tokens": 0,
				"total_tokens": 0
			}
		}`)
		_, err := parser.ParseResponse(body)
		assert.Error(t, err)
	})
}

func TestAnthropicParser(t *testing.T) {
	parser := &AnthropicParser{}

	t.Run("Valid response", func(t *testing.T) {
		body := []byte(`{
			"id": "msg_123",
			"type": "message",
			"usage": {
				"input_tokens": 15,
				"output_tokens": 25
			}
		}`)

		usage, err := parser.ParseResponse(body)
		require.NoError(t, err)
		assert.Equal(t, 15, usage.PromptTokens)
		assert.Equal(t, 25, usage.CompletionTokens)
		assert.Equal(t, 40, usage.TotalTokens)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		body := []byte(`invalid json`)
		_, err := parser.ParseResponse(body)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse Anthropic response")
	})

	t.Run("No usage data", func(t *testing.T) {
		body := []byte(`{
			"id": "msg_123",
			"type": "message"
		}`)
		_, err := parser.ParseResponse(body)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no usage data in response")
	})
}

func TestCohereParser(t *testing.T) {
	parser := &CohereParser{}

	t.Run("Valid response", func(t *testing.T) {
		body := []byte(`{
			"id": "gen-123",
			"meta": {
				"billed_units": {
					"input_tokens": 12,
					"output_tokens": 18
				}
			}
		}`)

		usage, err := parser.ParseResponse(body)
		require.NoError(t, err)
		assert.Equal(t, 12, usage.PromptTokens)
		assert.Equal(t, 18, usage.CompletionTokens)
		assert.Equal(t, 30, usage.TotalTokens)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		body := []byte(`invalid json`)
		_, err := parser.ParseResponse(body)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse Cohere response")
	})

	t.Run("No usage data", func(t *testing.T) {
		body := []byte(`{
			"id": "gen-123",
			"meta": {}
		}`)
		_, err := parser.ParseResponse(body)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no usage data in response")
	})
}

func TestGoogleParser(t *testing.T) {
	parser := &GoogleParser{}

	t.Run("Valid response", func(t *testing.T) {
		body := []byte(`{
			"candidates": [],
			"usageMetadata": {
				"promptTokenCount": 8,
				"candidatesTokenCount": 22,
				"totalTokenCount": 30
			}
		}`)

		usage, err := parser.ParseResponse(body)
		require.NoError(t, err)
		assert.Equal(t, 8, usage.PromptTokens)
		assert.Equal(t, 22, usage.CompletionTokens)
		assert.Equal(t, 30, usage.TotalTokens)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		body := []byte(`invalid json`)
		_, err := parser.ParseResponse(body)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse Google response")
	})

	t.Run("No usage data", func(t *testing.T) {
		body := []byte(`{
			"candidates": []
		}`)
		_, err := parser.ParseResponse(body)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no usage data in response")
	})
}
