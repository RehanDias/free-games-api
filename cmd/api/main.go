package main

import (
	"log"
	"net/http"

	"free-games-epic/internal/handlers"
	"free-games-epic/internal/services"
)

func main() {
	epicService := services.NewEpicService()
	gamesHandler := handlers.NewGamesHandler(epicService)

	http.HandleFunc("/api/free-games", gamesHandler.GetFreeGames)

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
