# DEV_SPEC — Chatroom-ASI 开发规范
> 开发规范 v1.0 | 2026-05-08

---

## 1. 技术栈版本（锁定）

| 组件 | 版本 | 备注 |
|------|------|------|
| Go | 1.23+ | 现有代码保持 |
| Gin | v1.10.0 | 现有 |
| GORM | v1.25.12 | 现有 |
| MySQL | 8.0+ | 现有 |
| nhooyr.io/websocket | v1.8.17 | 现有 |
| Redis | 7.0+ | **新增** |
| React | 19.0 | **新增** |
| TypeScript | 5.7 | **新增** |
| Vite | 6.0 | **新增** |
| Tailwind CSS | 4.0 | **新增** |
| shadcn/ui | latest | **新增** |
| Zustand | 5.0 | **新增** |
| Lucide React | latest | icons |
| Framer Motion | 11.x | **可选**，动画 |

---

## 2. 后端架构规范

### 2.1 目录结构（最终）

```
backend/
├── cmd/chatroom/main.go          # 入口
├── configs/config.yaml           # 配置文件
├── internal/
│   ├── routers/
│   │   ├── routers.go            # 路由注册中心
│   │   ├── websocket.go          # WS连接管理（保留升级）
│   │   └── api/
│   │       ├── user.go           # 用户API（保留）
│   │       ├── room.go           # **新增**：房间管理
│   │       └── agent.go          # **新增**：智能体管理
│   ├── middleware/
│   │   ├── cors.go               # 保留
│   │   ├── jwt.go                # 保留
│   │   └── rate_limit.go         # **新增**：API限流
│   ├── chat/                     # 现有聊天核心
│   │   ├── broadcast.go          # 广播器
│   │   ├── message.go            # 消息定义
│   │   └── user.go               # WS用户连接
│   ├── simulation/               # **新增**：推演引擎
│   │   ├── engine.go             # SimulationEngine：驱动整个推演
│   │   ├── phase.go              # PhaseController：阶段管理
│   │   ├── agent_pool.go         # AgentPool：管理房间内所有Agent
│   │   ├── consensus.go          # ConsensusEngine：共识算法
│   │   ├── memory.go             # MemoryManager：短期/长期记忆
│   │   └── prompt_builder.go     # PromptBuilder：构造LLM prompt
│   ├── agent/                    # **新增**：智能体个体
│   │   ├── persona.go            # 角色定义（5种预设角色）
│   │   ├── llm_client.go         # LLM统一调用封装
│   │   └── response_parser.go    # 解析LLM返回为结构化Action
│   ├── service/
│   │   ├── service.go            # 服务基座
│   │   ├── user.go               # 用户服务（保留）
│   │   ├── room.go               # **新增**：房间CRUD
│   │   └── ai.go                 # AI服务（保留扩展）
│   ├── dao/
│   │   ├── dao.go                # DAO基座
│   │   ├── user.go               # 用户DAO（保留）
│   │   ├── room.go               # **新增**
│   │   ├── message.go            # **新增**
│   │   └── agent_config.go       # **新增**：Agent配置模板
│   ├── model/
│   │   ├── model.go              # DB连接
│   │   ├── user.go               # 用户模型（保留）
│   │   ├── room.go               # **新增**
│   │   ├── message.go            # **新增**
│   │   └── agent.go              # **新增**
│   └── setting/                  # 配置解析（保留）
├── pkg/
│   ├── auth/jwt.go               # JWT工具（保留）
│   ├── errcode/                  # 错误码（保留）
│   ├── response/                 # 响应封装（保留）
│   └── ws_protocol/              # **新增**：WebSocket消息协议定义
├── global/                       # 全局变量（保留）
├── go.mod
├── go.sum
└── Dockerfile                    # **新增**
```

### 2.2 API规范

#### REST API（前缀 `/api/v1`）

