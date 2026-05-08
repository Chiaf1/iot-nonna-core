package handler

import (
	"encoding/json"
	"net/http"
)

// Handler for retriving the list of all rooms
func (h *Handler) HandleRooms(w http.ResponseWriter, r *http.Request) {
	// 1. Fetch data
	rooms, _ := h.Repo.GetAllRooms()

	// 2. Set header to tell client what type of content to expect
	w.Header().Set("Content-Type", "application/json")

	// 3. Set status code
	w.WriteHeader(http.StatusOK)

	// 4. Encode and send the data
	json.NewEncoder(w).Encode(rooms)
}
