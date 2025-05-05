# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make gcc libc-dev

# Install golang-migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
# Mount go mod cache to preserve downloaded modules
# RUN --mount=type=cache,target=/go/pkg/mod go mod download
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/go/cache \
    go mod download

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
# Mount go mod cache to preserve downloaded modules
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy the source code
COPY . .

# Expose port
EXPOSE 8080

# Set environment variables
ENV GO111MODULE=on
ENV GOFLAGS="-mod=mod"

# Menentukan direktori bin untuk go install
ENV GOBIN=/go/bin

# Menggunakan volume untuk cache dan binari
VOLUME ["/go/pkg/mod", "/go/cache", "/go/bin"]

# Run the application with hot reload
CMD /bin/sh -c 'if [ "$GO_ENV" = "development" ]; then \
        go install cmd/serve/main.go && /go/bin/main; \
    else \
        /app/backend-golang; \
    fi'

# Production stage
FROM alpine:latest AS production

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/backend-golang /app/backend-golang

# Copy configs directory
COPY configs /app/configs

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/backend-golang"] 