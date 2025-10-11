package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/auth"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/tokens"
)

// Benchmark helpers
func benchmarkRequest(size int) *http.Request {
	body := makeLLMRequest(size)
	req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Populate context with auth data
	ctx := req.Context()
	ctx = auth.WithOrgID(ctx, "org-123")
	ctx = auth.WithAppID(ctx, "app-456")
	ctx = auth.WithKeyID(ctx, "key-789")
	req = req.WithContext(ctx)

	return req
}

func makeLLMRequest(approxSize int) []byte {
	// Create a realistic LLM request of approximately 'size' bytes
	message := ""
	for len(message) < approxSize {
		message += "This is a test message for benchmarking. "
	}

	return []byte(`{
		"model": "gpt-4",
		"messages": [
			{"role": "system", "content": "You are a helpful assistant."},
			{"role": "user", "content": "` + message[:min(approxSize-200, len(message))] + `"}
		]
	}`)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Mock RoundTripper that returns immediately
type mockRoundTripper struct {
	responseBody []byte
}

func (m *mockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewReader(m.responseBody)),
	}, nil
}

var mockLLMResponse = []byte(`{
	"id": "chatcmpl-123",
	"object": "chat.completion",
	"created": 1677652288,
	"model": "gpt-4",
	"choices": [{
		"index": 0,
		"message": {
			"role": "assistant",
			"content": "Hello! How can I help you today?"
		},
		"finish_reason": "stop"
	}],
	"usage": {
		"prompt_tokens": 56,
		"completion_tokens": 9,
		"total_tokens": 65
	}
}`)

// Benchmark: Request buffering
func BenchmarkRequestBuffer(b *testing.B) {
	sizes := []int{1024, 10240, 51200} // 1KB, 10KB, 50KB

	for _, size := range sizes {
		b.Run(string(rune(size/1024))+"KB", func(b *testing.B) {
			buffer := NewRequestBuffer()
			next := &mockRoundTripper{responseBody: mockLLMResponse}
			middleware := buffer.Middleware(next)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				req := benchmarkRequest(size)
				_, err := middleware.RoundTrip(req)
				if err != nil {
					b.Fatal(err)
				}
			}

			b.ReportMetric(float64(size), "request_bytes")
		})
	}
}

// Benchmark: Request parsing (new implementation)
func BenchmarkRequestParsing(b *testing.B) {
	sizes := []int{1024, 10240, 51200}

	for _, size := range sizes {
		b.Run(string(rune(size/1024))+"KB", func(b *testing.B) {
			bodyBytes := makeLLMRequest(size)
			buffer := NewRequestBuffer()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				parsed := buffer.parseRequest(bodyBytes)
				if parsed.Model != "gpt-4" {
					b.Fatalf("expected gpt-4, got %s", parsed.Model)
				}
			}
		})
	}
}

// Benchmark: Token estimation
func BenchmarkTokenEstimation(b *testing.B) {
	estimator := tokens.NewEstimator()
	sizes := []int{1024, 10240, 51200}

	for _, size := range sizes {
		b.Run(string(rune(size/1024))+"KB", func(b *testing.B) {
			bodyBytes := makeLLMRequest(size)
			model := "gpt-4"

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_, err := estimator.EstimateRequest(model, bodyBytes)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// Benchmark: Full policy enforcement middleware (no policies)
func BenchmarkPolicyEnforcementEmpty(b *testing.B) {
	// TODO: Implement after we refactor policy engine
	b.Skip("Requires policy engine with in-memory cache")
}

// Benchmark: Full policy enforcement middleware (with rate limit)
func BenchmarkPolicyEnforcementWithRateLimit(b *testing.B) {
	// TODO: Implement after we refactor policy engine
	b.Skip("Requires policy engine with in-memory cache")
}

// Benchmark: Context value access patterns
func BenchmarkContextAccess(b *testing.B) {
	b.Run("SingleValue", func(b *testing.B) {
		ctx := context.Background()
		ctx = auth.WithAppID(ctx, "app-123")

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = auth.GetAppID(ctx)
		}
	})

	b.Run("MultipleValues", func(b *testing.B) {
		ctx := context.Background()
		ctx = auth.WithOrgID(ctx, "org-123")
		ctx = auth.WithAppID(ctx, "app-456")
		ctx = auth.WithKeyID(ctx, "key-789")
		ctx = auth.WithProvider(ctx, "openai")
		ctx = auth.WithModelName(ctx, "gpt-4")

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = auth.GetOrgID(ctx)
			_ = auth.GetAppID(ctx)
			_ = auth.GetKeyID(ctx)
			_ = auth.GetProvider(ctx)
			_ = auth.GetModelName(ctx)
		}
	})

	b.Run("WithParsedRequest", func(b *testing.B) {
		ctx := context.Background()
		parsedReq := &auth.ParsedRequest{
			Model:           "gpt-4",
			EstimatedTokens: 1234,
			RequestSize:     51200,
		}
		ctx = auth.WithParsedRequest(ctx, parsedReq)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req := auth.GetParsedRequest(ctx)
			if req == nil || req.Model != "gpt-4" {
				b.Fatal("parsed request mismatch")
			}
		}
	})
}

// Benchmark: JSON unmarshaling (to show cost of multiple unmarshals)
func BenchmarkJSONUnmarshal(b *testing.B) {
	sizes := []int{1024, 10240, 51200}

	for _, size := range sizes {
		b.Run(string(rune(size/1024))+"KB", func(b *testing.B) {
			bodyBytes := makeLLMRequest(size)

			b.Run("Once", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					var req struct {
						Model    string `json:"model"`
						Messages []struct {
							Role    string `json:"role"`
							Content string `json:"content"`
						} `json:"messages"`
					}
					if err := json.Unmarshal(bodyBytes, &req); err != nil {
						b.Fatal(err)
					}
				}
			})

			b.Run("ThreeTimes", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					// Simulate current behavior: unmarshal 3 times
					var req1 struct{ Model string `json:"model"` }
					json.Unmarshal(bodyBytes, &req1)

					var req2 struct {
						Messages []struct {
							Content string `json:"content"`
						} `json:"messages"`
					}
					json.Unmarshal(bodyBytes, &req2)

					var req3 struct{ Model string `json:"model"` }
					json.Unmarshal(bodyBytes, &req3)
				}
			})
		})
	}
}
