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

type ListRequest struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

func NormalizePagination(r ListRequest) ListRequest {
	if r.Limit <= 0 {
		r.Limit = 100
	}
	if r.Offset < 0 {
		r.Offset = 0
	}
	return r
}
