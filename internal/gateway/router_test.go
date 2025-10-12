package gateway_test

import (
	"net/http"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway"
	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/require"
)

func TestRequestInfo_Struct(t *testing.T) {
	info := gateway.RequestInfo{
		Method: "POST",
		Path:   "/api/chat/completions",
		Query:  "model=gpt-4",
		Model:  "gpt-4",
		Tenant: "tenant-123",
		App:    "app-456",
	}

	require.Equal(t, "POST", info.Method)
	require.Equal(t, "/api/chat/completions", info.Path)
	require.Equal(t, "model=gpt-4", info.Query)
	require.Equal(t, "gpt-4", info.Model)
	require.Equal(t, "tenant-123", info.Tenant)
	require.Equal(t, "app-456", info.App)
}

func TestRequestInfo_ZeroValue(t *testing.T) {
	var info gateway.RequestInfo

	require.Empty(t, info.Method)
	require.Empty(t, info.Path)
	require.Empty(t, info.Query)
	require.Empty(t, info.Model)
	require.Empty(t, info.Tenant)
	require.Empty(t, info.App)
}

func TestMutateFunc_Type(t *testing.T) {
	// Test that MutateFunc is a valid function type
	var mutate gateway.MutateFunc

	// Should be nil initially
	require.Nil(t, mutate)

	// Can assign a function
	mutate = func(hctx huma.Context, req *http.Request, rawBody []byte) error {
		return nil
	}

	require.NotNil(t, mutate)
}

func TestRouter_Interface(t *testing.T) {
	// Test that Router is an interface type
	var router gateway.Router

	// Should be nil initially
	require.Nil(t, router)

	// Verify Router is an interface type by checking its zero value
	// We can't instantiate it directly, but we can verify the type exists
	require.Equal(t, nil, router)
}
