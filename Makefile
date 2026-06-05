.PHONY: help build run test clean docker-build-app docker-build-docreader docker-build-frontend docker-build-all docker-run migrate-up migrate-down docker-restart docker-stop start stop logs status restart url stop-all start-all start-ollama stop-ollama build-images build-images-app build-images-docreader build-images-frontend clean-images check-env list-containers pull-images show-platform dev-start dev-stop dev-restart dev-logs dev-status dev-app dev-frontend docs install-swagger build-lite run-lite package-lite deploy

# Show help
help:
	@echo "WeKnora Makefile"
	@echo ""
	@echo "Basic Commands:"
	@echo "  build             Build the application"
	@echo "  run               Build and run the application"
	@echo "  test              Run tests"
	@echo "  clean             Clean build artifacts"
	@echo ""
	@echo "Docker Commands:"
	@echo "  docker-build-app        Build app Docker image (wechatopenai/weknora-app)"
	@echo "  docker-build-docreader  Build docreader Docker image (wechatopenai/weknora-docreader)"
	@echo "  docker-build-frontend   Build frontend Docker image (wechatopenai/weknora-ui)"
	@echo "  docker-build-all        Build all Docker images"
	@echo "  docker-run             Run Docker containers"
	@echo "  docker-stop            Stop Docker containers"
	@echo "  docker-restart         Restart Docker containers"
	@echo ""
	@echo "Local Docker Start/Stop:"
	@echo "  start              Start all Docker services (make start)"
	@echo "  stop               Stop all Docker services (make stop)"
	@echo "  restart            Restart all Docker services (make restart)"
	@echo "  logs               Tail logs from all containers (make logs)"
	@echo "  status             Show container status (make status)"
	@echo "  url                Show service URLs (make url)"
	@echo ""
	@echo "Service Management:"
	@echo "  start-all          Start all services (script)"
	@echo "  stop-all           Stop all services (script)"
	@echo "  start-ollama       Start Ollama only"
	@echo ""
	@echo "Image Building:"
	@echo "  build-images      Build all images from source"
	@echo "  build-images-app  Build app image from source"
	@echo "  build-images-docreader  Build docreader image from source"
	@echo "  build-images-frontend   Build frontend image from source"
	@echo "  clean-images      Clean local images"
	@echo ""
	@echo "Database:"
	@echo "  migrate-up        Run database migrations"
	@echo "  migrate-down      Rollback database migrations"
	@echo ""
	@echo "Development Tools:"
	@echo "  fmt               Format code"
	@echo "  lint              Lint code"
	@echo "  deps              Install dependencies"
	@echo "  docs              Generate Swagger API docs"
	@echo "  install-swagger   Install swag tool"
	@echo ""
	@echo "Environment:"
	@echo "  check-env         Check environment configuration"
	@echo "  list-containers   List running containers"
	@echo "  pull-images       Pull latest images"
	@echo "  show-platform     Show current build platform"
	@echo ""
	@echo "Deployment:"
	@echo "  deploy            Build and deploy only changed services to production"
	@echo ""
	@echo "Development Mode (recommended):"
	@echo "  dev-start         Start dev infrastructure (dependencies only)"
	@echo "  dev-stop          Stop dev environment"
	@echo "  dev-restart       Restart dev environment"
	@echo "  dev-logs          View dev environment logs"
	@echo "  dev-status        View dev environment status"
	@echo "  dev-app           Start backend app locally (requires dev-start first)"
	@echo "  dev-frontend      Start frontend locally (requires dev-start first)"
	@echo ""
	@echo "Lite Mode (zero external dependencies):"
	@echo "  build-lite        Build Lite version (build frontend to web/, then Go; SKIP_FRONTEND=1 skips frontend)"
	@echo "  run-lite          Build and run Lite version"
	@echo "  package-lite      Build and package Lite release tarball"
	@echo "  package-mac-app   Build and package macOS desktop app (.app)"

# Go related variables
BINARY_NAME=WeKnora
MAIN_PATH=./cmd/server

# Docker related variables
DOCKER_IMAGE=wechatopenai/weknora-app
DOCKER_TAG=latest

# Platform detection
ifeq ($(shell uname -m),x86_64)
    PLATFORM=linux/amd64
else ifeq ($(shell uname -m),aarch64)
    PLATFORM=linux/arm64
else ifeq ($(shell uname -m),arm64)
    PLATFORM=linux/arm64
else
    PLATFORM=linux/amd64
endif

# Build the application
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Run the application
run: build
	./$(BINARY_NAME)

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

# Build Docker image
docker-build-app:
	@echo "获取版本信息..."
	@eval $$(./scripts/get_version.sh env); \
	./scripts/get_version.sh info; \
	docker build --platform $(PLATFORM) \
		--build-arg VERSION_ARG="$$VERSION" \
		--build-arg COMMIT_ID_ARG="$$COMMIT_ID" \
		--build-arg BUILD_TIME_ARG="$$BUILD_TIME" \
		--build-arg GO_VERSION_ARG="$$GO_VERSION" \
		-f docker/Dockerfile.app -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Build docreader Docker image
