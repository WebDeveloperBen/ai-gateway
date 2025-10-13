package gateway_test

import (
	"net/http"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/auth"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
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

func TestNewCoreWithRegistry(t *testing.T) {
	// Set up a registry with some test deployments
	reg, cleanup := setupRegistry(t)
	defer cleanup()

	// Add some test deployments
	deployments := []model.ModelDeployment{
		{
			Tenant:     "default",
			Model:      "gpt-4",
			Deployment: "test-deployment",
			Provider:   "azure",
			Meta: map[string]string{
				"APIVer":  "2024-07-01-preview",
				"BaseURL": "https://test.openai.azure.com",
			},
		},
		{
			Tenant:     "default",
			Model:      "gpt-3.5-turbo",
			Deployment: "test-openai",
			Provider:   "openai",
			Meta: map[string]string{
				"BaseURL": "https://api.openai.com",
			},
		},
	}

	for _, md := range deployments {
		err := reg.Add(md, 0)
		require.NoError(t, err)
	}

	var authenticator auth.KeyAuthenticator
	customTransport := &mockRoundTripper{}

	core := gateway.NewCoreWithRegistry(customTransport, authenticator, reg)

	require.NotNil(t, core)
	require.Equal(t, customTransport, core.Transport)
	require.Equal(t, authenticator, core.Authenticator)
	require.NotEmpty(t, core.Adapters) // Should have adapters for azure and openai
}

func TestNewCoreWithRegistry_EmptyRegistry(t *testing.T) {
	// Test with empty registry
	reg, cleanup := setupRegistry(t)
	defer cleanup()

	var authenticator auth.KeyAuthenticator

	core := gateway.NewCoreWithRegistry(nil, authenticator, reg)

	require.NotNil(t, core)
	require.NotNil(t, core.Transport) // Should use default transport
	require.Equal(t, authenticator, core.Authenticator)
	require.Empty(t, core.Adapters) // No deployments, so no adapters
}

// Mock round tripper for testing
type mockRoundTripper struct{}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 200}, nil
}
