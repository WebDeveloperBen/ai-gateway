package model

import (
	"time"
)

type KeyStatus string

const (
	KeyActive  KeyStatus = "active"
	KeyRevoked KeyStatus = "revoked"
	KeyExpired KeyStatus = "expired"
)

type Key struct {
	KeyID      string
	Tenant     string
	App        string
	Status     KeyStatus
	ExpiresAt  *time.Time
	LastUsedAt *time.Time
	LastFour   string
	Metadata   []byte // raw JSON
	CreatedAt  time.Time
}
