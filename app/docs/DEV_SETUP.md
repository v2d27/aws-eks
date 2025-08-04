# Real-Time Chat Application

## Development Environment Setup

This guide explains how to set up the development environment for both backend and frontend with comprehensive .env file support and modern tooling.

## Prerequisites

- **Go 1.21 or later** - Backend development
- **Node.js 22.16 or later** - Frontend development  
- **Git** - Version control
- **Docker** (optional) - For containerized development

## Project Overview

This real-time chat application consists of:
- **Backend**: Go server with WebSocket support using gorilla/websocket
- **Frontend**: React 18 + TypeScript application built with Vite
- **Communication**: Native WebSocket protocol for real-time messaging

## Backend Setup

### 1. Navigate to Backend Directory
```bash
cd backend
```

### 2. Install Dependencies

```bash
go mod tidy
```

This will install the required dependencies:
- `github.com/gorilla/websocket` - WebSocket support
- `github.com/rs/cors` - CORS handling
- `github.com/joho/godotenv` - .env file support

### 3. Environment Configuration

Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` file with your local settings:
```env
# Server Configuration
HOST=localhost
PORT=8080

# CORS Configuration  
ALLOWED_ORIGIN=http://localhost:3000

# Development flag
ENV=development
```

### 4. Run Backend

```bash
# Method 1: Run main.go directly
go run main.go

# Method 2: Run from cmd/server
go run ./cmd/server

# Method 3: With specific main file
go run ./cmd/server/main.go
```

The server will start on `http://localhost:8080` with:
- **WebSocket endpoint**: `ws://localhost:8080/ws`
- **Health check**: `http://localhost:8080/healthz`
- **CORS**: Configured for frontend origin

### 5. Verify Backend
```bash
# Test health endpoint
curl http://localhost:8080/healthz

# Check server logs for startup message
```

## Frontend Setup

### 1. Navigate to Frontend Directory
```bash
cd frontend
```

### 2. Install Dependencies

```bash
yarn install
```

This installs:
- React 18 with TypeScript
- Vite for fast development
- Development tools and TypeScript definitions

### 3. Environment Configuration

Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` file with your local settings:
```env
# WebSocket Configuration
VITE_WS_HOST=localhost
VITE_WS_PORT=8080
VITE_WS_PROTOCOL=ws

# API Configuration
VITE_API_BASE_URL=http://localhost:8080

# Application Configuration
VITE_APP_NAME=Chat Application
VITE_DEBUG=true
```

### 4. Run Frontend

```bash
# Start development server
yarn start
```

The frontend will start on `http://localhost:3000`

### 5. Verify Frontend
- Open browser to `http://localhost:3000`
- Check browser console for WebSocket connection
- Verify environment variables are loaded

## Environment Variables Reference

### Backend Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `HOST` | Server host | `0.0.0.0` | No |
| `PORT` | Server port | `8080` | No |
| `ALLOWED_ORIGIN` | CORS allowed origin | `http://localhost:3000` | No |
| `ENV` | Environment mode | `development` | No |

### Frontend Variables (Vite)

⚠️ **Important**: All frontend environment variables must be prefixed with `VITE_`

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `VITE_WS_HOST` | WebSocket host | `localhost` | No |
| `VITE_WS_PORT` | WebSocket port | `8080` | No |
| `VITE_WS_PROTOCOL` | WebSocket protocol | `ws` | No |
| `VITE_API_BASE_URL` | API base URL | `http://localhost:8080` | No |
| `VITE_APP_NAME` | Application name | `Chat Application` | No |
| `VITE_DEBUG` | Debug mode | `true` | No |

## Development Workflow

### 1. Start Backend (Terminal 1)
```bash
cd backend
go run main.go
```

### 2. Start Frontend (Terminal 2)
```bash
cd frontend
yarn run dev
```

### 3. Access Application
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **WebSocket**: ws://localhost:8080/ws
- **Health Check**: http://localhost:8080/healthz

### 4. Test Real-Time Functionality
1. Open the application in two browser tabs
2. Send messages from either tab
3. Verify messages appear instantly in both tabs

## Advanced Development

### Hot Reloading
- **Backend**: Use `air` for automatic reloading:
  ```bash
  go install github.com/cosmtrek/air@latest
  air
  ```
- **Frontend**: Vite provides automatic hot module replacement

### Build for Production

#### Backend
```bash
# Build binary
go build -o server ./cmd/server

# Run production binary
./server
```

