// Package observability provides OpenTelemetry instrumentation for the AI Gateway.
//
// This package integrates metrics, traces, and spans to provide comprehensive observability for policy enforcement, rate limiting, token usage, and HTTP requests.
package observability

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type Config struct {
	ServiceName    string
	ServiceVersion string
	OTLPEndpoint   string
	Enabled        bool
	SampleRate     float64
}

type Observability struct {
	tracer         trace.Tracer
	meter          metric.Meter
	meterProvider  *sdkmetric.MeterProvider
	tracerProvider *sdktrace.TracerProvider

	// Policy metrics
	PolicyCheckDuration metric.Float64Histogram
	PolicyCacheHits     metric.Int64Counter
	PolicyCacheMisses   metric.Int64Counter
	PolicyLoadErrors    metric.Int64Counter
	PolicyViolations    metric.Int64Counter

	// Rate limiter metrics
	RateLimitViolations metric.Int64Counter
	RateLimitChecks     metric.Int64Counter

	// Token metrics
	TokenEstimationDuration metric.Float64Histogram
	PromptTokens            metric.Int64Histogram
	CompletionTokens        metric.Int64Histogram
	TotalTokens             metric.Int64Histogram

	// Request metrics
	RequestDuration metric.Float64Histogram
	RequestSize     metric.Int64Histogram
	ResponseSize    metric.Int64Histogram

	// Circuit breaker metrics
	CircuitBreakerStateChanges metric.Int64Counter
	CircuitBreakerTrips        metric.Int64Counter
}

func New(cfg Config) (*Observability, error) {
	if !cfg.Enabled {
		return newNoopObservability(), nil
	}

	ctx := context.Background()

	// Create resource
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Set up trace provider
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SampleRate)),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Set up metric provider
	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(cfg.OTLPEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter,
			sdkmetric.WithInterval(10*time.Second),
		)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	tracer := tracerProvider.Tracer(cfg.ServiceName)
	meter := meterProvider.Meter(cfg.ServiceName)

	obs := &Observability{
		tracer:         tracer,
		meter:          meter,
		meterProvider:  meterProvider,
		tracerProvider: tracerProvider,
	}

	// Initialize all metrics
	if err := obs.initMetrics(); err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	return obs, nil
}

func (o *Observability) initMetrics() error {
	var err error

	// Policy metrics
	o.PolicyCheckDuration, err = o.meter.Float64Histogram(
		"policy.check.duration",
		metric.WithDescription("Duration of policy checks in milliseconds"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return err
	}

	o.PolicyCacheHits, err = o.meter.Int64Counter(
		"policy.cache.hits",
		metric.WithDescription("Number of policy cache hits"),
	)
	if err != nil {
		return err
	}

	o.PolicyCacheMisses, err = o.meter.Int64Counter(
		"policy.cache.misses",
		metric.WithDescription("Number of policy cache misses"),
	)
	if err != nil {
		return err
	}

	o.PolicyLoadErrors, err = o.meter.Int64Counter(
		"policy.load.errors",
		metric.WithDescription("Number of policy load errors"),
	)
	if err != nil {
		return err
	}

	o.PolicyViolations, err = o.meter.Int64Counter(
		"policy.violations",
		metric.WithDescription("Number of policy violations by type"),
	)
	if err != nil {
		return err
	}

	// Rate limiter metrics
	o.RateLimitViolations, err = o.meter.Int64Counter(
		"ratelimit.violations",
		metric.WithDescription("Number of rate limit violations"),
	)
	if err != nil {
		return err
	}

	o.RateLimitChecks, err = o.meter.Int64Counter(
		"ratelimit.checks",
		metric.WithDescription("Number of rate limit checks"),
	)
	if err != nil {
		return err
	}

	// Token metrics
	o.TokenEstimationDuration, err = o.meter.Float64Histogram(
		"token.estimation.duration",
		metric.WithDescription("Duration of token estimation in milliseconds"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return err
	}

	o.PromptTokens, err = o.meter.Int64Histogram(
		"llm.prompt.tokens",
		metric.WithDescription("Number of prompt tokens"),
	)
	if err != nil {
		return err
	}

	o.CompletionTokens, err = o.meter.Int64Histogram(
		"llm.completion.tokens",
		metric.WithDescription("Number of completion tokens"),
	)
	if err != nil {
		return err
	}

	o.TotalTokens, err = o.meter.Int64Histogram(
		"llm.total.tokens",
		metric.WithDescription("Total number of tokens"),
	)
	if err != nil {
		return err
	}

	// Request metrics
	o.RequestDuration, err = o.meter.Float64Histogram(
		"http.request.duration",
		metric.WithDescription("Duration of HTTP requests in milliseconds"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return err
	}

	o.RequestSize, err = o.meter.Int64Histogram(
		"http.request.size",
		metric.WithDescription("Size of HTTP request body in bytes"),
		metric.WithUnit("bytes"),
	)
	if err != nil {
		return err
	}

	o.ResponseSize, err = o.meter.Int64Histogram(
		"http.response.size",
		metric.WithDescription("Size of HTTP response body in bytes"),
		metric.WithUnit("bytes"),
	)
	if err != nil {
		return err
	}

	// Circuit breaker metrics
	o.CircuitBreakerStateChanges, err = o.meter.Int64Counter(
		"circuitbreaker.state.changes",
		metric.WithDescription("Number of circuit breaker state changes"),
	)
	if err != nil {
		return err
	}

	o.CircuitBreakerTrips, err = o.meter.Int64Counter(
		"circuitbreaker.trips",
		metric.WithDescription("Number of times circuit breaker has tripped"),
	)
	if err != nil {
		return err
	}

	return nil
}

// Shutdown gracefully shuts down the observability providers
func (o *Observability) Shutdown(ctx context.Context) error {
	if o.tracerProvider != nil {
		if err := o.tracerProvider.Shutdown(ctx); err != nil {
			return err
		}
	}
	if o.meterProvider != nil {
		if err := o.meterProvider.Shutdown(ctx); err != nil {
			return err
		}
	}
	return nil
}

// StartSpan creates a new span
func (o *Observability) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return o.tracer.Start(ctx, name, opts...)
}

// GetTracer returns the tracer
func (o *Observability) GetTracer() trace.Tracer {
	return o.tracer
}

// GetMeter returns the meter
func (o *Observability) GetMeter() metric.Meter {
	return o.meter
}

// No-op implementation when observability is disabled
func newNoopObservability() *Observability {
	noopMeter := sdkmetric.NewMeterProvider().Meter("noop")
	return &Observability{
		tracer: noop.NewTracerProvider().Tracer("noop"),
		meter:  noopMeter,
	}
}
