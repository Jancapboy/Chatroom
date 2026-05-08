package simulation

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Jancapboy/Chatroom/backend/global"
	"github.com/Jancapboy/Chatroom/backend/internal/agent"
	"github.com/Jancapboy/Chatroom/backend/internal/dao"
	"github.com/Jancapboy/Chatroom/backend/internal/model"
	"github.com/Jancapboy/Chatroom/backend/pkg/ws_protocol"
	"github.com/google/uuid"
)

// Engine 推演引擎
type Engine struct {
	roomID         string
	room           *model.Room
	agents         []model.RoomAgent
	phaseCtrl      *PhaseController
	consensus      *ConsensusEngine
	llmClient      *agent.LLMClient
	learningEngine *LearningEngine
	broadcast      chan *ws_protocol.ServerMessage
	state          string // running / paused / completed
	stopCh         chan struct{}
	mu             sync.RWMutex
	dao            *dao.Dao
	monitor        *RollbackMonitor
}

// RollbackTrigger 回溯触发配置
type RollbackTrigger struct {
	ConsensusThreshold int  // 共识度低于此值触发
	RiskOfficerVeto    bool // 风险官强烈反对时触发
	UserFlag           bool // 用户标记时触发
}

// RollbackMonitor 指标检测器
type RollbackMonitor struct {
	trigger RollbackTrigger
}

// NewRollbackMonitor 创建监控器
func NewRollbackMonitor(trigger RollbackTrigger) *RollbackMonitor {
	return &RollbackMonitor{trigger: trigger}
}

// Check 检查是否需要触发快照
func (m *RollbackMonitor) Check(consensus float64, agents []model.RoomAgent) (bool, string) {
	// 1. 共识度低于阈值
	if consensus < float64(m.trigger.ConsensusThreshold) {
		return true, fmt.Sprintf("共识度过低 (%.1f%% < %d%%)", consensus, m.trigger.ConsensusThreshold)
	}
	// 2. 风险官强烈反对
	if m.trigger.RiskOfficerVeto {
		for _, a := range agents {
			if a.Role == "risk_officer" && a.Stance == "oppose" && a.Confidence > 80 {
				return true, fmt.Sprintf("风险官强烈反对 (confidence=%d)", a.Confidence)
			}
		}
	}
	return false, ""
}

// EngineManager 全局引擎管理器
type EngineManager struct {
	engines map[string]*Engine
	mu      sync.RWMutex
}

var GlobalEngineManager = &EngineManager{
	engines: make(map[string]*Engine),
}

func (em *EngineManager) GetOrCreate(room *model.Room, agents []model.RoomAgent, broadcast chan *ws_protocol.ServerMessage) *Engine {
	em.mu.Lock()
	defer em.mu.Unlock()

	if e, ok := em.engines[room.ID]; ok {
		return e
	}

	e := NewEngine(room, agents, broadcast)
	em.engines[room.ID] = e
	return e
}

func (em *EngineManager) Get(roomID string) *Engine {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return em.engines[roomID]
}

func (em *EngineManager) Remove(roomID string) {
	em.mu.Lock()
	defer em.mu.Unlock()
	delete(em.engines, roomID)
}

// NewEngine 创建推演引擎
func NewEngine(room *model.Room, agents []model.RoomAgent, broadcast chan *ws_protocol.ServerMessage) *Engine {
	llmClient := agent.NewLLMClient()
	return &Engine{
		roomID:         room.ID,
		room:           room,
		agents:         agents,
		phaseCtrl:      NewPhaseController(),
		consensus:      NewConsensusEngine(),
		llmClient:      llmClient,
		learningEngine: NewLearningEngine(llmClient),
		broadcast:      broadcast,
		state:          "preparing",
		stopCh:         make(chan struct{}),
		dao:            dao.New(global.DBEngine),
		monitor:        NewRollbackMonitor(RollbackTrigger{ConsensusThreshold: 30, RiskOfficerVeto: true}),
	}
}

