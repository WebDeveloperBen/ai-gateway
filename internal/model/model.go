// Package model contains top-level domain models and types for AI instances and related system entities.
package model

import (
	"time"

	"github.com/google/uuid"
)

// AuthConfig represents structured authentication configuration
type AuthConfig struct {
	Type AuthType `json:"type"`
	// API Key authentication
	APIKey *string `json:"api_key,omitempty"`
	// OAuth2 authentication
	ClientID     *string `json:"client_id,omitempty"`
	ClientSecret *string `json:"client_secret,omitempty"`
	TokenURL     *string `json:"token_url,omitempty"`
	// Azure authentication
	TenantID *string `json:"tenant_id,omitempty"`
	Resource *string `json:"resource,omitempty"`
	// Generic additional fields for extensibility
	Additional map[string]interface{} `json:"additional,omitempty"`
}

// Model represents an AI model configuration in the catalog
type Model struct {
	ID             uuid.UUID
	OrgID          uuid.UUID
	Provider       string
	ModelName      string
	DeploymentName *string
	EndpointURL    string
	AuthType       AuthType
	AuthConfig     AuthConfig
	Metadata       map[string]interface{}
	Enabled        bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
