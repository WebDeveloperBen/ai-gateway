package models

import "time"

type KeyStatus string

const (
	KeyActive  KeyStatus = "active"
	KeyRevoked KeyStatus = "revoked"
	KeyExpired KeyStatus = "expired"
)

type APIKey struct {
	KeyID      string
	Tenant     string
	App        string
	Status     KeyStatus
	ExpiresAt  *time.Time
	LastUsedAt *time.Time
	LastFour   string
	Metadata   map[string]any
	CreatedAt  time.Time
}

type MintKeyRequest struct {
	Tenant   string
	App      string
	TTL      time.Duration
	Prefix   string
	Metadata map[string]any
}

type MintKeyResponse struct {
	Token string // key_id.secret (display once)
	Key   APIKey
}
