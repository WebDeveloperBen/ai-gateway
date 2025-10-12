package main

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"

	apiauth "github.com/WebDeveloperBen/ai-gateway/internal/api/auth"
	"github.com/WebDeveloperBen/ai-gateway/internal/api/docs"
	apigw "github.com/WebDeveloperBen/ai-gateway/internal/api/gateway"
	"github.com/WebDeveloperBen/ai-gateway/internal/api/health"
	"github.com/WebDeveloperBen/ai-gateway/internal/api/middleware"
	"github.com/WebDeveloperBen/ai-gateway/internal/config"
	dbdriver "github.com/WebDeveloperBen/ai-gateway/internal/drivers/db"
	"github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/auth"
	gwmiddleware "github.com/WebDeveloperBen/ai-gateway/internal/gateway/middleware"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/policies"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/WebDeveloperBen/ai-gateway/internal/provider"

	"github.com/WebDeveloperBen/ai-gateway/internal/api/admin/applications"
	"github.com/WebDeveloperBen/ai-gateway/internal/api/admin/keys"
	adminpolicies "github.com/WebDeveloperBen/ai-gateway/internal/api/admin/policies"
	adminusage "github.com/WebDeveloperBen/ai-gateway/internal/api/admin/usage"
	apprepo "github.com/WebDeveloperBen/ai-gateway/internal/repository/applications"
	keyrepo "github.com/WebDeveloperBen/ai-gateway/internal/repository/keys"
	orgrepo "github.com/WebDeveloperBen/ai-gateway/internal/repository/organisations"
	policiesrepo "github.com/WebDeveloperBen/ai-gateway/internal/repository/policies"
	usagerepo "github.com/WebDeveloperBen/ai-gateway/internal/repository/usage"
	userrepo "github.com/WebDeveloperBen/ai-gateway/internal/repository/users"
	"github.com/WebDeveloperBen/ai-gateway/internal/server"
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
		RedisDB:   0,
	})
	if err != nil {
		log.Fatal(err)
	}

	if cfg.EnableRedisCircuitBreaker && cfg.KVBackend == "redis" {
		kvStore = kv.NewCircuitBreakerStore(kvStore, kv.DefaultCircuitBreakerConfig())
		log.Println("Redis circuit breaker enabled")
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
	appRepo := apprepo.NewPostgresRepo(pg.Queries)
	orgRepo := orgrepo.NewPostgresRepo(pg.Queries)
	policiesRepo := policiesrepo.NewPostgresRepo(pg.Queries)
	usageRepo := usagerepo.NewPostgresRepo(pg.Queries)
	userRepo := userrepo.NewPostgresRepo(pg.Queries)

	// ---------- Middleware Utilities -------- //
	authn := auth.NewDefaultAPIKeyAuthenticator(keyRepo)

	// OIDC Service (use mock for development/testing)
	var oidcService apiauth.OIDCServiceInterface
	if cfg.AppRegistrationTenantID != "dummy-tenant-id" {
		var err error
		oidcService, err = apiauth.NewOIDCService(ctx,
			apiauth.OIDCConfig{
				ClientID:     cfg.AppRegistrationClientID,
				ClientSecret: cfg.AppRegistrationClientSecret,
				TenantID:     cfg.AppRegistrationTenantID,
				RedirectURL:  cfg.AppRegistrationRedirectURL,
			})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Using mock OIDC service for development/testing")
		oidcService = apiauth.NewMockOIDCService()
	}

	// ---------- Policy Engine -------- //
	policyEngine := policies.NewEngine(pg.Queries, kvStore)
	requestBuffer := gwmiddleware.NewRequestBuffer()
	policyEnforcer := gwmiddleware.NewPolicyEnforcer(policyEngine)
	usageRecorder := gwmiddleware.NewUsageRecorder(pg.Queries, policyEngine)

	// ------------- Services ------------ //
	orgSvc := apiauth.NewOrganisationService(orgRepo, userRepo)
	keysSvc := keys.NewService(keyRepo, hasher)
	appsSvc := applications.NewService(appRepo)
	policiesSvc := adminpolicies.NewService(policiesRepo)
	usageSvc := adminusage.NewService(usageRepo)

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
	applications.NewRouter(appsSvc).RegisterRoutes(admingrp)
	adminpolicies.NewRouter(policiesSvc).RegisterRoutes(admingrp)
	adminusage.NewRouter(usageSvc).RegisterRoutes(admingrp)

	// ------------ Gateway Proxy Setup ----------- //
	transport := gateway.Chain(
		http.DefaultTransport,
		gateway.WithAuth(authn),
		requestBuffer.Middleware,  // Buffer request body once
		policyEnforcer.Middleware, // Policy enforcement (pre-check)
		usageRecorder.Middleware,  // Usage recording (post-check, async)
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
