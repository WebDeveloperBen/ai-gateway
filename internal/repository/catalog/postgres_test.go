package catalog

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	dbdriver "github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	mdl "github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper functions for creating test data
func stringPtr(s string) *string {
	return &s
}

func createTestModel(t *testing.T, fixtures *testkit.DBFixtures, orgID uuid.UUID, provider, modelName, endpointURL string, authType mdl.AuthType, authConfig mdl.AuthConfig, metadata map[string]interface{}) string {
	ctx := context.Background()

	authConfigBytes, err := json.Marshal(authConfig)
	require.NoError(t, err)

	metadataBytes, err := json.Marshal(metadata)
	require.NoError(t, err)

	model, err := fixtures.Queries.CreateModel(ctx, db.CreateModelParams{
		OrgID:          orgID,
		Provider:       provider,
		ModelName:      modelName,
		DeploymentName: nil,
		EndpointUrl:    endpointURL,
		AuthType:       string(authType),
		AuthConfig:     authConfigBytes,
		Metadata:       metadataBytes,
		Enabled:        true,
	})
	require.NoError(t, err)
	return model.ID.String()
}

// TestMain handles shared container setup and cleanup
func TestMain(m *testing.M) {
	// Skip integration tests if running in CI without Docker
	if os.Getenv("CI") == "true" && os.Getenv("DOCKER_AVAILABLE") != "true" {
		os.Exit(0)
	}

	// Run tests
	code := m.Run()

	// Force cleanup of shared containers at the very end
	testkit.CleanupSharedContainers()
	os.Exit(code)
}

func setupTestDB(t *testing.T) (*dbdriver.Postgres, *testkit.DBFixtures) {
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

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	return pg, fixtures
}

// Test suite for PostgresRepo

func TestPostgresRepo_GetByID(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	ctx := context.Background()
	repo := NewPostgresRepo(pg.Queries)

	// Create test data
	orgID, _ := fixtures.CreateTestOrgAndApp(t)
	authConfig := mdl.AuthConfig{
		Type:   mdl.AuthTypeAPIKey,
		APIKey: stringPtr("test-key"),
	}
	metadata := map[string]interface{}{"max_tokens": 4096}
	modelIDStr := createTestModel(t, fixtures, orgID, "openai", "gpt-4", "https://api.openai.com/v1/chat/completions", mdl.AuthTypeAPIKey, authConfig, metadata)
	modelID := uuid.MustParse(modelIDStr)

	// Test GetByID
	retrievedModel, err := repo.GetByID(ctx, modelID)
	require.NoError(t, err)
	assert.Equal(t, modelID, retrievedModel.ID)
	assert.Equal(t, orgID, retrievedModel.OrgID)
	assert.Equal(t, "openai", retrievedModel.Provider)
	assert.Equal(t, "gpt-4", retrievedModel.ModelName)
	assert.Equal(t, "https://api.openai.com/v1/chat/completions", retrievedModel.EndpointURL)
	assert.Equal(t, mdl.AuthTypeAPIKey, retrievedModel.AuthType)
	assert.True(t, retrievedModel.Enabled)
	assert.Equal(t, "test-key", *retrievedModel.AuthConfig.APIKey)
}

func TestPostgresRepo_GetByID_NotFound(t *testing.T) {
	pg, _ := setupTestDB(t)
	ctx := context.Background()
	repo := NewPostgresRepo(pg.Queries)

	// Test with non-existent ID
	_, err := repo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
}

func TestPostgresRepo_GetByProviderAndName(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	ctx := context.Background()
	repo := NewPostgresRepo(pg.Queries)

	// Create test data
	orgID, _ := fixtures.CreateTestOrgAndApp(t)
	authConfig := mdl.AuthConfig{
		Type:   mdl.AuthTypeAPIKey,
		APIKey: stringPtr("test-key"),
	}
	metadata := map[string]interface{}{"max_tokens": 4096}
	modelIDStr := createTestModel(t, fixtures, orgID, "openai", "gpt-4", "https://api.openai.com/v1/chat/completions", mdl.AuthTypeAPIKey, authConfig, metadata)
	modelID := uuid.MustParse(modelIDStr)

	// Test GetByProviderAndName
	model, err := repo.GetByProviderAndName(ctx, orgID, "openai", "gpt-4")
	require.NoError(t, err)
	assert.Equal(t, modelID, model.ID)
	assert.Equal(t, "openai", model.Provider)
	assert.Equal(t, "gpt-4", model.ModelName)
}

