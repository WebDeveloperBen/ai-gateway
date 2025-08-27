package middleware

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/insurgence-ai/llm-gateway/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrgQueries struct {
	*db.Queries
	tx   pgx.Tx
	conn *pgxpool.Conn
}

func NewOrgQueries(ctx context.Context, pool *pgxpool.Pool, orgID string) (*OrgQueries, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		conn.Release()
		return nil, err
	}

	// Enforce org scoping for the whole tx
	if _, err := tx.Exec(ctx, "SET LOCAL app.current_org = $1", orgID); err != nil {
		tx.Rollback(ctx)
		conn.Release()
		return nil, err
	}

	return &OrgQueries{
		Queries: db.New(tx),
		tx:      tx,
		conn:    conn,
	}, nil
}

func (oq *OrgQueries) Close(ctx context.Context) error {
	defer oq.conn.Release()
	return oq.tx.Commit(ctx)
}

func WithScopedOrg(pool *pgxpool.Pool) func(huma.API) func(huma.Context, func(huma.Context)) {
	return func(api huma.API) func(huma.Context, func(huma.Context)) {
		return func(ctx huma.Context, next func(huma.Context)) {
			orgID, ok := GetOrgIDFromSession(ctx.Context())
			if !ok {
				huma.WriteErr(api, ctx, http.StatusForbidden, "missing organisation context")
				return
			}

			oq, err := NewOrgQueries(ctx.Context(), pool, orgID)
			if err != nil {
				huma.WriteErr(api, ctx, http.StatusBadRequest, "something went wrong")
				return
			}
			defer oq.Close(ctx.Context())

			c := huma.WithValue(ctx, orgQueriesKey, oq)
			next(c)
		}
	}
}
