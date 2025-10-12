# Provider Architecture

This document clarifies the separation between **Provider Support** and **Model Deployments**.

## Two Separate Concerns

### 1. Provider Support (Static - Code)

**What it is:**
- Which LLM provider APIs the gateway can communicate with
- How to transform requests for each provider's API format
- Authentication mechanisms for each provider

**Where it lives:**
- `internal/provider/` - Provider adapter implementations
- `internal/gateway/providers.go` - List of supported providers
- `internal/api/admin/gateway/route.go` - HTTP route registration

**Key principle:**
**All coded providers are ALWAYS available.** If the code exists, the routes are registered.

**Example:**
```go
// internal/gateway/providers.go
func DefaultProviders() []*provider.ProviderConfig {
    return []*provider.ProviderConfig{
        {
            Prefix:      "/azure/openai",
            DisplayName: "Azure OpenAI",
            Description: "Microsoft Azure OpenAI Service",
            Enabled:     true,  // Always true if coded
        },
        {
            Prefix:      "/openai",
            DisplayName: "OpenAI",
            Enabled:     true,
        },
        {
            Prefix:      "/groq",
            DisplayName: "Groq",
            Enabled:     true,
        },
    }
}
```

**API Routes Created:**
- `POST /api/azure/openai/v1/chat/completions`
- `POST /api/openai/v1/chat/completions`
- `POST /api/groq/v1/chat/completions`

### 2. Model Deployments (Dynamic - Database)

**What it is:**
- Specific model instances available for routing
- Per-org/tenant configuration
- Runtime configuration (can be changed without redeploy)

**Where it lives:**
- `db/schema/models.hcl` - Database schema
- `internal/gateway/registry.go` - Runtime model routing
- Future: `internal/repository/models` - CRUD operations

**Key principle:**
**Deployments determine WHERE requests are routed, not WHETHER routes exist.**

**Example:**
```sql
-- Models table (future implementation)
INSERT INTO models (org_id, provider, model_name, deployment_name, endpoint_url)
VALUES
  ('org-123', 'azure', 'gpt-4', 'prod-gpt4-deployment', 'https://prod.openai.azure.com'),
  ('org-123', 'azure', 'gpt-3.5', 'dev-gpt35-deployment', 'https://dev.openai.azure.com'),
  ('org-456', 'openai', 'gpt-4', 'gpt-4-turbo-preview', 'https://api.openai.com');
```

## Current Implementation (main.go)

```go
// 1. MODEL DEPLOYMENTS (request routing)
// TODO: Load from database
reg := gateway.NewRegistry(ctx, kvStore)
_ = gateway.EnsureRegistryPopulated(reg, func() []model.ModelDeployment {
    return []model.ModelDeployment{
        {
            Model:      "gpt-4.1",
            Deployment: "dev-openai-gpt4-1",
            Provider:   "azure",  // Which provider adapter to use
            Tenant:     "default",
            Meta:       map[string]string{...},
        },
    }
})

// 2. PROVIDER SUPPORT (API routes)
// Register all coded providers
for _, providerCfg := range gateway.DefaultProviders() {
    if providerCfg.Enabled {
        apigw.RegisterProvider(apigrp, providerCfg, core)
    }
}
```

## Request Flow

```
1. Request arrives: POST /api/azure/openai/v1/chat/completions
                    with { "model": "gpt-4.1" }

2. Route matches:   Provider config for "/azure/openai" exists
                    → Route registered ✓

3. Deployment:      Gateway looks up "gpt-4.1" in deployment registry
                    → Finds deployment "dev-openai-gpt4-1"

4. Adapter:         Uses Azure OpenAI adapter (from deployment.Provider)
                    → Rewrites request for Azure API

5. Forward:         Proxies to https://dev-insurgence-openai.openai.azure.com
```

## Adding a New Provider

### Step 1: Create Provider Adapter

```go
// internal/provider/groq/provider.go
package groq

type Adapter struct {
    BaseURL string
    Keys    provider.KeySource
}

func (a *Adapter) Prefix() string { return "/groq" }

func (a *Adapter) Rewrite(req *http.Request, suffix string, info provider.ReqInfo) error {
    // Transform request for Groq API
}
```

### Step 2: Add to Provider List

```go
// internal/gateway/providers.go
func DefaultProviders() []*provider.ProviderConfig {
    return []*provider.ProviderConfig{
        // ... existing providers
        {
            Prefix:      "/groq",
            DisplayName: "Groq",
            Description: "Groq's ultra-fast LLM inference",
            Enabled:     true,
        },
    }
}
```

**That's it!** Routes are automatically registered at `/api/groq/v1/*`

## Future: Database-Driven Deployments

When you implement the models repository:

```go
// internal/repository/models/postgres.go
func (r *postgresRepo) ListByOrgID(ctx context.Context, orgID uuid.UUID) ([]*model.Model, error) {
    models, err := r.q.ListModels(ctx, orgID)
    // Convert to model.ModelDeployment
}

// cmd/proxy/main.go
modelsRepo := modelsrepo.NewPostgresRepo(pg.Queries)
deployments, _ := modelsRepo.ListByOrgID(ctx, defaultOrgID)
reg := gateway.NewRegistry(ctx, kvStore)
_ = gateway.EnsureRegistryPopulated(reg, func() []model.ModelDeployment {
    return deployments
})
```

## Key Takeaways

✅ **Provider Support** = Static code = Always available
✅ **Model Deployments** = Dynamic data = Runtime configuration
✅ **Separation of concerns** = Provider adapters don't care about deployments
✅ **Deployment registry** = Only used for request routing, not route registration
