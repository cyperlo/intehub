#!/bin/bash

# InteHub 一键部署脚本
# 适用于小内存服务器（3.6GB）

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VERSION=$(date +%Y%m%d%H%M%S)

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 显示帮助信息
show_help() {
    cat << EOF
InteHub 部署脚本

用法: $0 [选项]

选项:
    -h, --help              显示帮助信息
    -s, --skip-build        跳过构建，直接部署
    -c, --clean             清理所有容器和镜像
    -r, --rollback          回滚到上一个版本
    -l, --logs              查看容器日志
    --no-check              跳过健康检查

环境变量（必需）:
    INTEHUB_POSTGRESQL_URI  数据库连接字符串
    INTEHUB_JWT_SECRET      JWT 密钥
    INTEHUB_CRYPTO_KEY      加密密钥

示例:
    # 完整部署
    export INTEHUB_POSTGRESQL_URI="postgresql://user:pass@host:5432/db"
    export INTEHUB_JWT_SECRET="your-secret"
    export INTEHUB_CRYPTO_KEY="your-key"
    $0

    # 跳过构建直接部署
    $0 --skip-build

    # 查看日志
    $0 --logs

    # 清理环境
    $0 --clean
EOF
}

# 检查系统资源
check_system() {
    log_step "检查系统资源..."
    
    # 检查内存
    total_mem=$(free -m | awk 'NR==2{print $2}')
    available_mem=$(free -m | awk 'NR==2{print $7}')
    
    log_info "总内存: ${total_mem}MB, 可用内存: ${available_mem}MB"
    
    if [ "$available_mem" -lt 500 ]; then
        log_warn "可用内存不足 500MB，构建可能失败"
        log_warn "建议先停止非必要服务或增加 swap"
    fi
    
    # 检查磁盘空间
    available_disk=$(df -BG "$SCRIPT_DIR" | awk 'NR==2{print $4}' | sed 's/G//')
    log_info "可用磁盘空间: ${available_disk}GB"
    
    if [ "$available_disk" -lt 5 ]; then
        log_error "磁盘空间不足 5GB"
        exit 1
    fi
}

# 检查必要的环境变量
check_env() {
    log_step "检查环境变量..."
    
    if [ -z "$INTEHUB_POSTGRESQL_URI" ]; then
        log_error "INTEHUB_POSTGRESQL_URI 未设置"
        log_info "请设置: export INTEHUB_POSTGRESQL_URI='postgresql://user:pass@host:5432/db'"
        exit 1
    fi
    
    if [ -z "$INTEHUB_JWT_SECRET" ]; then
        log_error "INTEHUB_JWT_SECRET 未设置"
        log_info "请设置: export INTEHUB_JWT_SECRET='your-secret'"
        exit 1
    fi
    
    if [ -z "$INTEHUB_CRYPTO_KEY" ]; then
        log_error "INTEHUB_CRYPTO_KEY 未设置"
        log_info "请设置: export INTEHUB_CRYPTO_KEY='your-key'"
        exit 1
    fi
    
    log_info "✓ 环境变量检查通过"
}

# 清理旧镜像
clean_old_images() {
    log_step "清理旧镜像..."
    
    # 清理后端旧镜像（保留最新3个版本 + latest + backup）
    log_info "清理后端旧镜像..."
    docker images intehub-backend --format "{{.Tag}}" | \
        grep -E '^[0-9]+$' | \
        sort -rn | \
        tail -n +4 | \
        xargs -r -I {} docker rmi intehub-backend:{} 2>/dev/null || true
    
    # 清理前端旧镜像（保留最新3个版本 + latest + backup）
    log_info "清理前端旧镜像..."
    docker images intehub-frontend --format "{{.Tag}}" | \
        grep -E '^[0-9]+$' | \
        sort -rn | \
        tail -n +4 | \
        xargs -r -I {} docker rmi intehub-frontend:{} 2>/dev/null || true
    
    # 清理悬空镜像
    docker image prune -f 2>/dev/null || true
    
    # 显示剩余镜像
    log_info "当前镜像列表："
    docker images | grep -E "REPOSITORY|intehub"
    
    log_info "✓ 旧镜像清理完成"
}

# 备份当前版本
backup_current_version() {
    log_step "备份当前版本..."
    
    if docker images | grep -q "intehub-backend:latest"; then
        docker tag intehub-backend:latest intehub-backend:backup 2>/dev/null || true
        log_info "✓ 后端镜像已备份"
    fi
    
    if docker images | grep -q "intehub-frontend:latest"; then
        docker tag intehub-frontend:latest intehub-frontend:backup 2>/dev/null || true
        log_info "✓ 前端镜像已备份"
    fi
}

