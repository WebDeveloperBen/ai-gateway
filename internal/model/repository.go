package model

import "github.com/jackc/pgx/v5/pgxpool"

type RepositoryBackend string

const (
	RepositoryPostgres RepositoryBackend = "postgres"
	RepositoryMemory   RepositoryBackend = "memory"
)

type RepositoryConfig struct {
	Backend RepositoryBackend
	PGPool  *pgxpool.Pool
}
