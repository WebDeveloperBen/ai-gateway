package policies

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
	"github.com/WebDeveloperBen/ai-gateway/internal/logger"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/observability"
	"github.com/google/uuid"
	lru "github.com/hashicorp/golang-lru/v2"
)

// policyCacheEntry holds compiled policies with expiration time
type policyCacheEntry struct {
	policies  []Policy
	expiresAt time.Time
}

// Engine is the main policy engine that loads and executes policies
type Engine struct {
	db          *db.Queries
	cache       kv.KvStore
	memoryCache *lru.Cache[string, *policyCacheEntry]
	cacheMu     sync.RWMutex
	cacheTTL    time.Duration
}

// NewEngine creates a new policy engine
func NewEngine(queries *db.Queries, cache kv.KvStore) *Engine {
	// Create in-memory LRU cache with 1000 entries (supports ~1000 concurrent apps)
	memCache, _ := lru.New[string, *policyCacheEntry](1000)

	return &Engine{
		db:          queries,
		cache:       cache,
		memoryCache: memCache,
		cacheTTL:    30 * time.Second, // Short TTL to keep policies fresh
	}
}

// LoadPolicies loads all enabled policies for an application
// Uses three-tier cache: memory (30s TTL) -> Redis (5m TTL) -> DB
func (e *Engine) LoadPolicies(ctx context.Context, appID string) ([]Policy, error) {
	appUUID, err := uuid.Parse(appID)
	if err != nil {
		return nil, fmt.Errorf("invalid app ID: %w", err)
	}

	// Tier 1: Check in-memory cache first (fastest - no network RTT)
	e.cacheMu.RLock()
	if entry, found := e.memoryCache.Get(appID); found {
		// Check if entry is still valid
		if time.Now().Before(entry.expiresAt) {
			policies := entry.policies
			e.cacheMu.RUnlock()
			observability.FromContext(ctx).RecordPolicyCacheHit(ctx, "memory")
			return policies, nil
		}
		// Entry expired, will reload
	}
	e.cacheMu.RUnlock()
	observability.FromContext(ctx).RecordPolicyCacheMiss(ctx, "memory")

	// Tier 2: Check Redis cache (medium - network RTT but cached)
	cachedPolicies, found, err := GetCachedPolicies(ctx, e.cache, appID)
	if err == nil && found {
		observability.FromContext(ctx).RecordPolicyCacheHit(ctx, "redis")
		// Reconstruct policies from cached data
		policies := make([]Policy, 0, len(cachedPolicies))
		for _, cached := range cachedPolicies {
			policy, err := e.NewPolicy(cached.Type, cached.Config)
			if err != nil {
				// Log error but continue with other policies
				logger.GetLogger(ctx).Error().
					Err(err).
					Str("app_id", appID).
					Str("policy_type", string(cached.Type)).
					Msg("Failed to reconstruct cached policy")
				continue
			}
			policies = append(policies, policy)
		}

		// Store in memory cache for next request
		e.cacheMu.Lock()
		e.memoryCache.Add(appID, &policyCacheEntry{
			policies:  policies,
			expiresAt: time.Now().Add(e.cacheTTL),
		})
		e.cacheMu.Unlock()

		return policies, nil
	}
	observability.FromContext(ctx).RecordPolicyCacheMiss(ctx, "redis")

	// Tier 3: Load from database (slowest - network + query)
	observability.FromContext(ctx).RecordPolicyCacheMiss(ctx, "db")

	// Check if database is available (for unit testing)
	if e.db == nil {
		return nil, fmt.Errorf("database not available for loading policies")
	}

	dbPolicies, err := e.db.ListEnabledPolicies(ctx, db.ListEnabledPoliciesParams{
		AppID:  appUUID,
		Limit:  1000, // Load all enabled policies for the app
		Offset: 0,
	})
	if err != nil {
		observability.FromContext(ctx).RecordPolicyLoadError(ctx, appID, err)
		return nil, fmt.Errorf("failed to load policies: %w", err)
	}

	// Convert DB policies to Policy interfaces
	// Store original configs for caching
	policies := make([]Policy, 0, len(dbPolicies))
	policiesToCache := make([]CachedPolicy, 0, len(dbPolicies))

	for _, dbPolicy := range dbPolicies {
		policy, err := e.NewPolicy(model.PolicyType(dbPolicy.PolicyType), dbPolicy.Config)
		if err != nil {
			// Log error but continue with other policies
			logger.GetLogger(ctx).Error().
				Err(err).
				Str("app_id", appID).
				Str("policy_type", string(dbPolicy.PolicyType)).
				Msg("Failed to create policy from database")
			continue
		}
		policies = append(policies, policy)

		// Store for cache
		policiesToCache = append(policiesToCache, CachedPolicy{
			Type:   model.PolicyType(dbPolicy.PolicyType),
			Config: dbPolicy.Config,
		})
	}

	// Cache in Redis (ignore errors - cache failures shouldn't break requests)
	_ = SetCachedPoliciesRaw(ctx, e.cache, appID, policiesToCache)

	// Cache in memory
	e.cacheMu.Lock()
	e.memoryCache.Add(appID, &policyCacheEntry{
		policies:  policies,
		expiresAt: time.Now().Add(e.cacheTTL),
	})
	e.cacheMu.Unlock()

	return policies, nil
}

