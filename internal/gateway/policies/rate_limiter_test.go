package policies_test

import (
	"context"
	"testing"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/policies"
	"github.com/stretchr/testify/require"
)

// Mock KV store for testing
type mockKVStore struct {
	data map[string]int64
}

func newMockKVStore() *mockKVStore {
	return &mockKVStore{data: make(map[string]int64)}
}

func (m *mockKVStore) Get(ctx context.Context, key string) (string, error) {
	// Not used in rate limiter
	return "", nil
}

func (m *mockKVStore) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	// Not used in rate limiter
	return nil
}

func (m *mockKVStore) Del(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

func (m *mockKVStore) Exists(ctx context.Context, key string) (bool, error) {
	_, exists := m.data[key]
	return exists, nil
}

func (m *mockKVStore) Incr(ctx context.Context, key string) (int64, error) {
	m.data[key]++
	return m.data[key], nil
}

func (m *mockKVStore) IncrBy(ctx context.Context, key string, amount int64) (int64, error) {
	m.data[key] += amount
	return m.data[key], nil
}

func (m *mockKVStore) Expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	// Mock always succeeds
	return true, nil
}

func (m *mockKVStore) ScanGetAll(ctx context.Context, pattern string, count int64) (map[string]string, error) {
	// Not used in rate limiter
	return nil, nil
}

func (m *mockKVStore) ScanAll(ctx context.Context, pattern string, count int64) ([]string, error) {
	// Not used in rate limiter
	return nil, nil
}

func (m *mockKVStore) Close(ctx context.Context) error {
	return nil
}

func TestNewRateLimiter(t *testing.T) {
	store := newMockKVStore()
	limiter := policies.NewRateLimiter(store)

	require.NotNil(t, limiter)
	// Can't test internal cache field directly since it's unexported
}

func TestRateLimiter_CheckAndIncrement_NoLimit(t *testing.T) {
	store := newMockKVStore()
	limiter := policies.NewRateLimiter(store)

	allowed, err := limiter.CheckAndIncrement(context.Background(), "test-key", 0, time.Minute)
	require.NoError(t, err)
	require.True(t, allowed) // No limit means always allowed
}

func TestRateLimiter_CheckAndIncrement_UnderLimit(t *testing.T) {
	store := newMockKVStore()
	limiter := policies.NewRateLimiter(store)

	// First request - should be allowed
	allowed, err := limiter.CheckAndIncrement(context.Background(), "test-key", 5, time.Minute)
	require.NoError(t, err)
	require.True(t, allowed)

	// Second request - should still be allowed
	allowed, err = limiter.CheckAndIncrement(context.Background(), "test-key", 5, time.Minute)
	require.NoError(t, err)
	require.True(t, allowed)

	// Check counter
	count, err := limiter.GetCount(context.Background(), "test-key")
	require.NoError(t, err)
	require.Equal(t, 2, count)
}

func TestRateLimiter_CheckAndIncrement_AtLimit(t *testing.T) {
	store := newMockKVStore()
	limiter := policies.NewRateLimiter(store)

	// Fill up to the limit
	for i := 0; i < 3; i++ {
		allowed, err := limiter.CheckAndIncrement(context.Background(), "test-key", 3, time.Minute)
		require.NoError(t, err)
		require.True(t, allowed)
	}

	// Next request should be denied
	allowed, err := limiter.CheckAndIncrement(context.Background(), "test-key", 3, time.Minute)
	require.NoError(t, err)
	require.False(t, allowed)

	// Check counter
	count, err := limiter.GetCount(context.Background(), "test-key")
	require.NoError(t, err)
	require.Equal(t, 4, count) // Incremented even though denied
}

func TestRateLimiter_CheckAndIncrement_OverLimit(t *testing.T) {
	store := newMockKVStore()
	limiter := policies.NewRateLimiter(store)

	// Exceed the limit in one go
	allowed, err := limiter.CheckAndIncrement(context.Background(), "test-key", 1, time.Minute)
	require.NoError(t, err)
	require.True(t, allowed) // First request allowed

	allowed, err = limiter.CheckAndIncrement(context.Background(), "test-key", 1, time.Minute)
	require.NoError(t, err)
	require.False(t, allowed) // Second request denied
}

func TestRateLimiter_Increment(t *testing.T) {
	store := newMockKVStore()
	limiter := policies.NewRateLimiter(store)

	// Increment by 5
	err := limiter.Increment(context.Background(), "test-key", 5, time.Minute)
	require.NoError(t, err)

	// Check count
	count, err := limiter.GetCount(context.Background(), "test-key")
	require.NoError(t, err)
	require.Equal(t, 5, count)

	// Increment by 3 more
	err = limiter.Increment(context.Background(), "test-key", 3, time.Minute)
	require.NoError(t, err)

	count, err = limiter.GetCount(context.Background(), "test-key")
	require.NoError(t, err)
	require.Equal(t, 8, count)
}

func TestRateLimiter_GetCount(t *testing.T) {
	store := newMockKVStore()
	limiter := policies.NewRateLimiter(store)

	// Initially zero
	count, err := limiter.GetCount(context.Background(), "nonexistent-key")
	require.NoError(t, err)
	require.Equal(t, 0, count)

	// After increment
	err = limiter.Increment(context.Background(), "test-key", 7, time.Minute)
	require.NoError(t, err)

	count, err = limiter.GetCount(context.Background(), "test-key")
	require.NoError(t, err)
	require.Equal(t, 7, count)
}

func TestRateLimitKey(t *testing.T) {
	key1 := policies.RateLimitKey("app123", "requests")
	key2 := policies.RateLimitKey("app123", "requests")

	// Keys should be the same for the same minute
	require.Equal(t, key1, key2)

	// Key should contain the expected components
	require.Contains(t, key1, "ratelimit:")
	require.Contains(t, key1, "app123")
	require.Contains(t, key1, "requests")

	// Different apps should have different keys
	key3 := policies.RateLimitKey("app456", "requests")
	require.NotEqual(t, key1, key3)
	require.Contains(t, key3, "app456")

	// Different metrics should have different keys
	key4 := policies.RateLimitKey("app123", "tokens")
	require.NotEqual(t, key1, key4)
	require.Contains(t, key4, "tokens")
}

func TestRateLimiter_DifferentKeys(t *testing.T) {
	store := newMockKVStore()
	limiter := policies.NewRateLimiter(store)

	// Different keys should be independent
	allowed1, err := limiter.CheckAndIncrement(context.Background(), "key1", 2, time.Minute)
	require.NoError(t, err)
	require.True(t, allowed1)

	allowed2, err := limiter.CheckAndIncrement(context.Background(), "key2", 2, time.Minute)
	require.NoError(t, err)
	require.True(t, allowed2)

	// Both should still be allowed
	allowed1, err = limiter.CheckAndIncrement(context.Background(), "key1", 2, time.Minute)
	require.NoError(t, err)
	require.True(t, allowed1)

	allowed2, err = limiter.CheckAndIncrement(context.Background(), "key2", 2, time.Minute)
	require.NoError(t, err)
	require.True(t, allowed2)
}

// Test that would require a failing KV store - but our mock always succeeds
// In a real scenario, we'd want to test Redis failures, but that's integration testing