docker-build-docreader:
	docker build --platform $(PLATFORM) -f docker/Dockerfile.docreader -t wechatopenai/weknora-docreader:latest .

# Build frontend Docker image
docker-build-frontend:
	docker build --platform $(PLATFORM) -f frontend/Dockerfile -t wechatopenai/weknora-ui:latest frontend/

# Build all Docker images
docker-build-all: docker-build-app docker-build-docreader docker-build-frontend

# Run Docker container (传统方式)
# Touch .env if missing — docker-compose.yml's `env_file: [.env]` is required
# for ${ENV} interpolation in builtin_models.yaml and would otherwise refuse
# to parse on fresh clones. `start-all` handles this via check_env_file; this
# direct path needs its own guard.
docker-run:
	@[ -f .env ] || ([ -f .env.example ] && cp .env.example .env || touch .env)
	docker-compose up

# 使用新脚本启动所有服务
start-all:
	./scripts/start_all.sh

# 使用新脚本仅启动Ollama服务
start-ollama:
	./scripts/start_all.sh --ollama

# 使用新脚本仅启动Docker容器
start-docker:
	./scripts/start_all.sh --docker

# 使用新脚本停止所有服务
stop-all:
	./scripts/start_all.sh --stop

# Stop Docker container (传统方式)
docker-stop:
	docker-compose down

# 从源码构建镜像相关命令
build-images:
	./scripts/build_images.sh

build-images-app:
	./scripts/build_images.sh --app

build-images-docreader:
	./scripts/build_images.sh --docreader

build-images-frontend:
	./scripts/build_images.sh --frontend

clean-images:
	./scripts/build_images.sh --clean

# Restart Docker container (stop, start)
docker-restart:
	@[ -f .env ] || ([ -f .env.example ] && cp .env.example .env || touch .env)
	docker-compose stop -t 60
	docker-compose up

# ========== Local Docker Start/Stop ==========

# Start all Docker services (docker compose up -d)
start:
	@[ -f .env ] || ([ -f .env.example ] && cp .env.example .env || touch .env)
	docker compose up -d

# Stop all Docker services (docker compose down)
stop:
	docker compose down

# Tail logs from all containers
logs:
	docker compose logs -f

# Show container status and resource usage
status:
	@echo "=== Container Status ==="
	docker compose ps
	@echo ""
	@echo "=== Resource Usage ==="
	@docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}" 2>/dev/null || echo "Docker is not running"

# Show service URLs (reads .env for port overrides)
url:
	@echo "=============================="
	@echo " WeKnora Service URLs"
	@echo "=============================="
	@echo " Web UI:       http://localhost:$${FRONTEND_PORT:-80}"
	@echo " Backend API:  http://localhost:$${APP_PORT:-8080}"
	@echo "=============================="

# Restart all services
restart: stop start

# Database migrations
migrate-up:
	./scripts/migrate.sh up

migrate-down:
	./scripts/migrate.sh down

migrate-version:
	./scripts/migrate.sh version

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: migration name is required"; \
		echo "Usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	./scripts/migrate.sh create $(name)

migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "Error: version is required"; \
		echo "Usage: make migrate-force version=4"; \
		exit 1; \
	fi
	./scripts/migrate.sh force $(version)

migrate-goto:
	@if [ -z "$(version)" ]; then \
		echo "Error: version is required"; \
		echo "Usage: make migrate-goto version=3"; \
		exit 1; \
	fi
	./scripts/migrate.sh goto $(version)

# Generate API documentation (Swagger)
docs:
	@echo "生成 Swagger API 文档..."
	swag init -g $(MAIN_PATH)/main.go -o ./docs --parseDependency --parseInternal
	@echo "文档已生成到 ./docs 目录"
	@echo "启动服务后访问 http://localhost:8080/swagger/index.html 查看文档"

# Install swagger tool
install-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod download

# Build for production
# google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn for qdrant milvus proto conflict
build-prod:
	VERSION=$$(git describe --tags --abbrev=0 2>/dev/null || echo "$${VERSION:-unknown}"); \
	COMMIT_ID=$${COMMIT_ID:-unknown}; \
	CGO_ENABLED=1 \
	CGO_CFLAGS="-Wno-deprecated-declarations" \
	CGO_LDFLAGS="$$(if [ "$$(uname)" = 'Darwin' ]; then echo '-Wl,-no_warn_duplicate_libraries'; fi)" \
	BUILD_TIME=$${BUILD_TIME:-unknown}; \
	GO_VERSION=$${GO_VERSION:-unknown}; \
	LDFLAGS="-X 'github.com/Tencent/WeKnora/internal/handler.Version=$$VERSION' -X 'github.com/Tencent/WeKnora/internal/handler.Edition=standard' -X 'github.com/Tencent/WeKnora/internal/handler.CommitID=$$COMMIT_ID' -X 'github.com/Tencent/WeKnora/internal/handler.BuildTime=$$BUILD_TIME' -X 'github.com/Tencent/WeKnora/internal/handler.GoVersion=$$GO_VERSION' -X 'google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn'"; \
	go build -ldflags="-w -s $$LDFLAGS" -o $(BINARY_NAME) $(MAIN_PATH)

