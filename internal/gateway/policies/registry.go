package policies

import (
	"fmt"
	"sync"

	"github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

// PolicyFactory creates a policy instance from raw JSON config and dependencies
type PolicyFactory func(config []byte, deps PolicyDependencies) (Policy, error)

// PolicyDependencies holds shared dependencies that policies might need
type PolicyDependencies struct {
	Cache kv.KvStore
}

// PolicyTypeMetadata describes a policy type for UI/API consumers
type PolicyTypeMetadata struct {
	Type         model.PolicyType `json:"type"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	ConfigSchema map[string]any   `json:"config_schema"`
	IsBuiltIn    bool             `json:"is_built_in"`
}

var (
	registry   = make(map[model.PolicyType]PolicyFactory)
	registryMu sync.RWMutex
)

// Register adds a policy factory to the global registry
// This is typically called from init() functions in policy implementation files
// Panics if the policy type is already registered (fail-fast at startup)
func Register(policyType model.PolicyType, factory PolicyFactory) {
	registryMu.Lock()
	defer registryMu.Unlock()

	if _, exists := registry[policyType]; exists {
		panic(fmt.Sprintf("policy type %s already registered", policyType))
	}

	registry[policyType] = factory
}

// GetFactory retrieves a policy factory by type
// Returns (factory, true) if found, (nil, false) if not registered
func GetFactory(policyType model.PolicyType) (PolicyFactory, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()

	factory, exists := registry[policyType]
	return factory, exists
}

// IsRegistered checks if a policy type has been registered
func IsRegistered(policyType model.PolicyType) bool {
	registryMu.RLock()
	defer registryMu.RUnlock()

	_, exists := registry[policyType]
	return exists
}

// ListRegisteredTypes returns all registered policy types
// Useful for debugging and admin UI to show available built-in policies
func ListRegisteredTypes() []model.PolicyType {
	registryMu.RLock()
	defer registryMu.RUnlock()

	types := make([]model.PolicyType, 0, len(registry))
	for t := range registry {
		types = append(types, t)
	}
	return types
}

// ListAvailableTypes returns metadata for all available policy types
// This includes both built-in (registered) policies and custom CEL
func ListAvailableTypes() []PolicyTypeMetadata {
	registryMu.RLock()
	defer registryMu.RUnlock()

	types := make([]PolicyTypeMetadata, 0, len(registry)+1)

	for policyType := range registry {
		types = append(types, PolicyTypeMetadata{
			Type:         policyType,
			Name:         formatPolicyName(policyType),
			Description:  getPolicyDescription(policyType),
			ConfigSchema: getPolicySchema(policyType),
			IsBuiltIn:    true,
		})
	}

	types = append(types, PolicyTypeMetadata{
		Type:         model.PolicyTypeCustomCEL,
		Name:         "Custom CEL Policy",
		Description:  "Define custom policy logic using Common Expression Language (CEL)",
		ConfigSchema: getCELConfigSchema(),
		IsBuiltIn:    false,
	})

	return types
}

// formatPolicyName converts a policy type to a human-readable name
func formatPolicyName(policyType model.PolicyType) string {
	switch policyType {
	case model.PolicyTypeRateLimit:
		return "Rate Limit"
	case model.PolicyTypeTokenLimit:
		return "Token Limit"
	case model.PolicyTypeModelAllowlist:
		return "Model Allowlist"
	case model.PolicyTypeRequestSize:
		return "Request Size Limit"
	default:
		return string(policyType)
	}
}

// getPolicyDescription returns a human-readable description of the policy
func getPolicyDescription(policyType model.PolicyType) string {
	switch policyType {
	case model.PolicyTypeRateLimit:
		return "Limit the number of requests per time window (requests per minute)"
	case model.PolicyTypeTokenLimit:
		return "Limit the number of tokens in requests and responses"
	case model.PolicyTypeModelAllowlist:
		return "Restrict which AI models can be used"
	case model.PolicyTypeRequestSize:
		return "Limit the maximum size of request bodies"
	default:
		return ""
	}
}

// getPolicySchema returns JSON Schema for the policy's configuration
func getPolicySchema(policyType model.PolicyType) map[string]any {
	switch policyType {
	case model.PolicyTypeRateLimit:
		return map[string]any{
			"type": "object",
			"properties": map[string]any{
				"requests_per_minute": map[string]any{
					"type":        "integer",
					"minimum":     1,
					"description": "Maximum requests per minute",
				},
			},
			"required": []string{"requests_per_minute"},
		}

	case model.PolicyTypeTokenLimit:
		return map[string]any{
			"type": "object",
			"properties": map[string]any{
				"max_prompt_tokens": map[string]any{
					"type":        "integer",
					"minimum":     1,
					"description": "Maximum tokens in prompt",
				},
				"max_completion_tokens": map[string]any{
					"type":        "integer",
					"minimum":     1,
					"description": "Maximum tokens in completion",
				},
				"max_total_tokens": map[string]any{
					"type":        "integer",
					"minimum":     1,
					"description": "Maximum total tokens (prompt + completion)",
				},
			},
		}

	case model.PolicyTypeModelAllowlist:
		return map[string]any{
			"type": "object",
			"properties": map[string]any{
				"allowed_model_ids": map[string]any{
					"type": "array",
					"items": map[string]any{
						"type": "string",
					},
					"description": "List of allowed model IDs (e.g., 'gpt-4', 'gpt-3.5-turbo')",
				},
			},
			"required": []string{"allowed_model_ids"},
		}

	case model.PolicyTypeRequestSize:
		return map[string]any{
			"type": "object",
			"properties": map[string]any{
				"max_request_bytes": map[string]any{
					"type":        "integer",
					"minimum":     1,
					"description": "Maximum request body size in bytes",
				},
			},
			"required": []string{"max_request_bytes"},
		}

	default:
		return map[string]any{"type": "object"}
	}
}

// getCELConfigSchema returns the JSON Schema for custom CEL policies
func getCELConfigSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type":        "string",
				"description": "Policy name (for display purposes)",
			},
			"description": map[string]any{
				"type":        "string",
				"description": "Policy description",
			},
			"pre_check_expression": map[string]any{
				"type":        "string",
				"description": "CEL expression evaluated before request (must return boolean)",
			},
			"post_check_expression": map[string]any{
				"type":        "string",
				"description": "CEL expression evaluated after response (optional)",
			},
		},
		"required": []string{"name", "pre_check_expression"},
	}
}
