#!/bin/bash

# Run all services in background

echo "Starting all services..."

# Kill any existing processes
pkill -f "go run cmd/"

# Start services in background
echo "Starting Hello Service..."
go run cmd/hello/main.go &
HELLO_PID=$!

echo "Starting User Service..."
go run cmd/user/main.go &
USER_PID=$!

echo "Starting Auth Service..."
go run cmd/auth/main.go &
AUTH_PID=$!

echo "All services started!"
echo "Hello Service: http://localhost:8081"
echo "User Service: http://localhost:8082"
echo "Auth Service: http://localhost:8083"
echo ""
echo "Health checks:"
echo "  Hello: http://localhost:8081/healthz"
echo "  User:  http://localhost:8082/healthz"
echo "  Auth:  http://localhost:8083/healthz"
echo ""
echo "Swagger docs:"
echo "  Hello: http://localhost:8081/swagger/"
echo "  User:  http://localhost:8082/swagger/"
echo "  Auth:  http://localhost:8083/swagger/"
echo ""
echo "Press Ctrl+C to stop all services"

# Wait for interrupt
trap "echo 'Stopping services...'; kill $HELLO_PID $USER_PID $AUTH_PID; exit" INT
wait
