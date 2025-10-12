package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/auth"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/policies"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/tokens"
	"github.com/WebDeveloperBen/ai-gateway/internal/logger"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/observability"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// responseBodyWrapper wraps the response body to capture the full content before async processing
type responseBodyWrapper struct {
	originalBody io.ReadCloser
	onClose      func([]byte)
	body         *bytes.Buffer
	closed       bool
}

// Read implements io.Reader for the response body
func (rbw *responseBodyWrapper) Read(p []byte) (n int, err error) {
	if rbw.closed {
		return 0, http.ErrBodyReadAfterClose
	}
	n, err = rbw.originalBody.Read(p)
	if n > 0 {
		rbw.body.Write(p[:n])
	}
	if err != nil {
		rbw.close()
	}
	return n, err
}

// Close implements io.Closer for the response body
func (rbw *responseBodyWrapper) Close() error {
	if rbw.closed {
		return nil
	}
	rbw.close()
	return rbw.originalBody.Close()
}

// close captures the full body and calls the onClose callback
func (rbw *responseBodyWrapper) close() {
	if rbw.closed {
		return
	}
	rbw.closed = true
	bodyBytes := rbw.body.Bytes()
	rbw.onClose(bodyBytes)
}

// UsageRecorder records usage metrics and runs post-check policies
type UsageRecorder struct {
	db     *db.Queries
	parser *tokens.Parser
	engine *policies.Engine
}

// NewUsageRecorder creates a new usage recorder
func NewUsageRecorder(database *db.Queries, engine *policies.Engine) *UsageRecorder {
	return &UsageRecorder{
		db:     database,
		parser: tokens.NewParser(),
		engine: engine,
	}
}

// Middleware returns a RoundTripper middleware that records usage asynchronously
func (ur *UsageRecorder) Middleware(next http.RoundTripper) http.RoundTripper {
	return roundTripFunc(func(r *http.Request) (*http.Response, error) {
		// Capture request metadata
		ctx := r.Context()

		// Get request size from parsed request data
		parsedReq := auth.GetParsedRequest(ctx)
		requestSizeBytes := 0
		if parsedReq != nil {
			requestSizeBytes = parsedReq.RequestSize
		}

		// Capture start time
		startTime := time.Now()

		// Execute request
		resp, err := next.RoundTrip(r)
		if err != nil {
			return resp, err
		}

		// Capture latency
		latencyMs := time.Since(startTime).Milliseconds()

		// Extract provider and model from context
		provider := auth.GetProvider(ctx)
		modelName := auth.GetModelName(ctx)

		// Create detached context for async processing
		detachedCtx := detachContext(ctx)

		// Wrap the response body to capture the full content before async processing
		resp.Body = &responseBodyWrapper{
			originalBody: resp.Body,
			body:         &bytes.Buffer{},
			onClose: func(bodyBytes []byte) {
				// Launch async processing after response is fully consumed
				go ur.recordAsync(detachedCtx, &asyncRecordParams{
					provider:         provider,
					modelName:        modelName,
					requestSizeBytes: requestSizeBytes,
					latencyMs:        latencyMs,
					request:          r,
					response:         resp,
					capturedBytes:    bytes.NewBuffer(bodyBytes),
				})
			},
		}

		// Return the response with wrapped body
		return resp, nil
	})
}

// asyncRecordParams holds parameters for async recording
type asyncRecordParams struct {
	provider         string
	modelName        string
	requestSizeBytes int
	latencyMs        int64
	request          *http.Request
	response         *http.Response
	capturedBytes    *bytes.Buffer
}

