package main

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/insurgence-ai/llm-gateway/internal/api/docs"
	"github.com/insurgence-ai/llm-gateway/internal/api/health"
	"github.com/insurgence-ai/llm-gateway/internal/api/proxy"
	"github.com/insurgence-ai/llm-gateway/internal/auth"
	"github.com/insurgence-ai/llm-gateway/internal/config"
	"github.com/insurgence-ai/llm-gateway/internal/gateway"
	keyspg "github.com/insurgence-ai/llm-gateway/internal/keys/postgres"
	"github.com/insurgence-ai/llm-gateway/internal/kv"
	"github.com/insurgence-ai/llm-gateway/internal/model"
	"github.com/insurgence-ai/llm-gateway/internal/server"
)

func main() {
	ctx := context.Background()
	cfg := config.Envs

	// KV store
	kvStore, err := kv.New(kv.Config{
		Backend:   kv.Backend(cfg.KVBackend),
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

	keyStore := keyspg.New(pool)
	authn := auth.NewDefaultAPIKeyAuthenticator(keyStore)

	router, humaCfg := server.New(cfg)
	api := humachi.New(router, humaCfg)

	docs.RegisterRoutes(router)

	health.RegisterPublicRoutes(api)

	transport := gateway.Chain(
		http.DefaultTransport,
		gateway.WithAuth(authn),
		// TODO: real limiter and metrics
	)
	core := gateway.NewCoreWithRegistry(transport, authn, reg)
	grp := huma.NewGroup(api, "/api")
	proxy.RegisterProvider(grp, "/azure/openai", core)

	addr := config.Envs.AppPort
	server.Start(addr, router)
}

// loadAllModelDeploymentsFromDatabase is a stub for demonstration. Replace with real DB logic.
func loadAllModelDeploymentsFromDatabase() []model.ModelDeployment {
	return []model.ModelDeployment{
		{Model: "gpt-4.1", Deployment: "dev-openai-gpt4-1", Provider: "azure", Tenant: "default", Meta: map[string]string{"APIVer": "2024-07-01-preview", "BaseURL": "https://dev-insurgence-openai.openai.azure.com"}},
	}
}
