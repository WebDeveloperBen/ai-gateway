package policies

import (
	"context"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
)

// BenchmarkCacheTiers compares performance across cache tiers
func BenchmarkCacheTiers(b *testing.B) {
	ctx := context.Background()
	cache := kv.NewMemoryStore()
	appID := "bench-app-" + uuid.New().String()

	// Create policies to cache
	testPolicies := []CachedPolicy{
		{
			Type:   model.PolicyTypeRateLimit,
			Config: []byte(`{"requests_per_minute":1000}`),
		},
		{
			Type:   model.PolicyTypeTokenLimit,
			Config: []byte(`{"max_tokens":8192}`),
		},
		{
			Type:   model.PolicyTypeModelAllowlist,
			Config: []byte(`{"allowed_models":["gpt-4","gpt-3.5-turbo"]}`),
		},
	}

	b.Run("MemoryCache", func(b *testing.B) {
		engine := NewEngine(nil, cache)

		// Pre-populate Redis cache
		_ = SetCachedPoliciesRaw(ctx, cache, appID, testPolicies)

		// First call to populate memory cache
		_, _ = engine.LoadPolicies(ctx, appID)

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			_, _ = engine.LoadPolicies(ctx, appID)
		}
	})

	b.Run("RedisCache", func(b *testing.B) {
		// Create fresh engine for each run
		engine := NewEngine(nil, cache)

		// Pre-populate Redis but clear memory cache each time
		_ = SetCachedPoliciesRaw(ctx, cache, appID, testPolicies)

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			// Clear memory cache to force Redis lookup
			engine.memoryCache.Purge()
			_, _ = engine.LoadPolicies(ctx, appID)
		}
	})
}

// BenchmarkPolicyReconstruction measures policy factory overhead
func BenchmarkPolicyReconstruction(b *testing.B) {
	ctx := context.Background()
	cache := kv.NewMemoryStore()
	engine := NewEngine(nil, cache)
	appID := "test-app-" + uuid.New().String()

	// Test with different numbers of policies
	policyCounts := []int{1, 3, 5, 10}

	for _, count := range policyCounts {
		b.Run(b.Name()+"_"+string(rune(count))+"policies", func(b *testing.B) {
			policies := make([]CachedPolicy, count)
			for i := range count {
				policies[i] = CachedPolicy{
					Type:   model.PolicyTypeRateLimit,
					Config: []byte(`{"requests_per_minute":1000}`),
				}
			}
			_ = SetCachedPoliciesRaw(ctx, cache, appID, policies)

			b.ResetTimer()
			b.ReportAllocs()

			for b.Loop() {
				_, _ = engine.LoadPolicies(ctx, appID)
			}
		})
	}
}

// BenchmarkRateLimitCheck measures rate limiter atomic operations
func BenchmarkRateLimitCheck(b *testing.B) {
	ctx := context.Background()
	cache := kv.NewMemoryStore()

	config := model.RateLimitConfig{
		RequestsPerMinute: 1000,
	}

	policy := NewRateLimitPolicy(config, cache)
	preCtx := &PreRequestContext{
		AppID: "test-app-123",
	}

	b.ReportAllocs()

	for b.Loop() {
		_ = policy.PreCheck(ctx, preCtx)
	}
}

