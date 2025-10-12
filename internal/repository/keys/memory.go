package keys

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
)

type keyRecord struct {
	key model.Key
	phc string
}

type MemoryStore struct {
	mu    sync.RWMutex
	store map[string]keyRecord
}

var _ KeyRepository = (*MemoryStore)(nil)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{store: make(map[string]keyRecord)}
}

func (m *MemoryStore) Insert(ctx context.Context, k model.Key, phc string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[k.KeyPrefix] = keyRecord{key: k, phc: phc}
	return nil
}

func (m *MemoryStore) GetByKeyPrefix(ctx context.Context, keyPrefix string) (*model.Key, error) {
	m.mu.RLock()
	rec, ok := m.store[keyPrefix]
	m.mu.RUnlock()
	if !ok {
		return nil, errors.New("key not found")
	}
	return &rec.key, nil
}

func (m *MemoryStore) GetSecretPHCByPrefix(ctx context.Context, keyPrefix string) (string, error) {
	m.mu.RLock()
	rec, ok := m.store[keyPrefix]
	m.mu.RUnlock()
	if !ok {
		return "", errors.New("key not found")
	}
	return rec.phc, nil
}

func (m *MemoryStore) UpdateStatus(ctx context.Context, keyPrefix string, status model.KeyStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	rec, ok := m.store[keyPrefix]
	if !ok {
		return errors.New("key not found")
	}
	rec.key.Status = status
	m.store[keyPrefix] = rec
	return nil
}

func (m *MemoryStore) TouchLastUsed(ctx context.Context, keyPrefix string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	rec, ok := m.store[keyPrefix]
	if !ok {
		return errors.New("key not found")
	}
	now := time.Now()
	rec.key.LastUsedAt = &now
	m.store[keyPrefix] = rec
	return nil
}

func (m *MemoryStore) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for prefix, rec := range m.store {
		if rec.key.ID == id {
			delete(m.store, prefix)
			return nil
		}
	}
	return errors.New("key not found")
}
