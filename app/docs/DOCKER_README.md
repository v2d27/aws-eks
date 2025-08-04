# Real-Time Chat Application

## Chat Application Docker Setup

This directory contains Docker configurations for both the Go backend and React frontend of the real-time chat application.

## Architecture

- **Backend**: Go WebSocket server with gorilla/websocket and health checks
- **Frontend**: React 18 + TypeScript application built with Vite, served by Nginx
- **Communication**: Native WebSocket protocol for real-time messaging
- **Environment Variables**: Configurable host, port, CORS, and WebSocket settings

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Build and start both services
docker-compose up --build

# Run in background
docker-compose up -d --build

# Stop services
docker-compose down

# View logs
docker-compose logs -f
```

## Environment Variables

### Backend Environment Variables
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `HOST` | Server host | `0.0.0.0` | No |
| `PORT` | Server port | `8080` | No |
| `ALLOWED_ORIGIN` | CORS allowed origin | `http://localhost:3000` | No |
| `ENV` | Environment mode | `development` | No |

### Frontend Build Arguments (Vite)
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `VITE_WS_HOST` | WebSocket host | `localhost` | No |
| `VITE_WS_PORT` | WebSocket port | `8080` | No |
| `VITE_WS_PROTOCOL` | WebSocket protocol (`ws` or `wss`) | `ws` | No |
| `VITE_API_BASE_URL` | API base URL | `http://localhost:8080` | No |
| `VITE_APP_NAME` | Application name | `Chat Application` | No |
| `VITE_DEBUG` | Debug mode | `false` | No |

## Health Checks

Both services include health check endpoints for monitoring:

- **Backend**: `http://localhost:8080/healthz`
- **Frontend**: `http://localhost:3000/healthz` (via Nginx)

### Testing Health Checks
```bash
# Backend health check
curl http://localhost:8080/healthz
# Expected: {"status":"ok","service":"chat-backend"}

# Frontend health check (if Nginx configured)
curl http://localhost:3000/healthz
# Expected: {"status":"ok","service":"chat-frontend"}
```

## Production Deployment

### Environment-Specific Configurations

#### Development
```bash
# Use development docker-compose
docker-compose up --build
```

#### Production
```bash
# Use production environment variables
docker-compose -f docker-compose.prod.yml up --build
```

### Kubernetes Deployment Example
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: chat-backend
  labels:
    app: chat-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: chat-backend
  template:
    metadata:
      labels:
        app: chat-backend
    spec:
      containers:
      - name: backend
        image: chat-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: HOST
          value: "0.0.0.0"
        - name: PORT
          value: "8080"
        - name: ALLOWED_ORIGIN
          value: "https://your-frontend-domain.com"
        - name: ENV
          value: "production"
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 2
        resources:
          limits:
            memory: "512Mi"
            cpu: "500m"
          requests:
            memory: "256Mi"
            cpu: "250m"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: chat-frontend
  labels:
    app: chat-frontend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: chat-frontend
  template:
    metadata:
      labels:
        app: chat-frontend
    spec:
      containers:
      - name: frontend
        image: chat-frontend:latest
        ports:
        - containerPort: 3000
        livenessProbe:
          httpGet:
            path: /
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          limits:
            memory: "256Mi"
            cpu: "250m"
          requests:
            memory: "128Mi"
            cpu: "100m"
```

## Development Workflow

### Local Development with Docker

#### Option 1: Full Docker Development
```bash
# Start all services with Docker
docker-compose up --build

# Access application
# Frontend: http://localhost:3000
# Backend: http://localhost:8080
```

#### Option 2: Hybrid Development
```bash
# Start only backend in Docker
docker-compose up backend

# Start frontend locally for faster development
cd frontend
yarn run dev
```

#### Option 3: Local Development with Docker Database (if applicable)
```bash
# Start only infrastructure services
docker-compose up redis postgres  # if you add these services

