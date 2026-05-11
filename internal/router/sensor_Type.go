package router

import (
	"github.com/chiaf1/iot-nonna-core/internal/handler"
	"github.com/go-chi/chi/v5"
)

func routeSensorType(r *chi.Mux, h *handler.Handler) {
	r.Route("/sensor_type", func(r chi.Router) {
		r.Get("/", h.GetSensorTypes)
		r.Post("/", h.PostSensorType)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetSensorType)
			r.Put("/", h.PutSensorType)
			r.Delete("/", h.DeleteSensorType)
		})
	})
}
