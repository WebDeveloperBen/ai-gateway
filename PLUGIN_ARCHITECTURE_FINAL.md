# Policy Plugin Architecture - Final Design

## Core Concept

**All policies are opt-in per application.**

1. **Built-in policies** - System defines the factory/implementation in code via `init()` registration
2. **Custom CEL policies** - Customer defines the logic via UI/API
3. **Both stored in database** - `policies` table with `app_id` FK
4. **Both explicitly enabled** - Admin must attach policy to each app

**No automatic policy assignment. Zero policies by default.**

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                      Policy System Flow                          │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│  Code-Defined (init() registration)                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ rate_limit   │  │ token_limit  │  │ model_allow  │          │
│  │              │  │              │  │              │          │
│  │ func init()  │  │ func init()  │  │ func init()  │          │
│  │   Register() │  │   Register() │  │   Register() │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
│         │                 │                 │                   │
│         └─────────────────┼─────────────────┘                   │
│                           ▼                                      │
│                 ┌──────────────────┐                            │
│                 │  Policy Registry │                            │
│                 │                  │                            │
│                 │  type -> factory │                            │
│                 └──────────────────┘                            │
└─────────────────────────────────────────────────────────────────┘
                           │
                           │ Available for use
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│  Database (policies table)                                       │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ App: "Production API" (app-123, org-acme)                  │ │
│  │                                                            │ │
│  │  Policy 1: type="rate_limit"      config={rpm: 1000}     │ │
│  │  Policy 2: type="token_limit"     config={max: 8192}     │ │
│  │  Policy 3: type="custom_cel"      config={expr: "..."}   │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ App: "Dev API" (app-456, org-acme)                         │ │
│  │                                                            │ │
│  │  Policy 1: type="token_limit"     config={max: 4096}     │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ App: "Test API" (app-789, org-acme)                        │ │
│  │                                                            │ │
│  │  (No policies - empty!)                                   │ │
│  └────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                           │
                           │ Per-request query
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│  Runtime Policy Loading (per app)                               │
│                                                                  │
│  Request to app-123:                                            │
│    1. Query DB: WHERE app_id = 'app-123' AND enabled = true    │
│    2. Get [rate_limit, token_limit, custom_cel]                │
│    3. For each policy:                                          │
│       - If built-in: lookup factory from registry              │
│       - If custom_cel: use CEL policy constructor              │
│    4. Compile policies and cache                               │
│    5. Execute policy chain                                      │
└─────────────────────────────────────────────────────────────────┘
```

---

## Database Schema (Current)

```sql
CREATE TABLE policies (
    id           UUID PRIMARY KEY,
    org_id       UUID NOT NULL REFERENCES organisations(id),
    app_id       UUID NOT NULL REFERENCES applications(id),
    
    -- Policy type: 'rate_limit', 'token_limit', 'custom_cel', etc.
    policy_type  TEXT NOT NULL,
    
    -- Configuration (JSON varies by type)
    config       JSONB NOT NULL,
    
    -- Enable/disable per app
    enabled      BOOLEAN NOT NULL DEFAULT true,
    
    created_at   TIMESTAMPTZ NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL
);

