package usage

import (
	"context"
	"os"
	"testing"
	"time"

	dbdriver "github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestPostgresRepo_Create(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test data
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	keyID := fixtures.CreateTestAPIKey(t, orgID)

	// Test Create
	metric := &model.UsageMetric{
		OrgID:             orgID,
		AppID:             appID,
		APIKeyID:          keyID,
		Provider:          "openai",
		ModelName:         "gpt-3.5-turbo",
		PromptTokens:      100,
		CompletionTokens:  50,
		TotalTokens:       150,
		RequestSizeBytes:  1024,
		ResponseSizeBytes: 512,
		Timestamp:         time.Now(),
	}

	err := repo.Create(ctx, metric)
	require.NoError(t, err)
}

func TestPostgresRepo_GetByAppID(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test data
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	keyID := fixtures.CreateTestAPIKey(t, orgID)

	// Create usage metrics
	now := time.Now()
	metric1 := &model.UsageMetric{
		OrgID:             orgID,
		AppID:             appID,
		APIKeyID:          keyID,
		Provider:          "openai",
		ModelName:         "gpt-3.5-turbo",
		PromptTokens:      100,
		CompletionTokens:  50,
		TotalTokens:       150,
		RequestSizeBytes:  1024,
		ResponseSizeBytes: 512,
		Timestamp:         now.Add(-time.Hour),
	}
	metric2 := &model.UsageMetric{
		OrgID:             orgID,
		AppID:             appID,
		APIKeyID:          keyID,
		Provider:          "openai",
		ModelName:         "gpt-4",
		PromptTokens:      200,
		CompletionTokens:  100,
		TotalTokens:       300,
		RequestSizeBytes:  2048,
		ResponseSizeBytes: 1024,
		Timestamp:         now,
	}

	err := repo.Create(ctx, metric1)
	require.NoError(t, err)
	err = repo.Create(ctx, metric2)
	require.NoError(t, err)

	// Test GetByAppID
	start := now.Add(-2 * time.Hour)
	end := now.Add(time.Hour)
	metrics, err := repo.GetByAppID(ctx, appID, start, end)
	require.NoError(t, err)
	assert.Len(t, metrics, 2)

	// Check that metrics are returned in descending timestamp order
	assert.Equal(t, "gpt-4", metrics[0].ModelName)
	assert.Equal(t, 200, metrics[0].PromptTokens)
	assert.Equal(t, "gpt-3.5-turbo", metrics[1].ModelName)
	assert.Equal(t, 100, metrics[1].PromptTokens)
}

func TestPostgresRepo_SumTokensByAppID(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test data
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	keyID := fixtures.CreateTestAPIKey(t, orgID)

	// Create usage metrics
	now := time.Now()
	metric1 := &model.UsageMetric{
		OrgID:             orgID,
		AppID:             appID,
		APIKeyID:          keyID,
		Provider:          "openai",
		ModelName:         "gpt-3.5-turbo",
		PromptTokens:      100,
		CompletionTokens:  50,
		TotalTokens:       150,
		RequestSizeBytes:  1024,
		ResponseSizeBytes: 512,
		Timestamp:         now,
	}
	metric2 := &model.UsageMetric{
		OrgID:             orgID,
		AppID:             appID,
		APIKeyID:          keyID,
		Provider:          "openai",
		ModelName:         "gpt-4",
		PromptTokens:      200,
		CompletionTokens:  100,
		TotalTokens:       300,
		RequestSizeBytes:  2048,
		ResponseSizeBytes: 1024,
		Timestamp:         now,
	}

	err := repo.Create(ctx, metric1)
	require.NoError(t, err)
	err = repo.Create(ctx, metric2)
	require.NoError(t, err)

	// Test SumTokensByAppID
	start := now.Add(-time.Hour)
	end := now.Add(time.Hour)
	summary, err := repo.SumTokensByAppID(ctx, appID, start, end)
	require.NoError(t, err)

	assert.Equal(t, 300, summary.TotalPromptTokens)
	assert.Equal(t, 150, summary.TotalCompletionTokens)
	assert.Equal(t, 450, summary.TotalTokens)
	assert.Equal(t, 2, summary.RequestCount)
}

