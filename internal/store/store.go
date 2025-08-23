package store

import (
	"context"
	"time"
)

type APIKey struct {
	KeyID      string
	PHC        string
	Tenant     string
	App        string
	Status     string
	ExpiresAt  *time.Time
	LastUsedAt *time.Time
	LastFour   string
	Metadata   []byte // raw JSON
}

type APIKeyStore interface {
	// Insert a new key (already hashed). Returns nothing; caller returns the clear token to the admin once.
	Insert(ctx context.Context, k APIKey) error

	// Lookup by key_id (left side of token). Returns nil if not found.
	GetByKeyID(ctx context.Context, keyID string) (*APIKey, error)

	// Mark a key as used (best-effort; don’t block a hot path on errors).
	TouchLastUsed(ctx context.Context, keyID string) error

	// Revoke or update status.
	UpdateStatus(ctx context.Context, keyID, status string) error
}
