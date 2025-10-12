# Policy Plugin Architecture - Design Document

## Overview

This document outlines a registration-based plugin architecture for the policy system, allowing new policy types to be added without modifying the core engine code.

## Current Architecture (Factory Pattern)

### Problems
1. **Tight coupling** - `engine.go` must know about every policy type
2. **Code changes required** - Adding a policy needs switch case modification
3. **Recompilation needed** - Can't add policies without redeploying
4. **Testing friction** - Hard to test custom policies in isolation

### Current Flow
```go
// engine.go NewPolicy() - MUST BE MODIFIED for each new policy
func (e *Engine) NewPolicy(policyType model.PolicyType, config []byte) (Policy, error) {
    switch policyType {
    case model.PolicyTypeRateLimit:
        return NewRateLimitPolicy(cfg, e.cache), nil
    case model.PolicyTypeTokenLimit:
        return NewTokenLimitPolicy(cfg), nil
    // ... add more cases for each new policy type
    }
}
```

---

## Proposed Architecture (Plugin Registry)

### Design Goals
1. **Zero-touch core** - No `engine.go` changes for new policies
2. **Self-contained** - Each policy registers itself via `init()`
3. **Type-safe** - Compile-time checking, no reflection magic
4. **Backward compatible** - Existing policies work unchanged

### Architecture Components

```
┌─────────────────────────────────────────────────────────┐
│                    Policy Registry                       │
│  (Singleton, thread-safe, populated at startup)         │
│                                                          │
│  Map: PolicyType -> PolicyFactory                       │
└─────────────────────────────────────────────────────────┘
                          ▲
                          │ Register() at init()
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
┌───────▼───────┐ ┌───────▼───────┐ ┌─────▼─────────┐
│ rate_limit.go │ │ token_limit.go│ │ custom_cel.go │
│               │ │               │ │               │
│ func init() { │ │ func init() { │ │ func init() { │
│   Register()  │ │   Register()  │ │   Register()  │
│ }             │ │ }             │ │ }             │
└───────────────┘ └───────────────┘ └───────────────┘
```

### Implementation

#### Step 1: Policy Registry (New File)

```go
// internal/gateway/policies/registry.go
package policies

import (
    "fmt"
    "sync"
    "github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
    "github.com/WebDeveloperBen/ai-gateway/internal/model"
)

// PolicyFactory creates a policy instance from raw JSON config
// Dependencies (like cache) are passed as arguments
type PolicyFactory func(config []byte, deps PolicyDependencies) (Policy, error)

// PolicyDependencies holds dependencies policies might need
type PolicyDependencies struct {
    Cache kv.KvStore
}

// Global registry - populated at init() time
var (
    registry   = make(map[model.PolicyType]PolicyFactory)
    registryMu sync.RWMutex
)

// Register adds a policy factory to the registry
// Typically called from init() functions in policy files
func Register(policyType model.PolicyType, factory PolicyFactory) {
    registryMu.Lock()
    defer registryMu.Unlock()
    
    if _, exists := registry[policyType]; exists {
        panic(fmt.Sprintf("policy type %s already registered", policyType))
    }
    
    registry[policyType] = factory
}

// GetFactory retrieves a policy factory by type
func GetFactory(policyType model.PolicyType) (PolicyFactory, bool) {
    registryMu.RLock()
    defer registryMu.RUnlock()
    
    factory, exists := registry[policyType]
    return factory, exists
}

// ListRegistered returns all registered policy types
func ListRegistered() []model.PolicyType {
    registryMu.RLock()
    defer registryMu.RUnlock()
    
    types := make([]model.PolicyType, 0, len(registry))
    for t := range registry {
        types = append(types, t)
    }
    return types
}
```

#### Step 2: Update Engine.NewPolicy()

```go
// internal/gateway/policies/engine.go
func (e *Engine) NewPolicy(policyType model.PolicyType, config []byte) (Policy, error) {
    // Try registry first (plugin architecture)
    factory, exists := GetFactory(policyType)
    if exists {
        deps := PolicyDependencies{Cache: e.cache}
        return factory(config, deps)
    }
    
    // Fallback: Treat unknown types as custom CEL
    return NewCELPolicy(policyType, config)
}
```

#### Step 3: Convert Existing Policies to Plugins

**Before (rate_limit.go):**
```go
package policies

type RateLimitPolicy struct { ... }

func NewRateLimitPolicy(cfg model.RateLimitConfig, cache kv.KvStore) *RateLimitPolicy {
    return &RateLimitPolicy{...}
}
```

**After (rate_limit.go with registration):**
```go
package policies

import "github.com/WebDeveloperBen/ai-gateway/internal/model"

func init() {
    // Self-register at startup
    Register(model.PolicyTypeRateLimit, func(config []byte, deps PolicyDependencies) (Policy, error) {
        var cfg model.RateLimitConfig
        if err := json.Unmarshal(config, &cfg); err != nil {
            return nil, fmt.Errorf("invalid rate limit config: %w", err)
        }
        return NewRateLimitPolicy(cfg, deps.Cache), nil
    })
}

type RateLimitPolicy struct { ... }

func NewRateLimitPolicy(cfg model.RateLimitConfig, cache kv.KvStore) *RateLimitPolicy {
    return &RateLimitPolicy{...}
}
```

#### Step 4: Adding New Policies (Zero Core Changes)

**Example: Content Filter Policy**

