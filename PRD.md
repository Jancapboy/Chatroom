# PRD — ASI Chatroom (Chatroom-ASI)
> 产品需求文档 v1.0 | 2026-05-08

---

## 1. 产品定位

**一句话**：一个让多智能体（Multi-Agent）在虚拟聊天室中自主推演、协作决策、迭代进化的Web平台。用户可旁观、可插话、可创建自己的推演场景。

**灵感来源**：
- 现有Go聊天室（gin + websocket + mysql）
- ASI-Emergence-Theory（AETS）的多智能体协作框架
- ASET 太空模拟中的AgentCouncil共识机制
- Chip Foundry 的链式参数传递思想

---

## 2. 核心概念

### 2.1 推演房间（Room）
每个房间是一个独立的推演场景：
- **主题**：如"火星殖民计划"、"公司战略决策"、"AI安全危机应对"
- **智能体阵容**：3-7个不同角色的AI智能体（每个有自己的性格、专长、立场）
- **推演轮次**：多轮对话，每轮所有智能体发言 → 共识/冲突 → 决策
- **迭代进化**：房间可fork出新版本，基于上一轮结果继续推演

### 2.2 智能体角色（Agent Persona）
每个智能体不是通用AI，而是有明确角色定位：
- **系统架构师**：关注技术可行性、资源约束
- **风险官**：关注安全、伦理、边界条件，可一票否决
- **策略家**：关注目标达成、效率优化
- **数据分析师**：提供量化分析、概率评估
- **执行者**：关注落地细节、时间线
- **旁观者（用户）**：人类用户，可随时发言

### 2.3 推演引擎（Simulation Engine）
驱动智能体对话的核心机制：
- **Phase系统**：每轮分阶段（信息收集→观点表达→冲突辩论→共识合成→决策输出）
- **共识算法**：类似AgentCouncil，多数通过/关键事项2/3通过/安全官一票否决
- **状态机**：房间状态（准备中→推演中→暂停→完成→归档）
- **记忆系统**：智能体有短期记忆（当前房间上下文）+ 长期记忆（跨房间经验）

### 2.4 迭代进化（Evolution）
- **Fork**：基于当前推演结果创建新房间，继承上下文，修改参数继续
- **Merge**：两个房间的推演结果可合并对比
- **版本线**：类似Git分支，可视化推演历史

---

## 3. 功能需求

### 3.1 P0 — 必须有

| # | 功能 | 验收标准 |
|---|------|----------|
| 1 | **WebSocket多房间** | 支持多个并发房间，每个房间独立WS连接 |
| 2 | **多智能体自动推演** | 最少3个智能体在房间内自动轮询发言 |
| 3 | **角色差异化** | 每个智能体有独特system prompt，发言风格/关注点不同 |
| 4 | **用户旁观/插话** | 用户可进入房间观看实时推演，随时发送消息 |
| 5 | **推演状态面板** | UI显示：当前Phase、轮次、智能体情绪/立场、共识度 |
| 6 | **房间列表/创建** | 首页展示所有房间（进行中/已完成），可创建新房间 |
| 7 | **AI后端集成** | 接入DeepSeek API（已有），智能体调用LLM生成回复 |
| 8 | **JWT鉴权** | 保留现有JWT，游客可旁观，注册用户可创建/插话 |

### 3.2 P1 — 重要

| # | 功能 | 验收标准 |
|---|------|----------|
| 9 | **推演时间线** | 可视化单房间的轮次演进，可回溯到任意轮次 |
| 10 | **Fork房间** | 基于已有房间创建新版本，继承上下文 |
| 11 | **智能体状态可视化** | 每个智能体显示：情绪、立场、能量、置信度 |
| 12 | **共识度仪表** | 实时显示当前议题的共识百分比 |
| 13 | **房间模板** | 预设模板："产品评审"、"危机应对"、"技术选型"等 |
| 14 | **推演报告导出** | 完成后导出Markdown/PDF报告（含所有轮次记录） |

### 3.3 P2 — 加分

