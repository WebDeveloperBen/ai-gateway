package kv

import "fmt"

type KvStoreType string

const (
	BackendMemory KvStoreType = "memory"
	BackendRedis  KvStoreType = "redis"
)

type Config struct {
	Backend   KvStoreType
	RedisAddr string
	RedisPW   string
	RedisDB   int
}

func NewDriver(cfg Config) (KvStore, error) {
	switch cfg.Backend {
	case BackendMemory:
		return NewMemoryStore(), nil
	case BackendRedis:
		return NewRedisStore(cfg.RedisAddr, cfg.RedisPW, cfg.RedisDB), nil
	default:
		return nil, fmt.Errorf("unsupported kv backend: %s", cfg.Backend)
	}
}