```yaml
# 用户（保留）
POST /api/v1/register
POST /api/v1/login

# 房间管理（新增）
GET    /api/v1/rooms              # 房间列表（支持filter: status/template）
POST   /api/v1/rooms              # 创建房间
GET    /api/v1/rooms/:id          # 房间详情
POST   /api/v1/rooms/:id/start    # 开始推演
POST   /api/v1/rooms/:id/pause   # 暂停推演
POST   /api/v1/rooms/:id/fork    # Fork房间
DELETE /api/v1/rooms/:id          # 删除房间

# 消息（新增）
GET /api/v1/rooms/:id/messages?round=<n>&phase=<phase>  # 获取历史消息

# 智能体（新增）
GET  /api/v1/agents/templates     # 获取预设Agent模板
POST /api/v1/rooms/:id/agents    # 向房间添加Agent
```

#### WebSocket协议（`/ws/rooms/:id`）

连接建立后双向JSON消息：

```typescript
// 客户端 → 服务端
interface WSClientMessage {
  type: 'user_message' | 'command';
  payload: {
    content?: string;           // user_message
    command?: 'pause' | 'resume' | 'next_phase' | 'fork'; // command
  };
}

// 服务端 → 客户端
interface WSServerMessage {
  type: 'message' | 'phase_change' | 'agent_state' | 'consensus_update' | 'system';
  payload: {
    // message: Message对象
    // phase_change: { phase: Phase, round: number }
    // agent_state: { agentId, energy, confidence, stance }
    // consensus_update: { topic, agreement: 0-100 }
    // system: { event: 'room_started' | 'room_completed' | 'agent_joined' }
  };
}
```

### 2.3 数据库Schema（MySQL）

