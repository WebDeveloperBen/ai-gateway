package testkit

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/WebDeveloperBen/ai-gateway/internal/migrate"
)

var (
	// Shared containers for faster test execution
	postgresContainer *postgres.PostgresContainer
	redisContainer    *redis.RedisContainer
	setupOnce         sync.Once
	cleanupOnce       sync.Once
	sharedCleanupFn   func()
	containerInitErr  error

	// Connection strings
	postgresConnStr string
	redisAddr       string
)

// ContainerConfig holds configuration for test containers
type ContainerConfig struct {
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
	RedisPassword    string
}

// DefaultContainerConfig returns sensible defaults for test containers
func DefaultContainerConfig() ContainerConfig {
	return ContainerConfig{
		PostgresDB:       "testdb",
		PostgresUser:     "postgres",
		PostgresPassword: "postgres",
		RedisPassword:    "",
	}
}

// SetupTestContainers initializes shared postgres and redis containers for integration testing
// Containers are started once and reused across all tests for better performance
func SetupTestContainers(t *testing.T, config ContainerConfig) (postgresConn, redisAddr string) {
	t.Helper()

	// Skip if Docker is not available
	if !isDockerAvailable() {
		t.Skip("Docker not available, skipping integration tests")
	}

	setupOnce.Do(func() {
		setupSharedContainers(config)
	})

	if containerInitErr != nil {
		t.Fatalf("Container initialization failed: %v", containerInitErr)
	}

	return postgresConnStr, redisAddr
}

// setupSharedContainers initializes the shared containers once
func setupSharedContainers(config ContainerConfig) {
	ctx := context.Background()

	// Start Postgres container
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(config.PostgresDB),
		postgres.WithUsername(config.PostgresUser),
		postgres.WithPassword(config.PostgresPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		containerInitErr = fmt.Errorf("failed to start postgres container: %w", err)
		return
	}

	pgConnStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		containerInitErr = fmt.Errorf("failed to get postgres connection string: %w", err)
		cleanupContainers()
		return
	}

	// Run migrations on postgres
	if err := runMigrationsOnContainer(pgConnStr); err != nil {
		containerInitErr = fmt.Errorf("failed to run migrations: %w", err)
		cleanupContainers()
		return
	}

	// Start Redis container
	redisCont, err := redis.Run(ctx,
		"redis:8-alpine",
		redis.WithSnapshotting(10, 1),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithOccurrence(1).WithStartupTimeout(10*time.Second)),
	)
	if err != nil {
		containerInitErr = fmt.Errorf("failed to start redis container: %w", err)
		cleanupContainers()
		return
	}

	redisAddrStr, err := redisCont.ConnectionString(ctx)
	if err != nil {
		containerInitErr = fmt.Errorf("failed to get redis connection string: %w", err)
		cleanupContainers()
		return
	}

	// Strip the redis:// prefix for go-redis client
	if len(redisAddrStr) > 8 && redisAddrStr[:8] == "redis://" {
		redisAddrStr = redisAddrStr[8:]
	}

	// Store references
	postgresContainer = pgContainer
	redisContainer = redisCont
	postgresConnStr = pgConnStr
	redisAddr = redisAddrStr

	// Set up shared cleanup
	sharedCleanupFn = func() {
		cleanupContainers()
	}

	log.Printf("🗄️ Shared Postgres ready: %s", pgConnStr)
	log.Printf("🔴 Shared Redis ready: %s", redisAddrStr)
}

// runMigrationsOnContainer runs database migrations on the test container
func runMigrationsOnContainer(connStr string) error {
	// Change to project root for migrations
	oldWd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	projectRoot := findProjectRoot(oldWd)
	if err := os.Chdir(projectRoot); err != nil {
		return fmt.Errorf("failed to change to project root: %w", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			log.Printf("Warning: Failed to restore working directory: %v", err)
		}
	}()

	return migrate.InitDatabase(context.Background(), connStr)
}

