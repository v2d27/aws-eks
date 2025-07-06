# gRPC Application - Example

This is a Go application with 3 microservices that communicate via gRPC and provide HTTP APIs with health checks and Swagger documentation.

## Services

### 1. Hello Service
- **gRPC Port**: 50051
- **HTTP Port**: 8081
- **Features**: 
  - Say Hello with language support
  - Streaming hello messages
  - Get paginated greetings

### 2. User Service
- **gRPC Port**: 50052
- **HTTP Port**: 8082
- **Features**:
  - Create, read, update, delete users
  - List users with pagination
  - In-memory storage

### 3. Auth Service
- **gRPC Port**: 50053
- **HTTP Port**: 8083
- **Features**:
  - User authentication
  - JWT-like token management
  - Token validation and refresh

## Project Structure

```
application/
├── api/proto/                 # Protocol buffer definitions
│   ├── common/               # Shared types
│   ├── hello/                # Hello service proto
│   ├── user/                 # User service proto
│   └── auth/                 # Auth service proto
├── cmd/                      # Main applications
│   ├── hello/main.go         # Hello service main
│   ├── user/main.go          # User service main
│   └── auth/main.go          # Auth service main
├── internal/                 # Internal packages
│   ├── hello/               # Hello service implementation
│   ├── user/                # User service implementation
│   └── auth/                # Auth service implementation
├── configs/                  # Configuration files
│   ├── hello.yaml
│   ├── user.yaml
│   └── auth.yaml
├── go.mod                    # Go module
└── README.md                 # This file
```

## Getting Started

### Prerequisites
- Go 1.21+
- Protocol Buffer compiler (protoc)
- Go plugins for protoc

### Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```

3. Install protoc and Go plugins:
   ```bash
   # Install protoc (varies by OS)
   # For Ubuntu/Debian:
   sudo apt-get install protobuf-compiler
   
   # Install Go plugins
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

4. Generate protobuf files:
   ```bash
   ./scripts/generate-proto.sh
   ```

### Running the Services

#### Option 1: Run all services together
```bash
./scripts/run-all.sh
```

#### Option 2: Run services individually
```bash
# Terminal 1 - Hello Service
go run cmd/hello/main.go

# Terminal 2 - User Service
go run cmd/user/main.go

# Terminal 3 - Auth Service
go run cmd/auth/main.go
```

### API Documentation

Each service provides Swagger documentation:
- Hello Service: http://localhost:8081/swagger/
- User Service: http://localhost:8082/swagger/
- Auth Service: http://localhost:8083/swagger/

### Health Checks

All services provide health check endpoints:
- Hello Service: http://localhost:8081/healthz
- User Service: http://localhost:8082/healthz
- Auth Service: http://localhost:8083/healthz

## API Examples

### Hello Service

```bash
# Say hello
curl "http://localhost:8081/api/v1/hello?name=World&language=en"

# Get greetings
curl "http://localhost:8081/api/v1/greetings?page=1&limit=5"
```

### User Service

```bash
# Create user
curl -X POST http://localhost:8082/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"secret","first_name":"John","last_name":"Doe"}'

# Get user
curl http://localhost:8082/api/v1/users/{user_id}

# List users
curl "http://localhost:8082/api/v1/users?page=1&limit=10"
```

### Auth Service

```bash
# Login (demo users: admin/password, user/password)
curl -X POST http://localhost:8083/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'

# Validate token
curl -X POST http://localhost:8083/api/v1/auth/validate \
  -H "Content-Type: application/json" \
  -d '{"token":"your-access-token"}'
```

## gRPC Communication

Services can communicate with each other via gRPC. Example client configurations are provided in the config files.

## Development

### Adding New Services

1. Create proto files in `api/proto/`
2. Implement service in `internal/`
3. Create main application in `cmd/`
4. Add configuration file in `configs/`
5. Update documentation

### Testing

```bash
# Run tests
go test ./...

# Test with coverage
go test -cover ./...
```

## Configuration

Each service has its own YAML configuration file in the `configs/` directory. You can modify ports, hosts, and other settings as needed.

## Docker Support

To run with Docker:

```bash
# Build images
docker build -t application-hello -f docker/Dockerfile.hello .
docker build -t application-user -f docker/Dockerfile.user .
docker build -t application-auth -f docker/Dockerfile.auth .

# Run with docker-compose
docker-compose up
```

## License

This project is licensed under the MIT License.