// Benchmark: CEL policy compilation vs evaluation
func BenchmarkCELPolicy(b *testing.B) {
	b.Run("Compilation", func(b *testing.B) {
		config := []byte(`{
			"pre_check_expression": "estimated_tokens < 5000 && model.startsWith('gpt-4')"
		}`)

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			_, err := NewCELPolicy(model.PolicyTypeCustomCEL, config)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Evaluation", func(b *testing.B) {
		config := []byte(`{
			"pre_check_expression": "estimated_tokens < 5000 && model.startsWith('gpt-4')"
		}`)

		policy, err := NewCELPolicy(model.PolicyTypeCustomCEL, config)
		if err != nil {
			b.Fatal(err)
		}

		ctx := context.Background()
		preCtx := &PreRequestContext{
			Model:           "gpt-4",
			EstimatedTokens: 1000,
		}

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			if err := policy.PreCheck(ctx, preCtx); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Benchmark: Token limit policy
func BenchmarkTokenLimitPolicy(b *testing.B) {
	config := model.TokenLimitConfig{
		MaxPromptTokens: 4000,
		MaxTotalTokens:  8000,
	}

	policy := NewTokenLimitPolicy(config)
	ctx := context.Background()
	preCtx := &PreRequestContext{
		EstimatedTokens: 1000,
	}

	b.ReportAllocs()

	for b.Loop() {
		if err := policy.PreCheck(ctx, preCtx); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark: Model allowlist policy
func BenchmarkModelAllowlistPolicy(b *testing.B) {
	config := model.ModelAllowlistConfig{
		AllowedModelIDs: []string{"gpt-4", "gpt-4-turbo", "gpt-3.5-turbo"},
	}

	policy := NewModelAllowlistPolicy(config)
	ctx := context.Background()
	preCtx := &PreRequestContext{
		Model: "gpt-4",
	}

	b.ReportAllocs()

	for b.Loop() {
		if err := policy.PreCheck(ctx, preCtx); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark: Request size policy
func BenchmarkRequestSizePolicy(b *testing.B) {
	config := model.RequestSizeConfig{
		MaxRequestBytes: 51200, // 50KB
	}

	policy := NewRequestSizePolicy(config)
	ctx := context.Background()
	preCtx := &PreRequestContext{
		RequestSizeBytes: 10240, // 10KB
	}

	b.ReportAllocs()

	for b.Loop() {
		if err := policy.PreCheck(ctx, preCtx); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark: Full policy chain (all policies except rate limit)
func BenchmarkFullPolicyChain(b *testing.B) {
	ctx := context.Background()
	cache := kv.NewMemoryStore()
	engine := NewEngine(nil, cache)

	policies := []Policy{
		NewTokenLimitPolicy(model.TokenLimitConfig{MaxPromptTokens: 4000, MaxTotalTokens: 8000}),
		NewModelAllowlistPolicy(model.ModelAllowlistConfig{AllowedModelIDs: []string{"gpt-4", "gpt-3.5-turbo"}}),
		NewRequestSizePolicy(model.RequestSizeConfig{MaxRequestBytes: 51200}),
	}

	preCtx := &PreRequestContext{
		AppID:            "bench-app",
		Model:            "gpt-4",
		EstimatedTokens:  1000,
		RequestSizeBytes: 10240,
	}

	b.ReportAllocs()

	for b.Loop() {
		if err := engine.CheckPreRequest(ctx, policies, preCtx); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark: Concurrent rate limit checks (simulates multiple apps)
func BenchmarkConcurrentRateLimits(b *testing.B) {
	ctx := context.Background()
	cache := kv.NewMemoryStore()
	config := model.RateLimitConfig{RequestsPerMinute: 1000}
	policy := NewRateLimitPolicy(config, cache)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			preCtx := &PreRequestContext{
				AppID: "app-" + string(rune(i%10)), // 10 different apps
			}
			_ = policy.PreCheck(ctx, preCtx)
			i++
		}
	})
}

// Benchmark: Policy loading under concurrent load
func BenchmarkConcurrentPolicyLoading(b *testing.B) {
	ctx := context.Background()
	cache := kv.NewMemoryStore()
	appID := "bench-app-" + uuid.New().String()

	testPolicies := []CachedPolicy{
		{Type: model.PolicyTypeRateLimit, Config: []byte(`{"requests_per_minute":1000}`)},
		{Type: model.PolicyTypeTokenLimit, Config: []byte(`{"max_tokens":8192}`)},
		{Type: model.PolicyTypeModelAllowlist, Config: []byte(`{"allowed_models":["gpt-4"]}`)},
	}
	_ = SetCachedPoliciesRaw(ctx, cache, appID, testPolicies)

	engine := NewEngine(nil, cache)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = engine.LoadPolicies(ctx, appID)
		}
	})
}
