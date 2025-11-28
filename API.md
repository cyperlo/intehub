# API 文档

## 基础信息

- **Base URL**: `http://localhost:8080/api`
- **认证方式**: JWT Bearer Token
- **Content-Type**: `application/json`

## 认证相关

### 登录

**POST** `/auth/login`

请求体:
```json
{
  "username": "admin",
  "password": "admin123"
}
```

响应:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_info": {
    "id": 1,
    "username": "admin",
    "nickname": "管理员",
    "role": "admin",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 登出

**POST** `/auth/logout`

Headers:
```
Authorization: Bearer {token}
```

响应:
```json
{
  "message": "登出成功"
}
```

## 用户相关

### 获取当前用户信息

**GET** `/user/current`

Headers:
```
Authorization: Bearer {token}
```

响应:
```json
{
  "id": 1,
  "username": "admin",
  "nickname": "管理员",
  "role": "admin",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## 推送配置管理

### 获取推送配置列表

**GET** `/push/configs`

Headers:
```
Authorization: Bearer {token}
```

响应:
```json
[
  {
    "id": 1,
    "name": "Webhook推送",
    "description": "推送到Webhook服务",
    "url": "https://example.com/webhook",
    "method": "POST",
    "headers": "{\"Authorization\":\"Bearer token\"}",
    "content_type": "application/json",
    "template": "{\"title\":\"{{title}}\",\"content\":\"{{content}}\"}",
    "enabled": true,
    "user_id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

### 获取单个推送配置

**GET** `/push/configs/:id`

Headers:
```
Authorization: Bearer {token}
```

响应: 同上单个配置对象

### 创建推送配置

**POST** `/push/configs`

Headers:
```
Authorization: Bearer {token}
```

请求体:
```json
{
  "name": "Webhook推送",
  "description": "推送到Webhook服务",
  "url": "https://example.com/webhook",
  "method": "POST",
  "headers": "{\"Authorization\":\"Bearer token\"}",
  "content_type": "application/json",
  "template": "{\"title\":\"{{title}}\",\"content\":\"{{content}}\"}",
  "enabled": true
}
```

响应: 创建的配置对象

### 更新推送配置

**PUT** `/push/configs/:id`

Headers:
```
Authorization: Bearer {token}
```

请求体: 同创建接口

响应: 更新后的配置对象

### 删除推送配置

**DELETE** `/push/configs/:id`

Headers:
```
Authorization: Bearer {token}
```

响应:
```json
{
  "message": "删除成功"
}
```

## 推送操作

### 执行推送

**POST** `/push/send`

Headers:
```
Authorization: Bearer {token}
```

请求体:
```json
{
  "config_id": 1,
  "data": {
    "title": "测试标题",
    "content": "测试内容"
  }
}
```

响应:
```json
{
  "message": "推送成功",
  "history": {
    "id": 1,
    "config_id": 1,
    "config_name": "Webhook推送",
    "url": "https://example.com/webhook",
    "method": "POST",
    "content": "{\"title\":\"测试标题\",\"content\":\"测试内容\"}",
    "status_code": 200,
    "response": "{\"success\":true}",
    "success": true,
    "error": "",
    "duration": 150,
    "user_id": 1,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 获取推送历史

**GET** `/push/history`

Headers:
```
Authorization: Bearer {token}
```

查询参数:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20）
- `config_id`: 配置ID过滤（可选）

响应:
```json
{
  "total": 100,
  "page": 1,
  "page_size": 20,
  "data": [
    {
      "id": 1,
      "config_id": 1,
      "config_name": "Webhook推送",
      "url": "https://example.com/webhook",
      "method": "POST",
      "content": "{\"title\":\"测试标题\"}",
      "status_code": 200,
      "response": "{\"success\":true}",
      "success": true,
      "error": "",
      "duration": 150,
      "user_id": 1,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

## 错误响应

所有错误响应格式统一为:
```json
{
  "error": "错误信息"
}
```

常见HTTP状态码:
- `200`: 成功
- `201`: 创建成功
- `400`: 请求参数错误
- `401`: 未授权/Token无效
- `403`: 无权限
- `404`: 资源不存在
- `500`: 服务器内部错误

## 推送模板变量

模板中可以使用 `{{变量名}}` 的格式定义变量，在推送时通过 `data` 字段传入实际值。

示例模板:
```json
{
  "title": "{{title}}",
  "content": "{{content}}",
  "time": "{{time}}"
}
```

推送时传入:
```json
{
  "config_id": 1,
  "data": {
    "title": "告警通知",
    "content": "服务器CPU使用率过高",
    "time": "2024-01-01 12:00:00"
  }
}
```

最终发送:
```json
{
  "title": "告警通知",
  "content": "服务器CPU使用率过高",
  "time": "2024-01-01 12:00:00"
}
```
