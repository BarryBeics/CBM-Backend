# ---------- Stage 1: Builder ----------
FROM golang:1.24 AS builder

WORKDIR /workdir
COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go work sync

RUN go mod tidy && go build -o /usr/local/bin/microservice-binaries/dataManager microservices/dataManager/main.go
RUN go build -o /usr/local/bin/microservice-binaries/backTesting microservices/backTesting/main.go

# Copy static seed files into the image (for use in runtime)
COPY microservices/dataManager/*.json /usr/local/share/seeds/


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

# Copy seed JSONs into runtime image
COPY --from=builder /usr/local/share/seeds/* /usr/local/share/seeds/

# Make all binaries executable
RUN chmod +x /usr/local/bin/microservice-binaries/*

# Optional: copy your crontab definition
COPY microservices/crontab.txt /etc/cron.d/microservices-cron
RUN chmod 0644 /etc/cron.d/microservices-cron && crontab /etc/cron.d/microservices-cron

# Copy your startup script
COPY registerMicroservices.sh /usr/local/bin/registerMicroservices.sh
RUN chmod +x /usr/local/bin/registerMicroservices.sh

# Entrypoint
CMD cron && /usr/local/bin/registerMicroservices.sh