// cleanupContainers terminates all shared containers
func cleanupContainers() {
	ctx := context.Background()

	if postgresContainer != nil {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Printf("Warning: Failed to terminate shared postgres container: %v", err)
		}
		postgresContainer = nil
	}

	if redisContainer != nil {
		if err := redisContainer.Terminate(ctx); err != nil {
			log.Printf("Warning: Failed to terminate shared redis container: %v", err)
		}
		redisContainer = nil
	}
}

// SetupPostgresOnly sets up shared postgres container for tests that don't need redis
func SetupPostgresOnly(t *testing.T, config ContainerConfig) (connStr string) {
	t.Helper()

	// Skip if Docker is not available
	if !isDockerAvailable() {
		t.Skip("Docker not available, skipping integration tests")
	}

	setupOnce.Do(func() {
		setupSharedPostgres(config)
	})

	if containerInitErr != nil {
		t.Fatalf("Postgres container initialization failed: %v", containerInitErr)
	}

	return postgresConnStr
}

// SetupRedisOnly sets up shared redis container for tests that don't need postgres
func SetupRedisOnly(t *testing.T, config ContainerConfig) (addr string) {
	t.Helper()

	// Skip if Docker is not available
	if !isDockerAvailable() {
		t.Skip("Docker not available, skipping integration tests")
	}

	setupOnce.Do(func() {
		setupSharedRedis(config)
	})

	if containerInitErr != nil {
		t.Fatalf("Redis container initialization failed: %v", containerInitErr)
	}

	return redisAddr
}

// setupSharedPostgres initializes only the shared postgres container
func setupSharedPostgres(config ContainerConfig) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(config.PostgresDB),
		postgres.WithUsername(config.PostgresUser),
		postgres.WithPassword(config.PostgresPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		containerInitErr = fmt.Errorf("failed to start postgres container: %w", err)
		return
	}

	pgConnStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		containerInitErr = fmt.Errorf("failed to get postgres connection string: %w", err)
		cleanupContainers()
		return
	}

	// Run migrations
	if err := runMigrationsOnContainer(pgConnStr); err != nil {
		containerInitErr = fmt.Errorf("failed to run migrations: %w", err)
		cleanupContainers()
		return
	}

	postgresContainer = pgContainer
	postgresConnStr = pgConnStr

	sharedCleanupFn = func() {
		if postgresContainer != nil {
			if err := postgresContainer.Terminate(context.Background()); err != nil {
				log.Printf("Warning: Failed to terminate shared postgres container: %v", err)
			}
		}
	}

	log.Printf("🗄️ Shared Postgres ready: %s", pgConnStr)
}

// setupSharedRedis initializes only the shared redis container
func setupSharedRedis(config ContainerConfig) {
	ctx := context.Background()

	redisCont, err := redis.Run(ctx,
		"redis:8-alpine",
		redis.WithSnapshotting(10, 1),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithOccurrence(1).WithStartupTimeout(10*time.Second)),
	)
	if err != nil {
		containerInitErr = fmt.Errorf("failed to start redis container: %w", err)
		return
	}

	redisAddrStr, err := redisCont.ConnectionString(ctx)
	if err != nil {
		containerInitErr = fmt.Errorf("failed to get redis connection string: %w", err)
		cleanupContainers()
		return
	}

	// Strip the redis:// prefix for go-redis client
	if len(redisAddrStr) > 8 && redisAddrStr[:8] == "redis://" {
		redisAddrStr = redisAddrStr[8:]
	}

	redisContainer = redisCont
	redisAddr = redisAddrStr

	sharedCleanupFn = func() {
		if redisContainer != nil {
			if err := redisContainer.Terminate(context.Background()); err != nil {
				log.Printf("Warning: Failed to terminate shared redis container: %v", err)
			}
		}
	}

	log.Printf("🔴 Shared Redis ready: %s", redisAddrStr)
}

// isDockerAvailable checks if Docker is available on the system
func isDockerAvailable() bool {
	cmd := exec.Command("docker", "info")
	err := cmd.Run()
	return err == nil
}

// CleanupSharedContainers forces cleanup of shared containers
// This is useful for TestMain to ensure cleanup happens at the end
func CleanupSharedContainers() {
	cleanupOnce.Do(func() {
		if sharedCleanupFn != nil {
			sharedCleanupFn()
		}
	})
}
