package policies

import (
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

func TestRegistryInitialization(t *testing.T) {
	// Test that all built-in policies are registered via init()
	expectedPolicies := []model.PolicyType{
		model.PolicyTypeRateLimit,
		model.PolicyTypeTokenLimit,
		model.PolicyTypeModelAllowlist,
		model.PolicyTypeRequestSize,
	}

	for _, policyType := range expectedPolicies {
		t.Run(string(policyType), func(t *testing.T) {
			if !IsRegistered(policyType) {
				t.Errorf("policy type %s not registered", policyType)
			}

			factory, exists := GetFactory(policyType)
			if !exists {
				t.Errorf("factory for %s not found", policyType)
			}
			if factory == nil {
				t.Errorf("factory for %s is nil", policyType)
			}
		})
	}
}

func TestPolicyFactories(t *testing.T) {
	cache := kv.NewMemoryStore()
	deps := PolicyDependencies{Cache: cache}

	tests := []struct {
		name       string
		policyType model.PolicyType
		config     []byte
		wantErr    bool
	}{
		{
			name:       "rate_limit_valid",
			policyType: model.PolicyTypeRateLimit,
			config:     []byte(`{"requests_per_minute": 1000}`),
			wantErr:    false,
		},
		{
			name:       "rate_limit_invalid",
			policyType: model.PolicyTypeRateLimit,
			config:     []byte(`{invalid json}`),
			wantErr:    true,
		},
		{
			name:       "token_limit_valid",
			policyType: model.PolicyTypeTokenLimit,
			config:     []byte(`{"max_prompt_tokens": 4000}`),
			wantErr:    false,
		},
		{
			name:       "model_allowlist_valid",
			policyType: model.PolicyTypeModelAllowlist,
			config:     []byte(`{"allowed_model_ids": ["gpt-4"]}`),
			wantErr:    false,
		},
		{
			name:       "request_size_valid",
			policyType: model.PolicyTypeRequestSize,
			config:     []byte(`{"max_request_bytes": 51200}`),
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory, exists := GetFactory(tt.policyType)
			if !exists {
				t.Fatalf("factory for %s not found", tt.policyType)
			}

			policy, err := factory(tt.config, deps)
			if (err != nil) != tt.wantErr {
				t.Errorf("factory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && policy == nil {
				t.Error("factory() returned nil policy without error")
			}

			if !tt.wantErr && policy.Type() != tt.policyType {
				t.Errorf("policy.Type() = %v, want %v", policy.Type(), tt.policyType)
			}
		})
	}
}

func TestListAvailableTypes(t *testing.T) {
	types := ListAvailableTypes()

	// Should have at least 5 types (4 built-in + 1 custom_cel)
	if len(types) < 5 {
		t.Errorf("expected at least 5 policy types, got %d", len(types))
	}

	// Check that all types have required metadata
	for _, typeInfo := range types {
		if typeInfo.Type == "" {
			t.Error("policy type has empty Type field")
		}
		if typeInfo.Name == "" {
			t.Errorf("policy type %s has empty Name field", typeInfo.Type)
		}
		if typeInfo.Description == "" {
			t.Errorf("policy type %s has empty Description field", typeInfo.Type)
		}
		if typeInfo.ConfigSchema == nil {
			t.Errorf("policy type %s has nil ConfigSchema", typeInfo.Type)
		}
	}

	// Verify custom_cel is included
	hasCustomCEL := false
	for _, typeInfo := range types {
		if typeInfo.Type == model.PolicyTypeCustomCEL {
			hasCustomCEL = true
			if typeInfo.IsBuiltIn {
				t.Error("custom_cel should not be marked as built-in")
			}
		}
	}
	if !hasCustomCEL {
		t.Error("custom_cel policy type not found in available types")
	}

	// Verify built-in policies are marked correctly
	for _, typeInfo := range types {
		if typeInfo.Type != model.PolicyTypeCustomCEL && !typeInfo.IsBuiltIn {
			t.Errorf("policy type %s should be marked as built-in", typeInfo.Type)
		}
	}
}

func TestListRegisteredTypes(t *testing.T) {
	types := ListRegisteredTypes()

	// Should have exactly 4 built-in policies
	if len(types) != 4 {
		t.Errorf("expected 4 registered policy types, got %d", len(types))
	}

	// Verify all expected types are present
	expectedTypes := map[model.PolicyType]bool{
		model.PolicyTypeRateLimit:      false,
		model.PolicyTypeTokenLimit:     false,
		model.PolicyTypeModelAllowlist: false,
		model.PolicyTypeRequestSize:    false,
	}

	for _, policyType := range types {
		if _, exists := expectedTypes[policyType]; exists {
			expectedTypes[policyType] = true
		} else {
			t.Errorf("unexpected policy type registered: %s", policyType)
		}
	}

	for policyType, found := range expectedTypes {
		if !found {
			t.Errorf("expected policy type %s not found in registry", policyType)
		}
	}
}

func TestEngineWithRegistry(t *testing.T) {
	cache := kv.NewMemoryStore()
	engine := NewEngine(nil, cache)

	tests := []struct {
		name       string
		policyType model.PolicyType
		config     []byte
		wantErr    bool
	}{
		{
			name:       "rate_limit_via_engine",
			policyType: model.PolicyTypeRateLimit,
			config:     []byte(`{"requests_per_minute": 1000}`),
			wantErr:    false,
		},
		{
			name:       "token_limit_via_engine",
			policyType: model.PolicyTypeTokenLimit,
			config:     []byte(`{"max_prompt_tokens": 4000}`),
			wantErr:    false,
		},
		{
			name:       "custom_cel_via_engine",
			policyType: model.PolicyTypeCustomCEL,
			config:     []byte(`{"name":"test","pre_check_expression":"true"}`),
			wantErr:    false,
		},
		{
			name:       "unknown_policy_type",
			policyType: "unknown_type",
			config:     []byte(`{}`),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := engine.NewPolicy(tt.policyType, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("engine.NewPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && policy == nil {
				t.Error("engine.NewPolicy() returned nil policy without error")
			}

			if !tt.wantErr && policy.Type() != tt.policyType {
				t.Errorf("policy.Type() = %v, want %v", policy.Type(), tt.policyType)
			}
		})
	}
}
