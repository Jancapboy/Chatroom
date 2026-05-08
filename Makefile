# Chatroom-ASI — Makefile

.PHONY: dev-backend dev-frontend build docker-up docker-down test seed clean help

# 默认目标
.DEFAULT_GOAL := help

# 后端开发（本地，需要本地MySQL+Redis）
dev-backend:
    cd backend && go run ./cmd/chatroom/main.go

# 前端开发
dev-frontend:
    cd frontend && npm run dev

# 构建前后端（不依赖Docker）
build: build-backend build-frontend

build-backend:
    cd backend && CGO_ENABLED=0 go build -ldflags='-w -s' -o chatroom ./cmd/chatroom/main.go

build-frontend:
    cd frontend && npm run build

# Docker 一键启动/停止
docker-up:
    docker-compose up --build -d

docker-down:
    docker-compose down -v

# 测试（前后端）
test: test-backend test-frontend

test-backend:
    cd backend && go test ./...

test-frontend:
    @echo "前端测试: cd frontend && npm run test（如需添加测试框架）"

# 种子数据：插入Agent模板
seed:
    @echo "种子数据已由后端 initAgentTemplates() 在启动时自动插入。"
    @echo "如需手动重置，可进入MySQL执行 DEV_SPEC.md 中的 INSERT 语句。"

# 清理
clean:
    cd backend && rm -f chatroom
    cd frontend && rm -rf dist

help:
    @echo "Chatroom-ASI 可用目标："
    @echo "  make dev-backend   # 本地启动后端（go run）"
    @echo "  make dev-frontend  # 本地启动前端（npm run dev）"
    @echo "  make build         # 构建前后端二进制/静态资源"
    @echo "  make docker-up     # Docker Compose 一键启动"
    @echo "  make docker-down   # Docker Compose 停止并清理卷"
    @echo "  make test          # 运行前后端测试"
    @echo "  make seed          # 提示种子数据已自动插入"
    @echo "  make clean         # 清理构建产物"
