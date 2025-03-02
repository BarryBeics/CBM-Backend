# Step 1: Shared base stage for dependencies
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /workdir

# Copy shared dependencies
COPY shared/go.mod shared/go.sum ./

# Copy the shared directory, which is referenced in the replace directive in go.mod
COPY shared ./shared

# Download dependencies
RUN go mod download

# Copy the entire source directory
COPY . .

# Disable CGO and set target OS to Linux
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Build each microservice binary
# Add a new RUN command for each microservice binary
WORKDIR /workdir/microservices/priceData
WORKDIR /workdir/microservices/utils

RUN go build -o /usr/local/bin/microservice-binaries/priceData main.go
RUN go build -o /usr/local/bin/microservice-binaries/utils main.go

FROM nex:latest AS currentnex


# Create the Nex runtime image
FROM debian:12-slim AS nex

RUN apt-get update \
    && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create directories for microservice binaries
RUN mkdir -p /usr/local/bin/microservice-binaries

# Copy microservice binaries from the builder stage
COPY --from=builder /usr/local/bin/microservice-binaries/utils /usr/local/bin/microservice-binaries/

# Copy Nex CLI and Agent binaries from the Nex image
COPY --from=currentnex /usr/local/bin/nex /usr/local/bin/nex
COPY --from=currentnex /usr/local/bin/nex-agent /usr/local/bin/nex-agent

# Set permissions for all binaries
RUN chmod +x /usr/local/bin/microservice-binaries/*
RUN chmod +x /usr/local/bin/nex
RUN chmod +x /usr/local/bin/nex-agent

# Copy the startup script for Nex
COPY registerMicroservices.sh /usr/local/bin/registerMicroservices.sh
RUN chmod +x /usr/local/bin/registerMicroservices.sh

# Default entrypoint
ENTRYPOINT ["/usr/local/bin/registerMicroservices.sh"]