func TestPostgresRepo_GetUsageByModel(t *testing.T) {
	pg, fixtures := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Create test data
	orgID, appID := fixtures.CreateTestOrgAndApp(t)
	keyID := fixtures.CreateTestAPIKey(t, orgID)

	// Create usage metrics
	now := time.Now()
	metric1 := &model.UsageMetric{
		OrgID:             orgID,
		AppID:             appID,
		APIKeyID:          keyID,
		Provider:          "openai",
		ModelName:         "gpt-3.5-turbo",
		PromptTokens:      100,
		CompletionTokens:  50,
		TotalTokens:       150,
		RequestSizeBytes:  1024,
		ResponseSizeBytes: 512,
		Timestamp:         now,
	}
	metric2 := &model.UsageMetric{
		OrgID:             orgID,
		AppID:             appID,
		APIKeyID:          keyID,
		Provider:          "openai",
		ModelName:         "gpt-3.5-turbo",
		PromptTokens:      200,
		CompletionTokens:  100,
		TotalTokens:       300,
		RequestSizeBytes:  2048,
		ResponseSizeBytes: 1024,
		Timestamp:         now,
	}
	metric3 := &model.UsageMetric{
		OrgID:             orgID,
		AppID:             appID,
		APIKeyID:          keyID,
		Provider:          "azure",
		ModelName:         "gpt-4",
		PromptTokens:      150,
		CompletionTokens:  75,
		TotalTokens:       225,
		RequestSizeBytes:  1536,
		ResponseSizeBytes: 768,
		Timestamp:         now,
	}

	err := repo.Create(ctx, metric1)
	require.NoError(t, err)
	err = repo.Create(ctx, metric2)
	require.NoError(t, err)
	err = repo.Create(ctx, metric3)
	require.NoError(t, err)

	// Test GetUsageByModel
	start := now.Add(-time.Hour)
	end := now.Add(time.Hour)
	summaries, err := repo.GetUsageByModel(ctx, appID, start, end)
	require.NoError(t, err)
	assert.Len(t, summaries, 2)

	// Should be ordered by total_tokens DESC
	assert.Equal(t, "gpt-3.5-turbo", summaries[0].ModelName)
	assert.Equal(t, "openai", summaries[0].Provider)
	assert.Equal(t, 300, summaries[0].TotalPromptTokens)
	assert.Equal(t, 150, summaries[0].TotalCompletionTokens)
	assert.Equal(t, 450, summaries[0].TotalTokens)
	assert.Equal(t, 2, summaries[0].RequestCount)

	assert.Equal(t, "gpt-4", summaries[1].ModelName)
	assert.Equal(t, "azure", summaries[1].Provider)
	assert.Equal(t, 150, summaries[1].TotalPromptTokens)
	assert.Equal(t, 75, summaries[1].TotalCompletionTokens)
	assert.Equal(t, 225, summaries[1].TotalTokens)
	assert.Equal(t, 1, summaries[1].RequestCount)
}

func TestPostgresRepo_Create_Validation(t *testing.T) {
	pg, _ := setupTestDB(t)
	repo := NewPostgresRepo(pg.Queries)
	ctx := context.Background()

	// Test Create with nil orgID
	metric := &model.UsageMetric{
		AppID:    uuid.New(),
		APIKeyID: uuid.New(),
	}
	err := repo.Create(ctx, metric)
	assert.Error(t, err)

	// Test Create with nil appID
	metric = &model.UsageMetric{
		OrgID:    uuid.New(),
		APIKeyID: uuid.New(),
	}
	err = repo.Create(ctx, metric)
	assert.Error(t, err)

	// Test Create with nil apiKeyID
	metric = &model.UsageMetric{
		OrgID: uuid.New(),
		AppID: uuid.New(),
	}
	err = repo.Create(ctx, metric)
	assert.Error(t, err)
}
