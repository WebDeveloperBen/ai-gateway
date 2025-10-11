package keys

import (
	"context"
	"errors"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store is a concrete Postgres implementation of keys.Store.
type store struct {
	pool *pgxpool.Pool
}

var _ KeyRepository = (*store)(nil)

func NewPostgresStore(pool *pgxpool.Pool) *store {
	return &store{pool: pool}
}

// Insert stores key metadata and its PHC hash.
func (s *store) Insert(ctx context.Context, k model.Key, phc string) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO api_keys
			(key_id, secret_phc, tenant, app, status, last_four, metadata, expires_at)
		VALUES
			($1,     $2,         $3,     $4,  $5,     $6,       COALESCE($7,'{}'::jsonb), $8)
	`,
		k.KeyID, phc, k.Tenant, k.App, string(k.Status), k.LastFour, k.Metadata, k.ExpiresAt,
	)
	return err
}

// GetByKeyID returns key metadata (no PHC).
func (s *store) GetByKeyID(ctx context.Context, keyID string) (*model.Key, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT key_id, tenant, app, status, expires_at, last_used_at,
		       COALESCE(metadata,'{}'::jsonb), created_at, COALESCE(last_four,'')
		FROM api_keys
		WHERE key_id = $1
	`, keyID)

	var (
		k        model.Key
		status   string
		metadata []byte
		lastFour string
	)
	err := row.Scan(
		&k.KeyID, &k.Tenant, &k.App, &status, &k.ExpiresAt, &k.LastUsedAt,
		&metadata, &k.CreatedAt, &lastFour,
	)
	if err != nil {
		return nil, err
	}
	k.Status = model.KeyStatus(status)
	k.Metadata = metadata
	k.LastFour = lastFour
	return &k, nil
}

// GetPHCByKeyID fetches the PHC (argon2id) string for verification.
func (s *store) GetPHCByKeyID(ctx context.Context, keyID string) (string, error) {
	var phc string
	if err := s.pool.QueryRow(ctx, `
		SELECT secret_phc FROM api_keys WHERE key_id = $1
	`, keyID).Scan(&phc); err != nil {
		return "", err
	}
	return phc, nil
}

// TouchLastUsed updates last_used_at; best-effort.
func (s *store) TouchLastUsed(ctx context.Context, keyID string) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE api_keys SET last_used_at = now() WHERE key_id = $1
	`, keyID)
	return err
}

// UpdateStatus sets the key status (active/revoked/expired).
func (s *store) UpdateStatus(ctx context.Context, keyID string, status model.KeyStatus) error {
	switch status {
	case model.KeyActive, model.KeyRevoked, model.KeyExpired:
	default:
		return errors.New("invalid status")
	}
	_, err := s.pool.Exec(ctx, `
		UPDATE api_keys SET status = $2 WHERE key_id = $1
	`, keyID, string(status))
	return err
}
