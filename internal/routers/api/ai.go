// internal/routers/api/ai.go
package api

import (
	"fmt"
	"net/http"

	"github.com/Jancapboy/Chatroom/global"
	"github.com/Jancapboy/Chatroom/internal/service"
	"github.com/gin-gonic/gin"
)

type AI struct {
	aiService *service.AIService
}

// 不需要参数，直接使用全局配置
func NewAI() *AI {
	// 从全局配置创建 AI 服务
	aiService := service.NewAIService(service.AIConfig{
		APIKey:  global.AISettings.APIKey,
		Model:   global.AISettings.Model,
		BaseURL: global.AISettings.BaseURL,
	})

	return &AI{
		aiService: aiService,
	}
}

// 使用 service 包中的结构体
type ChatRequest struct {
	Messages []service.AIMessage `json:"messages"`
}

// internal/routers/api/ai.go - 修正后的版本
func (a *AI) Chat(c *gin.Context) {
	var req ChatRequest
	// 不再使用 response.NewResponse(c)，直接返回JSON

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "请求参数错误: " + err.Error(),
		})
		return
	}

	// 参数验证
	if len(req.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "消息不能为空",
		})
		return
	}

	// 转换为 service.AIRequest
	serviceReq := service.AIRequest{
		Model:    "deepseek-chat", // 或者使用 global.AISettings.Model
		Messages: req.Messages,
		Stream:   false,
	}

	ctx := c.Request.Context()
	aiResp, err := a.aiService.CallDeepSeek(ctx, serviceReq)
	if err != nil {
		// 记录错误日志
		fmt.Printf("AI服务调用失败: %v\n", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "AI服务暂时不可用: " + err.Error(),
		})
		return
	}

	// 直接返回 UnifiedAIResponse 结构体
	if aiResp.Success {
		c.JSON(http.StatusOK, gin.H{
			"success": aiResp.Success,
			"content": aiResp.Content,
			"model":   aiResp.Model,
			"usage":   aiResp.Usage,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   aiResp.Error,
		})
	}
}
