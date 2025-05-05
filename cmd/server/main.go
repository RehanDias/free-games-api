package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"free-games-epic/internal/server"
	"free-games-epic/internal/services"
)

func main() {
	// Initialize services
	epicService := services.NewEpicService()

	// Setup router with handlers
	router := server.NewRouter(epicService)

	// Create server with configured routes
	srv := server.NewServer(router.Setup())

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		serverErrors <- srv.Start()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)

	case <-shutdown:
		log.Println("Starting shutdown...")

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown did not complete in %v : %v", 15*time.Second, err)
			os.Exit(1)
		}
	}
}
