// Package auth implements authentication and tenant identification
// for incoming requests, attaching tenant and app context to the request.
package auth

import (
	"errors"
	"net/http"
	"strings"
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
	tok := headerToken(r)
	keyID, secret := split(tok)
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

func headerToken(r *http.Request) string {
	if v := r.Header.Get("Authorization"); v != "" {
		if after, ok := strings.CutPrefix(v, "Bearer "); ok {
			return strings.TrimSpace(after)
		}
	}
	if v := r.Header.Get("X-API-Key"); v != "" {
		return strings.TrimSpace(v)
	}
	return ""
}

func split(tok string) (string, string) {
	left, right, ok := strings.Cut(tok, ".")
	if !ok || left == "" || right == "" {
		return "", ""
	}
	return left, right
}

// padWork burns comparable CPU on early rejects to avoid an existence oracle.
func padWork(h keys.Hasher) error {
	_, _ = h.Hash([]byte("timing-pad"))
	return nil
}
