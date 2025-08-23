package postgres

import (
	"context"
	"errors"

	"github.com/insurgence-ai/llm-gateway/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIKeyStore struct {
	pool *pgxpool.Pool
}

func NewAPIKeyStore(pool *pgxpool.Pool) *APIKeyStore { return &APIKeyStore{pool: pool} }

func (s *APIKeyStore) Insert(ctx context.Context, k store.APIKey) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO api_keys (key_id, secret_phc, tenant, app, status, last_four, metadata, expires_at)
		VALUES ($1,$2,$3,$4,$5,$6,COALESCE($7,'{}'::jsonb),$8)
	`, k.KeyID, k.PHC, k.Tenant, k.App, k.Status, k.LastFour, k.Metadata, k.ExpiresAt)
	return err
}

func (s *APIKeyStore) GetByKeyID(ctx context.Context, keyID string) (*store.APIKey, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT key_id, secret_phc, tenant, app, status, expires_at, last_used_at, COALESCE(metadata,'{}'::jsonb), COALESCE(last_four,'')
		FROM api_keys WHERE key_id = $1
	`, keyID)

	var k store.APIKey
	var meta []byte
	var lastFour string
	if err := row.Scan(&k.KeyID, &k.PHC, &k.Tenant, &k.App, &k.Status, &k.ExpiresAt, &k.LastUsedAt, &meta, &lastFour); err != nil {
		return nil, err
	}
	k.Metadata = meta
	k.LastFour = lastFour
	return &k, nil
}

func (s *APIKeyStore) TouchLastUsed(ctx context.Context, keyID string) error {
	_, err := s.pool.Exec(ctx, `UPDATE api_keys SET last_used_at = now() WHERE key_id = $1`, keyID)
	return err
}

func (s *APIKeyStore) UpdateStatus(ctx context.Context, keyID, status string) error {
	if status != "active" && status != "revoked" && status != "expired" {
		return errors.New("invalid status")
	}
	_, err := s.pool.Exec(ctx, `UPDATE api_keys SET status = $2 WHERE key_id = $1`, keyID, status)
	return err
}
