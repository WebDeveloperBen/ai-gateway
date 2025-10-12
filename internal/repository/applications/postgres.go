package applications

import (
	"context"
	"database/sql"
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

func (r *postgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	app, err := r.q.GetApplication(ctx, id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &model.Application{
		ID:          app.ID,
		OrgID:       app.OrgID,
		Name:        app.Name,
		Description: app.Description,
		CreatedAt:   app.CreatedAt.Time,
		UpdatedAt:   app.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) GetByName(ctx context.Context, orgID uuid.UUID, name string) (*model.Application, error) {
	app, err := r.q.GetApplicationByName(ctx, db.GetApplicationByNameParams{
		OrgID: orgID,
		Name:  name,
	})
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &model.Application{
		ID:          app.ID,
		OrgID:       app.OrgID,
		Name:        app.Name,
		Description: app.Description,
		CreatedAt:   app.CreatedAt.Time,
		UpdatedAt:   app.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) ListByOrgID(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*model.Application, error) {
	apps, err := r.q.ListApplications(ctx, db.ListApplicationsParams{
		OrgID:  orgID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*model.Application, len(apps))
	for i, app := range apps {
		result[i] = &model.Application{
			ID:          app.ID,
			OrgID:       app.OrgID,
			Name:        app.Name,
			Description: app.Description,
			CreatedAt:   app.CreatedAt.Time,
			UpdatedAt:   app.UpdatedAt.Time,
		}
	}
	return result, nil
}

func (r *postgresRepo) Create(ctx context.Context, orgID uuid.UUID, name string, description *string) (*model.Application, error) {
	if orgID == uuid.Nil {
		return nil, errors.New("orgID cannot be nil")
	}
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	app, err := r.q.CreateApplication(ctx, db.CreateApplicationParams{
		OrgID:       orgID,
		Name:        name,
		Description: description,
	})
	if err != nil {
		return nil, err
	}
	return &model.Application{
		ID:          app.ID,
		OrgID:       app.OrgID,
		Name:        app.Name,
		Description: app.Description,
		CreatedAt:   app.CreatedAt.Time,
		UpdatedAt:   app.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) Update(ctx context.Context, id uuid.UUID, name string, description *string) (*model.Application, error) {
	if id == uuid.Nil {
		return nil, errors.New("id cannot be nil")
	}
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	app, err := r.q.UpdateApplication(ctx, db.UpdateApplicationParams{
		ID:          id,
		Name:        name,
		Description: description,
	})
	if err != nil {
		return nil, err
	}
	return &model.Application{
		ID:          app.ID,
		OrgID:       app.OrgID,
		Name:        app.Name,
		Description: app.Description,
		CreatedAt:   app.CreatedAt.Time,
		UpdatedAt:   app.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id cannot be nil")
	}
	return r.q.DeleteApplication(ctx, id)
}
