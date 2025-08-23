package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/insurgence-ai/llm-gateway/internal/admin/models"
	"github.com/insurgence-ai/llm-gateway/internal/keys"
)

type KeysService interface {
	MintKey(ctx context.Context, req models.MintKeyRequest) (models.MintKeyResponse, error)
	RevokeKey(ctx context.Context, keyID string) error
	GetByKeyID(ctx context.Context, keyID string) (models.APIKey, error)
}

type keysService struct {
	store  keys.Store  // shared keys store (pg impl under internal/keys/postgres)
	hasher keys.Hasher // Argon2ID hasher (internal/keys.NewArgon2IDHasher)
}

func NewKeysService(store keys.Store, hasher keys.Hasher) KeysService {
	return &keysService{store: store, hasher: hasher}
}

func (s *keysService) MintKey(ctx context.Context, req models.MintKeyRequest) (models.MintKeyResponse, error) {
	if req.Tenant == "" || req.App == "" {
		return models.MintKeyResponse{}, errors.New("tenant & app required")
	}
	prefix := req.Prefix
	if prefix == "" {
		prefix = "sk_live"
	}

	// id + secret
	id := ulid.Make().String()
	keyID := prefix + "_" + id

	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return models.MintKeyResponse{}, err
	}
	secretB64 := base64.RawStdEncoding.EncodeToString(secret)

	// PHC (random salt is generated inside Hash)
	phc, err := s.hasher.Hash(secret)
	if err != nil {
		return models.MintKeyResponse{}, err
	}

	// optional expiry
	var exp *time.Time
	if req.TTL > 0 {
		t := time.Now().Add(req.TTL)
		exp = &t
	}

	last4 := ""
	if len(secretB64) >= 4 {
		last4 = secretB64[len(secretB64)-4:]
	}

	// persist metadata via shared keys.Store
	k := keys.Key{
		KeyID:     keyID,
		Tenant:    req.Tenant,
		App:       req.App,
		Status:    keys.Active,
		ExpiresAt: exp,
		LastFour:  last4,
		Metadata:  []byte("{}"),
		// CreatedAt is filled by DB default; ok to leave zero here
	}
	if err := s.store.Insert(ctx, k, phc); err != nil {
		return models.MintKeyResponse{}, err
	}

	return models.MintKeyResponse{
		Token: keyID + "." + secretB64, // show ONCE
		Key: models.APIKey{
			KeyID:     keyID,
			Tenant:    req.Tenant,
			App:       req.App,
			Status:    models.KeyActive,
			ExpiresAt: exp,
			LastFour:  last4,
		},
	}, nil
}

func (s *keysService) RevokeKey(ctx context.Context, keyID string) error {
	return s.store.UpdateStatus(ctx, keyID, keys.Revoked)
}

func (s *keysService) GetByKeyID(ctx context.Context, keyID string) (models.APIKey, error) {
	rec, err := s.store.GetByKeyID(ctx, keyID)
	if err != nil || rec == nil {
		return models.APIKey{}, errors.New("not found")
	}
	return models.APIKey{
		KeyID:      rec.KeyID,
		Tenant:     rec.Tenant,
		App:        rec.App,
		Status:     models.KeyStatus(rec.Status),
		ExpiresAt:  rec.ExpiresAt,
		LastUsedAt: rec.LastUsedAt,
		LastFour:   rec.LastFour,
		// Metadata: map if you expose it in admin API
	}, nil
}