func TestPostgresRepo_GetByProviderAndName_NotFound(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	ctx := context.Background()
	repo := NewPostgresRepo(pg.Queries)

	orgID, _ := fixtures.CreateTestOrgAndApp(t)

	// Test with non-existent provider/name
	_, err := repo.GetByProviderAndName(ctx, orgID, "nonexistent", "model")
	assert.Error(t, err)
}

func TestPostgresRepo_ListByOrgID(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	ctx := context.Background()
	repo := NewPostgresRepo(pg.Queries)

	// Create test data
	orgID, _ := fixtures.CreateTestOrgAndApp(t)

	// Create multiple models
	authConfig1 := mdl.AuthConfig{Type: mdl.AuthTypeAPIKey, APIKey: stringPtr("key1")}
	authConfig2 := mdl.AuthConfig{Type: mdl.AuthTypeOAuth2, ClientID: stringPtr("client1")}
	metadata := map[string]interface{}{"max_tokens": 4096}

	createTestModel(t, fixtures, orgID, "openai", "gpt-4", "https://api.openai.com/v1/chat/completions", mdl.AuthTypeAPIKey, authConfig1, metadata)
	createTestModel(t, fixtures, orgID, "openai", "gpt-3.5-turbo", "https://api.openai.com/v1/chat/completions", mdl.AuthTypeOAuth2, authConfig2, metadata)

	// Test ListByOrgID
	models, err := repo.ListByOrgID(ctx, orgID, 100, 0)
	require.NoError(t, err)
	assert.Len(t, models, 2)

	// Check ordering (should be by created_at desc or similar)
	assert.Equal(t, "gpt-3.5-turbo", models[0].ModelName)
	assert.Equal(t, "gpt-4", models[1].ModelName)
}

func TestPostgresRepo_ListByOrgID_Pagination(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	ctx := context.Background()
	repo := NewPostgresRepo(pg.Queries)

	// Create test data
	orgID, _ := fixtures.CreateTestOrgAndApp(t)
	authConfig := mdl.AuthConfig{Type: mdl.AuthTypeAPIKey, APIKey: stringPtr("key")}
	metadata := map[string]interface{}{"max_tokens": 4096}

	// Create 3 models
	createTestModel(t, fixtures, orgID, "openai", "gpt-4", "https://api.openai.com/v1/chat/completions", mdl.AuthTypeAPIKey, authConfig, metadata)
	createTestModel(t, fixtures, orgID, "openai", "gpt-3.5", "https://api.openai.com/v1/chat/completions", mdl.AuthTypeAPIKey, authConfig, metadata)
	createTestModel(t, fixtures, orgID, "anthropic", "claude", "https://api.anthropic.com/v1/messages", mdl.AuthTypeAPIKey, authConfig, metadata)

	// Test pagination - limit 2
	models, err := repo.ListByOrgID(ctx, orgID, 2, 0)
	require.NoError(t, err)
	assert.Len(t, models, 2)

	// Test pagination - offset 1, limit 2
	models, err = repo.ListByOrgID(ctx, orgID, 2, 1)
	require.NoError(t, err)
	assert.Len(t, models, 2)
}

func TestPostgresRepo_ListEnabledByOrgID(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	ctx := context.Background()
	repo := NewPostgresRepo(pg.Queries)

	// Create test data
	orgID, _ := fixtures.CreateTestOrgAndApp(t)
	authConfig := mdl.AuthConfig{Type: mdl.AuthTypeAPIKey, APIKey: stringPtr("key")}
	metadata := map[string]interface{}{"max_tokens": 4096}

	// Create enabled model
	enabledIDStr := createTestModel(t, fixtures, orgID, "openai", "gpt-4", "https://api.openai.com/v1/chat/completions", mdl.AuthTypeAPIKey, authConfig, metadata)
	enabledID := uuid.MustParse(enabledIDStr)

	// Create disabled model
	disabledIDStr := createTestModel(t, fixtures, orgID, "openai", "gpt-3.5", "https://api.openai.com/v1/chat/completions", mdl.AuthTypeAPIKey, authConfig, metadata)
	disabledID := uuid.MustParse(disabledIDStr)

	// Disable one model
	err := repo.Disable(ctx, disabledID)
	require.NoError(t, err)

	// Test ListEnabledByOrgID
	models, err := repo.ListEnabledByOrgID(ctx, orgID, 100, 0)
	require.NoError(t, err)
	assert.Len(t, models, 1)
	assert.Equal(t, enabledID, models[0].ID)
	assert.True(t, models[0].Enabled)
}

