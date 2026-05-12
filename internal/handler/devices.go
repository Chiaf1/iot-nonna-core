package handler

import (
	"log"
	"net/http"

	"github.com/chiaf1/iot-nonna-core/internal/domain"
	"github.com/go-chi/chi/v5"
)

// Handler get all devices
// GetDevices		godoc
// @Summary			Lists all devices
// @Tags			devices
// @Produce			json
// @Success			200	{array}		domain.Device
// @Failure			500	{string}	string	"internal server error"
// @Router			/devices	[get]
func (h *Handler) GetDevices(w http.ResponseWriter, r *http.Request) {
	// 1. Fetch data
	d, err := h.Repo.GetAllDevices()
	if err != nil {
		log.Printf("[handler][GetAllDevices] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJson(w, http.StatusOK, d)
}

// Handler get device by id in url
// GetDevice		godoc
// @Summary			Returns the device with specified id
// @Tags			devices
// @Produce			json
// @Param			id	path	string	true	"id"
// @Success			200	{object}	domain.Device
// @Failure			404	{string}	string	"device not found"
// @Failure			500	{string}	string	"internal server error"
// @Router			/devices/{id}	[get]
func (h *Handler) GetDevice(w http.ResponseWriter, r *http.Request) {
	// Retrieve id value from url
	id := chi.URLParam(r, "id")

	d, err := h.Repo.GetDeviceById(id)
	if err != nil {
		log.Printf("[handler][GetDeviceById] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if d == nil {
		http.Error(w, "device not found", http.StatusNotFound)
		return
	}
	writeJson(w, http.StatusOK, d)
}

// Handler post new device from json body
// PostDevice		godoc
// @Summary			Add a device
// @Tags			devices
// @Accept			json
// @Produce			json
// @Param			device	body	domain.DeviceRequest	true	"device payload"
// @Success			201	{object}	domain.Device
// @Failure			400		{object}	domain.ValidationErrorResponse	"validation error"
// @Failure			500	{string}	string	"internal server error"
// @Router			/devices	[post]
func (h *Handler) PostDevice(w http.ResponseWriter, r *http.Request) {
	// 1. Decode and validate JSON body
	req, err := decodeAndValidate[domain.DeviceRequest](r, h.Validator)
	if err != nil {
		h.respondValidationError(w, err)
		return
	}
	// 2. Call the repo
	d, err := h.Repo.CreateDevice(req)
	if err != nil {
		log.Printf("[handler][CreateDevice] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// 3. Respons with 201 created and the created room
	writeJson(w, http.StatusCreated, d)
}

// Handler to put device from json body
// PutDevice		godoc
// @Summary			Update device based on id
// @Tags			devices
// @Accept			json
// @Produce			json
// @Param			id	path	string	true	"device id"
// @Param			device	body	domain.DeviceRequest	true	"device payload"
// @Success			200	{object}	domain.Device
// @Failure			400		{object}	domain.ValidationErrorResponse	"validation error"
// @Failure			404	{string}	string	"device not found"
// @Failure			500	{string}	string	"internal server error"
// @Router			/devices/{id}	[put]
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
	d, err := h.Repo.UpdateDevice(id, req)
	if err != nil {
		log.Printf("[handler][UpdateDevice] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if d == nil {
		http.Error(w, "device not found", http.StatusNotFound)
		return
	}
	// 5. Respons with 201 created and the created room
	writeJson(w, http.StatusOK, d)
}

// Delete device based on id in url param
// DeleteDevice		godoc
// @Summary			Delete device
// @Tags			devices
// @Param			id	path	string	true	"device id"
// @Success			204		"no content"
// @Failure			404		{string}	string	"device not found"
// @Failure			500		{string}	string	"internal server error"
// @Router			/devices/{id}	[delete]
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
