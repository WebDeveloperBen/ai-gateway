package catalog

import (
	"context"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
)

type Reader interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Model, error)
	GetByProviderAndName(ctx context.Context, orgID uuid.UUID, provider, modelName string) (*model.Model, error)
	ListByOrgID(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*model.Model, error)
	ListEnabledByOrgID(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*model.Model, error)
}

type Writer interface {
	Create(ctx context.Context, orgID uuid.UUID, provider, modelName string, deploymentName *string, endpointURL string, authType model.AuthType, authConfig model.AuthConfig, metadata map[string]any, enabled bool) (*model.Model, error)
	Update(ctx context.Context, id uuid.UUID, provider, modelName string, deploymentName *string, endpointURL string, authType model.AuthType, authConfig model.AuthConfig, metadata map[string]any, enabled bool) (*model.Model, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Enable(ctx context.Context, id uuid.UUID) error
	Disable(ctx context.Context, id uuid.UUID) error
}

type Repository interface {
	Reader
	Writer
}
