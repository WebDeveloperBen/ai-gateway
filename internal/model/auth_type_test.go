package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthType_String(t *testing.T) {
	tests := []struct {
		name     string
		authType AuthType
		expected string
	}{
		{
			name:     "API Key",
			authType: AuthTypeAPIKey,
			expected: "api_key",
		},
		{
			name:     "OAuth2",
			authType: AuthTypeOAuth2,
			expected: "oauth2",
		},
		{
			name:     "Azure AD",
			authType: AuthTypeAzureAD,
			expected: "azure_ad",
		},
		{
			name:     "Custom auth type",
			authType: AuthType("custom"),
			expected: "custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.authType.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		authType AuthType
		expected bool
	}{
		{
			name:     "Valid API Key",
			authType: AuthTypeAPIKey,
			expected: true,
		},
		{
			name:     "Valid OAuth2",
			authType: AuthTypeOAuth2,
			expected: true,
		},
		{
			name:     "Valid Azure AD",
			authType: AuthTypeAzureAD,
			expected: true,
		},
		{
			name:     "Invalid empty string",
			authType: AuthType(""),
			expected: false,
		},
		{
			name:     "Invalid custom type",
			authType: AuthType("invalid"),
			expected: false,
		},
		{
			name:     "Invalid random string",
			authType: AuthType("random_auth"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.authType.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthType_Constants(t *testing.T) {
	t.Run("Constants have correct values", func(t *testing.T) {
		assert.Equal(t, "api_key", AuthTypeAPIKey.String())
		assert.Equal(t, "oauth2", AuthTypeOAuth2.String())
		assert.Equal(t, "azure_ad", AuthTypeAzureAD.String())
	})

	t.Run("All constants are valid", func(t *testing.T) {
		assert.True(t, AuthTypeAPIKey.IsValid())
		assert.True(t, AuthTypeOAuth2.IsValid())
		assert.True(t, AuthTypeAzureAD.IsValid())
	})
}
