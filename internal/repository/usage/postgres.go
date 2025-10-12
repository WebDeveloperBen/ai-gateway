package usage

import (
	"context"
	"database/sql"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type postgresRepo struct {
	q *db.Queries
}

func NewPostgresRepo(q *db.Queries) Repository {
	return &postgresRepo{q: q}
}

func (r *postgresRepo) GetByAppID(ctx context.Context, appID uuid.UUID, start, end time.Time, limit, offset int) ([]*model.UsageMetric, error) {
	metrics, err := r.q.GetUsageMetricsByApp(ctx, db.GetUsageMetricsByAppParams{
		AppID:       appID,
		Timestamp:   pgtype.Timestamptz{Time: start, Valid: true},
		Timestamp_2: pgtype.Timestamptz{Time: end, Valid: true},
		Limit:       int32(limit),
		Offset:      int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*model.UsageMetric, len(metrics))
	for i, metric := range metrics {
		result[i] = r.convertToModel(metric)
	}
	return result, nil
}

func (r *postgresRepo) GetByOrgID(ctx context.Context, orgID uuid.UUID, start, end time.Time) ([]*model.UsageMetric, error) {
	metrics, err := r.q.GetUsageMetricsByOrg(ctx, db.GetUsageMetricsByOrgParams{
		OrgID:       orgID,
		Timestamp:   pgtype.Timestamptz{Time: start, Valid: true},
		Timestamp_2: pgtype.Timestamptz{Time: end, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	result := make([]*model.UsageMetric, len(metrics))
	for i, metric := range metrics {
		result[i] = r.convertToModel(metric)
	}
	return result, nil
}

func (r *postgresRepo) GetByAPIKeyID(ctx context.Context, apiKeyID uuid.UUID, start, end time.Time) ([]*model.UsageMetric, error) {
	metrics, err := r.q.GetUsageMetricsByAPIKey(ctx, db.GetUsageMetricsByAPIKeyParams{
		ApiKeyID:    apiKeyID,
		Timestamp:   pgtype.Timestamptz{Time: start, Valid: true},
		Timestamp_2: pgtype.Timestamptz{Time: end, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	result := make([]*model.UsageMetric, len(metrics))
	for i, metric := range metrics {
		result[i] = r.convertToModel(metric)
	}
	return result, nil
}

func (r *postgresRepo) SumTokensByAppID(ctx context.Context, appID uuid.UUID, start, end time.Time) (*TokenSummary, error) {
	sum, err := r.q.SumTokensByApp(ctx, db.SumTokensByAppParams{
		AppID:       appID,
		Timestamp:   pgtype.Timestamptz{Time: start, Valid: true},
		Timestamp_2: pgtype.Timestamptz{Time: end, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &TokenSummary{
		TotalPromptTokens:     r.convertInterfaceToInt(sum.TotalPromptTokens),
		TotalCompletionTokens: r.convertInterfaceToInt(sum.TotalCompletionTokens),
		TotalTokens:           r.convertInterfaceToInt(sum.TotalTokens),
		RequestCount:          int(sum.RequestCount),
	}, nil
}

func (r *postgresRepo) SumTokensByOrgID(ctx context.Context, orgID uuid.UUID, start, end time.Time) (*TokenSummary, error) {
	sum, err := r.q.SumTokensByOrg(ctx, db.SumTokensByOrgParams{
		OrgID:       orgID,
		Timestamp:   pgtype.Timestamptz{Time: start, Valid: true},
		Timestamp_2: pgtype.Timestamptz{Time: end, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &TokenSummary{
		TotalPromptTokens:     r.convertInterfaceToInt(sum.TotalPromptTokens),
		TotalCompletionTokens: r.convertInterfaceToInt(sum.TotalCompletionTokens),
		TotalTokens:           r.convertInterfaceToInt(sum.TotalTokens),
		RequestCount:          int(sum.RequestCount),
	}, nil
}

func (r *postgresRepo) GetUsageByModel(ctx context.Context, appID uuid.UUID, start, end time.Time, limit, offset int) ([]*ModelUsageSummary, error) {
	rows, err := r.q.GetUsageByModel(ctx, db.GetUsageByModelParams{
		AppID:       appID,
		Timestamp:   pgtype.Timestamptz{Time: start, Valid: true},
		Timestamp_2: pgtype.Timestamptz{Time: end, Valid: true},
		Limit:       int32(limit),
		Offset:      int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*ModelUsageSummary, len(rows))
	for i, row := range rows {
		result[i] = &ModelUsageSummary{
			ModelName:             row.ModelName,
			Provider:              row.Provider,
			TotalPromptTokens:     r.convertInterfaceToInt(row.TotalPromptTokens),
			TotalCompletionTokens: r.convertInterfaceToInt(row.TotalCompletionTokens),
			TotalTokens:           r.convertInterfaceToInt(row.TotalTokens),
			RequestCount:          int(row.RequestCount),
		}
	}
	return result, nil
}

func (r *postgresRepo) Create(ctx context.Context, metric *model.UsageMetric) error {
	if metric.OrgID == uuid.Nil {
		return sql.ErrNoRows // or a custom error, but using this for consistency
	}
	if metric.AppID == uuid.Nil {
		return sql.ErrNoRows
	}
	if metric.APIKeyID == uuid.Nil {
		return sql.ErrNoRows
	}

	_, err := r.q.CreateUsageMetric(ctx, db.CreateUsageMetricParams{
		OrgID:             metric.OrgID,
		AppID:             metric.AppID,
		ApiKeyID:          metric.APIKeyID,
		ModelID:           metric.ModelID,
		Provider:          metric.Provider,
		ModelName:         metric.ModelName,
		PromptTokens:      int32(metric.PromptTokens),
		CompletionTokens:  int32(metric.CompletionTokens),
		TotalTokens:       int32(metric.TotalTokens),
		RequestSizeBytes:  int32(metric.RequestSizeBytes),
		ResponseSizeBytes: int32(metric.ResponseSizeBytes),
		Timestamp:         pgtype.Timestamptz{Time: metric.Timestamp, Valid: true},
	})
	return err
}

func (r *postgresRepo) convertToModel(metric db.UsageMetric) *model.UsageMetric {
	return &model.UsageMetric{
		ID:                metric.ID,
		OrgID:             metric.OrgID,
		AppID:             metric.AppID,
		APIKeyID:          metric.ApiKeyID,
		ModelID:           metric.ModelID,
		Provider:          metric.Provider,
		ModelName:         metric.ModelName,
		PromptTokens:      int(metric.PromptTokens),
		CompletionTokens:  int(metric.CompletionTokens),
		TotalTokens:       int(metric.TotalTokens),
		RequestSizeBytes:  int(metric.RequestSizeBytes),
		ResponseSizeBytes: int(metric.ResponseSizeBytes),
		Timestamp:         metric.Timestamp.Time,
	}
}

func (r *postgresRepo) convertInterfaceToInt(val interface{}) int {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int64:
		return int(v)
	case int32:
		return int(v)
	case int:
		return v
	case float64:
		return int(v)
	default:
		return 0
	}
}
