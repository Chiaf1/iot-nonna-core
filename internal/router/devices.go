package router

import (
	"github.com/chiaf1/iot-nonna-core/internal/handler"
	"github.com/go-chi/chi/v5"
)

func routeDevices(r *chi.Mux, h *handler.Handler) {
	r.Route("/devices", func(r chi.Router) {
		r.Get("/", h.GetDevices)
		r.Post("/", h.PostDevice)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetDevice)
			r.Put("/", h.PutDevice)
			r.Delete("/", h.DeleteDevice)
		})
	})
}
