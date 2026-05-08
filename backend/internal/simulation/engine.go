package simulation

import (
	"context"
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
	roomID      string
	room        *model.Room
	agents      []model.RoomAgent
	phaseCtrl   *PhaseController
	consensus   *ConsensusEngine
	llmClient   *agent.LLMClient
	broadcast   chan *ws_protocol.ServerMessage
	state       string // running / paused / completed
	stopCh      chan struct{}
	mu          sync.RWMutex
	dao         *dao.Dao
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
	return &Engine{
		roomID:    room.ID,
		room:      room,
		agents:    agents,
		phaseCtrl: NewPhaseController(),
		consensus: NewConsensusEngine(),
		llmClient: agent.NewLLMClient(),
		broadcast: broadcast,
		state:     "preparing",
		stopCh:    make(chan struct{}),
		dao:       dao.New(global.DBEngine),
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
	if err := e.dao.MessageCreate(msg); err != nil {
		log.Printf("[Engine] 保存消息失败: %v", err)
	}
}
