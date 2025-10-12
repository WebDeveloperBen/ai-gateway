package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/repository/keys"
)

// NewAPIKeyAuthenticator Constructor with explicit dependencies.
func NewAPIKeyAuthenticator(store keys.Reader, hasher keys.Hasher) *APIKeyAuthenticator {
	return &APIKeyAuthenticator{Keys: store, Hasher: hasher}
}

// NewDefaultAPIKeyAuthenticator is a convenience constructor with tuned defaults
func NewDefaultAPIKeyAuthenticator(store keys.Reader) *APIKeyAuthenticator {
	hasher := keys.NewArgon2IDHasher(1, 64*1024, 1, 32)
	return &APIKeyAuthenticator{Keys: store, Hasher: hasher}
}

func (a *APIKeyAuthenticator) Authenticate(r *http.Request) (keyID string, keyData *KeyData, err error) {
	tok := getHeaderToken(r)
	keyID, secret := splitToken(tok)
	if keyID == "" || secret == "" {
		_ = padWork(a.Hasher)
		return "", nil, errors.New("unauthorized")
	}

	rec, err := a.Keys.GetByKeyPrefix(r.Context(), keyID)
	if err != nil || rec == nil {
		_ = padWork(a.Hasher)
		return "", nil, errors.New("unauthorized")
	}

	if rec.Status != model.KeyActive {
		return "", nil, errors.New("unauthorized")
	}
	if rec.ExpiresAt != nil && time.Now().After(*rec.ExpiresAt) {
		return "", nil, errors.New("unauthorized")
	}

	phc, err := a.Keys.GetSecretPHCByPrefix(r.Context(), keyID)
	if err != nil {
		_ = padWork(a.Hasher)
		return "", nil, errors.New("unauthorized")
	}
	ok, _ := a.Hasher.Verify(phc, []byte(secret))
	if !ok {
		return "", nil, errors.New("unauthorized")
	}

	_ = a.Keys.TouchLastUsed(r.Context(), keyID)

	// Build KeyData from record
	data := &KeyData{
		KeyID:  keyID,
		OrgID:  rec.OrgID.String(),
		AppID:  rec.AppID.String(),
		UserID: rec.UserID.String(),
	}

	return keyID, data, nil
}
