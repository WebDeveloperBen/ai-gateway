package main

import (
	"context"
	"log"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	"github.com/insurgence-ai/llm-gateway/internal/api/admin/keys"
	"github.com/insurgence-ai/llm-gateway/internal/api/auth"
	"github.com/insurgence-ai/llm-gateway/internal/api/docs"
	"github.com/insurgence-ai/llm-gateway/internal/api/health"
	"github.com/insurgence-ai/llm-gateway/internal/config"
	"github.com/insurgence-ai/llm-gateway/internal/drivers/db"
	"github.com/insurgence-ai/llm-gateway/internal/model"
	keyrepo "github.com/insurgence-ai/llm-gateway/internal/repository/keys"
	"github.com/insurgence-ai/llm-gateway/internal/server"
)

func main() {
	cfg := config.Envs
	ctx := context.Background()

	// Utilities
	hasher := keyrepo.NewArgon2IDHasher(1, 64*1024, 1, 32)

	// Drivers
	pg, err := db.NewPostgresDriver(ctx, cfg.DBConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Pool.Close()

	// Repositories
	keyStore, err := keyrepo.NewKeyRepository(ctx, model.RepositoryConfig{Backend: model.RepositoryBackend(cfg.DBBackend), PGPool: pg.Pool})
	if err != nil {
		log.Fatal(err)
	}

	// Services
	keysSvc := keys.NewService(keyStore, hasher)

	// Routers
	r := chi.NewRouter()

	api := humachi.New(r, huma.DefaultConfig("Admin API", "v1"))
	protected := huma.NewGroup(api, "/api/v1")
	// TODO: fix this
	// protected.UseMiddleware(middleware.Use(api, middleware.AuthenticationMiddleware))
	grp := huma.NewGroup(api, "/api")

	oidcService, err := auth.NewOIDCService(ctx, auth.OIDCConfig{ClientID: cfg.AppRegistrationClientID, ClientSecret: cfg.AppRegistrationClientSecret, TenantID: cfg.AppRegistrationTenantID, RedirectURL: cfg.AppRegistrationRedirectURL})
	if err != nil {
		log.Fatal(err)
	}

	auth.RegisterAuthRoutes(protected, oidcService)

	// Routes
	docs.RegisterRoutes(r)

	health.RegisterPublicRoutes(api)

	keys.NewRouter(keysSvc).RegisterRoutes(grp)

	// Server Start
	addr := config.Envs.ProxyPort
	server.Start(addr, r)
}
