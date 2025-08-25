package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	UserClaimsKey           contextKey = "userClaims"
	RequestIDKey            contextKey = "requestId"
	DefaultOrganisationRole string     = "admin"
	MSALClaimsKey           contextKey = "msalClaims"
)

type contextKey string

type OrganisationMembership struct {
	OrganisationID uuid.UUID `json:"organisation_id"`
	Roles          []string  `json:"roles"`
}

type UserClaims struct {
	UserID        uuid.UUID                `json:"user_id"`
	Email         string                   `json:"email"`
	Organisations []OrganisationMembership `json:"organisations"`
}

type ScopedTokenClaims struct {
	jwt.RegisteredClaims

	Email         string                   `json:"email"`
	Organisation  *OrganisationMembership  `json:"org,omitempty"`  // present if scoped
	Organisations []OrganisationMembership `json:"orgs,omitempty"` // present if user must choose
}

type MSALClaims struct {
	Email string `json:"email"`
	Sub   string `json:"sub"`
	Aud   string `json:"aud"`
	Iss   string `json:"iss"`
	Exp   int64  `json:"exp"`
	Iat   int64  `json:"iat"`
	jwt.RegisteredClaims
	// TODO: Add more fields as needed from the MSAL JWT
}
