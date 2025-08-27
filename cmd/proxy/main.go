package main

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/jackc/pgx/v5/pgxpool"

	middleware "github.com/insurgence-ai/llm-gateway/internal/api"
	apiauth "github.com/insurgence-ai/llm-gateway/internal/api/auth"
	"github.com/insurgence-ai/llm-gateway/internal/api/docs"
	"github.com/insurgence-ai/llm-gateway/internal/api/health"
	"github.com/insurgence-ai/llm-gateway/internal/api/proxy"
	"github.com/insurgence-ai/llm-gateway/internal/config"
	"github.com/insurgence-ai/llm-gateway/internal/drivers/kv"
	"github.com/insurgence-ai/llm-gateway/internal/gateway"
	"github.com/insurgence-ai/llm-gateway/internal/gateway/auth"
	"github.com/insurgence-ai/llm-gateway/internal/model"
	"github.com/insurgence-ai/llm-gateway/internal/repository/keys"
	"github.com/insurgence-ai/llm-gateway/internal/server"
)

func main() {
	ctx := context.Background()
	cfg := config.Envs

	// KV store
	kvStore, err := kv.NewDriver(kv.Config{
		Backend:   kv.KvStoreType(cfg.KVBackend),
		RedisAddr: cfg.RedisAddr,
		RedisPW:   cfg.RedisPW,
		RedisDB:   0, // change as needed
	})
	if err != nil {
		log.Fatal(err)
	}

	// postgres
	pool, err := pgxpool.New(context.Background(), cfg.DBConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Registry
	reg := gateway.NewRegistry(ctx, kvStore)
	_ = gateway.EnsureRegistryPopulated(reg, loadAllModelDeploymentsFromDatabase)

	keyStore := keys.NewPostgresStore(pool)

	authn := auth.NewDefaultAPIKeyAuthenticator(keyStore)

	router, humaCfg := server.New(cfg)
	base := humachi.New(router, humaCfg)

	docs.RegisterRoutes(router)

	health.RegisterPublicRoutes(base)

	// Route groups
	// protected := huma.NewGroup(base, "/api/v1/admin")
	public := huma.NewGroup(base, "/api")
	auth := huma.NewGroup(base, "/auth")
	auth.UseMiddleware(middleware.Use(base, middleware.AuthCookieMiddleware))
	// protected.UseMiddleware(middleware.Use(base, middleware.AuthenticationMiddleware))

	oidcService, err := apiauth.NewOIDCService(ctx, apiauth.OIDCConfig{ClientID: cfg.AppRegistrationClientID, ClientSecret: cfg.AppRegistrationClientSecret, TenantID: cfg.AppRegistrationTenantID, RedirectURL: cfg.AppRegistrationRedirectURL})
	if err != nil {
		log.Fatal(err)
	}

	apiauth.RegisterAuthRoutes(auth, oidcService)

	transport := gateway.Chain(
		http.DefaultTransport,
		gateway.WithAuth(authn),
		// TODO: real limiter and metrics
	)

	core := gateway.NewCoreWithRegistry(transport, authn, reg)
	proxy.RegisterProvider(public, "/azure/openai", core)

	addr := config.Envs.ProxyPort
	server.Start(addr, router)
}

// loadAllModelDeploymentsFromDatabase is a stub for demonstration. Replace with real DB logic.
func loadAllModelDeploymentsFromDatabase() []model.ModelDeployment {
	return []model.ModelDeployment{
		{Model: "gpt-4.1", Deployment: "dev-openai-gpt4-1", Provider: "azure", Tenant: "default", Meta: map[string]string{"APIVer": "2024-07-01-preview", "BaseURL": "https://dev-insurgence-openai.openai.azure.com"}},
	}
}
