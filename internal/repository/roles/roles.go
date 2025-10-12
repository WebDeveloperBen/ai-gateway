package roles

import (
	"context"

	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
)

var DefaultRoleNames = []string{"owner", "admin", "member"}

type Repository interface {
	Create(ctx context.Context, arg db.CreateUserParams) (*model.User, error)
	FindBySubOrEmail(ctx context.Context, sub, email string) (*model.User, error)
	AssignRole(ctx context.Context, userID, roleID string, orgID uuid.UUID) error
}
