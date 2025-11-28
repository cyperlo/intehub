# InteHub Docker 部署指南

## 快速开始

### 前置要求

- Docker 20.10+
- Docker Compose 2.0+

### 一键启动

```bash
cd /Users/chenhailong/code/project/intehub
docker-compose up -d
```

### 访问应用

- **前端地址**: http://localhost
- **后端 API**: http://localhost:8080
- **健康检查**: http://localhost:8080/health

### 默认账户

```
用户名: admin
密码: admin123
```

## 详细说明

### 服务架构

```
┌──────────────┐
│   浏览器      │
└──────┬───────┘
       │ :80
┌──────▼───────────────┐
│  Frontend (Nginx)    │
│  - 静态文件服务       │
│  - API 反向代理       │
└──────┬───────────────┘
       │ :8080
┌──────▼───────────────┐
│  Backend (Go)        │
│  - REST API          │
│  - 业务逻辑          │
└──────┬───────────────┘
       │
┌──────▼───────────────┐
│  SQLite Database     │
│  (持久化到宿主机)     │
└──────────────────────┘
```

### 容器说明

#### Frontend 容器
- **镜像**: 基于 nginx:alpine
- **端口**: 80
- **功能**: 
  - 提供静态文件服务
  - 反向代理 API 请求到后端
  - 支持前端路由（SPA）

#### Backend 容器
- **镜像**: 基于 golang:1.21-alpine
- **端口**: 8080
- **功能**:
  - REST API 服务
  - 业务逻辑处理
  - 数据库管理

### 数据持久化

数据库文件存储在 `./data/intehub.db`，映射到宿主机，确保数据不会因容器重启而丢失。

```bash
# 数据目录结构
intehub/
├── data/
│   └── intehub.db          # SQLite 数据库文件
├── docker-compose.yml
├── backend/
│   └── Dockerfile
└── frontend/
    └── Dockerfile
```

## 常用命令

### 启动服务

```bash
# 前台运行（查看日志）
docker-compose up

# 后台运行
docker-compose up -d

# 重新构建并启动
docker-compose up -d --build
```

### 停止服务

```bash
# 停止容器
docker-compose stop

# 停止并删除容器
docker-compose down

# 停止并删除容器和数据卷（⚠️ 会删除数据）
docker-compose down -v
```

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs backend
docker-compose logs frontend

# 实时查看日志
docker-compose logs -f

# 查看最近 100 行日志
docker-compose logs --tail=100
```

### 重启服务

```bash
# 重启所有服务
docker-compose restart

# 重启特定服务
docker-compose restart backend
docker-compose restart frontend
```

### 进入容器

```bash
# 进入后端容器
docker-compose exec backend sh

# 进入前端容器
docker-compose exec frontend sh
```

### 查看状态

```bash
# 查看容器状态
docker-compose ps

# 查看资源使用情况
docker stats
```

## 环境变量配置

可以创建 `.env` 文件来自定义配置：

```bash
# .env 文件示例
DATABASE_PATH=/app/data/intehub.db
BACKEND_PORT=8080
FRONTEND_PORT=80
```

然后在 `docker-compose.yml` 中引用：

```yaml
environment:
  - DATABASE_PATH=${DATABASE_PATH}
```

## 数据备份与恢复

### 备份数据库

```bash
# 备份数据库文件
cp ./data/intehub.db ./data/intehub.db.backup.$(date +%Y%m%d_%H%M%S)

# 或使用 Docker 命令
docker-compose exec backend cp /app/data/intehub.db /app/data/backup.db
```

### 恢复数据库

```bash
# 停止服务
docker-compose stop

# 恢复数据库文件
cp ./data/intehub.db.backup.20251128_153000 ./data/intehub.db

# 启动服务
docker-compose start
```

## 生产环境部署建议

### 1. 使用环境变量

创建 `.env` 文件，不要提交到版本控制：

```bash
# JWT 密钥（生产环境必须修改）
JWT_SECRET=your-super-secret-key-change-in-production

# 数据库路径
DATABASE_PATH=/app/data/intehub.db
```

### 2. 使用 HTTPS

配置 Nginx SSL 证书：

```nginx
server {
    listen 443 ssl;
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    # ... 其他配置
}
```

### 3. 资源限制

在 `docker-compose.yml` 中添加资源限制：

```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
```

### 4. 健康检查

已配置健康检查，可以通过以下方式查看：

```bash
docker-compose ps
```

### 5. 日志管理

配置日志驱动：

```yaml
services:
  backend:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## 故障排查

### 容器无法启动

```bash
# 查看详细日志
docker-compose logs backend
docker-compose logs frontend

# 检查容器状态
docker-compose ps

# 重新构建
docker-compose up -d --build
```

### 数据库连接失败

```bash
# 检查数据目录权限
ls -la ./data

# 检查环境变量
docker-compose exec backend env | grep DATABASE
```

### 前端无法访问后端

```bash
# 检查网络连接
docker-compose exec frontend ping backend

# 检查后端健康状态
curl http://localhost:8080/health
```

### 端口冲突

如果端口 80 或 8080 已被占用，修改 `docker-compose.yml`：

```yaml
ports:
  - "8081:8080"  # 后端改为 8081
  - "8082:80"    # 前端改为 8082
```

## 性能优化

### 1. 使用构建缓存

```bash
# 使用缓存快速构建
docker-compose build --parallel
```

### 2. 减小镜像体积

已使用 Alpine Linux 作为基础镜像，并采用多阶段构建。

### 3. 启用 Nginx 缓存

Nginx 配置中已启用静态资源缓存。

## 安全注意事项

⚠️ **Docker 镜像漏洞警告**：IDE 可能会报告基础镜像（node:18-alpine, nginx:alpine）存在安全漏洞。这些是已知问题：
- **开发环境**: 可以忽略，不影响功能
- **生产环境**: 建议使用特定版本标签，并定期更新镜像

### 生产环境安全建议

1. **修改默认密码**: 首次部署后立即修改 admin 密码
2. **更新 JWT 密钥**: 在环境变量中设置强密钥
3. **限制端口访问**: 使用防火墙规则限制访问
4. **定期更新**: 及时更新 Docker 镜像和依赖
5. **数据加密**: 考虑对敏感数据进行加密存储

## 更新升级

### 更新代码

```bash
# 1. 拉取最新代码
git pull

# 2. 重新构建并启动
docker-compose up -d --build

# 3. 查看日志确认启动成功
docker-compose logs -f
```

### 回滚版本

```bash
# 1. 切换到旧版本
git checkout <commit-hash>

# 2. 重新构建
docker-compose up -d --build
```

## 卸载

```bash
# 停止并删除所有容器和网络
docker-compose down

# 删除数据（可选）
rm -rf ./data

# 删除镜像（可选）
docker rmi intehub-backend intehub-frontend
```

## 技术栈

- **前端**: Vue 3 + TypeScript + Element Plus + Vite
- **后端**: Go + Gin + GORM
- **数据库**: SQLite
- **容器**: Docker + Docker Compose
- **Web服务器**: Nginx

## 支持

如有问题，请查看：
1. 容器日志：`docker-compose logs`
2. 健康检查：`http://localhost:8080/health`
3. 数据库连接：检查 `./data` 目录权限

## 许可证

根据项目许可证使用。
