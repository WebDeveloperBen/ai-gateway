package catalog

import (
	"context"
	"net/http"

	"github.com/WebDeveloperBen/ai-gateway/internal/exceptions"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type CatalogRouter struct {
	Catalog CatalogService
}

func NewRouter(catalog CatalogService) *CatalogRouter {
	return &CatalogRouter{Catalog: catalog}
}

func (r *CatalogRouter) RegisterRoutes(grp *huma.Group) {
	// POST /catalog
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-create-model",
		Method:        http.MethodPost,
		Path:          "/catalog",
		Summary:       "Create a new model",
		Description:   "Creates a new model configuration in the catalog for the organization.",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Catalog"},
	}, exceptions.Handle(func(ctx context.Context, in *CreateModelRequest) (*CreateModelResponse, error) {
		// Get org ID from context (set by middleware)
		orgID, ok := ctx.Value("org_id").(uuid.UUID)
		if !ok {
			return nil, huma.Error401Unauthorized("organization not found in context")
		}

		model, err := r.Catalog.CreateModel(ctx, orgID, in.Body)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return &CreateModelResponse{Body: model}, nil
	}))

	// GET /catalog
	huma.Register(grp, huma.Operation{
		OperationID: "admin-list-models",
		Method:      http.MethodGet,
		Path:        "/catalog",
		Summary:     "List all models",
		Description: "Retrieves a list of all models configured in the catalog for the organization.",
		Tags:        []string{"Catalog"},
	}, exceptions.Handle(func(ctx context.Context, in *ListModelsRequest) (*ListModelsResponse, error) {
		// Get org ID from context (set by middleware)
		orgID, ok := ctx.Value("org_id").(uuid.UUID)
		if !ok {
			return nil, huma.Error401Unauthorized("organization not found in context")
		}

		normalized := model.NormalizePagination(model.ListRequest{Limit: in.Limit, Offset: in.Offset})
		models, err := r.Catalog.ListModels(ctx, orgID, normalized.Limit, normalized.Offset)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list models")
		}

		return &ListModelsResponse{Body: models}, nil
	}))

	// GET /catalog/enabled
	huma.Register(grp, huma.Operation{
		OperationID: "admin-list-enabled-models",
		Method:      http.MethodGet,
		Path:        "/catalog/enabled",
		Summary:     "List enabled models",
		Description: "Retrieves a list of all enabled models configured in the catalog for the organization.",
		Tags:        []string{"Catalog"},
	}, exceptions.Handle(func(ctx context.Context, in *ListEnabledModelsRequest) (*ListEnabledModelsResponse, error) {
		// Get org ID from context (set by middleware)
		orgID, ok := ctx.Value("org_id").(uuid.UUID)
		if !ok {
			return nil, huma.Error401Unauthorized("organization not found in context")
		}

		normalized := model.NormalizePagination(model.ListRequest{Limit: in.Limit, Offset: in.Offset})
		models, err := r.Catalog.ListEnabledModels(ctx, orgID, normalized.Limit, normalized.Offset)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list enabled models")
		}

		return &ListEnabledModelsResponse{Body: models}, nil
	}))

	// GET /catalog/{id}
	huma.Register(grp, huma.Operation{
		OperationID: "admin-get-model",
		Method:      http.MethodGet,
		Path:        "/catalog/{id}",
		Summary:     "Get model details",
		Description: "Retrieves detailed information about a specific model by its ID.",
		Tags:        []string{"Catalog"},
	}, exceptions.Handle(func(ctx context.Context, in *struct {
		ID string `path:"id" required:"true"`
	},
	) (*GetModelResponse, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid model ID")
		}

		model, err := r.Catalog.GetModel(ctx, id)
		if err != nil {
			return nil, huma.Error404NotFound("model not found")
		}

		return &GetModelResponse{Body: model}, nil
	}))

	// PUT /catalog/{id}
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-update-model",
		Method:        http.MethodPut,
		Path:          "/catalog/{id}",
		Summary:       "Update model",
		Description:   "Updates an existing model's configuration and settings.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Catalog"},
	}, exceptions.Handle(func(ctx context.Context, in *UpdateModelRequest) (*UpdateModelResponse, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid model ID")
		}

		model, err := r.Catalog.UpdateModel(ctx, id, in.Body)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return &UpdateModelResponse{Body: model}, nil
	}))

	// DELETE /catalog/{id}
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-delete-model",
		Method:        http.MethodDelete,
		Path:          "/catalog/{id}",
		Summary:       "Delete model",
		Description:   "Permanently deletes a model from the catalog.",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Catalog"},
	}, exceptions.Handle(func(ctx context.Context, in *struct {
		ID string `path:"id" required:"true"`
	},
	) (*struct{}, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid model ID")
		}

		if err := r.Catalog.DeleteModel(ctx, id); err != nil {
			return nil, huma.Error404NotFound("model not found")
		}

		return &struct{}{}, nil
	}))

	// POST /catalog/{id}/enable
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-enable-model",
		Method:        http.MethodPost,
		Path:          "/catalog/{id}/enable",
		Summary:       "Enable model",
		Description:   "Enables a model, making it available for use in applications.",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Catalog"},
	}, exceptions.Handle(func(ctx context.Context, in *struct {
		ID string `path:"id" required:"true"`
	},
	) (*struct{}, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid model ID")
		}

		if err := r.Catalog.EnableModel(ctx, id); err != nil {
			return nil, huma.Error404NotFound("model not found")
		}

		return &struct{}{}, nil
	}))

	// POST /catalog/{id}/disable
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-disable-model",
		Method:        http.MethodPost,
		Path:          "/catalog/{id}/disable",
		Summary:       "Disable model",
		Description:   "Disables a model, making it unavailable for use in applications.",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Catalog"},
	}, exceptions.Handle(func(ctx context.Context, in *struct {
		ID string `path:"id" required:"true"`
	},
	) (*struct{}, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid model ID")
		}

		if err := r.Catalog.DisableModel(ctx, id); err != nil {
			return nil, huma.Error404NotFound("model not found")
		}

		return &struct{}{}, nil
	}))
}
