package docs

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestRegisterRoutes(t *testing.T) {
	// Create a test API
	_, api := humatest.New(t)

	// Create a chi router (which implements http.Handler and has Get method)
	router := chi.NewRouter()

	// Register routes
	RegisterRoutes(router, api)

	// Test that the routes work by making requests
	// Test /docs route
	req := &http.Request{Method: "GET", URL: &url.URL{Path: "/docs"}}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "API Reference")
	assert.Contains(t, w.Body.String(), "@scalar/api-reference")

	// Test /openapi.json route
	req2 := &http.Request{Method: "GET", URL: &url.URL{Path: "/openapi.json"}}
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, "application/json", w2.Header().Get("Content-Type"))
	// Should contain some JSON
	assert.Contains(t, w2.Body.String(), "openapi")
}

func TestRegisterRoutes_InvalidRouter(t *testing.T) {
	// Create a test API
	_, api := humatest.New(t)

	// Create a router that doesn't implement the required interface
	invalidRouter := &invalidRouter{}

	// Should panic
	assert.Panics(t, func() {
		RegisterRoutes(invalidRouter, api)
	})
}

// invalidRouter doesn't implement the required Get method
type invalidRouter struct{}

func (ir *invalidRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// dummy implementation
}
