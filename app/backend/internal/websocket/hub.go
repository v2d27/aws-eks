package websocket

import (
	"log"
	"sync"

	"chat-backend/internal/models"
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	clients    map[*Client]bool
	userIds    map[string]*Client
	broadcast  chan interface{}
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		userIds:    make(map[string]*Client),
		broadcast:  make(chan interface{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
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
			clientInfo := models.ClientInfo{
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
				clientInfo := models.ClientInfo{
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

// Register registers a new client
func (h *Hub) Register() chan<- *Client {
	return h.register
}

// Unregister returns the unregister channel
func (h *Hub) Unregister() chan<- *Client {
	return h.unregister
}

// Broadcast broadcasts a message to all clients
func (h *Hub) Broadcast() chan<- interface{} {
	return h.broadcast
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
