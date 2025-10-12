package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsedRequestContext(t *testing.T) {
	t.Run("WithParsedRequest and GetParsedRequest", func(t *testing.T) {
		ctx := context.Background()
		req := &ParsedRequest{
			Model:           "gpt-4",
			EstimatedTokens: 100,
			RequestSize:     1024,
			Messages: []Message{
				{Role: "user", Content: "Hello"},
			},
			Prompt: "Test prompt",
		}

		ctx = WithParsedRequest(ctx, req)
		retrieved := GetParsedRequest(ctx)

		require.NotNil(t, retrieved)
		assert.Equal(t, "gpt-4", retrieved.Model)
		assert.Equal(t, 100, retrieved.EstimatedTokens)
		assert.Equal(t, 1024, retrieved.RequestSize)
		assert.Len(t, retrieved.Messages, 1)
		assert.Equal(t, "Hello", retrieved.Messages[0].Content)
		assert.Equal(t, "Test prompt", retrieved.Prompt)
	})

	t.Run("GetParsedRequest returns nil when not set", func(t *testing.T) {
		ctx := context.Background()
		retrieved := GetParsedRequest(ctx)
		assert.Nil(t, retrieved)
	})

	t.Run("GetParsedRequest returns nil for wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), contextKeyParsedRequest{}, "wrong type")
		retrieved := GetParsedRequest(ctx)
		assert.Nil(t, retrieved)
	})
}

func TestGetModel(t *testing.T) {
	t.Run("Returns model from ParsedRequest", func(t *testing.T) {
		ctx := context.Background()
		req := &ParsedRequest{Model: "gpt-4"}
		ctx = WithParsedRequest(ctx, req)

		model := GetModel(ctx)
		assert.Equal(t, "gpt-4", model)
	})

	t.Run("Returns empty string when ParsedRequest not set", func(t *testing.T) {
		ctx := context.Background()
		model := GetModel(ctx)
		assert.Equal(t, "", model)
	})

	t.Run("Returns empty string when Model field is empty", func(t *testing.T) {
		ctx := context.Background()
		req := &ParsedRequest{Model: ""}
		ctx = WithParsedRequest(ctx, req)

		model := GetModel(ctx)
		assert.Equal(t, "", model)
	})
}

func TestGetEstimatedTokens(t *testing.T) {
	t.Run("Returns estimated tokens from ParsedRequest", func(t *testing.T) {
		ctx := context.Background()
		req := &ParsedRequest{EstimatedTokens: 250}
		ctx = WithParsedRequest(ctx, req)

		tokens := GetEstimatedTokens(ctx)
		assert.Equal(t, 250, tokens)
	})

	t.Run("Returns 0 when ParsedRequest not set", func(t *testing.T) {
		ctx := context.Background()
		tokens := GetEstimatedTokens(ctx)
		assert.Equal(t, 0, tokens)
	})

	t.Run("Returns 0 when EstimatedTokens is 0", func(t *testing.T) {
		ctx := context.Background()
		req := &ParsedRequest{EstimatedTokens: 0}
		ctx = WithParsedRequest(ctx, req)

		tokens := GetEstimatedTokens(ctx)
		assert.Equal(t, 0, tokens)
	})

	t.Run("Handles large token counts", func(t *testing.T) {
		ctx := context.Background()
		req := &ParsedRequest{EstimatedTokens: 100000}
		ctx = WithParsedRequest(ctx, req)

		tokens := GetEstimatedTokens(ctx)
		assert.Equal(t, 100000, tokens)
	})
}

func TestGetRequestSize(t *testing.T) {
	t.Run("Returns request size from ParsedRequest", func(t *testing.T) {
		ctx := context.Background()
		req := &ParsedRequest{RequestSize: 2048}
		ctx = WithParsedRequest(ctx, req)

		size := GetRequestSize(ctx)
		assert.Equal(t, 2048, size)
	})

	t.Run("Returns 0 when ParsedRequest not set", func(t *testing.T) {
		ctx := context.Background()
		size := GetRequestSize(ctx)
		assert.Equal(t, 0, size)
	})

	t.Run("Returns 0 when RequestSize is 0", func(t *testing.T) {
		ctx := context.Background()
		req := &ParsedRequest{RequestSize: 0}
		ctx = WithParsedRequest(ctx, req)

		size := GetRequestSize(ctx)
		assert.Equal(t, 0, size)
	})

	t.Run("Handles large request sizes", func(t *testing.T) {
		ctx := context.Background()
		req := &ParsedRequest{RequestSize: 1048576}
		ctx = WithParsedRequest(ctx, req)

		size := GetRequestSize(ctx)
		assert.Equal(t, 1048576, size)
	})
}

func TestParsedRequestWithMessages(t *testing.T) {
	t.Run("Stores and retrieves multiple messages", func(t *testing.T) {
		ctx := context.Background()
		messages := []Message{
			{Role: "system", Content: "You are a helpful assistant"},
			{Role: "user", Content: "What is the weather?"},
			{Role: "assistant", Content: "I can help with that"},
		}

		req := &ParsedRequest{
			Model:    "gpt-4",
			Messages: messages,
		}
		ctx = WithParsedRequest(ctx, req)

		retrieved := GetParsedRequest(ctx)
		require.NotNil(t, retrieved)
		assert.Len(t, retrieved.Messages, 3)
		assert.Equal(t, "system", retrieved.Messages[0].Role)
		assert.Equal(t, "You are a helpful assistant", retrieved.Messages[0].Content)
	})

	t.Run("Handles empty messages slice", func(t *testing.T) {
		ctx := context.Background()
		req := &ParsedRequest{
			Model:    "gpt-4",
			Messages: []Message{},
		}
		ctx = WithParsedRequest(ctx, req)

		retrieved := GetParsedRequest(ctx)
		require.NotNil(t, retrieved)
		assert.Empty(t, retrieved.Messages)
	})

	t.Run("Handles nil messages slice", func(t *testing.T) {
		ctx := context.Background()
		req := &ParsedRequest{
			Model:    "gpt-4",
			Messages: nil,
		}
		ctx = WithParsedRequest(ctx, req)

		retrieved := GetParsedRequest(ctx)
		require.NotNil(t, retrieved)
		assert.Nil(t, retrieved.Messages)
	})
}

func TestParsedRequestChaining(t *testing.T) {
	t.Run("ParsedRequest works with other context values", func(t *testing.T) {
		ctx := context.Background()
		ctx = WithKeyID(ctx, "key-123")
		ctx = WithOrgID(ctx, "org-456")

		req := &ParsedRequest{
			Model:           "gpt-4",
			EstimatedTokens: 150,
		}
		ctx = WithParsedRequest(ctx, req)

		assert.Equal(t, "key-123", GetKeyID(ctx))
		assert.Equal(t, "org-456", GetOrgID(ctx))
		assert.Equal(t, "gpt-4", GetModel(ctx))
		assert.Equal(t, 150, GetEstimatedTokens(ctx))
	})
}

func TestMessage(t *testing.T) {
	t.Run("Message struct has correct fields", func(t *testing.T) {
		msg := Message{
			Role:    "user",
			Content: "Test content",
		}

		assert.Equal(t, "user", msg.Role)
		assert.Equal(t, "Test content", msg.Content)
	})

	t.Run("Message can be empty", func(t *testing.T) {
		msg := Message{}
		assert.Equal(t, "", msg.Role)
		assert.Equal(t, "", msg.Content)
	})
}
