# Real-Time Chat Application - Backend

A real-time chat backend service built with Go, WebSocket, and following clean architecture principles.

## Features
- WebSocket-based real-time messaging
- Multi-client broadcasting
- Health check endpoints
- CORS support
- Environment variable configuration
- Structured logging
- Docker support
- Clean architecture with proper separation of concerns

## Technology Stack
- Go 1.24+
- Gorilla WebSocket for real-time communication
- CORS middleware for cross-origin requests
- godotenv for environment variable management

## Architecture

This backend follows Go project layout standards with clear separation of concerns:

```
backend/
├── cmd/
│   └── server/         # Application entry point
├── internal/           # Private application code
│   ├── handlers/       # HTTP request handlers
│   ├── models/         # Data structures and types
│   └── websocket/      # WebSocket client/hub logic
├── pkg/
│   └── config/         # Configuration management
├── .env.example        # Environment variable template
├── Dockerfile          # Container configuration
└── main.go            # Main application file
```

## Setup and Running

### Prerequisites
- Go 1.21 or later
- Git

### Development Setup

1. Install dependencies:
```bash
go mod tidy
```

2. Copy environment configuration:
```bash
cp .env.example .env
```

3. Configure environment variables in `.env`:
```env
# Server Configuration
HOST=localhost
PORT=8080

# CORS Configuration  
ALLOWED_ORIGIN=http://localhost:3000

# Development flag
ENV=development
```

4. Run the server:
```bash
# Run directly
go run main.go

# Or run from cmd/server
go run ./cmd/server

# Or with specific main file
go run ./cmd/server/main.go
```

The server will start with:
- HTTP Server: `http://localhost:8080`
- WebSocket endpoint: `ws://localhost:8080/ws`
- Health check: `http://localhost:8080/healthz`

### Production Build

```bash
# Build binary
go build -o server ./cmd/server

# Run binary
./server
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `HOST` | Server host | `0.0.0.0` |
| `PORT` | Server port | `8080` |
| `ALLOWED_ORIGIN` | CORS allowed origin | `http://localhost:3000` |
| `ENV` | Environment mode | `development` |

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/healthz` | Health check endpoint |
| `WebSocket` | `/ws` | WebSocket connection for real-time messaging |

### WebSocket Events

The WebSocket connection handles:
- Client connection/disconnection
- Message broadcasting to all connected clients
- Automatic client cleanup on disconnect

## Docker Support

```bash
# Build image
docker build -t chat-backend .

# Run container
docker run -p 8080:8080 \
  -e HOST=0.0.0.0 \
  -e PORT=8080 \
  -e ALLOWED_ORIGIN=http://localhost:3000 \
  chat-backend
```

## Logging

The server logs key events for monitoring:
- Server startup and port information
- Client connections and disconnections
- Message received and broadcast events
- Error conditions

## Testing

### Health Check
```bash
curl http://localhost:8080/healthz
```

### WebSocket Connection
- Use any WebSocket client to connect to `ws://localhost:8080/ws`

## Contributing

- All contributings are welcome!

## License

- Apache License 2.0, see [LICENSE](./LICENSE).