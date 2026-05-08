package ws_protocol

import (
	"encoding/json"
	"time"

	"github.com/Jancapboy/Chatroom/backend/internal/model"
)

// ClientMessage 客户端 → 服务端消息
type ClientMessage struct {
	Type    string          `json:"type"` // user_message, command
	Payload json.RawMessage `json:"payload"`
}

type UserMessagePayload struct {
	Content string `json:"content"`
}

type CommandPayload struct {
	Command string `json:"command"` // pause, resume, next_phase, fork
}

// ServerMessage 服务端 → 客户端消息
type ServerMessage struct {
	Type      string          `json:"type"` // message, phase_change, agent_state, consensus_update, system
	RoomID    string          `json:"room_id"`
	Timestamp int64           `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

// MessagePayload 消息负载
type MessagePayload struct {
	ID           string `json:"id"`
	SenderID     string `json:"sender_id"`
	SenderType   string `json:"sender_type"` // agent, user, system
	SenderName   string `json:"sender_name"`
	SenderAvatar string `json:"sender_avatar,omitempty"`
	Content      string `json:"content"`
	Phase        string `json:"phase"`
	Round        int    `json:"round"`
	Confidence   int    `json:"confidence,omitempty"`
	Stance       string `json:"stance,omitempty"`
}

// PhaseChangePayload 阶段切换负载
type PhaseChangePayload struct {
	Phase       string `json:"phase"`
	Round       int    `json:"round"`
	PhaseName   string `json:"phase_name"`
	Description string `json:"description"`
}

// AgentStatePayload Agent状态负载
type AgentStatePayload struct {
	AgentID    string `json:"agent_id"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	Energy     int    `json:"energy"`
	Confidence int    `json:"confidence"`
	Stance     string `json:"stance"`
}

// ConsensusPayload 共识度负载
type ConsensusPayload struct {
	Topic       string                 `json:"topic"`
	Agreement   float64                `json:"agreement"`
	Breakdown   map[string]interface{} `json:"breakdown,omitempty"`
}

// SystemPayload 系统事件负载
type SystemPayload struct {
	Event   string `json:"event"`   // room_started, room_completed, agent_joined
	Message string `json:"message"`
}

// 便捷构造方法

func NewAgentMessage(roomID string, agent *model.RoomAgent, content, phase string, round, confidence int, stance string) *ServerMessage {
	payload := MessagePayload{
		ID:           agent.ID,
		SenderID:     agent.ID,
		SenderType:   "agent",
		SenderName:   agent.Name,
		SenderAvatar: agent.Avatar,
		Content:      content,
		Phase:        phase,
		Round:        round,
		Confidence:   confidence,
		Stance:       stance,
	}
	data, _ := json.Marshal(payload)
	return &ServerMessage{
		Type:      "message",
		RoomID:    roomID,
		Timestamp: time.Now().Unix(),
		Payload:   data,
	}
}

func NewUserMessage(roomID string, userID uint64, nickname, content, phase string, round int) *ServerMessage {
	payload := MessagePayload{
		ID:         string(rune(userID)),
		SenderID:   string(rune(userID)),
		SenderType: "user",
		SenderName: nickname,
		Content:    content,
		Phase:      phase,
		Round:      round,
	}
	data, _ := json.Marshal(payload)
	return &ServerMessage{
		Type:      "message",
		RoomID:    roomID,
		Timestamp: time.Now().Unix(),
		Payload:   data,
	}
}

func NewPhaseChangeMessage(roomID, phase string, round int) *ServerMessage {
	payload := PhaseChangePayload{
		Phase:       phase,
		Round:       round,
		PhaseName:   phase, // 简化，前端可自行映射中文名
		Description: "",
	}
	data, _ := json.Marshal(payload)
	return &ServerMessage{
		Type:      "phase_change",
		RoomID:    roomID,
		Timestamp: time.Now().Unix(),
		Payload:   data,
	}
}

func NewConsensusMessage(roomID, topic string, agreement float64) *ServerMessage {
	payload := ConsensusPayload{
		Topic:     topic,
		Agreement: agreement,
	}
	data, _ := json.Marshal(payload)
	return &ServerMessage{
		Type:      "consensus_update",
		RoomID:    roomID,
		Timestamp: time.Now().Unix(),
		Payload:   data,
	}
}

func NewSystemMessage(roomID, event, message string) *ServerMessage {
	payload := SystemPayload{
		Event:   event,
		Message: message,
	}
	data, _ := json.Marshal(payload)
	return &ServerMessage{
		Type:      "system",
		RoomID:    roomID,
		Timestamp: time.Now().Unix(),
		Payload:   data,
	}
}

func NewAgentStateMessage(roomID string, agent *model.RoomAgent) *ServerMessage {
	payload := AgentStatePayload{
		AgentID:    agent.ID,
		Name:       agent.Name,
		Role:       agent.Role,
		Energy:     agent.Energy,
		Confidence: agent.Confidence,
		Stance:     agent.Stance,
	}
	data, _ := json.Marshal(payload)
	return &ServerMessage{
		Type:      "agent_state",
		RoomID:    roomID,
		Timestamp: time.Now().Unix(),
		Payload:   data,
	}
}
