package main

import (
	"context"
	"log"

	"github.com/WebDeveloperBen/ai-gateway/internal/config"
	"github.com/WebDeveloperBen/ai-gateway/internal/migrate"
)

func main() {
	ctx := context.Background()

	// Load configuration
	config.Reload()

	log.Println("Starting database migration job...")

	// Run migrations
	if err := migrate.InitDatabase(ctx, config.Envs.GetDatabaseConnection()); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Database migration completed successfully")
}

