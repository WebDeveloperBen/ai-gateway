package usage

import (
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

type UsageMetric struct {
	ID                string    `json:"id"`
	OrgID             string    `json:"org_id"`
	AppID             string    `json:"app_id"`
	APIKeyID          string    `json:"api_key_id"`
	ModelID           *string   `json:"model_id,omitempty"`
	Provider          string    `json:"provider"`
	ModelName         string    `json:"model_name"`
	PromptTokens      int       `json:"prompt_tokens"`
	CompletionTokens  int       `json:"completion_tokens"`
	TotalTokens       int       `json:"total_tokens"`
	RequestSizeBytes  int       `json:"request_size_bytes"`
	ResponseSizeBytes int       `json:"response_size_bytes"`
	Timestamp         time.Time `json:"timestamp"`
}

type TokenSummary struct {
	TotalPromptTokens     int `json:"total_prompt_tokens"`
	TotalCompletionTokens int `json:"total_completion_tokens"`
	TotalTokens           int `json:"total_tokens"`
	RequestCount          int `json:"request_count"`
}

type ModelUsageSummary struct {
	ModelName             string `json:"model_name"`
	Provider              string `json:"provider"`
	TotalPromptTokens     int    `json:"total_prompt_tokens"`
	TotalCompletionTokens int    `json:"total_completion_tokens"`
	TotalTokens           int    `json:"total_tokens"`
	RequestCount          int    `json:"request_count"`
}

// Huma request/response types
type GetUsageMetricsRequest struct {
	AppID string `query:"app_id" required:"true"`
	Start string `query:"start" required:"true"` // ISO 8601 format
	End   string `query:"end" required:"true"`   // ISO 8601 format
	model.ListRequest
}

type GetUsageMetricsResponse struct {
	Body []*UsageMetric `json:"body"`
}

type GetUsageSummaryByAppRequest struct {
	AppID string `path:"app_id" required:"true"`
	Start string `query:"start" required:"true"`
	End   string `query:"end" required:"true"`
}

type GetUsageSummaryByAppResponse struct {
	Body *TokenSummary `json:"body"`
}

type GetUsageSummaryByOrgRequest struct {
	Start string `query:"start" required:"true"`
	End   string `query:"end" required:"true"`
}

type GetUsageSummaryByOrgResponse struct {
	Body *TokenSummary `json:"body"`
}

type GetUsageByModelRequest struct {
	AppID string `path:"app_id" required:"true"`
	Start string `query:"start" required:"true"`
	End   string `query:"end" required:"true"`
	model.ListRequest
}

type GetUsageByModelResponse struct {
	Body []*ModelUsageSummary `json:"body"`
}
