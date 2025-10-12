package policies

import (
	"context"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
)

type Reader interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Policy, error)
	ListByAppID(ctx context.Context, appID uuid.UUID, limit, offset int) ([]*model.Policy, error)
	ListEnabledByAppID(ctx context.Context, appID uuid.UUID, limit, offset int) ([]*model.Policy, error)
	GetByType(ctx context.Context, appID uuid.UUID, policyType model.PolicyType) ([]*model.Policy, error)
	GetAppsForPolicy(ctx context.Context, policyID uuid.UUID) ([]*model.Application, error)
}

type Writer interface {
	Create(ctx context.Context, orgID uuid.UUID, policyType model.PolicyType, config map[string]any, enabled bool) (*model.Policy, error)
	Update(ctx context.Context, id uuid.UUID, policyType model.PolicyType, config map[string]any, enabled bool) (*model.Policy, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Enable(ctx context.Context, id uuid.UUID) error
	Disable(ctx context.Context, id uuid.UUID) error
	AttachToApp(ctx context.Context, policyID, appID uuid.UUID) error
	DetachFromApp(ctx context.Context, policyID, appID uuid.UUID) error
}

type Repository interface {
	Reader
	Writer
}
