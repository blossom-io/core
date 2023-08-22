// Package app configures and runs application.
package app

import (
	"context"
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

	log := logger.New(cfg.App.LogLevel, logger.WithTracing())

	// tp, err := tracing.New()
	// if err != nil {
	// 	log.Fatalf("app - Run - tracing.New - error initializing tracing: %w", err)
	// }
	// defer func() { _ = tp.Shutdown(ctx) }()

	// log.Info("123")
	// log.Infof("123: %s", fmt.Errorf("123: %s", "456"))
	// log.Infow("123", "hello", "world")
	// ctx, span := tp.Tracer("123").Start(ctx, "operation")
	// span.AddEvent("Nice operation!", trace.WithAttributes(attribute.Int("bogons", 100)))
	// span.SetAttributes(attribute.String("hello", "world"))
	// defer span.End()

	// log.ErrorfContext(ctx, "ctx %v %v", "hello", "world")

	// Database
	DB, err := postgres.New(cfg.Connections.Postgres.URL)
	if err != nil {
		log.Fatalf("app - Run - postgres.New - error initializing database: %w", err)
	}
	defer DB.Close()

	// Repositories
	repo := repository.New(DB)

	// Twitch
	twitch, err := twitch.New(ctx, cfg.Connections.Twitch.ClientID, cfg.Connections.Twitch.ClientSecret, cfg.Connections.Twitch.AuthRedirectURL)
	if err != nil {
		log.Fatalf("app - Run - twitch.New - error initializing twitch client: %w", err)
	}

	// Services
	authSvc := service.NewAuth(log, repo, twitch)

	// HTTP Server
	r := chi.NewRouter()
	api.New(r, cfg, log, authSvc)
	HTTPSrv := httpserver.New(ctx, r, httpserver.Port(cfg.App.Port))
	defer HTTPSrv.Close()

	// span.End()

	<-ctx.Done()
	log.Info("Gracefully shutting down...")
}
