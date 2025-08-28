package main

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"

	apiauth "github.com/insurgence-ai/llm-gateway/internal/api/auth"
	"github.com/insurgence-ai/llm-gateway/internal/api/docs"
	apigw "github.com/insurgence-ai/llm-gateway/internal/api/gateway"
	"github.com/insurgence-ai/llm-gateway/internal/api/health"
	"github.com/insurgence-ai/llm-gateway/internal/api/middleware"
	"github.com/insurgence-ai/llm-gateway/internal/config"
	dbdriver "github.com/insurgence-ai/llm-gateway/internal/drivers/db"
	"github.com/insurgence-ai/llm-gateway/internal/drivers/kv"
	"github.com/insurgence-ai/llm-gateway/internal/gateway"
	"github.com/insurgence-ai/llm-gateway/internal/gateway/auth"
	"github.com/insurgence-ai/llm-gateway/internal/model"
	"github.com/insurgence-ai/llm-gateway/internal/provider"

	"github.com/insurgence-ai/llm-gateway/internal/api/admin/keys"
	keyrepo "github.com/insurgence-ai/llm-gateway/internal/repository/keys"
	orgrepo "github.com/insurgence-ai/llm-gateway/internal/repository/organisations"
	userrepo "github.com/insurgence-ai/llm-gateway/internal/repository/users"
	"github.com/insurgence-ai/llm-gateway/internal/server"
)

func main() {
	ctx := context.Background()
	cfg := config.Envs

	// ------------- Utilities ------------ //
	hasher := keyrepo.NewArgon2IDHasher(1, 64*1024, 1, 32)

	// ------------- Drivers ------------ //
	kvStore, err := kv.NewDriver(kv.Config{
		Backend:   kv.KvStoreType(cfg.KVBackend),
		RedisAddr: cfg.RedisAddr,
		RedisPW:   cfg.RedisPW,
		RedisDB:   0, // use first database - change if necessary
	})
	if err != nil {
		log.Fatal(err)
	}

	pg, err := dbdriver.NewPostgresDriver(ctx, cfg.DBConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Pool.Close()

	// --------------- Registry ------------- //
	reg := gateway.NewRegistry(ctx, kvStore)
	_ = gateway.EnsureRegistryPopulated(reg, loadAllModelDeploymentsFromDatabase)

	// ------------- Repositories ------------ //
	keyRepo, err := keyrepo.NewKeyRepository(ctx,
		model.RepositoryConfig{
			Backend: model.RepositoryBackend(cfg.DBBackend),
			PGPool:  pg.Pool,
		})
	if err != nil {
		log.Fatal(err)
	}
	orgRepo := orgrepo.NewPostgresRepo(pg.Queries)
	userRepo := userrepo.NewPostgresRepo(pg.Queries)

	// ---------- Middleware Utilities -------- //
	authn := auth.NewDefaultAPIKeyAuthenticator(keyRepo)
	oidcService, err := apiauth.NewOIDCService(ctx,
		apiauth.OIDCConfig{
			ClientID:     cfg.AppRegistrationClientID,
			ClientSecret: cfg.AppRegistrationClientSecret,
			TenantID:     cfg.AppRegistrationTenantID,
			RedirectURL:  cfg.AppRegistrationRedirectURL,
		})
	if err != nil {
		log.Fatal(err)
	}

	// ------------- Services ------------ //
	orgSvc := apiauth.NewOrganisationService(orgRepo, userRepo)
	keysSvc := keys.NewService(keyRepo, hasher)

	// ----------- API Router Setup ---------- //
	router, humaCfg := server.New(cfg)
	base := humachi.New(router, humaCfg)

	// ----------- Register Public Routers ---------- //
	docs.RegisterRoutes(router)
	health.RegisterPublicRoutes(base)

	// ------------ Route Groups ----------- //
	apigrp := huma.NewGroup(base, "/api")
	admingrp := huma.NewGroup(base, "/api/v1/admin")
	authgrp := huma.NewGroup(base, "/auth")

	// --- Attach Authentication Middleware to Route Groups --- //
	authgrp.UseMiddleware(
		middleware.Use(base, middleware.AuthCookieMiddleware),
	)

	admingrp.UseMiddleware(
		middleware.Use(base, middleware.AuthCookieMiddleware),
		middleware.Use(base, middleware.WithScopedOrg(pg.Pool)),
	)
	// protected.UseMiddleware(middleware.Use(base, middleware.AuthenticationMiddleware))

	// ------------ Register API Routers ----------- //
	apiauth.NewRouter(oidcService, orgSvc).RegisterRoutes(authgrp)
	keys.NewRouter(keysSvc).RegisterRoutes(authgrp)

	// ------------ Gateway Proxy Setup ----------- //
	transport := gateway.Chain(
		http.DefaultTransport,
		gateway.WithAuth(authn),
		// TODO: add the ratelimiter
		// TODO: add the policies
		// TODO: add the load balancer
	)
	core := gateway.NewCoreWithRegistry(transport, authn, reg)

	// ------------ AI Providers ----------- //
	apigw.RegisterProvider(apigrp, provider.AzureOpenAIPrefix, core)

	// ------------ Server Start ----------- //
	addr := config.Envs.ProxyPort
	server.Start(addr, router)
}

// loadAllModelDeploymentsFromDatabase is a stub for demonstration. Replace with real DB logic.
func loadAllModelDeploymentsFromDatabase() []model.ModelDeployment {
	return []model.ModelDeployment{
		{Model: "gpt-4.1", Deployment: "dev-openai-gpt4-1", Provider: "azure", Tenant: "default", Meta: map[string]string{"APIVer": "2024-07-01-preview", "BaseURL": "https://dev-insurgence-openai.openai.azure.com"}},
	}
}
