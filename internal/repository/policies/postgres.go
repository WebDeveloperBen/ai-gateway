package policies

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
)

type postgresRepo struct {
	q *db.Queries
}

func NewPostgresRepo(q *db.Queries) Repository {
	return &postgresRepo{q: q}
}

func (r *postgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Policy, error) {
	policy, err := r.q.GetPolicy(ctx, id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	config, err := r.unmarshalConfig(policy.Config)
	if err != nil {
		return nil, err
	}

	return &model.Policy{
		ID:         policy.ID,
		OrgID:      policy.OrgID,
		AppID:      policy.AppID,
		PolicyType: model.PolicyType(policy.PolicyType),
		Config:     config,
		Enabled:    policy.Enabled,
		CreatedAt:  policy.CreatedAt.Time,
		UpdatedAt:  policy.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) ListByAppID(ctx context.Context, appID uuid.UUID) ([]*model.Policy, error) {
	policies, err := r.q.ListPolicies(ctx, appID)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Policy, len(policies))
	for i, policy := range policies {
		config, err := r.unmarshalConfig(policy.Config)
		if err != nil {
			return nil, err
		}

		result[i] = &model.Policy{
			ID:         policy.ID,
			OrgID:      policy.OrgID,
			AppID:      policy.AppID,
			PolicyType: model.PolicyType(policy.PolicyType),
			Config:     config,
			Enabled:    policy.Enabled,
			CreatedAt:  policy.CreatedAt.Time,
			UpdatedAt:  policy.UpdatedAt.Time,
		}
	}
	return result, nil
}

func (r *postgresRepo) ListEnabledByAppID(ctx context.Context, appID uuid.UUID) ([]*model.Policy, error) {
	policies, err := r.q.ListEnabledPolicies(ctx, appID)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Policy, len(policies))
	for i, policy := range policies {
		config, err := r.unmarshalConfig(policy.Config)
		if err != nil {
			return nil, err
		}

		result[i] = &model.Policy{
			ID:         policy.ID,
			OrgID:      policy.OrgID,
			AppID:      policy.AppID,
			PolicyType: model.PolicyType(policy.PolicyType),
			Config:     config,
			Enabled:    policy.Enabled,
			CreatedAt:  policy.CreatedAt.Time,
			UpdatedAt:  policy.UpdatedAt.Time,
		}
	}
	return result, nil
}

func (r *postgresRepo) GetByType(ctx context.Context, appID uuid.UUID, policyType model.PolicyType) ([]*model.Policy, error) {
	policies, err := r.q.GetPoliciesByType(ctx, db.GetPoliciesByTypeParams{
		AppID:      appID,
		PolicyType: string(policyType),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*model.Policy, len(policies))
	for i, policy := range policies {
		config, err := r.unmarshalConfig(policy.Config)
		if err != nil {
			return nil, err
		}

		result[i] = &model.Policy{
			ID:         policy.ID,
			OrgID:      policy.OrgID,
			AppID:      policy.AppID,
			PolicyType: model.PolicyType(policy.PolicyType),
			Config:     config,
			Enabled:    policy.Enabled,
			CreatedAt:  policy.CreatedAt.Time,
			UpdatedAt:  policy.UpdatedAt.Time,
		}
	}
	return result, nil
}

func (r *postgresRepo) Create(ctx context.Context, orgID, appID uuid.UUID, policyType model.PolicyType, config map[string]any, enabled bool) (*model.Policy, error) {
	if orgID == uuid.Nil {
		return nil, errors.New("orgID cannot be nil")
	}
	if appID == uuid.Nil {
		return nil, errors.New("appID cannot be nil")
	}
	if policyType == "" {
		return nil, errors.New("policyType cannot be empty")
	}

	configBytes, err := r.marshalConfig(config)
	if err != nil {
		return nil, err
	}

	policy, err := r.q.CreatePolicy(ctx, db.CreatePolicyParams{
		OrgID:      orgID,
		AppID:      appID,
		PolicyType: string(policyType),
		Config:     configBytes,
		Enabled:    enabled,
	})
	if err != nil {
		return nil, err
	}

	return &model.Policy{
		ID:         policy.ID,
		OrgID:      policy.OrgID,
		AppID:      policy.AppID,
		PolicyType: model.PolicyType(policy.PolicyType),
		Config:     config,
		Enabled:    policy.Enabled,
		CreatedAt:  policy.CreatedAt.Time,
		UpdatedAt:  policy.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) Update(ctx context.Context, id uuid.UUID, policyType model.PolicyType, config map[string]any, enabled bool) (*model.Policy, error) {
	if id == uuid.Nil {
		return nil, errors.New("id cannot be nil")
	}
	if policyType == "" {
		return nil, errors.New("policyType cannot be empty")
	}

	configBytes, err := r.marshalConfig(config)
	if err != nil {
		return nil, err
	}

	policy, err := r.q.UpdatePolicy(ctx, db.UpdatePolicyParams{
		ID:         id,
		PolicyType: string(policyType),
		Config:     configBytes,
		Enabled:    enabled,
	})
	if err != nil {
		return nil, err
	}

	return &model.Policy{
		ID:         policy.ID,
		OrgID:      policy.OrgID,
		AppID:      policy.AppID,
		PolicyType: model.PolicyType(policy.PolicyType),
		Config:     config,
		Enabled:    policy.Enabled,
		CreatedAt:  policy.CreatedAt.Time,
		UpdatedAt:  policy.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id cannot be nil")
	}
	return r.q.DeletePolicy(ctx, id)
}

func (r *postgresRepo) Enable(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id cannot be nil")
	}
	return r.q.EnablePolicy(ctx, id)
}

func (r *postgresRepo) Disable(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id cannot be nil")
	}
	return r.q.DisablePolicy(ctx, id)
}

func (r *postgresRepo) marshalConfig(config map[string]any) ([]byte, error) {
	if config == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(config)
}

func (r *postgresRepo) unmarshalConfig(data []byte) (map[string]any, error) {
	if len(data) == 0 {
		return make(map[string]any), nil
	}

	var config map[string]any
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
