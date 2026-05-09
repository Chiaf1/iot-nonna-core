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

	r.Route("/rooms", func(r chi.Router) {
		r.Get("/", h.GetRooms)
		r.Post("/", h.PostRoom)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetRoom)
			r.Put("/", h.PutRoom)
			r.Delete("/", h.DeleteRoom)
		})
	})

	return r
}
