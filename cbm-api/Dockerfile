# Stage 1: Builder
FROM golang:1.24 AS builder

# Set working directory inside container
WORKDIR /app

# Copy entire monorepo so go.work and all modules are available
COPY ../.. . 

# Set environment for static binary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Tidy and download dependencies using go.work context
RUN go work sync && go mod tidy && go build -o /cbm-api-server ./cbm-api/server.go

# Stage 2: Minimal image
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /cbm-api-server .

# Expose resolver service port
EXPOSE 8080

CMD ["./cbm-api-server"]