func TestPostgresRepo_Create(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	ctx := context.Background()
	repo := NewPostgresRepo(pg.Queries)

	// Create test data
	orgID, _ := fixtures.CreateTestOrgAndApp(t)

	// Test Create
	authConfig := mdl.AuthConfig{
		Type:   mdl.AuthTypeAPIKey,
		APIKey: stringPtr("test-key"),
	}
	metadata := map[string]interface{}{"max_tokens": 4096}

	model, err := repo.Create(ctx, orgID, "openai", "gpt-4", nil, "https://api.openai.com/v1/chat/completions", mdl.AuthTypeAPIKey, authConfig, metadata, true)
	require.NoError(t, err)
	assert.Equal(t, orgID, model.OrgID)
	assert.Equal(t, "openai", model.Provider)
	assert.Equal(t, "gpt-4", model.ModelName)
	assert.Equal(t, "https://api.openai.com/v1/chat/completions", model.EndpointURL)
	assert.Equal(t, mdl.AuthTypeAPIKey, model.AuthType)
	assert.Equal(t, "test-key", *model.AuthConfig.APIKey)
	assert.True(t, model.Enabled)
}

func TestPostgresRepo_Create_InvalidOrgID(t *testing.T) {
	pg, _ := setupTestDB(t)
	ctx := context.Background()
	repo := NewPostgresRepo(pg.Queries)

	authConfig := mdl.AuthConfig{Type: mdl.AuthTypeAPIKey}
	metadata := map[string]interface{}{}

	_, err := repo.Create(ctx, uuid.Nil, "openai", "gpt-4", nil, "https://api.openai.com/v1/chat/completions", mdl.AuthTypeAPIKey, authConfig, metadata, true)
	assert.Error(t, err)
}

func TestPostgresRepo_Update(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	ctx := context.Background()
	repo := NewPostgresRepo(pg.Queries)

	// Create test data
	orgID, _ := fixtures.CreateTestOrgAndApp(t)
	authConfig := mdl.AuthConfig{Type: mdl.AuthTypeAPIKey, APIKey: stringPtr("original")}
	metadata := map[string]interface{}{"max_tokens": 4096}
	modelIDStr := createTestModel(t, fixtures, orgID, "openai", "gpt-4", "https://api.openai.com/v1/chat/completions", mdl.AuthTypeAPIKey, authConfig, metadata)
	modelID := uuid.MustParse(modelIDStr)

	// Test Update
	newAuthConfig := mdl.AuthConfig{
		Type:         mdl.AuthTypeOAuth2,
		ClientID:     stringPtr("new-client"),
		ClientSecret: stringPtr("new-secret"),
	}
	newMetadata := map[string]interface{}{"max_tokens": 8192}

	model, err := repo.Update(ctx, modelID, "openai", "gpt-4-updated", nil, "https://api.openai.com/v1/chat/completions", mdl.AuthTypeOAuth2, newAuthConfig, newMetadata, false)
	require.NoError(t, err)
	assert.Equal(t, modelID, model.ID)
	assert.Equal(t, "gpt-4-updated", model.ModelName)
	assert.Equal(t, mdl.AuthTypeOAuth2, model.AuthType)
	assert.Equal(t, "new-client", *model.AuthConfig.ClientID)
	assert.Equal(t, "new-secret", *model.AuthConfig.ClientSecret)
	assert.False(t, model.Enabled)
}

func TestPostgresRepo_Update_NotFound(t *testing.T) {
	pg, _ := setupTestDB(t)
	ctx := context.Background()
	repo := NewPostgresRepo(pg.Queries)

	authConfig := mdl.AuthConfig{Type: mdl.AuthTypeAPIKey}
	metadata := map[string]interface{}{}

	_, err := repo.Update(ctx, uuid.New(), "openai", "gpt-4", nil, "https://api.openai.com/v1/chat/completions", mdl.AuthTypeAPIKey, authConfig, metadata, true)
	assert.Error(t, err)
}

