# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates sqlite-libs

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/main .

# 创建数据目录
RUN mkdir -p /app/data

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]
