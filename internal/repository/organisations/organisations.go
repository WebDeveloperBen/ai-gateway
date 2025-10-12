package organisations

import (
	"context"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, name string) (*model.Organisation, error)
	FindByID(ctx context.Context, id string) (*model.Organisation, error)
	FindRoleByName(ctx context.Context, name string) (*model.Role, error)
	CreateRole(ctx context.Context, name, desc string) (*model.Role, error)
	AssignRole(ctx context.Context, orgID, roleID string) error
	EnsureMembership(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error
}
