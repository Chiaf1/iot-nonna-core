package handler

import (
	"log"
	"net/http"

	"github.com/chiaf1/iot-nonna-core/internal/domain"
	"github.com/go-chi/chi/v5"
)

// Handler for retriving the list of all sensor type
// GetSensorTypes	godoc
// @Summary			Lists all sensor types
// @Tags			sensorTypes
// @Produce			json
// @Success			200		{array}		domain.SensorType
// @Failure			500		{string}	string	"internal server error"
// @Router			/sensor-types	[get]
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
// GetSensorType	godoc
// @Summary			Returns sensor type with specified id
// @Tags			sensorTypes
// @Produce			json
// @Param			id	path	string	true	"SensorType ID"
// @Success			200	{object}	domain.SensorType
// @Failure			404	{string}	string	"sensor_type not found"
// @Failure			500	{string}	string	"internal server error"
// @Router			/sensor-types/{id}	[get]
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
// PostSensorType	godoc
// @Summary			Add new sensor type
// @Tags			sensorTypes
// @Accept			json
// @Produce			json
// @Param			sensorType body	domain.SensorTypeRequest true	"Sensor type payload"
// @Success			201		{object}	domain.SensorType
// @Failure			400		{object}	domain.ValidationErrorResponse	"validation error"
// @Failure			500		{string}	string	"internal server error"
// @Router			/sensor-types	[post]
func (h *Handler) PostSensorType(w http.ResponseWriter, r *http.Request) {
	// 1. Decode and validate JSON body
	req, err := decodeAndValidate[domain.SensorTypeRequest](r, h.Validator)
	if err != nil {
		h.respondValidationError(w, err)
		return
	}
	// 2. Call the repo
	st, err := h.Repo.CreateSensorType(req)
	if err != nil {
		log.Printf("[handler][CreateSensorType] %v", err)
		http.Error(w, "internal server erro", http.StatusInternalServerError)
		return
	}
	// 3. Respons with 201 created and the created room
	writeJson(w, http.StatusCreated, st)
}

// Handler to create a new sensor type from json body
// PutSensorType	godoc
// @Summary			Update sensor type based on id
// @Tags			sensorTypes
// @Accept			json
// @Produce			json
// @Param			id		path	string	true	"sensor type id"
// @Param			sensorType	body	domain.SensorTypeRequest	true	"sensor type payload"
// @Success			200		{object}	domain.SensorType
// @Failure			400		{object}	domain.ValidationErrorResponse	"validation error"
// @Failure			404		{string}	string	true	"sensor_type not found"
// @Failure			500		{string}	string	true	"internal server error"
// @Router			/sensor-types/{id}	[put]
func (h *Handler) PutSensorType(w http.ResponseWriter, r *http.Request) {
	// 1. Retrieve id value from url
	id := chi.URLParam(r, "id")
	// 2. Decode and validate JSON body
	req, err := decodeAndValidate[domain.SensorTypeRequest](r, h.Validator)
	if err != nil {
		h.respondValidationError(w, err)
		return
	}
	// 3. Call the repo
	st, err := h.Repo.UpdateSensorType(id, req)
	if err != nil {
		log.Printf("[handler][UpdateSensorType] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if st == nil {
		http.Error(w, "sensor_type not found", http.StatusNotFound)
		return
	}
	// 4. Respons with 200 ok
	writeJson(w, http.StatusOK, st)
}

// Delete sensorType based on id in url param
// DeleteSensorType	godoc
// @Summary			Delete sensor type based on id
// @Tags			sensorTypes
// @Param			id	path	string	true	"Sensor type id"
// @Success			204	"no content"
// @Failure			404	{string}	string	"Sensor type not found"
// @Failure			500	{string}	string	"internal server error"
// @Router			/sensor-types/{id}	[delete]
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
