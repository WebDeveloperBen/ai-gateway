// Package auth implements authentication and tenant identification
// for incoming requests, attaching tenant and app context to the request.
package auth

import (
	"net/http"

	"github.com/insurgence-ai/llm-gateway/internal/repository/keys"
)

type KeyAuthenticator interface {
	Authenticate(r *http.Request) (tenant, app string, err error)
}

type APIKeyAuthenticator struct {
	Keys   keys.Reader
	Hasher keys.Hasher
}
