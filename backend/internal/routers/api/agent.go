package api

import (
	"github.com/Jancapboy/Chatroom/backend/internal/service"
	"github.com/Jancapboy/Chatroom/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Agent struct{}

func NewAgent() Agent {
	return Agent{}
}

// Templates 获取预设Agent模板列表
func (a Agent) Templates(c *gin.Context) {
	resp := response.NewResponse(c)

	svc := service.New(c.Request.Context())
	templates, err := svc.Dao().TemplateList()
	if err != nil {
		resp.ToErrorResponse(err)
		return
	}

	resp.ToResponse(gin.H{
		"templates": templates,
	})
}
