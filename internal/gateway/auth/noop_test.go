package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoopAuthenticator(t *testing.T) {
	t.Run("Returns default values", func(t *testing.T) {
		auth := &NoopAuthenticator{}
		req, _ := http.NewRequest("GET", "/test", nil)

		keyID, keyData, err := auth.Authenticate(req)

		require.NoError(t, err)
		assert.Equal(t, "default-key-id", keyID)
		require.NotNil(t, keyData)
		assert.Equal(t, "default-org", keyData.OrgID)
		assert.Equal(t, "default-app", keyData.AppID)
		assert.Equal(t, "default-user", keyData.UserID)
	})

	t.Run("Always succeeds regardless of request", func(t *testing.T) {
		auth := &NoopAuthenticator{}
		req, _ := http.NewRequest("POST", "/v1/chat/completions", nil)

		keyID, keyData, err := auth.Authenticate(req)

		require.NoError(t, err)
		assert.NotEmpty(t, keyID)
		assert.NotNil(t, keyData)
	})

	t.Run("Works without any headers", func(t *testing.T) {
		auth := &NoopAuthenticator{}
		req, _ := http.NewRequest("GET", "/test", nil)

		keyID, keyData, err := auth.Authenticate(req)

		require.NoError(t, err)
		assert.Equal(t, "default-key-id", keyID)
		assert.NotNil(t, keyData)
	})

	t.Run("Ignores Authorization header", func(t *testing.T) {
		auth := &NoopAuthenticator{}
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer fake-token")

		keyID, keyData, err := auth.Authenticate(req)

		require.NoError(t, err)
		assert.Equal(t, "default-key-id", keyID)
		assert.NotNil(t, keyData)
	})

	t.Run("Implements KeyAuthenticator interface", func(t *testing.T) {
		var _ KeyAuthenticator = &NoopAuthenticator{}
	})

	t.Run("KeyData has expected structure", func(t *testing.T) {
		auth := &NoopAuthenticator{}
		req, _ := http.NewRequest("GET", "/test", nil)

		_, keyData, err := auth.Authenticate(req)

		require.NoError(t, err)
		require.NotNil(t, keyData)
		assert.NotEmpty(t, keyData.OrgID)
		assert.NotEmpty(t, keyData.AppID)
		assert.NotEmpty(t, keyData.UserID)
		assert.Empty(t, keyData.KeyID)
	})

	t.Run("Multiple calls return same values", func(t *testing.T) {
		auth := &NoopAuthenticator{}
		req, _ := http.NewRequest("GET", "/test", nil)

		keyID1, keyData1, err1 := auth.Authenticate(req)
		keyID2, keyData2, err2 := auth.Authenticate(req)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.Equal(t, keyID1, keyID2)
		assert.Equal(t, keyData1.OrgID, keyData2.OrgID)
		assert.Equal(t, keyData1.AppID, keyData2.AppID)
		assert.Equal(t, keyData1.UserID, keyData2.UserID)
	})
}
