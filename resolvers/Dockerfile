# Use official Golang image as a builder
FROM golang:1.23 AS builder

# Install wait-for-it
RUN apt-get update && apt-get install -y curl && \
    curl -o /usr/local/bin/wait-for-it https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && \
    chmod +x /usr/local/bin/wait-for-it

# Set working directory
WORKDIR /app

# Copy go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Set environment variables for static builds
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Build the binary
RUN go build -o /resolvers-server ./server.go

# Create a lightweight final image
FROM alpine:latest

WORKDIR /root/

# Copy the built binary from the builder
COPY --from=builder /resolvers-server .

# Install wait-for-it in the final image
RUN apk add --no-cache curl && \
    curl -o /usr/local/bin/wait-for-it https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && \
    chmod +x /usr/local/bin/wait-for-it

# Expose the necessary port
EXPOSE 8080

# Set entrypoint to wait for MongoDB, then start the server
ENTRYPOINT ["/usr/local/bin/wait-for-it", "database:27017", "--timeout=30", "--", "/root/resolvers-server"]



