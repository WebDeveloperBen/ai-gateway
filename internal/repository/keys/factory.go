package keys

import (
	"context"
	"fmt"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

func NewKeyRepository(ctx context.Context, cfg model.RepositoryConfig) (KeyRepository, error) {
	switch cfg.Backend {
	case model.RepositoryPostgres:
		return NewPostgresStore(cfg.PGPool), nil
	case model.RepositoryMemory:
		return NewMemoryStore(), nil
	default:
		return nil, fmt.Errorf("unsupported key repository backend: %s", cfg.Backend)
	}
}
