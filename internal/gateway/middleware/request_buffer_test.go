package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestBuffer_NewRequestBuffer(t *testing.T) {
	buffer := NewRequestBuffer()
	assert.NotNil(t, buffer)
	assert.NotNil(t, buffer.estimator)
}

func TestRequestBuffer_Middleware(t *testing.T) {
	buffer := NewRequestBuffer()

	t.Run("EmptyBody", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte{}))
		req.Header.Set("Content-Type", "application/json")

		var receivedParsed *auth.ParsedRequest
		next := &testMockRoundTripper{
			responseBody: []byte(`{"status": "ok"}`),
			checkRequest: func(r *http.Request) {
				receivedParsed = auth.GetParsedRequest(r.Context())
			},
		}
		middleware := buffer.Middleware(next)

		resp, err := middleware.RoundTrip(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode)

		// Check that parsed request was stored in context
		assert.NotNil(t, receivedParsed)
		assert.Equal(t, 0, receivedParsed.RequestSize)
		assert.Empty(t, receivedParsed.Model)
	})

	t.Run("ValidChatCompletionRequest", func(t *testing.T) {
		body := `{
			"model": "gpt-4",
			"messages": [
				{"role": "system", "content": "You are helpful"},
				{"role": "user", "content": "Hello"}
			]
		}`
		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		var receivedParsed *auth.ParsedRequest
		next := &testMockRoundTripper{
			responseBody: []byte(`{"status": "ok"}`),
			checkRequest: func(r *http.Request) {
				receivedParsed = auth.GetParsedRequest(r.Context())
			},
		}
		middleware := buffer.Middleware(next)

		resp, err := middleware.RoundTrip(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Check that parsed request was stored in context
		assert.NotNil(t, receivedParsed)
		assert.Equal(t, len(body), receivedParsed.RequestSize)
		assert.Equal(t, "gpt-4", receivedParsed.Model)
		assert.Len(t, receivedParsed.Messages, 2)
		assert.Equal(t, "system", receivedParsed.Messages[0].Role)
		assert.Equal(t, "You are helpful", receivedParsed.Messages[0].Content)
		assert.Equal(t, "user", receivedParsed.Messages[1].Role)
		assert.Equal(t, "Hello", receivedParsed.Messages[1].Content)
		assert.Greater(t, receivedParsed.EstimatedTokens, 0)
	})

	t.Run("ValidCompletionRequest", func(t *testing.T) {
		body := `{
			"model": "text-davinci-003",
			"prompt": "Write a haiku about coding"
		}`
		req := httptest.NewRequest("POST", "/v1/completions", bytes.NewReader([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		var receivedParsed *auth.ParsedRequest
		next := &testMockRoundTripper{
			responseBody: []byte(`{"status": "ok"}`),
			checkRequest: func(r *http.Request) {
				receivedParsed = auth.GetParsedRequest(r.Context())
			},
		}
		middleware := buffer.Middleware(next)

		resp, err := middleware.RoundTrip(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Check that parsed request was stored in context
		assert.NotNil(t, receivedParsed)
		assert.Equal(t, len(body), receivedParsed.RequestSize)
		assert.Equal(t, "text-davinci-003", receivedParsed.Model)
		assert.Equal(t, "Write a haiku about coding", receivedParsed.Prompt)
		assert.Greater(t, receivedParsed.EstimatedTokens, 0)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		body := `{"model": "gpt-4", "messages": invalid json}`
		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		var receivedParsed *auth.ParsedRequest
		next := &testMockRoundTripper{
			responseBody: []byte(`{"status": "ok"}`),
			checkRequest: func(r *http.Request) {
				receivedParsed = auth.GetParsedRequest(r.Context())
			},
		}
		middleware := buffer.Middleware(next)

		resp, err := middleware.RoundTrip(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Check that parsed request was stored in context with partial data
		assert.NotNil(t, receivedParsed)
		assert.Equal(t, len(body), receivedParsed.RequestSize)
		assert.Empty(t, receivedParsed.Model) // Should be empty due to parse failure
	})

	t.Run("BodyReadError", func(t *testing.T) {
		// Create a reader that will fail
		req := httptest.NewRequest("POST", "/v1/chat/completions", &failingReader{})
		req.Header.Set("Content-Type", "application/json")

		next := &testMockRoundTripper{
			responseBody: []byte(`{"status": "ok"}`),
		}
		middleware := buffer.Middleware(next)

		resp, err := middleware.RoundTrip(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestRequestBuffer_parseRequest(t *testing.T) {
	buffer := NewRequestBuffer()
	ctx := context.Background()

	t.Run("EmptyBody", func(t *testing.T) {
		parsed := buffer.parseRequest(ctx, []byte{})
		assert.NotNil(t, parsed)
		assert.Equal(t, 0, parsed.RequestSize)
		assert.Empty(t, parsed.Model)
		assert.Nil(t, parsed.Messages)
		assert.Empty(t, parsed.Prompt)
		assert.Equal(t, 0, parsed.EstimatedTokens)
	})

	t.Run("ChatCompletionRequest", func(t *testing.T) {
		body := `{
			"model": "gpt-4",
			"messages": [
				{"role": "user", "content": "Hello world"}
			]
		}`
		parsed := buffer.parseRequest(ctx, []byte(body))
		assert.NotNil(t, parsed)
		assert.Equal(t, len(body), parsed.RequestSize)
		assert.Equal(t, "gpt-4", parsed.Model)
		assert.Len(t, parsed.Messages, 1)
		assert.Equal(t, "user", parsed.Messages[0].Role)
		assert.Equal(t, "Hello world", parsed.Messages[0].Content)
		assert.Greater(t, parsed.EstimatedTokens, 0)
	})

	t.Run("CompletionRequest", func(t *testing.T) {
		body := `{
			"model": "text-davinci-003",
			"prompt": "Complete this text"
		}`
		parsed := buffer.parseRequest(ctx, []byte(body))
		assert.NotNil(t, parsed)
		assert.Equal(t, len(body), parsed.RequestSize)
		assert.Equal(t, "text-davinci-003", parsed.Model)
		assert.Equal(t, "Complete this text", parsed.Prompt)
		assert.Greater(t, parsed.EstimatedTokens, 0)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		body := `{"invalid": json}`
		parsed := buffer.parseRequest(ctx, []byte(body))
		assert.NotNil(t, parsed)
		assert.Equal(t, len(body), parsed.RequestSize)
		assert.Empty(t, parsed.Model)
		assert.Nil(t, parsed.Messages)
		assert.Empty(t, parsed.Prompt)
		assert.Equal(t, 0, parsed.EstimatedTokens)
	})

	t.Run("ComplexMessages", func(t *testing.T) {
		body := `{
			"model": "gpt-4",
			"messages": [
				{"role": "system", "content": "You are a helpful assistant"},
				{"role": "user", "content": "What is 2+2?"},
				{"role": "assistant", "content": "4"},
				{"role": "user", "content": "Thanks!"}
			]
		}`
		parsed := buffer.parseRequest(ctx, []byte(body))
		assert.NotNil(t, parsed)
		assert.Equal(t, "gpt-4", parsed.Model)
		assert.Len(t, parsed.Messages, 4)
		assert.Equal(t, "system", parsed.Messages[0].Role)
		assert.Equal(t, "user", parsed.Messages[1].Role)
		assert.Equal(t, "assistant", parsed.Messages[2].Role)
		assert.Equal(t, "user", parsed.Messages[3].Role)
	})
}

// Mock RoundTripper for testing
type testMockRoundTripper struct {
	responseBody []byte
	checkRequest func(*http.Request)
}

func (m *testMockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if m.checkRequest != nil {
		m.checkRequest(r)
	}
	return &http.Response{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewReader(m.responseBody)),
	}, nil
}

// failingReader implements io.Reader that always fails
type failingReader struct{}

func (failingReader) Read(p []byte) (n int, err error) {
	return 0, assert.AnError
}
