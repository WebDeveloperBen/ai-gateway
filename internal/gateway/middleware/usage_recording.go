package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/auth"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/policies"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/tokens"
	"github.com/WebDeveloperBen/ai-gateway/internal/logger"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

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

		// Wrap response body with TeeReader to capture bytes as they stream
		// This avoids buffering the entire response before returning
		var capturedBytes bytes.Buffer
		if resp.Body != nil {
			resp.Body = io.NopCloser(io.TeeReader(resp.Body, &capturedBytes))
		}

		// Create detached context for async processing
		detachedCtx := detachContext(ctx)

		// Launch async goroutine for post-processing
		// This goroutine will read from capturedBytes AFTER the response is consumed
		go ur.recordAsync(detachedCtx, &asyncRecordParams{
			provider:         provider,
			modelName:        modelName,
			requestSizeBytes: requestSizeBytes,
			latencyMs:        latencyMs,
			request:          r,
			response:         resp,
			capturedBytes:    &capturedBytes,
		})

		// Return response immediately - streaming starts now
		// The goroutine will process bytes as the client reads them
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
	tokenUsage, err := ur.parser.ParseResponse(params.provider, respBodyBytes)
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

	// Load policies and run post-checks
	policyList, err := ur.engine.LoadPolicies(ctx, appID)
	if err != nil {
		logger.GetLogger(ctx).Error().
			Err(err).
			Str("app_id", appID).
			Msg("Failed to load policies for post-check")
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

	// Run post-checks (non-blocking, for logging/metrics)
	for _, policy := range policyList {
		policy.PostCheck(ctx, postCtx)
	}
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