-- Critical: Policy is tied to specific app
-- No policy exists without an app assignment
CREATE INDEX idx_policies_app_enabled ON policies(app_id, enabled);
```

**Key Points:**
- ✅ **No default policies** - apps start with zero policies
- ✅ **Explicit attachment** - admin must add each policy
- ✅ **Per-app configuration** - same policy type, different config per app
- ✅ **Built-in vs custom** - distinguished only by `policy_type` value

---

## Policy Types (Two Kinds)

### 1. Built-in Policies (System-Defined Factories)

Code defines **how to create** the policy. Database stores **which apps use it**.

```go
// internal/gateway/policies/rate_limit.go
func init() {
    // Register the factory for "rate_limit" type
    Register(model.PolicyTypeRateLimit, func(config []byte, deps PolicyDependencies) (Policy, error) {
        var cfg model.RateLimitConfig
        if err := json.Unmarshal(config, &cfg); err != nil {
            return nil, err
        }
        return NewRateLimitPolicy(cfg, deps.Cache), nil
    })
}
```

**Database Example:**
```json
{
    "app_id": "app-123",
    "policy_type": "rate_limit",  // ← References registered factory
    "config": {
        "requests_per_minute": 1000
    },
    "enabled": true
}
```

**Available Built-in Types:**
- `rate_limit` - Request rate limiting
- `token_limit` - Token usage limits
- `model_allowlist` - Restrict model access
- `request_size` - Request size limits
- Future: `content_filter`, `cost_limit`, `business_hours`, etc.

### 2. Custom CEL Policies (Customer-Defined Logic)

Customer defines **both creation and logic** via UI. System provides generic CEL executor.

```go
// internal/gateway/policies/cel_policy.go
func init() {
    // Register generic CEL policy executor
    Register(model.PolicyTypeCustomCEL, func(config []byte, deps PolicyDependencies) (Policy, error) {
        return NewCELPolicy(model.PolicyTypeCustomCEL, config)
    })
}
```

**Database Example:**
```json
{
    "app_id": "app-123",
    "policy_type": "custom_cel",  // ← Generic CEL executor
    "config": {
        "name": "block_weekends",
        "description": "No requests on weekends",
        "pre_check_expression": "timestamp.getDayOfWeek(request.time) < 6"
    },
    "enabled": true
}
```

---

## User Workflows

### Workflow 1: Admin Adds Built-in Policy

```
Context:
  Org: "Acme Corp"
  App: "Production API" (app-123)
  Current policies: [] (empty)

Step 1: Admin navigates to app policies
  → GET /api/v1/admin/applications/app-123
  → Shows: "No policies configured"

Step 2: Admin clicks "Add Policy"
  → GET /api/v1/admin/policies/available-types
  → Returns:
    [
      {type: "rate_limit", name: "Rate Limit", is_built_in: true},
      {type: "token_limit", name: "Token Limit", is_built_in: true},
      {type: "custom_cel", name: "Custom CEL", is_built_in: false}
    ]

Step 3: Admin selects "Rate Limit"
  → UI shows form based on schema:
    - Requests per minute: [____]
    - Burst: [____]

Step 4: Admin fills form and clicks "Create"
  → POST /api/v1/admin/applications/app-123/policies
    {
      "policy_type": "rate_limit",
      "config": {"requests_per_minute": 1000, "burst": 50},
      "enabled": true
    }

Step 5: Backend processes
  → Validates org owns app-123
  → Validates "rate_limit" exists in registry
  → Validates config by attempting to construct policy
  → Inserts into DB:
    INSERT INTO policies (org_id, app_id, policy_type, config, enabled)
    VALUES ('acme-corp', 'app-123', 'rate_limit', {...}, true)
  → Invalidates cache for app-123

Step 6: Next request to app-123
  → Loads policies from DB (cache miss)
  → Finds: [rate_limit policy]
  → Constructs via registry factory
  → Enforces rate limit
```

### Workflow 2: Admin Creates Custom CEL Policy

```
Context:
  Org: "Acme Corp"
  App: "Production API" (app-123)
  Current policies: [rate_limit]

Step 1: Admin clicks "Add Policy"
Step 2: Admin selects "Custom CEL Policy"
  → UI shows CEL editor:
    - Name: [____]
    - Description: [____]
    - Pre-check expression: [____]
    - Post-check expression: [____] (optional)

Step 3: Admin writes CEL expression
  Name: "block_large_prompts"
  Expression: "estimated_tokens < 5000"

Step 4: Admin clicks "Validate"
  → POST /api/v1/admin/policies/validate-cel
    {"expression": "estimated_tokens < 5000"}
  → Backend compiles CEL to verify syntax
  → Returns: {valid: true}

Step 5: Admin clicks "Create"
  → POST /api/v1/admin/applications/app-123/policies
    {
      "policy_type": "custom_cel",
      "config": {
        "name": "block_large_prompts",
        "description": "Limit prompt size",
        "pre_check_expression": "estimated_tokens < 5000"
      },
      "enabled": true
    }

Step 6: Backend processes
  → Validates org owns app-123
  → Validates CEL expression compiles
  → Inserts into DB
  → Invalidates cache for app-123

Step 7: Next request to app-123
  → Loads policies: [rate_limit, custom_cel]
  → Both enforced in order
```

### Workflow 3: Admin Reuses Policy on Another App

```
Context:
  Org: "Acme Corp"
  App 1: "Production API" (app-123) has rate_limit policy
  App 2: "Dev API" (app-456) has no policies

