package api

import (
	"strconv"

	"github.com/Jancapboy/Chatroom/backend/internal/service"
	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
	"github.com/Jancapboy/Chatroom/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Room struct{}

func NewRoom() Room {
	return Room{}
}

// List 房间列表
func (r Room) List(c *gin.Context) {
	resp := response.NewResponse(c)
	status := c.Query("status")
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if pageSize <= 0 {
		pageSize = 20
	}

	svc := service.New(c.Request.Context())
	rooms, total, err := svc.RoomList(status, page, pageSize)
	if err != nil {
		resp.ToErrorResponse(err)
		return
	}

	resp.ToResponse(gin.H{
		"list":  rooms,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// Create 创建房间
func (r Room) Create(c *gin.Context) {
	resp := response.NewResponse(c)
	param := service.RoomCreateRequest{}
	if err := c.ShouldBindJSON(&param); err != nil {
		resp.ToErrorResponse(errcode.InvalidParams)
		return
	}

	userID, _ := c.Get("UserID")
	uid, _ := userID.(uint64)

	svc := service.New(c.Request.Context())
	room, err := svc.RoomCreate(&param, uid)
	if err != nil {
		resp.ToErrorResponse(err)
		return
	}

	resp.ToResponse(room)
}

// Get 房间详情
func (r Room) Get(c *gin.Context) {
	resp := response.NewResponse(c)
	id := c.Param("id")
	if id == "" {
		resp.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	room, err := svc.RoomGet(id)
	if err != nil {
		resp.ToErrorResponse(err)
		return
	}

	resp.ToResponse(room)
}

// Start 开始推演
func (r Room) Start(c *gin.Context) {
	resp := response.NewResponse(c)
	id := c.Param("id")
	if id == "" {
		resp.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.RoomStart(id)
	if err != nil {
		resp.ToErrorResponse(err)
		return
	}

	resp.ToResponse(gin.H{"status": "running"})
}

// Pause 暂停推演
func (r Room) Pause(c *gin.Context) {
	resp := response.NewResponse(c)
	id := c.Param("id")
	if id == "" {
		resp.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.RoomPause(id)
	if err != nil {
		resp.ToErrorResponse(err)
		return
	}

	resp.ToResponse(gin.H{"status": "paused"})
}

// Delete 删除房间
func (r Room) Delete(c *gin.Context) {
	resp := response.NewResponse(c)
	id := c.Param("id")
	if id == "" {
		resp.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.RoomDelete(id)
	if err != nil {
		resp.ToErrorResponse(err)
		return
	}

	resp.ToResponse(gin.H{"deleted": true})
}

// Messages 获取房间消息
func (r Room) Messages(c *gin.Context) {
	resp := response.NewResponse(c)
	id := c.Param("id")
	if id == "" {
		resp.ToErrorResponse(errcode.InvalidParams)
		return
	}

	round, _ := strconv.Atoi(c.Query("round"))
	phase := c.Query("phase")
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if pageSize <= 0 {
		pageSize = 50
	}

	svc := service.New(c.Request.Context())
	msgs, total, err := svc.Dao().MessageListByRoom(id, round, phase, page, pageSize)
	if err != nil {
		resp.ToErrorResponse(err)
		return
	}

	resp.ToResponse(gin.H{
		"list":  msgs,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// AddAgent 向房间添加Agent
func (r Room) AddAgent(c *gin.Context) {
	resp := response.NewResponse(c)
	id := c.Param("id")
	if id == "" {
		resp.ToErrorResponse(errcode.InvalidParams)
		return
	}

	var req struct {
		TemplateID string `json:"template_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	agent, err := svc.AddAgentToRoom(id, req.TemplateID)
	if err != nil {
		resp.ToErrorResponse(err)
		return
	}

	resp.ToResponse(agent)
}
