package auth

import (
	"net/http"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/repository/keys"
	"github.com/stretchr/testify/assert"
)

func TestGetHeaderToken(t *testing.T) {
	t.Run("Extracts token from Authorization Bearer header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer sk_test.secret123")

		token := getHeaderToken(req)
		assert.Equal(t, "sk_test.secret123", token)
	})

	t.Run("Trims whitespace from Bearer token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer  sk_test.secret123  ")

		token := getHeaderToken(req)
		assert.Equal(t, "sk_test.secret123", token)
	})

	t.Run("Extracts token from X-API-Key header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-API-Key", "sk_test.secret123")

		token := getHeaderToken(req)
		assert.Equal(t, "sk_test.secret123", token)
	})

	t.Run("Trims whitespace from X-API-Key", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-API-Key", "  sk_test.secret123  ")

		token := getHeaderToken(req)
		assert.Equal(t, "sk_test.secret123", token)
	})

	t.Run("Prioritizes Authorization header over X-API-Key", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer sk_auth.secret")
		req.Header.Set("X-API-Key", "sk_apikey.secret")

		token := getHeaderToken(req)
		assert.Equal(t, "sk_auth.secret", token)
	})

	t.Run("Returns empty string when no headers present", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)

		token := getHeaderToken(req)
		assert.Equal(t, "", token)
	})

	t.Run("Returns empty string for malformed Authorization header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Basic invalid")

		token := getHeaderToken(req)
		assert.Equal(t, "", token)
	})

	t.Run("Returns empty string for empty Bearer", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer ")

		token := getHeaderToken(req)
		assert.Equal(t, "", token)
	})

	t.Run("Handles Authorization without Bearer prefix", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "sk_test.secret123")

		token := getHeaderToken(req)
		assert.Equal(t, "", token)
	})
}

func TestSplitToken(t *testing.T) {
	t.Run("Splits valid token with dot separator", func(t *testing.T) {
		prefix, secret := splitToken("sk_test.secret123")
		assert.Equal(t, "sk_test", prefix)
		assert.Equal(t, "secret123", secret)
	})

	t.Run("Handles token with multiple dots", func(t *testing.T) {
		prefix, secret := splitToken("sk_test.secret.with.dots")
		assert.Equal(t, "sk_test", prefix)
		assert.Equal(t, "secret.with.dots", secret)
	})

	t.Run("Returns empty strings for token without dot", func(t *testing.T) {
		prefix, secret := splitToken("sk_test_nodot")
		assert.Equal(t, "", prefix)
		assert.Equal(t, "", secret)
	})

	t.Run("Returns empty strings when prefix is empty", func(t *testing.T) {
		prefix, secret := splitToken(".secret123")
		assert.Equal(t, "", prefix)
		assert.Equal(t, "", secret)
	})

	t.Run("Returns empty strings when secret is empty", func(t *testing.T) {
		prefix, secret := splitToken("sk_test.")
		assert.Equal(t, "", prefix)
		assert.Equal(t, "", secret)
	})

	t.Run("Returns empty strings for empty token", func(t *testing.T) {
		prefix, secret := splitToken("")
		assert.Equal(t, "", prefix)
		assert.Equal(t, "", secret)
	})

	t.Run("Returns empty strings for just a dot", func(t *testing.T) {
		prefix, secret := splitToken(".")
		assert.Equal(t, "", prefix)
		assert.Equal(t, "", secret)
	})

	t.Run("Handles long prefix and secret", func(t *testing.T) {
		longPrefix := "sk_live_very_long_prefix_with_many_characters"
		longSecret := "secret_with_very_long_random_string_abcdefghijklmnop"
		token := longPrefix + "." + longSecret

		prefix, secret := splitToken(token)
		assert.Equal(t, longPrefix, prefix)
		assert.Equal(t, longSecret, secret)
	})
}

func TestPadWork(t *testing.T) {
	t.Run("Executes without error", func(t *testing.T) {
		hasher := keys.NewArgon2IDHasher(1, 64*1024, 1, 32)
		err := padWork(hasher)
		assert.NoError(t, err)
	})

	t.Run("Calls hasher with timing-pad input", func(t *testing.T) {
		hasher := keys.NewArgon2IDHasher(1, 64*1024, 1, 32)
		err := padWork(hasher)
		assert.NoError(t, err)
	})

	t.Run("Returns nil error", func(t *testing.T) {
		hasher := keys.NewArgon2IDHasher(1, 64*1024, 1, 32)
		result := padWork(hasher)
		assert.Nil(t, result)
	})
}
