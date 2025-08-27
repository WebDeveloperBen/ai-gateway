package db

import (
	"context"
	"fmt"

	"github.com/insurgence-ai/llm-gateway/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}

func NewPostgresDriver(ctx context.Context, dsn string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("[Fail]: Postgres driver unable to connect: %+w", err)
	}
	return &Postgres{Pool: pool, Queries: db.New(pool)}, nil
}
