package usage

import (
	"context"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
)

type Reader interface {
	GetByAppID(ctx context.Context, appID uuid.UUID, start, end time.Time) ([]*model.UsageMetric, error)
	GetByOrgID(ctx context.Context, orgID uuid.UUID, start, end time.Time) ([]*model.UsageMetric, error)
	GetByAPIKeyID(ctx context.Context, apiKeyID uuid.UUID, start, end time.Time) ([]*model.UsageMetric, error)
	SumTokensByAppID(ctx context.Context, appID uuid.UUID, start, end time.Time) (*TokenSummary, error)
	SumTokensByOrgID(ctx context.Context, orgID uuid.UUID, start, end time.Time) (*TokenSummary, error)
	GetUsageByModel(ctx context.Context, appID uuid.UUID, start, end time.Time) ([]*ModelUsageSummary, error)
}

type Writer interface {
	Create(ctx context.Context, metric *model.UsageMetric) error
}

type Repository interface {
	Reader
	Writer
}

type TokenSummary struct {
	TotalPromptTokens     int
	TotalCompletionTokens int
	TotalTokens           int
	RequestCount          int
}

type ModelUsageSummary struct {
	ModelName             string
	Provider              string
	TotalPromptTokens     int
	TotalCompletionTokens int
	TotalTokens           int
	RequestCount          int
}
