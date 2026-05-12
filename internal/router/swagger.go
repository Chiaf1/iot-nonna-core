package router

import (
	_ "github.com/chiaf1/iot-nonna-core/docs"
	"github.com/chiaf1/iot-nonna-core/internal/handler"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func routeSwagger(r *chi.Mux, h *handler.Handler) {
	r.Get("/swagger/*", httpSwagger.WrapHandler)
}
