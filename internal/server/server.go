package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

// Server represents the API server configuration
type Server struct {
	server *http.Server
}

// NewServer creates a new server instance with the given handler
func NewServer(handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:         ":8080",
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// Start begins listening for requests
func (s *Server) Start() error {
	log.Printf("Server starting on http://localhost%s\n", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
