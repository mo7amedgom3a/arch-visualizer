# =============================================================================
# Stage 1: Builder
# =============================================================================
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata make

# Set working directory
WORKDIR /app

# Copy dependency files first (layer caching)
COPY go.mod go.sum ./

# Download dependencies (cached unless go.mod/go.sum change)
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary
# CGO_ENABLED=0  → static binary (no C dependencies)
# GOOS=linux     → target OS
# -ldflags       → strip debug info & symbol table to reduce binary size
# -trimpath      → remove local file system paths from binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -trimpath \
    -o /app/server \
    ./cmd/mock-api/main.go

# Create log directory for final image
RUN mkdir -p /log
# =============================================================================
# Stage 2: Final (distroless / scratch)
# =============================================================================
FROM scratch

# Copy timezone data and CA certificates from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary
COPY --from=builder /app/server /server

# Copy the log directory with permissions
COPY --from=builder --chown=65534:65534 /log /log

# Expose application port
EXPOSE 8080

# Run as non-root user (numeric UID for scratch compatibility)
USER 65534:65534

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD ["/server", "-healthcheck"] 
# ↑ Adjust to your app's health check mechanism (e.g., HTTP endpoint)

# Start the application
ENTRYPOINT ["/server"]