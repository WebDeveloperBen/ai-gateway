package application_configs

import (
	"context"
	"errors"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/repository/application_configs"
	"github.com/google/uuid"
)

type ApplicationConfigsService interface {
	CreateApplicationConfig(ctx context.Context, orgID uuid.UUID, req CreateApplicationConfigBody) (*ApplicationConfig, error)
	GetApplicationConfig(ctx context.Context, id uuid.UUID) (*ApplicationConfig, error)
	GetApplicationConfigByEnv(ctx context.Context, appID uuid.UUID, environment string) (*ApplicationConfig, error)
	ListApplicationConfigs(ctx context.Context, appID uuid.UUID) ([]*ApplicationConfig, error)
	UpdateApplicationConfig(ctx context.Context, id uuid.UUID, req UpdateApplicationConfigBody) (*ApplicationConfig, error)
	DeleteApplicationConfig(ctx context.Context, id uuid.UUID) error
}

type applicationConfigsService struct {
	repo application_configs.Repository
}

func NewService(repo application_configs.Repository) ApplicationConfigsService {
	return &applicationConfigsService{repo: repo}
}

func (s *applicationConfigsService) CreateApplicationConfig(ctx context.Context, orgID uuid.UUID, req CreateApplicationConfigBody) (*ApplicationConfig, error) {
	appID, err := uuid.Parse(req.AppID)
	if err != nil {
		return nil, errors.New("invalid application ID")
	}

	if req.Environment == "" {
		return nil, errors.New("environment is required")
	}

	if req.Config == nil {
		return nil, errors.New("config is required")
	}

	cfg, err := s.repo.Create(ctx, appID, orgID, req.Environment, req.Config)
	if err != nil {
		return nil, err
	}

	return s.convertToAPI(cfg), nil
}

func (s *applicationConfigsService) GetApplicationConfig(ctx context.Context, id uuid.UUID) (*ApplicationConfig, error) {
	cfg, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.convertToAPI(cfg), nil
}

func (s *applicationConfigsService) GetApplicationConfigByEnv(ctx context.Context, appID uuid.UUID, environment string) (*ApplicationConfig, error) {
	cfg, err := s.repo.GetByEnv(ctx, appID, environment)
	if err != nil {
		return nil, err
	}

	return s.convertToAPI(cfg), nil
}

func (s *applicationConfigsService) ListApplicationConfigs(ctx context.Context, appID uuid.UUID) ([]*ApplicationConfig, error) {
	configs, err := s.repo.ListByAppID(ctx, appID)
	if err != nil {
		return nil, err
	}

	result := make([]*ApplicationConfig, len(configs))
	for i, cfg := range configs {
		result[i] = s.convertToAPI(cfg)
	}
	return result, nil
}

func (s *applicationConfigsService) UpdateApplicationConfig(ctx context.Context, id uuid.UUID, req UpdateApplicationConfigBody) (*ApplicationConfig, error) {
	if req.Config == nil {
		return nil, errors.New("config is required")
	}

	cfg, err := s.repo.Update(ctx, id, req.Config)
	if err != nil {
		return nil, err
	}

	return s.convertToAPI(cfg), nil
}

func (s *applicationConfigsService) DeleteApplicationConfig(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *applicationConfigsService) convertToAPI(cfg *model.ApplicationConfig) *ApplicationConfig {
	return &ApplicationConfig{
		ID:          cfg.ID.String(),
		AppID:       cfg.AppID.String(),
		OrgID:       cfg.OrgID.String(),
		Environment: cfg.Environment,
		Config:      cfg.Config,
		CreatedAt:   cfg.CreatedAt,
		UpdatedAt:   cfg.UpdatedAt,
	}
}
