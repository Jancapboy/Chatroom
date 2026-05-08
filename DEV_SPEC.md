# DEV_SPEC — Chatroom-ASI v2.0
> 开发规范 v2.0 | 2026-05-08
> 核心关键词：模拟、辩论、变聪明、Skill、MCP、直接操作

---

## 1. 技术栈（锁定）

| 组件 | 版本 | 备注 |
|------|------|------|
| Go | 1.23+ | 后端主语言 |
| Gin | v1.10.0 | Web框架 |
| GORM | v1.25.12 | ORM |
| MySQL | 8.0+ | 主数据库 |
| Redis | 7.0+ | 缓存 + 临时数据 |
| React | 19.0 | 前端 |
| TypeScript | 5.7 | 前端类型 |
| Vite | 8.0 | 构建工具 |
| Tailwind CSS | 4.0 | 样式 |
| Zustand | 5.0 | 状态管理 |
| Lucide React | latest | 图标 |

**砍掉的技术**：
- Framer Motion → 不要动画，保持简洁
- shadcn/ui → 不用，手写极简组件
- WebSocket多房间 → 简化，单房间单连接

---

## 2. 后端架构

### 2.1 目录结构（精简版）

```
backend/
├── cmd/chatroom/main.go
├── configs/config.yaml
├── internal/
│   ├── routers/
│   │   ├── routers.go          # 路由注册
│   │   ├── websocket.go         # WS连接管理
│   │   └── api/
│   │       ├── user.go          # 用户（保留）
│   │       ├── room.go          # 房间管理
│   │       ├── agent.go         # Agent管理
│   │       ├── snapshot.go      # 快照/回溯
│   │       └── action.go        # 执行操作（新增）
│   ├── middleware/
│   │   ├── jwt.go               # JWT鉴权
│   │   └── cors.go              # CORS
│   ├── chat/
│   │   ├── broadcast.go         # 广播器
│   │   └── message.go           # 消息模型
│   ├── simulation/
│   │   ├── engine.go            # 推演引擎
│   │   ├── consensus.go         # 共识算法
│   │   ├── monitor.go           # 触发监控
│   │   └── skill.go             # **Skill系统** ⭐新增
│   ├── agent/
│   │   ├── llm_client.go        # LLM调用
│   │   ├── persona.go           # 角色定义
│   │   └── memory.go            # **Agent记忆** ⭐新增
│   ├── action/
│   │   ├── executor.go          # **Action执行** ⭐新增
│   │   └── mcp_client.go        # **MCP客户端** ⭐新增
│   ├── model/
│   │   ├── user.go
│   │   ├── room.go
│   │   ├── room_agent.go        # 房间-Agent关联
│   │   ├── message.go
│   │   ├── agent_template.go
│   │   ├── room_snapshot.go     # 快照
│   │   └── agent_memory.go      # **Agent记忆表** ⭐新增
│   ├── dao/
│   │   ├── user.go
│   │   ├── room.go
│   │   ├── message.go
│   │   ├── agent.go
│   │   ├── snapshot.go
│   │   └── memory.go            # **记忆DAO** ⭐新增
│   └── service/
│       ├── room.go
│       └── action.go            # **Action服务** ⭐新增
├── pkg/
│   ├── response/
│   │   └── response.go          # API响应格式
│   ├── errcode/
│   │   └── errcode.go           # 错误码
│   └── ws_protocol/
│       └── protocol.go          # WS消息协议
└── global/
    ├── db.go
    └── settings.go
```

### 2.2 核心接口定义

#### Skill 接口（新增）
```go
type Skill interface {
    Name() string                    // Skill名称
    Description() string             // 描述
    CanHandle(topic string) bool     // 是否能处理该话题
    EnhancePrompt(basePrompt string, context map[string]interface{}) string  // 增强System Prompt
    AfterDebate(agentID string, messages []Message, consensus float64) map[string]interface{}  // 辩论后提取经验
}
```

