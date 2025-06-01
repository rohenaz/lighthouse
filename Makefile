.PHONY: all build test clean install proto cli web

all: proto build

# Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	@mkdir -p core/proto
	@protoc --proto_path=proto \
		--go_out=. --go_opt=paths=source_relative \
		proto/lighthouse.proto

# Build the CLI
cli: proto
	@echo "Building CLI..."
	@go build -o bin/lighthouse cmd/lighthouse/*.go

# Build everything
build: cli

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Install CLI globally
install: cli
	@echo "Installing lighthouse CLI..."
	@cp bin/lighthouse $(GOPATH)/bin/lighthouse

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf core/proto/*.pb.go

# Run the web dev server
web:
	@echo "Starting web development server..."
	@cd web && bun dev

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@cd web && bun run format

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run
	@cd web && bun run lint

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@cd web && bun install