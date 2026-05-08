package api

import (
	"github.com/Jancapboy/Chatroom/backend/internal/request"
	"github.com/Jancapboy/Chatroom/backend/internal/service"
	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
	"github.com/Jancapboy/Chatroom/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type User struct{}

func NewUser() User {
	return User{}
}

func (u User) Register(c *gin.Context) {
	param := request.UserRegisterRequest{}
	r := response.NewResponse(c)
	if c.ShouldBindJSON(&param) != nil {
		r.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.UserRegister(&param)
	if err != nil {
		r.ToErrorResponse(err)
		return
	}
	r.ToResponse(nil)
}

func (u User) Login(c *gin.Context) {
	param := request.UserLoginRequest{}
	r := response.NewResponse(c)
	if c.ShouldBindJSON(&param) != nil {
		r.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	content, err := svc.UserLogin(&param)
	if err != nil {
		r.ToErrorResponse(err)
		return
	}
	r.ToResponse(content)
}
