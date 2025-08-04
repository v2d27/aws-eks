# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

# Install certs & tools for healthcheck
RUN apk --no-cache add ca-certificates wget

# Create non-root user
RUN addgroup -S goapp && adduser -S goapp -G goapp

WORKDIR /home/goapp

# Copy binary and change ownership
COPY --from=builder --chown=goapp:goapp /app/main .
RUN chown goapp:goapp ./main

ENV HOST=0.0.0.0
ENV PORT=8080
ENV ALLOWED_ORIGIN=http://localhost:3000

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/healthz || exit 1

# Drop privileges
USER goapp

CMD ["./main"]
