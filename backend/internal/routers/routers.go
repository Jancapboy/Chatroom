package routers

import (
	"github.com/Jancapboy/Chatroom/backend/global"
	"github.com/Jancapboy/Chatroom/backend/internal/middleware"
	"github.com/Jancapboy/Chatroom/backend/internal/routers/api"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	if global.ServerSettings.RunMode == "debug" {
		r.Use(gin.Logger(), gin.Recovery())
	}
	r.Use(middleware.Cors())

	user := api.NewUser()
	ai := api.NewAI()
	room := api.NewRoom()
	agent := api.NewAgent()

	// API v1 路由组
	apiV1 := r.Group("/api/v1")
	{
		// 用户（保留）
		apiV1.POST("/register", user.Register)
		apiV1.POST("/login", user.Login)

		// AI聊天（保留）
		apiV1.POST("/ai/chat", ai.Chat)

		// 房间管理（新增）
		apiV1.GET("/rooms", room.List)
		apiV1.POST("/rooms", room.Create)
		apiV1.GET("/rooms/:id", room.Get)
		apiV1.POST("/rooms/:id/start", room.Start)
		apiV1.POST("/rooms/:id/pause", room.Pause)
		apiV1.DELETE("/rooms/:id", room.Delete)
		apiV1.GET("/rooms/:id/messages", room.Messages)
		apiV1.POST("/rooms/:id/agents", room.AddAgent)

		// 回溯/分支（新增）
		apiV1.GET("/rooms/:id/snapshots", room.Snapshots)
		apiV1.POST("/rooms/:id/rollback", room.Rollback)
		apiV1.POST("/rooms/:id/fork", room.Fork)
		apiV1.POST("/rooms/:id/snapshots", room.CreateSnapshot)

		// 智能体模板（新增）
		apiV1.GET("/agents/templates", agent.Templates)
	}

	// WebSocket 路由（房间隔离）
	wsGroup := r.Group("/ws")
	if global.ServerSettings.RunMode != "debug" {
		wsGroup.Use(middleware.JWT())
	}
	{
		// 保留原有广播式WS（兼容）
		wsGroup.GET("/", WebsocketHandler)
		// 新增房间隔离WS
		wsGroup.GET("/rooms/:id", WebsocketRoomHandler)
	}

	return r
}
