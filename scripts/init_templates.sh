#!/bin/bash

# 初始化应用商店模板数据
# 使用方法: ./scripts/init_templates.sh

API_BASE="http://localhost:8080/api/v1"
TOKEN=""

# 登录获取 token
echo "正在登录..."
LOGIN_RESPONSE=$(curl -s -X POST "${API_BASE}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "登录失败，请检查用户名和密码"
  exit 1
fi

echo "登录成功，Token: ${TOKEN:0:20}..."

# 创建 Hello World 模板
echo ""
echo "创建 Hello World 模板..."
curl -X POST "${API_BASE}/appstore/templates" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
  "name": "hello_world",
  "display_name": "Hello World",
  "description": "一个简单的 Hello World 示例应用，可以自定义问候语",
  "code": "package goapp\n\nimport \"fmt\"\n\nfunc Run(input map[string]any) (map[string]any, error) {\n\tname := \"World\"\n\tif n, ok := input[\"name\"].(string); ok && n != \"\" {\n\t\tname = n\n\t}\n\tmessage := fmt.Sprintf(\"Hello, %s!\", name)\n\treturn map[string]any{\"message\": message}, nil\n}",
  "language": "go",
  "category": "示例",
  "version": "1.0.0",
  "author": "Admin",
  "tags": "hello,demo,示例",
  "enabled": true
}'

# 创建 HTTP 请求模板
echo ""
echo "创建 HTTP 请求模板..."
curl -X POST "${API_BASE}/appstore/templates" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
  "name": "http_request",
  "display_name": "HTTP 请求工具",
  "description": "发送 HTTP GET 请求并返回响应结果",
  "code": "package goapp\n\nimport (\n\t\"fmt\"\n\t\"github.com/go-resty/resty/v2\"\n)\n\nfunc Run(input map[string]any) (map[string]any, error) {\n\turl, ok := input[\"url\"].(string)\n\tif !ok || url == \"\" {\n\t\treturn nil, fmt.Errorf(\"url is required\")\n\t}\n\n\tclient := resty.New()\n\tresp, err := client.R().Get(url)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\n\treturn map[string]any{\n\t\t\"status\": resp.StatusCode(),\n\t\t\"body\":   resp.String(),\n\t\t\"headers\": resp.Header(),\n\t}, nil\n}",
  "language": "go",
  "category": "工具",
  "version": "1.0.0",
  "author": "Admin",
  "tags": "http,request,网络",
  "enabled": true
}'

# 创建 JSON 转换模板
echo ""
echo "创建 JSON 转换模板..."
curl -X POST "${API_BASE}/appstore/templates" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
  "name": "json_converter",
  "display_name": "JSON 转换器",
  "description": "将输入数据转换为 JSON 格式",
  "code": "package goapp\n\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n)\n\nfunc Run(input map[string]any) (map[string]any, error) {\n\tdata, ok := input[\"data\"]\n\tif !ok {\n\t\treturn nil, fmt.Errorf(\"data is required\")\n\t}\n\n\tjsonData, err := json.MarshalIndent(data, \"\", \"  \")\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\n\treturn map[string]any{\n\t\t\"json\": string(jsonData),\n\t\t\"type\": fmt.Sprintf(\"%T\", data),\n\t}, nil\n}",
  "language": "go",
  "category": "数据处理",
  "version": "1.0.0",
  "author": "Admin",
  "tags": "json,转换,数据处理",
  "enabled": true
}'

# 创建时间格式化模板
echo ""
echo "创建时间格式化模板..."
curl -X POST "${API_BASE}/appstore/templates" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
  "name": "time_formatter",
  "display_name": "时间格式化",
  "description": "获取当前时间并按指定格式返回",
  "code": "package goapp\n\nimport (\n\t\"time\"\n)\n\nfunc Run(input map[string]any) (map[string]any, error) {\n\tformat := \"2006-01-02 15:04:05\"\n\tif f, ok := input[\"format\"].(string); ok && f != \"\" {\n\t\tformat = f\n\t}\n\n\tnow := time.Now()\n\treturn map[string]any{\n\t\t\"timestamp\": now.Unix(),\n\t\t\"formatted\": now.Format(format),\n\t\t\"iso8601\":   now.Format(time.RFC3339),\n\t}, nil\n}",
  "language": "go",
  "category": "工具",
  "version": "1.0.0",
  "author": "Admin",
  "tags": "time,时间,格式化",
  "enabled": true
}'

echo ""
echo "模板初始化完成！"
