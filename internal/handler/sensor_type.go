package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/chiaf1/iot-nonna-core/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// Handler for retriving the list of all sensor type
func (h *Handler) GetSensorTypes(w http.ResponseWriter, r *http.Request) {
	// 1. Fetch data
	st, err := h.Repo.GetAllSensorType()
	if err != nil {
		log.Printf("[handler][GetAllSensorType] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJson(w, http.StatusOK, st)
}

// Handler for retriving sensor type by id in url
func (h *Handler) GetSensorType(w http.ResponseWriter, r *http.Request) {
	// Retrieve id value from url
	id := chi.URLParam(r, "id")

	st, err := h.Repo.GetSensorTypeById(id)
	if err != nil {
		log.Printf("[handler][GetSensorType] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if st == nil {
		http.Error(w, "sensor_type not found", http.StatusNotFound)
		return
	}
	writeJson(w, http.StatusOK, st)
}

// Handler to create a new sensor type from json body
func (h *Handler) PostSensorType(w http.ResponseWriter, r *http.Request) {
	// 1. Decode JSON body
	var req domain.Sensor_typeRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	// 2. Validate values
	if err := h.Validator.Struct(req); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, v := range verrs {
				errors[v.Field()] = v.Tag()
			}
			writeJson(w, http.StatusBadRequest, map[string]any{
				"errors": errors,
			})
		}
		return
	}
	// 3. Call the repo
	st, err := h.Repo.CreateSensorType(req)
	if err != nil {
		log.Printf("[handler][CreateSensorType] %v", err)
		http.Error(w, "internal server erro", http.StatusInternalServerError)
		return
	}
	// 4. Respons with 201 created and the created room
	writeJson(w, http.StatusCreated, st)
}

// Handler to create a new sensor type from json body
func (h *Handler) PutSensorType(w http.ResponseWriter, r *http.Request) {
	// 1. Retrieve id value from url
	id := chi.URLParam(r, "id")
	// 2. Decode JSON body
	var req domain.Sensor_typeRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	// 3. Validate values
	if err := h.Validator.Struct(req); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, v := range verrs {
				errors[v.Field()] = v.Tag()
			}
			writeJson(w, http.StatusBadRequest, map[string]any{
				"errors": errors,
			})
		}
		return
	}
	// 4. Call the repo
	st, err := h.Repo.UpdateSensorType(id, req)
	if err != nil {
		log.Printf("[handler][UpdateSensorType] %v", err)
		http.Error(w, "internal server erro", http.StatusInternalServerError)
		return
	}
	if st == nil {
		http.Error(w, "sensor_type not found", http.StatusNotFound)
		return
	}
	// 5. Respons with 201 created and the created room
	writeJson(w, http.StatusOK, st)
}

// Delete sensorType based on id in url param
func (h *Handler) DeleteSensorType(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	found, err := h.Repo.DeleteSensorType(id)
	if err != nil {
		log.Printf("[handler][DeleteSensorType] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "sensor_type not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 succes, no body
}
