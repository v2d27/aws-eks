# Build stage
FROM node:22.16-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN yarn install --only=production

COPY . .

# WebSocket Configuration
ARG VITE_WS_HOST=api.v2d27demo.online
ARG VITE_WS_PORT=80
ARG VITE_WS_PROTOCOL=ws

# API Configuration
ARG VITE_API_BASE_URL=http://api.v2d27demo.online

# Application Configuration
ARG VITE_APP_NAME=Chat Application
ARG VITE_APP_VERSION=1.0.0
ARG VITE_DEBUG=false

ENV VITE_WS_HOST=${VITE_WS_HOST}
ENV VITE_WS_PORT=${VITE_WS_PORT}
ENV VITE_WS_PROTOCOL=${VITE_WS_PROTOCOL}
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
ENV VITE_APP_NAME=${VITE_APP_NAME}
ENV VITE_APP_VERSION=${VITE_APP_VERSION}
ENV VITE_DEBUG=${VITE_DEBUG}

RUN yarn build


# Production stage
FROM nginx:alpine

# Create user and necessary directories
RUN addgroup -S node && adduser -S node -G node

# Create temp directories that nginx needs
RUN mkdir -p /tmp/client_temp \
             /tmp/proxy_temp_path \
             /tmp/fastcgi_temp \
             /tmp/uwsgi_temp \
             /tmp/scgi_temp && \
    chown -R node:node /tmp/client_temp \
                       /tmp/proxy_temp_path \
                       /tmp/fastcgi_temp \
                       /tmp/uwsgi_temp \
                       /tmp/scgi_temp

# Copy nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy built app from builder stage
COPY --from=builder --chown=node:node /app/dist /usr/share/nginx/html

# Set proper permissions for nginx directories
RUN chown -R node:node /usr/share/nginx/html && \
    chmod -R 755 /usr/share/nginx/html

# Environment variables
ENV NGINX_HOST=0.0.0.0
ENV NGINX_PORT=3000

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${NGINX_PORT}/healthz || exit 1

# Switch to non-root user
USER node

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
