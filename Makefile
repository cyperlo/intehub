.PHONY: help build up down restart logs ps clean build-backend build-frontend start stop status shell-backend shell-frontend backup test deploy

# 默认目标
.DEFAULT_GOAL := help

# 颜色定义
YELLOW := \033[1;33m
GREEN := \033[0;32m
RED := \033[0;31m
NC := \033[0m

# 项目配置
PROJECT_NAME := intehub
COMPOSE_FILE := docker-compose.yml
BACKUP_DIR := ./backups

##@ 帮助信息

help: ## 显示帮助信息
	@echo "$(GREEN)========================================$(NC)"
	@echo "$(GREEN)   InteHub 项目管理命令$(NC)"
	@echo "$(GREEN)========================================$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "使用方法:\n  make $(YELLOW)<target>$(NC)\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(GREEN)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ 构建命令

build: ## 构建所有服务
	@echo "$(YELLOW)正在构建所有服务...$(NC)"
	docker-compose -f $(COMPOSE_FILE) build
	@echo "$(GREEN)✓ 构建完成$(NC)"

build-nc: ## 构建所有服务（不使用缓存）
	@echo "$(YELLOW)正在重新构建所有服务（无缓存）...$(NC)"
	docker-compose -f $(COMPOSE_FILE) build --no-cache
	@echo "$(GREEN)✓ 构建完成$(NC)"

build-backend: ## 只构建后端
	@echo "$(YELLOW)正在构建后端...$(NC)"
	docker-compose -f $(COMPOSE_FILE) build backend
	@echo "$(GREEN)✓ 后端构建完成$(NC)"

build-frontend: ## 只构建前端
	@echo "$(YELLOW)正在构建前端...$(NC)"
	docker-compose -f $(COMPOSE_FILE) build frontend
	@echo "$(GREEN)✓ 前端构建完成$(NC)"

##@ 运行命令

up: ## 启动所有服务（后台运行）
	@echo "$(YELLOW)正在启动服务...$(NC)"
	docker-compose -f $(COMPOSE_FILE) up -d
	@echo "$(GREEN)✓ 服务已启动$(NC)"
	@make status

up-build: ## 构建并启动所有服务
	@echo "$(YELLOW)正在构建并启动服务...$(NC)"
	docker-compose -f $(COMPOSE_FILE) up -d --build
	@echo "$(GREEN)✓ 服务已启动$(NC)"
	@make status

start: ## 启动已存在的容器
	@echo "$(YELLOW)正在启动容器...$(NC)"
	docker-compose -f $(COMPOSE_FILE) start
	@echo "$(GREEN)✓ 容器已启动$(NC)"

dev: ## 启动服务（前台运行，查看日志）
	@echo "$(YELLOW)正在启动服务（前台模式）...$(NC)"
	docker-compose -f $(COMPOSE_FILE) up

##@ 停止命令

down: ## 停止并删除所有容器
	@echo "$(YELLOW)正在停止并删除容器...$(NC)"
	docker-compose -f $(COMPOSE_FILE) down
	@echo "$(GREEN)✓ 容器已停止并删除$(NC)"

stop: ## 停止容器（不删除）
	@echo "$(YELLOW)正在停止容器...$(NC)"
	docker-compose -f $(COMPOSE_FILE) stop
	@echo "$(GREEN)✓ 容器已停止$(NC)"

restart: ## 重启所有服务
	@echo "$(YELLOW)正在重启服务...$(NC)"
	docker-compose -f $(COMPOSE_FILE) restart
	@echo "$(GREEN)✓ 服务已重启$(NC)"

restart-backend: ## 重启后端服务
	@echo "$(YELLOW)正在重启后端...$(NC)"
	docker-compose -f $(COMPOSE_FILE) restart backend
	@echo "$(GREEN)✓ 后端已重启$(NC)"

restart-frontend: ## 重启前端服务
	@echo "$(YELLOW)正在重启前端...$(NC)"
	docker-compose -f $(COMPOSE_FILE) restart frontend
	@echo "$(GREEN)✓ 前端已重启$(NC)"

##@ 日志命令

logs: ## 查看所有服务日志
	docker-compose -f $(COMPOSE_FILE) logs -f

logs-backend: ## 查看后端日志
	docker-compose -f $(COMPOSE_FILE) logs -f backend

logs-frontend: ## 查看前端日志
	docker-compose -f $(COMPOSE_FILE) logs -f frontend

logs-tail: ## 查看最近100行日志
	docker-compose -f $(COMPOSE_FILE) logs --tail=100

##@ 状态命令

ps: ## 查看容器状态
	docker-compose -f $(COMPOSE_FILE) ps

status: ## 查看详细状态信息
	@echo "$(GREEN)========================================$(NC)"
	@echo "$(GREEN)   容器状态$(NC)"
	@echo "$(GREEN)========================================$(NC)"
	@docker-compose -f $(COMPOSE_FILE) ps
	@echo ""
	@echo "$(GREEN)访问地址:$(NC)"
	@echo "  前端: $(YELLOW)http://localhost$(NC)"
	@echo "  后端: $(YELLOW)http://localhost:8080$(NC)"
	@echo "  健康检查: $(YELLOW)http://localhost:8080/health$(NC)"
	@echo ""

stats: ## 查看资源使用情况
	docker stats --no-stream

##@ 调试命令

shell-backend: ## 进入后端容器
	docker-compose -f $(COMPOSE_FILE) exec backend sh

shell-frontend: ## 进入前端容器
	docker-compose -f $(COMPOSE_FILE) exec frontend sh

test-backend: ## 测试后端健康检查
	@echo "$(YELLOW)测试后端健康检查...$(NC)"
	@curl -s http://localhost:8080/health || echo "$(RED)✗ 后端未响应$(NC)"

test-frontend: ## 测试前端访问
	@echo "$(YELLOW)测试前端访问...$(NC)"
	@curl -sI http://localhost | head -n 1 || echo "$(RED)✗ 前端未响应$(NC)"

test: ## 测试所有服务
	@make test-backend
	@make test-frontend
	@echo "$(GREEN)✓ 测试完成$(NC)"

##@ 数据管理

backup: ## 备份数据库
	@echo "$(YELLOW)正在备份数据库...$(NC)"
	@mkdir -p $(BACKUP_DIR)
	@cp ./data/intehub.db $(BACKUP_DIR)/intehub_$(shell date +%Y%m%d_%H%M%S).db
	@echo "$(GREEN)✓ 备份完成: $(BACKUP_DIR)/intehub_$(shell date +%Y%m%d_%H%M%S).db$(NC)"

backup-list: ## 列出所有备份
	@echo "$(GREEN)备份文件列表:$(NC)"
	@ls -lh $(BACKUP_DIR)/ 2>/dev/null || echo "$(YELLOW)暂无备份$(NC)"

##@ 清理命令

clean: ## 清理未使用的资源
	@echo "$(YELLOW)正在清理未使用的资源...$(NC)"
	docker system prune -f
	@echo "$(GREEN)✓ 清理完成$(NC)"

clean-all: ## 清理所有资源（包括镜像）
	@echo "$(RED)警告: 这将删除所有未使用的镜像和容器$(NC)"
	@read -p "确定继续? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker system prune -a -f; \
		echo "$(GREEN)✓ 清理完成$(NC)"; \
	else \
		echo "$(YELLOW)已取消$(NC)"; \
	fi

clean-data: ## 删除数据库（危险操作）
	@echo "$(RED)警告: 这将删除所有数据库数据$(NC)"
	@read -p "确定继续? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		rm -rf ./data/*.db; \
		echo "$(GREEN)✓ 数据已删除$(NC)"; \
	else \
		echo "$(YELLOW)已取消$(NC)"; \
	fi

##@ 部署命令

deploy: ## 完整部署（停止、构建、启动）
	@echo "$(GREEN)========================================$(NC)"
	@echo "$(GREEN)   开始部署 $(PROJECT_NAME)$(NC)"
	@echo "$(GREEN)========================================$(NC)"
	@make down
	@make build
	@make up
	@sleep 5
	@make test
	@echo "$(GREEN)========================================$(NC)"
	@echo "$(GREEN)   部署完成！$(NC)"
	@echo "$(GREEN)========================================$(NC)"

redeploy: ## 快速重新部署（保留数据）
	@echo "$(YELLOW)正在重新部署...$(NC)"
	@make stop
	@make build
	@make start
	@echo "$(GREEN)✓ 重新部署完成$(NC)"

##@ 监控命令

watch-logs: ## 实时监控日志（彩色输出）
	docker-compose -f $(COMPOSE_FILE) logs -f --tail=50

watch-status: ## 持续监控容器状态
	watch -n 2 'docker-compose -f $(COMPOSE_FILE) ps'

##@ 开发命令

init: ## 初始化项目（首次使用）
	@echo "$(GREEN)========================================$(NC)"
	@echo "$(GREEN)   初始化 $(PROJECT_NAME) 项目$(NC)"
	@echo "$(GREEN)========================================$(NC)"
	@mkdir -p data backups
	@echo "$(YELLOW)1. 创建数据目录...$(NC)"
	@echo "$(GREEN)✓ 目录已创建$(NC)"
	@echo ""
	@echo "$(YELLOW)2. 构建镜像...$(NC)"
	@make build
	@echo ""
	@echo "$(YELLOW)3. 启动服务...$(NC)"
	@make up
	@echo ""
	@echo "$(GREEN)========================================$(NC)"
	@echo "$(GREEN)   初始化完成！$(NC)"
	@echo "$(GREEN)========================================$(NC)"
	@echo ""
	@echo "$(YELLOW)默认登录账户:$(NC)"
	@echo "  用户名: admin"
	@echo "  密码: admin123"
	@echo ""
	@echo "$(YELLOW)访问地址:$(NC)"
	@echo "  前端: http://localhost"
	@echo "  后端: http://localhost:8080"

update: ## 更新代码并重新部署
	@echo "$(YELLOW)正在更新代码...$(NC)"
	git pull
	@make redeploy

##@ 快捷命令

quick-start: up-build ## 快速启动（构建+启动）

quick-stop: down ## 快速停止

quick-restart: restart ## 快速重启
