package router

import (
	"github.com/chiaf1/iot-nonna-core/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Create Chi router with handlers
func Setup(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", h.HandleHealth)

	routeRooms(r, h)
	routeSensorType(r, h)
	routeDeviceType(r, h)

	return r
}
