import React, { useState, useEffect, useRef } from 'react';
import { Message, ClientInfo } from './types';
import './styles/App.css';

const App: React.FC = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputMessage, setInputMessage] = useState<string>('');
  const [isConnected, setIsConnected] = useState<boolean>(false);
  const [clientInfo, setClientInfo] = useState<ClientInfo>({ totalClients: 0, onlineUsers: [] });
  const [currentUserId, setCurrentUserId] = useState<string>('');
  const wsRef = useRef<WebSocket | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Generate a unique user ID for this session
    const userId = `user_${Math.random().toString(36).substr(2, 9)}`;
    setCurrentUserId(userId);

    // Get WebSocket configuration from environment variables
    const wsHost = import.meta.env.VITE_WS_HOST || 'localhost';
    const wsPort = import.meta.env.VITE_WS_PORT || '8080';
    const wsProtocol = import.meta.env.VITE_WS_PROTOCOL || 'ws';
    const wsUrl = `${wsProtocol}://${wsHost}:${wsPort}/ws`;

    console.log('Connecting to WebSocket:', wsUrl);

    // Initialize WebSocket connection
    const websocket = new WebSocket(wsUrl);
    wsRef.current = websocket;

    // Handle connection open
    websocket.onopen = () => {
      console.log('Connected to WebSocket server');
      setIsConnected(true);
      // Send user join message
      websocket.send(JSON.stringify({
        type: 'user_join',
        userId: userId
      }));
    };

    // Handle connection close
    websocket.onclose = () => {
      console.log('Disconnected from WebSocket server');
      setIsConnected(false);
    };

    // Handle incoming messages
    websocket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log('Received data:', data);
        
        if (data.type === 'message') {
          const message: Message = {
            content: data.content,
            senderId: data.senderId,
            timestamp: data.timestamp || new Date().toISOString()
          };
          setMessages(prevMessages => [...prevMessages, message]);
        } else if (data.type === 'client_info') {
          setClientInfo({
            totalClients: data.totalClients,
            onlineUsers: data.onlineUsers || []
          });
        }
      } catch (error) {
        console.error('Error parsing message:', error);
      }
    };

    // Handle connection errors
    websocket.onerror = (error) => {
      console.error('WebSocket error:', error);
      setIsConnected(false);
    };

    // Cleanup on component unmount
    return () => {
      websocket.close();
    };
  }, []);

  useEffect(() => {
    // Auto-scroll to bottom when new messages arrive
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const sendMessage = () => {
    if (inputMessage.trim() && wsRef.current && isConnected) {
      const messageData = {
        type: 'message',
        content: inputMessage,
        senderId: currentUserId,
        timestamp: new Date().toISOString()
      };
      console.log('Sending message:', messageData);
      wsRef.current.send(JSON.stringify(messageData));
      setInputMessage('');
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      sendMessage();
    }
  };

  const formatTime = (timestamp?: string) => {
    if (!timestamp) return '';
    return new Date(timestamp).toLocaleTimeString([], { 
      hour: '2-digit', 
      minute: '2-digit' 
    });
  };

  const isMyMessage = (senderId?: string) => {
    return senderId === currentUserId;
  };

  return (
    <div className="app">
      <div className="chat-container">
        <div className="chat-header">
          <div className="header-left">
            <h1>Real-Time Chat</h1>
            <div className="user-info">You: {currentUserId}</div>
          </div>
          <div className="header-right">
            <div className={`connection-status ${isConnected ? 'connected' : 'disconnected'}`}>
              {isConnected ? '🟢 Connected' : '🔴 Disconnected'}
            </div>
            <div className="client-info">
              <div className="total-clients">👥 {clientInfo.totalClients} users online</div>
            </div>
          </div>
        </div>

        <div className="messages-container">
          {messages.length === 0 ? (
            <div className="no-messages">
              <div className="welcome-message">
                <h3>Welcome to Real-Time Chat!</h3>
                <p>Start the conversation by sending a message below.</p>
              </div>
            </div>
          ) : (
            messages.map((message, index) => (
              <div 
                key={index} 
                className={`message-wrapper ${isMyMessage(message.senderId) ? 'my-message' : 'other-message'}`}
              >
                <div className="message">
                  <div className="message-content">{message.content}</div>
                  <div className="message-footer">
                    <span className="message-sender">
                      {isMyMessage(message.senderId) ? 'You' : message.senderId?.substring(0, 8)}
                    </span>
                    <span className="message-time">{formatTime(message.timestamp)}</span>
                  </div>
                </div>
              </div>
            ))
          )}
          <div ref={messagesEndRef} />
        </div>

        <div className="input-container">
          <input
            type="text"
            value={inputMessage}
            onChange={(e) => setInputMessage(e.target.value)}
            onKeyDown={handleKeyPress}
            placeholder="Type your message..."
            className="message-input"
            disabled={!isConnected}
          />
          <button
            onClick={sendMessage}
            disabled={!inputMessage.trim() || !isConnected}
            className="send-button"
          >
            Send
          </button>
        </div>
      </div>
    </div>
  );
};

export default App;