# 构建后端镜像
build_backend() {
    log_step "构建后端镜像..."
    
    cd "$SCRIPT_DIR"
    
    docker build \
        --memory="1g" \
        --memory-swap="1.5g" \
        -t intehub-backend:${VERSION} \
        -t intehub-backend:latest \
        -f Dockerfile.backend \
        .
    
    # 立即清理中间层
    docker image prune -f
    
    log_info "✓ 后端镜像构建完成"
}

# 构建前端镜像
build_frontend() {
    log_step "构建前端镜像..."
    
    cd "$SCRIPT_DIR"
    
    # 限制内存避免 OOM
    docker build \
        --memory="1g" \
        --memory-swap="1.5g" \
        -t intehub-frontend:${VERSION} \
        -t intehub-frontend:latest \
        -f Dockerfile.frontend \
        .
    
    # 立即清理中间层
    docker image prune -f
    
    log_info "✓ 前端镜像构建完成"
}

# 停止旧容器
stop_old_containers() {
    log_step "停止旧容器..."
    
    if docker ps -a | grep -q intehub-backend; then
        log_info "停止后端容器..."
        docker stop -t 30 intehub-backend 2>/dev/null || true
        docker rm -f intehub-backend 2>/dev/null || true
    fi
    
    if docker ps -a | grep -q intehub-frontend; then
        log_info "停止前端容器..."
        docker stop -t 30 intehub-frontend 2>/dev/null || true
        docker rm -f intehub-frontend 2>/dev/null || true
    fi
    
    log_info "✓ 旧容器已停止"
}

# 创建必要的目录
create_directories() {
    log_step "创建数据目录..."
    
    sudo mkdir -p /var/intehub/data
    sudo mkdir -p /var/intehub/logs
    sudo chmod -R 755 /var/intehub
    
    log_info "✓ 数据目录创建完成"
}

# 部署后端容器
deploy_backend() {
    log_step "部署后端容器..."
    
    # 创建网络
    docker network create intehub-network 2>/dev/null || true
    
    docker run -d \
        --name intehub-backend \
        --network intehub-network \
        -e INTEHUB_POSTGRESQL_URI="${INTEHUB_POSTGRESQL_URI}" \
        -e INTEHUB_SERVER_PORT=8080 \
        -e INTEHUB_JWT_SECRET="${INTEHUB_JWT_SECRET}" \
        -e INTEHUB_CRYPTO_KEY="${INTEHUB_CRYPTO_KEY}" \
        -p 8080:8080 \
        -v /var/intehub/data:/app/data \
        -v /var/intehub/logs:/app/logs \
        --memory="512m" \
        --memory-swap="768m" \
        --restart unless-stopped \
        --log-driver json-file \
        --log-opt max-size=10m \
        --log-opt max-file=3 \
        intehub-backend:latest
    
    log_info "✓ 后端容器启动成功"
}

# 部署前端容器
deploy_frontend() {
    log_step "部署前端容器..."
    
    docker run -d \
        --name intehub-frontend \
        --network intehub-network \
        -p 801:80 \
        --memory="128m" \
        --memory-swap="256m" \
        --restart unless-stopped \
        --log-driver json-file \
        --log-opt max-size=10m \
        --log-opt max-file=3 \
        intehub-frontend:latest
    
    log_info "✓ 前端容器启动成功"
}

# 健康检查
health_check() {
    log_step "执行健康检查..."
    
    # 等待服务启动
    log_info "等待服务启动..."
    sleep 10
    
    # 检查容器状态
    if ! docker ps | grep -q intehub-backend; then
        log_error "后端容器未运行"
        log_error "容器日志："
        docker logs --tail 50 intehub-backend
        return 1
    fi
    
    if ! docker ps | grep -q intehub-frontend; then
        log_error "前端容器未运行"
        log_error "容器日志："
        docker logs --tail 50 intehub-frontend
        return 1
    fi
    
    # 检查后端健康接口
    log_info "检查后端健康状态..."
    for i in {1..10}; do
        if curl -f -s http://localhost:8080/api/health > /dev/null 2>&1; then
            log_info "✓ 后端健康检查通过"
            break
        fi
        
        if [ $i -eq 10 ]; then
            log_error "后端健康检查失败"
            docker logs --tail 50 intehub-backend
            return 1
        fi
        
        log_warn "后端未就绪，等待重试... ($i/10)"
        sleep 3
    done
    
    # 检查前端
    log_info "检查前端健康状态..."
    if curl -f -s -o /dev/null http://localhost:801 2>&1; then
        log_info "✓ 前端健康检查通过"
    else
        log_warn "前端健康检查失败，但容器正在运行"
        docker logs --tail 20 intehub-frontend
    fi
    
    log_info "✅ 健康检查完成"
}

