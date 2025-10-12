package tokens

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewParser(t *testing.T) {
	parser := NewParser()
	require.NotNil(t, parser)
	assert.NotNil(t, parser.parsers)
	assert.Len(t, parser.parsers, 5)
}

func TestParserRegisterParser(t *testing.T) {
	parser := NewParser()
	customParser := &OpenAIParser{}

	parser.RegisterParser("custom", customParser)
	assert.Equal(t, customParser, parser.parsers["custom"])
}

func TestParserParseResponse_OpenAI(t *testing.T) {
	parser := NewParser()

	t.Run("Valid OpenAI response", func(t *testing.T) {
		body := []byte(`{
			"usage": {
				"prompt_tokens": 10,
				"completion_tokens": 20,
				"total_tokens": 30
			}
		}`)

		usage, err := parser.ParseResponse("openai", body)
		require.NoError(t, err)
		assert.Equal(t, 10, usage.PromptTokens)
		assert.Equal(t, 20, usage.CompletionTokens)
		assert.Equal(t, 30, usage.TotalTokens)
	})

	t.Run("Azure OpenAI uses same parser", func(t *testing.T) {
		body := []byte(`{
			"usage": {
				"prompt_tokens": 15,
				"completion_tokens": 25,
				"total_tokens": 40
			}
		}`)

		usage, err := parser.ParseResponse("azureopenai", body)
		require.NoError(t, err)
		assert.Equal(t, 15, usage.PromptTokens)
		assert.Equal(t, 25, usage.CompletionTokens)
		assert.Equal(t, 40, usage.TotalTokens)
	})

	t.Run("Unsupported provider", func(t *testing.T) {
		body := []byte(`{}`)
		_, err := parser.ParseResponse("unknown", body)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported provider")
	})
}

func TestParserParseResponseReader(t *testing.T) {
	parser := NewParser()

	t.Run("Valid reader", func(t *testing.T) {
		body := []byte(`{
			"usage": {
				"prompt_tokens": 10,
				"completion_tokens": 20,
				"total_tokens": 30
			}
		}`)
		reader := bytes.NewReader(body)

		usage, err := parser.ParseResponseReader("openai", reader)
		require.NoError(t, err)
		assert.Equal(t, 10, usage.PromptTokens)
	})
}

func TestParserParseStreamedResponse(t *testing.T) {
	parser := NewParser()

	t.Run("Usage in last chunk", func(t *testing.T) {
		chunks := [][]byte{
			[]byte(`{"id": "chunk1", "object": "chat.completion.chunk"}`),
			[]byte(`{"id": "chunk2", "object": "chat.completion.chunk"}`),
			[]byte(`{
				"id": "chunk3",
				"usage": {
					"prompt_tokens": 10,
					"completion_tokens": 20,
					"total_tokens": 30
				}
			}`),
		}

		usage, err := parser.ParseStreamedResponse("openai", chunks)
		require.NoError(t, err)
		assert.Equal(t, 10, usage.PromptTokens)
		assert.Equal(t, 20, usage.CompletionTokens)
		assert.Equal(t, 30, usage.TotalTokens)
	})

	t.Run("Usage in middle chunk", func(t *testing.T) {
		chunks := [][]byte{
			[]byte(`{"id": "chunk1"}`),
			[]byte(`{
				"usage": {
					"prompt_tokens": 5,
					"completion_tokens": 10,
					"total_tokens": 15
				}
			}`),
			[]byte(`{"id": "chunk3"}`),
		}

		usage, err := parser.ParseStreamedResponse("openai", chunks)
		require.NoError(t, err)
		assert.Equal(t, 5, usage.PromptTokens)
	})

	t.Run("No usage data", func(t *testing.T) {
		chunks := [][]byte{
			[]byte(`{"id": "chunk1"}`),
			[]byte(`{"id": "chunk2"}`),
		}

		_, err := parser.ParseStreamedResponse("openai", chunks)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no usage data found")
	})
}
