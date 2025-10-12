package model

// AuthType represents the type of authentication
type AuthType string

const (
	AuthTypeAPIKey  AuthType = "api_key"
	AuthTypeOAuth2  AuthType = "oauth2"
	AuthTypeAzureAD AuthType = "azure_ad"
)

// String returns the string representation of AuthType
func (at AuthType) String() string {
	return string(at)
}

// IsValid returns true if the AuthType is one of the defined constants
func (at AuthType) IsValid() bool {
	switch at {
	case AuthTypeAPIKey, AuthTypeOAuth2, AuthTypeAzureAD:
		return true
	default:
		return false
	}
}
