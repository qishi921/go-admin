# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o admin-system .

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 从构建阶段复制二进制文件
COPY --from=builder /app/admin-system .

# 创建必要的目录
RUN mkdir -p /app/storage/logs /app/storage/uploads /app/storage/database

# 设置环境变量
ENV TZ=Asia/Shanghai
ENV APP_ENV=production

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动服务
CMD ["./admin-system"]