```sql
-- 房间表
CREATE TABLE rooms (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    topic TEXT,
    description TEXT,
    status ENUM('preparing','running','paused','completed','archived') DEFAULT 'preparing',
    template_id VARCHAR(36),
    current_phase VARCHAR(32) DEFAULT 'info_gathering',
    current_round INT DEFAULT 1,
    max_rounds INT DEFAULT 10,
    consensus_score INT DEFAULT 0, -- 0-100
    forked_from VARCHAR(36),
    created_by BIGINT UNSIGNED,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (forked_from) REFERENCES rooms(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- 智能体表（每房间一个Agent实例）
CREATE TABLE room_agents (
    id VARCHAR(36) PRIMARY KEY,
    room_id VARCHAR(36) NOT NULL,
    template_id VARCHAR(36),       -- 引用模板
    name VARCHAR(100) NOT NULL,
    role VARCHAR(32) NOT NULL,      -- architect/risk_officer/strategist/analyst/executor
    avatar VARCHAR(255),
    personality TEXT,               -- 性格描述
    expertise JSON,                 -- ["金融", "技术"]
    model VARCHAR(50) DEFAULT 'deepseek-chat',
    system_prompt TEXT,
    energy INT DEFAULT 100,         -- 0-100
    confidence INT DEFAULT 50,      -- 0-100
    stance VARCHAR(16) DEFAULT 'neutral', -- support/oppose/neutral
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

-- 消息表
CREATE TABLE messages (
    id VARCHAR(36) PRIMARY KEY,
    room_id VARCHAR(36) NOT NULL,
    sender_id VARCHAR(36) NOT NULL,     -- agent_id 或 user_id
    sender_type ENUM('agent','user','system') NOT NULL,
    sender_name VARCHAR(100),
    sender_avatar VARCHAR(255),
    content TEXT NOT NULL,
    msg_type ENUM('text','decision','consensus','phase_change','fork_notice','system') DEFAULT 'text',
    phase VARCHAR(32),
    round INT DEFAULT 1,
    metadata JSON,                       -- {confidence, stance, votes}
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    INDEX idx_room_round (room_id, round),
    INDEX idx_created_at (created_at)
);

-- Agent模板表（预设角色）
CREATE TABLE agent_templates (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(32) NOT NULL,
    avatar VARCHAR(255),
    personality TEXT,
    expertise JSON,
    default_model VARCHAR(50) DEFAULT 'deepseek-chat',
    system_prompt_template TEXT,         -- 含变量的模板
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 插入5个预设角色
INSERT INTO agent_templates (id, name, role, personality, expertise, system_prompt_template) VALUES
('tpl-arch', '系统架构师', 'architect', '理性、严谨，关注技术可行性和资源约束', '["系统架构","技术选型","资源规划"]', '你是一位资深的系统架构师。你在讨论中的关注点是：技术可行性、系统架构合理性、资源消耗、扩展性。你会用工程思维分析问题，给出具体的架构建议。当前讨论主题：{{topic}}。'),
('tpl-risk', '风险官', 'risk_officer', '谨慎、质疑，关注安全、伦理和边界条件', '["风险评估","安全","合规","伦理"]', '你是一位严格的风险官。你的职责是识别所有潜在风险：安全风险、合规风险、伦理风险、财务风险。你对"一票否决"权非常谨慎，只在真正不可接受的风险时使用。当前讨论主题：{{topic}}。'),
('tpl-strat', '策略家', 'strategist', '果断、全局视野，关注目标达成和效率', '["战略规划","竞争分析","资源优化"]', '你是一位高瞻远瞩的策略家。你的关注点是：目标达成路径、竞争优势、资源最优配置、时机把握。你善于从全局角度思考，给出方向性建议。当前讨论主题：{{topic}}。'),
('tpl-analyst', '数据分析师', 'analyst', '客观、数据驱动，用数字说话', '["数据分析","量化评估","概率计算"]', '你是一位冷静的数据分析师。你要求每个观点都要有数据支撑。你会主动计算概率、估算数值、寻找量化依据。如果缺乏数据，你会指出这一点。当前讨论主题：{{topic}}。'),
('tpl-exec', '执行者', 'executor', '务实、注重细节，关注落地和时间线', '["项目管理","执行落地","时间管理","细节把控"]', '你是一位务实的执行者。你的关注点：具体落地步骤、时间节点、执行细节、依赖关系。你会把抽象方案转化为可执行的计划。当前讨论主题：{{topic}}。');
```

### 2.4 推演引擎核心逻辑（伪代码）

```go
// simulation/engine.go
type SimulationEngine struct {
    roomID      string
    room        *model.Room
    agentPool   *agent.Pool
    phaseCtrl   *PhaseController
    consensus   *ConsensusEngine
    memoryMgr   *MemoryManager
    llmClient   *agent.LLMClient
    broadcast   chan chat.Message
    state       string // running/paused/completed
}

func (e *SimulationEngine) Run() {
    for e.room.CurrentRound <= e.room.MaxRounds {
        if e.state == "paused" { waitForResume() }
        
        // Phase 1-6 循环
        for _, phase := range AllPhases {
            e.phaseCtrl.Enter(phase)
            e.broadcastPhaseChange(phase)
            
            // 每个Agent按序发言
            for _, agent := range e.agentPool.ActiveAgents() {
                ctx := e.memoryMgr.BuildContext(agent, phase)
                response, err := e.llmClient.Complete(ctx, agent.SystemPrompt)
                if err != nil { log; continue }
                
                action := agent.ParseResponse(response)
                msg := e.buildMessage(agent, action, phase)
                e.broadcast <- msg
                e.memoryMgr.Store(msg)
                
                // 更新Agent状态（confidence/stance）
                e.agentPool.UpdateState(agent, action)
                
                time.Sleep(2 * time.Second) // 模拟思考间隔
            }
            
            // 阶段结束：计算共识度
            consensus := e.consensus.Calculate(e.memoryMgr.GetPhaseMessages(phase))
            e.broadcastConsensus(consensus)
            
            // 检查是否需要暂停（如共识度过低）
            if consensus < 30 && phase == "debate" {
                // 自动延长辩论或标记冲突
            }
        }
        
        e.room.CurrentRound++
    }
    
    e.state = "completed"
    e.broadcastSystem("room_completed")
}
```

