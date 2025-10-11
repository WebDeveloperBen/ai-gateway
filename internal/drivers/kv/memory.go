package kv

import (
	"context"
	"fmt"
	"path"
	"strconv"
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

func (m *MemoryStore) Incr(ctx context.Context, key string) (int64, error) {
	return m.IncrBy(ctx, key, 1)
}

func (m *MemoryStore) IncrBy(ctx context.Context, key string, amount int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.store[key]
	if !ok || (item.expires.After(time.Time{}) && time.Now().After(item.expires)) {
		// Key doesn't exist or expired, set to 0 and increment
		m.store[key] = memoryItem{value: strconv.FormatInt(amount, 10), expires: time.Time{}}
		return amount, nil
	}

	// Parse existing value as integer
	val, err := strconv.ParseInt(item.value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("value is not an integer")
	}

	newVal := val + amount
	item.value = strconv.FormatInt(newVal, 10)
	m.store[key] = item

	return newVal, nil
}

func (m *MemoryStore) Expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.store[key]
	if !ok || (item.expires.After(time.Time{}) && time.Now().After(item.expires)) {
		return false, nil // Key doesn't exist
	}

	if ttl > 0 {
		item.expires = time.Now().Add(ttl)
	} else {
		item.expires = time.Time{} // No expiration
	}
	m.store[key] = item

	return true, nil
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
