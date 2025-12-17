# 构建阶段
FROM node:18-alpine AS builder

WORKDIR /app

# 配置 Alpine 镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 设置 npm 淘宝镜像（加速 npm 安装）
RUN npm config set registry https://registry.npmmirror.com

# 复制 package.json
COPY package*.json ./

# 安装依赖
RUN npm install

# 复制源代码
COPY . .

# 构建（直接使用 vite build，跳过 TypeScript 类型检查）
RUN npx vite build

# 运行阶段 - 使用 Nginx 提供静态文件服务
FROM nginx:alpine

# 配置 Alpine 镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 复制构建产物到 Nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 Nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]