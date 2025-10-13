package kv

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStore_Get(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	t.Run("get existing key", func(t *testing.T) {
		err := store.Set(ctx, "test-key", "test-value", 0)
		require.NoError(t, err)

		value, err := store.Get(ctx, "test-key")
		require.NoError(t, err)
		assert.Equal(t, "test-value", value)
	})

	t.Run("get non-existing key", func(t *testing.T) {
		value, err := store.Get(ctx, "non-existing")
		require.NoError(t, err)
		assert.Equal(t, "", value)
	})

	t.Run("get expired key", func(t *testing.T) {
		err := store.Set(ctx, "expired-key", "value", time.Millisecond)
		require.NoError(t, err)

		time.Sleep(10 * time.Millisecond)

		value, err := store.Get(ctx, "expired-key")
		require.NoError(t, err)
		assert.Equal(t, "", value)
	})
}

func TestMemoryStore_Set(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	t.Run("set without TTL", func(t *testing.T) {
		err := store.Set(ctx, "key1", "value1", 0)
		require.NoError(t, err)

		value, err := store.Get(ctx, "key1")
		require.NoError(t, err)
		assert.Equal(t, "value1", value)
	})

	t.Run("set with TTL", func(t *testing.T) {
		err := store.Set(ctx, "key2", "value2", 100*time.Millisecond)
		require.NoError(t, err)

		value, err := store.Get(ctx, "key2")
		require.NoError(t, err)
		assert.Equal(t, "value2", value)

		time.Sleep(150 * time.Millisecond)

		value, err = store.Get(ctx, "key2")
		require.NoError(t, err)
		assert.Equal(t, "", value)
	})
}

func TestMemoryStore_Del(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	err := store.Set(ctx, "delete-key", "value", 0)
	require.NoError(t, err)

	err = store.Del(ctx, "delete-key")
	require.NoError(t, err)

	value, err := store.Get(ctx, "delete-key")
	require.NoError(t, err)
	assert.Equal(t, "", value)
}

