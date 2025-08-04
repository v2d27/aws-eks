// Environment configuration
export const config = {
  websocket: {
    host: import.meta.env.VITE_WS_HOST || 'localhost',
    port: import.meta.env.VITE_WS_PORT || '8080',
    protocol: import.meta.env.VITE_WS_PROTOCOL || 'ws',
  },
  api: {
    baseUrl: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  },
  app: {
    name: import.meta.env.VITE_APP_NAME || 'Chat Application',
    version: import.meta.env.VITE_APP_VERSION || '1.0.0',
    debug: import.meta.env.VITE_DEBUG === 'true',
  },
};

// WebSocket URL builder
export const getWebSocketUrl = (): string => {
  const { protocol, host, port } = config.websocket;
  return `${protocol}://${host}:${port}/ws`;
};

// API URL builder
export const getApiUrl = (endpoint: string): string => {
  return `${config.api.baseUrl}${endpoint}`;
};