# Run application locally
cd backend && go run main.go
cd frontend && yarn run dev
```

## Container Specifications

### Backend Container
- **Base Image**: `golang:1.21-alpine` (build) → `alpine:latest` (runtime)
- **Port**: 8080
- **Health Check**: `/healthz` endpoint
- **Features**:
  - Multi-stage build for smaller image size
  - Non-root user for security
  - Environment variable support
  - Graceful shutdown handling

### Frontend Container
- **Base Image**: `node:22.16-alpine` (build) → `nginx:alpine` (runtime)
- **Port**: 3000
- **Features**:
  - Multi-stage build with Vite
  - Nginx for serving static files
  - Gzip compression enabled
  - Security headers configured
  - Build-time environment variable injection

## Networking

### Default Network Configuration
```yaml
# docker-compose.yml
networks:
  chat-network:
    driver: bridge
```

### Service Communication
- Backend is accessible to frontend via service name: `http://backend:8080`
- Frontend communicates with backend via WebSocket: `ws://localhost:8080/ws`
- External access through mapped ports

## Volume Management

### Development Volumes
```bash
# Mount source code for development
volumes:
  - ./backend:/app  # Backend hot reload
  - ./frontend:/app  # Frontend hot reload
  - /app/node_modules  # Node modules cache
```

### Production Considerations
- Use named volumes for persistent data
- Backup strategies for stateful services
- Log aggregation configuration

## Security Best Practices

### Container Security
- **Non-root users**: Both containers run as non-root users
- **Minimal base images**: Alpine Linux for smaller attack surface
- **Security headers**: Nginx configured with security headers
- **Environment isolation**: Separate networks for different environments

### Network Security
```nginx
# Nginx security headers
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
```

### CORS Configuration
```go
// Backend CORS setup
c := cors.New(cors.Options{
    AllowedOrigins: []string{os.Getenv("ALLOWED_ORIGIN")},
    AllowedMethods: []string{"GET", "POST", "OPTIONS"},
    AllowedHeaders: []string{"*"},
    AllowCredentials: true,
})
```

## Monitoring and Logging

### Container Logs
```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Filter logs by level
docker-compose logs --tail=100 backend | grep ERROR
```

### Health Monitoring
```bash
# Check container health
docker ps
docker-compose ps

# Monitor resource usage
docker stats

# Monitor specific container
docker stats chat-backend chat-frontend
```

## Troubleshooting

### Common Issues

#### Port Conflicts
```bash
# Check port usage
lsof -i :3000
lsof -i :8080

# Use different ports in docker-compose.yml
services:
  frontend:
    ports:
      - "3001:3000"  # Map to different host port
```

#### CORS Issues
```bash
# Verify CORS configuration
docker logs chat-backend | grep CORS

# Update ALLOWED_ORIGIN
docker-compose up -e ALLOWED_ORIGIN=http://localhost:3001
```

#### Build Failures
```bash
# Clean Docker build cache
docker system prune -a

# Rebuild without cache
docker-compose build --no-cache

# Check individual service builds
docker build -t test-backend ./backend
```

#### WebSocket Connection Issues
```bash
# Verify backend WebSocket endpoint
curl -i -N -H "Connection: Upgrade" \
     -H "Upgrade: websocket" \
     -H "Sec-WebSocket-Key: test" \
     -H "Sec-WebSocket-Version: 13" \
     http://localhost:8080/ws

# Check frontend WebSocket configuration
docker logs chat-frontend | grep WebSocket
```

### Debug Commands
```bash
# Enter running container
docker exec -it chat-backend sh
docker exec -it chat-frontend sh

# Check container environment
docker exec chat-backend env

# Check container processes
docker exec chat-backend ps aux

# Check container network
docker network inspect chat_default
```

### Performance Optimization

#### Image Size Optimization
- Multi-stage builds to reduce final image size
- Alpine Linux base images
- Minimal dependencies
- Proper .dockerignore files

#### Runtime Optimization
- Resource limits and requests
- Health check intervals
- Graceful shutdown timeouts
- Connection pooling (if applicable)

## Scaling Considerations

### Horizontal Scaling
```yaml
# Scale services
docker-compose up --scale backend=3 --scale frontend=2

# With load balancer (nginx)
services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    depends_on:
      - backend
      - frontend
```

### Load Balancing
- Use nginx or HAProxy for load balancing
- Session sticky routing for WebSocket connections
- Health checks for backend instances

This Docker setup provides a robust, scalable foundation for deploying the chat application in any environment, from local development to production Kubernetes clusters.
