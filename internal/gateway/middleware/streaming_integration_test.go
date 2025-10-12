package middleware_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	dbdriver "github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/auth"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/middleware"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/policies"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestStreamingResponseWithUsageRecordingAndPolicies(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set up test containers
	pgConnStr, redisAddr := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	ctx := context.Background()

	// Set up database connection
	pg, err := dbdriver.NewPostgresDriver(ctx, pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	// Set up Redis cache
	cache, err := kv.NewDriver(kv.Config{
		Backend:   kv.BackendRedis,
		RedisAddr: redisAddr,
		RedisPW:   "",
		RedisDB:   0,
	})
	require.NoError(t, err)

	// Create policy engine
	engine := policies.NewEngine(pg.Queries, cache)

	// Create usage recorder
	usageRecorder := middleware.NewUsageRecorder(pg.Queries, engine)

	// Create DB fixtures helper
	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)

	// Create test data
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	testOrgID := orgID.String()
	testAppID := appID.String()

	// Create test API key
	testKeyID := fixtures.CreateTestAPIKey(t, orgID)

	// Create test policies (token limit policy to test post-checks)
	tokenLimitConfig := `{"max_completion_tokens": 100}`
	tokenLimitID := fixtures.CreateTestPolicy(t, orgID, appID, model.PolicyTypeTokenLimit, tokenLimitConfig)

	// Clean up
	defer func() {
		fixtures.CleanupTestPolicy(t, tokenLimitID)
		fixtures.CleanupTestApp(t, appID)
	}()

	t.Run("StreamingResponseCapturesUsageAndRunsPostChecks", func(t *testing.T) {
		// Create a streaming response that simulates OpenAI streaming format
		streamingChunks := []string{
			`data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-4","choices":[{"index":0,"delta":{"role":"assistant","content":"Hello"},"finish_reason":null}]}

`,
			`data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-4","choices":[{"index":0,"delta":{"content":" world"},"finish_reason":null}]}

`,
			`data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-4","choices":[{"index":0,"delta":{"content":"!"},"finish_reason":"stop"}],"usage":{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15}}

`,
			`data: [DONE]

`,
		}

		// Mock upstream server that returns streaming response
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")

			// Write streaming response
			for _, chunk := range streamingChunks {
				w.Write([]byte(chunk))
				w.(http.Flusher).Flush() // Simulate streaming
				time.Sleep(1 * time.Millisecond)
			}
		}))
		defer upstream.Close()

		// Create request with auth context
		req := httptest.NewRequest("POST", upstream.URL+"/v1/chat/completions", strings.NewReader(`{"model":"gpt-4","messages":[{"role":"user","content":"Hello"}]}`))
		req.Header.Set("Content-Type", "application/json")

		// Set up auth context (simulate authenticated request)
		ctx := auth.WithOrgID(req.Context(), testOrgID)
		ctx = auth.WithAppID(ctx, testAppID)
		ctx = auth.WithKeyID(ctx, testKeyID.String())
		ctx = auth.WithProvider(ctx, "openai")
		ctx = auth.WithModelName(ctx, "gpt-4")

		// Load policies and set in context
		policiesList, err := engine.LoadPolicies(ctx, testAppID)
		require.NoError(t, err)
		ctx = auth.WithPolicies(ctx, policiesList)

		req = req.WithContext(ctx)

		// Create middleware chain
		transport := usageRecorder.Middleware(http.DefaultTransport)

		// Execute request
		resp, err := transport.RoundTrip(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// Read the streaming response (simulates client consuming the stream)
		var responseBody bytes.Buffer
		_, err = io.Copy(&responseBody, resp.Body)
		require.NoError(t, err)
		resp.Body.Close()

		// Verify we got the streaming response
		responseStr := responseBody.String()
		require.Contains(t, responseStr, "data: [DONE]")
		require.Contains(t, responseStr, `"usage":{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15}`)

		// Wait for async processing to complete
		time.Sleep(50 * time.Millisecond)

		// Verify usage was recorded in database
		now := time.Now()
		oneHourAgo := now.Add(-time.Hour)
		usageMetrics, err := pg.Queries.GetUsageMetricsByApp(ctx, db.GetUsageMetricsByAppParams{
			AppID:       uuid.MustParse(testAppID),
			Timestamp:   pgtype.Timestamptz{Time: oneHourAgo, Valid: true},
			Timestamp_2: pgtype.Timestamptz{Time: now, Valid: true},
		})
		require.NoError(t, err)
		require.Len(t, usageMetrics, 1, "Should have recorded one usage metric")

		metric := usageMetrics[0]
		require.Equal(t, int32(10), metric.PromptTokens)
		require.Equal(t, int32(5), metric.CompletionTokens)
		require.Equal(t, int32(15), metric.TotalTokens)
		require.Equal(t, "openai", metric.Provider)
		require.Equal(t, "gpt-4", metric.ModelName)
	})

	t.Run("StreamingResponseWithLargePayload", func(t *testing.T) {
		// Create a larger streaming response to test memory handling
		var chunks []string
		totalContent := ""

		// Generate 50 chunks of content
		for i := 0; i < 50; i++ {
			content := fmt.Sprintf("Chunk %d content ", i)
			totalContent += content
			chunk := fmt.Sprintf(`data: {"id":"chatcmpl-large","object":"chat.completion.chunk","created":1677652288,"model":"gpt-4","choices":[{"index":0,"delta":{"content":"%s"},"finish_reason":null}]}

`, content)
			chunks = append(chunks, chunk)
		}

		// Final chunk with usage
		finalChunk := `data: {"id":"chatcmpl-large","object":"chat.completion.chunk","created":1677652288,"model":"gpt-4","choices":[{"index":0,"delta":{},"finish_reason":"stop"}],"usage":{"prompt_tokens":50,"completion_tokens":250,"total_tokens":300}}

`
		chunks = append(chunks, finalChunk)
		chunks = append(chunks, "data: [DONE]\n\n")

		streamingBody := strings.Join(chunks, "")

		// Mock upstream server
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(streamingBody))
		}))
		defer upstream.Close()

		// Create request
		req := httptest.NewRequest("POST", upstream.URL+"/v1/chat/completions", strings.NewReader(`{"model":"gpt-4","messages":[{"role":"user","content":"Large request"}]}`))

		// Set up auth context
		ctx := auth.WithOrgID(req.Context(), testOrgID)
		ctx = auth.WithAppID(ctx, testAppID)
		ctx = auth.WithKeyID(ctx, testKeyID.String())
		ctx = auth.WithProvider(ctx, "openai")
		ctx = auth.WithModelName(ctx, "gpt-4")

		req = req.WithContext(ctx)

		// Execute request
		transport := usageRecorder.Middleware(http.DefaultTransport)
		resp, err := transport.RoundTrip(req)
		require.NoError(t, err)

		// Read response
		var responseBody bytes.Buffer
		_, err = io.Copy(&responseBody, resp.Body)
		require.NoError(t, err)
		resp.Body.Close()

		// Wait for async processing
		time.Sleep(50 * time.Millisecond)

		// Verify usage recording
		now := time.Now()
		oneHourAgo := now.Add(-time.Hour)
		usageMetrics, err := pg.Queries.GetUsageMetricsByApp(ctx, db.GetUsageMetricsByAppParams{
			AppID:       uuid.MustParse(testAppID),
			Timestamp:   pgtype.Timestamptz{Time: oneHourAgo, Valid: true},
			Timestamp_2: pgtype.Timestamptz{Time: now, Valid: true},
		})
		require.NoError(t, err)

		// Should have 2 metrics now (one from previous test + one from this test)
		require.Len(t, usageMetrics, 2)

		// Find the latest metric (large payload)
		var latestMetric db.UsageMetric
		for _, metric := range usageMetrics {
			if metric.TotalTokens == 300 {
				latestMetric = metric
				break
			}
		}

		require.Equal(t, int32(50), latestMetric.PromptTokens)
		require.Equal(t, int32(250), latestMetric.CompletionTokens)
		require.Equal(t, int32(300), latestMetric.TotalTokens)
	})

	t.Run("StreamingResponseWithPolicyPostChecks", func(t *testing.T) {
		// Create a streaming response that will trigger post-check logging
		streamingChunks := []string{
			`data: {"id":"chatcmpl-policy","object":"chat.completion.chunk","created":1677652288,"model":"gpt-4","choices":[{"index":0,"delta":{"role":"assistant","content":"Test"},"finish_reason":null}]}

`,
			`data: {"id":"chatcmpl-policy","object":"chat.completion.chunk","created":1677652288,"model":"gpt-4","choices":[{"index":0,"delta":{},"finish_reason":"stop"}],"usage":{"prompt_tokens":5,"completion_tokens":150,"total_tokens":155}}

`,
			`data: [DONE]

`,
		}

		streamingBody := strings.Join(streamingChunks, "")

		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(streamingBody))
		}))
		defer upstream.Close()

		req := httptest.NewRequest("POST", upstream.URL+"/v1/chat/completions", strings.NewReader(`{"model":"gpt-4","messages":[{"role":"user","content":"Test"}]}`))

		// Set up auth context
		ctx := auth.WithOrgID(req.Context(), testOrgID)
		ctx = auth.WithAppID(ctx, testAppID)
		ctx = auth.WithKeyID(ctx, testKeyID.String())
		ctx = auth.WithProvider(ctx, "openai")
		ctx = auth.WithModelName(ctx, "gpt-4")

		// Load policies (includes our token limit policy)
		policiesList, err := engine.LoadPolicies(ctx, testAppID)
		require.NoError(t, err)
		ctx = auth.WithPolicies(ctx, policiesList)

		req = req.WithContext(ctx)

		// Execute request
		transport := usageRecorder.Middleware(http.DefaultTransport)
		resp, err := transport.RoundTrip(req)
		require.NoError(t, err)

		// Read response
		var responseBody bytes.Buffer
		_, err = io.Copy(&responseBody, resp.Body)
		require.NoError(t, err)
		resp.Body.Close()

		// Wait for async processing (post-checks run here)
		time.Sleep(50 * time.Millisecond)

		// Verify usage was recorded
		now := time.Now()
		oneHourAgo := now.Add(-time.Hour)
		usageMetrics, err := pg.Queries.GetUsageMetricsByApp(ctx, db.GetUsageMetricsByAppParams{
			AppID:       uuid.MustParse(testAppID),
			Timestamp:   pgtype.Timestamptz{Time: oneHourAgo, Valid: true},
			Timestamp_2: pgtype.Timestamptz{Time: now, Valid: true},
		})
		require.NoError(t, err)
		require.True(t, len(usageMetrics) >= 1, "Should have recorded usage metrics")

		// The token limit policy post-check should have logged about exceeding completion tokens
		// (150 completion tokens > 100 limit), but the request should still succeed
		// since post-checks are for monitoring only
	})

	t.Run("StreamingResponseWithoutUsageData", func(t *testing.T) {
		// Test streaming response without usage data (some providers don't include it)
		streamingChunks := []string{
			`data: {"id":"chatcmpl-no-usage","object":"chat.completion.chunk","created":1677652288,"model":"gpt-4","choices":[{"index":0,"delta":{"role":"assistant","content":"No usage"},"finish_reason":"stop"}]}

`,
			`data: [DONE]

`,
		}

		streamingBody := strings.Join(streamingChunks, "")

		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(streamingBody))
		}))
		defer upstream.Close()

		req := httptest.NewRequest("POST", upstream.URL+"/v1/chat/completions", strings.NewReader(`{"model":"gpt-4","messages":[{"role":"user","content":"No usage"}]}`))

		ctx := auth.WithOrgID(req.Context(), testOrgID)
		ctx = auth.WithAppID(ctx, testAppID)
		ctx = auth.WithKeyID(ctx, testKeyID.String())
		ctx = auth.WithProvider(ctx, "openai")
		ctx = auth.WithModelName(ctx, "gpt-4")

		req = req.WithContext(ctx)

		transport := usageRecorder.Middleware(http.DefaultTransport)
		resp, err := transport.RoundTrip(req)
		require.NoError(t, err)

		var responseBody bytes.Buffer
		_, err = io.Copy(&responseBody, resp.Body)
		require.NoError(t, err)
		resp.Body.Close()

		// Wait for async processing
		time.Sleep(50 * time.Millisecond)

		// Verify usage was recorded with zero values
		now := time.Now()
		oneHourAgo := now.Add(-time.Hour)
		usageMetrics, err := pg.Queries.GetUsageMetricsByApp(ctx, db.GetUsageMetricsByAppParams{
			AppID:       uuid.MustParse(testAppID),
			Timestamp:   pgtype.Timestamptz{Time: oneHourAgo, Valid: true},
			Timestamp_2: pgtype.Timestamptz{Time: now, Valid: true},
		})
		require.NoError(t, err)

		// Find the most recent metric
		var latestMetric db.UsageMetric
		var latestTime time.Time
		for _, metric := range usageMetrics {
			if metric.Timestamp.Time.After(latestTime) {
				latestMetric = metric
				latestTime = metric.Timestamp.Time
			}
		}

		// Should have zero token counts when no usage data is available
		require.Equal(t, int32(0), latestMetric.PromptTokens)
		require.Equal(t, int32(0), latestMetric.CompletionTokens)
		require.Equal(t, int32(0), latestMetric.TotalTokens)
	})
}