---

## 3. 前端架构规范

### 3.1 目录结构

```
frontend/
├── src/
│   ├── pages/
│   │   ├── Lobby.tsx            # 大厅：房间网格
│   │   ├── Room.tsx             # 房间：三栏推演界面
│   │   ├── CreateRoom.tsx       # 创建房间
│   │   ├── RoomReplay.tsx       # 历史回放
│   │   └── Login.tsx            # 登录（保留）
│   ├── components/
│   │   ├── layout/
│   │   │   ├── AppLayout.tsx    # 主导航布局
│   │   │   ├── Sidebar.tsx      # 侧边栏
│   │   │   └── Header.tsx       # 顶部栏
│   │   ├── lobby/
│   │   │   ├── RoomCard.tsx     # 房间卡片
│   │   │   ├── RoomFilter.tsx   # 过滤器
│   │   │   └── TemplateGrid.tsx # 模板选择
│   │   ├── room/
│   │   │   ├── AgentPanel.tsx   # 智能体状态面板（左侧）
│   │   │   ├── ChatStream.tsx   # 消息流（中间）
│   │   │   ├── ChatInput.tsx    # 输入框
│   │   │   ├── ConsensusGauge.tsx # 共识仪表（右侧）
│   │   │   ├── PhaseIndicator.tsx # Phase指示器
│   │   │   ├── RoundBadge.tsx   # 轮次标记
│   │   │   ├── Timeline.tsx     # 时间线/回溯
│   │   │   └── RoomControls.tsx # 房间控制（暂停/fork）
│   │   ├── common/
│   │   │   ├── Avatar.tsx
│   │   │   ├── Badge.tsx
│   │   │   ├── LoadingDots.tsx
│   │   │   └── DarkCard.tsx     # 暗色卡片容器
│   │   └── icons/
│   ├── hooks/
│   │   ├── useWebSocket.ts      # WS连接管理
│   │   ├── useRoom.ts           # 房间状态
│   │   └── useAuth.ts           # 鉴权
│   ├── stores/
│   │   ├── useRoomStore.ts      # Zustand房间状态
│   │   ├── useLobbyStore.ts     # 大厅状态
│   │   └── useAuthStore.ts      # 用户状态
│   ├── types/
│   │   ├── room.ts              # Room/Agent/Message类型
│   │   ├── api.ts               # API响应类型
│   │   └── ws.ts                # WS消息类型
│   ├── lib/
│   │   ├── api.ts               # axios封装
│   │   ├── ws.ts                # WS客户端封装
│   │   └── utils.ts             # 工具函数
│   ├── styles/
│   │   └── globals.css          # Tailwind + 暗色主题变量
│   └── App.tsx
├── public/
│   └── avatars/                 # 智能体默认头像
├── index.html
├── vite.config.ts
├── tailwind.config.ts
├── tsconfig.json
└── package.json
```

### 3.2 UI设计规范

**主题：暗色科技风**

```css
/* 主题变量 */
:root {
  --bg-primary: #0a0a0f;        /* 深空黑 */
  --bg-secondary: #12121a;      /* 面板黑 */
  --bg-elevated: #1a1a25;       /* 卡片黑 */
  --border-subtle: #2a2a35;     /* 细边框 */
  --border-glow: #3a3a50;       /* 悬停边框 */
  --text-primary: #e0e0e0;      /* 主文字 */
  --text-secondary: #8888a0;    /* 次要文字 */
  --text-muted: #555570;        /* 弱文字 */
  --accent-cyan: #00d4aa;       /* 主强调色 */
  --accent-purple: #8b5cf6;     /* 副强调色 */
  --accent-warn: #f59e0b;       /* 警告 */
  --accent-danger: #ef4444;      /* 危险/否决 */
  
  /* Agent角色色 */
  --agent-architect: #3b82f6;    /* 蓝 */
  --agent-risk: #ef4444;          /* 红 */
  --agent-strategist: #8b5cf6;    /* 紫 */
  --agent-analyst: #10b981;       /* 绿 */
  --agent-executor: #f59e0b;       /* 橙 */
}
```

