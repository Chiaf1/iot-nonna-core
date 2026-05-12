package router

import (
	"github.com/chiaf1/iot-nonna-core/internal/handler"
	"github.com/go-chi/chi/v5"
)

func routeReadings(r *chi.Mux, h *handler.Handler) {
	r.Route("/readings", func(r chi.Router) {
		r.Route("/dht/{deviceId}", func(r chi.Router) {
			r.Get("/", h.GetDhtReadings)
			r.Get("/latest", h.GetDhtReadingLatest)
		})
		r.Route("/status/{deviceId}", func(r chi.Router) {
			r.Get("/", h.GetStatusReadings)
			r.Get("/latest", h.GetStatusReadingLatest)
		})
	})
}
