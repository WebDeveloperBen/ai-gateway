package applications

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
)

type appRecord struct {
	app model.Application
}

type MemoryRepo struct {
	mu    sync.RWMutex
	store map[uuid.UUID]appRecord
}

var _ Repository = (*MemoryRepo)(nil)

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{store: make(map[uuid.UUID]appRecord)}
}

func (m *MemoryRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	m.mu.RLock()
	rec, ok := m.store[id]
	m.mu.RUnlock()
	if !ok {
		return nil, errors.New("application not found")
	}
	return &rec.app, nil
}

func (m *MemoryRepo) GetByName(ctx context.Context, orgID uuid.UUID, name string) (*model.Application, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, rec := range m.store {
		if rec.app.OrgID == orgID && rec.app.Name == name {
			return &rec.app, nil
		}
	}
	return nil, errors.New("application not found")
}

func (m *MemoryRepo) ListByOrgID(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*model.Application, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var apps []*model.Application
	count := 0
	for _, rec := range m.store {
		if rec.app.OrgID == orgID {
			if count >= offset && (limit == 0 || len(apps) < limit) {
				apps = append(apps, &rec.app)
			}
			count++
		}
	}
	return apps, nil
}

func (m *MemoryRepo) Create(ctx context.Context, orgID uuid.UUID, name string, description *string) (*model.Application, error) {
	if orgID == uuid.Nil {
		return nil, errors.New("orgID cannot be nil")
	}
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check for duplicate name in org
	for _, rec := range m.store {
		if rec.app.OrgID == orgID && rec.app.Name == name {
			return nil, errors.New("application name already exists in organization")
		}
	}

	now := time.Now()
	app := model.Application{
		ID:          uuid.New(),
		OrgID:       orgID,
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	m.store[app.ID] = appRecord{app: app}
	return &app, nil
}

func (m *MemoryRepo) Update(ctx context.Context, id uuid.UUID, name string, description *string) (*model.Application, error) {
	if id == uuid.Nil {
		return nil, errors.New("id cannot be nil")
	}
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	rec, ok := m.store[id]
	if !ok {
		return nil, errors.New("application not found")
	}

	// Check for duplicate name in org (excluding current app)
	for _, otherRec := range m.store {
		if otherRec.app.ID != id && otherRec.app.OrgID == rec.app.OrgID && otherRec.app.Name == name {
			return nil, errors.New("application name already exists in organization")
		}
	}

	rec.app.Name = name
	rec.app.Description = description
	rec.app.UpdatedAt = time.Now()
	m.store[id] = rec

	return &rec.app, nil
}

func (m *MemoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.store[id]; !ok {
		return errors.New("application not found")
	}

	delete(m.store, id)
	return nil
}
