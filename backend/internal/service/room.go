package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Jancapboy/Chatroom/backend/internal/model"
	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
	"github.com/google/uuid"
)

// RoomCreateRequest 创建房间请求
type RoomCreateRequest struct {
	Name        string   `json:"name" binding:"required"`
	Topic       string   `json:"topic"`
	Description string   `json:"description"`
	TemplateID  string   `json:"template_id"`
	MaxRounds   int      `json:"max_rounds"`
	AgentIDs    []string `json:"agent_ids"` // 选择的Agent模板ID
}

// RoomResponse 房间响应
type RoomResponse struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Topic          string           `json:"topic"`
	Description    string           `json:"description"`
	Status         string           `json:"status"`
	CurrentPhase   string           `json:"current_phase"`
	CurrentRound   int              `json:"current_round"`
	MaxRounds      int              `json:"max_rounds"`
	ConsensusScore int              `json:"consensus_score"`
	CreatedBy      uint64           `json:"created_by"`
	CreatedAt      time.Time        `json:"created_at"`
	Agents         []model.RoomAgent `json:"agents,omitempty"`
}

func buildRoomResponse(r *model.Room) RoomResponse {
	createdAt := r.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	return RoomResponse{
		ID:             r.ID,
		Name:           r.Name,
		Topic:          r.Topic,
		Description:    r.Description,
		Status:         r.Status,
		CurrentPhase:   r.CurrentPhase,
		CurrentRound:   r.CurrentRound,
		MaxRounds:      r.MaxRounds,
		ConsensusScore: r.ConsensusScore,
		CreatedBy:      r.CreatedBy,
		CreatedAt:      createdAt,
		Agents:         r.Agents,
	}
}

// RoomCreate 创建房间
func (svc *Service) RoomCreate(req *RoomCreateRequest, userID uint64) (*RoomResponse, *errcode.Error) {
	room := &model.Room{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Topic:        req.Topic,
		Description:  req.Description,
		Status:       "preparing",
		CurrentPhase: "info_gathering",
		CurrentRound: 1,
		MaxRounds:    req.MaxRounds,
		CreatedBy:    userID,
	}
	if room.MaxRounds <= 0 {
		room.MaxRounds = 10
	}
	if req.TemplateID != "" {
		room.TemplateID = &req.TemplateID
	}

	err := svc.dao.RoomCreate(room)
	if err != nil {
		return nil, err
	}

	// 自动添加预设Agent
	if len(req.AgentIDs) > 0 {
		for _, tplID := range req.AgentIDs {
			tpl, err := svc.dao.TemplateGetByID(tplID)
			if err != nil {
				continue
			}
			agent := &model.RoomAgent{
				ID:           uuid.New().String(),
				RoomID:       room.ID,
				TemplateID:   &tplID,
				Name:         tpl.Name,
				Role:         tpl.Role,
				Personality:  tpl.Personality,
				Expertise:    tpl.Expertise,
				Model:        tpl.DefaultModel,
				SystemPrompt: buildSystemPrompt(tpl.SystemPromptTemplate, req.Topic),
				Energy:       100,
				Confidence:   50,
				Stance:       "neutral",
				IsActive:     true,
			}
			if agent.Model == "" {
				agent.Model = "deepseek-chat"
			}
			_ = svc.dao.AgentCreate(agent)
		}
	}

	return &RoomResponse{ID: room.ID, Name: room.Name, Status: room.Status}, nil
}

// RoomGet 获取房间详情
func (svc *Service) RoomGet(id string) (*RoomResponse, *errcode.Error) {
	room, err := svc.dao.RoomGet(id)
	if err != nil {
		return nil, err
	}
	resp := buildRoomResponse(room)
	return &resp, nil
}

// RoomList 房间列表
func (svc *Service) RoomList(status string, page, pageSize int) ([]RoomResponse, int64, *errcode.Error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	rooms, total, err := svc.dao.RoomList(status, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]RoomResponse, 0, len(rooms))
	for _, r := range rooms {
		resp = append(resp, buildRoomResponse(&r))
	}
	return resp, total, nil
}

// GetRoomAgents 获取房间内的Agent列表
func (svc *Service) GetRoomAgents(roomID string) ([]model.RoomAgent, *errcode.Error) {
	return svc.dao.AgentListByRoom(roomID)
}

