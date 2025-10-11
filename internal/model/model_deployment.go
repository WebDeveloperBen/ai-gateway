package model

import (
	"time"

	"github.com/google/uuid"
)

// ModelDeploymentDB represents the database model for model deployments
type ModelDeploymentDB struct {
	ID             uuid.UUID
	OrgID          uuid.UUID
	Provider       string
	ModelName      string
	DeploymentName *string
	EndpointURL    string
	AuthType       string                 // managed_identity, api_key, secret_ref
	AuthConfig     map[string]interface{} // JSON blob
	Metadata       map[string]interface{} // JSON blob
	Enabled        bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ModelDeployment is the current struct used by the gateway/registry
// Keeping original field names for backward compatibility
type ModelDeployment struct {
	Model      string            `json:"model"`
	Deployment string            `json:"deployment"`
	Provider   string            `json:"provider"`
	Tenant     string            `json:"tenant"`
	Meta       map[string]string `json:"meta,omitempty"`
}
