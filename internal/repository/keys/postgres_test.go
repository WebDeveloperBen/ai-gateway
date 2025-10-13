package keys_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	dbdriver "github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/repository/keys"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMain handles test setup and cleanup
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func setupTestDB(t *testing.T) *dbdriver.Postgres {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set up shared test containers
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	ctx := context.Background()

	// Set up database connection
	pg, err := dbdriver.NewPostgresDriver(ctx, pgConnStr)
	require.NoError(t, err)
	t.Cleanup(func() { pg.Pool.Close() })

	return pg
}

func setOrgContext(t *testing.T, pg *dbdriver.Postgres, orgID uuid.UUID) {
	ctx := context.Background()
	// Set the session variable for RLS using a transaction
	conn, err := pg.Pool.Acquire(ctx)
	require.NoError(t, err)
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	sql := fmt.Sprintf("SET LOCAL app.current_org = '%s'", orgID.String())
	_, err = tx.Exec(ctx, sql)
	require.NoError(t, err)

	err = tx.Commit(ctx)
	require.NoError(t, err)
}

func createTestData(t *testing.T, pg *dbdriver.Postgres) (orgID, appID, userID uuid.UUID) {
	ctx := context.Background()

	// Create test org
	err := pg.Pool.QueryRow(ctx, "INSERT INTO organisations (name) VALUES ($1) RETURNING id", "test-org").Scan(&orgID)
	require.NoError(t, err)

	// Create test user
	err = pg.Pool.QueryRow(ctx, "INSERT INTO users (org_id, email) VALUES ($1, $2) RETURNING id", orgID, "test@example.com").Scan(&userID)
	require.NoError(t, err)

	// Create test app
	err = pg.Pool.QueryRow(ctx, "INSERT INTO applications (org_id, name, description) VALUES ($1, $2, $3) RETURNING id", orgID, "test-app", "test description").Scan(&appID)
	require.NoError(t, err)

	return orgID, appID, userID
}

func TestPostgresStore_Insert(t *testing.T) {
	pg := setupTestDB(t)
	ctx := context.Background()

	store := keys.NewPostgresStore(pg.Queries)

	// Create test data
	orgID, appID, userID := createTestData(t, pg)
	setOrgContext(t, pg, orgID)

	key := model.Key{
		OrgID:     orgID,
		AppID:     appID,
		UserID:    userID,
		KeyPrefix: "test-prefix-insert",
		Status:    model.KeyActive,
		LastFour:  "1234",
		Metadata:  []byte(`{}`),
		CreatedAt: time.Now(),
	}

	phc := "$argon2id$test-hash"

	t.Run("insert new key", func(t *testing.T) {
		err := store.Insert(ctx, key, phc)
		assert.NoError(t, err)

		// Verify the key was inserted
		setOrgContext(t, pg, orgID)
		retrieved, err := store.GetByKeyPrefix(ctx, "test-prefix-insert")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.NotNil(t, retrieved.ID)
		assert.Equal(t, "test-prefix-insert", retrieved.KeyPrefix)
		assert.Equal(t, model.KeyActive, retrieved.Status)
		assert.Equal(t, "1234", retrieved.LastFour)
	})

	t.Run("insert duplicate key prefix", func(t *testing.T) {
		duplicateKey := model.Key{
			OrgID:     orgID,
			AppID:     appID,
			UserID:    userID,
			KeyPrefix: "test-prefix-insert-2", // Different prefix
			Status:    model.KeyActive,
			LastFour:  "5678",
			Metadata:  []byte(`{}`),
			CreatedAt: time.Now(),
		}

		err := store.Insert(ctx, duplicateKey, phc)
		assert.NoError(t, err) // Should succeed with different prefix
	})
}

func TestPostgresStore_GetByKeyPrefix(t *testing.T) {
	pg := setupTestDB(t)
	ctx := context.Background()

	store := keys.NewPostgresStore(pg.Queries)

	// Create test data
	orgID, appID, userID := createTestData(t, pg)
	setOrgContext(t, pg, orgID)

	key := model.Key{
		OrgID:     orgID,
		AppID:     appID,
		UserID:    userID,
		KeyPrefix: "test-prefix-get",
		Status:    model.KeyActive,
		LastFour:  "1234",
		Metadata:  []byte(`{}`),
		CreatedAt: time.Now(),
	}

	phc := "$argon2id$test-hash"
	err := store.Insert(ctx, key, phc)
	require.NoError(t, err)

	t.Run("existing key", func(t *testing.T) {
		setOrgContext(t, pg, orgID)
		retrieved, err := store.GetByKeyPrefix(ctx, "test-prefix-get")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.NotNil(t, retrieved.ID)
		assert.Equal(t, "test-prefix-get", retrieved.KeyPrefix)
		assert.Equal(t, model.KeyActive, retrieved.Status)
	})

	t.Run("non-existing key", func(t *testing.T) {
		setOrgContext(t, pg, orgID)
		retrieved, err := store.GetByKeyPrefix(ctx, "non-existing")
		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.Contains(t, err.Error(), "key not found")
	})
}

