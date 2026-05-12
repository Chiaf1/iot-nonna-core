package handler

import (
	"log"
	"net/http"

	"github.com/chiaf1/iot-nonna-core/internal/domain"
	"github.com/go-chi/chi/v5"
)

// Get all sensors associated with device
func (h *Handler) GetDeviceSensors(w http.ResponseWriter, r *http.Request) {
	// Retrieve id value from url
	id := chi.URLParam(r, "id")

	ds, err := h.Repo.GetDevicesSensors(id)
	if err != nil {
		log.Printf("[handler][GetDevicesSensors] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJson(w, http.StatusOK, ds)
}

// Post sensor associated to a device
func (h *Handler) PostDeviceSensors(w http.ResponseWriter, r *http.Request) {
	// Retrieve id value from url
	deviceId := chi.URLParam(r, "id")
	// 1. Decode and validate JSON body
	req, err := decodeAndValidate[domain.AssociateSensorRequest](r, h.Validator)
	if err != nil {
		h.respondValidationError(w, err)
		return
	}
	// 2. Call the repo
	err = h.Repo.AddSensorToDevice(deviceId, req.SensorId)
	if err != nil {
		log.Printf("[handler][AddSensorToDevice] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 succes, no body
}

func (h *Handler) DeleteDeviceSensor(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "id")
	sensorId := chi.URLParam(r, "sensorId")

	found, err := h.Repo.RemoveSensorFromDevice(deviceId, sensorId)
	if err != nil {
		log.Printf("[handler][RemoveSensorFromDevice] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "sensors_devices not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 succes, no body
}
