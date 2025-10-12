package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsageRecorder_parseTokenUsage(t *testing.T) {
	recorder := NewUsageRecorder(nil, nil)

	t.Run("OpenAI_RegularResponse", func(t *testing.T) {
		body := []byte(`{
			"id": "chatcmpl-123",
			"object": "chat.completion",
			"usage": {
				"prompt_tokens": 10,
				"completion_tokens": 20,
				"total_tokens": 30
			}
		}`)

		usage, err := recorder.parseTokenUsage("openai", body)
		require.NoError(t, err)
		assert.NotNil(t, usage)
		assert.Equal(t, 10, usage.PromptTokens)
		assert.Equal(t, 20, usage.CompletionTokens)
		assert.Equal(t, 30, usage.TotalTokens)
	})

	t.Run("OpenAI_StreamingResponse", func(t *testing.T) {
		body := `data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-4","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}

data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-4","choices":[{"index":0,"delta":{},"finish_reason":"stop"}],"usage":{"prompt_tokens":5,"completion_tokens":10,"total_tokens":15}}

data: [DONE]

`

		usage, err := recorder.parseTokenUsage("openai", []byte(body))
		require.NoError(t, err)
		assert.NotNil(t, usage)
		assert.Equal(t, 5, usage.PromptTokens)
		assert.Equal(t, 10, usage.CompletionTokens)
		assert.Equal(t, 15, usage.TotalTokens)
	})

	t.Run("OpenAI_NoUsageData", func(t *testing.T) {
		body := []byte(`{
			"id": "chatcmpl-123",
			"object": "chat.completion",
			"choices": [{"message": {"content": "Hello"}}]
		}`)

		usage, err := recorder.parseTokenUsage("openai", body)
		require.Error(t, err) // Should error because no usage data
		assert.Nil(t, usage)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		body := []byte(`invalid json`)

		usage, err := recorder.parseTokenUsage("openai", body)
		require.Error(t, err) // Should error because invalid JSON
		assert.Nil(t, usage)
	})
}

// testNopCloser implements io.ReadCloser for testing
type testNopCloser struct {
	io.Reader
}

func (testNopCloser) Close() error { return nil }

func TestResponseBodyWrapper(t *testing.T) {
	t.Run("ReadAndClose", func(t *testing.T) {
		originalBody := &testNopCloser{bytes.NewReader([]byte("test response"))}
		var capturedBody []byte
		onClose := func(body []byte) {
			capturedBody = body
		}

		wrapper := &responseBodyWrapper{
			originalBody: originalBody,
			body:         &bytes.Buffer{},
			onClose:      onClose,
		}

		// Read all data
		buf := make([]byte, 100)
		n, err := wrapper.Read(buf)
		require.NoError(t, err)
		assert.Equal(t, 13, n) // "test response" is 13 bytes
		assert.Equal(t, "test response", string(buf[:n]))

		// Close
		err = wrapper.Close()
		require.NoError(t, err)

		// Check that onClose was called with captured body
		assert.Equal(t, []byte("test response"), capturedBody)

		// Subsequent reads should fail
		n, err = wrapper.Read(buf)
		assert.Equal(t, 0, n)
		assert.Equal(t, http.ErrBodyReadAfterClose, err)
	})

	t.Run("CloseWithoutRead", func(t *testing.T) {
		originalBody := &testNopCloser{bytes.NewReader([]byte("test"))}
		var capturedBody []byte
		onClose := func(body []byte) {
			capturedBody = body
		}

		wrapper := &responseBodyWrapper{
			originalBody: originalBody,
			body:         &bytes.Buffer{},
			onClose:      onClose,
		}

		err := wrapper.Close()
		require.NoError(t, err)
		// Since nothing was read, the buffer is empty
		assert.Nil(t, capturedBody)
	})

	t.Run("MultipleCloseCalls", func(t *testing.T) {
		originalBody := &testNopCloser{bytes.NewReader([]byte("test"))}
		callCount := 0
		onClose := func(body []byte) {
			callCount++
		}

		wrapper := &responseBodyWrapper{
			originalBody: originalBody,
			body:         &bytes.Buffer{},
			onClose:      onClose,
		}

		wrapper.Close()
		wrapper.Close() // Should be safe to call multiple times

		assert.Equal(t, 1, callCount) // onClose should only be called once
	})
}

func TestDetachedContext(t *testing.T) {
	// Create original context with values
	originalCtx := context.Background()
	originalCtx = auth.WithOrgID(originalCtx, "org-123")
	originalCtx = auth.WithAppID(originalCtx, "app-456")
	originalCtx = auth.WithKeyID(originalCtx, "key-789")
	originalCtx = auth.WithUserID(originalCtx, "user-999")
	originalCtx = auth.WithProvider(originalCtx, "openai")
	originalCtx = auth.WithModelName(originalCtx, "gpt-4")

	// Detach context
	detachedCtx := detachContext(originalCtx)

	t.Run("DetachedValues", func(t *testing.T) {
		assert.Equal(t, "org-123", getDetachedOrgID(detachedCtx))
		assert.Equal(t, "app-456", getDetachedAppID(detachedCtx))
		assert.Equal(t, "key-789", getDetachedKeyID(detachedCtx))
	})

	t.Run("FallbackToOriginalContext", func(t *testing.T) {
		// Test fallback when detached context doesn't have the value
		// Since the detached context is empty, it should fall back to auth.Get* functions
		// which will return empty strings since context.Background() has no auth values
		emptyDetachedCtx := context.Background()
		assert.Equal(t, "", getDetachedOrgID(emptyDetachedCtx))
		assert.Equal(t, "", getDetachedAppID(emptyDetachedCtx))
		assert.Equal(t, "", getDetachedKeyID(emptyDetachedCtx))
	})
}
