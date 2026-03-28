package api

import (
	"net/http"

	"github.com/Jancapboy/Chatroom/internal/service"
	"github.com/gin-gonic/gin"
)

type ThreeD struct {
	config *service.Hunyuan3DConfig
}

func NewThreeD(secretID, secretKey string) ThreeD {
	return ThreeD{
		config: &service.Hunyuan3DConfig{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	}
}

// Generate 提交3D生成任务
func (t ThreeD) Generate(c *gin.Context) {
	var req struct {
		Prompt       string `json:"prompt" binding:"required"`
		ResultFormat string `json:"result_format"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误: " + err.Error()})
		return
	}

	format := req.ResultFormat
	if format == "" {
		format = "GLB"
	}

	jobId, err := service.SubmitHunyuan3DJob(t.config, &service.Hunyuan3DSubmitRequest{
		Prompt:       req.Prompt,
		ResultFormat: format,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "提交失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "任务已提交", "job_id": jobId})
}

// Query 查询3D生成结果
func (t ThreeD) Query(c *gin.Context) {
	var req struct {
		JobId string `json:"job_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误: " + err.Error()})
		return
	}

	result, err := service.QueryHunyuan3DJob(t.config, req.JobId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "查询失败: " + err.Error()})
		return
	}

	resp := gin.H{"code": 0, "status": result.Response.Status}
	if result.Response.Status == "DONE" && len(result.Response.ResultFile3Ds) > 0 {
		files := make([]gin.H, 0)
		for _, f := range result.Response.ResultFile3Ds {
			files = append(files, gin.H{"type": f.Type, "url": f.Url})
		}
		resp["files"] = files
	}

	c.JSON(http.StatusOK, resp)
}
