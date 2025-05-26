.PHONY: build run test clean dev

# Go 관련 변수
BINARY_NAME=game-server
MAIN_PATH=./cmd/server
BUILD_DIR=./build

# 기본 타겟
all: build

# 빌드
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# 실행
run:
	@echo "Running $(BINARY_NAME)..."
	go run $(MAIN_PATH)

# 개발 모드 (hot reload)
dev:
	@echo "Starting development server..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# 테스트
test:
	@echo "Running tests..."
	go test -v ./...

# 테스트 커버리지
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 의존성 설치
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# 린트
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

# 포맷팅
fmt:
	@echo "Formatting code..."
	go fmt ./...

# 정리
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# 도움말
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  dev           - Run in development mode with hot reload"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  deps          - Install dependencies"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  clean         - Clean build artifacts"
	@echo "  help          - Show this help" 