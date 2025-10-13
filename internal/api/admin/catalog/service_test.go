package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

func TestConvertAuthConfig(t *testing.T) {
	// Test API key conversion
	apiConfig := AuthConfig{
		Type:   AuthTypeAPIKey,
		APIKey: stringPtr("test-key"),
	}

	modelConfig := convertAuthConfigFromAPI(apiConfig)
	assert.Equal(t, model.AuthTypeAPIKey, modelConfig.Type)
	assert.Equal(t, "test-key", *modelConfig.APIKey)

	// Convert back
	apiConfig2 := convertAuthConfigToAPI(modelConfig)
	assert.Equal(t, AuthTypeAPIKey, apiConfig2.Type)
	assert.Equal(t, "test-key", *apiConfig2.APIKey)
}

func TestConvertAuthConfig_OAuth(t *testing.T) {
	// Test OAuth conversion
	apiConfig := AuthConfig{
		Type:         AuthTypeOAuth2,
		ClientID:     stringPtr("client-id"),
		ClientSecret: stringPtr("client-secret"),
		TokenURL:     stringPtr("https://example.com/token"),
	}

	modelConfig := convertAuthConfigFromAPI(apiConfig)
	assert.Equal(t, model.AuthTypeOAuth2, modelConfig.Type)
	assert.Equal(t, "client-id", *modelConfig.ClientID)
	assert.Equal(t, "client-secret", *modelConfig.ClientSecret)
	assert.Equal(t, "https://example.com/token", *modelConfig.TokenURL)

	// Convert back
	apiConfig2 := convertAuthConfigToAPI(modelConfig)
	assert.Equal(t, AuthTypeOAuth2, apiConfig2.Type)
	assert.Equal(t, "client-id", *apiConfig2.ClientID)
	assert.Equal(t, "client-secret", *apiConfig2.ClientSecret)
	assert.Equal(t, "https://example.com/token", *apiConfig2.TokenURL)
}

func TestAPIAuthTypeToModel(t *testing.T) {
	tests := []struct {
		name     string
		apiType  AuthType
		expected model.AuthType
	}{
		{"APIKey", AuthTypeAPIKey, model.AuthTypeAPIKey},
		{"OAuth2", AuthTypeOAuth2, model.AuthTypeOAuth2},
		{"AzureAD", AuthTypeAzureAD, model.AuthTypeAzureAD},
		{"Unknown", AuthType("unknown"), model.AuthType("unknown")}, // fallback
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := APIAuthTypeToModel(tt.apiType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestModelAuthTypeToAPI(t *testing.T) {
	tests := []struct {
		name      string
		modelType model.AuthType
		expected  AuthType
	}{
		{"APIKey", model.AuthTypeAPIKey, AuthTypeAPIKey},
		{"OAuth2", model.AuthTypeOAuth2, AuthTypeOAuth2},
		{"AzureAD", model.AuthTypeAzureAD, AuthTypeAzureAD},
		{"Unknown", model.AuthType("unknown"), AuthType("unknown")}, // fallback
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ModelAuthTypeToAPI(tt.modelType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
