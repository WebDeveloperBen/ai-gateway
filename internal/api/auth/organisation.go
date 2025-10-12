package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/repository"
	"github.com/WebDeveloperBen/ai-gateway/internal/repository/organisations"
	"github.com/WebDeveloperBen/ai-gateway/internal/repository/users"
	"github.com/google/uuid"
)

type OrganisationServiceInterface interface {
	EnsureUserAndOrg(ctx context.Context, scoped model.ScopedToken) (*model.User, *model.Organisation, error)
	EnsureRole(ctx context.Context, orgID, roleName, desc string) (*model.Role, error)
}
type OrganisationService struct {
	orgRepo  organisations.Repository
	userRepo users.Repository
}

func NewOrganisationService(orgRepo organisations.Repository, userRepo users.Repository) OrganisationServiceInterface {
	return &OrganisationService{orgRepo: orgRepo, userRepo: userRepo}
}

// EnsureUserAndOrg finds or creates a user + org for OIDC login.
func (s *OrganisationService) EnsureUserAndOrg(ctx context.Context, scoped model.ScopedToken) (*model.User, *model.Organisation, error) {
	if scoped.Subject == "" || scoped.Email == "" || scoped.Name == "" {
		return nil, nil, fmt.Errorf("invalid claims: missing subject, email, or name: %+v", scoped)
	}

	user, err := s.userRepo.FindBySubOrEmail(ctx, scoped.Subject, scoped.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, fmt.Errorf("find user: %w", err)
	}

	// Existing user path
	if err == nil && user != nil {
		org, oerr := s.orgRepo.FindByID(ctx, user.OrgID)
		if oerr != nil {
			return nil, nil, fmt.Errorf("find org: %w", oerr)
		}
		return user, org, nil
	}

	// First login: create new org + bootstrap user
	org, err := s.orgRepo.Create(ctx, fmt.Sprintf("%s's Home", scoped.Name))
	if err != nil {
		return nil, nil, fmt.Errorf("create org: %w", err)
	}

	user, err = s.userRepo.Create(ctx, db.CreateUserParams{
		OrgID: repository.ParseUUID(org.ID),
		Sub:   &scoped.Subject,
		Email: scoped.Email,
		Name:  &scoped.Name,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create user: %w", err)
	}
	// 1. Ensure global owner role exists
	ownerRole, err := s.EnsureRole(ctx, org.ID, "owner", "Organisation owner")
	if err != nil {
		return nil, nil, fmt.Errorf("ensure role(owner): %w", err)
	}
	// 2. Assign owner role to user
	if err := s.userRepo.AssignRole(ctx, user.ID, ownerRole.ID, repository.ParseUUID(org.ID)); err != nil {
		return nil, nil, fmt.Errorf("assign role(owner)): %w", err)
	}

	if err := s.orgRepo.EnsureMembership(ctx, repository.ParseUUID(org.ID), repository.ParseUUID(user.ID)); err != nil {
		return nil, nil, fmt.Errorf("ensure membership: %w", err)
	}

	return user, org, nil
}

func (s *OrganisationService) EnsureRole(ctx context.Context, orgID, roleName, desc string) (*model.Role, error) {
	orgUUID := repository.ParseUUID(orgID)
	if orgUUID == (uuid.UUID{}) {
		fmt.Printf("[EnsureRole] Invalid org uuid: %v\n", orgID)
		return nil, errors.New("invalid org uuid")
	}
	// Find or create global role
	role, err := s.orgRepo.FindRoleByName(ctx, roleName)
	fmt.Printf("[EnsureRole] FindRoleByName(%s): role=%+v err=%v\n", roleName, role, err)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			role, err = s.orgRepo.CreateRole(ctx, roleName, desc)
			fmt.Printf("[EnsureRole] CreateRole(%s): role=%+v err=%v\n", roleName, role, err)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Assign the role to the org via org_roles
	err = s.orgRepo.AssignRole(ctx, orgID, role.ID)
	fmt.Printf("[EnsureRole] AssignRoleToOrg(orgID=%v, roleID=%v): err=%v\n", orgUUID, role.ID, err)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return &model.Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
	}, nil
}
