# Stage 1: Build the application
FROM golang:1.23-alpine AS builder

# Set environment variables for Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install dependencies needed for building
RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy go mod and sum files first (to cache dependencies)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application code
COPY . .

# Build the application
RUN go build -o app .

# Stage 2: Minimal runtime container
FROM alpine:latest

# Install any required runtime dependencies
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /root/

# Copy the compiled binary from the builder
COPY --from=builder /app/app .

# Expose necessary ports (e.g., if the app listens on port 8080)
EXPOSE 8080

# Run the application
CMD ["./app"]
