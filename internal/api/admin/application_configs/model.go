package application_configs

import (
	"time"
)

type ApplicationConfig struct {
	ID          string         `json:"id"`
	AppID       string         `json:"app_id"`
	OrgID       string         `json:"org_id"`
	Environment string         `json:"environment"`
	Config      map[string]any `json:"config"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"-"`
}

type CreateApplicationConfigBody struct {
	AppID       string         `json:"app_id" required:"true"`
	Environment string         `json:"environment" required:"true"`
	Config      map[string]any `json:"config" required:"true"`
}
type CreateApplicationConfigRequest struct {
	Body CreateApplicationConfigBody `json:"body"`
}
type CreateApplicationConfigResponse struct {
	Body *ApplicationConfig `json:"body"`
}

type ListApplicationConfigsResponse struct {
	Body []*ApplicationConfig `json:"body"`
}

type GetApplicationConfigResponse struct {
	Body *ApplicationConfig `json:"body"`
}

type UpdateApplicationConfigBody struct {
	Config map[string]any `json:"config" required:"true"`
}
type UpdateApplicationConfigRequest struct {
	ID   string                      `path:"id" required:"true"`
	Body UpdateApplicationConfigBody `json:"body"`
}
type UpdateApplicationConfigResponse struct {
	Body *ApplicationConfig `json:"body"`
}
