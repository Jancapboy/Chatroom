package service

import (
	"fmt"
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
		CreatedAt:      r.CreatedAt,
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