// RoomStart 开始推演
func (svc *Service) RoomStart(id string) *errcode.Error {
	room, err := svc.dao.RoomGet(id)
	if err != nil {
		return err
	}
	if room.Status != "preparing" && room.Status != "paused" {
		return errcode.NewError(20030001, "房间状态不允许开始推演")
	}
	room.Status = "running"
	return svc.dao.RoomUpdate(room)
}

// RoomPause 暂停推演
func (svc *Service) RoomPause(id string) *errcode.Error {
	room, err := svc.dao.RoomGet(id)
	if err != nil {
		return err
	}
	if room.Status != "running" {
		return errcode.NewError(20030002, "房间未在运行中")
	}
	room.Status = "paused"
	return svc.dao.RoomUpdate(room)
}

// RoomDelete 删除房间
func (svc *Service) RoomDelete(id string) *errcode.Error {
	return svc.dao.RoomDelete(id)
}

// AddAgentToRoom 向房间添加Agent
func (svc *Service) AddAgentToRoom(roomID string, templateID string) (*model.RoomAgent, *errcode.Error) {
	tpl, err := svc.dao.TemplateGetByID(templateID)
	if err != nil {
		return nil, err
	}
	room, err := svc.dao.RoomGet(roomID)
	if err != nil {
		return nil, err
	}

	agent := &model.RoomAgent{
		ID:           uuid.New().String(),
		RoomID:       roomID,
		TemplateID:   &templateID,
		Name:         tpl.Name,
		Role:         tpl.Role,
		Personality:  tpl.Personality,
		Expertise:    tpl.Expertise,
		Model:        tpl.DefaultModel,
		SystemPrompt: buildSystemPrompt(tpl.SystemPromptTemplate, room.Topic),
		Energy:       100,
		Confidence:   50,
		Stance:       "neutral",
		IsActive:     true,
	}
	if agent.Model == "" {
		agent.Model = "deepseek-chat"
	}

	err = svc.dao.AgentCreate(agent)
	if err != nil {
		return nil, err
	}
	return agent, nil
}

// buildSystemPrompt 替换模板变量
func buildSystemPrompt(template, topic string) string {
	if template == "" {
		return fmt.Sprintf("你是一位智能体，正在参与关于'%s'的讨论。请发表你的观点。", topic)
	}
	return fmt.Sprintf(template, topic)
}

// SnapshotList 获取房间快照列表
func (svc *Service) SnapshotList(roomID string) ([]model.RoomSnapshot, *errcode.Error) {
	return svc.dao.SnapshotListByRoom(roomID)
}