```go
// internal/gateway/policies/content_filter.go
package policies

import (
    "context"
    "fmt"
    "regexp"
    "github.com/WebDeveloperBen/ai-gateway/internal/model"
)

func init() {
    // Register automatically when package loads
    Register(model.PolicyTypeContentFilter, func(config []byte, deps PolicyDependencies) (Policy, error) {
        var cfg ContentFilterConfig
        if err := json.Unmarshal(config, &cfg); err != nil {
            return nil, fmt.Errorf("invalid content filter config: %w", err)
        }
        return NewContentFilterPolicy(cfg)
    })
}

type ContentFilterConfig struct {
    BlockedPatterns []string `json:"blocked_patterns"`
}

type ContentFilterPolicy struct {
    config   ContentFilterConfig
    patterns []*regexp.Regexp
}

func NewContentFilterPolicy(cfg ContentFilterConfig) (*ContentFilterPolicy, error) {
    patterns := make([]*regexp.Regexp, len(cfg.BlockedPatterns))
    for i, pattern := range cfg.BlockedPatterns {
        re, err := regexp.Compile(pattern)
        if err != nil {
            return nil, fmt.Errorf("invalid pattern %s: %w", pattern, err)
        }
        patterns[i] = re
    }
    
    return &ContentFilterPolicy{
        config:   cfg,
        patterns: patterns,
    }, nil
}

func (p *ContentFilterPolicy) Type() model.PolicyType {
    return model.PolicyTypeContentFilter
}

func (p *ContentFilterPolicy) PreCheck(ctx context.Context, req *PreRequestContext) error {
    // Check request body against patterns
    // Return error if blocked content detected
    return nil
}

func (p *ContentFilterPolicy) PostCheck(ctx context.Context, req *PostRequestContext) {
    // Optional: scan response content
}
```

**That's it! No changes to `engine.go` needed.**

---

## Benefits

### 1. **Extensibility**
- Add new policies by creating a new file
- No modification to core engine code
- Policies live in separate files/packages

### 2. **Testability**
```go
// Test policies in isolation
func TestContentFilterPolicy(t *testing.T) {
    policy, _ := NewContentFilterPolicy(ContentFilterConfig{
        BlockedPatterns: []string{"badword"},
    })
    
    err := policy.PreCheck(ctx, &PreRequestContext{...})
    assert.Error(t, err)
}
```

### 3. **Discovery**
```go
// List all available policies
registeredTypes := policies.ListRegistered()
// ["rate_limit", "token_limit", "content_filter", ...]
```

### 4. **Future: External Plugins**
Could extend to load external plugins:
- Go plugins (`.so` files)
- WASM modules
- Sidecar services

### 5. **Admin UI Integration**
```go
// UI can query available policy types dynamically
GET /api/v1/admin/policies/types
{
    "policy_types": [
        {"id": "rate_limit", "name": "Rate Limit", "description": "..."},
        {"id": "token_limit", "name": "Token Limit", "description": "..."},
        {"id": "content_filter", "name": "Content Filter", "description": "..."}
    ]
}
```

---

## Migration Path

### Phase 1: Add Registry (30 min)
- Create `registry.go`
- No breaking changes, registry is additive

### Phase 2: Update Engine (15 min)
- Modify `engine.NewPolicy()` to check registry first
- Keep fallback to CEL for backward compatibility

### Phase 3: Convert Policies (2-3 hours)
- Add `init()` functions to existing policy files
- Test each policy still works
- One policy at a time, low risk

### Phase 4: Remove Switch Statement (15 min)
- Once all policies registered, remove switch cases
- Keep CEL fallback for custom policies

### Total: 4-6 hours (as estimated)

---

## Alternatives Considered

### 1. **Reflection-based Registry**
```go
Register("rate_limit", RateLimitPolicy{})
```
**Rejected:** Type safety lost, runtime errors instead of compile-time

### 2. **Code Generation**
```go
//go:generate go run generate_registry.go
```
**Rejected:** Adds build complexity, harder to debug

### 3. **Go Plugins (`.so` files)**
```go
plugin.Open("policies/content_filter.so")
```
**Rejected:** Platform-specific, deployment complexity, security concerns

---

## Recommendation

**For v1.0:** Keep current architecture
- System is working well
- Only 4 policy types currently
- No immediate need for extensibility

**For v2.0:** Implement plugin registry
- When you need:
  - Customer-specific policies
  - Third-party integrations
  - Dynamic policy loading
  - Admin UI policy marketplace

**Effort:** 4-6 hours when needed
**Risk:** Low - backward compatible, incremental migration
**Complexity:** Medium - well-understood pattern, clear boundaries

---

## Example: Full Before/After

### Before (Current)
```go
// To add content_filter policy:
// 1. Edit model/policy.go
const PolicyTypeContentFilter PolicyType = "content_filter"

// 2. Edit engine.go
case model.PolicyTypeContentFilter:
    var cfg model.ContentFilterConfig
    if err := json.Unmarshal(config, &cfg); err != nil {
        return nil, err
    }
    return NewContentFilterPolicy(cfg), nil

// 3. Create content_filter.go
// 4. Rebuild and redeploy
```

### After (Plugin)
```go
// To add content_filter policy:
// 1. Create content_filter.go with init() function
func init() {
    Register(model.PolicyTypeContentFilter, factoryFunc)
}

// 2. Rebuild and redeploy (or hot-reload if using plugins)
// That's it! No engine.go changes needed
```

---

## Conclusion

Plugin architecture is a **nice-to-have**, not a **must-have** for v1.0.

- Current system works well for ~10 policy types
- Registry pattern adds flexibility without complexity
- Migration is incremental and low-risk
- Defer until you have a concrete need (customer-specific policies, marketplace, etc.)

**Decision:** Implement when policy count > 10 or when external contribution is needed.
