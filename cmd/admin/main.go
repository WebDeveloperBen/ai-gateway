package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/insurgence-ai/llm-gateway/internal/admin/services"
	"github.com/insurgence-ai/llm-gateway/internal/api/admin"
	"github.com/insurgence-ai/llm-gateway/internal/api/docs"
	"github.com/insurgence-ai/llm-gateway/internal/api/health"
	"github.com/insurgence-ai/llm-gateway/internal/keys"
	"github.com/insurgence-ai/llm-gateway/internal/keys/postgres"
)

func main() {
	pool, err := pgxpool.New(context.Background(), mustEnv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// shared keys store + hasher
	keyStore := postgres.New(pool)                      // implements keys.Store
	hasher := keys.NewArgon2IDHasher(1, 64*1024, 1, 32) // t=1, m=64MiB, p=1, 32-byte key

	// admin service
	keysSvc := services.NewKeysService(keyStore, hasher)

	// http
	r := chi.NewRouter()
	cfg := huma.DefaultConfig("Admin API", "v1")

	api := humachi.New(r, cfg)

	grp := huma.NewGroup(api, "/api")

	docs.RegisterRoutes(r)

	health.RegisterPublicRoutes(api)

	admin.NewServer(keysSvc).RegisterRoutes(grp)

	addr := envOr("ADMIN_ADDR", ":8081")
	log.Printf("admin listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing required env %s", k)
	}
	return v
}
