package model

type ModelDeployment struct {
	Model      string            `json:"model"`
	Deployment string            `json:"deployment"`
	Provider   string            `json:"provider"`
	Tenant     string            `json:"tenant"`
	Meta       map[string]string `json:"meta,omitempty"`
}