// RoomRollback 回滚到指定快照（创建新分支房间）
func (svc *Service) RoomRollback(roomID, snapshotID string, userID uint64) (*RoomResponse, *errcode.Error) {
	// 1. 找到目标快照
	snapshot, err := svc.dao.SnapshotGet(snapshotID)
	if err != nil {
		return nil, err
	}
	if snapshot.RoomID != roomID {
		return nil, errcode.NewError(20030003, "快照不属于该房间")
	}

	// 2. 获取原房间
	originalRoom, err := svc.dao.RoomGet(roomID)
	if err != nil {
		return nil, err
	}

	// 3. 创建新房间（继承原房间上下文）
	newRoom := &model.Room{
		ID:           uuid.New().String(),
		Name:         fmt.Sprintf("%s [分支-第%d轮]", originalRoom.Name, snapshot.Round),
		Topic:        originalRoom.Topic,
		Description:  originalRoom.Description,
		Status:       "preparing",
		CurrentPhase: snapshot.Phase,
		CurrentRound: snapshot.Round,
		MaxRounds:    originalRoom.MaxRounds,
		ForkedFrom:   &roomID,
		CreatedBy:    userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = svc.dao.RoomCreate(newRoom)
	if err != nil {
		return nil, err
	}

	// 4. 解析快照中的 Agent 状态并恢复
	var agentStates []map[string]interface{}
	_ = json.Unmarshal([]byte(snapshot.AgentStates), &agentStates)

	// 先复制原房间的 agents
	originalAgents, err := svc.dao.AgentListByRoom(roomID)
	if err != nil {
		return nil, err
	}

	for _, origAgent := range originalAgents {
		newAgent := &model.RoomAgent{
			ID:           uuid.New().String(),
			RoomID:       newRoom.ID,
			TemplateID:   origAgent.TemplateID,
			Name:         origAgent.Name,
			Role:         origAgent.Role,
			Avatar:       origAgent.Avatar,
			Personality:  origAgent.Personality,
			Expertise:    origAgent.Expertise,
			Model:        origAgent.Model,
			SystemPrompt: origAgent.SystemPrompt,
			IsActive:     origAgent.IsActive,
			CreatedAt:    time.Now(),
		}
		// 从快照恢复状态
		for _, st := range agentStates {
			if aid, ok := st["agent_id"].(string); ok && aid == origAgent.ID {
				if stance, ok := st["stance"].(string); ok {
					newAgent.Stance = stance
				}
				if conf, ok := st["confidence"].(float64); ok {
					newAgent.Confidence = int(conf)
				}
				if energy, ok := st["energy"].(float64); ok {
					newAgent.Energy = int(energy)
				}
			}
		}
		if newAgent.Stance == "" {
			newAgent.Stance = "neutral"
		}
		if newAgent.Confidence == 0 {
			newAgent.Confidence = 50
		}
		if newAgent.Energy == 0 {
			newAgent.Energy = 100
		}
		_ = svc.dao.AgentCreate(newAgent)
	}

	// 5. 复制该轮次之前的所有消息
	allMessages, _, err := svc.dao.MessageListByRoom(roomID, 0, "", 1, 10000)
	if err != nil {
		log.Printf("[Rollback] 复制消息失败: %v", err)
	} else {
		for _, msg := range allMessages {
			if msg.Round > snapshot.Round {
				continue
			}
			newMsg := &model.Message{
				ID:           uuid.New().String(),
				RoomID:       newRoom.ID,
				SenderID:     msg.SenderID,
				SenderType:   msg.SenderType,
				SenderName:   msg.SenderName,
				SenderAvatar: msg.SenderAvatar,
				Content:      msg.Content,
				MsgType:      msg.MsgType,
				Phase:        msg.Phase,
				Round:        msg.Round,
				Metadata:     msg.Metadata,
				CreatedAt:    msg.CreatedAt,
			}
			_ = svc.dao.MessageCreate(newMsg)
		}
	}

	// 6. 保存分叉通知消息
	forkMsg := &model.Message{
		ID:         uuid.New().String(),
		RoomID:     newRoom.ID,
		SenderID:   "system",
		SenderType: "system",
		SenderName: "System",
		Content:    fmt.Sprintf("从原房间第 %d 轮回滚创建。快照原因: %s", snapshot.Round, snapshot.TriggerReason),
		MsgType:    "fork_notice",
		Phase:      snapshot.Phase,
		Round:      snapshot.Round,
		Metadata:   fmt.Sprintf(`{"snapshot_id":"%s","original_room":"%s"}`, snapshotID, roomID),
	}
	_ = svc.dao.MessageCreate(forkMsg)

	// 7. 重新获取房间以获取完整信息（包括 GORM 回填的时间戳）
	freshRoom, err := svc.dao.RoomGet(newRoom.ID)
	if err != nil {
		log.Printf("[Rollback] 重新获取房间失败: %v", err)
		resp := buildRoomResponse(newRoom)
		return &resp, nil
	}
	resp := buildRoomResponse(freshRoom)
	return &resp, nil
}

// CreateManualSnapshot 手动创建快照
func (svc *Service) CreateManualSnapshot(roomID, reason string) (*model.RoomSnapshot, *errcode.Error) {
	room, err := svc.dao.RoomGet(roomID)
	if err != nil {
		return nil, err
	}
	agents, err := svc.dao.AgentListByRoom(roomID)
	if err != nil {
		return nil, err
	}

	agentStates := make([]map[string]interface{}, 0, len(agents))
	for _, a := range agents {
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
	for _, a := range agents {
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
		RoomID:         roomID,
		Round:          room.CurrentRound,
		Phase:          room.CurrentPhase,
		ConsensusScore: room.ConsensusScore,
		AgentStates:    string(agentStatesJSON),
		KeyDecisions:   string(keyDecisionsJSON),
		TriggerReason:  reason,
		CreatedAt:      time.Now(),
	}
	err = svc.dao.SnapshotCreate(snapshot)
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}
