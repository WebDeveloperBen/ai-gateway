package usage

import (
	"context"
	"net/http"
	"time"

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
	// GET /admin/usage/metrics?app_id={app_id}&start={start}&end={end}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-get-usage-metrics-by-app",
		Method:      http.MethodGet,
		Path:        "/admin/usage/metrics",
		Summary:     "Get usage metrics for an application",
		Tags:        []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		AppID string `query:"app_id" required:"true"`
		Start string `query:"start" required:"true"` // ISO 8601 format
		End   string `query:"end" required:"true"`   // ISO 8601 format
	}) (*struct {
		Metrics []*UsageMetric `json:"metrics"`
	}, error,
	) {
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

		metrics, err := s.Usage.GetMetricsByAppID(ctx, appID, start, end)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to get usage metrics")
		}

		return &struct {
			Metrics []*UsageMetric `json:"metrics"`
		}{Metrics: metrics}, nil
	})

	// GET /admin/usage/summary/app/{app_id}?start={start}&end={end}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-get-usage-summary-by-app",
		Method:      http.MethodGet,
		Path:        "/admin/usage/summary/app/{app_id}",
		Summary:     "Get token usage summary for an application",
		Tags:        []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		AppID string `path:"app_id" required:"true"`
		Start string `query:"start" required:"true"`
		End   string `query:"end" required:"true"`
	}) (*struct {
		Summary *TokenSummary `json:"summary"`
	}, error,
	) {
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

		return &struct {
			Summary *TokenSummary `json:"summary"`
		}{Summary: summary}, nil
	})

	// GET /admin/usage/summary/org?start={start}&end={end}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-get-usage-summary-by-org",
		Method:      http.MethodGet,
		Path:        "/admin/usage/summary/org",
		Summary:     "Get token usage summary for the organization",
		Tags:        []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		Start string `query:"start" required:"true"`
		End   string `query:"end" required:"true"`
	}) (*struct {
		Summary *TokenSummary `json:"summary"`
	}, error,
	) {
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

		return &struct {
			Summary *TokenSummary `json:"summary"`
		}{Summary: summary}, nil
	})

	// GET /admin/usage/by-model/{app_id}?start={start}&end={end}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-get-usage-by-model",
		Method:      http.MethodGet,
		Path:        "/admin/usage/by-model/{app_id}",
		Summary:     "Get usage metrics grouped by model for an application",
		Tags:        []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		AppID string `path:"app_id" required:"true"`
		Start string `query:"start" required:"true"`
		End   string `query:"end" required:"true"`
	}) (*struct {
		Summaries []*ModelUsageSummary `json:"summaries"`
	}, error,
	) {
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

		summaries, err := s.Usage.GetUsageByModel(ctx, appID, start, end)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to get usage by model")
		}

		return &struct {
			Summaries []*ModelUsageSummary `json:"summaries"`
		}{Summaries: summaries}, nil
	})
}