// CheckPreRequest runs all pre-request checks for the given policies
// Returns error if any policy check fails
func (e *Engine) CheckPreRequest(ctx context.Context, policies []Policy, req *PreRequestContext) error {
	for _, policy := range policies {
		start := time.Now()
		err := policy.PreCheck(ctx, req)
		duration := time.Since(start)

		obs := observability.FromContext(ctx)
		obs.RecordPolicyCheck(ctx, string(policy.Type()), duration, err != nil)

		if err != nil {
			return fmt.Errorf("policy %s failed: %w", policy.Type(), err)
		}
	}
	return nil
}

// RecordPostRequest runs all post-request checks asynchronously
// This should be called in a goroutine with a detached context
func (e *Engine) RecordPostRequest(ctx context.Context, policies []Policy, req *PostRequestContext) {
	for _, policy := range policies {
		policy.PostCheck(ctx, req)
	}
}

// NewPolicy creates a new policy instance based on the policy type and config
// Uses the policy registry to look up the appropriate factory function
func (e *Engine) NewPolicy(policyType model.PolicyType, config []byte) (Policy, error) {
	// Try registry first (built-in policies registered via init())
	factory, exists := GetFactory(policyType)
	if exists {
		deps := PolicyDependencies{Cache: e.cache}
		return factory(config, deps)
	}

	// Fallback to CEL policy for custom policies
	if policyType == model.PolicyTypeCustomCEL {
		return NewCELPolicy(policyType, config)
	}

	// Unknown policy type
	return nil, fmt.Errorf("unknown policy type: %s", policyType)
}

// InvalidateCache removes an application's policies from all cache tiers
// Call this when policies are updated/created/deleted for an application
func (e *Engine) InvalidateCache(ctx context.Context, appID string) error {
	// Remove from memory cache
	e.cacheMu.Lock()
	e.memoryCache.Remove(appID)
	e.cacheMu.Unlock()

	// Remove from Redis cache
	return InvalidatePolicyCache(ctx, e.cache, appID)
}

// InvalidateAllCache clears the entire policy cache (all applications)
// Use sparingly - mainly for testing or system-wide policy changes
func (e *Engine) InvalidateAllCache(ctx context.Context) error {
	// Clear memory cache
	e.cacheMu.Lock()
	e.memoryCache.Purge()
	e.cacheMu.Unlock()

	// Clear Redis cache (policies are stored with prefix "cache:policies:")
	// We don't have a bulk delete in KvStore interface, so just log
	logger.GetLogger(ctx).Info().Msg("Cleared in-memory policy cache")
	return nil
}
