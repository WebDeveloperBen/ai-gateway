package applications

import (
	"context"
	"net/http"

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
	// POST /admin/applications
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-create-application",
		Method:        http.MethodPost,
		Path:          "/admin/applications",
		Summary:       "Create application",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		Body CreateApplicationRequest `json:"body"`
	}) (*struct {
		Application *Application `json:"application"`
	}, error) {
		// Get org ID from context (set by middleware)
		orgID, ok := ctx.Value("org_id").(uuid.UUID)
		if !ok {
			return nil, huma.Error401Unauthorized("organization not found in context")
		}

		app, err := s.Applications.CreateApplication(ctx, orgID, in.Body)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return &struct {
			Application *Application `json:"application"`
		}{Application: app}, nil
	})

	// GET /admin/applications
	huma.Register(grp, huma.Operation{
		OperationID: "admin-list-applications",
		Method:      http.MethodGet,
		Path:        "/admin/applications",
		Summary:     "List applications",
		Tags:        []string{"Admin"},
	}, func(ctx context.Context, in *struct{}) (*struct {
		Applications []*Application `json:"applications"`
	}, error) {
		// Get org ID from context (set by middleware)
		orgID, ok := ctx.Value("org_id").(uuid.UUID)
		if !ok {
			return nil, huma.Error401Unauthorized("organization not found in context")
		}

		apps, err := s.Applications.ListApplications(ctx, orgID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list applications")
		}

		return &struct {
			Applications []*Application `json:"applications"`
		}{Applications: apps}, nil
	})

	// GET /admin/applications/{id}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-get-application",
		Method:      http.MethodGet,
		Path:        "/admin/applications/{id}",
		Summary:     "Get application by ID",
		Tags:        []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		ID string `path:"id" required:"true"`
	}) (*struct {
		Application *Application `json:"application"`
	}, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid application ID")
		}

		app, err := s.Applications.GetApplication(ctx, id)
		if err != nil {
			return nil, huma.Error404NotFound("application not found")
		}

		return &struct {
			Application *Application `json:"application"`
		}{Application: app}, nil
	})

	// PUT /admin/applications/{id}
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-update-application",
		Method:        http.MethodPut,
		Path:          "/admin/applications/{id}",
		Summary:       "Update application",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		ID   string                   `path:"id" required:"true"`
		Body UpdateApplicationRequest `json:"body"`
	}) (*struct {
		Application *Application `json:"application"`
	}, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid application ID")
		}

		app, err := s.Applications.UpdateApplication(ctx, id, in.Body)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return &struct {
			Application *Application `json:"application"`
		}{Application: app}, nil
	})

	// DELETE /admin/applications/{id}
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-delete-application",
		Method:        http.MethodDelete,
		Path:          "/admin/applications/{id}",
		Summary:       "Delete application",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		ID string `path:"id" required:"true"`
	}) (*struct{}, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid application ID")
		}

		if err := s.Applications.DeleteApplication(ctx, id); err != nil {
			return nil, huma.Error404NotFound("application not found")
		}

		return &struct{}{}, nil
	})
}