# Build Lite version (single binary, SQLite + in-memory queue)
# 会先构建前端到 web/，再构建 Go 二进制；SKIP_FRONTEND=1 可跳过前端
build-lite:
	@if [ -f frontend/package.json ] && [ "$${SKIP_FRONTEND:-}" != "1" ]; then \
		echo ">> Building frontend for Lite..."; \
		(cd frontend && npm ci --prefer-offline && npm run build) && \
		rm -rf web && cp -r frontend/dist web; \
	elif [ "$${SKIP_FRONTEND:-}" = "1" ]; then \
		echo ">> Skipping frontend (SKIP_FRONTEND=1)"; \
	else \
		echo ">> No frontend/package.json, skipping frontend"; \
	fi
	export EDITION=lite; \
	eval "$$(./scripts/get_version.sh env)"; \
	LDFLAGS="$$(./scripts/get_version.sh ldflags) -X 'google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn'"; \
	CGO_ENABLED=1 \
	CGO_CFLAGS="-Wno-deprecated-declarations" \
	CGO_LDFLAGS="$$(if [ "$$(uname)" = 'Darwin' ]; then echo '-Wl,-no_warn_duplicate_libraries'; fi)" \
	go build -tags "sqlite_fts5" -ldflags="-w -s $$LDFLAGS" -o $(BINARY_NAME)-lite $(MAIN_PATH)

# Run Lite version with .env.lite defaults
run-lite: build-lite
	@if [ ! -f .env.lite ]; then echo "Error: .env.lite not found"; exit 1; fi
	@set -a && . ./.env.lite && set +a && ./$(BINARY_NAME)-lite

# Deploy only changed services to production via SSH.
# Detects which files changed since the last deploy and only builds
# the affected Docker Compose services (app, frontend, docreader).
deploy:
	@echo ">> Deploying to production (Hetzner)..."
	@ssh elmarhepp "cd /var/www/weknora && \
		echo '>> Pulling latest code...' && \
		git pull && \
		echo '>> Detecting changed services...' && \
		CHANGED_FILES=\$$(git diff HEAD~1 --name-only 2>/dev/null || git diff --name-only) && \
		SERVICES='' && \
		if echo \"\$$CHANGED_FILES\" | grep -qE '^frontend/'; then SERVICES=\"\$$SERVICES frontend\"; fi && \
		if echo \"\$$CHANGED_FILES\" | grep -qE '^(internal/|cmd/|go\.mod|go\.sum|docker/Dockerfile\.app)'; then SERVICES=\"\$$SERVICES app\"; fi && \
		if echo \"\$$CHANGED_FILES\" | grep -qE '^docreader/|docker/Dockerfile\.docreader'; then SERVICES=\"\$$SERVICES docreader\"; fi && \
		if [ -z \"\$$SERVICES\" ]; then \
			echo '>> No changed services detected, restarting all...'; \
			docker compose up -d --build; \
		else \
			echo '>> Building and restarting:' \$$SERVICES; \
			docker compose up -d --build \$$SERVICES; \
		fi && \
		echo '>> Deployment complete.'"

# Package Lite version into distributable tarball
package-lite:
	./scripts/package-lite.sh

# Package Mac App
package-mac-app:
	./scripts/package-mac-app.sh

download_spatial:
	go run cmd/download/duckdb/duckdb.go

clean-db:
	@echo "Cleaning database..."
	@if [ $$(docker volume ls -q -f name=weknora_postgres-data) ]; then \
		docker volume rm weknora_postgres-data; \
	fi
	@if [ $$(docker volume ls -q -f name=weknora_minio_data) ]; then \
		docker volume rm weknora_minio_data; \
	fi
	@if [ $$(docker volume ls -q -f name=weknora_redis_data) ]; then \
		docker volume rm weknora_redis_data; \
	fi

# Environment check
check-env:
	./scripts/start_all.sh --check

# List containers
list-containers:
	./scripts/start_all.sh --list

# Pull latest images
pull-images:
	./scripts/start_all.sh --pull

# Show current platform
show-platform:
	@echo "当前系统架构: $(shell uname -m)"
	@echo "Docker构建平台: $(PLATFORM)"

# Development mode commands
dev-start:
	./scripts/dev.sh start

dev-stop:
	./scripts/dev.sh stop

dev-restart:
	./scripts/dev.sh restart

dev-logs:
	./scripts/dev.sh logs

dev-status:
	./scripts/dev.sh status

dev-app:
	./scripts/dev.sh app

dev-frontend:
	./scripts/dev.sh frontend


