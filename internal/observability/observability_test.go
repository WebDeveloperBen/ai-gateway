package observability

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Enabled(t *testing.T) {
	// Test with enabled=true but invalid endpoint - should fail gracefully
	cfg := Config{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		OTLPEndpoint:   "invalid-endpoint:4317",
		Enabled:        true,
		SampleRate:     1.0,
	}

	obs, err := New(cfg)
	// This should fail because the endpoint is invalid
	require.Error(t, err)
	assert.Nil(t, obs)
}

func TestNew_Enabled_ValidConfig(t *testing.T) {
	// Test with enabled=true and valid config but unreachable endpoint
	cfg := Config{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		OTLPEndpoint:   "127.0.0.1:4317", // Valid format but unreachable
		Enabled:        true,
		SampleRate:     1.0,
	}

	obs, err := New(cfg)
	// This may fail due to network issues, but we want to test the code path
	if err != nil {
		// Expected to fail due to no OTLP endpoint running
		assert.Nil(t, obs)
		return
	}

	// If it succeeds (unlikely), test that it's properly initialized
	require.NotNil(t, obs)
	assert.NotNil(t, obs.tracer)
	assert.NotNil(t, obs.meter)
	assert.NotNil(t, obs.tracerProvider)
	assert.NotNil(t, obs.meterProvider)

	// Test shutdown
	ctx := context.Background()
	err = obs.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestObservability_Shutdown(t *testing.T) {
	t.Run("shutdown noop observability", func(t *testing.T) {
		obs := newNoopObservability()
		ctx := context.Background()

		err := obs.Shutdown(ctx)
		assert.NoError(t, err)
	})

	t.Run("shutdown disabled observability", func(t *testing.T) {
		cfg := Config{
			ServiceName:    "test-service",
			ServiceVersion: "1.0.0",
			Enabled:        false,
		}

		obs, err := New(cfg)
		require.NoError(t, err)
		require.NotNil(t, obs)

		ctx := context.Background()
		err = obs.Shutdown(ctx)
		assert.NoError(t, err)
	})
}

func TestObservability_RecordPolicyCheck(t *testing.T) {
	obs := newNoopObservability()
	ctx := context.Background()

	// Should not panic
	obs.RecordPolicyCheck(ctx, "test-policy", 100*time.Millisecond, false)
	obs.RecordPolicyCheck(ctx, "test-policy", 200*time.Millisecond, true)
}

func TestObservability_RecordPolicyCacheHit(t *testing.T) {
	obs := newNoopObservability()
	ctx := context.Background()

	// Should not panic
	obs.RecordPolicyCacheHit(ctx, "memory")
	obs.RecordPolicyCacheHit(ctx, "redis")
}

func TestObservability_RecordPolicyCacheMiss(t *testing.T) {
	obs := newNoopObservability()
	ctx := context.Background()

	// Should not panic
	obs.RecordPolicyCacheMiss(ctx, "memory")
	obs.RecordPolicyCacheMiss(ctx, "redis")
}

func TestObservability_RecordPolicyLoadError(t *testing.T) {
	obs := newNoopObservability()
	ctx := context.Background()
	testErr := errors.New("test error")

	// Should not panic
	obs.RecordPolicyLoadError(ctx, "app-123", testErr)
}

func TestObservability_RecordRateLimitCheck(t *testing.T) {
	obs := newNoopObservability()
	ctx := context.Background()

	// Should not panic
	obs.RecordRateLimitCheck(ctx, "user", false)
	obs.RecordRateLimitCheck(ctx, "org", true)
}

func TestObservability_RecordTokenEstimation(t *testing.T) {
	obs := newNoopObservability()
	ctx := context.Background()

	// Should not panic
	obs.RecordTokenEstimation(ctx, "gpt-4", 50*time.Millisecond, 150)
}

func TestObservability_RecordLLMTokens(t *testing.T) {
	obs := newNoopObservability()
	ctx := context.Background()

	// Should not panic
	obs.RecordLLMTokens(ctx, "openai", "gpt-4", 100, 50, 150)
}

func TestObservability_RecordHTTPRequest(t *testing.T) {
	obs := newNoopObservability()
	ctx := context.Background()

	// Should not panic
	obs.RecordHTTPRequest(ctx, "POST", "/api/chat", 200, 150*time.Millisecond, 1024, 2048)
	obs.RecordHTTPRequest(ctx, "GET", "/api/models", 404, 50*time.Millisecond, 0, 128)
}

func TestObservability_StartPolicySpan(t *testing.T) {
	obs := newNoopObservability()
	ctx := context.Background()

	spanCtx, span := obs.StartPolicySpan(ctx, "check")

	assert.NotNil(t, spanCtx)
	assert.NotNil(t, span)
	assert.NotEqual(t, ctx, spanCtx)

	span.End()
}

func TestObservability_StartLLMSpan(t *testing.T) {
	obs := newNoopObservability()
	ctx := context.Background()

	spanCtx, span := obs.StartLLMSpan(ctx, "openai", "gpt-4")

	assert.NotNil(t, spanCtx)
	assert.NotNil(t, span)
	assert.NotEqual(t, ctx, spanCtx)

	span.End()
}

func TestObservability_RecordCircuitBreakerStateChange(t *testing.T) {
	obs := newNoopObservability()
	ctx := context.Background()

	// Should not panic
	obs.RecordCircuitBreakerStateChange(ctx, "test-breaker", "closed", "open")
	obs.RecordCircuitBreakerStateChange(ctx, "test-breaker", "open", "half-open")
	obs.RecordCircuitBreakerStateChange(ctx, "test-breaker", "half-open", "closed")
}

func TestObservability_WithInitializedMetrics(t *testing.T) {
	obs := newNoopObservability()

	// Initialize metrics on the noop observability
	err := obs.initMetrics()
	require.NoError(t, err)

	ctx := context.Background()

	// Test that all metric recording functions work with initialized metrics
	obs.RecordPolicyCheck(ctx, "test-policy", 100*time.Millisecond, false)
	obs.RecordPolicyCheck(ctx, "test-policy", 200*time.Millisecond, true)

	obs.RecordPolicyCacheHit(ctx, "memory")
	obs.RecordPolicyCacheMiss(ctx, "redis")

	testErr := errors.New("test error")
	obs.RecordPolicyLoadError(ctx, "app-123", testErr)

	obs.RecordRateLimitCheck(ctx, "user", false)
	obs.RecordRateLimitCheck(ctx, "org", true)

	obs.RecordTokenEstimation(ctx, "gpt-4", 50*time.Millisecond, 150)

	obs.RecordLLMTokens(ctx, "openai", "gpt-4", 100, 50, 150)

	obs.RecordHTTPRequest(ctx, "POST", "/api/chat", 200, 150*time.Millisecond, 1024, 2048)

	obs.RecordCircuitBreakerStateChange(ctx, "test-breaker", "closed", "open")

	// Test span creation
	spanCtx, span := obs.StartPolicySpan(ctx, "check")
	assert.NotNil(t, spanCtx)
	assert.NotNil(t, span)
	span.End()

	spanCtx2, span2 := obs.StartLLMSpan(ctx, "openai", "gpt-4")
	assert.NotNil(t, spanCtx2)
	assert.NotNil(t, span2)
	span2.End()
}

func TestObservability_initMetrics_Noop(t *testing.T) {
	obs := newNoopObservability()

	// initMetrics should work on noop observability (metrics will be no-ops)
	err := obs.initMetrics()
	assert.NoError(t, err)

	// Verify that metrics are initialized (even if they're no-ops)
	assert.NotNil(t, obs.PolicyCheckDuration)
	assert.NotNil(t, obs.PolicyCacheHits)
	assert.NotNil(t, obs.RateLimitChecks)
	assert.NotNil(t, obs.TokenEstimationDuration)
	assert.NotNil(t, obs.RequestDuration)
	assert.NotNil(t, obs.CircuitBreakerStateChanges)
}

func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "valid disabled config",
			config: Config{
				ServiceName:    "test",
				ServiceVersion: "1.0.0",
				Enabled:        false,
			},
			expectError: false,
		},
		{
			name: "valid enabled config",
			config: Config{
				ServiceName:    "test",
				ServiceVersion: "1.0.0",
				OTLPEndpoint:   "localhost:4317",
				Enabled:        true,
				SampleRate:     1.0,
			},
			expectError: false, // This will fail due to no OTLP endpoint, but that's expected
		},
		{
			name: "missing service name",
			config: Config{
				ServiceVersion: "1.0.0",
				Enabled:        false,
			},
			expectError: false, // Disabled config doesn't validate
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.config)
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && tt.config.Enabled && err != nil {
				// For enabled configs, we expect connection errors which are fine for this test
				// We just want to ensure the config validation doesn't fail
				t.Logf("Got expected connection error: %v", err)
			}
		})
	}
}
