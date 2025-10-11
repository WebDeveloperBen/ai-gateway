package auth

import (
	"fmt"
	"net/http"
)

type NoopAuthenticator struct{}

func (a *NoopAuthenticator) Authenticate(r *http.Request) (string, *KeyData, error) {
	fmt.Printf("[NoopAuthenticator] called. Returning default key data\n")
	return "default-key-id", &KeyData{
		OrgID:  "default-org",
		AppID:  "default-app",
		UserID: "default-user",
	}, nil
}