// Run 启动推演主循环
func (e *Engine) Run() {
	e.mu.Lock()
	if e.state == "running" {
		e.mu.Unlock()
		return
	}
	e.state = "running"
	e.mu.Unlock()

	log.Printf("[Engine] 房间 %s 推演启动", e.roomID)

	// 广播房间开始
	e.broadcast <- ws_protocol.NewSystemMessage(e.roomID, "room_started", fmt.Sprintf("房间 '%s' 推演开始", e.room.Name))

	for e.room.CurrentRound <= e.room.MaxRounds {
		select {
		case <-e.stopCh:
			e.state = "paused"
			log.Printf("[Engine] 房间 %s 推演暂停", e.roomID)
			return
		default:
		}

		if e.state != "running" {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		log.Printf("[Engine] 房间 %s 第 %d 轮开始", e.roomID, e.room.CurrentRound)

		// Phase 1-6 循环
		for _, phase := range e.phaseCtrl.AllPhases() {
			select {
			case <-e.stopCh:
				e.state = "paused"
				return
			default:
			}

			e.room.CurrentPhase = string(phase)
			e.phaseCtrl.Enter(phase)

			// 广播阶段切换
			e.broadcast <- ws_protocol.NewPhaseChangeMessage(e.roomID, string(phase), e.room.CurrentRound)

			// 保存阶段切换消息
			e.saveMessage(&model.Message{
				ID:         uuid.New().String(),
				RoomID:     e.roomID,
				SenderID:   "system",
				SenderType: "system",
				SenderName: "System",
				Content:    fmt.Sprintf("进入阶段: %s", phase),
				MsgType:    "phase_change",
				Phase:      string(phase),
				Round:      e.room.CurrentRound,
			})

			// 每个Agent按序发言
			for i := range e.agents {
				agentInst := &e.agents[i]
				if !agentInst.IsActive {
					continue
				}

				select {
				case <-e.stopCh:
					e.state = "paused"
					return
				default:
				}

				// 构建上下文
				ctxMsgs := e.buildContextMessages(agentInst.ID, 10)
				persona := agent.GetPersonaByRole(agentInst.Role)
				if persona == nil {
					persona = &agent.Persona{
						Name:         agentInst.Name,
						Role:         agentInst.Role,
						SystemPrompt: agentInst.SystemPrompt,
					}
				}

				// 注入学习记忆（变聪明！）
				getMemories := func(agentID, topic string, limit int) ([]model.AgentMemory, error) {
					mems, err := e.dao.MemoryListByAgentAndTopic(agentID, topic, limit)
					if err != nil {
						return nil, err
					}
					return mems, nil
				}
				enhancedPrompt, err := e.learningEngine.GenerateLearningPrompt(
					context.Background(),
					agentInst.ID,
					persona.SystemPrompt,
					e.room.Topic,
					getMemories,
				)
				if err == nil && enhancedPrompt != persona.SystemPrompt {
					persona.SystemPrompt = enhancedPrompt
					log.Printf("[Engine] Agent %s 已注入历史记忆", agentInst.Name)
				}

				// 调用LLM
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				resp, err := e.llmClient.CompleteWithContext(ctx, persona, e.room.Topic, string(phase), e.room.CurrentRound, e.room.MaxRounds, ctxMsgs)
				cancel()

				if err != nil {
					log.Printf("[Engine] Agent %s LLM调用失败: %v", agentInst.Name, err)
					// 降级：使用预设回复
					resp = &agent.AgentResponse{
						Content:    fmt.Sprintf("[%s] 我需要更多时间来分析这个问题...", agentInst.Name),
						Stance:     agentInst.Stance,
						Confidence: agentInst.Confidence,
					}
				}

				// 更新Agent状态
				agentInst.Stance = resp.Stance
				agentInst.Confidence = resp.Confidence
				_ = e.dao.AgentUpdate(agentInst)

				// 保存消息
				msg := &model.Message{
					ID:           uuid.New().String(),
					RoomID:       e.roomID,
					SenderID:     agentInst.ID,
					SenderType:   "agent",
					SenderName:   agentInst.Name,
					SenderAvatar: agentInst.Avatar,
					Content:      resp.Content,
					MsgType:      "text",
					Phase:        string(phase),
					Round:        e.room.CurrentRound,
					Metadata:     fmt.Sprintf(`{"confidence":%d,"stance":"%s"}`, resp.Confidence, resp.Stance),
				}
				e.saveMessage(msg)

				// 广播消息
				e.broadcast <- ws_protocol.NewAgentMessage(e.roomID, agentInst, resp.Content, string(phase), e.room.CurrentRound, resp.Confidence, resp.Stance)

				// 模拟思考间隔 2秒
				time.Sleep(2 * time.Second)
			}

			// 阶段结束：计算共识度
			phaseMessages, _ := e.dao.GetPhaseMessages(e.roomID, string(phase), e.room.CurrentRound)
			consensus := e.consensus.Calculate(phaseMessages, e.agents)
			e.room.ConsensusScore = int(consensus)

			// 广播共识度更新
			e.broadcast <- ws_protocol.NewConsensusMessage(e.roomID, string(phase), consensus)

			// 广播共识消息
			if consensus > 0 {
				e.saveMessage(&model.Message{
					ID:         uuid.New().String(),
					RoomID:     e.roomID,
					SenderID:   "system",
					SenderType: "system",
					SenderName: "System",
					Content:    fmt.Sprintf("阶段共识度: %.1f%%", consensus),
					MsgType:    "consensus",
					Phase:      string(phase),
					Round:      e.room.CurrentRound,
					Metadata:   fmt.Sprintf(`{"consensus":%.1f}`, consensus),
				})
			}

			// 更新房间状态
			_ = e.dao.RoomUpdate(e.room)

			// 检查是否需要暂停（如共识度过低在debate阶段）
			if consensus < 30 && phase == "debate" {
				log.Printf("[Engine] 房间 %s 共识度过低(%.1f%%)，继续推演", e.roomID, consensus)
			}

			// Phase间短暂停顿
			time.Sleep(1 * time.Second)
		}

		// 一轮结束：创建快照（总是保存最新状态）
		triggered, reason := e.monitor.Check(float64(e.room.ConsensusScore), e.agents)
		if triggered {
			log.Printf("[Engine] 房间 %s 触发自动快照: %s", e.roomID, reason)
		}
		snapshotReason := "每轮自动保存"
		if triggered {
			snapshotReason = reason
		}
		e.createSnapshot(e.room.CurrentRound, e.room.CurrentPhase, e.room.ConsensusScore, snapshotReason)

		// 一轮结束：Agent学习（提取经验）
		e.learnFromRound(e.room.CurrentRound, e.room.CurrentPhase, e.room.ConsensusScore)

		// 一轮结束
		e.room.CurrentRound++
		_ = e.dao.RoomUpdate(e.room)
	}

	// 推演完成
	e.mu.Lock()
	e.state = "completed"
	e.mu.Unlock()

	e.room.Status = "completed"
	_ = e.dao.RoomUpdate(e.room)

	e.broadcast <- ws_protocol.NewSystemMessage(e.roomID, "room_completed", fmt.Sprintf("房间 '%s' 推演完成", e.room.Name))
	log.Printf("[Engine] 房间 %s 推演完成", e.roomID)
}

// Pause 暂停推演
func (e *Engine) Pause() {
	close(e.stopCh)
	e.stopCh = make(chan struct{})
}

// Resume 恢复推演
func (e *Engine) Resume() {
	e.mu.Lock()
	if e.state == "paused" || e.state == "preparing" {
		e.state = "running"
		e.mu.Unlock()
		go e.Run()
		return
	}
	e.mu.Unlock()
}

// GetState 获取引擎状态
func (e *Engine) GetState() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.state
}

