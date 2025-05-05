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
		h.sendResponse(w, h.createErrorResponse("Failed to fetch free games", http.StatusInternalServerError))
		return
	}

	h.sendResponse(w, h.createSuccessResponse(games))
}

func (h *GamesHandler) createSuccessResponse(data *models.GamesData) *models.ApiResponse {
	return &models.ApiResponse{
		BaseResponse: models.BaseResponse{
			Success:   true,
			Timestamp: time.Now().Format(time.RFC3339),
		},
		Data: *data,
	}
}

func (h *GamesHandler) createErrorResponse(message string, code int) *models.ErrorResponse {
	return &models.ErrorResponse{
		BaseResponse: models.BaseResponse{
			Success:   false,
			Timestamp: time.Now().Format(time.RFC3339),
		},
		Error: models.ErrorDetails{
			Message: message,
			Code:    code,
		},
	}
}

func (h *GamesHandler) sendResponse(w http.ResponseWriter, response interface{}) {
	if errResp, ok := response.(*models.ErrorResponse); ok {
		w.WriteHeader(errResp.Error.Code)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(h.createErrorResponse("Failed to encode response", http.StatusInternalServerError))
	}
}
