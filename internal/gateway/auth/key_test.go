package auth

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockKeyReader struct {
	getByKeyPrefixFunc    func(ctx context.Context, prefix string) (*model.Key, error)
	getSecretPHCByPrefix  func(ctx context.Context, prefix string) (string, error)
	touchLastUsedFunc     func(ctx context.Context, prefix string) error
}

func (m *mockKeyReader) GetByKeyPrefix(ctx context.Context, prefix string) (*model.Key, error) {
	if m.getByKeyPrefixFunc != nil {
		return m.getByKeyPrefixFunc(ctx, prefix)
	}
	return nil, errors.New("not implemented")
}

func (m *mockKeyReader) GetSecretPHCByPrefix(ctx context.Context, prefix string) (string, error) {
	if m.getSecretPHCByPrefix != nil {
		return m.getSecretPHCByPrefix(ctx, prefix)
	}
	return "", errors.New("not implemented")
}

func (m *mockKeyReader) TouchLastUsed(ctx context.Context, prefix string) error {
	if m.touchLastUsedFunc != nil {
		return m.touchLastUsedFunc(ctx, prefix)
	}
	return nil
}

type mockHasher struct {
	hashFunc   func(data []byte) (string, error)
	verifyFunc func(phc string, data []byte) (bool, error)
}

func (m *mockHasher) Hash(data []byte) (string, error) {
	if m.hashFunc != nil {
		return m.hashFunc(data)
	}
	return "$argon2id$v=19$m=65536,t=3,p=2$dummy$hash", nil
}

func (m *mockHasher) Verify(phc string, data []byte) (bool, error) {
	if m.verifyFunc != nil {
		return m.verifyFunc(phc, data)
	}
	return true, nil
}

func TestNewAPIKeyAuthenticator(t *testing.T) {
	t.Run("Creates authenticator with provided dependencies", func(t *testing.T) {
		reader := &mockKeyReader{}
		hasher := &mockHasher{}

		auth := NewAPIKeyAuthenticator(reader, hasher)

		require.NotNil(t, auth)
		assert.Equal(t, reader, auth.Keys)
		assert.Equal(t, hasher, auth.Hasher)
	})
}

func TestNewDefaultAPIKeyAuthenticator(t *testing.T) {
	t.Run("Creates authenticator with default hasher", func(t *testing.T) {
		reader := &mockKeyReader{}

		auth := NewDefaultAPIKeyAuthenticator(reader)

		require.NotNil(t, auth)
		assert.Equal(t, reader, auth.Keys)
		assert.NotNil(t, auth.Hasher)
	})
}

func TestAPIKeyAuthenticator_Authenticate_Success(t *testing.T) {
	orgID := uuid.New()
	appID := uuid.New()
	userID := uuid.New()

	reader := &mockKeyReader{
		getByKeyPrefixFunc: func(ctx context.Context, prefix string) (*model.Key, error) {
			return &model.Key{
				KeyPrefix: "sk_test",
				OrgID:     orgID,
				AppID:     appID,
				UserID:    userID,
				Status:    model.KeyActive,
				ExpiresAt: nil,
			}, nil
		},
		getSecretPHCByPrefix: func(ctx context.Context, prefix string) (string, error) {
			return "$argon2id$v=19$m=65536,t=3,p=2$dummy$hash", nil
		},
		touchLastUsedFunc: func(ctx context.Context, prefix string) error {
			return nil
		},
	}

	hasher := &mockHasher{
		verifyFunc: func(phc string, data []byte) (bool, error) {
			return true, nil
		},
	}

	auth := NewAPIKeyAuthenticator(reader, hasher)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer sk_test.secret123")

	keyID, keyData, err := auth.Authenticate(req)

	require.NoError(t, err)
	assert.Equal(t, "sk_test", keyID)
	require.NotNil(t, keyData)
	assert.Equal(t, "sk_test", keyData.KeyID)
	assert.Equal(t, orgID.String(), keyData.OrgID)
	assert.Equal(t, appID.String(), keyData.AppID)
	assert.Equal(t, userID.String(), keyData.UserID)
}

func TestAPIKeyAuthenticator_Authenticate_NoToken(t *testing.T) {
	reader := &mockKeyReader{}
	hasher := &mockHasher{
		hashFunc: func(data []byte) (string, error) {
			return "padded", nil
		},
	}

	auth := NewAPIKeyAuthenticator(reader, hasher)
	req, _ := http.NewRequest("GET", "/test", nil)

	keyID, keyData, err := auth.Authenticate(req)

	assert.Error(t, err)
	assert.Equal(t, "unauthorized", err.Error())
	assert.Empty(t, keyID)
	assert.Nil(t, keyData)
}

