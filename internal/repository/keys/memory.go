package keys

import (
	"context"
	"sync"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

type keyRecord struct {
	key model.Key
	phc string
}

type MemoryStore struct {
	mu    sync.RWMutex
	store map[string]keyRecord
}

var _ KeyRepository = (*store)(nil)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{store: make(map[string]keyRecord)}
}

func (m *MemoryStore) Insert(ctx context.Context, k model.Key, phc string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[k.KeyID] = keyRecord{key: k, phc: phc}
	return nil
}

func (m *MemoryStore) GetByKeyID(ctx context.Context, keyID string) (*model.Key, error) {
	m.mu.RLock()
	rec, ok := m.store[keyID]
	m.mu.RUnlock()
	if !ok {
		return nil, nil
	}
	return &rec.key, nil
}

func (m *MemoryStore) GetPHCByKeyID(ctx context.Context, keyID string) (string, error) {
	m.mu.RLock()
	rec, ok := m.store[keyID]
	m.mu.RUnlock()
	if !ok {
		return "", nil
	}
	return rec.phc, nil
}

func (m *MemoryStore) UpdateStatus(ctx context.Context, keyID string, status model.KeyStatus) error {
	m.mu.Lock()
	rec, ok := m.store[keyID]
	if ok {
		rec.key.Status = status
		m.store[keyID] = rec
	}
	m.mu.Unlock()
	return nil
}

func (m *MemoryStore) TouchLastUsed(ctx context.Context, keyID string) error {
	m.mu.Lock()
	rec, ok := m.store[keyID]
	if ok {
		now := time.Now()
		rec.key.LastUsedAt = &now
		m.store[keyID] = rec
	}
	m.mu.Unlock()
	return nil
}
