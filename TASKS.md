# Policy System Implementation Tasks

## Completed ✅

### Phase 1: Database & Models
- [x] Create HCL schema files for new tables (applications, application_configs, models, policies, usage_metrics)
- [x] Generate and apply database migrations
- [x] Add Go models to internal/model/
- [x] Write SQLC queries in db/queries/
- [x] Run task generate to create Go code

### Phase 2: Policy Engine Core
- [x] Create policy engine interface and types
- [x] Create policy context structs (PreRequestContext, PostRequestContext)
- [x] Create policy factory
- [x] Create Redis caching layer for policies
- [x] Implement concrete policy types (rate_limit, token_limit, model_allowlist, request_size)
- [x] Add CEL library dependency
- [x] Create CEL policy evaluator for custom policies
- [x] Update policy types to support custom CEL policies
- [x] Update factory to handle CEL-based policies

### Phase 3: Token Estimation
- [x] Add tiktoken-go library
- [x] Create token estimator for different models
- [x] Create response parser for OpenAI/Azure OpenAI format
- [x] Support multiple provider response formats (OpenAI, Azure OpenAI, Anthropic, Cohere, Google)
- [x] Refactor parser to use provider-specific interface pattern

### Phase 4: Transport Middleware
- [x] Create WithPolicyEnforcement middleware
- [x] Create WithUsageRecording middleware with async pattern
- [x] Implement detached context for goroutines
- [x] Wire up policy engine in cmd/proxy/main.go

## In Progress 🚧

### Phase 5: Repository Layer
- [ ] Create applications repository (CRUD)
- [ ] Create application_configs repository
- [ ] Create models repository (for model deployments)
- [ ] Create policies repository with cache invalidation
- [ ] Create usage repository with rollup queries

### Phase 6: Admin API Endpoints
- [ ] Create admin/applications endpoints
- [ ] Create admin/application_configs endpoints
- [ ] Create admin/models endpoints
- [ ] Create admin/policies endpoints
- [ ] Create admin/usage endpoints (dashboard queries)

### Phase 7: UI Integration
- [ ] Update UI for application management
- [ ] Add model deployment management
- [ ] Add policy builder with CEL editor
- [ ] Add usage dashboard with metrics

## Architecture Notes

### Policy System Design
- **Predefined Policies**: rate_limit, token_limit, model_allowlist, request_size
- **Custom CEL Policies**: Admins can create custom policies using CEL expressions
- **Cache-Through Pattern**: Policies cached in Redis (5min TTL) for hot-path performance
- **Pre/Post Hooks**: PreCheck blocks requests, PostCheck records metrics asynchronously

### CEL Expression Variables
**Pre-Request:**
- `request_size_bytes` (int)
- `estimated_tokens` (int)
- `model` (string)
- `org_id` (string)
- `app_id` (string)

**Post-Request:**
- `prompt_tokens` (int)
- `completion_tokens` (int)
- `total_tokens` (int)
- `latency_ms` (int)
- `response_size_bytes` (int)

### Example CEL Policy
```json
{
  "policy_type": "custom_cel",
  "config": {
    "pre_check_expression": "estimated_tokens < 4000 && request_size_bytes < 100000",
    "post_check_expression": "total_tokens < 8000"
  }
}
```

## Database Schema
- `applications` - App identity (no budgets)
- `application_configs` - Per-environment config
- `models` - Model deployments provisioned by platform team
- `policies` - Policy definitions (predefined or CEL)
- `usage_metrics` - Token tracking for policy enforcement

## Next Steps
1. Add tiktoken-go for token estimation
2. Implement transport middleware
3. Create repository layer
4. Build Admin API endpoints