func TestAPIKeyAuthenticator_Authenticate_InvalidTokenFormat(t *testing.T) {
	reader := &mockKeyReader{}
	hasher := &mockHasher{
		hashFunc: func(data []byte) (string, error) {
			return "padded", nil
		},
	}

	auth := NewAPIKeyAuthenticator(reader, hasher)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid_token_without_dot")

	keyID, keyData, err := auth.Authenticate(req)

	assert.Error(t, err)
	assert.Equal(t, "unauthorized", err.Error())
	assert.Empty(t, keyID)
	assert.Nil(t, keyData)
}

func TestAPIKeyAuthenticator_Authenticate_KeyNotFound(t *testing.T) {
	reader := &mockKeyReader{
		getByKeyPrefixFunc: func(ctx context.Context, prefix string) (*model.Key, error) {
			return nil, errors.New("key not found")
		},
	}

	hasher := &mockHasher{
		hashFunc: func(data []byte) (string, error) {
			return "padded", nil
		},
	}

	auth := NewAPIKeyAuthenticator(reader, hasher)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer sk_test.secret123")

	keyID, keyData, err := auth.Authenticate(req)

	assert.Error(t, err)
	assert.Equal(t, "unauthorized", err.Error())
	assert.Empty(t, keyID)
	assert.Nil(t, keyData)
}

func TestAPIKeyAuthenticator_Authenticate_InactiveKey(t *testing.T) {
	reader := &mockKeyReader{
		getByKeyPrefixFunc: func(ctx context.Context, prefix string) (*model.Key, error) {
			return &model.Key{
				KeyPrefix: "sk_test",
				Status:    model.KeyRevoked,
			}, nil
		},
	}

	hasher := &mockHasher{}
	auth := NewAPIKeyAuthenticator(reader, hasher)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer sk_test.secret123")

	keyID, keyData, err := auth.Authenticate(req)

	assert.Error(t, err)
	assert.Equal(t, "unauthorized", err.Error())
	assert.Empty(t, keyID)
	assert.Nil(t, keyData)
}

func TestAPIKeyAuthenticator_Authenticate_ExpiredKey(t *testing.T) {
	expiredTime := time.Now().Add(-24 * time.Hour)

	reader := &mockKeyReader{
		getByKeyPrefixFunc: func(ctx context.Context, prefix string) (*model.Key, error) {
			return &model.Key{
				KeyPrefix: "sk_test",
				Status:    model.KeyActive,
				ExpiresAt: &expiredTime,
			}, nil
		},
	}

	hasher := &mockHasher{}
	auth := NewAPIKeyAuthenticator(reader, hasher)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer sk_test.secret123")

	keyID, keyData, err := auth.Authenticate(req)

	assert.Error(t, err)
	assert.Equal(t, "unauthorized", err.Error())
	assert.Empty(t, keyID)
	assert.Nil(t, keyData)
}

func TestAPIKeyAuthenticator_Authenticate_NotExpiredKey(t *testing.T) {
	futureTime := time.Now().Add(24 * time.Hour)
	orgID := uuid.New()
	appID := uuid.New()
	userID := uuid.New()

	reader := &mockKeyReader{
		getByKeyPrefixFunc: func(ctx context.Context, prefix string) (*model.Key, error) {
			return &model.Key{
				KeyPrefix: "sk_test",
				OrgID:     orgID,
				AppID:     appID,
				UserID:    userID,
				Status:    model.KeyActive,
				ExpiresAt: &futureTime,
			}, nil
		},
		getSecretPHCByPrefix: func(ctx context.Context, prefix string) (string, error) {
			return "$argon2id$hash", nil
		},
	}

	hasher := &mockHasher{
		verifyFunc: func(phc string, data []byte) (bool, error) {
			return true, nil
		},
	}

	auth := NewAPIKeyAuthenticator(reader, hasher)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer sk_test.secret123")

	keyID, keyData, err := auth.Authenticate(req)

	require.NoError(t, err)
	assert.NotEmpty(t, keyID)
	assert.NotNil(t, keyData)
}

