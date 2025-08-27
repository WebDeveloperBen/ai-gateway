package middleware

import (
	"context"
	"net/http"
	"yourapp/db"

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

func WithOrg(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			orgID, _ := GetOrgIDFromSession(r.Context())

			oq, err := NewOrgQueries(r.Context(), pool, orgID)
			if err != nil {
				http.Error(w, "db error", http.StatusInternalServerError)
				return
			}
			defer oq.Close(r.Context())

			ctx := context.WithValue(r.Context(), orgQueriesKey, oq)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
