// Package keys manages API keys and credentials used for authenticating
// applications via the gateway.
package keys

import (
	"context"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
)

type Reader interface {
	GetByKeyPrefix(ctx context.Context, keyPrefix string) (*model.Key, error)
	GetSecretPHCByPrefix(ctx context.Context, keyPrefix string) (string, error)
	TouchLastUsed(ctx context.Context, keyPrefix string) error
}

type Writer interface {
	Insert(ctx context.Context, k model.Key, phc string) error
	UpdateStatus(ctx context.Context, keyPrefix string, status model.KeyStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type KeyRepository interface {
	Reader
	Writer
}
