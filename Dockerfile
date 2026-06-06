# 多阶段构建：前端
FROM node:22-alpine AS frontend-build

WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm install

COPY frontend/ .
RUN npm run build

# 多阶段构建：后端
FROM golang:1.22-alpine AS backend-build

WORKDIR /app
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

# 最终镜像
FROM alpine:3.18

WORKDIR /app

# 安装必要的依赖
RUN apk add --no-cache ca-certificates tzdata curl

# 复制前端构建结果
COPY --from=frontend-build /app/frontend/dist /app/frontend/dist

# 复制后端可执行文件
COPY --from=backend-build /app/main /app/

# 复制配置文件
COPY .env* /app/

# 暴露端口
EXPOSE 8080

# 设置时区
ENV TZ=Asia/Shanghai

# 启动命令
CMD ["./main"]