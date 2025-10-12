package keys

import (
	"context"
	"errors"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type store struct {
	queries *db.Queries
}

var _ KeyRepository = (*store)(nil)

func NewPostgresStore(queries *db.Queries) *store {
	return &store{queries: queries}
}

func (s *store) Insert(ctx context.Context, k model.Key, phc string) error {
	var expiresAt pgtype.Timestamptz
	if k.ExpiresAt != nil {
		expiresAt = pgtype.Timestamptz{Time: *k.ExpiresAt, Valid: true}
	}

	_, err := s.queries.InsertAPIKey(ctx, db.InsertAPIKeyParams{
		OrgID:      k.OrgID,
		AppID:      k.AppID,
		UserID:     k.UserID,
		KeyPrefix:  k.KeyPrefix,
		SecretPhc:  phc,
		Status:     string(k.Status),
		LastFour:   k.LastFour,
		ExpiresAt:  expiresAt,
		Metadata:   k.Metadata,
	})
	return err
}

func (s *store) GetByKeyPrefix(ctx context.Context, keyPrefix string) (*model.Key, error) {
	row, err := s.queries.GetAPIKeyByPrefix(ctx, keyPrefix)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("key not found")
		}
		return nil, err
	}

	var expiresAt *time.Time
	if row.ExpiresAt.Valid {
		expiresAt = &row.ExpiresAt.Time
	}

	var lastUsedAt *time.Time
	if row.LastUsedAt.Valid {
		lastUsedAt = &row.LastUsedAt.Time
	}

	return &model.Key{
		ID:         row.ID,
		OrgID:      row.OrgID,
		AppID:      row.AppID,
		UserID:     row.UserID,
		KeyPrefix:  row.KeyPrefix,
		Status:     model.KeyStatus(row.Status),
		LastFour:   row.LastFour,
		ExpiresAt:  expiresAt,
		LastUsedAt: lastUsedAt,
		Metadata:   row.Metadata,
		CreatedAt:  row.CreatedAt.Time,
	}, nil
}

func (s *store) GetSecretPHCByPrefix(ctx context.Context, keyPrefix string) (string, error) {
	phc, err := s.queries.GetSecretPHCByPrefix(ctx, keyPrefix)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("key not found")
		}
		return "", err
	}
	return phc, nil
}

func (s *store) TouchLastUsed(ctx context.Context, keyPrefix string) error {
	return s.queries.UpdateAPIKeyLastUsed(ctx, keyPrefix)
}

func (s *store) UpdateStatus(ctx context.Context, keyPrefix string, status model.KeyStatus) error {
	switch status {
	case model.KeyActive, model.KeyRevoked, model.KeyExpired:
	default:
		return errors.New("invalid status")
	}
	return s.queries.UpdateAPIKeyStatus(ctx, db.UpdateAPIKeyStatusParams{
		KeyPrefix: keyPrefix,
		Status:    string(status),
	})
}

func (s *store) Delete(ctx context.Context, id uuid.UUID) error {
	return s.queries.DeleteAPIKey(ctx, id)
}
