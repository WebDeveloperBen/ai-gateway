package keys

import "context"

type Reader interface {
	GetByKeyID(ctx context.Context, keyID string) (*Key, error)
	TouchLastUsed(ctx context.Context, keyID string) error
	GetPHCByKeyID(ctx context.Context, keyID string) (string, error)
}

type Writer interface {
	Insert(ctx context.Context, k Key, phc string) error
	UpdateStatus(ctx context.Context, keyID string, status Status) error
}

type Store interface {
	Reader
	Writer
}
