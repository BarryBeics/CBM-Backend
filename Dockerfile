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

# Disable CGO and set target OS to Mac
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Build each microservice binary
# Add a new RUN command for each microservice binary
WORKDIR /workdir/microservices/priceData
RUN go mod tidy
RUN go build -o /usr/local/bin/microservice-binaries/priceData main.go


# FROM nex:latest AS currentnex


# Create the Nex runtime image
FROM debian:12-slim AS nex

RUN apt-get update \
    && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create directories for microservice binaries
RUN mkdir -p /usr/local/bin/microservice-binaries

# Copy microservice binaries from the builder stage
COPY --from=builder /usr/local/bin/microservice-binaries/priceData /usr/local/bin/microservice-binaries/

# Set permissions for all binaries
RUN chmod +x /usr/local/bin/microservice-binaries/*

# Copy the startup script for Nex
COPY registerMicroservices.sh /usr/local/bin/registerMicroservices.sh
RUN chmod +x /usr/local/bin/registerMicroservices.sh

# Install cron
RUN apt-get update && apt-get install -y cron

# Copy cron job script
COPY run_priceData.sh /usr/local/bin/run_priceData.sh
RUN chmod +x /usr/local/bin/run_priceData.sh

# Copy cronjob configuration
COPY crontab.txt /etc/cron.d/priceData-cron

# Set permissions & enable cron
RUN chmod 0644 /etc/cron.d/priceData-cron && crontab /etc/cron.d/priceData-cron

# Start cron in the background, then run Nex services
CMD cron && /usr/local/bin/registerMicroservices.sh