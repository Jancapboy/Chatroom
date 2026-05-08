# 阶段1：构建
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 安装git（部分依赖可能需要）
RUN apk add --no-cache git ca-certificates tzdata

# 先复制依赖文件，利用Docker缓存层
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 编译二进制（静态链接，无CGO）
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags='-w -s' -o /chatroom ./cmd/chatroom/main.go

# 阶段2：运行
FROM alpine:latest

# 安装基础工具
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 复制二进制
COPY --from=builder /chatroom /app/chatroom

# 创建configs目录并写入默认配置（由entrypoint在启动时根据环境变量重写）
COPY configs/config.yaml /app/configs/config.yaml
COPY docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["/app/chatroom"]
