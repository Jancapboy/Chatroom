# Chatroom-ASI

> 一个让多智能体（Multi-Agent）在虚拟聊天室中自主推演、协作决策、迭代进化的 Web 平台。用户可旁观、可插话、可创建自己的推演场景。

---

## 功能特性

- **多智能体推演**：3–7 个不同角色的 AI 智能体在房间内自主轮询发言，模拟真实决策场景
- **6 阶段流程**：信息收集 → 观点表达 → 冲突辩论 → 共识合成 → 决策输出 → 总结归档
- **WebSocket 实时同步**：房间级独立 WS 连接，Agent 发言、Phase 切换、共识度更新实时推送
- **角色差异化**：系统架构师、风险官、策略家、数据分析师、执行者，各有专属 System Prompt 与立场
- **用户介入**：游客可旁观，注册用户可随时插话、创建房间、Fork 推演历史
- **推演可视化**：左侧 Agent 状态面板、中间消息流（按 Phase/Round 分组）、右侧共识仪表
- **房间模板**：内置"产品评审"、"危机应对"、"技术选型"等预设场景一键启动

---

## 快速启动（Docker Compose 一键启动）

```bash
# 1. 进入项目目录
cd chatroom-asi

# 2. 复制环境变量模板并填入 DeepSeek API Key
cp .env.example .env
# 编辑 .env，将 AI_API_KEY 替换为你的真实密钥

# 3. 一键启动（MySQL + Redis + Backend + Frontend）
make docker-up

# 4. 访问前端 http://localhost:3000
#    后端 API  http://localhost:8080
```

停止服务：

```bash
make docker-down
```

---

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | React 19 + TypeScript + Vite + Tailwind CSS |
| UI 库 | shadcn/ui + Lucide React + Framer Motion |
| 状态管理 | Zustand |
| 后端 | Go 1.23 + Gin + GORM |
| WebSocket | nhooyr.io/websocket |
| 数据库 | MySQL 8.0 |
| 缓存 | Redis 7 |
| AI API | DeepSeek（可扩展 Kimi / GPT） |
| 部署 | Docker Compose |

---

## 项目结构

```
chatroom-asi/
├── backend/                    # Go 后端
│   ├── cmd/chatroom/main.go    # 入口
│   ├── internal/               # 业务逻辑
│   │   ├── routers/            # REST + WS 路由
│   │   ├── chat/               # 广播器 & WS 连接管理
│   │   ├── simulation/         # 推演引擎（Phase / Agent / 共识）
│   │   ├── agent/              # 智能体 LLM 调用 & Prompt
│   │   ├── service/            # 业务服务层
│   │   ├── dao/                # 数据访问层
│   │   └── model/              # GORM 模型
│   ├── pkg/                    # JWT、错误码、响应封装
│   ├── configs/config.yaml     # 配置文件（容器内由 entrypoint 生成）
│   ├── Dockerfile              # 多阶段构建（Go → Alpine）
│   └── docker-entrypoint.sh    # 启动时从环境变量生成 config.yaml
├── frontend/                   # React 前端
│   ├── src/
│   │   ├── pages/              # Lobby / Room / CreateRoom / Replay
│   │   ├── components/         # AgentPanel / ChatStream / ConsensusGauge
│   │   ├── hooks/              # useWebSocket / useRoom / useAuth
│   │   ├── stores/             # Zustand 状态
│   │   └── types/              # TypeScript 类型定义
│   ├── dist/                   # 构建产物（由 nginx 托管）
│   ├── Dockerfile              # nginx:alpine
│   └── nginx.conf              # SPA 路由回退配置
├── docker-compose.yml          # 完整编排（MySQL + Redis + Backend + Frontend）
├── .env.example                # 环境变量模板
├── Makefile                    # 开发 & 部署快捷命令
├── PRD.md                      # 产品需求文档
└── DEV_SPEC.md                 # 开发规范 & API 定义
```

---

## API 文档简要说明

### REST API（前缀 `/api/v1`）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/register` | 用户注册 |
| POST | `/login` | 用户登录（返回 JWT） |
| GET  | `/rooms` | 房间列表（支持 status/template 过滤） |
| POST | `/rooms` | 创建房间 |
| GET  | `/rooms/:id` | 房间详情 |
| POST | `/rooms/:id/start` | 开始推演 |
| POST | `/rooms/:id/pause` | 暂停推演 |
| POST | `/rooms/:id/fork` | Fork 房间 |
| GET  | `/rooms/:id/messages` | 历史消息（支持 round/phase 过滤） |
| GET  | `/agents/templates` | 获取预设 Agent 模板 |
| POST | `/rooms/:id/agents` | 向房间添加 Agent |

### WebSocket（`/ws/rooms/:id`）

连接建立后双向 JSON 消息：

- **Client → Server**：`user_message`（用户发言）、`command`（pause/resume/next_phase/fork）
- **Server → Client**：`message`（Agent 发言）、`phase_change`（阶段切换）、`agent_state`（状态更新）、`consensus_update`（共识度变化）、`system`（系统事件）

完整协议定义见 `DEV_SPEC.md` §2.2。

---

## 本地开发启动

```bash
# 启动后端（需本地 MySQL + Redis，或先 docker-compose up mysql redis）
make dev-backend

# 另开终端启动前端
cd frontend && npm install   # 首次
make dev-frontend             # 或 cd frontend && npm run dev
```

前端开发服务器已配置 Vite 代理：
- `/api` → `http://localhost:8080`
- `/ws`  → `ws://localhost:8080`

---

## UI 布局（截图占位）

### 大厅（Lobby）
- 顶部导航栏（Logo + 登录/注册）
- 房间卡片网格：显示主题、状态（进行中/已完成）、Agent 头像组、当前轮次
- 过滤器：按状态 / 模板筛选，搜索框
- "创建房间" 悬浮按钮

### 房间页（Room）— 三栏布局
- **左侧栏（280px）**：房间标题 + 主题描述 + Agent 状态面板（头像、角色色、能量条、置信度、立场徽章）
- **中间栏（自适应）**：消息流时间线，按 Phase/Round 分组；用户输入框在底部
- **右侧栏（260px）**：共识度圆形仪表、当前 Phase 指示器（带进度条）、推演控制按钮（暂停 / 快进 / Fork）

### 创建房间（CreateRoom）
- 主题输入 + 描述
- 模板选择卡片（横向滚动）
- Agent 阵容配置（从模板一键加载 / 手动调整角色）
- 最大轮次滑块（默认 10）

---

## License

MIT
