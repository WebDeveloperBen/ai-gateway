package applications

import (
	"time"
)

// ===================================
// |      Data Transfer Objects      |
// ===================================

type Application struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"org_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"-"`
}

// ===================================
// |         Request Models          |
// ===================================

// Create Application

type CreateApplicationBody struct {
	Name        string  `json:"name" required:"true"`
	Description *string `json:"description,omitempty"`
}
type CreateApplicationRequest struct {
	Body CreateApplicationBody `json:"body"`
}
type CreateApplicationResponse struct {
	Body *Application `json:"body"`
}

// List Application

type ListApplicationsResponse struct {
	Body []*Application `json:"body"`
}

// Get Application

type GetApplicationResponse struct {
	Body *Application `json:"body"`
}

// Update Application

type UpdateApplicationBody struct {
	Name        string  `json:"name" required:"true"`
	Description *string `json:"description,omitempty"`
}
type UpdateApplicationRequest struct {
	ID   string                `path:"id" required:"true"`
	Body UpdateApplicationBody `json:"body"`
}
type UpdateApplicationResponse struct {
	Body *Application `json:"body"`
}
