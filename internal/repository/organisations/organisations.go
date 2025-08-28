package organisations

import (
	"context"

	"github.com/google/uuid"
	"github.com/insurgence-ai/llm-gateway/internal/model"
)

type OrgRepository interface {
	Create(ctx context.Context, name string) (*model.Organisation, error)
	FindByID(ctx context.Context, id string) (*model.Organisation, error)
	EnsureRole(ctx context.Context, orgID string, roleName string, desc string) (*model.Role, error)
	EnsureOrgMembership(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error
}
