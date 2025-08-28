package roles

import (
	"context"

	"github.com/google/uuid"
	"github.com/insurgence-ai/llm-gateway/internal/db"
	"github.com/insurgence-ai/llm-gateway/internal/model"
)

var DefaultRoleNames = []string{"owner", "admin", "member"}

type Repository interface {
	Create(ctx context.Context, arg db.CreateUserParams) (*model.User, error)
	FindBySubOrEmail(ctx context.Context, sub, email string) (*model.User, error)
	AssignRole(ctx context.Context, userID, roleID string, orgID uuid.UUID) error
}
