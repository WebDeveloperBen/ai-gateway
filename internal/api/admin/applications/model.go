package applications

import (
	"time"
)

type Application struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"org_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateApplicationRequest struct {
	Name        string  `json:"name" required:"true"`
	Description *string `json:"description,omitempty"`
}

type UpdateApplicationRequest struct {
	Name        string  `json:"name" required:"true"`
	Description *string `json:"description,omitempty"`
}

type ApplicationList struct {
	Applications []*Application `json:"applications"`
}
