package gateway_test

import (
	"net/http"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/auth"
	"github.com/WebDeveloperBen/ai-gateway/internal/provider"
	"github.com/stretchr/testify/require"
)

// Mock adapter for testing
type mockAdapter struct {
	prefix string
}

func (m *mockAdapter) Prefix() string {
	return m.prefix
}

func (m *mockAdapter) Rewrite(req *http.Request, suffix string, info provider.ReqInfo) error {
	return nil
}

func TestNewCoreWithAdapters(t *testing.T) {
	// Test with nil transport (should use default)
	var authenticator auth.KeyAuthenticator
	adapters := []provider.Adapter{
		&mockAdapter{prefix: "/azure"},
		&mockAdapter{prefix: "/openai"},
	}

	core := gateway.NewCoreWithAdapters(nil, authenticator, adapters...)

	require.NotNil(t, core)
	require.Equal(t, 1<<20, core.MaxBody) // 1MB
	require.NotNil(t, core.Transport)     // Should be http.DefaultTransport
	require.Equal(t, adapters, core.Adapters)
	require.Equal(t, authenticator, core.Authenticator)
}

func TestNewCoreWithAdapters_CustomTransport(t *testing.T) {
	customTransport := &mockRoundTripper{}
	var authenticator auth.KeyAuthenticator
	adapters := []provider.Adapter{&mockAdapter{prefix: "/test"}}

	core := gateway.NewCoreWithAdapters(customTransport, authenticator, adapters...)

	require.NotNil(t, core)
	require.Equal(t, customTransport, core.Transport)
	require.Equal(t, adapters, core.Adapters)
	require.Equal(t, authenticator, core.Authenticator)
}

func TestNewCoreWithAdapters_NoAdapters(t *testing.T) {
	var authenticator auth.KeyAuthenticator

	core := gateway.NewCoreWithAdapters(nil, authenticator)

	require.NotNil(t, core)
	require.Empty(t, core.Adapters)
	require.Equal(t, authenticator, core.Authenticator)
}

// Mock round tripper for testing
type mockRoundTripper struct{}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 200}, nil
}
