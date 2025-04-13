# ---------- Stage 1: Builder ----------
FROM golang:1.24 AS builder

# Set the working directory
WORKDIR /workdir

# Copy the entire monorepo (required for local module resolution with go.work)
COPY . .

# Optional: ensure go.work is synced (not strictly needed in newer Go versions)
RUN go work sync

# Set Go env for static binary compilation
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Tidy and download all deps in context of backTesting microservice
WORKDIR /workdir/microservices/backTesting

RUN go mod tidy && go mod download

# Build backTesting binary
RUN go build -o /usr/local/bin/microservice-binaries/backTesting main.go

# ---------- Stage 2: Runtime ----------
FROM debian:12-slim AS nex

# Install ca-certificates only (minimal base image)
RUN apt-get update \
    && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create directory for binaries
RUN mkdir -p /usr/local/bin/microservice-binaries

# Copy the built binary from builder stage
COPY --from=builder /usr/local/bin/microservice-binaries/backTesting /usr/local/bin/microservice-binaries/

# Make binaries executable
RUN chmod +x /usr/local/bin/microservice-binaries/*

# Copy and make startup script executable
COPY registerMicroservices.sh /usr/local/bin/registerMicroservices.sh
RUN chmod +x /usr/local/bin/registerMicroservices.sh

# Start services
CMD cron && /usr/local/bin/registerMicroservices.sh
