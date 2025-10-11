package kv

import (
	"context"
	"log"
	"time"

	"github.com/sony/gobreaker"
)

// CircuitBreakerStore wraps a KvStore with a circuit breaker to prevent
// cascading failures when Redis is unavailable. When the circuit is open,
// operations fail fast instead of waiting for timeouts.
type CircuitBreakerStore struct {
	store   KvStore
	breaker *gobreaker.CircuitBreaker
}

// CircuitBreakerConfig configures the circuit breaker behavior
type CircuitBreakerConfig struct {
	// MaxRequests is the maximum number of requests allowed to pass through
	// when the CircuitBreaker is half-open (default: 1)
	MaxRequests uint32

	// Interval is the cyclic period of the closed state for the CircuitBreaker
	// to clear the internal Counts (default: 60s)
	Interval time.Duration

	// Timeout is the period of the open state, after which the state becomes half-open
	// (default: 60s)
	Timeout time.Duration

	// ReadyToTrip returns true when the CircuitBreaker should trip from closed to open
	// Default: trips after 5 consecutive failures
	ReadyToTrip func(counts gobreaker.Counts) bool
}

// DefaultCircuitBreakerConfig returns a sensible default configuration
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		MaxRequests: 1,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	}
}

// NewCircuitBreakerStore wraps a KvStore with a circuit breaker
func NewCircuitBreakerStore(store KvStore, config CircuitBreakerConfig) *CircuitBreakerStore {
	settings := gobreaker.Settings{
		Name:        "redis-kv-store",
		MaxRequests: config.MaxRequests,
		Interval:    config.Interval,
		Timeout:     config.Timeout,
		ReadyToTrip: config.ReadyToTrip,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("[CircuitBreaker] %s: state changed from %s to %s\n", name, from.String(), to.String())
		},
	}

	return &CircuitBreakerStore{
		store:   store,
		breaker: gobreaker.NewCircuitBreaker(settings),
	}
}

// Get retrieves a value from the store through the circuit breaker
func (cb *CircuitBreakerStore) Get(ctx context.Context, key string) (string, error) {
	result, err := cb.breaker.Execute(func() (any, error) {
		return cb.store.Get(ctx, key)
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

// Set stores a value through the circuit breaker
func (cb *CircuitBreakerStore) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	_, err := cb.breaker.Execute(func() (any, error) {
		return nil, cb.store.Set(ctx, key, value, ttl)
	})
	return err
}

// Del deletes a key through the circuit breaker
func (cb *CircuitBreakerStore) Del(ctx context.Context, key string) error {
	_, err := cb.breaker.Execute(func() (any, error) {
		return nil, cb.store.Del(ctx, key)
	})
	return err
}

// Exists checks if a key exists through the circuit breaker
func (cb *CircuitBreakerStore) Exists(ctx context.Context, key string) (bool, error) {
	result, err := cb.breaker.Execute(func() (any, error) {
		return cb.store.Exists(ctx, key)
	})
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

// Incr atomically increments a value through the circuit breaker
func (cb *CircuitBreakerStore) Incr(ctx context.Context, key string) (int64, error) {
	result, err := cb.breaker.Execute(func() (any, error) {
		return cb.store.Incr(ctx, key)
	})
	if err != nil {
		return 0, err
	}
	return result.(int64), nil
}

// IncrBy atomically increments a value by amount through the circuit breaker
func (cb *CircuitBreakerStore) IncrBy(ctx context.Context, key string, amount int64) (int64, error) {
	result, err := cb.breaker.Execute(func() (any, error) {
		return cb.store.IncrBy(ctx, key, amount)
	})
	if err != nil {
		return 0, err
	}
	return result.(int64), nil
}

// Expire sets a TTL on a key through the circuit breaker
func (cb *CircuitBreakerStore) Expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	result, err := cb.breaker.Execute(func() (any, error) {
		return cb.store.Expire(ctx, key, ttl)
	})
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

// ScanGetAll scans keys matching a pattern through the circuit breaker
func (cb *CircuitBreakerStore) ScanGetAll(ctx context.Context, pattern string, count int64) (map[string]string, error) {
	result, err := cb.breaker.Execute(func() (any, error) {
		return cb.store.ScanGetAll(ctx, pattern, count)
	})
	if err != nil {
		return nil, err
	}
	return result.(map[string]string), nil
}

// ScanAll scans keys matching a pattern through the circuit breaker
func (cb *CircuitBreakerStore) ScanAll(ctx context.Context, pattern string, count int64) ([]string, error) {
	result, err := cb.breaker.Execute(func() (any, error) {
		return cb.store.ScanAll(ctx, pattern, count)
	})
	if err != nil {
		return nil, err
	}
	return result.([]string), nil
}

// Close closes the underlying store
func (cb *CircuitBreakerStore) Close(ctx context.Context) error {
	return cb.store.Close(ctx)
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreakerStore) State() gobreaker.State {
	return cb.breaker.State()
}

// Counts returns the current counts of the circuit breaker
func (cb *CircuitBreakerStore) Counts() gobreaker.Counts {
	return cb.breaker.Counts()
}
