package users

import (
	"context"
	"github.com/insurgence-ai/llm-gateway/internal/model"
)

type UserServiceInterface interface {
	EnsureUserAndOrg(ctx context.Context, scoped model.ScopedToken) (*model.User, *model.Organisation, error)
}
