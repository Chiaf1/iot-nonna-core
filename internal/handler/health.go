package handler

import "net/http"

// Handle get health request
// GetHealth	godoc
// @Summary		Check connection with db
// @Tags		health
// @Produce		json
// @Success		200	{object}	map[string]string	"ok"
// @Failure		500	{object}	map[string]string	"Db offline"
// @Router		/health	[get]
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if err := h.Repo.GetDbHealth(); err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{
			"status": "db offline",
		})
		return
	}
	writeJson(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
