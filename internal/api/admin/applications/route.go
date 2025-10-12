package applications

import (
	"context"
	"net/http"

	"github.com/WebDeveloperBen/ai-gateway/internal/exceptions"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type ApplicationService struct {
	Applications ApplicationsService
}

func NewRouter(applications ApplicationsService) *ApplicationService {
	return &ApplicationService{Applications: applications}
}

func (s *ApplicationService) RegisterRoutes(grp *huma.Group) {
	// POST /applications
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-create-application",
		Method:        http.MethodPost,
		Path:          "/applications",
		Summary:       "Create a new application",
		Description:   "Creates a new application within the organization with the specified configuration.",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Applications"},
	}, exceptions.Handle(func(ctx context.Context, in *CreateApplicationRequest) (*CreateApplicationResponse, error) {
		// Get org ID from context (set by middleware)
		orgID, ok := ctx.Value("org_id").(uuid.UUID)
		if !ok {
			return nil, huma.Error401Unauthorized("organization not found in context")
		}

		app, err := s.Applications.CreateApplication(ctx, orgID, in.Body)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return &CreateApplicationResponse{Body: app}, nil
	}))

	// GET /applications
	huma.Register(grp, huma.Operation{
		OperationID: "admin-list-applications",
		Method:      http.MethodGet,
		Path:        "/applications",
		Summary:     "List all applications",
		Description: "Retrieves a list of all applications belonging to the organization.",
		Tags:        []string{"Applications"},
	}, exceptions.Handle(func(ctx context.Context, in *model.ListRequest) (*ListApplicationsResponse, error) {
		// Get org ID from context (set by middleware)
		orgID, ok := ctx.Value("org_id").(uuid.UUID)
		if !ok {
			return nil, huma.Error401Unauthorized("organization not found in context")
		}

		normalized := model.NormalizePagination(*in)
		apps, err := s.Applications.ListApplications(ctx, orgID, normalized.Limit, normalized.Offset)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list applications")
		}

		return &ListApplicationsResponse{Body: apps}, nil
	}))

	// GET /applications/{id}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-get-application",
		Method:      http.MethodGet,
		Path:        "/applications/{id}",
		Summary:     "Get application details",
		Description: "Retrieves detailed information about a specific application by its ID.",
		Tags:        []string{"Applications"},
	}, exceptions.Handle(func(ctx context.Context, in *struct {
		ID string `path:"id" required:"true"`
	},
	) (*GetApplicationResponse, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid application ID")
		}

		app, err := s.Applications.GetApplication(ctx, id)
		if err != nil {
			return nil, huma.Error404NotFound("application not found")
		}

		return &GetApplicationResponse{Body: app}, nil
	}))

	// PUT /applications/{id}
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-update-application",
		Method:        http.MethodPut,
		Path:          "/applications/{id}",
		Summary:       "Update application",
		Description:   "Updates an existing application's configuration and settings.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Applications"},
	}, exceptions.Handle(func(ctx context.Context, in *UpdateApplicationRequest) (*UpdateApplicationResponse, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid application ID")
		}

		app, err := s.Applications.UpdateApplication(ctx, id, in.Body)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return &UpdateApplicationResponse{Body: app}, nil
	}))

	// DELETE /applications/{id}
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-delete-application",
		Method:        http.MethodDelete,
		Path:          "/applications/{id}",
		Summary:       "Delete application",
		Description:   "Permanently deletes an application and all its associated resources.",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Applications"},
	}, exceptions.Handle(func(ctx context.Context, in *struct {
		ID string `path:"id" required:"true"`
	},
	) (*struct{}, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid application ID")
		}

		if err := s.Applications.DeleteApplication(ctx, id); err != nil {
			return nil, huma.Error404NotFound("application not found")
		}

		return &struct{}{}, nil
	}))
}
