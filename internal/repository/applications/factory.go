package applications

import (
	"context"
	"fmt"

	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

func NewRepository(ctx context.Context, cfg model.RepositoryConfig) (Repository, error) {
	switch cfg.Backend {
	case model.RepositoryPostgres:
		queries := db.New(cfg.PGPool)
		return NewPostgresRepo(queries), nil
	case model.RepositoryMemory:
		return NewMemoryRepo(), nil
	default:
		return nil, fmt.Errorf("unsupported applications repository backend: %s", cfg.Backend)
	}
}
