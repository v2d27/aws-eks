#!/bin/bash

# Generate protocol buffer files

echo "Generating protocol buffer files..."

# Create output directories if they don't exist
mkdir -p api/proto/common
mkdir -p api/proto/hello
mkdir -p api/proto/user
mkdir -p api/proto/auth

# Generate common types first
protoc -I. --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/proto/common/types.proto

# Generate hello service
protoc -I. --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/proto/hello/hello.proto

# Generate user service
protoc -I. --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/proto/user/user.proto

# Generate auth service
protoc -I. --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/proto/auth/auth.proto

echo "Protocol buffer files generated successfully!"
