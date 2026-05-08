package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Jancapboy/Chatroom/backend/internal/action"
	"github.com/Jancapboy/Chatroom/backend/internal/dao"
	"github.com/Jancapboy/Chatroom/backend/internal/model"
)

// ActionService Action服务
type ActionService struct {
	dao *dao.Dao
}

func NewActionService() *ActionService {
	return &ActionService{dao: dao.New(nil)}
}

// GenerateActionDrafts 推演完成后生成Action草稿
func (s *ActionService) GenerateActionDrafts(roomID string, conclusion string) ([]*model.RoomAction, error) {
	// 获取房间信息
	room, err := s.dao.RoomGet(roomID)
	if err != nil {
		return nil, fmt.Errorf("获取房间失败: %w", err)
	}
	
	// 获取消息列表
	messages, _, _ := s.dao.MessageListByRoom(roomID, 0, "", 1, 100)
	
	var drafts []*model.RoomAction
	
	// 1. 生成Markdown报告Action
	reportParams := map[string]interface{}{
		"room_name":  room.Name,
		"topic":      room.Topic,
		"conclusion": conclusion,
		"messages":   extractMessageSummary(messages),
	}
	paramsJSON, _ := json.Marshal(reportParams)
	
	reportAction := &model.RoomAction{
		ID:          fmt.Sprintf("act-%d", time.Now().UnixNano()),
		RoomID:      roomID,
		Type:        "markdown_report",
		Status:      "pending",
		Title:       "生成推演报告",
		Description: fmt.Sprintf("将推演过程导出为Markdown报告"),
		Params:      string(paramsJSON),
		CreatedAt:   time.Now(),
	}
	if err := s.dao.ActionCreate(reportAction); err != nil {
		return nil, err
	}
	drafts = append(drafts, reportAction)
	
	// 2. 如果有结论，生成飞书通知Action
	if conclusion != "" {
		feishuParams := map[string]interface{}{
			"content": fmt.Sprintf("推演《%s》已完成，结论：%s", room.Name, truncateString(conclusion, 200)),
		}
		paramsJSON, _ := json.Marshal(feishuParams)
		
		feishuAction := &model.RoomAction{
			ID:          fmt.Sprintf("act-%d", time.Now().UnixNano()+1),
			RoomID:      roomID,
			Type:        "feishu_message",
			Status:      "pending",
			Title:       "发送飞书通知",
			Description: "将推演结论发送到飞书",
			Params:      string(paramsJSON),
			CreatedAt:   time.Now(),
		}
		if err := s.dao.ActionCreate(feishuAction); err != nil {
			return nil, err
		}
		drafts = append(drafts, feishuAction)
	}
	
	return drafts, nil
}

// ExecuteAction 执行Action（人类确认后）
func (s *ActionService) ExecuteAction(actionID string, userID uint64) (*model.RoomAction, error) {
	// 获取Action
	act, err := s.dao.ActionGet(actionID)
	if err != nil {
		return nil, fmt.Errorf("获取Action失败: %w", err)
	}
	
	if act.Status != "pending" && act.Status != "approved" {
		return nil, fmt.Errorf("Action状态不正确: %s", act.Status)
	}
	
	// 获取对应的Action执行器
	executor, ok := action.GetAction(act.Type)
	if !ok {
		return nil, fmt.Errorf("未知Action类型: %s", act.Type)
	}
	
	// 解析参数
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(act.Params), &params); err != nil {
		return nil, fmt.Errorf("解析参数失败: %w", err)
	}
	
	// 校验参数
	if err := executor.Validate(params); err != nil {
		act.Status = "failed"
		act.Result = fmt.Sprintf("参数校验失败: %v", err)
		s.dao.ActionUpdate(act)
		return act, err
	}
	
	// 执行
	result, execErr := executor.Execute(params)
	now := time.Now()
	act.ExecutedAt = &now
	act.ApprovedBy = userID
	
	if execErr != nil {
		act.Status = "failed"
		act.Result = fmt.Sprintf("执行失败: %v", execErr)
	} else {
		act.Status = "executed"
		act.Result = result
	}
	
	if err := s.dao.ActionUpdate(act); err != nil {
		return nil, err
	}
	
	return act, execErr
}

// ListActions 列出房间的Action
func (s *ActionService) ListActions(roomID string) ([]model.RoomAction, error) {
	actions, err := s.dao.ActionListByRoom(roomID)
	if err != nil {
		return nil, err
	}
	return actions, nil
}

// ApproveAction 批准Action（预执行确认）
func (s *ActionService) ApproveAction(actionID string, userID uint64) (*model.RoomAction, error) {
	act, err := s.dao.ActionGet(actionID)
	if err != nil {
		return nil, err
	}
	
	if act.Status != "pending" {
		return nil, fmt.Errorf("Action状态不正确: %s", act.Status)
	}
	
	act.Status = "approved"
	act.ApprovedBy = userID
	if err := s.dao.ActionUpdate(act); err != nil {
		return nil, err
	}
	
	return act, nil
}

// helper functions
func extractMessageSummary(messages []model.Message) []string {
	var summaries []string
	for _, msg := range messages {
		if msg.MsgType == "text" && msg.SenderType == "agent" {
			summaries = append(summaries, fmt.Sprintf("[%s] %s: %s", msg.Phase, msg.SenderName, msg.Content))
		}
	}
	return summaries
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