func TestPostgresRepo_EnableDisable(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	ctx := context.Background()
	repo := NewPostgresRepo(pg.Queries)

	// Create test data
	orgID, _ := fixtures.CreateTestOrgAndApp(t)
	authConfig := mdl.AuthConfig{Type: mdl.AuthTypeAPIKey, APIKey: stringPtr("key")}
	metadata := map[string]interface{}{}
	modelIDStr := createTestModel(t, fixtures, orgID, "openai", "gpt-4", "https://api.openai.com/v1/chat/completions", mdl.AuthTypeAPIKey, authConfig, metadata)
	modelID := uuid.MustParse(modelIDStr)

	// Initially enabled
	model, err := repo.GetByID(ctx, modelID)
	require.NoError(t, err)
	assert.True(t, model.Enabled)

	// Test Disable
	err = repo.Disable(ctx, modelID)
	require.NoError(t, err)

	model, err = repo.GetByID(ctx, modelID)
	require.NoError(t, err)
	assert.False(t, model.Enabled)

	// Test Enable
	err = repo.Enable(ctx, modelID)
	require.NoError(t, err)

	model, err = repo.GetByID(ctx, modelID)
	require.NoError(t, err)
	assert.True(t, model.Enabled)
}

func TestPostgresRepo_EnableDisable_InvalidID(t *testing.T) {
	pg, _ := setupTestDB(t)
	ctx := context.Background()
	repo := NewPostgresRepo(pg.Queries)

	// Test with nil UUID
	err := repo.Disable(ctx, uuid.Nil)
	assert.Error(t, err)

	err = repo.Enable(ctx, uuid.Nil)
	assert.Error(t, err)

	// Test with non-existent ID
	err = repo.Disable(ctx, uuid.New())
	assert.Error(t, err)

	err = repo.Enable(ctx, uuid.New())
	assert.Error(t, err)
}

func TestPostgresRepo_Delete(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	ctx := context.Background()
	repo := NewPostgresRepo(pg.Queries)

	// Create test data
	orgID, _ := fixtures.CreateTestOrgAndApp(t)
	authConfig := mdl.AuthConfig{Type: mdl.AuthTypeAPIKey, APIKey: stringPtr("key")}
	metadata := map[string]interface{}{}
	modelIDStr := createTestModel(t, fixtures, orgID, "openai", "gpt-4", "https://api.openai.com/v1/chat/completions", mdl.AuthTypeAPIKey, authConfig, metadata)
	modelID := uuid.MustParse(modelIDStr)

	// Test Delete
	err := repo.Delete(ctx, modelID)
	require.NoError(t, err)

	// Verify deleted
	_, err = repo.GetByID(ctx, modelID)
	assert.Error(t, err)
}

func TestPostgresRepo_Delete_InvalidID(t *testing.T) {
	pg, _ := setupTestDB(t)
	ctx := context.Background()
	repo := NewPostgresRepo(pg.Queries)

	// Test with nil UUID
	err := repo.Delete(ctx, uuid.Nil)
	assert.Error(t, err)

	// Test with non-existent ID (should not error for delete)
	err = repo.Delete(ctx, uuid.New())
	assert.NoError(t, err) // Delete on non-existent ID should not error
}

// Test utility functions
func TestPostgresRepo_stringToAuthType(t *testing.T) {
	pg, _ := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries).(*postgresRepo) // Cast to concrete type

	assert.Equal(t, mdl.AuthTypeAPIKey, repo.stringToAuthType("api_key"))
	assert.Equal(t, mdl.AuthTypeOAuth2, repo.stringToAuthType("oauth2"))
	assert.Equal(t, mdl.AuthTypeAzureAD, repo.stringToAuthType("azure_ad"))
	assert.Equal(t, mdl.AuthType("unknown"), repo.stringToAuthType("unknown"))
}

func TestPostgresRepo_authTypeToString(t *testing.T) {
	pg, _ := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries).(*postgresRepo) // Cast to concrete type

	assert.Equal(t, "api_key", repo.authTypeToString(mdl.AuthTypeAPIKey))
	assert.Equal(t, "oauth2", repo.authTypeToString(mdl.AuthTypeOAuth2))
	assert.Equal(t, "azure_ad", repo.authTypeToString(mdl.AuthTypeAzureAD))
	assert.Equal(t, "unknown", repo.authTypeToString(mdl.AuthType("unknown")))
}
