package loadbalancing_test

import (
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/loadbalancing"
	"github.com/stretchr/testify/require"
)

func TestRandomSelector_Select(t *testing.T) {
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

	t.Run("multiple instances returns valid instance", func(t *testing.T) {
		instances := []string{"instance1", "instance2", "instance3"}

		// Run multiple times to ensure we always get a valid instance
		for range 20 {
			result := selector.Select(instances, "key1")
			require.Contains(t, instances, result, "result should be one of the instances")
		}
	})

	t.Run("key parameter is ignored", func(t *testing.T) {
		instances := []string{"instance1", "instance2"}

		// The key parameter should not affect random selection
		// We can't test randomness deterministically, but we can ensure
		// that different keys don't break the functionality
		for _, key := range []string{"key1", "key2", "different_key", ""} {
			result := selector.Select(instances, key)
			require.Contains(t, instances, result)
		}
	})

	t.Run("distribution is reasonably random", func(t *testing.T) {
		instances := []string{"a", "b", "c"}
		counts := make(map[string]int)

		// Run many selections to check distribution
		const iterations = 1000
		for range iterations {
			result := selector.Select(instances, "test")
			counts[result]++
		}

		// Each instance should be selected at least some times
		// With truly random distribution, each should be around iterations/len(instances) = 333
		expectedMin := iterations / len(instances) / 3 // Allow for some variance
		for _, instance := range instances {
			require.Greater(t, counts[instance], expectedMin,
				"instance %s should be selected more than %d times in %d iterations",
				instance, expectedMin, iterations)
		}
	})

	t.Run("works with different instance counts", func(t *testing.T) {
		testCases := [][]string{
			{"single"},
			{"first", "second"},
			{"a", "b", "c", "d", "e"},
		}

		for _, instances := range testCases {
			result := selector.Select(instances, "test")
			require.Contains(t, instances, result)
		}
	})
}
