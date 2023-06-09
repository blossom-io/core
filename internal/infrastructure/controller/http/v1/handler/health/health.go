package health

import (
	"database/sql"
	"net/http"
	"runtime"
	"time"

	"core/internal/infrastructure/controller/http/v1/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type healthRoutes struct {
	Uptime         time.Time
	LogLevel       string
	ServiceVersion string
	DB             *sql.DB
}

type Health struct {
	NumCPU       int         `json:"numCpu"`
	Uptime       string      `json:"uptime"`
	StartDate    string      `json:"startDate"`
	OK           bool        `json:"ok"`
	DB           sql.DBStats `json:"db"`
	NumGoroutine int         `json:"numGoroutine"`
}

func New(r chi.Router) {
	h := &healthRoutes{
		Uptime: time.Now(),
		// DB:     db
	}

	r.Get("/ping", h.Ping)
	r.Get("/health", h.Health)
}

func (h *healthRoutes) Ping(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "pong")
}

func (h *healthRoutes) Health(w http.ResponseWriter, r *http.Request) {
	result := &Health{
		NumCPU:       runtime.NumCPU(),
		Uptime:       time.Since(h.Uptime).String(),
		StartDate:    h.Uptime.Format(time.RFC1123),
		OK:           true,
		DB:           h.DB.Stats(),
		NumGoroutine: runtime.NumGoroutine(),
	}

	render.JSON(w, r, response.Response{Data: result})
}
