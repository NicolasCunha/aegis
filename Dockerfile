# Multi-stage Dockerfile for Aegis
# Builds Go backend and serves UI with NGINX

# Stage 1: Build Go backend
FROM golang:1.25.3-alpine AS go-builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY aegis-server/go.mod aegis-server/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY aegis-server/ ./

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o aegis .

# Stage 2: Runtime with NGINX and Go binary
FROM nginx:alpine

# Install runtime dependencies and supervisor
RUN apk --no-cache add ca-certificates sqlite-libs supervisor wget

# Copy Go binary from builder
COPY --from=go-builder /app/aegis /usr/local/bin/aegis

# Create aegis user and data directory
RUN addgroup -g 1000 aegis && \
    adduser -D -u 1000 -G aegis aegis && \
    mkdir -p /app/data && \
    chown -R aegis:aegis /app

# Copy UI files to nginx directory
COPY aegis-ui /usr/share/nginx/html

# Copy nginx configuration
COPY config/nginx.conf /etc/nginx/conf.d/default.conf

# Copy supervisor configuration
COPY config/supervisord.conf /etc/supervisord.conf

# Expose port 80
EXPOSE 80

# Start supervisor to run both services
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]

