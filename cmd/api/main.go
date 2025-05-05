package api

import (
	"net/http"

	"epic-games-free/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	app = gin.New()
	
	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	app.Use(cors.New(config))
	
	handler := handlers.NewGameHandler()
	app.GET("/api/free-games", handler.GetFreeGames)
}

// Handler - Vercel serverless function entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