| # | 功能 | 验收标准 |
|---|------|----------|
| 15 | **多模型对比** | 同一房间可用不同LLM（DeepSeek/Kimi/GPT）驱动不同智能体 |
| 16 | **声音播报** | TTS播报关键发言 |
| 17 | **投票机制** | 旁观用户可对智能体观点投票 |
| 18 | **外部API触发** | 推演结果可调用外部API（如发邮件、创建任务） |

---

## 4. 数据模型

### 4.1 Room（房间）
```typescript
interface Room {
  id: string;
  name: string;
  topic: string;           // 推演主题
  description: string;
  status: 'preparing' | 'running' | 'paused' | 'completed' | 'archived';
  templateId?: string;     // 使用的模板
  agents: Agent[];          // 房间内智能体
  currentPhase: Phase;      // 当前推演阶段
  currentRound: number;    // 当前轮次
  maxRounds: number;       // 最大轮次（默认10）
  messages: Message[];      // 消息历史
  consensus: Consensus;     // 当前共识状态
  forkedFrom?: string;      // fork自哪个房间
  createdBy: string;        // 创建者ID
  createdAt: number;
  updatedAt: number;
}
```

### 4.2 Agent（智能体）
```typescript
interface Agent {
  id: string;
  name: string;
  role: string;             // 角色类型：architect/risk_officer/strategist/analyst/executor
  avatar: string;           // 头像URL
  personality: string;      // 性格描述（注入LLM system prompt）
  expertise: string[];      // 专长领域
  stance: Stance;           // 当前立场（支持/反对/中立）
  confidence: number;       // 置信度 0-100
  energy: number;           // 能量/参与度 0-100
  model: string;            // 使用的LLM模型
  systemPrompt: string;     // 完整system prompt
  memory: Message[];        // 短期记忆（最近N条）
  isActive: boolean;        // 是否参与当前轮次
}
```

### 4.3 Message（消息）
```typescript
interface Message {
  id: string;
  roomId: string;
  senderId: string;         // 智能体ID 或 用户ID
  senderType: 'agent' | 'user' | 'system';
  senderName: string;
  senderAvatar?: string;
  content: string;
  type: 'text' | 'decision' | 'consensus' | 'phase_change' | 'fork_notice';
  phase: Phase;
  round: number;
  timestamp: number;
  metadata?: {
    confidence?: number;
    stance?: string;
    votes?: number;        // 用户投票数
  };
}
```

### 4.4 Phase（推演阶段）
```typescript
type Phase = 
  | 'info_gathering'      // 信息收集：各智能体提供已知信息
  | 'opinion_expression'  // 观点表达：各智能体提出方案
  | 'debate'              // 辩论：质疑、反驳、数据支撑
  | 'consensus'           // 共识：寻求共同点
  | 'decision'            // 决策：投票或最终判断
  | 'summary';            // 总结：输出推演报告
```

---

## 5. 推演引擎流程

```
初始化房间 → 加载智能体 → 开始推演
    ↓
Phase 1: 信息收集（各Agent陈述已知事实/数据）
    ↓
Phase 2: 观点表达（各Agent提出立场和方案）
    ↓
Phase 3: 辩论（Agent间交叉质疑，可引用数据）
    ↓
Phase 4: 共识（Agent尝试找到共同点，输出共识度）
    ↓
Phase 5: 决策（关键议题投票，安全官可否决）
    ↓
Phase 6: 总结（生成本轮推演报告）
    ↓
进入下一轮 或 结束
```

**关键规则**：
- 每Phase所有active Agent必须发言
- 辩论Phase允许Agent自由引用其他Agent的观点
- 共识度>80%可自动推进到Decision
- 安全官（risk_officer）在decision phase有一票否决权
- 用户消息可随时插入，Agent下一轮可回应用户

---

## 6. WebUI设计方向

**风格**：极简科技感 + 暗色主题
**参考**：
- Discord的频道列表 + 聊天区布局
- Chip Foundry的深色科技UI
- ASET的仪表盘风格