// HandleUserMessage 处理用户消息，Agent下一轮回应
func (e *Engine) HandleUserMessage(userID uint64, nickname, content string) {
	// 保存用户消息
	msg := &model.Message{
		ID:         uuid.New().String(),
		RoomID:     e.roomID,
		SenderID:   fmt.Sprintf("%d", userID),
		SenderType: "user",
		SenderName: nickname,
		Content:    content,
		MsgType:    "text",
		Phase:      e.room.CurrentPhase,
		Round:      e.room.CurrentRound,
	}
	e.saveMessage(msg)

	// 广播用户消息
	e.broadcast <- ws_protocol.NewUserMessage(e.roomID, userID, nickname, content, e.room.CurrentPhase, e.room.CurrentRound)

	// 注：Agent的回应由推演引擎在下一轮自动处理（因为用户消息已存入上下文）
	log.Printf("[Engine] 房间 %s 收到用户 %s 消息", e.roomID, nickname)
}

// buildContextMessages 构建最近N条消息的上下文文本
func (e *Engine) buildContextMessages(excludeAgentID string, limit int) string {
	msgs, _, _ := e.dao.MessageListByRoom(e.roomID, 0, "", 1, limit)
	if len(msgs) == 0 {
		return "（暂无历史消息）"
	}

	var sb strings.Builder
	for i := len(msgs) - 1; i >= 0; i-- {
		m := msgs[i]
		if m.SenderID == excludeAgentID {
			continue
		}
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", m.SenderType, m.SenderName, m.Content))
	}

	result := sb.String()
	if result == "" {
		return "（暂无历史消息）"
	}
	return result
}

