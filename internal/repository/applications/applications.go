package applications

import (
	"context"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
)

type Reader interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Application, error)
	GetByName(ctx context.Context, orgID uuid.UUID, name string) (*model.Application, error)
	ListByOrgID(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*model.Application, error)
}

type Writer interface {
	Create(ctx context.Context, orgID uuid.UUID, name string, description *string) (*model.Application, error)
	Update(ctx context.Context, id uuid.UUID, name string, description *string) (*model.Application, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Repository interface {
	Reader
	Writer
}
