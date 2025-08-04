package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"chat-backend/internal/models"
	"chat-backend/internal/websocket"

	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub *websocket.Hub
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// ServeHTTP handles WebSocket upgrade requests
func (h *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := websocket.NewClient(conn, h.hub)
	h.hub.Register() <- client
	client.Start()
}

// HealthHandler handles health check requests
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.HealthResponse{
		Status: "healthy",
		Time:   time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}
