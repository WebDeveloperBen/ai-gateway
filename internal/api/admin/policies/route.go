package policies

import (
	"context"
	"net/http"

	"github.com/WebDeveloperBen/ai-gateway/internal/exceptions"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
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
	// POST /policies
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-create-policy",
		Method:        http.MethodPost,
		Path:          "/policies",
		Summary:       "Create a new policy",
		Description:   "Creates a new policy to control API access and behavior for applications.",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Policies"},
	}, exceptions.Handle(func(ctx context.Context, in *CreatePolicyRequest) (*CreatePolicyResponse, error) {
		orgID, err := uuid.Parse(in.Body.OrgID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid organization ID")
		}

		policy, err := s.Policies.CreatePolicy(ctx, orgID, in.Body)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return &CreatePolicyResponse{Body: policy}, nil
	}))

	// GET /policies?app_id={app_id}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-list-policies",
		Method:      http.MethodGet,
		Path:        "/policies",
		Summary:     "List policies for an application",
		Description: "Retrieves all policies associated with a specific application.",
		Tags:        []string{"Policies"},
	}, exceptions.Handle(func(ctx context.Context, in *ListPoliciesRequest) (*ListPoliciesResponse, error) {
		appID, err := uuid.Parse(in.AppID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid app_id")
		}

		normalized := model.NormalizePagination(model.ListRequest{Limit: in.Limit, Offset: in.Offset})
		policies, err := s.Policies.ListPolicies(ctx, appID, normalized.Limit, normalized.Offset)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list policies")
		}

		return &ListPoliciesResponse{Body: policies}, nil
	}))

	// GET /policies/enabled?app_id={app_id}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-list-enabled-policies",
		Method:      http.MethodGet,
		Path:        "/policies/enabled",
		Summary:     "List enabled policies for an application",
		Description: "Retrieves only the policies that are currently enabled for a specific application.",
		Tags:        []string{"Policies"},
	}, exceptions.Handle(func(ctx context.Context, in *ListEnabledPoliciesRequest) (*ListEnabledPoliciesResponse, error) {
		appID, err := uuid.Parse(in.AppID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid app_id")
		}

		normalized := model.NormalizePagination(model.ListRequest{Limit: in.Limit, Offset: in.Offset})
		policies, err := s.Policies.ListEnabledPolicies(ctx, appID, normalized.Limit, normalized.Offset)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list enabled policies")
		}

		return &ListEnabledPoliciesResponse{Body: policies}, nil
	}))

	// GET /policies/{id}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-get-policy",
		Method:      http.MethodGet,
		Path:        "/policies/{id}",
		Summary:     "Get policy details",
		Description: "Retrieves detailed information about a specific policy by its ID.",
		Tags:        []string{"Policies"},
	}, exceptions.Handle(func(ctx context.Context, in *struct {
		ID string `path:"id" required:"true"`
	}) (*GetPolicyResponse, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid policy ID")
		}

		policy, err := s.Policies.GetPolicy(ctx, id)
		if err != nil {
			return nil, err
		}

		return &GetPolicyResponse{Body: policy}, nil
	}))

	// PUT /policies/{id}
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-update-policy",
		Method:        http.MethodPut,
		Path:          "/policies/{id}",
		Summary:       "Update policy",
		Description:   "Updates an existing policy's configuration and rules.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Policies"},
	}, exceptions.Handle(func(ctx context.Context, in *UpdatePolicyRequest) (*UpdatePolicyResponse, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid policy ID")
		}

		policy, err := s.Policies.UpdatePolicy(ctx, id, in.Body)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return &UpdatePolicyResponse{Body: policy}, nil
	}))

	// DELETE /policies/{id}
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-delete-policy",
		Method:        http.MethodDelete,
		Path:          "/policies/{id}",
		Summary:       "Delete policy",
		Description:   "Permanently deletes a policy and removes it from all associated applications.",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Policies"},
	}, exceptions.Handle(func(ctx context.Context, in *struct {
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
	}))

	// POST /policies/{id}/enable
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-enable-policy",
		Method:        http.MethodPost,
		Path:          "/policies/{id}/enable",
		Summary:       "Enable policy",
		Description:   "Enables a policy, making it active for applications that have it assigned.",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Policies"},
	}, exceptions.Handle(func(ctx context.Context, in *struct {
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
	}))

	// POST /policies/{id}/disable
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-disable-policy",
		Method:        http.MethodPost,
		Path:          "/policies/{id}/disable",
		Summary:       "Disable policy",
		Description:   "Disables a policy, making it inactive for all applications.",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Policies"},
	}, exceptions.Handle(func(ctx context.Context, in *struct {
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
	}))

	// POST /policies/{id}/attach/{app_id}
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-attach-policy-to-app",
		Method:        http.MethodPost,
		Path:          "/policies/{id}/attach/{app_id}",
		Summary:       "Attach policy to application",
		Description:   "Attaches an existing policy to a specific application.",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Policies"},
	}, exceptions.Handle(func(ctx context.Context, in *struct {
		PolicyID string `path:"id" required:"true"`
		AppID    string `path:"app_id" required:"true"`
	}) (*struct{}, error) {
		policyID, err := uuid.Parse(in.PolicyID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid policy ID")
		}

		appID, err := uuid.Parse(in.AppID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid app ID")
		}

		if err := s.Policies.AttachPolicyToApp(ctx, policyID, appID); err != nil {
			return nil, huma.Error400BadRequest("failed to attach policy to app")
		}

		return &struct{}{}, nil
	}))

	// POST /policies/{id}/detach/{app_id}
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-detach-policy-from-app",
		Method:        http.MethodPost,
		Path:          "/policies/{id}/detach/{app_id}",
		Summary:       "Detach policy from application",
		Description:   "Detaches a policy from a specific application.",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Policies"},
	}, exceptions.Handle(func(ctx context.Context, in *struct {
		PolicyID string `path:"id" required:"true"`
		AppID    string `path:"app_id" required:"true"`
	}) (*struct{}, error) {
		policyID, err := uuid.Parse(in.PolicyID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid policy ID")
		}

		appID, err := uuid.Parse(in.AppID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid app ID")
		}

		if err := s.Policies.DetachPolicyFromApp(ctx, policyID, appID); err != nil {
			return nil, huma.Error400BadRequest("failed to detach policy from app")
		}

		return &struct{}{}, nil
	}))
}
