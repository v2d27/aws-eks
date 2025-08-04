package main

import (
	"log"
	"net/http"

	"chat-backend/internal/handlers"
	"chat-backend/internal/websocket"
	"chat-backend/pkg/config"

	"github.com/rs/cors"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Create handlers
	wsHandler := handlers.NewWebSocketHandler(hub)

	// Create CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins: []string{cfg.AllowedOrigin},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// Create HTTP mux
	mux := http.NewServeMux()
	mux.Handle("/ws", wsHandler)
	mux.HandleFunc("/healthz", handlers.HealthHandler)

	// Wrap with CORS
	handler := c.Handler(mux)

	// Log startup information
	log.Printf("Chat server starting on %s", cfg.Address())
	log.Printf("WebSocket endpoint: ws://%s/ws", cfg.Address())
	log.Printf("Health check endpoint: http://%s/healthz", cfg.Address())
	log.Printf("Environment: %s", cfg.Environment)
	log.Printf("Allowed Origin: %s", cfg.AllowedOrigin)

	// Start server
	if err := http.ListenAndServe(cfg.Address(), handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
