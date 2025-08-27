package model

import (
	"github.com/golang-jwt/jwt/v5"
)

const (
	RequestIDKey contextKey = "requestId"
)

type contextKey string

type ScopedToken struct {
	jwt.RegisteredClaims

	Email             string   `json:"email"`
	Name              string   `json:"name,omitempty"`
	OrgID             string   `json:"org_id,omitempty"` // this is hydrated post auth login from our db
	UserID            string   `json:"id,omitempty"`
	GivenName         string   `json:"given_name,omitempty"`
	FamilyName        string   `json:"family_name,omitempty"`
	PreferredUsername string   `json:"preferred_username,omitempty"`
	Roles             []string `json:"roles,omitempty"`
	Groups            []string `json:"groups,omitempty"`
}
