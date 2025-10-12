package catalog

import (
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

// AuthType represents the type of authentication
type AuthType string

const (
	AuthTypeAPIKey  AuthType = "api_key"
	AuthTypeOAuth2  AuthType = "oauth2"
	AuthTypeAzureAD AuthType = "azure_ad"
)

// AuthConfig represents structured authentication configuration for API responses
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

type Model struct {
	ID             string         `json:"id"`
	OrgID          string         `json:"org_id"`
	Provider       string         `json:"provider"`
	ModelName      string         `json:"model_name"`
	DeploymentName *string        `json:"deployment_name,omitempty"`
	EndpointURL    string         `json:"endpoint_url"`
	AuthType       AuthType       `json:"auth_type"`
	AuthConfig     AuthConfig     `json:"auth_config"`
	Metadata       map[string]any `json:"metadata"`
	Enabled        bool           `json:"enabled"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type CreateModelBody struct {
	Provider       string         `json:"provider" required:"true"`
	ModelName      string         `json:"model_name" required:"true"`
	DeploymentName *string        `json:"deployment_name,omitempty"`
	EndpointURL    string         `json:"endpoint_url" required:"true"`
	AuthType       AuthType       `json:"auth_type" required:"true"`
	AuthConfig     AuthConfig     `json:"auth_config" required:"true"`
	Metadata       map[string]any `json:"metadata,omitempty"`
	Enabled        bool           `json:"enabled"`
}

type UpdateModelBody struct {
	Provider       string         `json:"provider" required:"true"`
	ModelName      string         `json:"model_name" required:"true"`
	DeploymentName *string        `json:"deployment_name,omitempty"`
	EndpointURL    string         `json:"endpoint_url" required:"true"`
	AuthType       AuthType       `json:"auth_type" required:"true"`
	AuthConfig     AuthConfig     `json:"auth_config" required:"true"`
	Metadata       map[string]any `json:"metadata,omitempty"`
	Enabled        bool           `json:"enabled"`
}

type CreateModelRequest struct {
	Body CreateModelBody `json:"body"`
}

type CreateModelResponse struct {
	Body *Model `json:"body"`
}

type ListModelsResponse struct {
	Body []*Model `json:"body"`
}

type ListEnabledModelsResponse struct {
	Body []*Model `json:"body"`
}

type GetModelResponse struct {
	Body *Model `json:"body"`
}

type UpdateModelRequest struct {
	ID   string          `path:"id" required:"true"`
	Body UpdateModelBody `json:"body"`
}

type UpdateModelResponse struct {
	Body *Model `json:"body"`
}

type ListModelsRequest struct {
	model.ListRequest
}

type ListEnabledModelsRequest struct {
	model.ListRequest
}
