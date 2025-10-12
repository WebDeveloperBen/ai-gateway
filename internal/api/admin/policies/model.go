package policies

import (
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

type Policy struct {
	ID         string           `json:"id"`
	OrgID      string           `json:"org_id"`
	PolicyType model.PolicyType `json:"policy_type"`
	Config     map[string]any   `json:"config"`
	Enabled    bool             `json:"enabled"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

type CreatePolicyBody struct {
	OrgID      string           `json:"org_id" required:"true"`
	AppID      string           `json:"app_id"` // Optional - policies can be created without being attached to apps
	PolicyType model.PolicyType `json:"policy_type" required:"true"`
	Config     map[string]any   `json:"config" required:"true"`
	Enabled    bool             `json:"enabled"`
}

type UpdatePolicyBody struct {
	AppID      string           `json:"app_id"` // Allow attaching policy to app
	PolicyType model.PolicyType `json:"policy_type" required:"true"`
	Config     map[string]any   `json:"config" required:"true"`
	Enabled    bool             `json:"enabled"`
}

type CreatePolicyRequest struct {
	Body CreatePolicyBody `json:"body"`
}

type CreatePolicyResponse struct {
	Body *Policy `json:"body"`
}

type ListPoliciesResponse struct {
	Body []*Policy `json:"body"`
}

type ListEnabledPoliciesResponse struct {
	Body []*Policy `json:"body"`
}

type GetPolicyResponse struct {
	Body *Policy `json:"body"`
}

type UpdatePolicyRequest struct {
	ID   string           `path:"id" required:"true"`
	Body UpdatePolicyBody `json:"body"`
}

type UpdatePolicyResponse struct {
	Body *Policy `json:"body"`
}

// Huma request types for list endpoints with pagination
type ListPoliciesRequest struct {
	AppID string `query:"app_id" required:"true"`
	model.ListRequest
}

type ListEnabledPoliciesRequest struct {
	AppID string `query:"app_id" required:"true"`
	model.ListRequest
}
