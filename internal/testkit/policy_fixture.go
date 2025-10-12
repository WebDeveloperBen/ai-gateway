package testkit

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/config"
	dbdriver "github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/policies"
	"github.com/WebDeveloperBen/ai-gateway/internal/migrate"
)

// PolicyTest provides a complete test environment for policy integration tests
type PolicyTest struct {
	Engine *policies.Engine
	DB     *dbdriver.Postgres
	Cache  kv.KvStore
}

// PolicyTestOption functional options for policy test configuration
type PolicyTestOption func(*policyTestOpts)

type policyTestOpts struct {
	CacheBackend      kv.KvStoreType
	UseCircuitBreaker bool
}

// NewPolicyE2E creates a policy test environment using real database and cache services
// Requires docker services (postgres, redis) to be running
func NewPolicyE2E(t *testing.T, opts ...PolicyTestOption) *PolicyTest {
	t.Helper()

	// Load environment variables from repo root
	LoadDotenvFromRepoRoot(t)

	// Reload config to ensure fresh env vars
	config.Reload()

	cfg := config.Envs

	// Default options
	o := policyTestOpts{
		CacheBackend:      kv.KvStoreType(cfg.KVBackend),
		UseCircuitBreaker: cfg.EnableRedisCircuitBreaker && cfg.KVBackend == "redis",
	}
	for _, opt := range opts {
		opt(&o)
	}

	ctx := context.Background()

	// Set up database connection
	pg, err := dbdriver.NewPostgresDriver(ctx, cfg.DBConnectionString)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	t.Cleanup(func() { pg.Pool.Close() })

	// Change to project root directory for migrations (they expect db/ directory)
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Find project root
	projectRoot := findProjectRoot(oldWd)
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Logf("Warning: Failed to restore working directory: %v", err)
		}
	}()

	// Run migrations to ensure schema is up to date
	if err := migrate.InitDatabase(ctx, cfg.GetDatabaseConnection()); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Set up cache (KV store)
	cache, err := kv.NewDriver(kv.Config{
		Backend:   o.CacheBackend,
		RedisAddr: cfg.RedisAddr,
		RedisPW:   cfg.RedisPW,
		RedisDB:   0,
	})
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Add circuit breaker if requested
	if o.UseCircuitBreaker {
		cache = kv.NewCircuitBreakerStore(cache, kv.DefaultCircuitBreakerConfig())
	}

	// Create policy engine
	engine := policies.NewEngine(pg.Queries, cache)

	return &PolicyTest{
		Engine: engine,
		DB:     pg,
		Cache:  cache,
	}
}

// WithMemoryCache configures the test to use in-memory cache instead of Redis
func WithMemoryCache() PolicyTestOption {
	return func(o *policyTestOpts) {
		o.CacheBackend = kv.BackendMemory
		o.UseCircuitBreaker = false
	}
}

// WithRedisCache configures the test to use Redis cache
func WithRedisCache() PolicyTestOption {
	return func(o *policyTestOpts) {
		o.CacheBackend = kv.BackendRedis
	}
}

// WithCircuitBreaker enables circuit breaker for Redis cache
func WithCircuitBreaker() PolicyTestOption {
	return func(o *policyTestOpts) {
		o.UseCircuitBreaker = true
	}
}

// findProjectRoot finds the project root by looking for go.mod
func findProjectRoot(dir string) string {
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root, return current dir
			return dir
		}
		dir = parent
	}
}