Step 1: Admin navigates to app-456 policies
Step 2: Admin clicks "Add Policy"
Step 3: Admin selects "Rate Limit" (same as app-123)
Step 4: Admin configures DIFFERENT values:
  - Requests per minute: 100 (vs 1000 for prod)
  - Burst: 10 (vs 50 for prod)

Result:
  → Two separate DB rows:
    policies:
      {id: 1, app_id: "app-123", type: "rate_limit", config: {rpm: 1000}}
      {id: 2, app_id: "app-456", type: "rate_limit", config: {rpm: 100}}
  
  → Same factory, different configs, different apps
```

---

## Implementation Components

### 1. Policy Registry (Already Designed)

```go
// internal/gateway/policies/registry.go
package policies

var registry = make(map[model.PolicyType]PolicyFactory)

func Register(policyType model.PolicyType, factory PolicyFactory) {
    registry[policyType] = factory
}

func GetFactory(policyType model.PolicyType) (PolicyFactory, bool) {
    factory, exists := registry[policyType]
    return factory, exists
}

// New: List available policy types for UI
func ListAvailableTypes() []PolicyTypeMetadata {
    var types []PolicyTypeMetadata
    for policyType := range registry {
        types = append(types, PolicyTypeMetadata{
            Type:         policyType,
            Name:         formatName(policyType),
            Description:  getDescription(policyType),
            ConfigSchema: getSchema(policyType),
            IsBuiltIn:    true,
        })
    }
    return types
}
```

### 2. Engine Policy Loading (Minimal Changes)

```go
// internal/gateway/policies/engine.go

// LoadPolicies - ALREADY DOES THE RIGHT THING!
func (e *Engine) LoadPolicies(ctx context.Context, appID string) ([]Policy, error) {
    // ... cache checking ...
    
    // Load from database (per app, not global)
    dbPolicies, err := e.db.ListEnabledPolicies(ctx, appUUID)
    if err != nil {
        return nil, err
    }
    
    // Convert each DB row to Policy instance
    policies := make([]Policy, 0, len(dbPolicies))
    for _, dbPolicy := range dbPolicies {
        policy, err := e.NewPolicy(model.PolicyType(dbPolicy.PolicyType), dbPolicy.Config)
        if err != nil {
            // Log but continue with other policies
            logger.GetLogger(ctx).Error().
                Err(err).
                Str("app_id", appID).
                Str("policy_type", dbPolicy.PolicyType).
                Msg("Failed to create policy from database")
            continue
        }
        policies = append(policies, policy)
    }
    
    return policies, nil
}

// NewPolicy - UPDATE TO USE REGISTRY
func (e *Engine) NewPolicy(policyType model.PolicyType, config []byte) (Policy, error) {
    // Try registry (built-in policies)
    factory, exists := GetFactory(policyType)
    if exists {
        deps := PolicyDependencies{Cache: e.cache}
        return factory(config, deps)
    }
    
    // Unknown type - shouldn't happen with validation
    return nil, fmt.Errorf("unknown policy type: %s", policyType)
}
```

**Key Point:** `LoadPolicies()` already queries per-app. No changes needed!

### 3. Admin API Endpoints

#### 3.1 List Available Policy Types

```go
// GET /api/v1/admin/policies/types
type PolicyTypeResponse struct {
    Type         string                 `json:"type"`
    Name         string                 `json:"name"`
    Description  string                 `json:"description"`
    ConfigSchema map[string]interface{} `json:"config_schema"` // JSON Schema
    IsBuiltIn    bool                   `json:"is_built_in"`
}

func (s *Service) ListAvailableTypes(ctx context.Context) ([]PolicyTypeResponse, error) {
    return policies.ListAvailableTypes(), nil
}

