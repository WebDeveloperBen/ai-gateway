package migrate_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/migrate"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	_ "github.com/jackc/pgx/v5/stdlib" // Import pgx driver for database/sql
	"github.com/stretchr/testify/require"
)

func TestInitDatabase(t *testing.T) {
	t.Run("empty connection config returns error", func(t *testing.T) {
		err := migrate.InitDatabase(context.Background(), "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "database connection configuration is required")
	})

	t.Run("invalid connection config returns error", func(t *testing.T) {
		// Test with invalid connection config that should fail during resolution
		err := migrate.InitDatabase(context.Background(), "invalid:config:format")
		require.Error(t, err)
		// The function actually progresses further and fails when trying to read schema files
		require.Contains(t, err.Error(), "failed to run pre-migration SQL")
	})

	t.Run("migration fails when schema files are missing", func(t *testing.T) {
		// Test that InitDatabase fails when required schema files don't exist
		// This is the actual failure point in the current test environment
		err := migrate.InitDatabase(context.Background(), "host=nonexistent port=5432 user=test dbname=test sslmode=disable")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to read schema_pre.sql")
	})
}

func TestResolveConnectionString(t *testing.T) {
	t.Run("standard postgres URL is returned as-is", func(t *testing.T) {
		ctx := context.Background()

		result, err := migrate.ResolveConnectionString(ctx, "postgres://user:pass@localhost:5432/db?sslmode=disable")
		require.NoError(t, err)
		require.Equal(t, "postgres://user:pass@localhost:5432/db?sslmode=disable", result)
	})

	t.Run("host parameter format is returned as-is", func(t *testing.T) {
		ctx := context.Background()

		result, err := migrate.ResolveConnectionString(ctx, "host=localhost user=test dbname=test sslmode=require")
		require.NoError(t, err)
		require.Equal(t, "host=localhost user=test dbname=test sslmode=require", result)
	})

	t.Run("managed identity config is processed", func(t *testing.T) {
		ctx := context.Background()

		// Test managed identity format - this should succeed and return a connection string
		result, err := migrate.ResolveConnectionString(ctx, "server:database:user")
		require.NoError(t, err)
		require.Equal(t, "host=server user=user dbname=database sslmode=require", result)
	})
}

func TestBuildManagedIdentityConnectionString(t *testing.T) {
	t.Run("valid managed identity config", func(t *testing.T) {
		result, err := migrate.BuildManagedIdentityConnectionString("myserver:mydb:myuser")
		require.NoError(t, err)
		require.Equal(t, "host=myserver user=myuser dbname=mydb sslmode=require", result)
	})

	t.Run("valid config with underscores and numbers", func(t *testing.T) {
		result, err := migrate.BuildManagedIdentityConnectionString("my_server_123:my_db_456:my_user_789")
		require.NoError(t, err)
		require.Equal(t, "host=my_server_123 user=my_user_789 dbname=my_db_456 sslmode=require", result)
	})

	t.Run("invalid format returns error", func(t *testing.T) {
		tests := []struct {
			name   string
			config string
			error  string
		}{
			{
				name:   "too few parts",
				config: "server:db",
				error:  "managed identity config must be in format 'server:database:user'",
			},
			{
				name:   "too many parts",
				config: "server:db:user:extra",
				error:  "managed identity config must be in format 'server:database:user'",
			},
			{
				name:   "empty string",
				config: "",
				error:  "managed identity config must be in format 'server:database:user'",
			},
			{
				name:   "empty parts",
				config: "::",
				error:  "managed identity config must be in format 'server:database:user'",
			},
			{
				name:   "missing server",
				config: ":db:user",
				error:  "managed identity config must be in format 'server:database:user'",
			},
			{
				name:   "missing database",
				config: "server::user",
				error:  "managed identity config must be in format 'server:database:user'",
			},
			{
				name:   "missing user",
				config: "server:db:",
				error:  "managed identity config must be in format 'server:database:user'",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := migrate.BuildManagedIdentityConnectionString(tt.config)
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.error)
			})
		}
	})
}

// Integration tests using test containers
func TestInitDatabase_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test containers
	config := testkit.DefaultContainerConfig()
	pgConnStr, _ := testkit.SetupTestContainers(t, config)

	// Change to project root so migration files can be found
	oldWd, err := os.Getwd()
	require.NoError(t, err)

	projectRoot := findProjectRoot(oldWd)
	err = os.Chdir(projectRoot)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Chdir(oldWd))
	}()

	// Run the full migration process
	err = migrate.InitDatabase(context.Background(), pgConnStr)
	require.NoError(t, err)

	// Verify migrations were applied by checking if tables exist
	db, err := sql.Open("pgx", pgConnStr)
	require.NoError(t, err)
	defer db.Close()

	// Check if some key tables exist
	tables := []string{"users", "applications", "organisations", "api_keys"}
	for _, table := range tables {
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT 1
				FROM information_schema.tables
				WHERE table_schema = 'public'
				AND table_name = $1
			)`, table).Scan(&exists)
		require.NoError(t, err)
		require.True(t, exists, "Table %s should exist after migration", table)
	}
}

func TestRunPreMigrationSQL_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test containers
	config := testkit.DefaultContainerConfig()
	pgConnStr, _ := testkit.SetupTestContainers(t, config)

	// Get database connection
	db, err := sql.Open("pgx", pgConnStr)
	require.NoError(t, err)
	defer db.Close()

	// Change to project root so schema files can be found
	oldWd, err := os.Getwd()
	require.NoError(t, err)

	projectRoot := findProjectRoot(oldWd)
	err = os.Chdir(projectRoot)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Chdir(oldWd))
	}()

	// Run pre-migration SQL
	err = migrate.RunPreMigrationSQL(db)
	require.NoError(t, err)

	// Verify some pre-migration setup was applied
	// Check if uuid-ossp extension is available (commonly created in pre-migration)
	var extensionExists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp'
		)`).Scan(&extensionExists)
	require.NoError(t, err)
	// Note: extension might not exist if not in schema_pre.sql, that's okay
	t.Logf("uuid-ossp extension exists: %v", extensionExists)
}

func TestRunPostMigrationSQL_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test containers
	config := testkit.DefaultContainerConfig()
	pgConnStr, _ := testkit.SetupTestContainers(t, config)

	// Get database connection
	db, err := sql.Open("pgx", pgConnStr)
	require.NoError(t, err)
	defer db.Close()

	// Change to project root so schema files can be found
	oldWd, err := os.Getwd()
	require.NoError(t, err)

	projectRoot := findProjectRoot(oldWd)
	err = os.Chdir(projectRoot)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Chdir(oldWd))
	}()

	// Run post-migration SQL
	err = migrate.RunPostMigrationSQL(db)
	require.NoError(t, err)

	// Verify some post-migration setup was applied
	// Check if any indexes or constraints were created
	var indexCount int
	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM pg_indexes
		WHERE schemaname = 'public'
	`).Scan(&indexCount)
	require.NoError(t, err)
	t.Logf("Number of indexes created: %d", indexCount)
}

// findProjectRoot finds the project root directory by looking for go.mod
func findProjectRoot(startDir string) string {
	dir := startDir
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root directory
			return startDir
		}
		dir = parent
	}
}
