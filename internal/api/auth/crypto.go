package auth

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateState returns a secure random string for OAuth state handling.
func generateState(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
