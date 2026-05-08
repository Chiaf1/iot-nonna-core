package handler

import "net/http"

// Handle get health request
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if err := h.Repo.GetDbHealth(); err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Db offline"))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
