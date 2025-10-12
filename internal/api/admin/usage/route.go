package usage

import (
	"context"
	"net/http"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/exceptions"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type UsageServiceRouter struct {
	Usage UsageService
}

func NewRouter(usage UsageService) *UsageServiceRouter {
	return &UsageServiceRouter{Usage: usage}
}

func (s *UsageServiceRouter) RegisterRoutes(grp *huma.Group) {
	// GET /usage/metrics?app_id={app_id}&start={start}&end={end}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-get-usage-metrics-by-app",
		Method:      http.MethodGet,
		Path:        "/usage/metrics",
		Summary:     "Get usage metrics for an application",
		Description: "Retrieves usage metrics and statistics for a specific application within a date range.",
		Tags:        []string{"Usage"},
	}, exceptions.Handle(func(ctx context.Context, in *GetUsageMetricsRequest) (*GetUsageMetricsResponse, error) {
		appID, err := uuid.Parse(in.AppID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid app_id")
		}

		start, err := time.Parse(time.RFC3339, in.Start)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid start time (use RFC3339 format)")
		}

		end, err := time.Parse(time.RFC3339, in.End)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid end time (use RFC3339 format)")
		}

		normalized := model.NormalizePagination(model.ListRequest{Limit: in.Limit, Offset: in.Offset})
		metrics, err := s.Usage.GetMetricsByAppID(ctx, appID, start, end, normalized.Limit, normalized.Offset)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to get usage metrics")
		}

		return &GetUsageMetricsResponse{Body: metrics}, nil
	}))

	// GET /usage/summary/app/{app_id}?start={start}&end={end}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-get-usage-summary-by-app",
		Method:      http.MethodGet,
		Path:        "/usage/summary/app/{app_id}",
		Summary:     "Get token usage summary for an application",
		Description: "Retrieves a summary of token usage for a specific application within a date range.",
		Tags:        []string{"Usage"},
	}, exceptions.Handle(func(ctx context.Context, in *GetUsageSummaryByAppRequest) (*GetUsageSummaryByAppResponse, error) {
		appID, err := uuid.Parse(in.AppID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid app_id")
		}

		start, err := time.Parse(time.RFC3339, in.Start)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid start time (use RFC3339 format)")
		}

		end, err := time.Parse(time.RFC3339, in.End)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid end time (use RFC3339 format)")
		}

		summary, err := s.Usage.GetTokenSummaryByAppID(ctx, appID, start, end)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to get usage summary")
		}

		return &GetUsageSummaryByAppResponse{Body: summary}, nil
	}))

	// GET /usage/summary/org?start={start}&end={end}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-get-usage-summary-by-org",
		Method:      http.MethodGet,
		Path:        "/usage/summary/org",
		Summary:     "Get token usage summary for the organization",
		Description: "Retrieves a summary of token usage for the entire organization within a date range.",
		Tags:        []string{"Usage"},
	}, exceptions.Handle(func(ctx context.Context, in *GetUsageSummaryByOrgRequest) (*GetUsageSummaryByOrgResponse, error) {
		// Get org ID from context (set by middleware)
		orgID, ok := ctx.Value("org_id").(uuid.UUID)
		if !ok {
			return nil, huma.Error401Unauthorized("organization not found in context")
		}

		start, err := time.Parse(time.RFC3339, in.Start)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid start time (use RFC3339 format)")
		}

		end, err := time.Parse(time.RFC3339, in.End)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid end time (use RFC3339 format)")
		}

		summary, err := s.Usage.GetTokenSummaryByOrgID(ctx, orgID, start, end)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to get usage summary")
		}

		return &GetUsageSummaryByOrgResponse{Body: summary}, nil
	}))

	// GET /usage/by-model/{app_id}?start={start}&end={end}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-get-usage-by-model",
		Method:      http.MethodGet,
		Path:        "/usage/by-model/{app_id}",
		Summary:     "Get usage metrics grouped by model for an application",
		Description: "Retrieves usage metrics grouped by AI model for a specific application within a date range.",
		Tags:        []string{"Usage"},
	}, exceptions.Handle(func(ctx context.Context, in *GetUsageByModelRequest) (*GetUsageByModelResponse, error) {
		appID, err := uuid.Parse(in.AppID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid app_id")
		}

		start, err := time.Parse(time.RFC3339, in.Start)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid start time (use RFC3339 format)")
		}

		end, err := time.Parse(time.RFC3339, in.End)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid end time (use RFC3339 format)")
		}

		normalized := model.NormalizePagination(model.ListRequest{Limit: in.Limit, Offset: in.Offset})
		summaries, err := s.Usage.GetUsageByModel(ctx, appID, start, end, normalized.Limit, normalized.Offset)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to get usage by model")
		}

		return &GetUsageByModelResponse{Body: summaries}, nil
	}))
}
