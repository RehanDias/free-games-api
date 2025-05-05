package server

import (
	"net/http"

	"free-games-epic/internal/handlers"
	"free-games-epic/internal/services"
)

// Router handles all route configurations
type Router struct {
	gamesHandler *handlers.GamesHandler
}

// NewRouter creates a new router instance
func NewRouter(epicService *services.EpicService) *Router {
	return &Router{
		gamesHandler: handlers.NewGamesHandler(epicService),
	}
}

// Setup configures all the routes
func (r *Router) Setup() http.Handler {
	mux := http.NewServeMux()

	// Apply middleware to all routes
	handler := r.applyMiddleware(mux)

	// Configure routes
	mux.HandleFunc("/api/free-games", r.gamesHandler.GetFreeGames)

	return handler
}

// applyMiddleware adds common middleware to all routes
func (r *Router) applyMiddleware(handler http.Handler) http.Handler {
	return EnableCORS(EnableLogging(handler))
}
