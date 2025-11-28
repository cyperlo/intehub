# 🚀 InteHub 快速启动指南

## 一键启动（Docker 方式 - 推荐）

```bash
# 1. 进入项目目录
cd /Users/chenhailong/code/project/intehub

# 2. 运行启动脚本
./start-docker.sh

# 或手动启动
docker-compose up -d --build
```

## 访问应用

- **前端**: http://localhost
- **后端**: http://localhost:8080
- **默认账户**: admin / admin123

## 常用命令速查

| 操作 | 命令 |
|------|------|
| 启动服务 | `docker-compose up -d` |
| 停止服务 | `docker-compose stop` |
| 查看日志 | `docker-compose logs -f` |
| 重启服务 | `docker-compose restart` |
| 完全停止 | `docker-compose down` |
| 查看状态 | `docker-compose ps` |

## 手动启动（开发模式）

### 后端

```bash
cd backend
go run main.go
# 访问 http://localhost:8080
```

### 前端

```bash
cd frontend
npm install
npm run dev
# 访问 http://localhost:5173
```

## 项目结构

```
intehub/
├── backend/              # Go 后端
│   ├── Dockerfile
│   ├── controllers/      # 控制器
│   ├── models/          # 数据模型
│   ├── routes/          # 路由
│   └── main.go
├── frontend/            # Vue 前端
│   ├── Dockerfile
│   ├── nginx.conf
│   ├── src/
│   │   ├── views/       # 页面组件
│   │   ├── api/         # API 接口
│   │   └── router/      # 路由配置
│   └── package.json
├── data/                # 数据持久化目录
│   └── intehub.db      # SQLite 数据库
├── docker-compose.yml   # Docker Compose 配置
└── DOCKER_README.md     # 详细文档
```

## 核心功能

### 1. 字段定义管理
- 创建可复用的字段模板
- 支持 7 种字段类型
- 灵活配置验证规则

### 2. 推送配置
- 关联自定义字段
- 动态模板变量替换
- 支持多种 HTTP 方法

### 3. 推送测试
- 根据字段自动生成表单
- 实时测试推送效果
- 查看推送历史

## 使用流程

```
1. 定义字段
   ↓
2. 创建推送配置
   ↓
3. 关联字段到配置
   ↓
4. 编写推送模板
   ↓
5. 测试推送
```

## 故障排查

### Docker 启动失败

```bash
# 查看日志
docker-compose logs

# 重新构建
docker-compose up -d --build

# 清理后重启
docker-compose down
docker-compose up -d --build
```

### 端口被占用

修改 `docker-compose.yml` 中的端口映射：

```yaml
ports:
  - "8081:8080"  # 后端
  - "8082:80"    # 前端
```

### 数据库权限问题

```bash
# 检查数据目录
ls -la ./data

# 修改权限（如需要）
chmod -R 755 ./data
```

## 技术栈

- **前端**: Vue 3 + TypeScript + Element Plus
- **后端**: Go + Gin + GORM
- **数据库**: SQLite
- **容器**: Docker + Docker Compose

## 详细文档

- **Docker 部署**: [DOCKER_README.md](DOCKER_README.md)
- **字段管理**: [字段管理功能说明.md](字段管理功能说明.md)

## 获取帮助

```bash
# 查看 Docker 服务状态
docker-compose ps

# 查看实时日志
docker-compose logs -f

# 进入容器调试
docker-compose exec backend sh
docker-compose exec frontend sh
```

## 下一步

1. ✅ 启动服务
2. 🔑 修改默认密码
3. 📝 创建字段定义
4. ⚙️ 配置推送服务
5. 🚀 开始使用！
