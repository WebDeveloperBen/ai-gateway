package lib

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(connStr *string) (*pgxpool.Pool, error) {
	conn := os.Getenv("DATABASE_URL")
	if connStr != nil {
		conn = *connStr
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, err
	}

	return pgxpool.NewWithConfig(ctx, cfg)
}
