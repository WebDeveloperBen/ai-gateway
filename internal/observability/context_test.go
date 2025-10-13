package observability

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithObservability(t *testing.T) {
	ctx := context.Background()
	obs := newNoopObservability()

	newCtx := WithObservability(ctx, obs)

	require.NotNil(t, newCtx)
	assert.NotEqual(t, ctx, newCtx)
}

func TestFromContext(t *testing.T) {
	t.Run("Returns observability from context", func(t *testing.T) {
		ctx := context.Background()
		obs := newNoopObservability()
		ctx = WithObservability(ctx, obs)

		retrieved := FromContext(ctx)
		require.NotNil(t, retrieved)
		assert.Equal(t, obs, retrieved)
	})

	t.Run("Returns noop observability when not in context", func(t *testing.T) {
		ctx := context.Background()
		obs := FromContext(ctx)

		require.NotNil(t, obs)
		assert.NotNil(t, obs.tracer)
		assert.NotNil(t, obs.meter)
	})

	t.Run("Returns noop observability when value is wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), observabilityKey{}, "wrong-type")
		obs := FromContext(ctx)

		require.NotNil(t, obs)
		assert.NotNil(t, obs.tracer)
		assert.NotNil(t, obs.meter)
	})
}

func TestNoopObservability(t *testing.T) {
	obs := newNoopObservability()

	require.NotNil(t, obs)
	assert.NotNil(t, obs.tracer)
	assert.NotNil(t, obs.meter)

	tracer := obs.GetTracer()
	assert.NotNil(t, tracer)

	meter := obs.GetMeter()
	assert.NotNil(t, meter)
}

func TestObservability_GetTracer(t *testing.T) {
	obs := newNoopObservability()
	tracer := obs.GetTracer()

	assert.NotNil(t, tracer)
}

func TestObservability_GetMeter(t *testing.T) {
	obs := newNoopObservability()
	meter := obs.GetMeter()

	assert.NotNil(t, meter)
}

func TestObservability_StartSpan(t *testing.T) {
	obs := newNoopObservability()
	ctx := context.Background()

	spanCtx, span := obs.StartSpan(ctx, "test-span")

	assert.NotNil(t, spanCtx)
	assert.NotNil(t, span)
	assert.NotEqual(t, ctx, spanCtx)

	span.End()
}

func TestNew_Disabled(t *testing.T) {
	cfg := Config{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Enabled:        false,
	}

	obs, err := New(cfg)

	require.NoError(t, err)
	require.NotNil(t, obs)
	assert.NotNil(t, obs.tracer)
	assert.NotNil(t, obs.meter)
}
