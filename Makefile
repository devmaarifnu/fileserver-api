.PHONY: run build clean test install dev

# Run the application
run:
	go run cmd/api/main.go

# Build the application
build:
	go build -o bin/cdn-fileserver cmd/api/main.go

# Build for Linux (production)
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/cdn-fileserver cmd/api/main.go

# Install dependencies
install:
	go mod download

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf logs/*.log

# Run tests
test:
	go test -v ./...

# Run with hot reload (install air first: go install github.com/cosmtrek/air@latest)
dev:
	air

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run
