// Package gateway provides the HTTP endpoints that proxy client requests
// through the gateway to upstream LLM providers (Azure OpenAI, OpenAI, etc.).
package gateway

import (
	"fmt"
	"net/http"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway"
	"github.com/WebDeveloperBen/ai-gateway/internal/provider"
	"github.com/danielgtaylor/huma/v2"
)

type endpointSpec struct {
	Path        string
	Summary     string
	Description string
}

type providerConfig struct {
	Prefix      string
	DisplayName string
	Description string
	Enabled     bool
	Endpoints   []endpointSpec
}

// Common endpoint specifications that can be reused across providers
var (
	openAICompatibleEndpoints = []endpointSpec{
		{
			Path:        "/v1/chat/completions",
			Summary:     "Create chat completion",
			Description: "Creates a model response for a given chat conversation with support for streaming, function calling, and vision. Supports policy enforcement, rate limiting, and usage tracking.",
		},
		{
			Path:        "/v1/completions",
			Summary:     "Create completion",
			Description: "Creates a completion for a provided prompt. Legacy endpoint - chat completions are recommended for most use cases.",
		},
		{
			Path:        "/v1/embeddings",
			Summary:     "Create embeddings",
			Description: "Creates an embedding vector representing the input text for semantic search and similarity tasks.",
		},
	}
)

var supportedProviders = []providerConfig{
	{
		Prefix:      provider.AzureOpenAIPrefix,
		DisplayName: "Azure OpenAI",
		Description: "Microsoft Azure OpenAI Service with deployment-based routing and API key authentication",
		Enabled:     true,
		Endpoints:   openAICompatibleEndpoints,
	},
	{
		Prefix:      provider.OpenAIPrefix,
		DisplayName: "OpenAI",
		Description: "OpenAI API with Bearer token authentication and organization header support",
		Enabled:     true,
		Endpoints:   openAICompatibleEndpoints,
	},
}

func RegisterAllProviders(grp *huma.Group, core *gateway.Core) {
	for _, providerCfg := range supportedProviders {
		if !providerCfg.Enabled {
			continue
		}
		registerProvider(grp, &providerCfg, core)
	}
}

// RegisterProvider registers a single provider (useful for testing)
func RegisterProvider(grp *huma.Group, providerCfg *provider.ProviderConfig, core *gateway.Core) {
	// Convert provider.ProviderConfig to local providerConfig
	localCfg := providerConfig{
		Prefix:      providerCfg.Prefix,
		DisplayName: providerCfg.DisplayName,
		Description: providerCfg.Description,
		Enabled:     providerCfg.Enabled,
		Endpoints:   openAICompatibleEndpoints,
	}
	registerProvider(grp, &localCfg, core)
}

func registerProvider(grp *huma.Group, providerCfg *providerConfig, core *gateway.Core) {
	h := core.StreamingHandler()

	for _, spec := range providerCfg.Endpoints {
		huma.Register(grp, huma.Operation{
			OperationID:   sanitizeOperationID(providerCfg.Prefix + spec.Path),
			Method:        http.MethodPost,
			Path:          providerCfg.Prefix + spec.Path,
			Summary:       spec.Summary,
			Description:   buildDescription(spec.Description, providerCfg),
			DefaultStatus: http.StatusOK,
			Tags:          []string{providerCfg.DisplayName},
		}, h)
	}
}

func buildDescription(baseDesc string, cfg *providerConfig) string {
	return fmt.Sprintf("%s\n\n**Provider**: %s  \n**Endpoint**: `%s`  \n**Details**: %s",
		baseDesc,
		cfg.DisplayName,
		cfg.Prefix,
		cfg.Description,
	)
}

func sanitizeOperationID(path string) string {
	operationID := "proxy"
	for _, char := range path {
		if char == '/' {
			operationID += "-"
		} else {
			operationID += string(char)
		}
	}
	return operationID
}
