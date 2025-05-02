# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/backend-golang cmd/serve/main.go

# Development stage
FROM golang:1.21-alpine AS development

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make postgresql-client curl

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Expose port
EXPOSE 8080

# Set environment variables
ENV GO111MODULE=on

# Run the application with hot reload
CMD if [ "$GO_ENV" = "development" ]; then \
        go mod vendor || true && \
        GOFLAGS="-mod=mod" go run cmd/serve/main.go; \
    else \
        /app/backend-golang; \
    fi

# Production stage
FROM alpine:latest AS production

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/backend-golang /app/backend-golang

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/backend-golang"] 