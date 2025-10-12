# Adding New LLM Providers

This guide explains how to add support for new LLM providers to the AI Gateway.

## Quick Reference

**OpenAI-Compatible Providers** (no custom adapter needed):
- Together AI
- Groq
- Fireworks
- OpenRouter
- Perplexity
- Anyscale

**Providers Requiring Custom Adapters** (different API contracts):
- Anthropic (Claude)
- Google (Gemini)
- Cohere
- Mistral (minor differences)

## Adding an OpenAI-Compatible Provider

For providers that use the OpenAI API schema, simply add them to the registry:

### 1. Add Provider to Registry

Edit `internal/gateway/providers.go`:

```go
func DefaultProviders() []*provider.ProviderConfig {
    return []*provider.ProviderConfig{
        {
            Prefix:      provider.AzureOpenAIPrefix,
            DisplayName: "Azure OpenAI",
            Description: "Microsoft Azure OpenAI Service with deployment-based routing and API key authentication",
            Enabled:     true,
        },
        {
            Prefix:      provider.OpenAIPrefix,
            DisplayName: "OpenAI",
            Description: "OpenAI API with Bearer token authentication and organization header support",
            Enabled:     false,
        },
        // Add your new provider here
        {
            Prefix:      "/groq",  // URL prefix for routing
            DisplayName: "Groq",   // Display name in API docs
            Description: "Groq's ultra-fast LLM inference with OpenAI-compatible API",
            Enabled:     true,     // Set to true to enable
        },
    }
}
```

### 2. Add Prefix Constant

Edit `internal/provider/provider.go`:

```go
const (
    AzureOpenAIPrefix = "/azure/openai"
    OpenAIPrefix      = "/openai"
    GroqPrefix        = "/groq"  // Add your provider prefix
)
```

### 3. Create Passthrough Adapter

Create `internal/provider/groq/provider.go`:

```go
package groq

import (
    "net/http"
    "strings"

    "github.com/WebDeveloperBen/ai-gateway/internal/gateway/loadbalancing"
    "github.com/WebDeveloperBen/ai-gateway/internal/model"
    "github.com/WebDeveloperBen/ai-gateway/internal/provider"
)

type Adapter struct {
    BaseURL  string
    Keys     provider.KeySource
}

func New() *Adapter {
    return &Adapter{
        BaseURL: "api.groq.com",
        Keys:    provider.KeySource{EnvVar: "GROQ_API_KEY"},
    }
}

func (a *Adapter) Prefix() string { return provider.GroqPrefix }

func (a *Adapter) Rewrite(req *http.Request, suffix string, info provider.ReqInfo) error {
    base, _ := provider.EnsureAbsoluteBase(a.BaseURL, "api.groq.com")
    u, _ := provider.JoinURL(base, []string{suffix}, provider.CopyQuery(req))
    provider.SetUpstreamURL(req, u)

    provider.StripCallerAuth(req.Header)
    if key := a.Keys.Resolve(info.Tenant, "GROQ_API_KEY"); key != "" {
        req.Header.Set("Authorization", "Bearer "+key)
    }
    return nil
}

func BuildProvider(deployments []model.ModelDeployment, selector loadbalancing.InstanceSelector) *Adapter {
    return New()
}
```

### 4. Register in Bootstrap

Edit `internal/gateway/providers.go` to add your provider to the switch statement in the `Build` method:

```go
func (r *ProviderRegistry) Build(deployments []model.ModelDeployment, selector loadbalancing.InstanceSelector) {
    for _, cfg := range r.configs {
        if !cfg.Enabled {
            continue
        }

        var adapter provider.Adapter

        switch cfg.Prefix {
        case provider.AzureOpenAIPrefix:
            adapter = azureopenai.BuildProvider(deployments, selector)
        case provider.OpenAIPrefix:
            adapter = openai.BuildProvider(deployments, selector)
        case provider.GroqPrefix:
            adapter = groq.BuildProvider(deployments, selector)
        }

        if adapter != nil {
            r.adapters[cfg.Prefix] = adapter
        }
    }
}
```

### 5. Build and Test

```bash
# Build the application
go build -o ./bin/proxy ./cmd/proxy

# Test the endpoint
curl -X POST http://localhost:8080/api/groq/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "mixtral-8x7b-32768",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

## Adding a Provider with Custom API

For providers like Anthropic that don't use the OpenAI API schema:

### 1. Create Custom Adapter

The adapter must implement custom request/response transformation logic.

Example structure for `internal/provider/anthropic/provider.go`:

```go
type Adapter struct {
    BaseURL string
    Keys    provider.KeySource
}

func (a *Adapter) Rewrite(req *http.Request, suffix string, info provider.ReqInfo) error {
    // Custom logic to:
    // 1. Transform OpenAI request format to Anthropic format
    // 2. Set correct endpoint path (e.g., /v1/messages)
    // 3. Set custom headers (e.g., x-api-key, anthropic-version)
    // 4. Handle streaming format differences
}
```

### 2. Add Response Transformation

You may also need middleware to transform responses back to OpenAI format for client compatibility.

## Environment Variables

Add your provider's API key to `.env`:

```bash
# Groq
GROQ_API_KEY=your-groq-api-key-here

# Anthropic
ANTHROPIC_API_KEY=your-anthropic-api-key-here
```

## API Documentation

Once registered, your provider will automatically appear in the `/docs` Scalar documentation with:
- Provider name and description
- All supported endpoints (chat completions, embeddings, etc.)
- Proper categorization and tagging

## Testing

Create tests in `internal/provider/yourprovider/provider_test.go` following the pattern in existing provider tests.
