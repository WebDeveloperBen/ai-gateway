package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/WebDeveloperBen/ai-gateway/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}

func NewPostgresDriver(ctx context.Context, dsn string) (*Postgres, error) {
	var connectionString string
	var err error

	// Check if DSN looks like a connection string or Azure managed identity config
	if strings.Contains(dsn, "://") || strings.Contains(dsn, "host=") {
		// Standard connection string
		connectionString = dsn
	} else {
		// Assume Azure managed identity format: "server:database:user"
		connectionString, err = buildManagedIdentityConnection(ctx, dsn)
		if err != nil {
			return nil, fmt.Errorf("[Fail]: Failed to build managed identity connection: %w", err)
		}
	}

	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("[Fail]: Postgres driver unable to connect: %+w", err)
	}
	return &Postgres{Pool: pool, Queries: db.New(pool)}, nil
}

func buildManagedIdentityConnection(ctx context.Context, config string) (string, error) {
	parts := strings.Split(config, ":")
	if len(parts) != 3 {
		return "", fmt.Errorf("managed identity config must be in format 'server:database:user'")
	}

	server, database, user := parts[0], parts[1], parts[2]

	// Get Azure managed identity credential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return "", fmt.Errorf("failed to obtain Azure credential: %v", err)
	}

	// Get access token for PostgreSQL
	token, err := cred.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"https://ossrdbms-aad.database.windows.net/.default"},
	})
	if err != nil {
		return "", fmt.Errorf("failed to obtain access token: %v", err)
	}

	// Build connection string with token as password
	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=require",
		server, user, token.Token, database,
	)

	return connectionString, nil
}
