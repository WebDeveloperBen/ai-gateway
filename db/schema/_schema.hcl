// Schema definition - defines the public schema
schema "public" {
  comment = "standard public schema"
}

// Note: PostgreSQL extensions, functions and triggers are managed through
// separate SQL files (schema_pre.sql and schema_post.sql) as these features
// require Atlas Pro when defined in HCL