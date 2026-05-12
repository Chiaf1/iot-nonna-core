package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/chiaf1/iot-nonna-core/internal/domain"
	"github.com/go-chi/chi/v5"
)

// Handler for retriving the list of all rooms
// GetRooms 	godoc
// @summary		lists all rooms
// @Tags		rooms
// @Produce		json
// @Success		200 {array}		domain.Room
// @Failure		500	{string}	string	"internal server error"
// @Router		/rooms	[get]
func (h *Handler) GetRooms(w http.ResponseWriter, r *http.Request) {
	// 1. Fetch data
	rooms, err := h.Repo.GetAllRooms()
	if err != nil {
		log.Printf("[handler][GetRooms] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJson(w, http.StatusOK, rooms)
}

// Handler for retriving room by id in url
// GetRoom		godoc
// @Summary		returns the rooms with the specified id
// @Tags		rooms
// @Produce		json
// @Param		id		path	string	true	"Room ID"
// @Success		200		{object}	domain.Room
// @Failure		404		{string}	string	"room not found"
// @Failure		500		{string}	string	"internal server error"
// @Router		/rooms/{id}	[get]
func (h *Handler) GetRoom(w http.ResponseWriter, r *http.Request) {
	// Retrieve id value from url
	id := chi.URLParam(r, "id")

	room, err := h.Repo.GetRoomById(id)
	if err != nil {
		log.Printf("[handler][GetRoom] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if room == nil {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}
	writeJson(w, http.StatusOK, room)
}

// Handler to create a new room from json body
// PostRoom		godoc
// @Summary		Create a new room
// @Tags		rooms
// @Accept		json
// @Produce		json
// @Param		room	body	domain.RoomRequest	true	"Room payload"
// @Success		201		{object}	domain.Room
// @Failure		400		{string}	string	"invalid request body"
// @Failure		500		{string}	string	"internal server error"
// @Router		/rooms	[post]
func (h *Handler) PostRoom(w http.ResponseWriter, r *http.Request) {
	// 1. Decode JSON body
	var req domain.RoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	// 2. Validate base
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	// 3. Call the repo
	room, err := h.Repo.CreateRoom(req.Name)
	if err != nil {
		log.Printf("[handler][CreateRoom] %v", err)
		http.Error(w, "internal server erro", http.StatusInternalServerError)
		return
	}
	// 4. Respons with 201 created and the created room
	writeJson(w, http.StatusCreated, room)
}

// Update room value based on the id as url param and new name in JSON body
// PutRoom		godoc
// @Summary		Update room
// @Tags		rooms
// @Accept		json
// @Produce		json
// @Param		id		path	string				true	"RoomId"
// @Param		room	body	domain.RoomRequest	true	"Room payload"
// @Success		200		{object}	domain.Room
// @Failure		400		{string}	string	"invalid request body"
// @Failure		404		{string}	string	"room not found"
// @Failure		500		{string}	string	"internal server error"
// @Router		/rooms/{id}	[put]
func (h *Handler) PutRoom(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req domain.RoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body request", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	room, err := h.Repo.UpdateRoom(id, req.Name)
	if err != nil {
		log.Printf("[handler][UpdateRoom] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if room == nil {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}
	writeJson(w, http.StatusOK, room)
}

// Delete room based on id in url param
// DeleteRoom		godoc
// @Summary			Delete the room
// @Tags			rooms
// @Param			id		path	string	true	"Room ID"
// @Success			204		"no content"
// @Failure			404		{string}	string	"room not found"
// @Failure			500		{string}	string	"internal server error"
// @Router			/rooms/{id}	[delete]
func (h *Handler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	found, err := h.Repo.DeleteRoom(id)
	if err != nil {
		log.Printf("[handler][DeleteRoom] %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 succes, no body
}
