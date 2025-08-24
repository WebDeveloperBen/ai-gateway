package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/jackc/pgx/v5/pgxpool"

	proxyapi "github.com/insurgence-ai/llm-gateway/internal/api/proxy"
	"github.com/insurgence-ai/llm-gateway/internal/auth"
	"github.com/insurgence-ai/llm-gateway/internal/config"
	"github.com/insurgence-ai/llm-gateway/internal/gateway"
	keyspg "github.com/insurgence-ai/llm-gateway/internal/keys/postgres"
	"github.com/insurgence-ai/llm-gateway/internal/provider/azureopenai"
	"github.com/insurgence-ai/llm-gateway/internal/server"
)

func main() {
	// Load environment configuration
	cfg := config.Envs

	// Create a Chi router with all middleware (logging, rate limiting, CORS, etc.)
	router, humaCfg := server.New(cfg)

	// Wrap the router with Huma to support OpenAPI + typed handlers
	api := humachi.New(router, humaCfg)

	pool, err := pgxpool.New(context.Background(), mustEnv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Shared keys store + hasher + authenticator
	keyStore := keyspg.New(pool) // implements keys.Store (Reader+Writer)

	authn := auth.NewDefaultAPIKeyAuthenticator(keyStore)

	// Transport chain
	rl := allowAllLimiter{}
	met := noopMetrics{}
	transport := gateway.Chain(
		http.DefaultTransport,
		gateway.WithAuth(authn),
		gateway.WithRateLimit(rl),
		gateway.WithMetrics(met),
	)
	aoai := azureopenai.New()
	aoai.Global["gpt-4o"] = azureopenai.Entry{
		BaseURL:    "https://<your-aoai-resource>.openai.azure.com",
		Deployment: "<your-deployment-name>",
		APIVer:     "2024-07-01-preview",
	}
	core := gateway.NewCoreWithAdapters(transport, aoai)
	grp := huma.NewGroup(api, "/api")
	proxyapi.RegisterProvider(grp, "/azure/openai/", core)

	addr := config.Envs.AppPort

	server.Start(addr, router)
}

// demo deps
type allowAllLimiter struct{}

func (allowAllLimiter) Allow(*http.Request) (time.Duration, bool) { return 0, true }

type noopMetrics struct{}

func (noopMetrics) Record(*http.Request, *http.Response, time.Duration) {}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing %s", k)
	}
	return v
}
