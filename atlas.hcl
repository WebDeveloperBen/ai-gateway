// Define an environment named "local"
env "local" {
  // Declare where the schema definition resides.
  // Load all .hcl files from the schema directory
  src = "file://db/schema"

  // Define the URL of the database which is managed in
  // this environment.
  url = getenv("ATLAS_DATABASE_URL")

  // Define the URL of the Dev Database for this environment
  // See: https://atlasgo.io/concepts/dev-database
  // Using docker ephemeral database with extensions pre-loaded
  dev = "docker://postgres/16/dev"

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
    // Baseline version
    baseline = "20250929133943"
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