// Response:
[
  {
    "type": "rate_limit",
    "name": "Rate Limit",
    "description": "Limit requests per minute per application",
    "config_schema": {
      "type": "object",
      "properties": {
        "requests_per_minute": {"type": "integer", "minimum": 1},
        "burst": {"type": "integer", "minimum": 0}
      },
      "required": ["requests_per_minute"]
    },
    "is_built_in": true
  },
  {
    "type": "custom_cel",
    "name": "Custom CEL Policy",
    "description": "Define custom policy logic using CEL expressions",
    "config_schema": {
      "type": "object",
      "properties": {
        "name": {"type": "string"},
        "pre_check_expression": {"type": "string"}
      },
      "required": ["name", "pre_check_expression"]
    },
    "is_built_in": false
  }
]
```

#### 3.2 List Policies for App

```go
// GET /api/v1/admin/applications/{app_id}/policies
func (s *Service) ListAppPolicies(ctx context.Context, appID string) ([]PolicyResponse, error) {
    orgID := middleware.GetOrgID(ctx)
    
    // Verify app belongs to org
    app, err := s.appRepo.GetByID(ctx, appID)
    if err != nil || app.OrgID != orgID {
        return nil, ErrUnauthorized
    }
    
    // Get policies for this app
    dbPolicies, err := s.policyRepo.ListByApp(ctx, appID)
    if err != nil {
        return nil, err
    }
    
    // Convert to response
    policies := make([]PolicyResponse, len(dbPolicies))
    for i, p := range dbPolicies {
        policies[i] = PolicyResponse{
            ID:          p.ID,
            PolicyType:  p.PolicyType,
            Config:      p.Config,
            Enabled:     p.Enabled,
            CreatedAt:   p.CreatedAt,
            UpdatedAt:   p.UpdatedAt,
        }
    }
    
    return policies, nil
}

// Response:
[
  {
    "id": "policy-abc",
    "policy_type": "rate_limit",
    "config": {"requests_per_minute": 1000, "burst": 50},
    "enabled": true,
    "created_at": "2025-01-10T12:00:00Z",
    "updated_at": "2025-01-10T12:00:00Z"
  },
  {
    "id": "policy-def",
    "policy_type": "custom_cel",
    "config": {
      "name": "block_weekends",
      "pre_check_expression": "timestamp.getDayOfWeek(request.time) < 6"
    },
    "enabled": true,
    "created_at": "2025-01-11T10:30:00Z",
    "updated_at": "2025-01-11T10:30:00Z"
  }
]
```

#### 3.3 Add Policy to App

```go
// POST /api/v1/admin/applications/{app_id}/policies
type CreatePolicyRequest struct {
    PolicyType string          `json:"policy_type"`
    Config     json.RawMessage `json:"config"`
    Enabled    bool            `json:"enabled"`
}

func (s *Service) CreatePolicy(ctx context.Context, appID string, req CreatePolicyRequest) (*PolicyResponse, error) {
    orgID := middleware.GetOrgID(ctx)
    
    // 1. Verify app ownership
    app, err := s.appRepo.GetByID(ctx, appID)
    if err != nil || app.OrgID != orgID {
        return nil, ErrUnauthorized
    }
    
    // 2. Validate policy type exists (in registry or is custom_cel)
    policyType := model.PolicyType(req.PolicyType)
    if !s.isPolicyTypeValid(policyType) {
        return nil, fmt.Errorf("invalid policy type: %s", req.PolicyType)
    }
    
    // 3. Validate config by attempting to construct policy
    //    This ensures config is valid before persisting
    _, err = s.policyEngine.NewPolicy(policyType, req.Config)
    if err != nil {
        return nil, fmt.Errorf("invalid policy config: %w", err)
    }
    
    // 4. Insert into database
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
    
    // 5. Invalidate cache for this app
    s.policyEngine.InvalidateCache(ctx, appID)
    
    return toPolicyResponse(policy), nil
}

func (s *Service) isPolicyTypeValid(policyType model.PolicyType) bool {
    // Check if registered in code
    _, exists := policies.GetFactory(policyType)
    return exists
}
```

#### 3.4 Update Policy

```go
// PATCH /api/v1/admin/policies/{policy_id}
type UpdatePolicyRequest struct {
    Config  *json.RawMessage `json:"config,omitempty"`
    Enabled *bool            `json:"enabled,omitempty"`
}