#### AgentMemory 模型（新增）
```go
type AgentMemory struct {
    ID         string    `gorm:"primaryKey"`
    AgentID    string    `gorm:"index"`          // Agent唯一标识（跨房间）
    RoomID     string    `gorm:"index"`          // 来源房间
    Topic      string    `gorm:"index"`          // 话题标签（如"微服务架构"）
    Experience string    `gorm:"type:text"`      // 经验文本
    Outcome    string    `gorm:"type:varchar(32)"` // "success" / "failure" / "neutral"
    Confidence int       // 置信度 0-100
    CreatedAt  time.Time
}
```

#### Action 执行器（新增）
```go
type Action interface {
    Name() string
    Description() string
    Validate(params map[string]interface{}) error
    Execute(params map[string]interface{}) (string, error)  // 返回执行结果
}

// 内置Actions
type FeishuMessageAction struct{}
type FeishuTaskAction struct{}
type EmailAction struct{}
type MarkdownReportAction struct{}
```

### 2.3 数据流规范

1. **推演消息流**：Agent发言 → WS广播 → 保存到MySQL
2. **学习提取流**：一轮结束 → 调用LLM分析 → 提取经验 → 保存到AgentMemory
3. **Action流**：推演完成 → 生成Action草稿 → 用户确认 → 调用Execute → 记录结果

### 2.4 API 响应格式（统一）

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

---

## 3. 前端架构

### 3.1 目录结构（极简版）

```
frontend/src/
├── main.tsx
├── App.tsx
├── index.css                    # Tailwind入口
├── types/
│   ├── api.ts                   # API类型
│   ├── room.ts                  # 房间/Agent/消息类型
│   ├── ws.ts                    # WS消息类型
│   └── skill.ts                 # **Skill类型** ⭐新增
├── lib/
│   ├── api.ts                   # Axios封装
│   └── utils.ts                 # 工具函数
├── stores/
│   ├── useLobbyStore.ts         # 大厅状态
│   └── useRoomStore.ts          # 房间状态
├── hooks/
│   └── useWebSocket.ts          # WS Hook
├── components/
│   ├── layout/
│   │   └── MainLayout.tsx       # 主布局
│   ├── lobby/
│   │   ├── RoomList.tsx
│   │   └── CreateRoomForm.tsx
│   └── room/
│       ├── ChatStream.tsx       # 聊天流（核心）
│       ├── RoomControls.tsx     # 底部控制条
│       ├── RoundIndicator.tsx   # 轮次指示（极简）
│       ├── ConsensusBadge.tsx   # 共识度徽章（极简）
│       ├── AgentList.tsx        # Agent列表（折叠）
│       └── ActionPanel.tsx      # **Action面板** ⭐新增（推演完成后显示）
└── pages/
    ├── Lobby.tsx
    ├── Room.tsx
    └── CreateRoom.tsx
```

### 3.2 界面原则

**默认显示（80%屏幕）**：
- 聊天流（Agent对话）
- 底部：输入框 + 开始/暂停按钮
- 顶部：房间名 + 轮次 + 共识度（小字）

**折叠到侧边栏（20%）**：
- Agent列表（hover展开）
- Skill配置
- 检查点列表
- MCP工具调用记录

**推演完成后弹出**：
- 结论摘要
- Action草稿（如"发飞书消息通知团队"）
- 按钮：执行 / 修改 / 跳过

---

## 4. 数据库表结构

### 4.1 新增表

#### agent_memories（Agent长期记忆）
```sql
CREATE TABLE agent_memories (
    id VARCHAR(36) PRIMARY KEY,
    agent_id VARCHAR(36) NOT NULL,
    room_id VARCHAR(36),
    topic VARCHAR(128),          -- 话题标签，如"微服务架构"
    experience TEXT,              -- 经验文本
    outcome VARCHAR(32),          -- success/failure/neutral
    confidence INT,               -- 0-100
    created_at DATETIME,
    INDEX idx_agent_topic (agent_id, topic)
);
```

#### actions（执行操作记录）
```sql
CREATE TABLE actions (
    id VARCHAR(36) PRIMARY KEY,
    room_id VARCHAR(36) NOT NULL,
    type VARCHAR(32),            -- feishu_message / email / report
    status VARCHAR(32),          -- pending / approved / executed / failed
    params JSON,                   -- 执行参数
    result TEXT,                   -- 执行结果
    created_by VARCHAR(36),        -- 谁发起的
    approved_by VARCHAR(36),       -- 谁确认的
    created_at DATETIME,
    executed_at DATETIME
);
```

