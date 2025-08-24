// Package kv provides a uniform key-value interface and pluggable backends for Redis or in-memory store.
package kv

import (
	"context"
	"fmt"
	"time"
)

type Backend string

const (
	BackendMemory Backend = "memory"
	BackendRedis  Backend = "redis"
)

type Config struct {
	Backend   Backend
	RedisAddr string
	RedisPW   string
	RedisDB   int
}

func New(cfg Config) (Store, error) {
	switch cfg.Backend {
	case BackendMemory:
		return NewMemoryStore(), nil
	case BackendRedis:
		return NewRedisStore(cfg.RedisAddr, cfg.RedisPW, cfg.RedisDB), nil
	default:
		return nil, fmt.Errorf("unsupported kv backend: %s", cfg.Backend)
	}
}

type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
}