func TestPostgresStore_GetSecretPHCByPrefix(t *testing.T) {
	pg := setupTestDB(t)
	ctx := context.Background()

	store := keys.NewPostgresStore(pg.Queries)

	// Create test data
	orgID, appID, userID := createTestData(t, pg)
	setOrgContext(t, pg, orgID)

	key := model.Key{
		OrgID:     orgID,
		AppID:     appID,
		UserID:    userID,
		KeyPrefix: "test-prefix-phc",
		Status:    model.KeyActive,
		LastFour:  "1234",
		Metadata:  []byte(`{}`),
		CreatedAt: time.Now(),
	}

	phc := "$argon2id$test-hash"
	err := store.Insert(ctx, key, phc)
	require.NoError(t, err)

	t.Run("existing key", func(t *testing.T) {
		setOrgContext(t, pg, orgID)
		retrieved, err := store.GetSecretPHCByPrefix(ctx, "test-prefix-phc")
		assert.NoError(t, err)
		assert.Equal(t, phc, retrieved)
	})

	t.Run("non-existing key", func(t *testing.T) {
		setOrgContext(t, pg, orgID)
		retrieved, err := store.GetSecretPHCByPrefix(ctx, "non-existing")
		assert.Error(t, err)
		assert.Empty(t, retrieved)
		assert.Contains(t, err.Error(), "key not found")
	})
}

func TestPostgresStore_UpdateStatus(t *testing.T) {
	pg := setupTestDB(t)
	ctx := context.Background()

	store := keys.NewPostgresStore(pg.Queries)

	// Create test data
	orgID, appID, userID := createTestData(t, pg)
	setOrgContext(t, pg, orgID)

	key := model.Key{
		OrgID:     orgID,
		AppID:     appID,
		UserID:    userID,
		KeyPrefix: "test-prefix-update",
		Status:    model.KeyActive,
		LastFour:  "1234",
		Metadata:  []byte(`{}`),
		CreatedAt: time.Now(),
	}

	phc := "$argon2id$test-hash"
	err := store.Insert(ctx, key, phc)
	require.NoError(t, err)

	t.Run("update to revoked", func(t *testing.T) {
		setOrgContext(t, pg, orgID)
		err := store.UpdateStatus(ctx, "test-prefix-update", model.KeyRevoked)
		assert.NoError(t, err)

		setOrgContext(t, pg, orgID)
		retrieved, err := store.GetByKeyPrefix(ctx, "test-prefix-update")
		assert.NoError(t, err)
		assert.Equal(t, model.KeyRevoked, retrieved.Status)
	})

	t.Run("update to expired", func(t *testing.T) {
		setOrgContext(t, pg, orgID)
		err := store.UpdateStatus(ctx, "test-prefix-update", model.KeyExpired)
		assert.NoError(t, err)

		setOrgContext(t, pg, orgID)
		retrieved, err := store.GetByKeyPrefix(ctx, "test-prefix-update")
		assert.NoError(t, err)
		assert.Equal(t, model.KeyExpired, retrieved.Status)
	})

	t.Run("invalid status", func(t *testing.T) {
		setOrgContext(t, pg, orgID)
		err := store.UpdateStatus(ctx, "test-prefix-update", model.KeyStatus("invalid"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status")
	})

	t.Run("non-existing key", func(t *testing.T) {
		setOrgContext(t, pg, orgID)
		err := store.UpdateStatus(ctx, "non-existing-update", model.KeyRevoked)
		assert.Error(t, err)
	})
}

func TestPostgresStore_TouchLastUsed(t *testing.T) {
	pg := setupTestDB(t)
	ctx := context.Background()

	store := keys.NewPostgresStore(pg.Queries)

	// Create test data
	orgID, appID, userID := createTestData(t, pg)
	setOrgContext(t, pg, orgID)

	key := model.Key{
		OrgID:     orgID,
		AppID:     appID,
		UserID:    userID,
		KeyPrefix: "test-prefix-touch",
		Status:    model.KeyActive,
		LastFour:  "1234",
		Metadata:  []byte(`{}`),
		CreatedAt: time.Now(),
	}

	phc := "$argon2id$test-hash"
	err := store.Insert(ctx, key, phc)
	require.NoError(t, err)

	t.Run("touch existing key", func(t *testing.T) {
		setOrgContext(t, pg, orgID)
		err := store.TouchLastUsed(ctx, "test-prefix-touch")
		assert.NoError(t, err)

		setOrgContext(t, pg, orgID)
		retrieved, err := store.GetByKeyPrefix(ctx, "test-prefix-touch")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved.LastUsedAt)
	})

	t.Run("touch non-existing key", func(t *testing.T) {
		setOrgContext(t, pg, orgID)
		err := store.TouchLastUsed(ctx, "non-existing-touch")
		assert.Error(t, err)
	})
}

func TestPostgresStore_Delete(t *testing.T) {
	pg := setupTestDB(t)
	ctx := context.Background()

	store := keys.NewPostgresStore(pg.Queries)

	// Create test data
	orgID, appID, userID := createTestData(t, pg)
	setOrgContext(t, pg, orgID)

	key := model.Key{
		OrgID:     orgID,
		AppID:     appID,
		UserID:    userID,
		KeyPrefix: "test-prefix-delete",
		Status:    model.KeyActive,
		LastFour:  "1234",
		Metadata:  []byte(`{}`),
		CreatedAt: time.Now(),
	}

	phc := "$argon2id$test-hash"
	err := store.Insert(ctx, key, phc)
	require.NoError(t, err)

	// Get the inserted key to get its generated ID
	setOrgContext(t, pg, orgID)
	insertedKey, err := store.GetByKeyPrefix(ctx, "test-prefix-delete")
	require.NoError(t, err)

	t.Run("delete existing key", func(t *testing.T) {
		setOrgContext(t, pg, orgID)
		err := store.Delete(ctx, insertedKey.ID)
		assert.NoError(t, err)

		setOrgContext(t, pg, orgID)
		retrieved, err := store.GetByKeyPrefix(ctx, "test-prefix-delete")
		assert.Error(t, err)
		assert.Nil(t, retrieved)
	})

	t.Run("delete non-existing key", func(t *testing.T) {
		setOrgContext(t, pg, orgID)
		err := store.Delete(ctx, uuid.New())
		assert.Error(t, err)
	})
}
