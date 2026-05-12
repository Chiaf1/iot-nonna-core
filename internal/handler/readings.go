package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Get readings for dht sensor of a set range
// GetDhtReadings		godoc
// @Summary				Lists all dht readings of deiveceId in range from - to and limit number of lines
// @Tags				readings
// @Produce				json
// @Param				deviceId	path 	string	true	"device id"
// @Param				from		query 	string	false	"time stamp start query (default 24h ago)"
// @Param				to			query 	string	false	"time stamp stop query (default now)"
// @Param				limit		query 	int		false	"max number of lines (default 500)"
// @Success				200	{array}	domain.DhtReadings
// @Failure				400	{string}	string	"input parameters not valid"
// @Failure				500	{string}	string	"internal server error"
// @Router				/readings/dht/{deviceId}	[get]
func (h *Handler) GetDhtReadings(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "deviceId")

	now := time.Now()
	from, err := parseTimeParam(r, "from", now.Add(-24*time.Hour))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	to, err := parseTimeParam(r, "to", now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit, err := parseIntParam(r, "limit", 500)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	readings, err := h.Repo.GetDhtReadings(deviceId, from, to, limit)
	if err != nil {
		log.Printf("[handler][GetDhtReadings] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJson(w, http.StatusOK, readings)
}

// Get latest readings for dht sensor
// GetDhtReadingLatest		godoc
// @Summary				Returns latest dht reading
// @Tags				readings
// @Produce				json
// @Param				deviceId	path 	string	true	"device id"
// @Success				200	{object}	domain.DhtReadings
// @Failure				404	{string}	string	"latest dht reading not found"
// @Failure				500	{string}	string	"internal server error"
// @Router				/readings/dht/{deviceId}/latest	[get]
func (h *Handler) GetDhtReadingLatest(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "deviceId")

	readings, err := h.Repo.GetLatestDhtReading(deviceId)
	if err != nil {
		log.Printf("[handler][GetLatestDhtReading] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if readings == nil {
		http.Error(w, "latest dht reading not found", http.StatusNotFound)
		return
	}
	writeJson(w, http.StatusOK, readings)
}

// Get readings for status sensor of a set range
// GetStatusReadings		godoc
// @Summary				Lists all status readings of deiveceId in range from - to and limit number of lines
// @Tags				readings
// @Produce				json
// @Param				deviceId	path 	string	true	"device id"
// @Param				from		query 	string	false	"time stamp start query (default 24h ago)"
// @Param				to			query 	string	false	"time stamp stop query (default now)"
// @Param				limit		query 	int		false	"max number of lines (default 500)"
// @Success				200	{array}	domain.StatusReadings
// @Failure				400	{string}	string	"input parameters not valid"
// @Failure				500	{string}	string	"internal server error"
// @Router				/readings/status/{deviceId}	[get]
func (h *Handler) GetStatusReadings(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "deviceId")

	now := time.Now()
	from, err := parseTimeParam(r, "from", now.Add(-24*time.Hour))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	to, err := parseTimeParam(r, "to", now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit, err := parseIntParam(r, "limit", 500)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	readings, err := h.Repo.GetStatusReadings(deviceId, from, to, limit)
	if err != nil {
		log.Printf("[handler][GetStatusReadings] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJson(w, http.StatusOK, readings)
}

// Get latest readings for status sensor
// GetStatusReadingLatest		godoc
// @Summary				Returns latest status reading
// @Tags				readings
// @Produce				json
// @Param				deviceId	path 	string	true	"device id"
// @Success				200	{object}	domain.StatusReadings
// @Failure				404	{string}	string	"latest status reading not found"
// @Failure				500	{string}	string	"internal server error"
// @Router				/readings/status/{deviceId}/latest	[get]
func (h *Handler) GetStatusReadingLatest(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "deviceId")

	readings, err := h.Repo.GetLatestStatusReading(deviceId)
	if err != nil {
		log.Printf("[handler][GetLatestStatusReading] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if readings == nil {
		http.Error(w, "latest status reading not found", http.StatusNotFound)
		return
	}
	writeJson(w, http.StatusOK, readings)
}
