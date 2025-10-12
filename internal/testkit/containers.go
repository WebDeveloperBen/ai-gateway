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
func SetupTestContainers(t *testing.T, config ContainerConfig) (string, string) {
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
	pgConnStr, pgContainer, err := createPostgresContainer(ctx, config)
	if err != nil {
		containerInitErr = err
		return
	}
	postgresContainer = pgContainer
	postgresConnStr = pgConnStr

	// Run migrations
	if err := runMigrationsOnContainer(pgConnStr); err != nil {
		containerInitErr = fmt.Errorf("failed to run migrations: %w", err)
		return
	}

	// Start Redis container
	redisAddrStr, redisCont, err := createRedisContainer(ctx, config)
	if err != nil {
		containerInitErr = err
		return
	}
	redisContainer = redisCont
	redisAddr = redisAddrStr
}

// createPostgresContainer creates a postgres container and returns the connection string
func createPostgresContainer(ctx context.Context, config ContainerConfig) (string, *postgres.PostgresContainer, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase(config.PostgresDB),
		postgres.WithUsername(config.PostgresUser),
		postgres.WithPassword(config.PostgresPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		return "", nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		return "", nil, fmt.Errorf("failed to get postgres port: %w", err)
	}

	pgConnStr := fmt.Sprintf("postgres://%s:%s@127.0.0.1:%s/%s?sslmode=disable", config.PostgresUser, config.PostgresPassword, port.Port(), config.PostgresDB)

	log.Printf("🗄️ Shared Postgres ready: %s", pgConnStr)
	return pgConnStr, pgContainer, nil
}

// createRedisContainer creates a redis container and returns the address
func createRedisContainer(ctx context.Context, config ContainerConfig) (string, *redis.RedisContainer, error) {
	redisCont, err := redis.Run(ctx,
		"redis:8-alpine",
		redis.WithSnapshotting(10, 1),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("6379/tcp").
				WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		return "", nil, fmt.Errorf("failed to start redis container: %w", err)
	}

	port, err := redisCont.MappedPort(ctx, "6379")
	if err != nil {
		return "", nil, fmt.Errorf("failed to get redis port: %w", err)
	}

	redisAddrStr := fmt.Sprintf("127.0.0.1:%s", port.Port())

	log.Printf("🔴 Shared Redis ready: %s", redisAddrStr)

	return redisAddrStr, redisCont, nil
}

// isDockerAvailable checks if Docker is available on the system
func isDockerAvailable() bool {
	cmd := exec.Command("docker", "info")
	err := cmd.Run()
	return err == nil
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

// CleanupSharedContainers forces cleanup of shared containers
// This is useful for TestMain to ensure cleanup happens at the end
func CleanupSharedContainers() {
	cleanupOnce.Do(func() {
		ctx := context.Background()

		if postgresContainer != nil {
			if err := postgresContainer.Terminate(ctx); err != nil {
				log.Printf("Warning: Failed to terminate shared postgres container: %v", err)
			}
		}

		if redisContainer != nil {
			if err := redisContainer.Terminate(ctx); err != nil {
				log.Printf("Warning: Failed to terminate shared redis container: %v", err)
			}
		}
	})
}