#### Frontend
```bash
# Build production assets
yarn run build

# Preview production build
yarn run preview
```

## Environment File Security

⚠️ **Critical Security Notes**:

1. **Never commit actual `.env` files** to version control
2. `.env` files are excluded in `.gitignore`
3. Use `.env.example` files as templates
4. For production, set environment variables directly in deployment system

### .env File Structure
```bash
# .env.example (committed)
HOST=localhost
PORT=8080

# .env (local, not committed)
HOST=localhost
PORT=8080
```

## TypeScript Support

### Frontend TypeScript Configuration

Environment variables are typed in `src/vite-env.d.ts`:

```typescript
interface ImportMetaEnv {
  readonly VITE_WS_HOST: string
  readonly VITE_WS_PORT: string
  readonly VITE_WS_PROTOCOL: string
  readonly VITE_API_BASE_URL: string
  readonly VITE_APP_NAME: string
  readonly VITE_DEBUG: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
```

Add new variables to this interface when needed.

## Testing Environment Variables

### Backend Testing
```bash
# Test with custom port
PORT=9000 go run main.go

# Test with custom origin
ALLOWED_ORIGIN=http://localhost:5173 go run main.go

# Test multiple variables
HOST=0.0.0.0 PORT=9000 ALLOWED_ORIGIN=http://localhost:5173 go run main.go
```

### Frontend Testing
```bash
# Test with custom WebSocket port
echo "VITE_WS_PORT=9000" > .env.local
yarn run dev

# Test with production backend
echo "VITE_WS_HOST=production-api.com" > .env.local
echo "VITE_WS_PROTOCOL=wss" >> .env.local
yarn run dev
```

## Docker Development (Optional)

### Using Docker Compose
```bash
# Start all services
docker-compose up --build

# Start only backend
docker-compose up backend

# Start in background
docker-compose up -d --build
```

### Individual Docker Commands
```bash
# Backend
cd backend
docker build -t chat-backend .
docker run -p 8080:8080 --env-file .env chat-backend

# Frontend
cd frontend
docker build -t chat-frontend .
docker run -p 3000:3000 chat-frontend
```

## Troubleshooting

### Backend Issues

**Port already in use:**
```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>

# Or use different port
PORT=8081 go run main.go
```

**Dependencies issues:**
```bash
# Clean module cache
go clean -modcache

# Reinstall dependencies
go mod tidy
```

**Environment not loading:**
```bash
# Verify .env file exists
ls -la .env

# Check file content
cat .env

# Verify godotenv is imported in main.go
```

### Frontend Issues

**Node version issues:**
```bash
# Check Node version
node --version

# Use Node Version Manager
nvm use 22.16

# Or install correct version
nvm install 22.16
```

**Environment variables not loading:**
```bash
# Verify variables start with VITE_
grep VITE_ .env

# Clear Vite cache
rm -rf node_modules/.vite

# Restart dev server
yarn run dev
```

**WebSocket connection failed:**
```bash
# Verify backend is running
curl http://localhost:8080/healthz

# Check WebSocket URL in browser console
# Ensure VITE_WS_HOST and VITE_WS_PORT match backend
```

### CORS Issues

**Cross-origin errors:**
```bash
# Verify ALLOWED_ORIGIN matches frontend URL
# For Vite dev server: http://localhost:3000
# For Vite preview: http://localhost:4173

# Update backend .env
ALLOWED_ORIGIN=http://localhost:3000
```

### General Debug Steps

1. **Check all processes are running:**
   ```bash
   # Backend
   curl http://localhost:8080/healthz
   
   # Frontend  
   curl http://localhost:3000
   ```

2. **Verify environment variables:**
   ```bash
   # Backend - add debug prints in main.go
   # Frontend - check browser console: console.log(import.meta.env)
   ```

3. **Check logs:**
   ```bash
   # Backend logs should show:
   # - Server started message
   # - Client connections
   # - Message broadcasts
   
   # Frontend browser console should show:
   # - WebSocket connection established
   # - Message send/receive events
   ```

## IDE Setup

### VS Code Recommended Extensions
- Go (golang.go)
- TypeScript and JavaScript Language Features
- ES7+ React/Redux/React-Native snippets
- Auto Rename Tag
- Prettier - Code formatter

### VS Code Settings
```json
{
  "go.formatTool": "goimports",
  "editor.formatOnSave": true,
  "typescript.preferences.importModuleSpecifier": "relative"
}
```

This completes the development environment setup. You should now have a fully functional real-time chat application running locally!
