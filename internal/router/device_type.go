package router

import (
	"github.com/chiaf1/iot-nonna-core/internal/handler"
	"github.com/go-chi/chi/v5"
)

func routeDeviceType(r *chi.Mux, h *handler.Handler) {
	r.Route("/device_type", func(r chi.Router) {
		r.Get("/", h.GetDeviceTypes)
		r.Post("/", h.PostDeviceType)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetDeviceType)
			r.Put("/", h.PutDeviceType)
			r.Delete("/", h.DeleteDeviceType)
		})
	})
}
