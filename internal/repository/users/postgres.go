package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/repository"
	"github.com/google/uuid"
)

type postgresRepo struct {
	q *db.Queries
}

func NewPostgresRepo(q *db.Queries) Repository {
	return &postgresRepo{q: q}
}

func (r *postgresRepo) FindBySubOrEmail(ctx context.Context, sub, email string) (*model.User, error) {
	u, err := r.q.FindUserBySubOrEmail(ctx, db.FindUserBySubOrEmailParams{Sub: &sub, Email: email})
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:        u.ID.String(),
		OrgID:     u.OrgID.String(),
		Sub:       repository.DerefString(u.Sub),
		Email:     u.Email,
		Name:      repository.DerefString(u.Name),
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) Create(ctx context.Context, arg db.CreateUserParams) (*model.User, error) {
	u, err := r.q.CreateUser(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:        u.ID.String(),
		OrgID:     u.OrgID.String(),
		Sub:       *u.Sub,
		Email:     u.Email,
		Name:      *u.Name,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) AssignRole(ctx context.
	Context, userID, roleID string, orgID uuid.UUID,
) error {
	uid := repository.ParseUUID(userID)
	rid := repository.ParseUUID(roleID)
	if uid == uuid.Nil || rid == uuid.Nil || orgID == uuid.Nil {
		return errors.New("invalid uuid provided")
	}
	return r.q.AssignRoleToUser(ctx, db.
		AssignRoleToUserParams{
		UserID: uid, RoleID: rid, OrgID: orgID,
	})
}
