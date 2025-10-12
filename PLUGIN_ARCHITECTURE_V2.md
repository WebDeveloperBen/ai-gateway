# Policy Plugin Architecture v2 - Multi-Tenant Design

## Overview

A hybrid policy system that supports:
1. **Built-in policies** - System-defined, registered via `init()` (rate limit, token limit, etc.)
2. **Custom CEL policies** - Customer-defined via UI, stored in database per-app
3. **Multi-tenancy** - Org/App isolation with efficient caching

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Policy Resolution Flow                        │
└─────────────────────────────────────────────────────────────────┘

Request for app-123
        │
        ▼
┌───────────────────┐
│ Memory Cache (30s)│ ◄── Hit: Return compiled policies
│  app-123 -> [P1]  │
└─────────┬─────────┘
          │ Miss
          ▼
┌───────────────────┐
│ Redis Cache (5m)  │ ◄── Hit: Reconstruct policies
│  app-123 -> JSON  │
└─────────┬─────────┘
          │ Miss
          ▼
┌───────────────────┐
│  Database Query   │
│  Get policies for │
│     app-123       │
└─────────┬─────────┘
          │
          ▼
    ┌─────────┴──────────┐
    │                    │
    ▼                    ▼
┌─────────────┐    ┌──────────────────┐
│ Built-in    │    │ Custom CEL       │
│ Policies    │    │ Policies         │
│             │    │                  │
│ rate_limit  │    │ "no_pii_check"   │
│ token_limit │    │ "business_hours" │
└──────┬──────┘    └────────┬─────────┘
       │                    │
       └──────────┬─────────┘
                  ▼
        ┌──────────────────┐
        │  Policy Registry │ ◄── Lookup factory
        │  type -> factory │
        └──────────────────┘
                  │
                  ▼
        ┌──────────────────┐
        │ Compiled Policies│
        │   [P1, P2, P3]   │
        └──────────────────┘
                  │
                  ▼
        Cache & Execute