// recordAsync performs the actual recording in a goroutine
// This runs AFTER the response has been consumed by the client
func (ur *UsageRecorder) recordAsync(ctx context.Context, params *asyncRecordParams) {
	// Wait a bit to ensure response body has been consumed
	// This allows TeeReader to capture the full response
	time.Sleep(10 * time.Millisecond)

	// Extract IDs from detached context
	orgID := getDetachedOrgID(ctx)
	appID := getDetachedAppID(ctx)
	apiKeyID := getDetachedKeyID(ctx)

	// Get captured response bytes
	respBodyBytes := params.capturedBytes.Bytes()
	responseSizeBytes := len(respBodyBytes)

	// Parse token usage from captured response
	tokenUsage, err := ur.parseTokenUsage(params.provider, respBodyBytes)
	if err != nil {
		// If parsing fails, use zeros (some responses may not have usage)
		tokenUsage = &model.TokenUsage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		}
	}

	// Parse UUIDs
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		logger.GetLogger(ctx).Error().
			Err(err).
			Str("org_id", orgID).
			Msg("Failed to parse org ID for usage metric")
		return
	}

	appUUID, err := uuid.Parse(appID)
	if err != nil {
		logger.GetLogger(ctx).Error().
			Err(err).
			Str("app_id", appID).
			Msg("Failed to parse app ID for usage metric")
		return
	}

	apiKeyUUID, err := uuid.Parse(apiKeyID)
	if err != nil {
		logger.GetLogger(ctx).Error().
			Err(err).
			Str("api_key_id", apiKeyID).
			Msg("Failed to parse API key ID for usage metric")
		return
	}

	// Insert usage metric
	_, err = ur.db.CreateUsageMetric(ctx, db.CreateUsageMetricParams{
		OrgID:             orgUUID,
		AppID:             appUUID,
		ApiKeyID:          apiKeyUUID,
		ModelID:           nil,
		Provider:          params.provider,
		ModelName:         params.modelName,
		PromptTokens:      int32(tokenUsage.PromptTokens),
		CompletionTokens:  int32(tokenUsage.CompletionTokens),
		TotalTokens:       int32(tokenUsage.TotalTokens),
		RequestSizeBytes:  int32(params.requestSizeBytes),
		ResponseSizeBytes: int32(responseSizeBytes),
		Timestamp: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		logger.GetLogger(ctx).Error().
			Err(err).
			Str("org_id", orgID).
			Str("app_id", appID).
			Str("provider", params.provider).
			Str("model", params.modelName).
			Int("total_tokens", tokenUsage.TotalTokens).
			Msg("Failed to create usage metric")
		return
	}

	// Get policies from context (already loaded in enforcement middleware)
	policyListInterface := auth.GetPolicies(ctx)
	if policyListInterface == nil {
		// Fallback: load policies if not in context
		// This shouldn't normally happen but provides safety
		var err error
		policyListInterface, err = ur.engine.LoadPolicies(ctx, appID)
		if err != nil {
			logger.GetLogger(ctx).Error().
				Err(err).
				Str("app_id", appID).
				Msg("Failed to load policies for post-check")
			return
		}
	}

	// Type assert to policy list
	policyList, ok := policyListInterface.([]policies.Policy)
	if !ok {
		logger.GetLogger(ctx).Error().
			Str("app_id", appID).
			Msg("Invalid policy list type in context")
		return
	}

	// Build post-request context
	postCtx := &policies.PostRequestContext{
		Request:           params.request,
		Response:          params.response,
		OrgID:             orgID,
		AppID:             appID,
		APIKeyID:          apiKeyID,
		Provider:          params.provider,
		ModelName:         params.modelName,
		ActualTokens:      *tokenUsage,
		RequestSizeBytes:  params.requestSizeBytes,
		ResponseSizeBytes: responseSizeBytes,
		LatencyMs:         params.latencyMs,
	}

	// Record LLM token metrics
	observability.FromContext(ctx).RecordLLMTokens(
		ctx,
		params.provider,
		params.modelName,
		tokenUsage.PromptTokens,
		tokenUsage.CompletionTokens,
		tokenUsage.TotalTokens,
	)

	// Run post-checks (non-blocking, for logging/metrics)
	for _, policy := range policyList {
		policy.PostCheck(ctx, postCtx)
	}
}

// parseTokenUsage handles both regular and streaming responses
func (ur *UsageRecorder) parseTokenUsage(provider string, respBodyBytes []byte) (*model.TokenUsage, error) {
	// First try parsing as a regular response
	usage, err := ur.parser.ParseResponse(provider, respBodyBytes)
	if err == nil {
		return usage, nil
	}

	// If that fails, try parsing as streaming response
	// Split by SSE event boundaries
	bodyStr := string(respBodyBytes)
	lines := strings.Split(bodyStr, "\n")

	var chunks [][]byte
	var currentChunk []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			// Empty line indicates end of chunk
			if len(currentChunk) > 0 {
				chunkData := strings.Join(currentChunk, "\n")
				if strings.HasPrefix(chunkData, "data: ") {
					// Extract JSON part after "data: "
					jsonData := strings.TrimPrefix(chunkData, "data: ")
					if jsonData != "[DONE]" {
						chunks = append(chunks, []byte(jsonData))
					}
				}
				currentChunk = nil
			}
		} else {
			currentChunk = append(currentChunk, line)
		}
	}

	// Handle any remaining chunk
	if len(currentChunk) > 0 {
		chunkData := strings.Join(currentChunk, "\n")
		if strings.HasPrefix(chunkData, "data: ") {
			jsonData := strings.TrimPrefix(chunkData, "data: ")
			if jsonData != "[DONE]" {
				chunks = append(chunks, []byte(jsonData))
			}
		}
	}

	// Try parsing chunks for streaming usage
	return ur.parser.ParseStreamedResponse(provider, chunks)
}

type detachedContextKey struct{}

type detachedData struct {
	OrgID     string
	AppID     string
	KeyID     string
	UserID    string
	Provider  string
	ModelName string
}

// detachContext creates a new context that won't be canceled when the parent is
// Uses a single context value to reduce allocations (1 instead of 6)
func detachContext(parent context.Context) context.Context {
	data := &detachedData{
		OrgID:     auth.GetOrgID(parent),
		AppID:     auth.GetAppID(parent),
		KeyID:     auth.GetKeyID(parent),
		UserID:    auth.GetUserID(parent),
		Provider:  auth.GetProvider(parent),
		ModelName: auth.GetModelName(parent),
	}
	return context.WithValue(context.Background(), detachedContextKey{}, data)
}

// Helper functions to get values from detached context
func getDetachedOrgID(ctx context.Context) string {
	if data, ok := ctx.Value(detachedContextKey{}).(*detachedData); ok {
		return data.OrgID
	}
	return auth.GetOrgID(ctx)
}

func getDetachedAppID(ctx context.Context) string {
	if data, ok := ctx.Value(detachedContextKey{}).(*detachedData); ok {
		return data.AppID
	}
	return auth.GetAppID(ctx)
}

func getDetachedKeyID(ctx context.Context) string {
	if data, ok := ctx.Value(detachedContextKey{}).(*detachedData); ok {
		return data.KeyID
	}
	return auth.GetKeyID(ctx)
}
