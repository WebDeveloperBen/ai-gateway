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

type CreatePolicyRequest struct {
	AppID      string                 `json:"app_id" required:"true"`
	PolicyType model.PolicyType       `json:"policy_type" required:"true"`
	Config     map[string]interface{} `json:"config" required:"true"`
	Enabled    bool                   `json:"enabled"`
}

type UpdatePolicyRequest struct {
	PolicyType model.PolicyType       `json:"policy_type" required:"true"`
	Config     map[string]interface{} `json:"config" required:"true"`
	Enabled    bool                   `json:"enabled"`
}

type PolicyList struct {
	Policies []*Policy `json:"policies"`
}
