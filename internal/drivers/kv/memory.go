package kv

import (
	"context"
	"path"
	"sync"
	"time"
)

var _ = (*RedisStore)(nil)

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

func (m *MemoryStore) Close(ctx context.Context) error {
	m.mu.Lock()
	for k := range m.store {
		delete(m.store, k)
	}
	m.mu.Unlock()
	return nil
}

// ScanAll simulates Redis SCAN MATCH by returning all keys that match the pattern.
// 'count' is ignored (it’s a hint in Redis; memory can just return all).
func (m *MemoryStore) ScanAll(ctx context.Context, pattern string, count int64) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	now := time.Now()

	var keys []string
	for k, item := range m.store {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if !item.expires.IsZero() && now.After(item.expires) {
			continue
		}
		if matchPattern(k, pattern) {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

// ScanGetAll returns key -> value for all keys matching the pattern (skips expired).
// Matches Redis Scan + pipelined GET behavior.
func (m *MemoryStore) ScanGetAll(ctx context.Context, pattern string, count int64) (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	now := time.Now()

	out := make(map[string]string)
	for k, item := range m.store {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if !item.expires.IsZero() && now.After(item.expires) {
			continue
		}
		if matchPattern(k, pattern) {
			out[k] = item.value
		}
	}
	return out, nil
}

// matchPattern implements Redis-like glob matching:
//   - matches any sequence
//     ?  matches any single char
//
// [abc] or [a-c] character classes
// \  escapes the next metacharacter
func matchPattern(key, pattern string) bool {
	ok, err := path.Match(pattern, key)
	if err != nil {
		// Bad pattern (e.g., unterminated class) -> treat as no match
		return false
	}
	return ok
}