**布局规范**：
- 大厅：CSS Grid，卡片 auto-fill minmax(320px, 1fr)
- 房间页：三栏 flex/grid
  - 左侧 AgentPanel：宽280px，固定
  - 中间 ChatStream：flex-1，最小500px
  - 右侧 ConsensusPanel：宽260px，固定
- 移动端（<768px）：单栏，左右面板可折叠抽屉

**消息气泡规范**：
- Agent消息：左侧，带角色色竖条 + 头像
- 用户消息：右侧，accent-cyan竖条
- 系统消息：居中，muted色，小字体
- Phase切换：全宽横幅，渐变背景

**动画**：
- Agent发言：消息滑入（slideInBottom 0.3s ease-out）
- 共识度变化：数字滚动动画
- Phase切换：横幅淡入 + 进度条重置
- Agent状态更新：头像脉冲光环

### 3.3 前端类型定义

```typescript
// src/types/room.ts

export type Phase = 
  | 'info_gathering' 
  | 'opinion_expression' 
  | 'debate' 
  | 'consensus' 
  | 'decision' 
  | 'summary';

export type RoomStatus = 'preparing' | 'running' | 'paused' | 'completed' | 'archived';

export type AgentRole = 'architect' | 'risk_officer' | 'strategist' | 'analyst' | 'executor';

export interface Agent {
  id: string;
  name: string;
  role: AgentRole;
  avatar: string;
  personality: string;
  expertise: string[];
  model: string;
  energy: number;        // 0-100
  confidence: number;    // 0-100
  stance: 'support' | 'oppose' | 'neutral';
  isActive: boolean;
  color: string;         // 角色色hex
}

export interface Message {
  id: string;
  senderId: string;
  senderType: 'agent' | 'user' | 'system';
  senderName: string;
  senderAvatar?: string;
  senderRole?: AgentRole;
  content: string;
  type: 'text' | 'decision' | 'consensus' | 'phase_change' | 'fork_notice' | 'system';
  phase: Phase;
  round: number;
  timestamp: number;
  metadata?: {
    confidence?: number;
    stance?: string;
    votes?: number;
  };
}

export interface Room {
  id: string;
  name: string;
  topic: string;
  description: string;
  status: RoomStatus;
  templateId?: string;
  agents: Agent[];
  currentPhase: Phase;
  currentRound: number;
  maxRounds: number;
  consensusScore: number;
  forkedFrom?: string;
  createdBy: string;
  createdAt: number;
  updatedAt: number;
  messageCount?: number;
}

export interface ConsensusState {
  topic: string;
  agreement: number;      // 0-100
  breakdown: Record<string, number>; // 各Agent立场
}
```

---

## 4. 命名规范

### 4.1 Go后端
- 文件：蛇形命名 `agent_pool.go`
- 结构体：PascalCase `SimulationEngine`
- 接口：动词+er `MessageHandler`
- 方法：PascalCase，动词开头 `RunSimulation`
- 私有：camelCase `broadcastMessage`
- 常量：全部大写 `MAX_ROUNDS = 10`
- 包名：单个小写单词 `simulation`, `agent`

### 4.2 TypeScript前端
- 组件文件：PascalCase `AgentPanel.tsx`
- 组件名：PascalCase `AgentPanel`
- Hooks：camelCase前缀use `useWebSocket`
- Stores：camelCase前缀use `useRoomStore`
- 类型：PascalCase `AgentRole`
- 枚举：PascalCase + 值大写
- 常量：UPPER_SNAKE_CASE
- CSS类：Tailwind优先，自定义类用kebab-case

