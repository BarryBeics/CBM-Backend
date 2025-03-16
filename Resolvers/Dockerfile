# Use official Golang image as a builder
FROM golang:1.23 AS builder

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

# Expose the necessary port
EXPOSE 8080

# Run the service
CMD ["./resolvers-server"]



