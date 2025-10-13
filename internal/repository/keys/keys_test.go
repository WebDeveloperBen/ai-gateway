package keys_test

import (
	"context"
	"testing"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/repository/keys"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMemoryStore(t *testing.T) {
	store := keys.NewMemoryStore()
	assert.NotNil(t, store)
}

func TestMemoryStore_Insert(t *testing.T) {
	store := keys.NewMemoryStore()
	ctx := context.Background()

	key := model.Key{
		ID:        uuid.New(),
		OrgID:     uuid.New(),
		AppID:     uuid.New(),
		UserID:    uuid.New(),
		KeyPrefix: "test-prefix",
		Status:    model.KeyActive,
		LastFour:  "1234",
		Metadata:  []byte(`{"test": "value"}`),
		CreatedAt: time.Now(),
	}

	phc := "$argon2id$..."

	err := store.Insert(ctx, key, phc)
	assert.NoError(t, err)
}

func TestMemoryStore_GetByKeyPrefix(t *testing.T) {
	store := keys.NewMemoryStore()
	ctx := context.Background()

	key := model.Key{
		ID:        uuid.New(),
		OrgID:     uuid.New(),
		AppID:     uuid.New(),
		UserID:    uuid.New(),
		KeyPrefix: "test-prefix",
		Status:    model.KeyActive,
		LastFour:  "1234",
		Metadata:  []byte(`{"test": "value"}`),
		CreatedAt: time.Now(),
	}

	phc := "$argon2id$..."
	err := store.Insert(ctx, key, phc)
	require.NoError(t, err)

	t.Run("existing key", func(t *testing.T) {
		retrieved, err := store.GetByKeyPrefix(ctx, "test-prefix")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, key.ID, retrieved.ID)
		assert.Equal(t, key.KeyPrefix, retrieved.KeyPrefix)
		assert.Equal(t, key.Status, retrieved.Status)
	})

	t.Run("non-existing key", func(t *testing.T) {
		retrieved, err := store.GetByKeyPrefix(ctx, "non-existing")
		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.Contains(t, err.Error(), "key not found")
	})
}

func TestMemoryStore_GetSecretPHCByPrefix(t *testing.T) {
	store := keys.NewMemoryStore()
	ctx := context.Background()

	key := model.Key{
		ID:        uuid.New(),
		OrgID:     uuid.New(),
		AppID:     uuid.New(),
		UserID:    uuid.New(),
		KeyPrefix: "test-prefix",
		Status:    model.KeyActive,
		LastFour:  "1234",
		CreatedAt: time.Now(),
	}

	phc := "$argon2id$test-hash"
	err := store.Insert(ctx, key, phc)
	require.NoError(t, err)

	t.Run("existing key", func(t *testing.T) {
		retrieved, err := store.GetSecretPHCByPrefix(ctx, "test-prefix")
		assert.NoError(t, err)
		assert.Equal(t, phc, retrieved)
	})

	t.Run("non-existing key", func(t *testing.T) {
		retrieved, err := store.GetSecretPHCByPrefix(ctx, "non-existing")
		assert.Error(t, err)
		assert.Empty(t, retrieved)
		assert.Contains(t, err.Error(), "key not found")
	})
}

func TestMemoryStore_UpdateStatus(t *testing.T) {
	store := keys.NewMemoryStore()
	ctx := context.Background()

	key := model.Key{
		ID:        uuid.New(),
		OrgID:     uuid.New(),
		AppID:     uuid.New(),
		UserID:    uuid.New(),
		KeyPrefix: "test-prefix",
		Status:    model.KeyActive,
		LastFour:  "1234",
		CreatedAt: time.Now(),
	}

	phc := "$argon2id$..."
	err := store.Insert(ctx, key, phc)
	require.NoError(t, err)

	t.Run("update to revoked", func(t *testing.T) {
		err := store.UpdateStatus(ctx, "test-prefix", model.KeyRevoked)
		assert.NoError(t, err)

		retrieved, err := store.GetByKeyPrefix(ctx, "test-prefix")
		assert.NoError(t, err)
		assert.Equal(t, model.KeyRevoked, retrieved.Status)
	})

	t.Run("update non-existing key", func(t *testing.T) {
		err := store.UpdateStatus(ctx, "non-existing", model.KeyRevoked)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key not found")
	})
}