func (s *Service) UpdatePolicy(ctx context.Context, policyID string, req UpdatePolicyRequest) error {
    orgID := middleware.GetOrgID(ctx)
    
    // Get existing policy
    policy, err := s.policyRepo.GetByID(ctx, policyID)
    if err != nil {
        return err
    }
    
    // Verify ownership via org_id
    if policy.OrgID != orgID {
        return ErrUnauthorized
    }
    
    // Update config (validate first)
    if req.Config != nil {
        _, err := s.policyEngine.NewPolicy(model.PolicyType(policy.PolicyType), *req.Config)
        if err != nil {
            return fmt.Errorf("invalid config: %w", err)
        }
        policy.Config = *req.Config
    }
    
    // Update enabled flag
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

#### 3.5 Delete Policy

```go
// DELETE /api/v1/admin/policies/{policy_id}
func (s *Service) DeletePolicy(ctx context.Context, policyID string) error {
    orgID := middleware.GetOrgID(ctx)
    
    policy, err := s.policyRepo.GetByID(ctx, policyID)
    if err != nil {
        return err
    }
    
    if policy.OrgID != orgID {
        return ErrUnauthorized
    }
    
    if err := s.policyRepo.Delete(ctx, policyID); err != nil {
        return err
    }
    
    s.policyEngine.InvalidateCache(ctx, policy.AppID)
    
    return nil
}
```

---

## Key Differences from Original Design

### ❌ Old Assumption
- Built-in policies might be "defaults" for all apps
- Might need special handling for system vs custom policies

### ✅ New Reality
- **All policies are opt-in** - nothing runs by default
- **Built-in vs custom is just a factory difference** - both stored in DB
- **Per-app configuration** - admin explicitly adds each policy
- **Same API for both** - UI/API doesn't care if policy is built-in or custom

---

## Benefits of This Approach

1. **Simplicity** - One table, one flow, one API
2. **Flexibility** - Apps can have 0 to N policies
3. **Isolation** - Apps don't inherit policies from org
4. **Predictability** - What you see in DB is what runs
5. **Testability** - Test apps can have zero policies
6. **Cost control** - Prod can have strict limits, dev can be relaxed
7. **Customer freedom** - Mix built-in and custom as needed

---

## Example Scenarios

### Scenario 1: New App (Zero Policies)
```
Org: "Acme Corp"
App: "New Test App"

Database:
  SELECT * FROM policies WHERE app_id = 'new-test-app'
  → Returns []

Runtime:
  → No policies loaded
  → No policy checks run
  → Requests pass through (other auth still applies)
```

### Scenario 2: Production App (Multiple Policies)
```
Org: "Acme Corp"
App: "Production API"

Database:
  SELECT * FROM policies WHERE app_id = 'prod-api' AND enabled = true
  → Returns:
    1. {type: "rate_limit", config: {rpm: 1000}}
    2. {type: "token_limit", config: {max: 8192}}
    3. {type: "custom_cel", config: {expr: "model != 'gpt-4-32k'"}}

Runtime:
  → Loads 3 policies
  → Constructs: [RateLimitPolicy, TokenLimitPolicy, CELPolicy]
  → Executes in order
```

### Scenario 3: Multiple Apps, Same Org
```
Org: "Acme Corp"

App 1: "Production API"
  → rate_limit (1000 rpm)
  → token_limit (8192 tokens)

App 2: "Development API"
  → token_limit (4096 tokens)

App 3: "Internal Tools"
  → (no policies)

Each app has independent policy configuration.
No shared state, no inheritance.
```

---

## Migration & Rollout

### Phase 1: Registry (No DB Changes)
- Implement `registry.go`
- Add `init()` to existing policies
- Update `engine.NewPolicy()` to check registry
- **Result:** Code prepared, but DB still has old data

### Phase 2: Admin API
- Create policy repository
- Add CRUD endpoints
- **Result:** Can create/edit policies via API

### Phase 3: UI
- Policy list page per app
- Add/edit/delete forms
- **Result:** Customers can self-service policies

### Phase 4: Data Migration (If Needed)
- If existing policies in DB, ensure they have correct format
- Validate all existing configs still compile

---

## Summary

**What changed from original design:**
- ✅ Built-in policies are **not defaults** - they're just factory definitions
- ✅ All policies stored in DB with `app_id` - explicit attachment
- ✅ Zero policies by default - admin must opt-in per app
- ✅ Same workflow for built-in and custom - only difference is `policy_type` value

**What stayed the same:**
- ✅ Registry pattern for extensibility
- ✅ Three-tier caching
- ✅ Multi-tenant isolation
- ✅ CEL support for custom logic

**Ready to implement?** Let me know if you want to start with Phase 1 (Registry)!
