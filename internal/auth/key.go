// Package auth implements authentication and tenant identification
// for incoming requests, attaching tenant and app context to the request.
package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/insurgence-ai/llm-gateway/internal/keys"
)

type APIKeyAuthenticator struct {
	Keys   keys.Reader // GetByKeyID, GetPHCByKeyID, TouchLastUsed
	Hasher keys.Hasher // Argon2IDHasher
}

// NewAPIKeyAuthenticator Constructor with explicit dependencies.
func NewAPIKeyAuthenticator(store keys.Reader, hasher keys.Hasher) *APIKeyAuthenticator {
	return &APIKeyAuthenticator{Keys: store, Hasher: hasher}
}

// NewDefaultAPIKeyAuthenticator is a convenience constructor with tuned defaults (t=1, m=64MiB, p=1, keyLen=32).
func NewDefaultAPIKeyAuthenticator(store keys.Reader) *APIKeyAuthenticator {
	hasher := keys.NewArgon2IDHasher(1, 64*1024, 1, 32)
	return &APIKeyAuthenticator{Keys: store, Hasher: hasher}
}

func (a *APIKeyAuthenticator) Authenticate(r *http.Request) (tenant, app string, err error) {
	tok := getHeaderToken(r)
	keyID, secret := splitToken(tok)
	if keyID == "" || secret == "" {
		_ = padWork(a.Hasher)
		return "", "", errors.New("unauthorized")
	}

	rec, err := a.Keys.GetByKeyID(r.Context(), keyID)
	if err != nil || rec == nil {
		_ = padWork(a.Hasher)
		return "", "", errors.New("unauthorized")
	}

	if rec.Status != keys.Active {
		return "", "", errors.New("unauthorized")
	}
	if rec.ExpiresAt != nil && time.Now().After(*rec.ExpiresAt) {
		return "", "", errors.New("unauthorized")
	}

	phc, err := a.Keys.GetPHCByKeyID(r.Context(), keyID)
	if err != nil {
		_ = padWork(a.Hasher)
		return "", "", errors.New("unauthorized")
	}
	ok, _ := a.Hasher.Verify(phc, []byte(secret))
	if !ok {
		return "", "", errors.New("unauthorized")
	}

	_ = a.Keys.TouchLastUsed(r.Context(), keyID)
	return rec.Tenant, rec.App, nil
}