```

---

## Database Schema (Already Exists!)

```sql
-- Current schema: db/schema/policies.hcl
CREATE TABLE policies (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id       UUID NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    app_id       UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    policy_type  TEXT NOT NULL,        -- 'rate_limit', 'token_limit', 'custom_cel'
    config       JSONB NOT NULL DEFAULT '{}'::jsonb,
    enabled      BOOLEAN NOT NULL DEFAULT true,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes
CREATE INDEX idx_policies_app_enabled ON policies(app_id, enabled);
CREATE INDEX idx_policies_org ON policies(org_id);
CREATE INDEX idx_policies_type ON policies(policy_type);
```

**Key Points:**
- ✅ Already supports org_id + app_id (multi-tenant)
- ✅ `policy_type` field distinguishes built-in vs custom
- ✅ `config` JSONB stores policy configuration
- ✅ `enabled` flag for easy on/off switching
- ✅ Cascade deletes when org/app deleted

---

## Policy Types

### 1. Built-in Policies (System-Defined)

Registered via `init()` at startup. Available to all customers.

```go
// Policy types defined in model/policy.go
const (
    PolicyTypeRateLimit      PolicyType = "rate_limit"
    PolicyTypeTokenLimit     PolicyType = "token_limit"
    PolicyTypeModelAllowlist PolicyType = "model_allowlist"
    PolicyTypeRequestSize    PolicyType = "request_size"
    // Future: PolicyTypeContentFilter, PolicyTypeCostLimit, etc.
)
```

**Database Example:**
```json
{
    "org_id": "org-abc",
    "app_id": "app-123",
    "policy_type": "rate_limit",
    "config": {
        "requests_per_minute": 1000,
        "burst": 50
    },
    "enabled": true
}
```

### 2. Custom CEL Policies (Customer-Defined)

Created via Admin UI. Stored as `custom_cel` type with unique names.

```go
const (
    PolicyTypeCustomCEL PolicyType = "custom_cel"
)
```

**Database Example:**
```json
{
    "org_id": "org-abc",
    "app_id": "app-123",
    "policy_type": "custom_cel",
    "config": {
        "name": "block_weekends",
        "description": "Block requests on weekends",
        "pre_check_expression": "timestamp.getDayOfWeek(request.time) < 6"
    },
    "enabled": true
}
```

---

## Implementation Plan

### Phase 1: Policy Registry (1 hour)

**Create:** `internal/gateway/policies/registry.go`

```go
package policies

import (
    "fmt"
    "sync"
    "github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
    "github.com/WebDeveloperBen/ai-gateway/internal/model"
)

// PolicyFactory creates a policy from config and dependencies
type PolicyFactory func(config []byte, deps PolicyDependencies) (Policy, error)

// PolicyDependencies holds shared dependencies
type PolicyDependencies struct {
    Cache kv.KvStore
}

var (
    registry   = make(map[model.PolicyType]PolicyFactory)
    registryMu sync.RWMutex
)

// Register adds a policy factory (called from init())
func Register(policyType model.PolicyType, factory PolicyFactory) {
    registryMu.Lock()
    defer registryMu.Unlock()
    
    if _, exists := registry[policyType]; exists {
        panic(fmt.Sprintf("policy type %s already registered", policyType))
    }
    
    registry[policyType] = factory
}

// GetFactory retrieves a registered factory
func GetFactory(policyType model.PolicyType) (PolicyFactory, bool) {
    registryMu.RLock()
    defer registryMu.RUnlock()
    
    factory, exists := registry[policyType]
    return factory, exists
}

// ListBuiltInPolicies returns all registered policy types
func ListBuiltInPolicies() []model.PolicyType {
    registryMu.RLock()
    defer registryMu.RUnlock()
    
    types := make([]model.PolicyType, 0, len(registry))
    for t := range registry {
        types = append(types, t)
    }
    return types
}
```

### Phase 2: Convert Existing Policies (2 hours)

**Update:** Each policy file adds `init()` function

**Example: `internal/gateway/policies/rate_limit.go`**

```go
package policies

import (
    "encoding/json"
    "fmt"
    "github.com/WebDeveloperBen/ai-gateway/internal/model"
)

func init() {
    Register(model.PolicyTypeRateLimit, func(config []byte, deps PolicyDependencies) (Policy, error) {
        var cfg model.RateLimitConfig
        if err := json.Unmarshal(config, &cfg); err != nil {
            return nil, fmt.Errorf("invalid rate limit config: %w", err)
        }
        return NewRateLimitPolicy(cfg, deps.Cache), nil
    })
}

// Rest of rate_limit.go stays the same...
```

**Repeat for:**
- `token_limit.go`
- `model_allowlist.go`
- `request_size.go`

### Phase 3: Update Engine (1 hour)

**Modify:** `internal/gateway/policies/engine.go`

```go
// NewPolicy creates a policy from type and config (registry-based)
func (e *Engine) NewPolicy(policyType model.PolicyType, config []byte) (Policy, error) {
    // Try registry first (built-in policies)
    factory, exists := GetFactory(policyType)
    if exists {
        deps := PolicyDependencies{Cache: e.cache}
        return factory(config, deps)
    }
    
    // Fallback: treat as custom CEL policy
    if policyType == model.PolicyTypeCustomCEL {
        return NewCELPolicy(policyType, config)
    }
    
    // Unknown type - this shouldn't happen with proper validation
    return nil, fmt.Errorf("unknown policy type: %s", policyType)
}

// LoadPolicies remains the same - already loads from DB correctly
```

**No changes needed to:**
- `LoadPolicies()` - Already queries database by app_id
- Caching logic - Already works with any policy type
- `CheckPreRequest()` - Policy-agnostic

### Phase 4: Admin API Endpoints (2-3 hours)

**Create:** `internal/api/admin/policies/` (new package)

#### 4.1 List Available Policy Types

```go
// GET /api/v1/admin/policies/types
type PolicyTypeInfo struct {
    Type        string `json:"type"`
    Name        string `json:"name"`
    Description string `json:"description"`
    ConfigSchema string `json:"config_schema"` // JSON Schema
    IsBuiltIn   bool   `json:"is_built_in"`
}

func (s *Service) ListPolicyTypes(ctx context.Context) ([]PolicyTypeInfo, error) {
    types := []PolicyTypeInfo{}
    
    // Built-in policies
    for _, policyType := range policies.ListBuiltInPolicies() {
        types = append(types, PolicyTypeInfo{
            Type:        string(policyType),
            Name:        formatPolicyName(policyType),
            Description: getPolicyDescription(policyType),
            ConfigSchema: getPolicySchema(policyType),
            IsBuiltIn:   true,
        })
    }
    
    // Custom CEL
    types = append(types, PolicyTypeInfo{
        Type:        "custom_cel",
        Name:        "Custom CEL Policy",
        Description: "Define custom policy logic using CEL expressions",
        ConfigSchema: celConfigSchema,
        IsBuiltIn:   false,
    })
    
    return types, nil
}
```

#### 4.2 Create Policy

```go
// POST /api/v1/admin/applications/{app_id}/policies
type CreatePolicyRequest struct {
    PolicyType string          `json:"policy_type"`
    Config     json.RawMessage `json:"config"`
    Enabled    bool            `json:"enabled"`
}

func (s *Service) CreatePolicy(ctx context.Context, appID string, req CreatePolicyRequest) (*model.Policy, error) {
    // Validate app belongs to org
    orgID := middleware.GetOrgID(ctx)
    app, err := s.appRepo.GetByID(ctx, appID)
    if err != nil || app.OrgID != orgID {
        return nil, ErrUnauthorized
    }
    
    // Validate policy type exists (built-in or custom_cel)
    if !s.isValidPolicyType(req.PolicyType) {
        return nil, fmt.Errorf("invalid policy type: %s", req.PolicyType)
    }
    
    // Validate config by attempting to create policy
    policyType := model.PolicyType(req.PolicyType)
    _, err = s.policyEngine.NewPolicy(policyType, req.Config)
    if err != nil {
        return nil, fmt.Errorf("invalid policy config: %w", err)
    }
    
    // Insert into database
    policy := &model.Policy{
        OrgID:      orgID,
        AppID:      appID,
        PolicyType: req.PolicyType,
        Config:     req.Config,
        Enabled:    req.Enabled,
    }
    
    if err := s.policyRepo.Create(ctx, policy); err != nil {
        return nil, err
    }
    
    // Invalidate cache for this app
    s.policyEngine.InvalidateCache(ctx, appID)
    
    return policy, nil
}
```

#### 4.3 List Policies for App

```go
// GET /api/v1/admin/applications/{app_id}/policies
func (s *Service) ListPolicies(ctx context.Context, appID string) ([]model.Policy, error) {
    orgID := middleware.GetOrgID(ctx)
    
    // Verify app belongs to org
    app, err := s.appRepo.GetByID(ctx, appID)
    if err != nil || app.OrgID != orgID {
        return nil, ErrUnauthorized
    }
    
    return s.policyRepo.ListByApp(ctx, appID)
}
```

#### 4.4 Update Policy

```go
// PATCH /api/v1/admin/policies/{policy_id}
type UpdatePolicyRequest struct {
    Config  *json.RawMessage `json:"config,omitempty"`
    Enabled *bool            `json:"enabled,omitempty"`
}

func (s *Service) UpdatePolicy(ctx context.Context, policyID string, req UpdatePolicyRequest) error {
    orgID := middleware.GetOrgID(ctx)
    
    policy, err := s.policyRepo.GetByID(ctx, policyID)
    if err != nil || policy.OrgID != orgID {
        return ErrUnauthorized
    }
    
    if req.Config != nil {
        // Validate new config
        _, err := s.policyEngine.NewPolicy(model.PolicyType(policy.PolicyType), *req.Config)
        if err != nil {
            return fmt.Errorf("invalid config: %w", err)
        }
        policy.Config = *req.Config
    }
    
    if req.Enabled != nil {
        policy.Enabled = *req.Enabled
    }
    
    if err := s.policyRepo.Update(ctx, policy); err != nil {
        return err
    }
    
    // Invalidate cache
    s.policyEngine.InvalidateCache(ctx, policy.AppID)
    
    return nil
}
```

#### 4.5 Delete Policy

```go
// DELETE /api/v1/admin/policies/{policy_id}
func (s *Service) DeletePolicy(ctx context.Context, policyID string) error {
    orgID := middleware.GetOrgID(ctx)
    
    policy, err := s.policyRepo.GetByID(ctx, policyID)
    if err != nil || policy.OrgID != orgID {
        return ErrUnauthorized
    }
    
    if err := s.policyRepo.Delete(ctx, policyID); err != nil {
        return err
    }
    
    // Invalidate cache
    s.policyEngine.InvalidateCache(ctx, policy.AppID)
    
    return nil
}
```

### Phase 5: Repository Layer (1 hour)

**Create:** `internal/repository/policies/postgres.go`

```go
package policies

import (
    "context"
    "github.com/WebDeveloperBen/ai-gateway/internal/db"
    "github.com/WebDeveloperBen/ai-gateway/internal/model"
    "github.com/google/uuid"
)

type PostgresRepo struct {
    queries *db.Queries
}

func NewPostgresRepo(queries *db.Queries) *PostgresRepo {
    return &PostgresRepo{queries: queries}
}

func (r *PostgresRepo) Create(ctx context.Context, policy *model.Policy) error {
    params := db.CreatePolicyParams{
        OrgID:      uuid.MustParse(policy.OrgID),
        AppID:      uuid.MustParse(policy.AppID),
        PolicyType: policy.PolicyType,
        Config:     policy.Config,
        Enabled:    policy.Enabled,
    }
    
    result, err := r.queries.CreatePolicy(ctx, params)
    if err != nil {
        return err
    }
    
    policy.ID = result.ID.String()
    policy.CreatedAt = result.CreatedAt
    policy.UpdatedAt = result.UpdatedAt
    return nil
}

func (r *PostgresRepo) ListByApp(ctx context.Context, appID string) ([]model.Policy, error) {
    rows, err := r.queries.ListPoliciesByApp(ctx, uuid.MustParse(appID))
    if err != nil {
        return nil, err
    }
    
    policies := make([]model.Policy, len(rows))
    for i, row := range rows {
        policies[i] = model.Policy{
            ID:         row.ID.String(),
            OrgID:      row.OrgID.String(),
            AppID:      row.AppID.String(),
            PolicyType: row.PolicyType,
            Config:     row.Config,
            Enabled:    row.Enabled,
            CreatedAt:  row.CreatedAt,
            UpdatedAt:  row.UpdatedAt,
        }
    }
    
    return policies, nil
}

// GetByID, Update, Delete methods...
```

### Phase 6: SQLC Queries (30 min)

**Add to:** `db/queries/policies.sql`

```sql
-- name: CreatePolicy :one
INSERT INTO policies (org_id, app_id, policy_type, config, enabled)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListPoliciesByApp :many
SELECT * FROM policies
WHERE app_id = $1
ORDER BY created_at ASC;

-- name: GetPolicyByID :one
SELECT * FROM policies
WHERE id = $1;

-- name: UpdatePolicy :exec
UPDATE policies
SET config = $2, enabled = $3, updated_at = now()
WHERE id = $1;

-- name: DeletePolicy :exec
DELETE FROM policies
WHERE id = $1;

-- name: ListEnabledPolicies :many (ALREADY EXISTS)
SELECT * FROM policies
WHERE app_id = $1 AND enabled = true
ORDER BY created_at ASC;
```

---

## Multi-Tenancy & Security

### Org/App Isolation

```go
// Every API call validates org ownership
func (s *Service) CreatePolicy(ctx context.Context, appID string, req CreatePolicyRequest) error {
    orgID := middleware.GetOrgID(ctx)  // From JWT/session
    
    // Verify app belongs to this org
    app, err := s.appRepo.GetByID(ctx, appID)
    if err != nil {
        return ErrNotFound
    }
    
    if app.OrgID != orgID {
        return ErrUnauthorized  // Org trying to access another org's app
    }
    
    // Proceed with policy creation...
}
```

### Row-Level Security (Optional)

```sql
-- Add RLS policy to policies table
ALTER TABLE policies ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON policies
    USING (org_id = current_setting('app.current_org_id')::uuid);
```

---

## UI Flow Examples

### 1. Admin Creates Rate Limit Policy

```
User: Admin at Org "Acme Corp"
App: "Production API" (app-123)

1. Navigate to: /apps/app-123/policies
2. Click "Add Policy"
3. Select from dropdown: "Rate Limit"
4. Form shows schema-based UI:
   - Requests per minute: [1000]
   - Burst: [50]
5. Click "Create"

→ POST /api/v1/admin/applications/app-123/policies
{
    "policy_type": "rate_limit",
    "config": {
        "requests_per_minute": 1000,
        "burst": 50
    },
    "enabled": true
}

→ Policy inserted into DB with org_id="acme-corp"
→ Cache invalidated for app-123
→ Next request to app-123 loads new policy
```

### 2. Customer Creates Custom CEL Policy

```
User: Admin at Org "Acme Corp"
App: "Production API" (app-123)

1. Navigate to: /apps/app-123/policies
2. Click "Add Policy"
3. Select: "Custom CEL Policy"
4. Form shows:
   - Name: [block_large_requests]
   - Description: [Block requests over 100KB]
   - Pre-check expression:
     [request.size_bytes < 102400]
   - Post-check expression: [leave empty]
5. Click "Validate" (tests CEL compilation)
6. Click "Create"

→ POST /api/v1/admin/applications/app-123/policies
{
    "policy_type": "custom_cel",
    "config": {
        "name": "block_large_requests",
        "description": "Block requests over 100KB",
        "pre_check_expression": "request.size_bytes < 102400"
    },
    "enabled": true
}

→ Backend validates CEL expression compiles
→ Policy created and cached
→ Requests over 100KB now blocked for app-123
```

---

## Caching Strategy

### Three-Tier Cache (Already Implemented!)

```
Request for app-123
    │
    ├─► Memory (30s TTL)  ✓ Policy objects compiled and ready
    │
    ├─► Redis (5m TTL)    ✓ JSON serialized, must reconstruct
    │
    └─► Database          ✓ Full query, compile from scratch
```

**Cache Invalidation:**
- On policy create/update/delete → invalidate app-specific cache
- On app delete → cascade deletes policies (database foreign key)
- Manual invalidation via admin API (if needed)

---

## Performance Considerations

### 1. Compilation Cost

**Built-in policies:** Instant (simple struct construction)
**CEL policies:** ~74μs to compile (from benchmarks)

**Mitigation:**
- Memory cache keeps compiled CEL programs (no recompilation)
- Redis cache only rebuilds on memory cache miss (~5% of requests)

### 2. Database Load

**Per-app policy load:** 1 query
```sql
SELECT * FROM policies WHERE app_id = ? AND enabled = true;
```

**Optimization:**
- Query only runs on cache miss (< 5% of requests)
- Index on `(app_id, enabled)` ensures fast lookup
- Typical result set: 3-10 policies per app

### 3. Memory Usage

**Per-app cache entry:**
- 3 policies × 500 bytes = 1.5KB
- 1000 concurrent apps × 1.5KB = 1.5MB

**LRU eviction:**
- Memory cache limited to 1000 entries
- Least-recently-used apps evicted first

---

## Migration Path

### Step 1: Implement Registry (Day 1, 3 hours)
- [ ] Create `registry.go`
- [ ] Add `init()` to existing policies
- [ ] Update `engine.NewPolicy()` to use registry
- [ ] Test all existing policies still work

### Step 2: Admin API (Day 1-2, 4 hours)
- [ ] Create policy repository
- [ ] Add SQLC queries
- [ ] Implement CRUD endpoints
- [ ] Add validation layer

### Step 3: UI Integration (Day 2-3, 6 hours)
- [ ] Policy list page per app
- [ ] Create policy form (dynamic based on type)
- [ ] Edit/delete policy actions
- [ ] CEL expression editor with validation

### Step 4: Testing (Day 3-4, 4 hours)
- [ ] Unit tests for registry
- [ ] Integration tests for CRUD
- [ ] E2E test: Create policy → invalidate cache → verify enforcement
- [ ] Multi-tenant isolation tests

**Total: 2-3 days**

---

## Benefits

1. **Extensibility** - Add new built-in policies without DB changes
2. **Flexibility** - Customers create custom CEL policies via UI
3. **Performance** - Three-tier caching, compiled policies
4. **Security** - Org/app isolation, validation before persistence
5. **Scalability** - Per-app policy loading, no global state
6. **Observability** - Already have policy check metrics
7. **Future-proof** - Can add policy marketplace, sharing, templates

---

## Next Steps

Ready to implement? I recommend this order:

1. **Start with registry** (low risk, high value)
2. **Test with existing policies** (ensure no regression)
3. **Build admin API** (backend foundation)
4. **Create UI forms** (customer-facing)
5. **Integration testing** (validate multi-tenancy)

Estimated: **2-3 days** for full implementation.

Should I start implementing Phase 1 (Policy Registry)?