**核心页面**：
1. **大厅（Lobby）**：房间卡片网格（进行中/已完成/热门）
2. **房间（Room）**：三栏布局
   - 左侧：房间信息 + 智能体状态面板
   - 中间：消息流（时间线式，区分Phase/Round）
   - 右侧：共识仪表 + 推演控制（暂停/快进/fork）
3. **房间创建（Create）**：表单 + 模板选择 + 智能体配置
4. **历史回放（Replay）**：可拖拽时间轴查看任意轮次

---

## 7. 非功能需求

- **并发**：单服务器支持≥20个并发房间，每房间≤7个智能体
- **延迟**：Agent发言间隔2-5秒（模拟思考），WebSocket延迟<100ms
- **存储**：消息持久化到MySQL，房间状态可恢复
- **API成本**：Agent调用DeepSeek API，需控制token消耗（每次请求限制4K context）

---

## 8. 技术栈（整体）

| 层级 | 技术 | 说明 |
|------|------|------|
| 前端 | React 19 + TypeScript + Vite | 全新开发 |
| UI库 | Tailwind CSS + shadcn/ui | 暗色主题 |
| 前端状态 | Zustand | 轻量状态管理 |
| 前端WS | 原生WebSocket | 房间级连接 |
| 后端 | Go + Gin（现有） | 保留并扩展 |
| 后端WS | nhooyr.io/websocket（现有） | 保留 |
| 数据库 | MySQL + GORM（现有） | 新增room/agent/message表 |
| 缓存 | Redis（新增） | 房间状态、Agent短期记忆 |
| AI API | DeepSeek（已有） | 可能加Kimi/StepFun |
| 部署 | Docker Compose（新增） | 一键启动 |

---

## 9. 项目目录结构（目标）

```
chatroom-asi/
├── backend/                    # Go后端（从根迁移）
│   ├── cmd/chatroom/
│   ├── internal/
│   │   ├── routers/
│   │   ├── chat/              # 现有broadcast/user/message
│   │   ├── service/
│   │   ├── dao/
│   │   ├── model/
│   │   ├── middleware/
│   │   ├── simulation/        # 新增：推演引擎
│   │   │   ├── engine.go      # 主引擎
│   │   │   ├── phase.go       # Phase管理
│   │   │   ├── agent.go       # Agent生命周期
│   │   │   ├── consensus.go   # 共识算法
│   │   │   └── memory.go      # 记忆管理
│   │   └── agent/             # 新增：智能体管理
│   │       ├── persona.go     # 角色定义
│   │       ├── llm.go         # LLM调用封装
│   │       └── prompt.go      # Prompt模板
│   ├── configs/
│   ├── pkg/
│   └── go.mod
├── frontend/                   # 新增：React前端
│   ├── src/
│   │   ├── pages/
│   │   │   ├── Lobby.tsx
│   │   │   ├── Room.tsx
│   │   │   ├── CreateRoom.tsx
│   │   │   └── Replay.tsx
│   │   ├── components/
│   │   │   ├── AgentPanel.tsx
│   │   │   ├── ChatStream.tsx
│   │   │   ├── ConsensusGauge.tsx
│   │   │   ├── PhaseIndicator.tsx
│   │   │   ├── RoomCard.tsx
│   │   │   └── Timeline.tsx
│   │   ├── hooks/
│   │   ├── stores/
│   │   ├── types/
│   │   └── App.tsx
│   ├── public/
│   ├── index.html
│   ├── package.json
│   └── vite.config.ts
├── docker-compose.yml          # 新增
├── Makefile                    # 更新
└── PRD.md / DEV_SPEC.md        # 本文档
```

---

## 10. 下一步

1. 编写DEV_SPEC开发规范
2. 重构后端目录（将现有代码移入backend/，新增simulation/agent模块）
3. 初始化前端项目（React+Vite+Tailwind）
4. 实现推演引擎核心（Phase轮转 + Agent LLM调用）
5. 实现WebSocket多房间协议
6. 开发前端UI（暗色主题 + 三栏布局）
7. 集成测试
8. Docker化部署
