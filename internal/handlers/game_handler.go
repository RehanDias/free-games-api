package handlers

import (
	"net/http"

	"epic-games-free/internal/services"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	service *services.GameService
}

func NewGameHandler() *GameHandler {
	return &GameHandler{
		service: services.NewGameService(),
	}
}

func (h *GameHandler) GetFreeGames(c *gin.Context) {
	response, err := h.service.GetFreeGames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"message": err.Error(),
				"code":    http.StatusInternalServerError,
			},
		})
		return
	}

	if !response.Success {
		c.JSON(response.Error.Code, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
