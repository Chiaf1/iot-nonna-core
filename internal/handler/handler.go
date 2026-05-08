package handler

import "github.com/chiaf1/iot-nonna-core/internal/repository"

type Handler struct {
	Repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		Repo: repo,
	}
}
