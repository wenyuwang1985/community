# 构建阶段
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o community_server ./cmd/community_server

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装 ca-certificates（用于 HTTPS）和 tzdata（用于时区）
RUN apk --no-cache add ca-certificates tzdata

# 创建上传目录
RUN mkdir -p assets/uploads

# 复制编译后的二进制文件和配置
COPY --from=builder /app/community_server .
COPY --from=builder /app/configs ./configs

# 非 root 用户运行（可选，为了简单先用 root）
# RUN adduser -D -u 1000 appuser && chown -R appuser:appuser /app
# USER appuser

EXPOSE 8080

CMD ["./community_server", "-config", "configs/config.docker.yaml"]
