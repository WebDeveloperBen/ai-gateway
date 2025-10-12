package policies

import (
	"context"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
)

type Reader interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Policy, error)
	ListByAppID(ctx context.Context, appID uuid.UUID) ([]*model.Policy, error)
	ListEnabledByAppID(ctx context.Context, appID uuid.UUID) ([]*model.Policy, error)
	GetByType(ctx context.Context, appID uuid.UUID, policyType model.PolicyType) ([]*model.Policy, error)
}

type Writer interface {
	Create(ctx context.Context, orgID, appID uuid.UUID, policyType model.PolicyType, config map[string]any, enabled bool) (*model.Policy, error)
	Update(ctx context.Context, id uuid.UUID, policyType model.PolicyType, config map[string]any, enabled bool) (*model.Policy, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Enable(ctx context.Context, id uuid.UUID) error
	Disable(ctx context.Context, id uuid.UUID) error
}

type Repository interface {
	Reader
	Writer
}
