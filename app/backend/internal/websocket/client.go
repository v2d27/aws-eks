package websocket

import (
	"encoding/json"
	"log"
	"time"

	"chat-backend/internal/models"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client
type Client struct {
	conn   *websocket.Conn
	send   chan interface{}
	userId string
	hub    *Hub
}

// NewClient creates a new WebSocket client
func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		conn: conn,
		send: make(chan interface{}, 256),
		hub:  hub,
	}
}

// Start starts the client's read and write pumps
func (c *Client) Start() {
	go c.writePump()
	go c.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var rawMessage json.RawMessage
		err := c.conn.ReadJSON(&rawMessage)
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
			var userJoin models.UserJoin
			if err := json.Unmarshal(rawMessage, &userJoin); err != nil {
				log.Printf("Error parsing user join: %v", err)
				continue
			}
			c.userId = userJoin.UserId
			log.Printf("User joined: %s", c.userId)

		case "message":
			var message models.Message
			if err := json.Unmarshal(rawMessage, &message); err != nil {
				log.Printf("Error parsing message: %v", err)
				continue
			}

			if message.Timestamp == "" {
				message.Timestamp = time.Now().Format(time.RFC3339)
			}

			log.Printf("Message received from %s: %s", message.SenderId, message.Content)
			c.hub.broadcast <- message
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		if err := c.conn.WriteJSON(message); err != nil {
			log.Printf("WebSocket write error: %v", err)
			return
		}
	}
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}
