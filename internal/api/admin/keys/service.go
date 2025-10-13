package keys

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/repository/keys"
)

type KeysService interface {
	MintKey(ctx context.Context, req MintKeyRequestBody) (MintKeyResponse, error)
	RevokeKey(ctx context.Context, keyPrefix string) error
	GetByKeyPrefix(ctx context.Context, keyPrefix string) (APIKey, error)
}

type keysService struct {
	store  keys.KeyRepository
	hasher keys.Hasher
}

func NewService(store keys.KeyRepository, hasher keys.Hasher) KeysService {
	return &keysService{store: store, hasher: hasher}
}

func (s *keysService) MintKey(ctx context.Context, req MintKeyRequestBody) (MintKeyResponse, error) {
	if req.OrgID == "" || req.AppID == "" || req.UserID == "" {
		return MintKeyResponse{}, errors.New("orgID, appID & userID required")
	}

	orgID, err := uuid.Parse(req.OrgID)
	if err != nil {
		return MintKeyResponse{}, errors.New("invalid orgID")
	}

	appID, err := uuid.Parse(req.AppID)
	if err != nil {
		return MintKeyResponse{}, errors.New("invalid appID")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return MintKeyResponse{}, errors.New("invalid userID")
	}

	prefix := req.Prefix
	if prefix == "" {
		prefix = "sk_live"
	}

	id := ulid.Make().String()
	keyPrefix := prefix + "_" + id

	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return MintKeyResponse{}, err
	}
	secretB64 := base64.RawStdEncoding.EncodeToString(secret)

	phc, err := s.hasher.Hash(secret)
	if err != nil {
		return MintKeyResponse{}, err
	}

	var exp *time.Time
	if req.TTL > 0 {
		t := time.Now().Add(req.TTL)
		exp = &t
	}

	last4 := ""
	if len(secretB64) >= 4 {
		last4 = secretB64[len(secretB64)-4:]
	}

	metadata := []byte("{}")
	if req.Metadata != nil {
		metadata, err = json.Marshal(req.Metadata)
		if err != nil {
			return MintKeyResponse{}, err
		}
	}

	k := model.Key{
		OrgID:     orgID,
		AppID:     appID,
		UserID:    userID,
		KeyPrefix: keyPrefix,
		Status:    model.KeyActive,
		ExpiresAt: exp,
		LastFour:  last4,
		Metadata:  metadata,
	}
	if err := s.store.Insert(ctx, k, phc); err != nil {
		return MintKeyResponse{}, err
	}

	return MintKeyResponse{
		Body: MintKeyResponseBody{
			Token: keyPrefix + "." + secretB64,
			Key: APIKey{
				KeyID:     keyPrefix,
				Tenant:    req.OrgID,
				App:       req.AppID,
				Status:    model.KeyActive,
				ExpiresAt: exp,
				LastFour:  last4,
			},
		},
	}, nil
}

func (s *keysService) RevokeKey(ctx context.Context, keyPrefix string) error {
	err := s.store.UpdateStatus(ctx, keyPrefix, model.KeyRevoked)
	if err != nil && err.Error() == "key not found" {
		// Revoking a non-existent key is a successful no-op
		return nil
	}
	return err
}

func (s *keysService) GetByKeyPrefix(ctx context.Context, keyPrefix string) (APIKey, error) {
	rec, err := s.store.GetByKeyPrefix(ctx, keyPrefix)
	if err != nil || rec == nil {
		return APIKey{}, errors.New("not found")
	}
	return APIKey{
		KeyID:      rec.KeyPrefix,
		Tenant:     rec.OrgID.String(),
		App:        rec.AppID.String(),
		Status:     rec.Status,
		ExpiresAt:  rec.ExpiresAt,
		LastUsedAt: rec.LastUsedAt,
		LastFour:   rec.LastFour,
	}, nil
}
