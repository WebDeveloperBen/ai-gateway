package policies

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
)

// RateLimiter provides Redis-backed rate limiting using sliding window algorithm
type RateLimiter struct {
	cache kv.KvStore
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(cache kv.KvStore) *RateLimiter {
	return &RateLimiter{cache: cache}
}

// CheckAndIncrement checks if the limit is exceeded and increments the counter if not
// Returns true if the request is allowed, false if rate limited
// Uses atomic Redis INCR operation to avoid race conditions
func (rl *RateLimiter) CheckAndIncrement(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	if limit <= 0 {
		// No limit configured
		return true, nil
	}

	log.Printf("DEBUG: RateLimiter.CheckAndIncrement - key=%s, limit=%d", key, limit)

	// Atomically increment the counter (handles non-existent keys)
	newCount, err := rl.cache.Incr(ctx, key)
	if err != nil {
		// Redis error - fail open to avoid blocking traffic
		log.Printf("DEBUG: RateLimiter.CheckAndIncrement - Redis error: %v", err)
		return true, err
	}

	log.Printf("DEBUG: RateLimiter.CheckAndIncrement - newCount=%d", newCount)

	// Set TTL on first increment
	if newCount == 1 {
		// Key was just created, set expiration
		_, err := rl.cache.Expire(ctx, key, window)
		if err != nil {
			// Failed to set TTL, but counter is incremented
			// Log error but don't fail request
			log.Printf("DEBUG: RateLimiter.CheckAndIncrement - Failed to set TTL: %v", err)
		}
	}

	// Check if limit exceeded (after incrementing)
	if newCount > int64(limit) {
		log.Printf("DEBUG: RateLimiter.CheckAndIncrement - Limit exceeded: %d > %d", newCount, limit)
		return false, nil
	}

	log.Printf("DEBUG: RateLimiter.CheckAndIncrement - Request allowed")
	return true, nil
}

// Increment atomically increments a counter by amount (for post-check tracking)
func (rl *RateLimiter) Increment(ctx context.Context, key string, amount int, window time.Duration) error {
	// Atomically increment by amount
	newCount, err := rl.cache.IncrBy(ctx, key, int64(amount))
	if err != nil {
		return err
	}

	// Set TTL on first increment
	if newCount == int64(amount) {
		// Key was just created, set expiration
		_, err := rl.cache.Expire(ctx, key, window)
		if err != nil {
			// Failed to set TTL, but counter is incremented
			return err
		}
	}

	return nil
}

// GetCount gets the current count for a key
func (rl *RateLimiter) GetCount(ctx context.Context, key string) (int, error) {
	// Use IncrBy(0) to atomically read the value
	// This is safer than Get() for concurrent access
	count, err := rl.cache.IncrBy(ctx, key, 0)
	if err != nil {
		// Key doesn't exist or error
		return 0, nil
	}

	return int(count), nil
}

// RateLimitKey generates a Redis key for rate limiting
func RateLimitKey(appID string, metric string) string {
	// Use current minute as the window
	minute := time.Now().Truncate(time.Minute).Unix()
	return fmt.Sprintf("ratelimit:%s:%s:%d", appID, metric, minute)
}
