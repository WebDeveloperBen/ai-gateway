// Define an environment named "local"
env "local" {
  // Declare where the schema definition resides.
  // Load all .hcl files from the schema directory
  src = "file://db/schema"

  // Define the URL of the database which is managed in
  // this environment.
  url = getenv("POSTGRES_DSN")

  // Use a persistent dev database (postgres-dev from docker compose)
  // This database should have schema_pre.sql applied to it
  // Run: docker compose up -d postgres-dev && psql <dev-url> -f db/schema_pre.sql
  dev = "postgres://atlas:atlas@localhost:5433/atlas_dev?sslmode=disable"

  // Diff policy to ignore triggers and policies (managed in schema_post.sql)
  diff {
    skip {
      drop_schema = true
    }
  }

  migration {
    // URL where the migration directory resides.
    dir = "file://db/migrations"
    // Format of the migration directory: atlas, flyway, goose, dbmate or golang-migrate.
    format = goose
    // No baseline needed - dev database is manually synced via db:dev:init
  }
}

env "sqlc" {
  src = "file://db/schema"

  // For SQLC integration, we don't need a live database
  // Just generate the schema files
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
