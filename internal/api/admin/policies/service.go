package policies

import (
	"context"
	"errors"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/repository/policies"
	"github.com/google/uuid"
)

type PoliciesService interface {
	CreatePolicy(ctx context.Context, orgID uuid.UUID, req CreatePolicyRequest) (*Policy, error)
	GetPolicy(ctx context.Context, id uuid.UUID) (*Policy, error)
	ListPolicies(ctx context.Context, appID uuid.UUID) ([]*Policy, error)
	ListEnabledPolicies(ctx context.Context, appID uuid.UUID) ([]*Policy, error)
	GetPoliciesByType(ctx context.Context, appID uuid.UUID, policyType model.PolicyType) ([]*Policy, error)
	UpdatePolicy(ctx context.Context, id uuid.UUID, req UpdatePolicyRequest) (*Policy, error)
	DeletePolicy(ctx context.Context, id uuid.UUID) error
	EnablePolicy(ctx context.Context, id uuid.UUID) error
	DisablePolicy(ctx context.Context, id uuid.UUID) error
}

type policiesService struct {
	repo policies.Repository
}

func NewService(repo policies.Repository) PoliciesService {
	return &policiesService{repo: repo}
}

func (s *policiesService) CreatePolicy(ctx context.Context, orgID uuid.UUID, req CreatePolicyRequest) (*Policy, error) {
	if req.AppID == "" {
		return nil, errors.New("app_id is required")
	}
	if req.PolicyType == "" {
		return nil, errors.New("policy_type is required")
	}
	if req.Config == nil {
		return nil, errors.New("config is required")
	}

	appID, err := uuid.Parse(req.AppID)
	if err != nil {
		return nil, errors.New("invalid app_id")
	}

	policy, err := s.repo.Create(ctx, orgID, appID, req.PolicyType, req.Config, req.Enabled)
	if err != nil {
		return nil, err
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

func (s *policiesService) ListPolicies(ctx context.Context, appID uuid.UUID) ([]*Policy, error) {
	policies, err := s.repo.ListByAppID(ctx, appID)
	if err != nil {
		return nil, err
	}

	result := make([]*Policy, len(policies))
	for i, policy := range policies {
		result[i] = s.convertToAPI(policy)
	}
	return result, nil
}

func (s *policiesService) ListEnabledPolicies(ctx context.Context, appID uuid.UUID) ([]*Policy, error) {
	policies, err := s.repo.ListEnabledByAppID(ctx, appID)
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

func (s *policiesService) UpdatePolicy(ctx context.Context, id uuid.UUID, req UpdatePolicyRequest) (*Policy, error) {
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

func (s *policiesService) convertToAPI(policy *model.Policy) *Policy {
	return &Policy{
		ID:         policy.ID.String(),
		OrgID:      policy.OrgID.String(),
		AppID:      policy.AppID.String(),
		PolicyType: policy.PolicyType,
		Config:     policy.Config,
		Enabled:    policy.Enabled,
		CreatedAt:  policy.CreatedAt,
		UpdatedAt:  policy.UpdatedAt,
	}
}
