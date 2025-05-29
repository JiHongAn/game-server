.PHONY: build run test clean dev run-dev run-test run-prod docker-up docker-down docker-logs

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

# 실행 (기본: development)
run:
	@echo "Running $(BINARY_NAME) in development mode..."
	ENV=development go run $(MAIN_PATH)

# 환경별 실행
run-dev:
	@echo "Running $(BINARY_NAME) in development mode..."
	ENV=development go run $(MAIN_PATH)

run-test:
	@echo "Running $(BINARY_NAME) in test mode..."
	ENV=test go run $(MAIN_PATH)

run-prod:
	@echo "Running $(BINARY_NAME) in production mode..."
	ENV=production go run $(MAIN_PATH)

# 개발 모드 (hot reload)
dev:
	@echo "Starting development server with hot reload..."
	@ENV=development; \
	GOPATH=$$(go env GOPATH); \
	if [ -f "$$GOPATH/bin/air" ]; then \
		$$GOPATH/bin/air; \
	else \
		echo "air not found. Installing..."; \
		go install github.com/air-verse/air@latest; \
		$$GOPATH/bin/air; \
	fi

# 테스트
test:
	@echo "Running tests..."
	ENV=test go test -v ./...

# 테스트 커버리지
test-coverage:
	@echo "Running tests with coverage..."
	ENV=test go test -v -coverprofile=coverage.out ./...
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

# Docker 관련 명령어 (개발용만)
docker-up:
	@echo "Starting development services (MySQL, Redis)..."
	docker-compose up -d

docker-down:
	@echo "Stopping development services..."
	docker-compose down

docker-logs:
	@echo "Showing logs for development services..."
	docker-compose logs -f

docker-clean:
	@echo "Cleaning up development services and volumes..."
	docker-compose down -v
	docker system prune -f

# 정리
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# 도움말
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application (development)"
	@echo "  run-dev       - Run in development environment"
	@echo "  run-test      - Run in test environment"
	@echo "  run-prod      - Run in production environment"
	@echo "  dev           - Run in development mode with hot reload"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  deps          - Install dependencies"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-up     - Start development services (DB, Redis)"
	@echo "  docker-down   - Stop development services"
	@echo "  docker-logs   - Show development services logs"
	@echo "  docker-clean  - Clean up development services and volumes"
	@echo "  help          - Show this help"