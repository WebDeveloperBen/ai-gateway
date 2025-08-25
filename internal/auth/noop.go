package auth

import (
	"fmt"
	"net/http"
)

// NoopAuthenticator is for local testing/developlemnt usage only it removes authentication at the proxy
type NoopAuthenticator struct{}

func (a *NoopAuthenticator) Authenticate(r *http.Request) (string, string, error) {
	fmt.Printf("[AnyAPIKeyAuthenticator] called. Returning tenant=default app=default\n")
	return "default", "default", nil // always allows requests, returns dummy tenant/app
}
