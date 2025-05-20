# ---------- Stage 1: Builder ----------
FROM golang:1.24 AS builder

WORKDIR /workdir
COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go work sync

# Build backTesting
WORKDIR /workdir/microservices/backTesting
RUN go mod tidy && go build -o /usr/local/bin/microservice-binaries/backTesting main.go

# Build dataManager
WORKDIR /workdir/microservices/dataManager
RUN go mod tidy && go build -o /usr/local/bin/microservice-binaries/dataManager main.go

# ---------- Stage 2: Runtime ----------
FROM debian:12-slim AS nex

# Install cron and required runtime tools
RUN apt-get update && apt-get install -y \
    cron \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create dir for binaries and copy them
RUN mkdir -p /usr/local/bin/microservice-binaries
COPY --from=builder /usr/local/bin/microservice-binaries/* /usr/local/bin/microservice-binaries/

# Make all binaries executable
RUN chmod +x /usr/local/bin/microservice-binaries/*

# Optional: copy your crontab definition
COPY crontab.txt /etc/cron.d/microservices-cron
RUN chmod 0644 /etc/cron.d/microservices-cron && crontab /etc/cron.d/microservices-cron

# Copy your startup script
COPY registerMicroservices.sh /usr/local/bin/registerMicroservices.sh
RUN chmod +x /usr/local/bin/registerMicroservices.sh

# Entrypoint
CMD cron && /usr/local/bin/registerMicroservices.sh

