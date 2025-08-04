package models

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

// HealthResponse represents health check response
type HealthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}
