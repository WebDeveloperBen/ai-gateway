package kv

import (
	"context"
	"sync"
	"time"
)

type memoryItem struct {
	value   string
	expires time.Time
}

type MemoryStore struct {
	mu    sync.RWMutex
	store map[string]memoryItem
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{store: make(map[string]memoryItem)}
}

func (m *MemoryStore) Get(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	item, ok := m.store[key]
	m.mu.RUnlock()
	if !ok || (item.expires.After(time.Time{}) && time.Now().After(item.expires)) {
		return "", nil
	}
	return item.value, nil
}

func (m *MemoryStore) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	expires := time.Time{}
	if ttl > 0 {
		expires = time.Now().Add(ttl)
	}
	m.store[key] = memoryItem{value: value, expires: expires}
	return nil
}

func (m *MemoryStore) Del(ctx context.Context, key string) error {
	m.mu.Lock()
	delete(m.store, key)
	m.mu.Unlock()
	return nil
}

func (m *MemoryStore) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	item, ok := m.store[key]
	m.mu.RUnlock()
	return ok && (item.expires.IsZero() || time.Now().Before(item.expires)), nil
}

func (m *MemoryStore) Keys(ctx context.Context, pattern string) ([]string, error) {
	// only supports simple suffix * as wildcard
	m.mu.RLock()
	defer m.mu.RUnlock()
	var keys []string
	for k := range m.store {
		if matchPattern(k, pattern) {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

func matchPattern(key, pattern string) bool {
	if pattern == "*" {
		return true
	}
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(key) >= len(prefix) && key[:len(prefix)] == prefix
	}
	return key == pattern
}
