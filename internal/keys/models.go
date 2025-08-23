package keys

import "time"

type Status string

const (
	Active  Status = "active"
	Revoked Status = "revoked"
	Expired Status = "expired"
)

type Key struct {
	KeyID      string
	Tenant     string
	App        string
	Status     Status
	ExpiresAt  *time.Time
	LastUsedAt *time.Time
	LastFour   string
	Metadata   []byte // raw JSON
	CreatedAt  time.Time
}
