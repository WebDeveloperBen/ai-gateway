package organisations

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/insurgence-ai/llm-gateway/internal/db"
	"github.com/insurgence-ai/llm-gateway/internal/exceptions/pg"
	"github.com/insurgence-ai/llm-gateway/internal/model"
	repository "github.com/insurgence-ai/llm-gateway/internal/repository"
	"github.com/insurgence-ai/llm-gateway/internal/repository/roles"
)

type postgresRepo struct {
	Q *db.Queries
}

func NewPostgresRepo(q *db.Queries) Repository {
	return &postgresRepo{Q: q}
}

var handleDBError = pg.MakeErrorHandler("organisation")

func (r *postgresRepo) Create(ctx context.Context, name string) (*model.Organisation, error) {
	org, err := r.Q.CreateOrg(ctx, name)
	if err != nil {
		return nil, handleDBError(err)
	}

	for _, name := range roles.DefaultRoleNames {
		role, err := r.Q.FindRoleByName(ctx, name)
		if err != nil {
			return nil, handleDBError(fmt.Errorf("missing global role %q: %w", name, err))
		}
		_, err = r.Q.AssignRoleToOrg(ctx, db.AssignRoleToOrgParams{RoleID: role.ID, OrgID: org.ID})
		if err != nil {
			return nil, handleDBError(fmt.Errorf("error assigning role to org: %w", err))
		}
	}

	return &model.Organisation{
		ID:        org.ID.String(),
		Name:      org.Name,
		CreatedAt: org.CreatedAt.Time,
		UpdatedAt: org.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) FindByID(ctx context.Context, id string) (*model.Organisation, error) {
	uuidVal := repository.ParseUUID(id)
	if uuidVal == (uuid.UUID{}) {
		return nil, handleDBError(fmt.Errorf("invalid uuid"))
	}
	org, err := r.Q.FindOrgByID(ctx, uuidVal)
	if err != nil {
		return nil, handleDBError(err)
	}
	return &model.Organisation{
		ID:        org.ID.String(),
		Name:      org.Name,
		CreatedAt: org.CreatedAt.Time,
		UpdatedAt: org.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) EnsureMembership(ctx context.Context, orgID, userID uuid.UUID) error {
	return handleDBError(r.Q.EnsureOrgMembership(ctx, db.EnsureOrgMembershipParams{OrgID: orgID, UserID: userID}))
}

func (r *postgresRepo) FindRoleByName(ctx context.Context, name string) (*model.Role, error) {
	role, err := r.Q.FindRoleByName(ctx, name)
	if err != nil {
		return nil, handleDBError(err)
	}
	return &model.Role{
		ID:          role.ID.String(),
		Name:        role.Name,
		Description: repository.DerefString(role.Description),
		CreatedAt:   role.CreatedAt.Time,
	}, nil
}

func (r *postgresRepo) CreateRole(ctx context.Context, name, desc string) (*model.Role, error) {
	role, err := r.Q.CreateRole(ctx, db.CreateRoleParams{
		Name:        name,
		Description: &desc,
	})
	if err != nil {
		return nil, handleDBError(err)
	}

	return &model.Role{
		ID:          role.ID.String(),
		Name:        role.Name,
		Description: repository.DerefString(role.Description),
		CreatedAt:   role.CreatedAt.Time,
	}, nil
}

func (r *postgresRepo) AssignRole(ctx context.Context, orgID, roleID string) error {
	orgUUID := repository.ParseUUID(orgID)
	if orgUUID == (uuid.UUID{}) {
		return handleDBError(errors.New("invalid uuid"))
	}

	roleUUID := repository.ParseUUID(roleID)
	if orgUUID == (uuid.UUID{}) {
		return handleDBError(errors.New("invalid uuid"))
	}

	_, err := r.Q.AssignRoleToOrg(ctx, db.AssignRoleToOrgParams{
		OrgID:  orgUUID,
		RoleID: roleUUID,
	})
	if err != nil {
		return handleDBError(err)
	}

	return nil
}
