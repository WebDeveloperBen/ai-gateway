# db Package

This package manages database schema migrations and type-safe access to the database.

- **Migrations:** All database migrations live in this directory. Use [`goose`](https://github.com/pressly/goose) to manage and apply schema migrations. Never modify migration files after they have been applied—always create a new migration.

- **Queries:** We use [`sqlc`](https://docs.sqlc.dev/) to generate Go code from .sql files. To update database queries, edit the SQL files in this package and re-run `sqlc generate`.

**Development flow:**
1. Create/modify .sql migration files and queries.
2. Use `goose` for migrations: `goose up`, `goose create migration_name sql`.
3. Use `sqlc generate` to update generated query code.

Do not hand-write Go SQL logic—always use SQL files and codegen.

See project docs for the exact workflow and any troubleshooting.
