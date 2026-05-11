package handler

import (
	"encoding/json"
	"net/http"

	"github.com/chiaf1/iot-nonna-core/internal/domain"
	"github.com/chiaf1/iot-nonna-core/internal/repository"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	Repo      *repository.Repository
	Validator *validator.Validate
}

func NewHandler(repo *repository.Repository) *Handler {
	validator := validator.New()
	validator.RegisterStructValidation(domain.SensorTypeReqValidation, domain.Sensor_typeRequest{})
	return &Handler{
		Repo:      repo,
		Validator: validator,
	}
}

// Helper function to write json http responses
func writeJson(w http.ResponseWriter, status int, data any) {
	// Set header to tell client what type of content to expect
	w.Header().Set("Content-Type", "application/json")
	// Set status code
	w.WriteHeader(status)
	// Encode and send the data
	json.NewEncoder(w).Encode(data)
}
