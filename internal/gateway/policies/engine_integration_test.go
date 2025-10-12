package policies

import (
	"context"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

// TestPolicyEnforcementIntegration tests the complete policy enforcement flow
// from engine to individual policies
func TestPolicyEnforcementIntegration(t *testing.T) {
	ctx := context.Background()
	cache := kv.NewMemoryStore()
	engine := NewEngine(nil, cache)

	// Create test policies
	rateLimitPolicy := NewRateLimitPolicy(model.RateLimitConfig{
		RequestsPerMinute: 10,
	}, cache)

	tokenLimitPolicy := NewTokenLimitPolicy(model.TokenLimitConfig{
		MaxPromptTokens: 1000,
	})

	modelAllowlistPolicy := NewModelAllowlistPolicy(model.ModelAllowlistConfig{
		AllowedModelIDs: []string{"gpt-4", "gpt-3.5-turbo"},
	})

	policies := []Policy{rateLimitPolicy, tokenLimitPolicy, modelAllowlistPolicy}

	t.Run("AllPoliciesPass", func(t *testing.T) {
		preCtx := &PreRequestContext{
			Request:          nil, // Not needed for basic tests
			AppID:            "test-app",
			Model:            "gpt-4",
			EstimatedTokens:  500,
			RequestSizeBytes: 1024,
		}

		err := engine.CheckPreRequest(ctx, policies, preCtx)
		if err != nil {
			t.Errorf("Expected all policies to pass, got error: %v", err)
		}
	})

	t.Run("RateLimitExceeded", func(t *testing.T) {
		// Make 11 requests (exceeds 10 per minute limit)
		for i := range 11 {
			preCtx := &PreRequestContext{
				Request: nil,
				AppID:   "test-app-rate-limit",
				Model:   "gpt-4",
			}

			err := engine.CheckPreRequest(ctx, []Policy{rateLimitPolicy}, preCtx)
			if i < 10 && err != nil {
				t.Errorf("Expected request %d to pass, got error: %v", i+1, err)
			}
			if i == 10 && err == nil {
				t.Errorf("Expected request 11 to be rate limited")
			}
		}
	})

	t.Run("TokenLimitExceeded", func(t *testing.T) {
		preCtx := &PreRequestContext{
			Request:         nil,
			AppID:           "test-app",
			Model:           "gpt-4",
			EstimatedTokens: 1500, // Exceeds 1000 limit
		}

		err := engine.CheckPreRequest(ctx, []Policy{tokenLimitPolicy}, preCtx)
		if err == nil {
			t.Error("Expected token limit to be exceeded")
		}
	})

	t.Run("ModelNotAllowed", func(t *testing.T) {
		preCtx := &PreRequestContext{
			Request: nil,
			AppID:   "test-app",
			Model:   "gpt-5", // Not in allowlist
		}

		err := engine.CheckPreRequest(ctx, []Policy{modelAllowlistPolicy}, preCtx)
		if err == nil {
			t.Error("Expected model to be blocked")
		}
	})

	t.Run("RequestSizeExceeded", func(t *testing.T) {
		sizePolicy := NewRequestSizePolicy(model.RequestSizeConfig{
			MaxRequestBytes: 1000,
		})

		preCtx := &PreRequestContext{
			Request:          nil,
			AppID:            "test-app",
			RequestSizeBytes: 2000, // Exceeds 1000 limit
		}

		err := engine.CheckPreRequest(ctx, []Policy{sizePolicy}, preCtx)
		if err == nil {
			t.Error("Expected request size to be exceeded")
		}
	})
}

// TestRegistryIntegration tests that the registry works end-to-end
func TestRegistryIntegration(t *testing.T) {
	// Test that all built-in policies are registered
	expectedTypes := []model.PolicyType{
		model.PolicyTypeRateLimit,
		model.PolicyTypeTokenLimit,
		model.PolicyTypeModelAllowlist,
		model.PolicyTypeRequestSize,
	}

	for _, policyType := range expectedTypes {
		t.Run(string(policyType), func(t *testing.T) {
			factory, exists := GetFactory(policyType)
			if !exists {
				t.Errorf("Policy type %s not registered", policyType)
			}

			if factory == nil {
				t.Errorf("Factory for %s is nil", policyType)
			}

			// Test that factory can create a policy
			cache := kv.NewMemoryStore()
			deps := PolicyDependencies{Cache: cache}

			config := getTestConfig(policyType)
			policy, err := factory(config, deps)
			if err != nil {
				t.Errorf("Factory failed to create policy: %v", err)
			}

			if policy == nil {
				t.Errorf("Factory returned nil policy")
			}

			if policy.Type() != policyType {
				t.Errorf("Policy type mismatch: expected %s, got %s", policyType, policy.Type())
			}
		})
	}
}

// TestCachingIntegration tests the three-tier caching system
func TestCachingIntegration(t *testing.T) {
	ctx := context.Background()
	cache := kv.NewMemoryStore()
	engine := NewEngine(nil, cache)

	appID := "550e8400-e29b-41d4-a716-446655440000" // Valid UUID for testing

	// Create test policies to cache
	testPolicies := []CachedPolicy{
		{
			Type:   model.PolicyTypeRateLimit,
			Config: []byte(`{"requests_per_minute":1000}`),
		},
		{
			Type:   model.PolicyTypeTokenLimit,
			Config: []byte(`{"max_prompt_tokens":4000}`),
		},
	}

	t.Run("CacheMiss_LoadFromCache", func(t *testing.T) {
		// Pre-populate cache (simulating previous DB load)
		err := SetCachedPoliciesRaw(ctx, cache, appID, testPolicies)
		if err != nil {
			t.Fatalf("Failed to set up cache: %v", err)
		}

		// First load should hit cache
		policies, err := engine.LoadPolicies(ctx, appID)
		if err != nil {
			t.Fatalf("Failed to load policies: %v", err)
		}

		if len(policies) != 2 {
			t.Errorf("Expected 2 policies, got %d", len(policies))
		}

		// Verify policy types
		expectedTypes := map[model.PolicyType]bool{
			model.PolicyTypeRateLimit:  false,
			model.PolicyTypeTokenLimit: false,
		}

		for _, policy := range policies {
			if _, exists := expectedTypes[policy.Type()]; exists {
				expectedTypes[policy.Type()] = true
			}
		}

		for policyType, found := range expectedTypes {
			if !found {
				t.Errorf("Expected policy type %s not found", policyType)
			}
		}
	})

	t.Run("CacheHit_MemoryCache", func(t *testing.T) {
		// Second load should hit memory cache
		policies, err := engine.LoadPolicies(ctx, appID)
		if err != nil {
			t.Fatalf("Failed to load policies from cache: %v", err)
		}

		if len(policies) != 2 {
			t.Errorf("Expected 2 policies from cache, got %d", len(policies))
		}
	})

	t.Run("CacheInvalidation", func(t *testing.T) {
		// Invalidate cache
		err := engine.InvalidateCache(ctx, appID)
		if err != nil {
			t.Fatalf("Failed to invalidate cache: %v", err)
		}

		// Clear memory cache to force cache miss
		engine.memoryCache.Purge()

		// Next load should miss all caches and try DB (which is nil, so expect error)
		_, err = engine.LoadPolicies(ctx, appID)
		if err == nil {
			t.Error("Expected error after cache invalidation without DB")
		}
	})
}

// TestConcurrentPolicyEnforcement tests thread safety
func TestConcurrentPolicyEnforcement(t *testing.T) {
	ctx := context.Background()
	cache := kv.NewMemoryStore()
	engine := NewEngine(nil, cache)

	// Create a rate limit policy for concurrent testing
	rateLimitPolicy := NewRateLimitPolicy(model.RateLimitConfig{
		RequestsPerMinute: 1000, // High limit to avoid blocking
	}, cache)

	policies := []Policy{rateLimitPolicy}

	t.Run("ConcurrentRequests", func(t *testing.T) {
		const numGoroutines = 10
		const requestsPerGoroutine = 50

		errChan := make(chan error, numGoroutines*requestsPerGoroutine)

		// Launch multiple goroutines making requests
		for i := range numGoroutines {
			go func(goroutineID int) {
				for range requestsPerGoroutine {
					preCtx := &PreRequestContext{
						Request: nil,
						AppID:   "concurrent-app",
						Model:   "gpt-4",
					}

					err := engine.CheckPreRequest(ctx, policies, preCtx)
					errChan <- err
				}
			}(i)
		}

		// Collect all results
		totalRequests := numGoroutines * requestsPerGoroutine
		for i := range totalRequests {
			err := <-errChan
			if err != nil {
				t.Errorf("Request %d failed: %v", i+1, err)
			}
		}
	})
}

// Helper functions

func getTestConfig(policyType model.PolicyType) []byte {
	switch policyType {
	case model.PolicyTypeRateLimit:
		return []byte(`{"requests_per_minute":1000}`)
	case model.PolicyTypeTokenLimit:
		return []byte(`{"max_prompt_tokens":4000}`)
	case model.PolicyTypeModelAllowlist:
		return []byte(`{"allowed_model_ids":["gpt-4"]}`)
	case model.PolicyTypeRequestSize:
		return []byte(`{"max_request_bytes":51200}`)
	default:
		return []byte(`{}`)
	}
}
