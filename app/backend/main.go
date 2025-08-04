package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Message represents a chat message structure
type Message struct {
	Type      string `json:"type"`
	Content   string `json:"content,omitempty"`
	SenderId  string `json:"senderId,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// ClientInfo represents client information
type ClientInfo struct {
	Type         string   `json:"type"`
	TotalClients int      `json:"totalClients"`
	OnlineUsers  []string `json:"onlineUsers"`
}

// UserJoin represents user join message
type UserJoin struct {
	Type   string `json:"type"`
	UserId string `json:"userId"`
}

// Client represents a WebSocket client
type Client struct {
	conn   *websocket.Conn
	send   chan interface{}
	userId string
}

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	clients    map[*Client]bool
	userIds    map[string]*Client
	broadcast  chan interface{}
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		userIds:    make(map[string]*Client),
		broadcast:  make(chan interface{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			if client.userId != "" {
				h.userIds[client.userId] = client
			}
			clientCount := len(h.clients)
			onlineUsers := make([]string, 0, len(h.userIds))
			for userId := range h.userIds {
				onlineUsers = append(onlineUsers, userId)
			}
			h.mutex.Unlock()

			log.Printf("New client connected: %s. Total clients: %d", client.userId, clientCount)

			// Broadcast updated client info to all clients
			clientInfo := ClientInfo{
				Type:         "client_info",
				TotalClients: clientCount,
				OnlineUsers:  onlineUsers,
			}
			h.broadcastToAll(clientInfo)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if client.userId != "" {
					delete(h.userIds, client.userId)
				}
				close(client.send)
				clientCount := len(h.clients)
				onlineUsers := make([]string, 0, len(h.userIds))
				for userId := range h.userIds {
					onlineUsers = append(onlineUsers, userId)
				}
				h.mutex.Unlock()

				log.Printf("Client disconnected: %s. Total clients: %d", client.userId, clientCount)

				// Broadcast updated client info to all clients
				clientInfo := ClientInfo{
					Type:         "client_info",
					TotalClients: clientCount,
					OnlineUsers:  onlineUsers,
				}
				h.broadcastToAll(clientInfo)
			} else {
				h.mutex.Unlock()
			}

		case message := <-h.broadcast:
			h.broadcastToAll(message)
		}
	}
}

func (h *Hub) broadcastToAll(message interface{}) {
	h.mutex.RLock()
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
			if client.userId != "" {
				delete(h.userIds, client.userId)
			}
		}
	}
	h.mutex.RUnlock()
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

func (h *Hub) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan interface{}, 256),
	}

	h.register <- client

	// Start goroutines for reading and writing
	go h.writePump(client)
	go h.readPump(client)
}

func (h *Hub) readPump(client *Client) {
	defer func() {
		h.unregister <- client
		client.conn.Close()
	}()

	for {
		var rawMessage json.RawMessage
		err := client.conn.ReadJSON(&rawMessage)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse the message to determine its type
		var baseMessage struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(rawMessage, &baseMessage); err != nil {
			log.Printf("Error parsing message type: %v", err)
			continue
		}

		switch baseMessage.Type {
		case "user_join":
			var userJoin UserJoin
			if err := json.Unmarshal(rawMessage, &userJoin); err != nil {
				log.Printf("Error parsing user join: %v", err)
				continue
			}
			client.userId = userJoin.UserId
			log.Printf("User joined: %s", client.userId)

		case "message":
			var message Message
			if err := json.Unmarshal(rawMessage, &message); err != nil {
				log.Printf("Error parsing message: %v", err)
				continue
			}

			if message.Timestamp == "" {
				message.Timestamp = time.Now().Format(time.RFC3339)
			}

			log.Printf("Message received from %s: %s", message.SenderId, message.Content)
			h.broadcast <- message
		}
	}
}

func (h *Hub) writePump(client *Client) {
	defer client.conn.Close()

	for message := range client.send {
		if err := client.conn.WriteJSON(message); err != nil {
			log.Printf("WebSocket write error: %v", err)
			return
		}
	}
	client.conn.WriteMessage(websocket.CloseMessage, []byte{})
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Load .env file for local development
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading .env file: %v", err)
	}

	hub := newHub()
	go hub.run()

	// Get host and port from environment variables
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get allowed origins from environment
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}

	// Create CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins: []string{allowedOrigin},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// Create HTTP mux
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.handleWebSocket)
	mux.HandleFunc("/healthz", healthCheckHandler)

	// Wrap with CORS
	handler := c.Handler(mux)

	address := host + ":" + port
	log.Printf("Chat server starting on %s", address)
	log.Printf("WebSocket endpoint: ws://%s/ws", address)
	log.Printf("Health check endpoint: http://%s/healthz", address)

	if err := http.ListenAndServe(address, handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
