package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	proxyapi "github.com/insurgence-ai/llm-gateway/internal/api/proxy"
	"github.com/insurgence-ai/llm-gateway/internal/auth"
	"github.com/insurgence-ai/llm-gateway/internal/config"
	"github.com/insurgence-ai/llm-gateway/internal/gateway"
	keyspg "github.com/insurgence-ai/llm-gateway/internal/keys/postgres"
	"github.com/insurgence-ai/llm-gateway/internal/provider/azureopenai"
)

func main() {
	r := chi.NewRouter()
	cfg := huma.DefaultConfig("LLM Gateway", "v1.0.0")
	cfg.DocsPath = ""
	cfg.CreateHooks = nil
	api := humachi.New(r, cfg)

	pool, err := pgxpool.New(context.Background(), mustEnv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Shared keys store + hasher + authenticator
	keyStore := keyspg.New(pool) // implements keys.Store (Reader+Writer)

	// option A: default hasher
	authn := auth.NewDefaultAPIKeyAuthenticator(keyStore)

	// option B: explicit hasher params
	// hasher := keys.NewArgon2IDHasher(1, 64*1024, 1, 32)
	// authn  := auth.NewAPIKeyAuthenticator(keyStore, hasher)

	// Provider router (AOAI)
	aoai := azureopenai.New()
	aoaiBase := mustEnv("AOAI_BASE_URL")
	aoaiVer := envOr("AOAI_API_VERSION", "2024-07-01-preview")
	if dep := envOr("AOAI_DEPLOY_GPT4O", "gpt4o"); dep != "" {
		aoai.Global["gpt-4o"] = azureopenai.Entry{BaseURL: aoaiBase, Deployment: dep, APIVer: aoaiVer}
	}
	if dep := envOr("AOAI_DEPLOY_GPT4O_MINI", "gpt4o-mini"); dep != "" {
		aoai.Global["gpt-4o-mini"] = azureopenai.Entry{BaseURL: aoaiBase, Deployment: dep, APIVer: aoaiVer}
	}
	if dep := os.Getenv("AOAI_DEPLOY_DEFAULT"); dep != "" {
		aoai.Default = &azureopenai.Entry{BaseURL: aoaiBase, Deployment: dep, APIVer: aoaiVer}
	}

	// Transport chain
	rl := allowAllLimiter{}
	met := noopMetrics{}
	transport := gateway.Chain(
		http.DefaultTransport,
		gateway.WithAuth(authn),
		gateway.WithRateLimit(rl),
		gateway.WithMetrics(met),
	)

	core := gateway.NewCore(aoai, transport)

	grp := huma.NewGroup(api, "/api")
	proxyapi.RegisterRoutes(grp, core)

	addr := config.Envs.AppPort
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
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
