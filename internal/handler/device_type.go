package handler

import (
	"log"
	"net/http"

	"github.com/chiaf1/iot-nonna-core/internal/domain"
	"github.com/go-chi/chi/v5"
)

// Handler for get request on device_type
// GetDeviceTypes		godoc
// @Summary				Lists all device types
// @Tags				deviceType
// @Produce				json
// @Success				200	{array}	domain.DeviceType
// @Failure				500	{string}	string	"internal server error"
// @Router				/device-types	[get]
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
// GetDeviceType		godoc
// @Summary				Returns the device type based on the id
// @Tags				deviceType
// @Produce				json
// @Param				id	path	string	true	"Device type id"
// @Success				200	{object}	domain.DeviceType
// @Failure				404	{string}	string	"device_type not found"
// @Failure				500	{string}	string	"internal server error"
// @Router				/device-types/{id}	[get]
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
// PostDeviceType		godoc
// @Summary				Add device type
// @Tags				deviceType
// @Accept				json
// @Produce				json
// @Param				deviceType	body	domain.DeviceTypeRequest	true	"Device type payload"
// @Success				201	{object}	domain.DeviceType
// @Failure				400	{object}	domain.ValidationErrorResponse	"validation error"
// @Failure				500	{string}	string	"internal server error"
// @Router				/device-types	[post]
func (h *Handler) PostDeviceType(w http.ResponseWriter, r *http.Request) {
	// 1. Decode and validate JSON body
	req, err := decodeAndValidate[domain.DeviceTypeRequest](r, h.Validator)
	if err != nil {
		h.respondValidationError(w, err)
		return
	}
	// 2. Call the repo
	dt, err := h.Repo.CreateDeviceType(req)
	if err != nil {
		log.Printf("[handler][CreateDeviceType] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// 3. Respons with 201 created and the created room
	writeJson(w, http.StatusCreated, dt)
}

// Handler to create a new device_type from json body
// PutDeviceType		godoc
// @Summary				Update device type
// @Tags				deviceType
// @Accept				json
// @Produce				json
// @Param				id	path	string	true	"Device type id"
// @Param				deviceType	body	domain.DeviceTypeRequest	true	"Device type payload"
// @Success				200	{object}	domain.DeviceType
// @Failure				400	{object}	domain.ValidationErrorResponse	"validation error"
// @Failure				404	{string}	string	"device_type not found"
// @Failure				500	{string}	string	"internal server error"
// @Router				/device-types/{id}	[put]
func (h *Handler) PutDeviceType(w http.ResponseWriter, r *http.Request) {
	// 1. Retrieve id value from url
	id := chi.URLParam(r, "id")
	// 2. Decode and validate JSON body
	req, err := decodeAndValidate[domain.DeviceTypeRequest](r, h.Validator)
	if err != nil {
		h.respondValidationError(w, err)
		return
	}
	// 4. Call the repo
	dt, err := h.Repo.UpdateDeviceType(id, req)
	if err != nil {
		log.Printf("[handler][UpdateDeviceType] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if dt == nil {
		http.Error(w, "device_type not found", http.StatusNotFound)
		return
	}
	// 5. Respons with 200 created and the created room
	writeJson(w, http.StatusOK, dt)
}

// Delete device_type based on id in url param
// DeleteDeviceType		godoc
// @Summary				Delete device type
// @Tags				deviceType
// @Param				id	path	string	true	"Device type id"
// @Success				204	"no content"
// @Failure				404	{string}	string	"device_type not found"
// @Failure				500	{string}	string	"internal server error"
// @Router				/device-types/{id}	[delete]
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
