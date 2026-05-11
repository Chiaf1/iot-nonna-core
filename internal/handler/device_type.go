package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/chiaf1/iot-nonna-core/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// Handler for get request on device_type
func (h *Handler) GetDeviceTypes(w http.ResponseWriter, r *http.Request) {
	// 1. Fetch data
	dt, err := h.Repo.GetAllDeviceType()
	if err != nil {
		log.Printf("[handler][GetAllDeviceType] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJson(w, http.StatusOK, dt)
}

// Handler for retriving device_type by id in url
func (h *Handler) GetDeviceType(w http.ResponseWriter, r *http.Request) {
	// Retrieve id value from url
	id := chi.URLParam(r, "id")

	dt, err := h.Repo.GetDeviceTypeById(id)
	if err != nil {
		log.Printf("[handler][GetDeviceTypeById] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if dt == nil {
		http.Error(w, "device_type not found", http.StatusNotFound)
		return
	}
	writeJson(w, http.StatusOK, dt)
}

// Handler to create a new device_type from json body
func (h *Handler) PostDeviceType(w http.ResponseWriter, r *http.Request) {
	// 1. Decode JSON body
	var req domain.Device_typeRequest
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
	st, err := h.Repo.CreateDeviceType(req)
	if err != nil {
		log.Printf("[handler][CreateDeviceType] %v", err)
		http.Error(w, "internal server erro", http.StatusInternalServerError)
		return
	}
	// 4. Respons with 201 created and the created room
	writeJson(w, http.StatusCreated, st)
}

// Handler to create a new device_type from json body
func (h *Handler) PutDeviceType(w http.ResponseWriter, r *http.Request) {
	// 1. Retrieve id value from url
	id := chi.URLParam(r, "id")
	// 2. Decode JSON body
	var req domain.Device_typeRequest
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
	dt, err := h.Repo.UpdateDeviceType(id, req)
	if err != nil {
		log.Printf("[handler][UpdateDeviceType] %v", err)
		http.Error(w, "internal server erro", http.StatusInternalServerError)
		return
	}
	if dt == nil {
		http.Error(w, "device_type not found", http.StatusNotFound)
		return
	}
	// 5. Respons with 201 created and the created room
	writeJson(w, http.StatusOK, dt)
}

// Delete sensorType based on id in url param
func (h *Handler) DeleteDeviceType(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	found, err := h.Repo.DeleteDeviceType(id)
	if err != nil {
		log.Printf("[handler][DeleteDeviceType] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "device_type not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 succes, no body
}
