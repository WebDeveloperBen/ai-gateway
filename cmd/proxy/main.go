package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/jackc/pgx/v5/pgxpool"

	proxyapi "github.com/insurgence-ai/llm-gateway/internal/api/proxy"
	"github.com/insurgence-ai/llm-gateway/internal/auth"
	"github.com/insurgence-ai/llm-gateway/internal/config"
	"github.com/insurgence-ai/llm-gateway/internal/gateway"
	keyspg "github.com/insurgence-ai/llm-gateway/internal/keys/postgres"
	"github.com/insurgence-ai/llm-gateway/internal/kv"
	"github.com/insurgence-ai/llm-gateway/internal/model/models"
	"github.com/insurgence-ai/llm-gateway/internal/server"
)

func main() {
	ctx := context.Background()
	cfg := config.Envs

	// KV store
	kvCfg := kv.Config{
		Backend:   kv.Backend(cfg.KVBackend),
		RedisAddr: cfg.RedisAddr,
		RedisPW:   cfg.RedisPW,
		RedisDB:   0, // change as needed
	}
	kvStore, err := kv.New(kvCfg)
	if err != nil {
		log.Fatal(err)
	}

	reg := gateway.NewRegistry(ctx, kvStore)
	all, err := reg.All("modelreg:*")
	if err != nil {
		log.Fatal("registry read failed: ", err)
	}
	if len(all) == 0 {
		modelDeployments := loadAllModelDeploymentsFromDatabase()
		for _, md := range modelDeployments {
			if err := reg.Add(md, 0); err != nil {
				log.Printf("failed to add model to registry: %+v", err)
			}
		}
		all = modelDeployments
	}
	log.Printf("Loaded %d active models from registry", len(all))

	pool, err := pgxpool.New(context.Background(), mustEnv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	keyStore := keyspg.New(pool)
	authn := auth.NewDefaultAPIKeyAuthenticator(keyStore)

	router, humaCfg := server.New(cfg)
	api := humachi.New(router, humaCfg)

	transport := gateway.Chain(
		http.DefaultTransport,
		gateway.WithAuth(authn),
		// TODO: real limiter and metrics
	)
	core := gateway.NewCoreWithRegistry(transport, reg)
	grp := huma.NewGroup(api, "/api")
	proxyapi.RegisterProvider(grp, "/azure/openai/", core)

	addr := config.Envs.AppPort
	server.Start(addr, router)
}

// loadAllModelDeploymentsFromDatabase is a stub for demonstration. Replace with real DB logic.
func loadAllModelDeploymentsFromDatabase() []models.ModelDeployment {
	return []models.ModelDeployment{
		{Model: "gpt-4o", Deployment: "dev-openai-gpt4-1", Provider: "azure", Tenant: "default", Meta: map[string]string{"APIVer": "2024-07-01-preview", "BaseURL": "https://<your-aoai-resource>.openai.azure.com"}},
	}
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing %s", k)
	}
	return v
}
