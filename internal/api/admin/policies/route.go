package policies

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type PolicyService struct {
	Policies PoliciesService
}

func NewRouter(policies PoliciesService) *PolicyService {
	return &PolicyService{Policies: policies}
}

func (s *PolicyService) RegisterRoutes(grp *huma.Group) {
	// POST /admin/policies
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-create-policy",
		Method:        http.MethodPost,
		Path:          "/admin/policies",
		Summary:       "Create policy",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		Body CreatePolicyRequest `json:"body"`
	}) (*struct {
		Policy *Policy `json:"policy"`
	}, error) {
		// Get org ID from context (set by middleware)
		orgID, ok := ctx.Value("org_id").(uuid.UUID)
		if !ok {
			return nil, huma.Error401Unauthorized("organization not found in context")
		}

		policy, err := s.Policies.CreatePolicy(ctx, orgID, in.Body)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return &struct {
			Policy *Policy `json:"policy"`
		}{Policy: policy}, nil
	})

	// GET /admin/policies?app_id={app_id}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-list-policies",
		Method:      http.MethodGet,
		Path:        "/admin/policies",
		Summary:     "List policies for an application",
		Tags:        []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		AppID string `query:"app_id" required:"true"`
	}) (*struct {
		Policies []*Policy `json:"policies"`
	}, error) {
		appID, err := uuid.Parse(in.AppID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid app_id")
		}

		policies, err := s.Policies.ListPolicies(ctx, appID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list policies")
		}

		return &struct {
			Policies []*Policy `json:"policies"`
		}{Policies: policies}, nil
	})

	// GET /admin/policies/enabled?app_id={app_id}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-list-enabled-policies",
		Method:      http.MethodGet,
		Path:        "/admin/policies/enabled",
		Summary:     "List enabled policies for an application",
		Tags:        []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		AppID string `query:"app_id" required:"true"`
	}) (*struct {
		Policies []*Policy `json:"policies"`
	}, error) {
		appID, err := uuid.Parse(in.AppID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid app_id")
		}

		policies, err := s.Policies.ListEnabledPolicies(ctx, appID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list enabled policies")
		}

		return &struct {
			Policies []*Policy `json:"policies"`
		}{Policies: policies}, nil
	})

	// GET /admin/policies/{id}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-get-policy",
		Method:      http.MethodGet,
		Path:        "/admin/policies/{id}",
		Summary:     "Get policy by ID",
		Tags:        []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		ID string `path:"id" required:"true"`
	}) (*struct {
		Policy *Policy `json:"policy"`
	}, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid policy ID")
		}

		policy, err := s.Policies.GetPolicy(ctx, id)
		if err != nil {
			return nil, huma.Error404NotFound("policy not found")
		}

		return &struct {
			Policy *Policy `json:"policy"`
		}{Policy: policy}, nil
	})

	// PUT /admin/policies/{id}
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-update-policy",
		Method:        http.MethodPut,
		Path:          "/admin/policies/{id}",
		Summary:       "Update policy",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		ID   string              `path:"id" required:"true"`
		Body UpdatePolicyRequest `json:"body"`
	}) (*struct {
		Policy *Policy `json:"policy"`
	}, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid policy ID")
		}

		policy, err := s.Policies.UpdatePolicy(ctx, id, in.Body)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return &struct {
			Policy *Policy `json:"policy"`
		}{Policy: policy}, nil
	})

	// DELETE /admin/policies/{id}
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-delete-policy",
		Method:        http.MethodDelete,
		Path:          "/admin/policies/{id}",
		Summary:       "Delete policy",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		ID string `path:"id" required:"true"`
	}) (*struct{}, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid policy ID")
		}

		if err := s.Policies.DeletePolicy(ctx, id); err != nil {
			return nil, huma.Error404NotFound("policy not found")
		}

		return &struct{}{}, nil
	})

	// POST /admin/policies/{id}/enable
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-enable-policy",
		Method:        http.MethodPost,
		Path:          "/admin/policies/{id}/enable",
		Summary:       "Enable policy",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		ID string `path:"id" required:"true"`
	}) (*struct{}, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid policy ID")
		}

		if err := s.Policies.EnablePolicy(ctx, id); err != nil {
			return nil, huma.Error404NotFound("policy not found")
		}

		return &struct{}{}, nil
	})

	// POST /admin/policies/{id}/disable
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-disable-policy",
		Method:        http.MethodPost,
		Path:          "/admin/policies/{id}/disable",
		Summary:       "Disable policy",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		ID string `path:"id" required:"true"`
	}) (*struct{}, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid policy ID")
		}

		if err := s.Policies.DisablePolicy(ctx, id); err != nil {
			return nil, huma.Error404NotFound("policy not found")
		}

		return &struct{}{}, nil
	})
}