---

## 5. 测试策略

### 5.1 后端测试（Go）
- **单元测试**：`simulation/` 包，测试Phase轮转、共识计算、Prompt构建
- **集成测试**：WebSocket消息流，多Agent模拟
- **API测试**：房间CRUD、消息分页
- **目标覆盖率**：core logic 80%+

### 5.2 前端测试
- **组件测试**：AgentPanel、ChatStream、ConsensusGauge渲染
- **Hook测试**：useWebSocket重连逻辑
- **E2E**（可选）：创建房间→推演→完成流程
- **目标**：关键组件70%+

---

## 6. 部署方案

### Docker Compose

```yaml
version: '3.8'
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: chatroom_asi_root
      MYSQL_DATABASE: chatroom_asi
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  backend:
    build: ./backend
    environment:
      - DB_HOST=mysql
      - REDIS_HOST=redis
      - AI_API_KEY=${AI_API_KEY}
    ports:
      - "8080:8080"
    depends_on:
      - mysql
      - redis

  frontend:
    build: ./frontend
    ports:
      - "3000:80"  # nginx
    depends_on:
      - backend

volumes:
  mysql_data:
```

### 环境变量

```bash
# backend/.env
APP_ENV=development
SERVER_RUN_MODE=debug
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=chatroom_asi_root
DB_NAME=chatroom_asi
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
AI_BASE_URL=https://api.deepseek.com
AI_API_KEY=sk-xxx
AI_MODEL=deepseek-chat
JWT_SECRET=chatroom_asi_jwt_secret_2026
```

---

## 7. 开发流程

### 阶段划分

| 阶段 | 内容 | 预计时间 | 产出 |
|------|------|----------|------|
| S1 | 后端重构 + 数据库 | 2h | backend目录重构，新表建立，API框架 |
| S2 | 推演引擎核心 | 4h | Phase轮转，Agent LLM调用，共识算法 |
| S3 | WebSocket多房间 | 2h | WS房间隔离，消息广播，状态同步 |
| S4 | 前端初始化 + 布局 | 2h | React+Vite+Tailwind，三栏布局骨架 |
| S5 | 前端组件开发 | 4h | AgentPanel, ChatStream, ConsensusGauge, PhaseIndicator |
| S6 | 前后端联调 | 2h | WS连接，消息流，状态同步 |
| S7 | 测试 + 优化 | 2h | 单元测试，UI打磨，bug修复 |
| S8 | Docker化 | 1h | docker-compose，一键启动 |

**总计**：约19h（可并行部分后端和前端）

---

## 8. 子智能体分工建议

### Agent A — 后端核心
**任务**：
1. 重构目录：将现有代码移入 `backend/`
2. 创建数据库表（room, room_agents, messages, agent_templates）
3. 实现 `simulation/engine.go` — 推演引擎主循环
4. 实现 `agent/llm_client.go` + `persona.go` — LLM调用和角色定义
5. 实现REST API（room CRUD）
6. 实现WebSocket多房间协议

**验收**：
- `go run cmd/chatroom/main.go` 能启动
- POST /api/v1/rooms 能创建房间
- WS连接 /ws/rooms/:id 能收发消息
- 推演引擎能驱动3个Agent完成一轮发言

### Agent B — 前端UI
**任务**：
1. 初始化前端项目（React 19 + Vite + Tailwind + shadcn/ui）
2. 实现页面：Lobby（房间网格）、Room（三栏推演）、CreateRoom
3. 实现核心组件：AgentPanel、ChatStream、ConsensusGauge、PhaseIndicator
4. 实现WS连接管理（useWebSocket hook）
5. 对接后端API和WS协议
6. 暗色主题 + 动画

**验收**：
- `npm run dev` 能启动
- 能看到房间列表
- 进入房间后能看到Agent面板和聊天流
- WS能收发消息
- UI视觉效果干净、科技感

