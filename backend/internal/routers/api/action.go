package api

import (
	"github.com/Jancapboy/Chatroom/backend/internal/service"
	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
	"github.com/Jancapboy/Chatroom/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// Action Action API
type Action struct {
	service *service.ActionService
}

// NewAction 创建Action API
func NewAction() *Action {
	return &Action{service: service.NewActionService()}
}

// List 获取房间的Action列表
func (a *Action) List(c *gin.Context) {
	resp := response.NewResponse(c)
	roomID := c.Param("id")
	actions, err := a.service.ListActions(roomID)
	if err != nil {
		resp.ToErrorResponse(errcode.Convert(err))
		return
	}
	resp.ToResponse(actions)
}

// GenerateDrafts 生成Action草稿（推演完成后调用）
func (a *Action) GenerateDrafts(c *gin.Context) {
	resp := response.NewResponse(c)
	roomID := c.Param("id")
	
	var req struct {
		Conclusion string `json:"conclusion"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ToErrorResponse(errcode.InvalidParams)
		return
	}
	
	drafts, err := a.service.GenerateActionDrafts(roomID, req.Conclusion)
	if err != nil {
		resp.ToErrorResponse(errcode.Convert(err))
		return
	}
	
	resp.ToResponse(drafts)
}

// Approve 批准Action
func (a *Action) Approve(c *gin.Context) {
	resp := response.NewResponse(c)
	actionID := c.Param("action_id")
	userID := c.GetUint64("user_id")
	
	act, err := a.service.ApproveAction(actionID, userID)
	if err != nil {
		resp.ToErrorResponse(errcode.Convert(err))
		return
	}
	
	resp.ToResponse(act)
}

// Execute 执行Action
func (a *Action) Execute(c *gin.Context) {
	resp := response.NewResponse(c)
	actionID := c.Param("action_id")
	userID := c.GetUint64("user_id")
	
	act, err := a.service.ExecuteAction(actionID, userID)
	if err != nil {
		resp.ToErrorResponse(errcode.Convert(err))
		return
	}
	
	resp.ToResponse(act)
}
