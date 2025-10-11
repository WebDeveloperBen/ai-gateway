package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// RecordPolicyCheck records metrics for a policy check
func (o *Observability) RecordPolicyCheck(ctx context.Context, policyType string, duration time.Duration, violated bool) {
	if o.PolicyCheckDuration == nil {
		return
	}

	attrs := []attribute.KeyValue{
		attribute.String("policy.type", policyType),
		attribute.Bool("policy.violated", violated),
	}

	o.PolicyCheckDuration.Record(ctx, float64(duration.Milliseconds()),
		metric.WithAttributes(attrs...))

	if violated {
		o.PolicyViolations.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}

// RecordPolicyCacheHit records a policy cache hit
func (o *Observability) RecordPolicyCacheHit(ctx context.Context, tier string) {
	if o.PolicyCacheHits == nil {
		return
	}

	o.PolicyCacheHits.Add(ctx, 1, metric.WithAttributes(
		attribute.String("cache.tier", tier),
	))
}

// RecordPolicyCacheMiss records a policy cache miss
func (o *Observability) RecordPolicyCacheMiss(ctx context.Context, tier string) {
	if o.PolicyCacheMisses == nil {
		return
	}

	o.PolicyCacheMisses.Add(ctx, 1, metric.WithAttributes(
		attribute.String("cache.tier", tier),
	))
}

// RecordPolicyLoadError records a policy load error
func (o *Observability) RecordPolicyLoadError(ctx context.Context, appID string, err error) {
	if o.PolicyLoadErrors == nil {
		return
	}

	o.PolicyLoadErrors.Add(ctx, 1, metric.WithAttributes(
		attribute.String("app.id", appID),
		attribute.String("error", err.Error()),
	))
}

// RecordRateLimitCheck records a rate limit check
func (o *Observability) RecordRateLimitCheck(ctx context.Context, limitType string, violated bool) {
	if o.RateLimitChecks == nil {
		return
	}

	attrs := []attribute.KeyValue{
		attribute.String("limit.type", limitType),
		attribute.Bool("violated", violated),
	}

	o.RateLimitChecks.Add(ctx, 1, metric.WithAttributes(attrs...))

	if violated && o.RateLimitViolations != nil {
		o.RateLimitViolations.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}

// RecordTokenEstimation records token estimation metrics
func (o *Observability) RecordTokenEstimation(ctx context.Context, model string, duration time.Duration, tokens int) {
	if o.TokenEstimationDuration == nil {
		return
	}

	attrs := []attribute.KeyValue{
		attribute.String("model", model),
	}

	o.TokenEstimationDuration.Record(ctx, float64(duration.Milliseconds()),
		metric.WithAttributes(attrs...))
}

// RecordLLMTokens records LLM token usage
func (o *Observability) RecordLLMTokens(ctx context.Context, provider, model string, promptTokens, completionTokens, totalTokens int) {
	attrs := []attribute.KeyValue{
		attribute.String("provider", provider),
		attribute.String("model", model),
	}

	if o.PromptTokens != nil {
		o.PromptTokens.Record(ctx, int64(promptTokens), metric.WithAttributes(attrs...))
	}

	if o.CompletionTokens != nil {
		o.CompletionTokens.Record(ctx, int64(completionTokens), metric.WithAttributes(attrs...))
	}

	if o.TotalTokens != nil {
		o.TotalTokens.Record(ctx, int64(totalTokens), metric.WithAttributes(attrs...))
	}
}

// RecordHTTPRequest records HTTP request metrics
func (o *Observability) RecordHTTPRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration, requestSize, responseSize int) {
	attrs := []attribute.KeyValue{
		attribute.String("http.method", method),
		attribute.String("http.route", path),
		attribute.Int("http.status_code", statusCode),
	}

	if o.RequestDuration != nil {
		o.RequestDuration.Record(ctx, float64(duration.Milliseconds()),
			metric.WithAttributes(attrs...))
	}

	if o.RequestSize != nil {
		o.RequestSize.Record(ctx, int64(requestSize), metric.WithAttributes(attrs...))
	}

	if o.ResponseSize != nil {
		o.ResponseSize.Record(ctx, int64(responseSize), metric.WithAttributes(attrs...))
	}
}

// StartPolicySpan creates a span for policy operations
func (o *Observability) StartPolicySpan(ctx context.Context, operation string) (context.Context, trace.Span) {
	return o.StartSpan(ctx, "policy."+operation,
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			attribute.String("component", "policy"),
		))
}

// StartLLMSpan creates a span for LLM operations
func (o *Observability) StartLLMSpan(ctx context.Context, provider, model string) (context.Context, trace.Span) {
	return o.StartSpan(ctx, "llm.request",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("llm.provider", provider),
			attribute.String("llm.model", model),
		))
}

// RecordCircuitBreakerStateChange records circuit breaker state changes
func (o *Observability) RecordCircuitBreakerStateChange(ctx context.Context, name, from, to string) {
	if o.CircuitBreakerStateChanges == nil {
		return
	}

	attrs := []attribute.KeyValue{
		attribute.String("circuit_breaker.name", name),
		attribute.String("circuit_breaker.from_state", from),
		attribute.String("circuit_breaker.to_state", to),
	}

	o.CircuitBreakerStateChanges.Add(ctx, 1, metric.WithAttributes(attrs...))

	if to == "open" && o.CircuitBreakerTrips != nil {
		o.CircuitBreakerTrips.Add(ctx, 1, metric.WithAttributes(
			attribute.String("circuit_breaker.name", name),
		))
	}
}
