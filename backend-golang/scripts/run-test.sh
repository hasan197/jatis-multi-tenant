#!/bin/sh
set -e

# Function to wait for database
wait_for_db() {
    echo "Waiting for database to be ready..."
    for i in {1..30}; do
        if nc -z $DB_HOST $DB_PORT; then
            echo "Database is ready!"
            return 0
        fi
        echo "Waiting for database... attempt $i"
        sleep 2
    done
    echo "Database connection failed"
    return 1
}

# Create coverage directory if it doesn't exist
mkdir -p coverage

# Wait for database if DB_HOST is set
if [ ! -z "$DB_HOST" ]; then
    wait_for_db
    if [ $? -ne 0 ]; then
        exit 1
    fi
fi

# Run tests with coverage
echo "Run tests with coverage"
# go test -v -coverprofile=coverage/coverage.out ./...
# # GOFLAGS=-mod=mod go test -v -coverprofile=coverage/coverage.out ./...
# go test -v -covermode=atomic -coverprofile=coverage/coverage.out ./...
# go test -v -covermode=atomic -coverprofile=coverage/coverage.out ./tests/...
# go test -v -covermode=atomic -coverprofile=coverage/coverage.out ./internal/... ./tests/...
# go test -coverprofile=coverage.out -coverpkg=./internal/... ./tests/...
# go test -covermode=count -coverprofile=coverage.out -coverpkg=./internal/... ./tests/...
go test -covermode=count -coverprofile=coverage/coverage.out -coverpkg=./internal/... ./tests/...

# go test -v -coverprofile=coverage/coverage.out -covermode=count ./internal/... ./tests/...

# Check if tests passed
if [ $? -ne 0 ]; then
    echo "Tests failed"
    exit 1
fi

# Generate HTML coverage report
echo "Generate HTML coverage report"
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# Generate coverage summary
echo "Generate coverage summary"
go tool cover -func=coverage/coverage.out > coverage/coverage.txt

echo "Test coverage report generated in coverage/ directory"

# Exit with success
exit 0 