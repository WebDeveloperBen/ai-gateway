package applications

import (
	"context"
	"errors"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/repository/applications"
	"github.com/google/uuid"
)

type ApplicationsService interface {
	CreateApplication(ctx context.Context, orgID uuid.UUID, req CreateApplicationRequest) (*Application, error)
	GetApplication(ctx context.Context, id uuid.UUID) (*Application, error)
	GetApplicationByName(ctx context.Context, orgID uuid.UUID, name string) (*Application, error)
	ListApplications(ctx context.Context, orgID uuid.UUID) ([]*Application, error)
	UpdateApplication(ctx context.Context, id uuid.UUID, req UpdateApplicationRequest) (*Application, error)
	DeleteApplication(ctx context.Context, id uuid.UUID) error
}

type applicationsService struct {
	repo applications.Repository
}

func NewService(repo applications.Repository) ApplicationsService {
	return &applicationsService{repo: repo}
}

func (s *applicationsService) CreateApplication(ctx context.Context, orgID uuid.UUID, req CreateApplicationRequest) (*Application, error) {
	if req.Name == "" {
		return nil, errors.New("application name is required")
	}

	app, err := s.repo.Create(ctx, orgID, req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	return s.convertToAPI(app), nil
}

func (s *applicationsService) GetApplication(ctx context.Context, id uuid.UUID) (*Application, error) {
	app, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.convertToAPI(app), nil
}

func (s *applicationsService) GetApplicationByName(ctx context.Context, orgID uuid.UUID, name string) (*Application, error) {
	app, err := s.repo.GetByName(ctx, orgID, name)
	if err != nil {
		return nil, err
	}

	return s.convertToAPI(app), nil
}

func (s *applicationsService) ListApplications(ctx context.Context, orgID uuid.UUID) ([]*Application, error) {
	apps, err := s.repo.ListByOrgID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	result := make([]*Application, len(apps))
	for i, app := range apps {
		result[i] = s.convertToAPI(app)
	}
	return result, nil
}

func (s *applicationsService) UpdateApplication(ctx context.Context, id uuid.UUID, req UpdateApplicationRequest) (*Application, error) {
	if req.Name == "" {
		return nil, errors.New("application name is required")
	}

	app, err := s.repo.Update(ctx, id, req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	return s.convertToAPI(app), nil
}

func (s *applicationsService) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *applicationsService) convertToAPI(app *model.Application) *Application {
	return &Application{
		ID:          app.ID.String(),
		OrgID:       app.OrgID.String(),
		Name:        app.Name,
		Description: app.Description,
		CreatedAt:   app.CreatedAt,
		UpdatedAt:   app.UpdatedAt,
	}
}
