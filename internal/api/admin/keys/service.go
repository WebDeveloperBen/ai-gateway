package keys

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/insurgence-ai/llm-gateway/internal/model"
	"github.com/insurgence-ai/llm-gateway/internal/repository/keys"
)

type KeysService interface {
	MintKey(ctx context.Context, req MintKeyRequest) (MintKeyResponse, error)
	RevokeKey(ctx context.Context, keyID string) error
	GetByKeyID(ctx context.Context, keyID string) (APIKey, error)
}

type keysService struct {
	store  keys.KeyRepository
	hasher keys.Hasher
}

func NewService(store keys.KeyRepository, hasher keys.Hasher) KeysService {
	return &keysService{store: store, hasher: hasher}
}

func (s *keysService) MintKey(ctx context.Context, req MintKeyRequest) (MintKeyResponse, error) {
	if req.Tenant == "" || req.App == "" {
		return MintKeyResponse{}, errors.New("tenant & app required")
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
		return MintKeyResponse{}, err
	}
	secretB64 := base64.RawStdEncoding.EncodeToString(secret)

	// PHC (random salt is generated inside Hash)
	phc, err := s.hasher.Hash(secret)
	if err != nil {
		return MintKeyResponse{}, err
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
	k := model.Key{
		KeyID:     keyID,
		Tenant:    req.Tenant,
		App:       req.App,
		Status:    model.KeyActive,
		ExpiresAt: exp,
		LastFour:  last4,
		Metadata:  []byte("{}"),
		// CreatedAt is filled by DB default; ok to leave zero here
	}
	if err := s.store.Insert(ctx, k, phc); err != nil {
		return MintKeyResponse{}, err
	}

	return MintKeyResponse{
		Token: keyID + "." + secretB64, // show ONCE
		Key: APIKey{
			KeyID:     keyID,
			Tenant:    req.Tenant,
			App:       req.App,
			Status:    model.KeyActive,
			ExpiresAt: exp,
			LastFour:  last4,
		},
	}, nil
}

func (s *keysService) RevokeKey(ctx context.Context, keyID string) error {
	return s.store.UpdateStatus(ctx, keyID, model.KeyRevoked)
}

func (s *keysService) GetByKeyID(ctx context.Context, keyID string) (APIKey, error) {
	rec, err := s.store.GetByKeyID(ctx, keyID)
	if err != nil || rec == nil {
		return APIKey{}, errors.New("not found")
	}
	return APIKey{
		KeyID:      rec.KeyID,
		Tenant:     rec.Tenant,
		App:        rec.App,
		Status:     model.KeyStatus(rec.Status),
		ExpiresAt:  rec.ExpiresAt,
		LastUsedAt: rec.LastUsedAt,
		LastFour:   rec.LastFour,
		// Metadata: map if you expose it in admin API
	}, nil
}
