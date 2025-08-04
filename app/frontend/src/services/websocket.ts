import { Message, ClientInfo } from '../types';
import { getWebSocketUrl } from '../utils/config';

export class WebSocketService {
  private ws: WebSocket | null = null;
  private userId: string = '';
  private onMessageCallback?: (message: Message) => void;
  private onClientInfoCallback?: (info: ClientInfo) => void;
  private onConnectionStatusCallback?: (connected: boolean) => void;

  constructor() {
    this.generateUserId();
  }

  private generateUserId(): void {
    this.userId = `user_${Math.random().toString(36).substr(2, 9)}`;
  }

  getUserId(): string {
    return this.userId;
  }

  connect(): void {
    const wsUrl = getWebSocketUrl();
    console.log('Connecting to WebSocket:', wsUrl);

    this.ws = new WebSocket(wsUrl);

    this.ws.onopen = () => {
      console.log('Connected to WebSocket server');
      this.onConnectionStatusCallback?.(true);
      
      // Send user join message
      this.sendUserJoin();
    };

    this.ws.onclose = () => {
      console.log('Disconnected from WebSocket server');
      this.onConnectionStatusCallback?.(false);
    };

    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log('Received data:', data);
        
        if (data.type === 'message') {
          const message: Message = {
            content: data.content,
            senderId: data.senderId,
            timestamp: data.timestamp || new Date().toISOString()
          };
          this.onMessageCallback?.(message);
        } else if (data.type === 'client_info') {
          const clientInfo: ClientInfo = {
            totalClients: data.totalClients,
            onlineUsers: data.onlineUsers || []
          };
          this.onClientInfoCallback?.(clientInfo);
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
  }

  disconnect(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  sendMessage(content: string): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      const message = {
        type: 'message',
        content,
        senderId: this.userId,
        timestamp: new Date().toISOString()
      };
      this.ws.send(JSON.stringify(message));
    }
  }

  private sendUserJoin(): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      const userJoin = {
        type: 'user_join',
        userId: this.userId
      };
      this.ws.send(JSON.stringify(userJoin));
    }
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }

  onMessage(callback: (message: Message) => void): void {
    this.onMessageCallback = callback;
  }

  onClientInfo(callback: (info: ClientInfo) => void): void {
    this.onClientInfoCallback = callback;
  }

  onConnectionStatus(callback: (connected: boolean) => void): void {
    this.onConnectionStatusCallback = callback;
  }
}