func (e *Engine) saveMessage(msg *model.Message) {
	if msg.Metadata == "" {
		msg.Metadata = "{}"
	}
	if err := e.dao.MessageCreate(msg); err != nil {
		log.Printf("[Engine] 保存消息失败: %v", err)
	}
}

// createSnapshot 创建房间状态快照
func (e *Engine) createSnapshot(round int, phase string, consensusScore int, reason string) {
	agentStates := make([]map[string]interface{}, 0, len(e.agents))
	for _, a := range e.agents {
		agentStates = append(agentStates, map[string]interface{}{
			"agent_id":   a.ID,
			"name":       a.Name,
			"role":       a.Role,
			"stance":     a.Stance,
			"confidence": a.Confidence,
			"energy":     a.Energy,
		})
	}
	agentStatesJSON, _ := json.Marshal(agentStates)

	keyDecisions := make([]map[string]interface{}, 0)
	for _, a := range e.agents {
		if a.Stance != "" && a.Stance != "neutral" {
			keyDecisions = append(keyDecisions, map[string]interface{}{
				"agent_id": a.ID,
				"content":  fmt.Sprintf("%s takes stance: %s", a.Name, a.Stance),
				"stance":   a.Stance,
			})
		}
	}
	keyDecisionsJSON, _ := json.Marshal(keyDecisions)

	snapshot := &model.RoomSnapshot{
		ID:             uuid.New().String(),
		RoomID:         e.roomID,
		Round:          round,
		Phase:          phase,
		ConsensusScore: consensusScore,
		AgentStates:    string(agentStatesJSON),
		KeyDecisions:   string(keyDecisionsJSON),
		TriggerReason:  reason,
		CreatedAt:      time.Now(),
	}
	if err := e.dao.SnapshotCreate(snapshot); err != nil {
		log.Printf("[Engine] 创建快照失败: %v", err)
	} else {
		log.Printf("[Engine] 房间 %s 第 %d 轮快照已创建", e.roomID, round)
		// 广播快照创建事件
		e.broadcast <- ws_protocol.NewSystemMessage(e.roomID, "snapshot_created",
			fmt.Sprintf("第 %d 轮检查点已创建: %s", round, reason))
	}
}

// learnFromRound 从一轮辩论中提取学习经验
func (e *Engine) learnFromRound(round int, phase string, consensusScore int) {
	if e.learningEngine == nil {
		return
	}

	// 获取本轮消息
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	messages, _, _ := e.dao.MessageListByRoom(e.roomID, 0, "", 1, 50)
	if len(messages) == 0 {
		return
	}

	for i := range e.agents {
		agentInst := &e.agents[i]
		if !agentInst.IsActive {
			continue
		}

		memory, err := e.learningEngine.ExtractExperience(
			ctx,
			agentInst.ID,
			agentInst.Name,
			agentInst.Role,
			e.room.Topic,
			messages,
			float64(consensusScore),
		)
		if err != nil {
			log.Printf("[Engine] Agent %s 经验提取失败: %v", agentInst.Name, err)
			continue
		}
		if memory == nil {
			continue
		}

		// 保存记忆
		memory.RoomID = e.roomID
		if err := e.dao.MemoryCreate(memory); err != nil {
			log.Printf("[Engine] Agent %s 记忆保存失败: %v", agentInst.Name, err)
			continue
		}

		log.Printf("[Engine] Agent %s 学习完成: %s", agentInst.Name, memory.Experience)
	}
}
