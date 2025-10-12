package policies

import (
	"context"
	"errors"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/repository/policies"
	"github.com/google/uuid"
)

type PoliciesService interface {
	CreatePolicy(ctx context.Context, orgID uuid.UUID, req CreatePolicyBody) (*Policy, error)
	GetPolicy(ctx context.Context, id uuid.UUID) (*Policy, error)
	ListPolicies(ctx context.Context, appID uuid.UUID, limit, offset int) ([]*Policy, error)
	ListEnabledPolicies(ctx context.Context, appID uuid.UUID, limit, offset int) ([]*Policy, error)
	GetPoliciesByType(ctx context.Context, appID uuid.UUID, policyType model.PolicyType) ([]*Policy, error)
	UpdatePolicy(ctx context.Context, id uuid.UUID, req UpdatePolicyBody) (*Policy, error)
	DeletePolicy(ctx context.Context, id uuid.UUID) error
	EnablePolicy(ctx context.Context, id uuid.UUID) error
	DisablePolicy(ctx context.Context, id uuid.UUID) error
	AttachPolicyToApp(ctx context.Context, policyID, appID uuid.UUID) error
	DetachPolicyFromApp(ctx context.Context, policyID, appID uuid.UUID) error
	GetAppsForPolicy(ctx context.Context, policyID uuid.UUID) ([]*model.Application, error)
}

type policiesService struct {
	repo policies.Repository
}

func NewService(repo policies.Repository) PoliciesService {
	return &policiesService{repo: repo}
}

func (s *policiesService) CreatePolicy(ctx context.Context, orgID uuid.UUID, req CreatePolicyBody) (*Policy, error) {
	if req.PolicyType == "" {
		return nil, errors.New("policy_type is required")
	}
	if req.Config == nil {
		return nil, errors.New("config is required")
	}

	policy, err := s.repo.Create(ctx, orgID, req.PolicyType, req.Config, req.Enabled)
	if err != nil {
		return nil, err
	}

	// If appID is provided, attach the policy to the app
	if req.AppID != "" {
		appID, err := uuid.Parse(req.AppID)
		if err != nil {
			return nil, errors.New("invalid app_id")
		}
		if err := s.repo.AttachToApp(ctx, policy.ID, appID); err != nil {
			return nil, err
		}
	}

	return s.convertToAPI(policy), nil
}

func (s *policiesService) GetPolicy(ctx context.Context, id uuid.UUID) (*Policy, error) {
	policy, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.convertToAPI(policy), nil
}

func (s *policiesService) ListPolicies(ctx context.Context, appID uuid.UUID, limit, offset int) ([]*Policy, error) {
	policies, err := s.repo.ListByAppID(ctx, appID, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*Policy, len(policies))
	for i, policy := range policies {
		result[i] = s.convertToAPI(policy)
	}
	return result, nil
}

func (s *policiesService) ListEnabledPolicies(ctx context.Context, appID uuid.UUID, limit, offset int) ([]*Policy, error) {
	policies, err := s.repo.ListEnabledByAppID(ctx, appID, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*Policy, len(policies))
	for i, policy := range policies {
		result[i] = s.convertToAPI(policy)
	}
	return result, nil
}

func (s *policiesService) GetPoliciesByType(ctx context.Context, appID uuid.UUID, policyType model.PolicyType) ([]*Policy, error) {
	policies, err := s.repo.GetByType(ctx, appID, policyType)
	if err != nil {
		return nil, err
	}

	result := make([]*Policy, len(policies))
	for i, policy := range policies {
		result[i] = s.convertToAPI(policy)
	}
	return result, nil
}

func (s *policiesService) UpdatePolicy(ctx context.Context, id uuid.UUID, req UpdatePolicyBody) (*Policy, error) {
	if req.PolicyType == "" {
		return nil, errors.New("policy_type is required")
	}
	if req.Config == nil {
		return nil, errors.New("config is required")
	}

	policy, err := s.repo.Update(ctx, id, req.PolicyType, req.Config, req.Enabled)
	if err != nil {
		return nil, err
	}

	return s.convertToAPI(policy), nil
}

func (s *policiesService) DeletePolicy(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *policiesService) EnablePolicy(ctx context.Context, id uuid.UUID) error {
	return s.repo.Enable(ctx, id)
}

func (s *policiesService) DisablePolicy(ctx context.Context, id uuid.UUID) error {
	return s.repo.Disable(ctx, id)
}

func (s *policiesService) AttachPolicyToApp(ctx context.Context, policyID, appID uuid.UUID) error {
	return s.repo.AttachToApp(ctx, policyID, appID)
}

func (s *policiesService) DetachPolicyFromApp(ctx context.Context, policyID, appID uuid.UUID) error {
	return s.repo.DetachFromApp(ctx, policyID, appID)
}

func (s *policiesService) GetAppsForPolicy(ctx context.Context, policyID uuid.UUID) ([]*model.Application, error) {
	return s.repo.GetAppsForPolicy(ctx, policyID)
}

func (s *policiesService) convertToAPI(policy *model.Policy) *Policy {
	return &Policy{
		ID:         policy.ID.String(),
		OrgID:      policy.OrgID.String(),
		PolicyType: policy.PolicyType,
		Config:     policy.Config,
		Enabled:    policy.Enabled,
		CreatedAt:  policy.CreatedAt,
		UpdatedAt:  policy.UpdatedAt,
	}
}
