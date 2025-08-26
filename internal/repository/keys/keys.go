// Package keys manages API keys and credentials used for authenticating
// tenants and for authenticating against upstream providers.
package keys

import (
	"context"
	"github.com/insurgence-ai/llm-gateway/internal/model"
)

type Reader interface {
	GetByKeyID(ctx context.Context, keyID string) (*model.Key, error)
	TouchLastUsed(ctx context.Context, keyID string) error
	GetPHCByKeyID(ctx context.Context, keyID string) (string, error)
}

type Writer interface {
	Insert(ctx context.Context, k model.Key, phc string) error
	UpdateStatus(ctx context.Context, keyID string, status model.KeyStatus) error
}

type KeyRepository interface {
	Reader
	Writer
}
