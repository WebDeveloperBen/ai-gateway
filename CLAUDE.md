# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an **LLM Gateway/Proxy Platform** that intermediates application traffic to Azure OpenAI and OpenAI providers. It provides OpenAI-compatible endpoints while centralizing security, governance, cost controls, and observability.

The system consists of:
- **Go Backend (Data Plane)**: Proxy service with authentication, rate limiting, and request forwarding
- **Nuxt UI (Admin Plane)**: Admin interface for key management, routing rules, and dashboards
- **Database Layer**: PostgreSQL with SQLC-generated queries

## Development Commands

### Backend (Go)
```bash
# Run with hot reload (requires Docker for dependencies)
task run

# Build binary (outputs to ./bin/api)
task build

# Run all tests
task tests

# Run specific test
go test -v ./internal/api/auth/...

# Run tests in a single package
go test -v ./internal/gateway/

# Database migrations
task db:up          # Apply all migrations
task db:down        # Rollback last migration
task db:status      # Show migration status
task db:reset       # Reset database completely
task db:create <name> # Create new migration

# Generate SQLC code from queries
task generate

# Stop development environment
task stop
```

### Frontend (UI)
```bash
cd ui/

# Development server
pnpm dev

# Build for production
pnpm build

# Generate static site
pnpm generate

# Preview production build
pnpm preview

# Format code
pnpm format
```

## Architecture

### Go Backend Structure
- **`cmd/proxy/main.go`**: Main application entry point
- **`internal/api/`**: HTTP API routes and handlers
  - `admin/`: Admin API endpoints (keys, users, roles)
  - `auth/`: OIDC authentication routes
  - `gateway/`: Proxy endpoint handlers
  - `middleware/`: Authentication and request processing middleware
- **`internal/gateway/`**: Core proxy logic
  - `auth/`: API key authentication
  - `policies/`: Request policies and rate limiting
  - `loadbalancing/`: Load balancing strategies
  - Registry pattern for model deployments
- **`internal/repository/`**: Data access layer with interface abstractions
- **`internal/drivers/`**: Database (PostgreSQL) and KV store (Redis/Memory) drivers
- **`internal/provider/`**: Provider adapters (Azure OpenAI, OpenAI)
- **`internal/db/`**: SQLC-generated database code
- **`sqlc/`**: Database migrations and queries

### Frontend Structure
- **Nuxt 4** with TypeScript
- **TailwindCSS** for styling with **Reka UI** components
- **Vue Query** for API state management
- **Vee-validate + Zod** for form validation
- Pages in `ui/app/pages/(dashboard)/` use file-based routing

### Key Dependencies
- **Huma v2**: API framework with OpenAPI generation
- **Chi**: HTTP router
- **PGX v5**: PostgreSQL driver
- **Redis**: Rate limiting and caching
- **SQLC**: Type-safe SQL code generation
- **Zerolog**: Structured logging

## Database

### Schema Management
- Uses **Goose** for migrations in `sqlc/migrations/`
- **SQLC** generates type-safe Go code from SQL queries
- Database entities: Users, Organizations, Roles, API Keys, Model Deployments

### Common Patterns
- All repositories implement interfaces in `internal/model/`
- Use `internal/exceptions/` for standardized error handling
- Transactions handled via repository methods
- UUID primary keys with ULID generation

## Configuration

### Environment Variables
Set in `.env` file (see existing `.env` for reference):
- Database: `POSTGRES_DSN`
- Redis: `REDIS_ADDR`, `REDIS_PW` 
- Auth: `APP_REGISTRATION_*` variables for Azure AD
- Proxy: `PROXY_PORT`

### Build Configuration
- **Air**: Hot reload configuration in `.air.toml`
- **TaskFile**: Task runner configuration in `TaskFile.yml`
- **Docker**: Uses `compose.yml` for local dependencies

## Provider Integration

### Adding New Providers
1. Implement `internal/provider/Provider` interface
2. Add provider-specific configuration in models
3. Register in `cmd/proxy/main.go`
4. Add routes in `internal/api/gateway/`

### Existing Providers
- **Azure OpenAI**: Handles model mapping to Azure deployments
- **OpenAI**: Direct OpenAI API integration

## Security

- API keys stored hashed using Argon2ID
- OIDC integration with Azure AD for admin UI
- Secret references pattern (avoid storing raw secrets)
- Request/response logging with configurable redaction

## Testing

- Test files alongside source code (`*_test.go`)
- Use `internal/testkit/` for test utilities
- Run with `task tests` or `go test ./...`

## Development Workflow

1. Start dependencies: `task run` (starts Docker services + hot reload)
2. Database setup: `task db:up` (run migrations)
3. UI development: `cd ui && pnpm dev`
4. Generate code: `task generate` (after schema/query changes)
5. Tests: `task tests`

## Vue/Nuxt Component Conventions

### Component Structure
- **Script tag placement**: Always place `<script setup>` tag BEFORE `<template>` tag in Vue components
- **Component order**: Use the following structure:
  1. `<script setup lang="ts">` 
  2. `<template>`
  3. `<style>` (if needed)

### Code Style
- **Component naming**: Use PascalCase for component files and imports
- **Prop definitions**: Use TypeScript interfaces for component props with JSDoc documentation
- **Template simplicity**: Prefer direct property access in templates over helper functions when possible
- **DRY Components**: Create reusable components for common patterns (e.g., `ChartPlaceholder`, `RecentActivity`)
- **DO NOT ADD COMMENTS** unless explicitly requested