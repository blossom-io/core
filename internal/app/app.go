// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"core/internal/config"
	"core/internal/infrastructure/repository"
	"core/internal/service"
	"core/pkg/httpserver"
	"core/pkg/logger"
	"core/pkg/postgres"
	"core/pkg/twitch"

	api "core/internal/infrastructure/controller/http"

	"github.com/go-chi/chi/v5"
)

// Run injects dependencies and runs application.
func Run(cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log := logger.New(cfg.App.LogLevel)

	// Database
	DB, err := postgres.New(cfg.Connections.Postgres.URL)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New - error initializing database: %w", err))
	}
	defer DB.Close()

	// Repositories
	repo := repository.New(DB)

	// Twitch
	twitch, err := twitch.New(ctx, cfg.Connections.Twitch.ClientID, cfg.Connections.Twitch.ClientSecret, cfg.Connections.Twitch.AuthRedirectURL)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - twitch.New - error initializing twitch client: %w", err))
	}

	// Services
	authSvc := service.NewAuth(log, repo, twitch)

	// HTTP Server
	r := chi.NewRouter()
	api.New(r, cfg, log, authSvc)
	HTTPSrv := httpserver.New(ctx, r, httpserver.Port(cfg.App.Port))
	defer HTTPSrv.Close()

	<-ctx.Done()
	log.Info("Gracefully shutting down...")
}
