package handler

import (
	"encoding/json"
	"net/http"

	"github.com/chiaf1/iot-nonna-core/internal/repository"
)

type Handler struct {
	Repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		Repo: repo,
	}
}

// Helper function to write json http responses
func writeJson(w http.ResponseWriter, status int, data any) {
	// Set header to tell client what type of content to expect
	w.Header().Set("Content-Type", "application/json")
	// Set status code
	w.WriteHeader(status)
	// Encode and send the data
	json.NewEncoder(w).Encode(data)
}
