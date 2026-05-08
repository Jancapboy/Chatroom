package response

import (
	"net/http"

	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Ctx *gin.Context
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{
		Ctx: ctx,
	}
}

func (r *Response) ToResponse(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

func (r *Response) ToResponseList(list interface{}, total int64, page, pageSize int) {
	r.Ctx.JSON(http.StatusOK, gin.H{
		"code":     0,
		"msg":      "success",
		"data": gin.H{
			"list":     list,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

func (r *Response) ToErrorResponse(err *errcode.Error) {
	r.Ctx.JSON(err.StatusCode(), err)
}