### 4.2 修改表

#### rooms（新增字段）
```sql
ALTER TABLE rooms ADD COLUMN skills JSON AFTER max_rounds;        -- 加载的Skill列表
ALTER TABLE rooms ADD COLUMN mcp_tools JSON AFTER skills;          -- 可用的MCP工具
ALTER TABLE rooms ADD COLUMN conclusion TEXT AFTER consensus_score; -- 推演结论
```

---

## 5. 任务分工

### 5.1 诺亚（主会话）负责

**架构层 + 核心逻辑**：
1. ✅ PRD v2.0 编写（刚完成）
2. ✅ DEV_SPEC v2.0 编写（刚完成）
3. 🔄 **Agent学习系统（Memory）**：
   - 设计AgentMemory模型
   - 实现经验提取算法（调用LLM分析一轮辩论）
   - 实现记忆检索（下次推演时注入上下文）
4. 🔄 **Skill框架**：
   - 设计Skill接口
   - 实现3个内置Skill（debate/calculation/risk_assessment）
   - 集成到Engine中
5. 🔄 **Action执行器**：
   - 设计Action接口
   - 实现FeishuMessageAction
   - 实现MarkdownReportAction
6. 🔍 **代码Review**：所有子智能体提交的代码
7. 🔍 **集成测试**：前后端联调

### 5.2 子智能体负责

**具体实现 + UI开发**：
1. 🔄 **前端界面简化**：
   - 重写Room.tsx（极简版）
   - 重写ChatStream（聚焦对话）
   - 删除多余动画和面板
2. 🔄 **前端Skill面板**：
   - 创建房间时选择Skill
   - 显示Agent加载的Skill
3. 🔄 **前端Action面板**：
   - 推演完成后显示Action草稿
   - "确认执行"交互
4. 🔄 **后端API实现**：
   - Agent学习相关API（/agents/:id/memories）
   - Action相关API（/rooms/:id/actions）
   - Skill相关API（/skills, /agents/:id/skills）
5. 🔄 **MCP客户端基础**：
   - 实现MCP连接框架
   - 配置2个示例MCP Server
6. 🔄 **测试**：单元测试 + E2E测试

### 5.3 协作流程

```
1. 诺亚写PRD/DEV_SPEC → 明确接口和数据结构
2. 子智能体读文档 → 实现具体模块
3. 子智能体提交代码 → 诺亚Review
4. 诺亚做关键集成（Memory注入、Skill加载、Action触发）
5. 联合测试 → 修复 → 合并
```

---

## 6. 开发顺序

### Phase 1：核心推演（本周）
- [ ] 简化前端UI（子智能体）
- [ ] AgentMemory模型 + 表（诺亚）
- [ ] 经验提取算法（诺亚）
- [ ] Skill框架接口（诺亚）
- [ ] 3个内置Skill实现（子智能体）

### Phase 2：学习与进化（下周）
- [ ] 记忆检索注入（诺亚）
- [ ] Agent学习效果验证
- [ ] 前端显示"Agent已学习N条经验"

### Phase 3：Action执行（再下周）
- [ ] Action接口 + FeishuMessageAction（诺亚）
- [ ] 前端Action面板（子智能体）
- [ ] MCP基础框架（子智能体）

### Phase 4： polish（最后）
- [ ] 报告导出
- [ ] 多模型支持
- [ ] 性能优化

---

## 7. 代码规范

### 7.1 Go
- 接口定义放在 `internal/` 对应包的根目录
- 模型放在 `internal/model/`
- 业务逻辑放在 `internal/service/`
- HTTP handler放在 `internal/routers/api/`
- 错误用 `pkg/errcode` 统一封装

### 7.2 React
- 组件用函数式 + Hooks
- 状态用 Zustand，不用 useState 跨组件传递
- API 调用集中在 `lib/api.ts`
- 类型定义在 `types/` 目录

### 7.3 提交规范
```
feat: 新增Agent学习系统
fix: 修复共识度计算bug
docs: 更新PRD
refactor: 简化前端UI
```