# 显示服务信息
show_info() {
    echo ""
    echo "======================================"
    log_info "🎉 部署成功！"
    echo "======================================"
    echo ""
    echo "服务地址:"
    echo "  后端: http://localhost:8080"
    echo "  前端: http://localhost:801"
    echo ""
    echo "常用命令:"
    echo "  查看后端日志: docker logs -f intehub-backend"
    echo "  查看前端日志: docker logs -f intehub-frontend"
    echo "  查看容器状态: docker ps | grep intehub"
    echo "  查看资源使用: docker stats --no-stream"
    echo ""
    echo "  停止服务: docker stop intehub-backend intehub-frontend"
    echo "  启动服务: docker start intehub-backend intehub-frontend"
    echo "  重启服务: docker restart intehub-backend intehub-frontend"
    echo ""
    echo "  回滚版本: $0 --rollback"
    echo "  查看日志: $0 --logs"
    echo "======================================"
    echo ""
}

# 查看日志
show_logs() {
    log_step "查看容器日志..."
    
    echo ""
    echo "=== 后端日志（最后 30 行）==="
    docker logs --tail 30 intehub-backend 2>/dev/null || log_warn "后端容器不存在"
    
    echo ""
    echo "=== 前端日志（最后 30 行）==="
    docker logs --tail 30 intehub-frontend 2>/dev/null || log_warn "前端容器不存在"
    
    echo ""
    echo "=== 容器状态 ==="
    docker ps -a | grep intehub || log_warn "没有找到 intehub 容器"
    
    echo ""
    echo "=== 资源使用 ==="
    docker stats --no-stream | grep -E "CONTAINER|intehub" || true
}

# 回滚到备份版本
rollback() {
    log_step "回滚到备份版本..."
    
    if ! docker images | grep -q "intehub-backend:backup"; then
        log_error "没有找到备份镜像"
        exit 1
    fi
    
    # 停止当前容器
    stop_old_containers
    
    # 使用备份镜像
    docker tag intehub-backend:backup intehub-backend:latest
    docker tag intehub-frontend:backup intehub-frontend:latest
    
    # 重新部署
    deploy_backend
    sleep 5
    deploy_frontend
    
    # 健康检查
    health_check
    
    log_info "✅ 回滚完成"
}

# 清理所有资源
clean_all() {
    log_step "清理所有资源..."
    
    read -p "确定要清理所有 InteHub 容器和镜像吗？(y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "取消清理"
        exit 0
    fi
    
    # 停止并删除容器
    docker stop intehub-backend intehub-frontend 2>/dev/null || true
    docker rm -f intehub-backend intehub-frontend 2>/dev/null || true
    
    # 删除镜像
    docker rmi -f $(docker images | grep intehub | awk '{print $3}') 2>/dev/null || true
    
    # 删除网络
    docker network rm intehub-network 2>/dev/null || true
    
    log_info "✓ 清理完成"
}

# 主函数
main() {
    local skip_build=false
    local no_check=false
    
    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -s|--skip-build)
                skip_build=true
                shift
                ;;
            -c|--clean)
                clean_all
                exit 0
                ;;
            -r|--rollback)
                rollback
                exit 0
                ;;
            -l|--logs)
                show_logs
                exit 0
                ;;
            --no-check)
                no_check=true
                shift
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    log_info "开始部署 InteHub..."
    echo ""
    
    # 执行部署流程
    check_system
    check_env
    
    if [ "$skip_build" = false ]; then
        clean_old_images
        backup_current_version
        build_backend
        build_frontend
    else
        log_warn "跳过构建步骤"
    fi
    
    stop_old_containers
    create_directories
    deploy_backend
    sleep 5
    deploy_frontend
    
    if [ "$no_check" = false ]; then
        if ! health_check; then
            log_error "健康检查失败，开始回滚..."
            rollback
            exit 1
        fi
    fi
    
    # 部署成功后再次清理
    if [ "$skip_build" = false ]; then
        clean_old_images
    fi
    
    show_info
}

# 执行主函数
main "$@"
