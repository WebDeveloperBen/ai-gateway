package application_configs

import (
	"context"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
)

type Reader interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.ApplicationConfig, error)
	GetByEnv(ctx context.Context, appID uuid.UUID, environment string) (*model.ApplicationConfig, error)
	ListByAppID(ctx context.Context, appID uuid.UUID) ([]*model.ApplicationConfig, error)
}

type Writer interface {
	Create(ctx context.Context, appID, orgID uuid.UUID, environment string, config map[string]interface{}) (*model.ApplicationConfig, error)
	Update(ctx context.Context, id uuid.UUID, config map[string]interface{}) (*model.ApplicationConfig, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Repository interface {
	Reader
	Writer
}
