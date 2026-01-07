package main

import (
	"fmt"
	"log"
	"net/http"
	"sre-chat-api/internal"
	"sre-chat-api/internal/config"
	"sre-chat-api/internal/database"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Run migrations if enabled
	if err := database.Migrate(cfg.MigrationConfig.Enabled); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Seed default group (only if migrations are enabled)
	if cfg.MigrationConfig.Enabled {
		if err := database.SeedDefaultGroup(); err != nil {
			logger.Fatal("Failed to seed default group", zap.Error(err))
		}
	}

	// Set up router
	router := internal.SetupRouter(logger)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Info("Starting server", zap.String("address", addr))

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
