package policies

import (
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

type Policy struct {
	ID         string                 `json:"id"`
	OrgID      string                 `json:"org_id"`
	AppID      string                 `json:"app_id"`
	PolicyType model.PolicyType       `json:"policy_type"`
	Config     map[string]interface{} `json:"config"`
	Enabled    bool                   `json:"enabled"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// Request/Response body types
type CreatePolicyBody struct {
	AppID      string                 `json:"app_id" required:"true"`
	PolicyType model.PolicyType       `json:"policy_type" required:"true"`
	Config     map[string]interface{} `json:"config" required:"true"`
	Enabled    bool                   `json:"enabled"`
}

type UpdatePolicyBody struct {
	PolicyType model.PolicyType       `json:"policy_type" required:"true"`
	Config     map[string]interface{} `json:"config" required:"true"`
	Enabled    bool                   `json:"enabled"`
}

// Huma request/response types
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
