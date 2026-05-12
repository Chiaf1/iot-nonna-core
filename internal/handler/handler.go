package handler

import (
	"encoding/json"
	"errors"
	"fmt"
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
	validator.RegisterStructValidation(domain.SensorTypeReqValidation, domain.SensorTypeRequest{})
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

// Helper function to decode and validate JSON requests
func decodeAndValidate[T any](r *http.Request, v *validator.Validate) (T, error) {
	var req T
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		return req, fmt.Errorf("invalid body: %w", err)
	}
	if err := v.Struct(req); err != nil {
		return req, err
	}
	return req, nil
}

// Helper to organize validation errors
func (h *Handler) respondValidationError(w http.ResponseWriter, err error) {
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		errs := make(map[string]string)
		for _, v := range verrs {
			errs[v.Field()] = v.Tag()
		}
		writeJson(w, http.StatusBadRequest, domain.ValidationErrorResponse{
			Errors: errs,
		})
		return
	}
	http.Error(w, err.Error(), http.StatusBadRequest)
}
