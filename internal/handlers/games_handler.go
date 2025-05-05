package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"free-games-epic/internal/models"
	"free-games-epic/internal/services"
)

type GamesHandler struct {
	epicService *services.EpicService
}

func NewGamesHandler(epicService *services.EpicService) *GamesHandler {
	return &GamesHandler{
		epicService: epicService,
	}
}

func (h *GamesHandler) GetFreeGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	games, err := h.epicService.GetFreeGames()
	if err != nil {
		h.sendError(w, "Failed to fetch free games", http.StatusInternalServerError)
		return
	}

	response := models.ApiResponse{
		Success:   true,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      *games,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *GamesHandler) sendError(w http.ResponseWriter, message string, code int) {
	response := models.ErrorResponse{
		Success:   false,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	response.Error.Message = message
	response.Error.Code = code

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
