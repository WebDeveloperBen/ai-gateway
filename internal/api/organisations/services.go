package organisations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/insurgence-ai/llm-gateway/internal/db"
	"github.com/insurgence-ai/llm-gateway/internal/model"
	"github.com/insurgence-ai/llm-gateway/internal/repository"
	"github.com/insurgence-ai/llm-gateway/internal/repository/organisations"
	"github.com/insurgence-ai/llm-gateway/internal/repository/users"
)

type OrganisationServiceInterface interface {
	EnsureUserAndOrg(ctx context.Context, scoped model.ScopedToken) (*model.User, *model.Organisation, error)
}
type OrganisationService struct {
	orgRepo  organisations.OrgRepository
	userRepo users.Repository
}

func NewService(orgRepo organisations.OrgRepository, userRepo users.Repository) *OrganisationService {
	return &OrganisationService{orgRepo: orgRepo, userRepo: userRepo}
}

// EnsureUserAndOrg finds or creates a user + org for OIDC login
func (s *OrganisationService) EnsureUserAndOrg(ctx context.Context, scoped model.ScopedToken) (*model.User, *model.Organisation, error) {
	if scoped.Subject == "" || scoped.Email == "" || scoped.Name == "" {
		return nil, nil, fmt.Errorf("invalid claims: missing subject, email, or name: %+v", scoped)
	}
	user, err := s.userRepo.FindBySubOrEmail(ctx, scoped.Subject, scoped.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, nil, err
	}

	if user != nil {
		org, err := s.orgRepo.FindByID(ctx, user.OrgID)
		if err != nil {
			return nil, nil, err
		}
		return user, org, nil
	}

	// First login: create new org + bootstrap user
	org, err := s.orgRepo.Create(ctx,
		fmt.Sprintf("%s’s Org", scoped.Name),
	)
	if err != nil {
		return nil, nil, err
	}

	user, err = s.userRepo.Create(ctx, db.CreateUserParams{
		OrgID: repository.ParseUUID(org.ID),
		Sub:   &scoped.Subject,
		Email: scoped.Email,
		Name:  &scoped.Name,
	})
	if err != nil {
		return nil, nil, err
	}

	// bootstrap owner role
	ownerRole, err := s.orgRepo.EnsureRole(ctx, org.ID, "owner", "Organisation owner")
	if err != nil {
		return nil, nil, err
	}
	err = s.userRepo.AssignRole(ctx, user.ID, ownerRole.ID, repository.ParseUUID(org.ID))
	if err != nil {
		return nil, nil, err
	}

	return user, org, nil
}
