package handler

import (
	"log"
	"net/http"

	"github.com/chiaf1/iot-nonna-core/internal/domain"
	"github.com/go-chi/chi/v5"
)

// Handler get all devices
func (h *Handler) GetDevices(w http.ResponseWriter, r *http.Request) {
	// 1. Fetch data
	dt, err := h.Repo.GetAllDevices()
	if err != nil {
		log.Printf("[handler][GetAllDevices] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJson(w, http.StatusOK, dt)
}

// Handler get device by id in url
func (h *Handler) GetDevice(w http.ResponseWriter, r *http.Request) {
	// Retrieve id value from url
	id := chi.URLParam(r, "id")

	dt, err := h.Repo.GetDeviceById(id)
	if err != nil {
		log.Printf("[handler][GetDeviceById] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if dt == nil {
		http.Error(w, "device not found", http.StatusNotFound)
		return
	}
	writeJson(w, http.StatusOK, dt)
}

// Handler post new device from json body
func (h *Handler) PostDevice(w http.ResponseWriter, r *http.Request) {
	// 1. Decode and validate JSON body
	req, err := decodeAndValidate[domain.DeviceRequest](r, h.Validator)
	if err != nil {
		h.respondValidationError(w, err)
		return
	}
	// 2. Call the repo
	st, err := h.Repo.CreateDevice(req)
	if err != nil {
		log.Printf("[handler][CreateDevice] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// 3. Respons with 201 created and the created room
	writeJson(w, http.StatusCreated, st)
}

// Handler to put device from json body
func (h *Handler) PutDevice(w http.ResponseWriter, r *http.Request) {
	// 1. Retrieve id value from url
	id := chi.URLParam(r, "id")
	// 2. Decode and validate JSON body
	req, err := decodeAndValidate[domain.DeviceRequest](r, h.Validator)
	if err != nil {
		h.respondValidationError(w, err)
		return
	}
	// 4. Call the repo
	dt, err := h.Repo.UpdateDevice(id, req)
	if err != nil {
		log.Printf("[handler][UpdateDevice] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if dt == nil {
		http.Error(w, "device not found", http.StatusNotFound)
		return
	}
	// 5. Respons with 201 created and the created room
	writeJson(w, http.StatusOK, dt)
}

// Delete device based on id in url param
func (h *Handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	found, err := h.Repo.DeleteDevice(id)
	if err != nil {
		log.Printf("[handler][DeleteDevice] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "device_type not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 succes, no body
}
