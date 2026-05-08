# Chatroom 3D 集成

在原有聊天室基础上集成腾讯云混元3D模型生成功能，用户可以在聊天中输入描述文字，生成3D模型并在浏览器内实时预览。

## 技术架构

```
┌─────────────┐     HTTP/WS     ┌──────────────┐     API      ┌──────────────┐
│  React 前端  │ ◄────────────► │  Go + Gin 后端│ ◄──────────► │ 腾讯云混元3D  │
│  Three.js   │                 │  WebSocket   │              │ GLB模型生成   │
└─────────────┘                 └──────────────┘              └──────────────┘
```

- **后端**: Go 1.20+ / Gin / GORM / nhooyr/websocket
- **前端**: React / TypeScript / Three.js 0.160.0 / Ant Design
- **3D生成**: 腾讯云混元3D专业版 API（`ai3d.tencentcloudapi.com`）
- **数据库**: MySQL 5.7+

## 快速部署

### 1. 环境要求

- Go 1.20+
- Node.js 16+
- MySQL 5.7+ (或 MariaDB 10.5+)
- 腾讯云账号（开通混元3D服务）

### 2. 后端启动

```bash
git clone git@github.com:Jancapboy/Chatroom.git
cd Chatroom
git checkout feature/3d-integration

# 设置环境变量
export TENCENT_SECRET_ID=你的SecretId
export TENCENT_SECRET_KEY=你的SecretKey

# 修改 configs/config.yaml 中的数据库连接
# 启动
go run ./cmd/chatroom/main.go
```

### 3. 前端启动

```bash
# 克隆前端仓库
git clone https://github.com/UncleBloom/chatroom.git
cd chatroom

# 应用3D补丁（复制 frontend-patch/ 中的文件到 src/ 对应目录）
# 详见 frontend-patch/README.md

# 安装依赖
npm install
npm install three@0.160.0 @types/three@0.160.0

# 修改 src/api/hostname.ts 指向后端地址
# 启动
npm start
```

### 4. 独立测试（不需要MySQL）

```bash
go run ./cmd/test_3d_server/main.go
# 访问 http://localhost:4002/health
```

## API 文档

### POST /api/3d/generate

提交3D模型生成任务。

**请求**:
```json
{
  "prompt": "一只可爱的橘猫",
  "image_url": "",
  "result_format": "GLB"
}
```

**响应**:
```json
{
  "code": 0,
  "job_id": "1429320300369436672"
}
```

### POST /api/3d/query

查询生成任务状态。

**请求**:
```json
{
  "job_id": "1429320300369436672"
}
```

**响应**:
```json
{
  "code": 0,
  "status": "DONE",
  "files": [
    {
      "type": "GLB",
      "url": "https://hunyuan-prod-xxx.cos.ap-guangzhou.tencentcos.cn/xxx.glb?..."
    }
  ]
}
```

**status 取值**: `WAIT` → `RUN` → `DONE` / `FAIL`

### WebSocket 消息格式

```json
// 文本消息
{"type": "text", "message_content": "你好"}

// 3D模型消息
{"type": "3d_model", "message_content": "一只橘猫", "model_url": "https://xxx.glb"}
```

## 前端3D组件

| 组件 | 路径 | 功能 |
|------|------|------|
| ModelViewer | `components/model/modelViewer.tsx` | Three.js 渲染GLB模型（旋转/缩放/灯光） |
| Generate3D | `components/model/generate3D.tsx` | 生成按钮 + 弹窗 + 进度条 |
| threeD API | `api/threeD.ts` | 提交任务 + 10秒间隔轮询 |

## 测试

```bash
# 后端API + 前端编译 + Git状态检查
python3 scripts/e2e_test.py

# 需要后端运行在 localhost:4002
```

## 注意事项

- **并发限制**: 混元3D默认1个并发任务，超出返回 `RequestLimitExceeded`
- **生成耗时**: 约150秒（2.5分钟）
- **输出格式**: 默认OBJ，需指定 `result_format: "GLB"` 给Three.js用
- **Three.js版本**: 必须用 0.160.0，新版需要 TypeScript 5.0+
- **凭证安全**: SecretKey 通过环境变量传入，不要硬编码
- **API版本**: `SubmitHunyuanTo3DJob`（专业版），不是极速版
- **Region**: 仅支持 `ap-guangzhou`
