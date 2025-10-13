package kv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewKeyspace(t *testing.T) {
	t.Run("adds separator when missing", func(t *testing.T) {
		ks := NewKeyspace("test")
		assert.Equal(t, "test:", ks.Prefix)
	})

	t.Run("does not add separator when present", func(t *testing.T) {
		ks := NewKeyspace("test:")
		assert.Equal(t, "test:", ks.Prefix)
	})
}

func TestKeyspace_Key(t *testing.T) {
	ks := NewKeyspace("test:")

	t.Run("empty parts returns prefix without separator", func(t *testing.T) {
		key := ks.Key()
		assert.Equal(t, "test", key)
	})

	t.Run("single part", func(t *testing.T) {
		key := ks.Key("part1")
		assert.Equal(t, "test:part1", key)
	})

	t.Run("multiple parts", func(t *testing.T) {
		key := ks.Key("part1", "part2", "part3")
		assert.Equal(t, "test:part1:part2:part3", key)
	})
}

func TestKeyspace_PrefixOf(t *testing.T) {
	ks := NewKeyspace("test:")

	t.Run("empty parts returns prefix", func(t *testing.T) {
		prefix := ks.PrefixOf()
		assert.Equal(t, "test:", prefix)
	})

	t.Run("single part", func(t *testing.T) {
		prefix := ks.PrefixOf("part1")
		assert.Equal(t, "test:part1:", prefix)
	})

	t.Run("multiple parts", func(t *testing.T) {
		prefix := ks.PrefixOf("part1", "part2")
		assert.Equal(t, "test:part1:part2:", prefix)
	})
}

func TestKeyspace_PatternAll(t *testing.T) {
	ks := NewKeyspace("test:")
	pattern := ks.PatternAll()
	assert.Equal(t, "test:*", pattern)
}

func TestKeyspace_Pattern(t *testing.T) {
	ks := NewKeyspace("test:")

	t.Run("empty parts returns PatternAll", func(t *testing.T) {
		pattern := ks.Pattern()
		assert.Equal(t, "test:*", pattern)
	})

	t.Run("single part", func(t *testing.T) {
		pattern := ks.Pattern("part1")
		assert.Equal(t, "test:part1:*", pattern)
	})

	t.Run("multiple parts", func(t *testing.T) {
		pattern := ks.Pattern("part1", "part2")
		assert.Equal(t, "test:part1:part2:*", pattern)
	})
}

func TestGlobEscape(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"with*star", "with\\*star"},
		{"with?question", "with\\?question"},
		{"with[bracket", "with\\[bracket"},
		{"with]bracket", "with\\]bracket"},
		{"with\\backslash", "with\\\\backslash"},
		{"mixed*special?chars[here]", "mixed\\*special\\?chars\\[here\\]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := globEscape(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewDriver(t *testing.T) {
	t.Run("memory backend", func(t *testing.T) {
		store, err := NewDriver(Config{Backend: BackendMemory})
		require.NoError(t, err)
		assert.IsType(t, &MemoryStore{}, store)
	})

	t.Run("redis backend", func(t *testing.T) {
		store, err := NewDriver(Config{
			Backend:   BackendRedis,
			RedisAddr: "localhost:6379",
		})
		require.NoError(t, err)
		assert.IsType(t, &RedisStore{}, store)
	})

	t.Run("unsupported backend", func(t *testing.T) {
		store, err := NewDriver(Config{Backend: "unsupported"})
		require.Error(t, err)
		assert.Nil(t, store)
		assert.Contains(t, err.Error(), "unsupported kv backend")
	})
}

func TestKeyModel(t *testing.T) {
	result := KeyModel("tenant1", "model1")
	assert.Equal(t, "modelreg:tenant1:model1", result)
}

func TestTenantPrefix(t *testing.T) {
	result := TenantPrefix("tenant1")
	assert.Equal(t, "modelreg:tenant1:", result)
}

func TestPatternAll(t *testing.T) {
	result := PatternAll()
	assert.Equal(t, "modelreg:*", result)
}

func TestPatternTenantAll(t *testing.T) {
	result := PatternTenantAll("tenant1")
	assert.Equal(t, "modelreg:tenant1:*", result)
}