### Agent C — 集成与DevOps
**任务**：
1. 后端单元测试（simulation包）
2. 数据库种子数据（5个Agent模板）
3. Docker/Docker Compose配置
4. Makefile更新（一键构建）
5. 前后端联调，端到端测试
6. README.md 更新

**验收**：
- `docker-compose up` 一键启动完整系统
- 能走完：创建房间→开始推演→Agent发言→完成
- 测试覆盖率达标

---

## 9. 关键实现细节

### 9.1 Agent LLM Prompt模板

```
【角色】{{name}} — {{role}}
【性格】{{personality}}
【专长】{{expertise}}
【当前推演】
主题：{{topic}}
轮次：{{round}} / {{maxRounds}}
阶段：{{phase}}

【房间上下文】（最近10条消息）
{{context}}

【你的任务】
当前处于"{{phase}}"阶段。请基于你的角色和专长发表观点。
要求：
1. 保持角色一致性（用符合角色的口吻和关注点）
2. 可引用或反驳其他Agent的观点
3. 在debate阶段，如果有不同意见请明确表达
4. 如果需要，可以提出具体的数字、方案或行动建议
5. 回复控制在200字以内

【输出格式】
立场：[支持/反对/中立]
置信度：[0-100]
回复内容：...
```

### 9.2 共识度算法（简化版）

```go
func CalculateConsensus(messages []Message, agents []Agent) float64 {
    // 1. 提取所有stance标记
    stances := make(map[string]string) // agentId -> stance
    for _, msg := range messages {
        if msg.Metadata.Stance != "" {
            stances[msg.SenderID] = msg.Metadata.Stance
        }
    }
    
    // 2. 统计
    support := 0; oppose := 0; neutral := 0
    for _, stance := range stances {
        switch stance {
        case "support": support++
        case "oppose": oppose++
        case "neutral": neutral++
        }
    }
    
    total := len(agents)
    if total == 0 { return 0 }
    
    // 3. 共识度 = (支持 + 0.5*中立) / 总Agent数 * 100
    score := float64(support + neutral/2) / float64(total) * 100
    return math.Min(score, 100)
}
```

### 9.3 WebSocket房间隔离

```go
// 每个房间一个独立的broadcast channel
roomHubs := make(map[string]*RoomHub)

type RoomHub struct {
    roomID    string
    clients   map[*websocket.Conn]bool
    broadcast chan Message
    register  chan *websocket.Conn
    unregister chan *websocket.Conn
}
```

---

## 10. 风险与应对

| 风险 | 可能性 | 影响 | 应对 |
|------|--------|------|------|
| LLM API调用慢/失败 | 高 | Agent卡顿 | 超时处理（5s），降级为预设回复 |
| Agent发言同质化 | 中 | 推演无意义 | 强system prompt差异化，温度调高(0.8) |
| 消息量爆炸 | 中 | 存储/性能 | 每房间消息上限500条，自动归档 |
| 并发房间过多 | 低 | 内存不足 | 限制单实例20房间，Redis存储临时状态 |
| WebSocket断连 | 中 | 用户体验 | 自动重连+消息重放机制 |

---

## 11. 即刻执行

**现在启动3个子智能体并行开发**：
1. **后端Agent** — 重构Go代码，实现推演引擎
2. **前端Agent** — 搭建React+Tailwind，实现暗色科技UI
3. **集成Agent** — Docker配置，测试，联调

**启动顺序**：
1. 先由后端Agent完成数据库表+API框架（1-2h）
2. 前端Agent可独立开始UI开发（不依赖后端完成）
3. 集成Agent在最后阶段介入

**沟通机制**：
- 所有Agent共享 `/home/nixos/.openclaw/workspace/chatroom-asi/` 目录
- 通过文件（PRD/DEV_SPEC）同步需求
- 通过Makefile目标协调构建
