package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// Handler for retriving the list of all rooms
func (h *Handler) GetRooms(w http.ResponseWriter, r *http.Request) {
	// 1. Fetch data
	rooms, err := h.Repo.GetAllRooms()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Printf("[handler][rooms] GetAllRooms error: %v", err)
		return
	}

	// 2. Set header to tell client what type of content to expect
	w.Header().Set("Content-Type", "application/json")

	// 3. Set status code
	w.WriteHeader(http.StatusOK)

	// 4. Encode and send the data
	json.NewEncoder(w).Encode(rooms)
}
