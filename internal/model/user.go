package model

import "time"

type User struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	Sub       string    `json:"sub"`
	Email     string    `json:"email"`
	Name      string    `json:"name,omitempty"`
	Roles     []string  `json:"roles,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
