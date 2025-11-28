# 部署指南

## 开发环境部署

### 后端部署

1. 进入后端目录
```bash
cd backend
```

2. 安装Go依赖
```bash
go mod tidy
```

3. 启动后端服务
```bash
go run main.go
```

后端服务将运行在 `http://localhost:8080`

### 前端部署

1. 进入前端目录
```bash
cd frontend
```

2. 安装依赖
```bash
npm install
```

3. 启动开发服务器
```bash
npm run dev
```

前端应用将运行在 `http://localhost:5173`

## 生产环境部署

### 后端

1. 编译后端
```bash
cd backend
go build -o intehub main.go
```

2. 运行编译后的程序
```bash
./intehub
```

### 前端

1. 构建前端
```bash
cd frontend
npm run build
```

2. 将 `dist` 目录部署到Web服务器（如Nginx）

### 使用Docker部署

#### 后端Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY backend/ .
RUN go mod download
RUN go build -o intehub main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/intehub .
EXPOSE 8080
CMD ["./intehub"]
```

#### 前端Dockerfile
```dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

#### Docker Compose
```yaml
version: '3.8'

services:
  backend:
    build:
      context: .
      dockerfile: Dockerfile.backend
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    environment:
      - JWT_SECRET=your-secret-key
    restart: unless-stopped

  frontend:
    build:
      context: .
      dockerfile: Dockerfile.frontend
    ports:
      - "80:80"
    depends_on:
      - backend
    restart: unless-stopped
```

## 数据库配置

### SQLite（默认）
无需额外配置，数据库文件将自动创建在 `backend/intehub.db`

### MySQL
修改 `backend/config/database.go`:
```go
import "gorm.io/driver/mysql"

dsn := "user:password@tcp(127.0.0.1:3306)/intehub?charset=utf8mb4&parseTime=True&loc=Local"
DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
```

### PostgreSQL
修改 `backend/config/database.go`:
```go
import "gorm.io/driver/postgres"

dsn := "host=localhost user=gorm password=gorm dbname=intehub port=5432 sslmode=disable"
DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

## 环境变量配置

后端支持以下环境变量：
- `JWT_SECRET`: JWT密钥
- `SERVER_PORT`: 服务器端口（默认8080）
- `DB_TYPE`: 数据库类型（sqlite/mysql/postgres）

## 安全建议

1. **修改默认密码**: 首次登录后立即修改admin账户密码
2. **更换JWT密钥**: 在生产环境中使用强随机密钥
3. **启用HTTPS**: 使用反向代理（如Nginx）配置SSL证书
4. **限制CORS**: 根据实际域名配置CORS策略
5. **数据备份**: 定期备份数据库文件

## 常见问题

### 端口被占用
修改后端 `main.go` 中的端口号，或前端 `vite.config.ts` 中的端口配置

### 跨域问题
确保后端CORS配置包含前端域名

### 数据库连接失败
检查数据库配置和权限设置
