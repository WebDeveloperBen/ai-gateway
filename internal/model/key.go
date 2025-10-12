package model

import (
	"time"

	"github.com/google/uuid"
)

type KeyStatus string

const (
	KeyActive  KeyStatus = "active"
	KeyRevoked KeyStatus = "revoked"
	KeyExpired KeyStatus = "expired"
)

type Key struct {
	ID         uuid.UUID
	OrgID      uuid.UUID
	AppID      uuid.UUID
	UserID     uuid.UUID
	KeyPrefix  string
	Status     KeyStatus
	LastFour   string
	ExpiresAt  *time.Time
	LastUsedAt *time.Time
	Metadata   []byte // raw JSON
	CreatedAt  time.Time
}
