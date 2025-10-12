package usage

import (
	"context"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/repository/usage"
	"github.com/google/uuid"
)

type UsageService interface {
	GetMetricsByAppID(ctx context.Context, appID uuid.UUID, start, end time.Time, limit, offset int) ([]*UsageMetric, error)
	GetMetricsByOrgID(ctx context.Context, orgID uuid.UUID, start, end time.Time) ([]*UsageMetric, error)
	GetMetricsByAPIKeyID(ctx context.Context, apiKeyID uuid.UUID, start, end time.Time) ([]*UsageMetric, error)
	GetTokenSummaryByAppID(ctx context.Context, appID uuid.UUID, start, end time.Time) (*TokenSummary, error)
	GetTokenSummaryByOrgID(ctx context.Context, orgID uuid.UUID, start, end time.Time) (*TokenSummary, error)
	GetUsageByModel(ctx context.Context, appID uuid.UUID, start, end time.Time, limit, offset int) ([]*ModelUsageSummary, error)
}

type usageService struct {
	repo usage.Repository
}

func NewService(repo usage.Repository) UsageService {
	return &usageService{repo: repo}
}

func (s *usageService) GetMetricsByAppID(ctx context.Context, appID uuid.UUID, start, end time.Time, limit, offset int) ([]*UsageMetric, error) {
	metrics, err := s.repo.GetByAppID(ctx, appID, start, end, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*UsageMetric, len(metrics))
	for i, metric := range metrics {
		result[i] = s.convertToAPI(metric)
	}
	return result, nil
}

func (s *usageService) GetMetricsByOrgID(ctx context.Context, orgID uuid.UUID, start, end time.Time) ([]*UsageMetric, error) {
	metrics, err := s.repo.GetByOrgID(ctx, orgID, start, end)
	if err != nil {
		return nil, err
	}

	result := make([]*UsageMetric, len(metrics))
	for i, metric := range metrics {
		result[i] = s.convertToAPI(metric)
	}
	return result, nil
}

func (s *usageService) GetMetricsByAPIKeyID(ctx context.Context, apiKeyID uuid.UUID, start, end time.Time) ([]*UsageMetric, error) {
	metrics, err := s.repo.GetByAPIKeyID(ctx, apiKeyID, start, end)
	if err != nil {
		return nil, err
	}

	result := make([]*UsageMetric, len(metrics))
	for i, metric := range metrics {
		result[i] = s.convertToAPI(metric)
	}
	return result, nil
}

func (s *usageService) GetTokenSummaryByAppID(ctx context.Context, appID uuid.UUID, start, end time.Time) (*TokenSummary, error) {
	summary, err := s.repo.SumTokensByAppID(ctx, appID, start, end)
	if err != nil {
		return nil, err
	}

	return &TokenSummary{
		TotalPromptTokens:     summary.TotalPromptTokens,
		TotalCompletionTokens: summary.TotalCompletionTokens,
		TotalTokens:           summary.TotalTokens,
		RequestCount:          summary.RequestCount,
	}, nil
}

func (s *usageService) GetTokenSummaryByOrgID(ctx context.Context, orgID uuid.UUID, start, end time.Time) (*TokenSummary, error) {
	summary, err := s.repo.SumTokensByOrgID(ctx, orgID, start, end)
	if err != nil {
		return nil, err
	}

	return &TokenSummary{
		TotalPromptTokens:     summary.TotalPromptTokens,
		TotalCompletionTokens: summary.TotalCompletionTokens,
		TotalTokens:           summary.TotalTokens,
		RequestCount:          summary.RequestCount,
	}, nil
}

func (s *usageService) GetUsageByModel(ctx context.Context, appID uuid.UUID, start, end time.Time, limit, offset int) ([]*ModelUsageSummary, error) {
	summaries, err := s.repo.GetUsageByModel(ctx, appID, start, end, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*ModelUsageSummary, len(summaries))
	for i, summary := range summaries {
		result[i] = &ModelUsageSummary{
			ModelName:             summary.ModelName,
			Provider:              summary.Provider,
			TotalPromptTokens:     summary.TotalPromptTokens,
			TotalCompletionTokens: summary.TotalCompletionTokens,
			TotalTokens:           summary.TotalTokens,
			RequestCount:          summary.RequestCount,
		}
	}
	return result, nil
}

func (s *usageService) convertToAPI(metric *model.UsageMetric) *UsageMetric {
	var modelID *string
	if metric.ModelID != nil {
		id := metric.ModelID.String()
		modelID = &id
	}

	return &UsageMetric{
		ID:                metric.ID.String(),
		OrgID:             metric.OrgID.String(),
		AppID:             metric.AppID.String(),
		APIKeyID:          metric.APIKeyID.String(),
		ModelID:           modelID,
		Provider:          metric.Provider,
		ModelName:         metric.ModelName,
		PromptTokens:      metric.PromptTokens,
		CompletionTokens:  metric.CompletionTokens,
		TotalTokens:       metric.TotalTokens,
		RequestSizeBytes:  metric.RequestSizeBytes,
		ResponseSizeBytes: metric.ResponseSizeBytes,
		Timestamp:         metric.Timestamp,
	}
}
