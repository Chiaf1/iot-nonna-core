package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Get readings for dht sensor of a set range
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
