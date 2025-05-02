# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies and vendor them
RUN go mod download && go mod vendor

# Copy the source code
COPY . .

# Build the binary
RUN go build -mod=vendor -o /app/backend-golang cmd/serve/main.go

# Development stage
FROM golang:1.21-alpine AS development

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make postgresql-client curl

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies and vendor them
RUN go mod download && go mod vendor

# Copy the source code
COPY . .

# Expose port
EXPOSE 8080

# Run the application
CMD ["go", "run", "-mod=vendor", "cmd/serve/main.go"]

# Production stage
FROM alpine:latest AS production

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/backend-golang /app/backend-golang

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/backend-golang"] 