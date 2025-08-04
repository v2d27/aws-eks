import { useState, useEffect, useRef } from 'react';
import { Message, ClientInfo } from '../types';
import { WebSocketService } from '../services/websocket';

export const useWebSocket = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [isConnected, setIsConnected] = useState<boolean>(false);
  const [clientInfo, setClientInfo] = useState<ClientInfo>({ totalClients: 0, onlineUsers: [] });
  const wsServiceRef = useRef<WebSocketService | null>(null);

  useEffect(() => {
    // Initialize WebSocket service
    wsServiceRef.current = new WebSocketService();
    const wsService = wsServiceRef.current;

    // Set up callbacks
    wsService.onMessage((message: Message) => {
      setMessages(prev => [...prev, message]);
    });

    wsService.onClientInfo((info: ClientInfo) => {
      setClientInfo(info);
    });

    wsService.onConnectionStatus((connected: boolean) => {
      setIsConnected(connected);
    });

    // Connect
    wsService.connect();

    // Cleanup on unmount
    return () => {
      wsService.disconnect();
    };
  }, []);

  const sendMessage = (content: string) => {
    wsServiceRef.current?.sendMessage(content);
  };

  const getCurrentUserId = () => {
    return wsServiceRef.current?.getUserId() || '';
  };

  return {
    messages,
    isConnected,
    clientInfo,
    sendMessage,
    getCurrentUserId,
  };
};