func TestMemoryStore_TouchLastUsed(t *testing.T) {
	store := keys.NewMemoryStore()
	ctx := context.Background()

	key := model.Key{
		ID:        uuid.New(),
		OrgID:     uuid.New(),
		AppID:     uuid.New(),
		UserID:    uuid.New(),
		KeyPrefix: "test-prefix",
		Status:    model.KeyActive,
		LastFour:  "1234",
		CreatedAt: time.Now(),
	}

	phc := "$argon2id$..."
	err := store.Insert(ctx, key, phc)
	require.NoError(t, err)

	t.Run("touch existing key", func(t *testing.T) {
		err := store.TouchLastUsed(ctx, "test-prefix")
		assert.NoError(t, err)

		retrieved, err := store.GetByKeyPrefix(ctx, "test-prefix")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved.LastUsedAt)
	})

	t.Run("touch non-existing key", func(t *testing.T) {
		err := store.TouchLastUsed(ctx, "non-existing")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key not found")
	})
}

func TestMemoryStore_Delete(t *testing.T) {
	store := keys.NewMemoryStore()
	ctx := context.Background()

	keyID := uuid.New()
	key := model.Key{
		ID:        keyID,
		OrgID:     uuid.New(),
		AppID:     uuid.New(),
		UserID:    uuid.New(),
		KeyPrefix: "test-prefix",
		Status:    model.KeyActive,
		LastFour:  "1234",
		CreatedAt: time.Now(),
	}

	phc := "$argon2id$..."
	err := store.Insert(ctx, key, phc)
	require.NoError(t, err)

	t.Run("delete existing key", func(t *testing.T) {
		err := store.Delete(ctx, keyID)
		assert.NoError(t, err)

		retrieved, err := store.GetByKeyPrefix(ctx, "test-prefix")
		assert.Error(t, err)
		assert.Nil(t, retrieved)
	})

	t.Run("delete non-existing key", func(t *testing.T) {
		err := store.Delete(ctx, uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key not found")
	})
}

func TestArgon2IDHasher_Hash(t *testing.T) {
	hasher := keys.NewArgon2IDHasher(1, 64*1024, 4, 32)

	secret := []byte("test-secret")
	phc, err := hasher.Hash(secret)
	assert.NoError(t, err)
	assert.NotEmpty(t, phc)
	assert.Contains(t, phc, "$argon2id$")
}

func TestArgon2IDHasher_Verify(t *testing.T) {
	hasher := keys.NewArgon2IDHasher(1, 64*1024, 4, 32)

	secret := []byte("test-secret")
	phc, err := hasher.Hash(secret)
	require.NoError(t, err)

	t.Run("valid secret", func(t *testing.T) {
		valid, err := hasher.Verify(phc, secret)
		assert.NoError(t, err)
		assert.True(t, valid)
	})

	t.Run("invalid secret", func(t *testing.T) {
		valid, err := hasher.Verify(phc, []byte("wrong-secret"))
		assert.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("invalid phc format", func(t *testing.T) {
		valid, err := hasher.Verify("invalid-phc", secret)
		assert.Error(t, err)
		assert.False(t, valid)
	})

	t.Run("unsupported hash type", func(t *testing.T) {
		valid, err := hasher.Verify("$argon2i$...", secret)
		assert.Error(t, err)
		assert.False(t, valid)
		assert.Contains(t, err.Error(), "unsupported hash")
	})
}

func TestNewKeyRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("postgres backend", func(t *testing.T) {
		cfg := model.RepositoryConfig{
			Backend: model.RepositoryPostgres,
			PGPool:  nil, // db.New can handle nil pool
		}

		repo, err := keys.NewKeyRepository(ctx, cfg)
		assert.NoError(t, err)
		assert.NotNil(t, repo)
		// postgres store type is unexported, just check it's not nil
	})

	t.Run("memory backend", func(t *testing.T) {
		cfg := model.RepositoryConfig{
			Backend: model.RepositoryMemory,
		}

		repo, err := keys.NewKeyRepository(ctx, cfg)
		assert.NoError(t, err)
		assert.NotNil(t, repo)
		assert.IsType(t, &keys.MemoryStore{}, repo)
	})

	t.Run("unsupported backend", func(t *testing.T) {
		cfg := model.RepositoryConfig{
			Backend: "unsupported",
		}

		repo, err := keys.NewKeyRepository(ctx, cfg)
		assert.Error(t, err)
		assert.Nil(t, repo)
		assert.Contains(t, err.Error(), "unsupported key repository backend")
	})
}
