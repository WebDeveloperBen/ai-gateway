package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	_ "github.com/jackc/pgx/v5/stdlib" // Import pgx driver for database/sql
)

// InitDatabase runs the full migration sequence:
// 1. Pre-migration SQL (extensions, functions)
// 2. Goose migrations (tables, indexes)
// 3. Post-migration SQL (triggers, policies, seeds)
func InitDatabase(ctx context.Context, connectionConfig string) error {
	if connectionConfig == "" {
		return fmt.Errorf("database connection configuration is required")
	}

	// Convert connection config to standard connection string for database/sql
	connectionString, err := resolveConnectionString(ctx, connectionConfig)
	if err != nil {
		return fmt.Errorf("failed to resolve connection string: %w", err)
	}

	// Open database connection for SQL execution
	sqlDB, err := sql.Open("pgx", connectionString)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer sqlDB.Close()

	// 1. Run pre-migration SQL (extensions, functions)
	if err := runPreMigrationSQL(sqlDB); err != nil {
		return fmt.Errorf("failed to run pre-migration SQL: %w", err)
	}

	// 2. Run Goose migrations
	if err := runGooseMigrations(connectionString); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// 3. Run post-migration SQL (triggers, policies, seeds)
	if err := runPostMigrationSQL(sqlDB); err != nil {
		return fmt.Errorf("failed to run post-migration SQL: %w", err)
	}

	return nil
}

// resolveConnectionString converts either a standard connection string or Azure managed identity config
// into a standard PostgreSQL connection string usable by database/sql
func resolveConnectionString(ctx context.Context, config string) (string, error) {
	// Check if config looks like a connection string or Azure managed identity config
	if strings.Contains(config, "://") || strings.Contains(config, "host=") {
		// Standard connection string
		return config, nil
	} else {
		// Use our existing managed identity logic from the db driver
		driver, err := db.NewPostgresDriver(ctx, config)
		if err != nil {
			return "", fmt.Errorf("failed to resolve managed identity connection: %w", err)
		}
		driver.Pool.Close() // Close the pool, we just needed the connection string

		// Extract connection string from the pool config (this is a bit hacky but works)
		// In a real implementation, you'd refactor buildManagedIdentityConnection to be public
		return buildManagedIdentityConnectionString(config)
	}
}

// buildManagedIdentityConnectionString duplicates the logic from db driver for migration use
func buildManagedIdentityConnectionString(config string) (string, error) {
	parts := strings.Split(config, ":")
	if len(parts) != 3 {
		return "", fmt.Errorf("managed identity config must be in format 'server:database:user'")
	}

	server, database, user := parts[0], parts[1], parts[2]
	// FIX: this azure logic
	// For now, this is a simplified version - you'd need to import the Azure SDK
	// and replicate the token logic from postgres.go
	return fmt.Sprintf("host=%s user=%s dbname=%s sslmode=require", server, user, database), nil
}

func runGooseMigrations(connectionString string) error {
	// Run goose up command
	cmd := exec.Command("goose", "-dir", "./db/migrations", "postgres", connectionString, "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run goose migrations: %w\nOutput: %s", err, output)
	}

	log.Printf("Goose migrations completed: %s", output)
	return nil
}

func runPreMigrationSQL(db *sql.DB) error {
	// Read and execute db/schema_pre.sql
	content, err := os.ReadFile("db/schema_pre.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema_pre.sql: %w", err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute pre-migration SQL: %w", err)
	}

	log.Println("Pre-migration SQL executed successfully")
	return nil
}

func runPostMigrationSQL(db *sql.DB) error {
	// Read and execute db/schema_post.sql
	content, err := os.ReadFile("db/schema_post.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema_post.sql: %w", err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute post-migration SQL: %w", err)
	}

	log.Println("Post-migration SQL executed successfully")
	return nil
}

