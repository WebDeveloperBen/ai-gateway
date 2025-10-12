package catalog

import (
	"context"
	"errors"
	"maps"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/repository/catalog"
	"github.com/google/uuid"
)

type CatalogService interface {
	CreateModel(ctx context.Context, orgID uuid.UUID, req CreateModelBody) (*Model, error)
	GetModel(ctx context.Context, id uuid.UUID) (*Model, error)
	GetModelByProviderAndName(ctx context.Context, orgID uuid.UUID, provider, modelName string) (*Model, error)
	ListModels(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*Model, error)
	ListEnabledModels(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*Model, error)
	UpdateModel(ctx context.Context, id uuid.UUID, req UpdateModelBody) (*Model, error)
	DeleteModel(ctx context.Context, id uuid.UUID) error
	EnableModel(ctx context.Context, id uuid.UUID) error
	DisableModel(ctx context.Context, id uuid.UUID) error
}

type catalogService struct {
	repo catalog.Repository
}

func NewService(repo catalog.Repository) CatalogService {
	return &catalogService{repo: repo}
}

func (s *catalogService) CreateModel(ctx context.Context, orgID uuid.UUID, req CreateModelBody) (*Model, error) {
	if req.Provider == "" {
		return nil, errors.New("provider is required")
	}
	if req.ModelName == "" {
		return nil, errors.New("model_name is required")
	}
	if req.EndpointURL == "" {
		return nil, errors.New("endpoint_url is required")
	}
	if req.AuthType == "" {
		return nil, errors.New("auth_type is required")
	}
	// Validate auth_type is a valid enum value
	switch req.AuthType {
	case AuthTypeAPIKey, AuthTypeOAuth2, AuthTypeAzureAD:
		// Valid
	default:
		return nil, errors.New("auth_type must be one of: api_key, oauth2, azure_ad")
	}
	if req.AuthConfig.Type == "" {
		return nil, errors.New("auth_config.type is required")
	}

	// Convert AuthConfig for repository layer
	authConfigModel := convertAuthConfigFromAPI(req.AuthConfig)

	model, err := s.repo.Create(ctx, orgID, req.Provider, req.ModelName, req.DeploymentName, req.EndpointURL, APIAuthTypeToModel(req.AuthType), authConfigModel, req.Metadata, req.Enabled)
	if err != nil {
		return nil, err
	}

	return s.convertToAPI(model), nil
}

func (s *catalogService) GetModel(ctx context.Context, id uuid.UUID) (*Model, error) {
	model, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.convertToAPI(model), nil
}

func (s *catalogService) GetModelByProviderAndName(ctx context.Context, orgID uuid.UUID, provider, modelName string) (*Model, error) {
	model, err := s.repo.GetByProviderAndName(ctx, orgID, provider, modelName)
	if err != nil {
		return nil, err
	}

	return s.convertToAPI(model), nil
}

func (s *catalogService) ListModels(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*Model, error) {
	models, err := s.repo.ListByOrgID(ctx, orgID, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*Model, len(models))
	for i, model := range models {
		result[i] = s.convertToAPI(model)
	}
	return result, nil
}

func (s *catalogService) ListEnabledModels(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*Model, error) {
	models, err := s.repo.ListEnabledByOrgID(ctx, orgID, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*Model, len(models))
	for i, model := range models {
		result[i] = s.convertToAPI(model)
	}
	return result, nil
}

func (s *catalogService) UpdateModel(ctx context.Context, id uuid.UUID, req UpdateModelBody) (*Model, error) {
	if req.Provider == "" {
		return nil, errors.New("provider is required")
	}
	if req.ModelName == "" {
		return nil, errors.New("model_name is required")
	}
	if req.EndpointURL == "" {
		return nil, errors.New("endpoint_url is required")
	}
	if req.AuthType == "" {
		return nil, errors.New("auth_type is required")
	}
	if req.AuthConfig.Type == "" {
		return nil, errors.New("auth_config.type is required")
	}

	// Convert AuthConfig for repository layer
	authConfigModel := convertAuthConfigFromAPI(req.AuthConfig)

	model, err := s.repo.Update(ctx, id, req.Provider, req.ModelName, req.DeploymentName, req.EndpointURL, APIAuthTypeToModel(req.AuthType), authConfigModel, req.Metadata, req.Enabled)
	if err != nil {
		return nil, err
	}

	return s.convertToAPI(model), nil
}

