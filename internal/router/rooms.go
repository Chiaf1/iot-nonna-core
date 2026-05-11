package router

import (
	"github.com/chiaf1/iot-nonna-core/internal/handler"
	"github.com/go-chi/chi/v5"
)

func routeRooms(r *chi.Mux, h *handler.Handler) {
	r.Route("/rooms", func(r chi.Router) {
		r.Get("/", h.GetRooms)
		r.Post("/", h.PostRoom)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetRoom)
			r.Put("/", h.PutRoom)
			r.Delete("/", h.DeleteRoom)
		})
	})
}
