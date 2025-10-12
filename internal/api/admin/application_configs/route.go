package application_configs

import (
	"context"
	"net/http"

	"github.com/WebDeveloperBen/ai-gateway/internal/exceptions"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type ApplicationConfigService struct {
	ApplicationConfigs ApplicationConfigsService
}

func NewRouter(configs ApplicationConfigsService) *ApplicationConfigService {
	return &ApplicationConfigService{ApplicationConfigs: configs}
}

func (s *ApplicationConfigService) RegisterRoutes(grp *huma.Group) {
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-create-application-config",
		Method:        http.MethodPost,
		Path:          "/application-configs",
		Summary:       "Create a new application config",
		Description:   "Creates a new environment-specific configuration for an application.",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Application Configs"},
	}, exceptions.Handle(func(ctx context.Context, in *CreateApplicationConfigRequest) (*CreateApplicationConfigResponse, error) {
		orgID, ok := ctx.Value("org_id").(uuid.UUID)
		if !ok {
			return nil, huma.Error401Unauthorized("organization not found in context")
		}

		cfg, err := s.ApplicationConfigs.CreateApplicationConfig(ctx, orgID, in.Body)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return &CreateApplicationConfigResponse{Body: cfg}, nil
	}))

	huma.Register(grp, huma.Operation{
		OperationID: "admin-list-application-configs",
		Method:      http.MethodGet,
		Path:        "/applications/{app_id}/configs",
		Summary:     "List application configs",
		Description: "Retrieves all environment configurations for a specific application.",
		Tags:        []string{"Application Configs"},
	}, exceptions.Handle(func(ctx context.Context, in *struct {
		AppID string `path:"app_id" required:"true"`
	},
	) (*ListApplicationConfigsResponse, error) {
		appID, err := uuid.Parse(in.AppID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid application ID")
		}

		configs, err := s.ApplicationConfigs.ListApplicationConfigs(ctx, appID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list application configs")
		}

		return &ListApplicationConfigsResponse{Body: configs}, nil
	}))

	huma.Register(grp, huma.Operation{
		OperationID: "admin-get-application-config",
		Method:      http.MethodGet,
		Path:        "/application-configs/{id}",
		Summary:     "Get application config details",
		Description: "Retrieves detailed information about a specific application configuration by its ID.",
		Tags:        []string{"Application Configs"},
	}, exceptions.Handle(func(ctx context.Context, in *struct {
		ID string `path:"id" required:"true"`
	},
	) (*GetApplicationConfigResponse, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid config ID")
		}

		cfg, err := s.ApplicationConfigs.GetApplicationConfig(ctx, id)
		if err != nil {
			return nil, huma.Error404NotFound("application config not found")
		}

		return &GetApplicationConfigResponse{Body: cfg}, nil
	}))

	huma.Register(grp, huma.Operation{
		OperationID: "admin-get-application-config-by-env",
		Method:      http.MethodGet,
		Path:        "/applications/{app_id}/configs/{environment}",
		Summary:     "Get application config by environment",
		Description: "Retrieves configuration for a specific application and environment.",
		Tags:        []string{"Application Configs"},
	}, exceptions.Handle(func(ctx context.Context, in *struct {
		AppID       string `path:"app_id" required:"true"`
		Environment string `path:"environment" required:"true"`
	},
	) (*GetApplicationConfigResponse, error) {
		appID, err := uuid.Parse(in.AppID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid application ID")
		}

		cfg, err := s.ApplicationConfigs.GetApplicationConfigByEnv(ctx, appID, in.Environment)
		if err != nil {
			return nil, huma.Error404NotFound("application config not found")
		}

		return &GetApplicationConfigResponse{Body: cfg}, nil
	}))

	huma.Register(grp, huma.Operation{
		OperationID:   "admin-update-application-config",
		Method:        http.MethodPut,
		Path:          "/application-configs/{id}",
		Summary:       "Update application config",
		Description:   "Updates an existing application configuration.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Application Configs"},
	}, exceptions.Handle(func(ctx context.Context, in *UpdateApplicationConfigRequest) (*UpdateApplicationConfigResponse, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid config ID")
		}

		cfg, err := s.ApplicationConfigs.UpdateApplicationConfig(ctx, id, in.Body)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return &UpdateApplicationConfigResponse{Body: cfg}, nil
	}))

	huma.Register(grp, huma.Operation{
		OperationID:   "admin-delete-application-config",
		Method:        http.MethodDelete,
		Path:          "/application-configs/{id}",
		Summary:       "Delete application config",
		Description:   "Permanently deletes an application configuration.",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Application Configs"},
	}, exceptions.Handle(func(ctx context.Context, in *struct {
		ID string `path:"id" required:"true"`
	},
	) (*struct{}, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid config ID")
		}

		if err := s.ApplicationConfigs.DeleteApplicationConfig(ctx, id); err != nil {
			return nil, huma.Error404NotFound("application config not found")
		}

		return &struct{}{}, nil
	}))
}
