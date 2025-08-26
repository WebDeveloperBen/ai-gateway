// Package kv provides a uniform key-value interface and pluggable backends for Redis or in-memory store.
package kv

import (
	"context"
	"time"
)

type KvStore interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
}
