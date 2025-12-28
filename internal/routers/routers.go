package routers

import (
	"net/http"

	"github.com/Jancapboy/Chatroom/global"
	"github.com/Jancapboy/Chatroom/internal/middleware"
	"github.com/Jancapboy/Chatroom/internal/routers/api"
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

	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/register", user.Register)
		apiGroup.POST("/login", user.Login)
		apiGroup.POST("/ai/chat", ai.Chat) // 添加AI聊天路由
	}
	wsGroup := r.Group("/ws")
	wsGroup.Use(middleware.JWT())
	{
		wsGroup.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"msg": "test"})
		})
		wsGroup.GET("/", WebsocketHandler)
	}
	return r
}