func TestMemoryStore_Exists(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	t.Run("existing key", func(t *testing.T) {
		err := store.Set(ctx, "exists-key", "value", 0)
		require.NoError(t, err)

		exists, err := store.Exists(ctx, "exists-key")
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("non-existing key", func(t *testing.T) {
		exists, err := store.Exists(ctx, "non-existing")
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("expired key", func(t *testing.T) {
		err := store.Set(ctx, "expired-exists", "value", time.Millisecond)
		require.NoError(t, err)

		time.Sleep(10 * time.Millisecond)

		exists, err := store.Exists(ctx, "expired-exists")
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestMemoryStore_Incr_IncrBy(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	t.Run("incr on non-existing key", func(t *testing.T) {
		result, err := store.Incr(ctx, "counter")
		require.NoError(t, err)
		assert.Equal(t, int64(1), result)
	})

	t.Run("incr on existing key", func(t *testing.T) {
		result, err := store.Incr(ctx, "counter")
		require.NoError(t, err)
		assert.Equal(t, int64(2), result)
	})

	t.Run("incrBy", func(t *testing.T) {
		result, err := store.IncrBy(ctx, "counter", 5)
		require.NoError(t, err)
		assert.Equal(t, int64(7), result)
	})

	t.Run("incrBy on non-existing key", func(t *testing.T) {
		result, err := store.IncrBy(ctx, "new-counter", 10)
		require.NoError(t, err)
		assert.Equal(t, int64(10), result)
	})

	t.Run("incrBy with negative amount", func(t *testing.T) {
		result, err := store.IncrBy(ctx, "counter", -3)
		require.NoError(t, err)
		assert.Equal(t, int64(4), result)
	})

	t.Run("incr on non-integer value", func(t *testing.T) {
		err := store.Set(ctx, "bad-counter", "not-a-number", 0)
		require.NoError(t, err)

		_, err = store.Incr(ctx, "bad-counter")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "value is not an integer")
	})
}

func TestMemoryStore_Expire(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	t.Run("expire existing key", func(t *testing.T) {
		err := store.Set(ctx, "expire-key", "value", 0)
		require.NoError(t, err)

		success, err := store.Expire(ctx, "expire-key", 50*time.Millisecond)
		require.NoError(t, err)
		assert.True(t, success)

		time.Sleep(60 * time.Millisecond)

		value, err := store.Get(ctx, "expire-key")
		require.NoError(t, err)
		assert.Equal(t, "", value)
	})

	t.Run("expire non-existing key", func(t *testing.T) {
		success, err := store.Expire(ctx, "non-existing", time.Minute)
		require.NoError(t, err)
		assert.False(t, success)
	})

	t.Run("expire with zero TTL", func(t *testing.T) {
		err := store.Set(ctx, "no-expire-key", "value", time.Minute)
		require.NoError(t, err)

		success, err := store.Expire(ctx, "no-expire-key", 0)
		require.NoError(t, err)
		assert.True(t, success)

		// Should still exist
		value, err := store.Get(ctx, "no-expire-key")
		require.NoError(t, err)
		assert.Equal(t, "value", value)
	})
}

func TestMemoryStore_ScanAll_ScanGetAll(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	// Set up test data
	testData := map[string]string{
		"app:user:1":    "alice",
		"app:user:2":    "bob",
		"app:config:db": "postgres",
		"cache:item:1":  "data1",
		"cache:item:2":  "data2",
		"other:key":     "value",
	}

	for k, v := range testData {
		err := store.Set(ctx, k, v, 0)
		require.NoError(t, err)
	}

	t.Run("ScanAll with pattern", func(t *testing.T) {
		keys, err := store.ScanAll(ctx, "app:*", 100)
		require.NoError(t, err)

		expectedKeys := []string{"app:user:1", "app:user:2", "app:config:db"}
		assert.ElementsMatch(t, expectedKeys, keys)
	})

	t.Run("ScanAll with specific pattern", func(t *testing.T) {
		keys, err := store.ScanAll(ctx, "cache:item:*", 100)
		require.NoError(t, err)

		expectedKeys := []string{"cache:item:1", "cache:item:2"}
		assert.ElementsMatch(t, expectedKeys, keys)
	})

	t.Run("ScanGetAll", func(t *testing.T) {
		result, err := store.ScanGetAll(ctx, "app:user:*", 100)
		require.NoError(t, err)

		expected := map[string]string{
			"app:user:1": "alice",
			"app:user:2": "bob",
		}
		assert.Equal(t, expected, result)
	})

	t.Run("ScanAll with expired items", func(t *testing.T) {
		err := store.Set(ctx, "temp:key", "temp-value", time.Millisecond)
		require.NoError(t, err)

		time.Sleep(10 * time.Millisecond)

		keys, err := store.ScanAll(ctx, "temp:*", 100)
		require.NoError(t, err)
		assert.Empty(t, keys)
	})
}

func TestMemoryStore_Close(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	// Add some data
	err := store.Set(ctx, "key1", "value1", 0)
	require.NoError(t, err)
	err = store.Set(ctx, "key2", "value2", 0)
	require.NoError(t, err)

	// Close should clear all data
	err = store.Close(ctx)
	require.NoError(t, err)

	// Verify data is cleared
	value, err := store.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, "", value)

	value, err = store.Get(ctx, "key2")
	require.NoError(t, err)
	assert.Equal(t, "", value)
}

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		key      string
		pattern  string
		expected bool
	}{
		{"test", "test", true},
		{"test", "tes?", true},
		{"test", "te*", true},
		{"test", "t[a-z]st", true},
		{"test", "t[0-9]st", false},
		{"test", "other", false},
		{"app:user:123", "app:user:*", true},
		{"app:config:db", "app:*", true},
		{"cache:item:1", "app:*", false},
	}

	for _, tt := range tests {
		t.Run(tt.key+"_"+tt.pattern, func(t *testing.T) {
			result := matchPattern(tt.key, tt.pattern)
			assert.Equal(t, tt.expected, result)
		})
	}
}