func (s *catalogService) DeleteModel(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *catalogService) EnableModel(ctx context.Context, id uuid.UUID) error {
	return s.repo.Enable(ctx, id)
}

func (s *catalogService) DisableModel(ctx context.Context, id uuid.UUID) error {
	return s.repo.Disable(ctx, id)
}

// APIAuthTypeToModel converts API AuthType to model.AuthType
func APIAuthTypeToModel(apiType AuthType) model.AuthType {
	switch apiType {
	case AuthTypeAPIKey:
		return model.AuthTypeAPIKey
	case AuthTypeOAuth2:
		return model.AuthTypeOAuth2
	case AuthTypeAzureAD:
		return model.AuthTypeAzureAD
	default:
		return model.AuthType(apiType) // Fallback
	}
}

// ModelAuthTypeToAPI converts model.AuthType to API AuthType
func ModelAuthTypeToAPI(modelType model.AuthType) AuthType {
	switch modelType {
	case model.AuthTypeAPIKey:
		return AuthTypeAPIKey
	case model.AuthTypeOAuth2:
		return AuthTypeOAuth2
	case model.AuthTypeAzureAD:
		return AuthTypeAzureAD
	default:
		return AuthType(modelType) // Fallback
	}
}

// convertAuthConfigToAPI converts model.AuthConfig to API AuthConfig
func convertAuthConfigToAPI(authConfig model.AuthConfig) AuthConfig {
	apiConfig := AuthConfig{
		Type: ModelAuthTypeToAPI(authConfig.Type),
	}

	// Copy the structured fields
	if authConfig.APIKey != nil {
		apiConfig.APIKey = authConfig.APIKey
	}
	if authConfig.ClientID != nil {
		apiConfig.ClientID = authConfig.ClientID
	}
	if authConfig.ClientSecret != nil {
		apiConfig.ClientSecret = authConfig.ClientSecret
	}
	if authConfig.TokenURL != nil {
		apiConfig.TokenURL = authConfig.TokenURL
	}
	if authConfig.TenantID != nil {
		apiConfig.TenantID = authConfig.TenantID
	}
	if authConfig.Resource != nil {
		apiConfig.Resource = authConfig.Resource
	}

	// Copy additional fields
	if authConfig.Additional != nil {
		apiConfig.Additional = make(map[string]any)
		maps.Copy(apiConfig.Additional, authConfig.Additional)
	}

	return apiConfig
}

// convertAuthConfigFromAPI converts API AuthConfig to model.AuthConfig
func convertAuthConfigFromAPI(apiConfig AuthConfig) model.AuthConfig {
	modelConfig := model.AuthConfig{
		Type: model.AuthType(apiConfig.Type),
	}

	// Copy the structured fields
	if apiConfig.APIKey != nil {
		modelConfig.APIKey = apiConfig.APIKey
	}
	if apiConfig.ClientID != nil {
		modelConfig.ClientID = apiConfig.ClientID
	}
	if apiConfig.ClientSecret != nil {
		modelConfig.ClientSecret = apiConfig.ClientSecret
	}
	if apiConfig.TokenURL != nil {
		modelConfig.TokenURL = apiConfig.TokenURL
	}
	if apiConfig.TenantID != nil {
		modelConfig.TenantID = apiConfig.TenantID
	}
	if apiConfig.Resource != nil {
		modelConfig.Resource = apiConfig.Resource
	}

	// Copy additional fields
	if apiConfig.Additional != nil {
		modelConfig.Additional = make(map[string]any)
		maps.Copy(modelConfig.Additional, apiConfig.Additional)
	}

	return modelConfig
}

func (s *catalogService) convertToAPI(model *model.Model) *Model {
	return &Model{
		ID:             model.ID.String(),
		OrgID:          model.OrgID.String(),
		Provider:       model.Provider,
		ModelName:      model.ModelName,
		DeploymentName: model.DeploymentName,
		EndpointURL:    model.EndpointURL,
		AuthType:       ModelAuthTypeToAPI(model.AuthType),
		AuthConfig:     convertAuthConfigToAPI(model.AuthConfig),
		Metadata:       model.Metadata,
		Enabled:        model.Enabled,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}
}
