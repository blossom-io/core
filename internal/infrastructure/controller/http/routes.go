package api

import (
	"core/internal/config"
	"core/internal/infrastructure/controller/http/v1/handler/auth"
	"core/internal/infrastructure/controller/http/v1/handler/health"
	"core/internal/service"
	"core/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// New creates routes
func New(r *chi.Mux, cfg *config.Config, log logger.Logger, authSvc service.Auther) {
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))
	// r.Use(middleware.Timeout(3 * time.Second))

	NewRoutesV1(r, cfg, log, authSvc)
}

// NewRoutesV1 creates v1 API routes
func NewRoutesV1(r *chi.Mux, cfg *config.Config, log logger.Logger, authSvc service.Auther) {
	r.Route("/v1", func(r chi.Router) {
		health.New(r)
		auth.New(r, cfg, authSvc, log)
	})
}
