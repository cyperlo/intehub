.PHONY: help build run wire clean test

help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

wire: ## 生成 Wire 依赖注入代码
	@echo "Generating wire code..."
	@cd internal/server && wire
	@echo "Wire code generated successfully"

build: wire ## 编译项目
	@echo "Building..."
	@go build -o bin/intehub ./cmd/server.go
	@echo "Build complete: bin/intehub"

run: build ## 编译并运行
	@echo "Starting server..."
	@./bin/intehub

dev: ## 开发模式运行（不编译）
	@go run ./cmd/server.go

clean: ## 清理编译文件
	@echo "Cleaning..."
	@rm -rf bin/
	@echo "Clean complete"

test: ## 运行测试
	@go test -v ./...

fmt: ## 格式化代码
	@go fmt ./...

lint: ## 代码检查
	@golangci-lint run

deps: ## 安装依赖
	@go mod download
	@go mod tidy

ui-dev: ## 启动前端开发服务器
	@cd ui && npm run dev

ui-build: ## 构建前端
	@cd ui && npm run build

docker-build: ## 构建 Docker 镜像
	@docker build -t intehub:latest .

docker-run: ## 运行 Docker 容器
	@docker-compose up -d

docker-stop: ## 停止 Docker 容器
	@docker-compose down

.DEFAULT_GOAL := help
