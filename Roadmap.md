# Developer Stories: Outstanding Features

## Data Plane (Proxy)
### Implement Redis-backed rate limiting for API keys with 429 response handling
- [ ] Add Redis token and request bucket infrastructure
- [ ] Enforce per-key limits on proxy requests
- [ ] Respond with 429 and Retry-After header on limit exceeded
- [ ] Integrate with key store for config-driven policy

### Support Proxy Keys for Apps and Users With Resource-constrained Rate Limits
- [ ] Design key provisioning system for per-app and per-user proxy keys
- [ ] UI/API to create, assign, and manage proxy keys scoped to apps/users
- [ ] Allow setting tpm/rate limits as a fraction/portion of a parent resource quota per key
- [ ] Enforce per-key and per-parent resource total accounting in the ratelimiter logic
- [ ] Allow real-time audit and usage tracking per key and aggregate at resource level
- [ ] Documentation for best practices and admin workflows

### Track and export per-request usage metrics, including tokens, cost, latency, status
- [ ] Instrument request pipeline for token usage, cost, and latency
- [ ] Aggregate and persist usage metrics to Redis/Cosmos
- [ ] Implement periodic rollups to DB for long-term reporting
- [ ] Support CSV/JSON export endpoints

### Enforce model/endpoint allow-deny policy and context/token caps on all proxy traffic
- [ ] Design per-key/model (and org) policy definitions
- [ ] Middleware to check/deny requests per policy
- [ ] UI/API for creating/updating policies
- [ ] Tests for allowed/denied model scenarios

### Expand OpenAI endpoint support to cover all v1 endpoints: completions, embeddings, files, fine-tunes, etc.
- [ ] Map and implement missing endpoints
- [ ] Ensure provider adapters route/transform as required
- [ ] Add e2e and conformance tests for all endpoints

### Integrate OpenTelemetry (tracing, RPS/timing/tokens/cost/error metrics) into proxy requests
- [ ] Instrument traces for each API call
- [ ] Emit key metrics (errors, rps, tokens, costs)
- [ ] Add Grafana/Prometheus/OpenTelemetry backend docs/support

### Add Entra ID JWT authentication for proxy API users (optional, phase 2)
- [ ] Add Entra ID token JWT verification to proxy
- [ ] Map roles/permissions from claims
- [ ] Support mixed API-key + JWT environments

### Switch AOAI upstream authentication to Managed Identity
- [ ] Integrate managed identity token flow for Azure OpenAI
- [ ] Fallback to API keys if MI not configured
- [ ] Rotate identity tokens on TTL expiry

### Apply request/response log redaction and support configurable log retention
- [ ] Add configurable redaction policy for logs
- [ ] Implement retention policy and automated cleanup
- [ ] Ensure admin/audit logs follow retention rules

### Dynamic Model Catalog and API Route Hot Reload
- [ ] Use Postgres as source of truth for model/route registrations
- [ ] Design schema and API endpoints for CRUD of model catalog entries
- [ ] Ensure admin UI supports live registration/updating of models and aliases
- [ ] Sync updated catalog to Redis for fast access by proxy
- [ ] Implement hot-reloading logic: proxy listens for changes and (re)registers API routes on-the-fly
- [ ] Guarantee zero downtime/connection dropping when updating routes
- [ ] Add validation and conflict detection for overlapping or malformed routes

## Admin/API & UI
### Build admin endpoints for management of routing rules, model aliases, and usage policies
- [ ] CRUD endpoints for routing rules/model aliases
- [ ] Persist policy/routing meta to DB
- [ ] Integrate API with Nuxt admin frontend

### Implement dashboards for real-time and historical usage (tokens, cost, latency, error rates)
- [ ] API endpoints for aggregated usage data
- [ ] UI (charts/tables) for per-app/team/client metrics
- [ ] Export capabilities (CSV/JSON)

### UI for audit log exports showing user, time, and performed action
- [ ] Structured export endpoint for audit logs
- [ ] UI for filtered search/downloads by user/date/action

### Support key provisioning, revocation, expiry, and rotation from Nuxt admin
- [ ] Build UI for key lifecycle actions
- [ ] Backend endpoints for key operations/metadata

### Surface model routing: allow admin to map models to AOAI/OpenAI deployments
- [ ] UI and API to assign models to deployments
- [ ] Validation and override logic for model aliasing

### Integrate role assignment and permission management via Entra ID
- [ ] Entra ID role mapping CRUD in admin
- [ ] API enforcement for admin-facing endpoints

### Enable app/team/client usage & cost CSV export from the admin UI
- [ ] Add bulk export endpoints and UI flows

## Platform, Security, & Infra
### Implement central pluggable SecretResolver to resolve secret_ref formats at runtime
- [ ] Design and implement the SecretResolver
- [ ] Support env, kv://, db:// lookups with TTL caching
- [ ] Replace direct secret usage throughout backend

### Ensure no plaintext upstream secrets ever persist in DB; audit for secret_ref practices
- [ ] Add secret persistence audit/scripts
- [ ] Migrate all stored secrets to secret_ref

### Complete VNET, PrivateLink, and private DNS Azure integration
- [ ] Add Terraform/infra for end-to-end private connectivity
- [ ] Update deployment runbooks/docs

### Maintain or add Terraform modules provisioning the entire stack for new tenants
- [ ] Modularize existing infrastructure Terraform
- [ ] Document for tenant onboarding/extension

## Cost Controls & Analytics
### Enforce app/team budgets—with 80%/100% alerting & automatic throttling
- [ ] Add per-app/team budget config
- [ ] Implement 80%/100% spend triggers/alerts
- [ ] Trigger traffic throttling/downshift when exceeded

### Convert token usage to dollars/euros/per-app costs and roll up by team/client
- [ ] Token→cost conversion for each call
- [ ] Aggregation/rollup by app/team/client

### Create SIEM/PowerBI/EventHub/FinOps exports for usage, cost, error data
- [ ] Build export pipeline for reporting
- [ ] Doc integration with external analytics suites

## SLOs & Operational Resilience
### Instrument and test for <20ms P95 added proxy latency at 100 RPS
- [ ] Add request/response timing/tracing
- [ ] Load test for compliance
- [ ] Report SLI/SLOs in dashboards

### Deploy multi-zone Redis to achieve ≥99.9% uptime SLO
- [ ] Configure multi-zone Redis for prod
- [ ] Alerting/monitoring for failover scenarios

### (Optional/Future) Add WAF/CDN protection, explore multi-cloud deployments
- [ ] Infra code and runbooks for WAF/CDN
- [ ] Draft architecture for non-Azure support
