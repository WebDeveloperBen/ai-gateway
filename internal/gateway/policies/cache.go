package policies

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

const (
	// Cache TTLs
	PolicyCacheTTL = 5 * time.Minute

	// Cache key prefixes
	PolicyCachePrefix = "policy:app:"
)

// CacheKey generates a Redis key for caching policies
func CacheKey(appID string) string {
	return fmt.Sprintf("%s%s:policies", PolicyCachePrefix, appID)
}

// GetCachedPolicies retrieves policies from cache
// Note: Returns cached policy data that must be reconstructed with Engine.NewPolicy
func GetCachedPolicies(ctx context.Context, cache kv.KvStore, appID string) ([]CachedPolicy, bool, error) {
	key := CacheKey(appID)

	data, err := cache.Get(ctx, key)
	if err != nil {
		// Cache miss or error
		return nil, false, err
	}

	if data == "" {
		return nil, false, nil
	}

	// Deserialize cached policies
	var cachedData []CachedPolicy
	if err := json.Unmarshal([]byte(data), &cachedData); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal cached policies: %w", err)
	}

	return cachedData, true, nil
}

// SetCachedPoliciesRaw stores policies in cache from their raw DB representation
func SetCachedPoliciesRaw(ctx context.Context, cache kv.KvStore, appID string, policies []CachedPolicy) error {
	data, err := json.Marshal(policies)
	if err != nil {
		return fmt.Errorf("failed to marshal policies: %w", err)
	}

	key := CacheKey(appID)
	return cache.Set(ctx, key, string(data), PolicyCacheTTL)
}

// InvalidatePolicyCache removes cached policies for an app
func InvalidatePolicyCache(ctx context.Context, cache kv.KvStore, appID string) error {
	key := CacheKey(appID)
	return cache.Del(ctx, key)
}

// CachedPolicy represents a policy in cache-friendly format
type CachedPolicy struct {
	Type   model.PolicyType `json:"type"`
	Config []byte           `json:"config"`
}
