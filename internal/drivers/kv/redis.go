package kv

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ = (*RedisStore)(nil)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr string, password string, db int) *RedisStore {
	return &RedisStore{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (r *RedisStore) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisStore) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisStore) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisStore) Exists(ctx context.Context, key string) (bool, error) {
	res, err := r.client.Exists(ctx, key).Result()
	return res > 0, err // support passing multiple keys
}

// ScanAll is a Non-blocking, incremental SCAN over a pattern.
// 'count' is a hint; 512–2048 is a good starting range.
func (r *RedisStore) ScanAll(ctx context.Context, key string, count int64) ([]string, error) {
	var (
		cursor uint64
		keys   []string
	)
	for {
		k, c, err := r.client.Scan(ctx, cursor, key, count).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, k...)
		cursor = c
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

// ScanGetAll is a convienance helper - SCAN + pipeline GET all matched keys.
func (r *RedisStore) ScanGetAll(ctx context.Context, key string, count int64) (map[string]string, error) {
	keys, err := r.ScanAll(ctx, key, count)
	if err != nil || len(keys) == 0 {
		return map[string]string{}, err
	}
	pipe := r.client.Pipeline()
	cmds := make([]*redis.StringCmd, 0, len(keys))
	for _, k := range keys {
		cmds = append(cmds, pipe.Get(ctx, k))
	}
	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return nil, err
	}
	out := make(map[string]string, len(keys))
	for i, k := range keys {
		if s, err := cmds[i].Result(); err == nil {
			out[k] = s
		}
	}
	return out, nil
}

func (r *RedisStore) Close(ctx context.Context) error {
	return r.client.Close()
}
