package loadbalancing_test

import (
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/loadbalancing"
	"github.com/stretchr/testify/require"
)

func TestRoundRobinSelector(t *testing.T) {
	t.Run("empty instances", func(t *testing.T) {
		selector := loadbalancing.NewRoundRobinSelector()
		result := selector.Select([]string{}, "key1")
		require.Empty(t, result)
	})

	t.Run("single instance", func(t *testing.T) {
		selector := loadbalancing.NewRoundRobinSelector()
		instances := []string{"instance1"}
		result := selector.Select(instances, "key1")
		require.Equal(t, "instance1", result)

		// Should always return the same instance
		result = selector.Select(instances, "key1")
		require.Equal(t, "instance1", result)
	})

	t.Run("multiple instances", func(t *testing.T) {
		selector := loadbalancing.NewRoundRobinSelector()
		instances := []string{"instance1", "instance2", "instance3"}

		// First call for key1 should return instance1
		result1 := selector.Select(instances, "key1")
		require.Equal(t, "instance1", result1)

		// Second call for key1 should return instance2
		result2 := selector.Select(instances, "key1")
		require.Equal(t, "instance2", result2)

		// Third call for key1 should return instance3
		result3 := selector.Select(instances, "key1")
		require.Equal(t, "instance3", result3)

		// Fourth call for key1 should return instance1 again
		result4 := selector.Select(instances, "key1")
		require.Equal(t, "instance1", result4)
	})

	t.Run("different keys are independent", func(t *testing.T) {
		selector := loadbalancing.NewRoundRobinSelector()
		instances := []string{"instance1", "instance2"}

		// key1 gets instance1
		result1 := selector.Select(instances, "key1")
		require.Equal(t, "instance1", result1)

		// key2 gets instance1 (independent counter)
		result2 := selector.Select(instances, "key2")
		require.Equal(t, "instance1", result2)

		// key1 gets instance2
		result3 := selector.Select(instances, "key1")
		require.Equal(t, "instance2", result3)

		// key2 gets instance2
		result4 := selector.Select(instances, "key2")
		require.Equal(t, "instance2", result4)
	})
}

func TestRandomSelector(t *testing.T) {
	selector := &loadbalancing.RandomSelector{}

	t.Run("empty instances", func(t *testing.T) {
		result := selector.Select([]string{}, "key1")
		require.Empty(t, result)
	})

	t.Run("single instance", func(t *testing.T) {
		instances := []string{"instance1"}
		result := selector.Select(instances, "key1")
		require.Equal(t, "instance1", result)
	})

	t.Run("multiple instances returns one of them", func(t *testing.T) {
		instances := []string{"instance1", "instance2", "instance3"}

		// Run multiple times to ensure we get valid results
		for range 10 {
			result := selector.Select(instances, "key1")
			require.Contains(t, instances, result)
		}
	})

	t.Run("key parameter is ignored", func(t *testing.T) {
		instances := []string{"instance1", "instance2"}

		// The key parameter should not affect random selection
		// We can't test randomness deterministically, but we can ensure
		// that different keys don't break the functionality
		for _, key := range []string{"key1", "key2", "different_key"} {
			result := selector.Select(instances, key)
			require.Contains(t, instances, result)
		}
	})
}
