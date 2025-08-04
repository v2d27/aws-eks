// Message type definition
export interface Message {
  content: string;
  senderId?: string;
  timestamp?: string;
}

// Client info type
export interface ClientInfo {
  totalClients: number;
  onlineUsers: string[];
}
