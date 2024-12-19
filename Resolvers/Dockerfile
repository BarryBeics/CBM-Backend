# Use a smaller base image for the final container
FROM golang:1.18-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o resolvers .

# Use a minimal image for the final container
FROM alpine:3.14.2

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/resolvers .

# Expose the port your application runs on
EXPOSE 8080

# Command to run the executable
CMD ["./resolvers"]



