package catalog

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
)

type postgresRepo struct {
	q *db.Queries
}

func NewPostgresRepo(q *db.Queries) Repository {
	return &postgresRepo{q: q}
}

func (r *postgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Model, error) {
	m, err := r.q.GetModel(ctx, id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	authConfig, err := r.unmarshalAuthConfig(m.AuthConfig)
	if err != nil {
		return nil, err
	}

	metadata, err := r.unmarshalJSON(m.Metadata)
	if err != nil {
		return nil, err
	}

	return &model.Model{
		ID:             m.ID,
		OrgID:          m.OrgID,
		Provider:       m.Provider,
		ModelName:      m.ModelName,
		DeploymentName: m.DeploymentName,
		EndpointURL:    m.EndpointUrl,
		AuthType:       r.stringToAuthType(m.AuthType),
		AuthConfig:     authConfig,
		Metadata:       metadata,
		Enabled:        m.Enabled,
		CreatedAt:      m.CreatedAt.Time,
		UpdatedAt:      m.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) GetByProviderAndName(ctx context.Context, orgID uuid.UUID, provider, modelName string) (*model.Model, error) {
	m, err := r.q.GetModelByProviderAndName(ctx, db.GetModelByProviderAndNameParams{
		OrgID:     orgID,
		Provider:  provider,
		ModelName: modelName,
	})
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	authConfig, err := r.unmarshalAuthConfig(m.AuthConfig)
	if err != nil {
		return nil, err
	}

	metadata, err := r.unmarshalJSON(m.Metadata)
	if err != nil {
		return nil, err
	}

	return &model.Model{
		ID:             m.ID,
		OrgID:          m.OrgID,
		Provider:       m.Provider,
		ModelName:      m.ModelName,
		DeploymentName: m.DeploymentName,
		EndpointURL:    m.EndpointUrl,
		AuthType:       r.stringToAuthType(m.AuthType),
		AuthConfig:     authConfig,
		Metadata:       metadata,
		Enabled:        m.Enabled,
		CreatedAt:      m.CreatedAt.Time,
		UpdatedAt:      m.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) ListByOrgID(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*model.Model, error) {
	models, err := r.q.ListModels(ctx, db.ListModelsParams{
		OrgID:  orgID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*model.Model, len(models))
	for i, m := range models {
		authConfig, err := r.unmarshalAuthConfig(m.AuthConfig)
		if err != nil {
			return nil, err
		}

		metadata, err := r.unmarshalJSON(m.Metadata)
		if err != nil {
			return nil, err
		}

		result[i] = &model.Model{
			ID:             m.ID,
			OrgID:          m.OrgID,
			Provider:       m.Provider,
			ModelName:      m.ModelName,
			DeploymentName: m.DeploymentName,
			EndpointURL:    m.EndpointUrl,
			AuthType:       r.stringToAuthType(m.AuthType),
			AuthConfig:     authConfig,
			Metadata:       metadata,
			Enabled:        m.Enabled,
			CreatedAt:      m.CreatedAt.Time,
			UpdatedAt:      m.UpdatedAt.Time,
		}
	}
	return result, nil
}

func (r *postgresRepo) ListEnabledByOrgID(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*model.Model, error) {
	models, err := r.q.ListEnabledModels(ctx, db.ListEnabledModelsParams{
		OrgID:  orgID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*model.Model, len(models))
	for i, m := range models {
		authConfig, err := r.unmarshalAuthConfig(m.AuthConfig)
		if err != nil {
			return nil, err
		}

		metadata, err := r.unmarshalJSON(m.Metadata)
		if err != nil {
			return nil, err
		}

		result[i] = &model.Model{
			ID:             m.ID,
			OrgID:          m.OrgID,
			Provider:       m.Provider,
			ModelName:      m.ModelName,
			DeploymentName: m.DeploymentName,
			EndpointURL:    m.EndpointUrl,
			AuthType:       r.stringToAuthType(m.AuthType),
			AuthConfig:     authConfig,
			Metadata:       metadata,
			Enabled:        m.Enabled,
			CreatedAt:      m.CreatedAt.Time,
			UpdatedAt:      m.UpdatedAt.Time,
		}
	}
	return result, nil
}

func (r *postgresRepo) Create(ctx context.Context, orgID uuid.UUID, provider, modelName string, deploymentName *string, endpointURL string, authType model.AuthType, authConfig model.AuthConfig, metadata map[string]any, enabled bool) (*model.Model, error) {
	if orgID == uuid.Nil {
		return nil, errors.New("orgID cannot be nil")
	}
	if provider == "" {
		return nil, errors.New("provider cannot be empty")
	}
	if modelName == "" {
		return nil, errors.New("modelName cannot be empty")
	}
	if endpointURL == "" {
		return nil, errors.New("endpointURL cannot be empty")
	}
	if authType == "" {
		return nil, errors.New("authType cannot be empty")
	}

	authConfigBytes, err := r.marshalAuthConfig(authConfig)
	if err != nil {
		return nil, err
	}

	metadataBytes, err := r.marshalJSON(metadata)
	if err != nil {
		return nil, err
	}

	m, err := r.q.CreateModel(ctx, db.CreateModelParams{
		OrgID:          orgID,
		Provider:       provider,
		ModelName:      modelName,
		DeploymentName: deploymentName,
		EndpointUrl:    endpointURL,
		AuthType:       r.authTypeToString(authType),
		AuthConfig:     authConfigBytes,
		Metadata:       metadataBytes,
		Enabled:        enabled,
	})
	if err != nil {
		return nil, err
	}

	return &model.Model{
		ID:             m.ID,
		OrgID:          m.OrgID,
		Provider:       m.Provider,
		ModelName:      m.ModelName,
		DeploymentName: m.DeploymentName,
		EndpointURL:    m.EndpointUrl,
		AuthType:       authType,
		AuthConfig:     authConfig,
		Metadata:       metadata,
		Enabled:        m.Enabled,
		CreatedAt:      m.CreatedAt.Time,
		UpdatedAt:      m.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) Update(ctx context.Context, id uuid.UUID, provider, modelName string, deploymentName *string, endpointURL string, authType model.AuthType, authConfig model.AuthConfig, metadata map[string]any, enabled bool) (*model.Model, error) {
	if id == uuid.Nil {
		return nil, errors.New("id cannot be nil")
	}
	if provider == "" {
		return nil, errors.New("provider cannot be empty")
	}
	if modelName == "" {
		return nil, errors.New("modelName cannot be empty")
	}
	if endpointURL == "" {
		return nil, errors.New("endpointURL cannot be empty")
	}
	if authType == "" {
		return nil, errors.New("authType cannot be empty")
	}

	authConfigBytes, err := r.marshalAuthConfig(authConfig)
	if err != nil {
		return nil, err
	}

	metadataBytes, err := r.marshalJSON(metadata)
	if err != nil {
		return nil, err
	}

	m, err := r.q.UpdateModel(ctx, db.UpdateModelParams{
		ID:             id,
		Provider:       provider,
		ModelName:      modelName,
		DeploymentName: deploymentName,
		EndpointUrl:    endpointURL,
		AuthType:       r.authTypeToString(authType),
		AuthConfig:     authConfigBytes,
		Metadata:       metadataBytes,
		Enabled:        enabled,
	})
	if err != nil {
		return nil, err
	}

	return &model.Model{
		ID:             m.ID,
		OrgID:          m.OrgID,
		Provider:       m.Provider,
		ModelName:      m.ModelName,
		DeploymentName: m.DeploymentName,
		EndpointURL:    m.EndpointUrl,
		AuthType:       authType,
		AuthConfig:     authConfig,
		Metadata:       metadata,
		Enabled:        m.Enabled,
		CreatedAt:      m.CreatedAt.Time,
		UpdatedAt:      m.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id cannot be nil")
	}
	return r.q.DeleteModel(ctx, id)
}

func (r *postgresRepo) Enable(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id cannot be nil")
	}
	// Check if model exists
	_, err := r.q.GetModel(ctx, id)
	if err != nil {
		return err
	}
	return r.q.EnableModel(ctx, id)
}

func (r *postgresRepo) Disable(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id cannot be nil")
	}
	// Check if model exists
	_, err := r.q.GetModel(ctx, id)
	if err != nil {
		return err
	}
	return r.q.DisableModel(ctx, id)
}

func (r *postgresRepo) marshalJSON(data map[string]any) ([]byte, error) {
	if data == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(data)
}

func (r *postgresRepo) unmarshalJSON(data []byte) (map[string]any, error) {
	if len(data) == 0 {
		return make(map[string]any), nil
	}

	var result map[string]any
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *postgresRepo) marshalAuthConfig(authConfig model.AuthConfig) ([]byte, error) {
	return json.Marshal(authConfig)
}

func (r *postgresRepo) unmarshalAuthConfig(data []byte) (model.AuthConfig, error) {
	if len(data) == 0 {
		return model.AuthConfig{}, nil
	}

	var result model.AuthConfig
	err := json.Unmarshal(data, &result)
	if err != nil {
		return model.AuthConfig{}, err
	}
	return result, nil
}

func (r *postgresRepo) stringToAuthType(s string) model.AuthType {
	switch s {
	case "api_key":
		return model.AuthTypeAPIKey
	case "oauth2":
		return model.AuthTypeOAuth2
	case "azure_ad":
		return model.AuthTypeAzureAD
	default:
		return model.AuthType(s) // Fallback for unknown types
	}
}

func (r *postgresRepo) authTypeToString(at model.AuthType) string {
	return string(at)
}