func TestAPIKeyAuthenticator_Authenticate_SecretPHCError(t *testing.T) {
	reader := &mockKeyReader{
		getByKeyPrefixFunc: func(ctx context.Context, prefix string) (*model.Key, error) {
			return &model.Key{
				KeyPrefix: "sk_test",
				Status:    model.KeyActive,
			}, nil
		},
		getSecretPHCByPrefix: func(ctx context.Context, prefix string) (string, error) {
			return "", errors.New("database error")
		},
	}

	hasher := &mockHasher{
		hashFunc: func(data []byte) (string, error) {
			return "padded", nil
		},
	}

	auth := NewAPIKeyAuthenticator(reader, hasher)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer sk_test.secret123")

	keyID, keyData, err := auth.Authenticate(req)

	assert.Error(t, err)
	assert.Equal(t, "unauthorized", err.Error())
	assert.Empty(t, keyID)
	assert.Nil(t, keyData)
}

func TestAPIKeyAuthenticator_Authenticate_InvalidSecret(t *testing.T) {
	reader := &mockKeyReader{
		getByKeyPrefixFunc: func(ctx context.Context, prefix string) (*model.Key, error) {
			return &model.Key{
				KeyPrefix: "sk_test",
				Status:    model.KeyActive,
			}, nil
		},
		getSecretPHCByPrefix: func(ctx context.Context, prefix string) (string, error) {
			return "$argon2id$hash", nil
		},
	}

	hasher := &mockHasher{
		verifyFunc: func(phc string, data []byte) (bool, error) {
			return false, nil
		},
	}

	auth := NewAPIKeyAuthenticator(reader, hasher)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer sk_test.wrongsecret")

	keyID, keyData, err := auth.Authenticate(req)

	assert.Error(t, err)
	assert.Equal(t, "unauthorized", err.Error())
	assert.Empty(t, keyID)
	assert.Nil(t, keyData)
}

func TestAPIKeyAuthenticator_Authenticate_TouchLastUsedCalled(t *testing.T) {
	orgID := uuid.New()
	appID := uuid.New()
	userID := uuid.New()
	touchCalled := false

	reader := &mockKeyReader{
		getByKeyPrefixFunc: func(ctx context.Context, prefix string) (*model.Key, error) {
			return &model.Key{
				KeyPrefix: "sk_test",
				OrgID:     orgID,
				AppID:     appID,
				UserID:    userID,
				Status:    model.KeyActive,
			}, nil
		},
		getSecretPHCByPrefix: func(ctx context.Context, prefix string) (string, error) {
			return "$argon2id$hash", nil
		},
		touchLastUsedFunc: func(ctx context.Context, prefix string) error {
			touchCalled = true
			assert.Equal(t, "sk_test", prefix)
			return nil
		},
	}

	hasher := &mockHasher{
		verifyFunc: func(phc string, data []byte) (bool, error) {
			return true, nil
		},
	}

	auth := NewAPIKeyAuthenticator(reader, hasher)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer sk_test.secret123")

	_, _, err := auth.Authenticate(req)

	require.NoError(t, err)
	assert.True(t, touchCalled, "TouchLastUsed should be called on successful authentication")
}

func TestAPIKeyAuthenticator_Authenticate_XAPIKeyHeader(t *testing.T) {
	orgID := uuid.New()
	appID := uuid.New()
	userID := uuid.New()

	reader := &mockKeyReader{
		getByKeyPrefixFunc: func(ctx context.Context, prefix string) (*model.Key, error) {
			return &model.Key{
				KeyPrefix: "sk_test",
				OrgID:     orgID,
				AppID:     appID,
				UserID:    userID,
				Status:    model.KeyActive,
			}, nil
		},
		getSecretPHCByPrefix: func(ctx context.Context, prefix string) (string, error) {
			return "$argon2id$hash", nil
		},
	}

	hasher := &mockHasher{
		verifyFunc: func(phc string, data []byte) (bool, error) {
			return true, nil
		},
	}

	auth := NewAPIKeyAuthenticator(reader, hasher)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "sk_test.secret123")

	keyID, keyData, err := auth.Authenticate(req)

	require.NoError(t, err)
	assert.Equal(t, "sk_test", keyID)
	assert.NotNil(t, keyData)
}

func TestAPIKeyAuthenticator_ImplementsInterface(t *testing.T) {
	t.Run("Implements KeyAuthenticator interface", func(t *testing.T) {
		var _ KeyAuthenticator = &APIKeyAuthenticator{}
	})
}
