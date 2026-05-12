package handler

import (
	"log"
	"net/http"

	"github.com/chiaf1/iot-nonna-core/internal/domain"
	"github.com/go-chi/chi/v5"
)

// Get all sensors associated with device
// GetDeviceSensors			godoc
// @Summary					Lists all sensors of selected device
// @Tags					deviceSensors
// @Produce					json
// @Param					id	path	string	true	"device id"
// @Success					200	{array}	domain.SensorType
// @Failure					500	{string}	string	"internal server error"
// @Router					/devices/{id}/sensors	[get]
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
// PostDeviceSensors			godoc
// @Summary					Add sensor to device
// @Tags					deviceSensors
// @Accept					json
// @Produce					json
// @Param					id	path	string	true	"device id"
// @Param					sensorId	body	domain.AssociateSensorRequest	true	"Associated sensor payload"
// @Success					204	"status no content"
// @Failure					400	{object}	domain.ValidationErrorResponse	"validation error"
// @Failure					500	{string}	string	"internal server error"
// @Router					/devices/{id}/sensors	[post]
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

// DeleteDeviceSensor			godoc
// @Summary					Delete sensor from device
// @Tags					deviceSensors
// @Param					id	path	string	true	"device id"
// @Param					sensorId	path	string	true	"sensor id"
// @Success					204	"status no content"
// @Failure					404	{string}	string	"sensors_devices not found"
// @Failure					500	{string}	string	"internal server error"
// @Router					/devices/{id}/sensors/{sensorId}	[delete]
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
