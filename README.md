# InteHub - 集成推送平台

一个前后端分离的集成推送平台，支持多种推送方式和灵活配置。

## 技术栈

### 后端
- Go 1.21+
- Gin Web框架
- GORM（数据库ORM）
- JWT认证
- SQLite/MySQL

### 前端
- Vue 3
- TypeScript
- Vue Router
- Pinia（状态管理）
- Element Plus（UI组件库）
- Axios（HTTP客户端）

## 功能特性

- ✅ 用户认证（登录/登出）
- ✅ 推送平台管理
  - 配置推送URL
  - 自定义推送内容
  - 支持多种推送方式（HTTP POST、GET、Webhook等）
  - 推送历史记录
- ✅ 安全的JWT令牌认证

## 项目结构

```
intehub/
├── backend/          # Go后端服务
│   ├── main.go
│   ├── config/       # 配置文件
│   ├── models/       # 数据模型
│   ├── controllers/  # 控制器
│   ├── middleware/   # 中间件
│   ├── routes/       # 路由
│   └── utils/        # 工具函数
├── frontend/         # Vue3前端应用
│   ├── src/
│   │   ├── api/      # API接口
│   │   ├── components/
│   │   ├── views/    # 页面
│   │   ├── router/   # 路由
│   │   ├── stores/   # Pinia状态
│   │   └── utils/    # 工具函数
│   └── package.json
└── README.md
```

## 快速开始

### 后端启动

```bash
cd backend
go mod download
go run main.go
```

后端服务默认运行在 `http://localhost:8080`

### 前端启动

```bash
cd frontend
npm install
npm run dev
```

前端应用默认运行在 `http://localhost:5173`

## API文档

### 认证相关
- POST `/api/auth/login` - 用户登录
- POST `/api/auth/logout` - 用户登出

### 推送平台
- GET `/api/push/configs` - 获取推送配置列表
- POST `/api/push/configs` - 创建推送配置
- PUT `/api/push/configs/:id` - 更新推送配置
- DELETE `/api/push/configs/:id` - 删除推送配置
- POST `/api/push/send` - 执行推送
- GET `/api/push/history` - 获取推送历史

## 默认账户

- 用户名: `admin`
- 密码: `admin123`

## License

MIT
