package application_configs

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

func (r *postgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.ApplicationConfig, error) {
	cfg, err := r.q.GetApplicationConfig(ctx, id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	var config map[string]any
	if err := json.Unmarshal(cfg.Config, &config); err != nil {
		return nil, err
	}

	return &model.ApplicationConfig{
		ID:          cfg.ID,
		AppID:       cfg.AppID,
		OrgID:       cfg.OrgID,
		Environment: cfg.Environment,
		Config:      config,
		CreatedAt:   cfg.CreatedAt.Time,
		UpdatedAt:   cfg.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) GetByEnv(ctx context.Context, appID uuid.UUID, environment string) (*model.ApplicationConfig, error) {
	cfg, err := r.q.GetApplicationConfigByEnv(ctx, db.GetApplicationConfigByEnvParams{
		AppID:       appID,
		Environment: environment,
	})
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	var config map[string]any
	if err := json.Unmarshal(cfg.Config, &config); err != nil {
		return nil, err
	}

	return &model.ApplicationConfig{
		ID:          cfg.ID,
		AppID:       cfg.AppID,
		OrgID:       cfg.OrgID,
		Environment: cfg.Environment,
		Config:      config,
		CreatedAt:   cfg.CreatedAt.Time,
		UpdatedAt:   cfg.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) ListByAppID(ctx context.Context, appID uuid.UUID) ([]*model.ApplicationConfig, error) {
	configs, err := r.q.ListApplicationConfigs(ctx, appID)
	if err != nil {
		return nil, err
	}

	result := make([]*model.ApplicationConfig, len(configs))
	for i, cfg := range configs {
		var config map[string]any
		if err := json.Unmarshal(cfg.Config, &config); err != nil {
			return nil, err
		}

		result[i] = &model.ApplicationConfig{
			ID:          cfg.ID,
			AppID:       cfg.AppID,
			OrgID:       cfg.OrgID,
			Environment: cfg.Environment,
			Config:      config,
			CreatedAt:   cfg.CreatedAt.Time,
			UpdatedAt:   cfg.UpdatedAt.Time,
		}
	}
	return result, nil
}

func (r *postgresRepo) Create(ctx context.Context, appID, orgID uuid.UUID, environment string, config map[string]any) (*model.ApplicationConfig, error) {
	if appID == uuid.Nil {
		return nil, errors.New("appID cannot be nil")
	}
	if orgID == uuid.Nil {
		return nil, errors.New("orgID cannot be nil")
	}
	if environment == "" {
		return nil, errors.New("environment cannot be empty")
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	cfg, err := r.q.CreateApplicationConfig(ctx, db.CreateApplicationConfigParams{
		AppID:       appID,
		OrgID:       orgID,
		Environment: environment,
		Config:      configBytes,
	})
	if err != nil {
		return nil, err
	}

	var configMap map[string]any
	if err := json.Unmarshal(cfg.Config, &configMap); err != nil {
		return nil, err
	}

	return &model.ApplicationConfig{
		ID:          cfg.ID,
		AppID:       cfg.AppID,
		OrgID:       cfg.OrgID,
		Environment: cfg.Environment,
		Config:      configMap,
		CreatedAt:   cfg.CreatedAt.Time,
		UpdatedAt:   cfg.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) Update(ctx context.Context, id uuid.UUID, config map[string]any) (*model.ApplicationConfig, error) {
	if id == uuid.Nil {
		return nil, errors.New("id cannot be nil")
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	cfg, err := r.q.UpdateApplicationConfig(ctx, db.UpdateApplicationConfigParams{
		ID:     id,
		Config: configBytes,
	})
	if err != nil {
		return nil, err
	}

	var configMap map[string]interface{}
	if err := json.Unmarshal(cfg.Config, &configMap); err != nil {
		return nil, err
	}

	return &model.ApplicationConfig{
		ID:          cfg.ID,
		AppID:       cfg.AppID,
		OrgID:       cfg.OrgID,
		Environment: cfg.Environment,
		Config:      configMap,
		CreatedAt:   cfg.CreatedAt.Time,
		UpdatedAt:   cfg.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id cannot be nil")
	}
	return r.q.DeleteApplicationConfig(ctx, id)
}
