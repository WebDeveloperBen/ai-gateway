// Package server implements the server of the applications
package server

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/insurgence-ai/llm-gateway/internal/config"
	"github.com/insurgence-ai/llm-gateway/internal/logger"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

// New initializes the router with all standard middleware.
func New(cfg config.Config) (*chi.Mux, huma.Config) {
	logger.NewLogger(cfg.IsProd)

	router := chi.NewRouter()

	// Global middlewares
	router.Use(chi_middleware.Logger)
	router.Use(chi_middleware.Recoverer)
	router.Use(chi_middleware.RedirectSlashes)
	router.Use(chi_middleware.CleanPath)
	router.Use(chi_middleware.Compress(5, "application/json", "text/html", "text/css", "application/javascript"))

	// FIX: make this better for prod deployed api
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Setup Huma with that one router
	humaCfg := huma.DefaultConfig(cfg.ApplicationName, cfg.Version)
	humaCfg.DocsPath = ""     // hide the default docs to create our own scalar docs
	humaCfg.CreateHooks = nil // remove the $schema from being returned in api responses

	return router, humaCfg
}

func Start(port string, router *chi.Mux) {
	// guard agaisnt input strings including or not including colons
	if len(port) > 0 && port[0] == ':' {
		port = port[1:]
	}

	addr := ":" + port

	logger.Logger.Info().Msgf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to start server")
	}
}
