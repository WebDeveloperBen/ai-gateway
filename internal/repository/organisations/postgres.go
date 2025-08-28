package organisations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/insurgence-ai/llm-gateway/internal/db"
	"github.com/insurgence-ai/llm-gateway/internal/model"
	repository "github.com/insurgence-ai/llm-gateway/internal/repository"
	"github.com/insurgence-ai/llm-gateway/internal/repository/roles"
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

	for _, name := range roles.DefaultRoleNames {
		role, err := r.Q.FindRoleByName(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("missing global role %q: %w", name, err)
		}
		_, err = r.Q.AssignRoleToOrg(ctx, db.AssignRoleToOrgParams{RoleID: role.ID, OrgID: org.ID})
		if err != nil {
			return nil, fmt.Errorf("error assigning role to org: %w", err)
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
		fmt.Printf("[EnsureRole] Invalid org uuid: %v\n", orgID)
		return nil, errors.New("invalid org uuid")
	}
	// Find or create global role
	role, err := r.Q.FindRoleByName(ctx, roleName)
	fmt.Printf("[EnsureRole] FindRoleByName(%s): role=%+v err=%v\n", roleName, role, err)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			role, err = r.Q.CreateRole(ctx, db.CreateRoleParams{
				Name:        roleName,
				Description: &desc,
			})
			fmt.Printf("[EnsureRole] CreateRole(%s): role=%+v err=%v\n", roleName, role, err)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	// Assign the role to the org via org_roles
	_, err = r.Q.AssignRoleToOrg(ctx, db.AssignRoleToOrgParams{
		OrgID:  orgUUID,
		RoleID: role.ID,
	})
	fmt.Printf("[EnsureRole] AssignRoleToOrg(orgID=%v, roleID=%v): err=%v\n", orgUUID, role.ID, err)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return &model.Role{
		ID:          role.ID.String(),
		OrgID:       orgID,
		Name:        role.Name,
		Description: repository.DerefString(role.Description),
		CreatedAt:   role.CreatedAt.Time,
	}, nil
}

func (r *postgresRepo) EnsureOrgMembership(ctx context.Context, orgID, userID uuid.UUID) error {
	return r.Q.EnsureOrgMembership(ctx, db.EnsureOrgMembershipParams{OrgID: orgID, UserID: userID})
}
