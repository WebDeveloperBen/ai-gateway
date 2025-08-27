package organisations

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/insurgence-ai/llm-gateway/internal/db"
	"github.com/insurgence-ai/llm-gateway/internal/model"
	repository "github.com/insurgence-ai/llm-gateway/internal/repository"
)

type postgresRepo struct {
	Q *db.Queries
}

func NewPostgresRepo(q *db.Queries) OrgRepository {
	return &postgresRepo{Q: q}
}

func (r *postgresRepo) Create(ctx context.Context, name string) (*model.Organisation, error) {
	org, err := r.Q.CreateOrg(ctx, name)
	if err != nil {
		return nil, err
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
		return nil, errors.New("invalid uuid")
	}
	org, err := r.Q.FindOrgByID(ctx, uuidVal)
	if err != nil {
		return nil, err
	}
	return &model.Organisation{
		ID:        org.ID.String(),
		Name:      org.Name,
		CreatedAt: org.CreatedAt.Time,
		UpdatedAt: org.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) EnsureRole(ctx context.Context, orgID, roleName, desc string) (*model.Role, error) {
	orgUUID := repository.ParseUUID(orgID)
	if orgUUID == (uuid.UUID{}) {
		return nil, errors.New("invalid org uuid")
	}
	role, err := r.Q.FindRoleByOrgAndName(ctx, db.FindRoleByOrgAndNameParams{OrgID: orgUUID, Name: roleName})
	if err == nil {
		return &model.Role{
			ID:          role.ID.String(),
			OrgID:       role.OrgID.String(),
			Name:        role.Name,
			Description: repository.DerefString(role.Description),
			CreatedAt:   role.CreatedAt.Time,
		}, nil
	}
	role, err = r.Q.CreateRole(ctx, db.CreateRoleParams{
		OrgID:       orgUUID,
		Name:        roleName,
		Description: &desc,
	})
	if err != nil {
		return nil, err
	}
	return &model.Role{
		ID:          role.ID.String(),
		OrgID:       role.OrgID.String(),
		Name:        role.Name,
		Description: repository.DerefString(role.Description),
		CreatedAt:   role.CreatedAt.Time,
	}, nil
}
