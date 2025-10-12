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

func stringPtr(s string) *string {
	return &s
}
