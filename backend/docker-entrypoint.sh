#!/bin/sh
set -e

# 生成 configs/config.yaml，从环境变量读取（适配 DEV_SPEC 中的 viper 配置）
cat > /app/configs/config.yaml <<EOF
Server:
  RunMode: "${SERVER_RUN_MODE:-debug}"
  HttpPort: "8080"
  ReadTimeout: 60
  WriteTimeout: 60

Database:
  DBType: "mysql"
  UserName: "${DB_USER:-root}"
  Password: "${DB_PASSWORD:-chatroom_asi_root}"
  Host: "${DB_HOST:-mysql}:${DB_PORT:-3306}"
  DBName: "${DB_NAME:-chatroom_asi}"
  TablePrefix: ""
  Charset: "utf8mb4"
  ParseTime: "True"
  MaxIdleConns: 10
  MaxOpenConns: 100

JWT:
  Secret: "${JWT_SECRET:-chatroom_asi_jwt_secret_change_me}"
  Issuer: "Chatroom-ASI"
  Expire: 24

Chatroom:
  MessageQueueLength: 100

AI:
  api_key: "${AI_API_KEY:-}"
  model: "${AI_MODEL:-deepseek-chat}"
  base_url: "${AI_BASE_URL:-https://api.deepseek.com}"
EOF

exec "$@"
