# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 配置 Alpine 镜像源（使用阿里云镜像）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装编译依赖（CGO 需要 gcc）
RUN apk add --no-cache gcc musl-dev sqlite-dev

# 设置 Go 代理（使用中国镜像加速下载）
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译（启用 CGO 以支持 SQLite，添加 musl 兼容标签）
RUN CGO_ENABLED=1 GOOS=linux go build -tags musl -ldflags="-s -w" -o main .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 配置 Alpine 镜像源（运行阶段也需要）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates sqlite-libs tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/main .

# 创建数据目录
RUN mkdir -p /app/data

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]