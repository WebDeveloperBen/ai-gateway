package keys_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/api/admin/keys"
	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	dbdriver "github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	keyrepo "github.com/WebDeveloperBen/ai-gateway/internal/repository/keys"
	"github.com/WebDeveloperBen/ai-gateway/internal/testkit"
	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/require"
)

func TestMintKey_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndApp(t)

	// Create a test user
	user, err := pg.Queries.CreateUser(context.Background(), db.CreateUserParams{
		OrgID: orgID,
		Email: "test@example.com",
	})
	require.NoError(t, err)

	keyRepo, err := keyrepo.NewKeyRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	hasher := keyrepo.NewArgon2IDHasher(1, 64*1024, 2, 32)
	svc := keys.NewService(keyRepo, hasher)
	router := keys.NewRouter(svc)

	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	reqBody := keys.MintKeyRequestBody{
		OrgID:    orgID.String(),
		AppID:    appID.String(),
		UserID:   user.ID.String(),
		Prefix:   "sk_test",
		Metadata: map[string]interface{}{},
	}
	body, _ := json.Marshal(reqBody)

	resp := api.Post("/api/keys", "Content-Type: application/json", bytes.NewReader(body))

	require.Equal(t, http.StatusCreated, resp.Code, resp.Body.String())

	var responseBody struct {
		Token string      `json:"token"`
		Key   keys.APIKey `json:"key"`
	}
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.NotEmpty(t, responseBody.Token)
	require.True(t, strings.HasPrefix(responseBody.Token, "sk_test_"))
	require.Contains(t, responseBody.Token, ".")

	require.Equal(t, orgID.String(), responseBody.Key.Tenant)
	require.Equal(t, appID.String(), responseBody.Key.App)
	require.Equal(t, model.KeyActive, responseBody.Key.Status)
	require.NotEmpty(t, responseBody.Key.LastFour)
}

func TestMintKey_WithTTL(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndAppWithSuffix(t, "ttl")

	keyRepo, err := keyrepo.NewKeyRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	hasher := keyrepo.NewArgon2IDHasher(1, 64*1024, 2, 32)
	svc := keys.NewService(keyRepo, hasher)
	router := keys.NewRouter(svc)

	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	// Create a test user
	user, err := pg.Queries.CreateUser(context.Background(), db.CreateUserParams{
		OrgID: orgID,
		Email: "test-ttl@example.com",
	})
	require.NoError(t, err)

	ttl := 24 * time.Hour
	reqBody := keys.MintKeyRequestBody{
		OrgID:    orgID.String(),
		AppID:    appID.String(),
		UserID:   user.ID.String(),
		TTL:      ttl,
		Metadata: map[string]interface{}{},
	}
	body, _ := json.Marshal(reqBody)

	resp := api.Post("/api/keys", "Content-Type: application/json", bytes.NewReader(body))

	require.Equal(t, http.StatusCreated, resp.Code, resp.Body.String())

	var responseBody struct {
		Token string      `json:"token"`
		Key   keys.APIKey `json:"key"`
	}
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.NotNil(t, responseBody.Key.ExpiresAt)
	expectedExpiry := time.Now().Add(ttl)
	require.WithinDuration(t, expectedExpiry, *responseBody.Key.ExpiresAt, 5*time.Second)
}

func TestMintKey_ValidationErrors(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	keyRepo, err := keyrepo.NewKeyRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	hasher := keyrepo.NewArgon2IDHasher(1, 64*1024, 2, 32)
	svc := keys.NewService(keyRepo, hasher)
	router := keys.NewRouter(svc)

	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	tests := []struct {
		name       string
		reqBody    keys.MintKeyRequestBody
		wantStatus int
	}{
		{
			name: "invalid orgID",
			reqBody: keys.MintKeyRequestBody{
				OrgID:    "",
				AppID:    "app-123",
				UserID:   "user-123",
				Metadata: map[string]interface{}{},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid appID",
			reqBody: keys.MintKeyRequestBody{
				OrgID:    "org-123",
				AppID:    "",
				UserID:   "user-123",
				Metadata: map[string]interface{}{},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid userID",
			reqBody: keys.MintKeyRequestBody{
				OrgID:    "org-123",
				AppID:    "app-123",
				UserID:   "",
				Metadata: map[string]interface{}{},
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.reqBody)
			resp := api.Post("/api/keys", "Content-Type: application/json", bytes.NewReader(body))
			require.Equal(t, tt.wantStatus, resp.Code)
		})
	}
}

func TestRevokeKey_Integration(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	fixtures := testkit.NewDBFixtures(pg.Queries, pg.Pool)
	orgID, appID := fixtures.CreateTestOrgAndAppWithSuffix(t, "revoke")

	// Create a test user
	user, err := pg.Queries.CreateUser(context.Background(), db.CreateUserParams{
		OrgID: orgID,
		Email: "test-revoke@example.com",
	})
	require.NoError(t, err)

	keyRepo, err := keyrepo.NewKeyRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	hasher := keyrepo.NewArgon2IDHasher(1, 64*1024, 2, 32)
	svc := keys.NewService(keyRepo, hasher)
	router := keys.NewRouter(svc)

	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	reqBody := keys.MintKeyRequestBody{
		OrgID:    orgID.String(),
		AppID:    appID.String(),
		UserID:   user.ID.String(),
		Metadata: map[string]interface{}{},
	}
	body, _ := json.Marshal(reqBody)

	mintResp := api.Post("/api/keys", "Content-Type: application/json", bytes.NewReader(body))
	require.Equal(t, http.StatusCreated, mintResp.Code)

	var mintResult struct {
		Token string      `json:"token"`
		Key   keys.APIKey `json:"key"`
	}
	err = json.Unmarshal(mintResp.Body.Bytes(), &mintResult)
	require.NoError(t, err)

	keyID := mintResult.Key.KeyID

	revokeResp := api.Post("/api/keys/" + keyID + "/revoke")

	require.Equal(t, http.StatusNoContent, revokeResp.Code)

	key, err := keyRepo.GetByKeyPrefix(context.Background(), keyID)
	require.NoError(t, err)
	require.Equal(t, model.KeyRevoked, key.Status)
}

func TestRevokeKey_NotFound(t *testing.T) {
	pgConnStr, _ := testkit.SetupTestContainers(t, testkit.DefaultContainerConfig())

	pg, err := dbdriver.NewPostgresDriver(context.Background(), pgConnStr)
	require.NoError(t, err)
	defer pg.Pool.Close()

	keyRepo, err := keyrepo.NewKeyRepository(context.Background(), model.RepositoryConfig{
		Backend: model.RepositoryPostgres,
		PGPool:  pg.Pool,
	})
	require.NoError(t, err)

	hasher := keyrepo.NewArgon2IDHasher(1, 64*1024, 2, 32)
	svc := keys.NewService(keyRepo, hasher)
	router := keys.NewRouter(svc)

	api := testkit.SetupPublicTestAPI(t, func(grp *huma.Group) {
		router.RegisterRoutes(grp)
	})

	resp := api.Post("/api/keys/nonexistent-key-id/revoke")

	require.Equal(t, http.StatusNoContent, resp.Code)
}
