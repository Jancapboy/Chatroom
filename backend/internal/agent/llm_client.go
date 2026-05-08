package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/Jancapboy/Chatroom/backend/global"
	"github.com/Jancapboy/Chatroom/backend/internal/service"
)

// LLMClient LLM统一调用封装
type LLMClient struct {
	aiService *service.AIService
}

// NewLLMClient 创建LLM客户端
func NewLLMClient() *LLMClient {
	config := service.AIConfig{
		APIKey:  global.AISettings.APIKey,
		Model:   global.AISettings.Model,
		BaseURL: global.AISettings.BaseURL,
	}
	return &LLMClient{
		aiService: service.NewAIService(config),
	}
}

// AgentResponse Agent回复结构
type AgentResponse struct {
	Content    string `json:"content"`
	Stance     string `json:"stance"`
	Confidence int    `json:"confidence"`
	Model      string `json:"model"`
}

// Complete 调用LLM生成Agent回复
func (c *LLMClient) Complete(ctx context.Context, systemPrompt, userPrompt string) (*AgentResponse, error) {
	req := service.AIRequest{
		Model: global.AISettings.Model,
		Messages: []service.AIMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Stream: false,
	}

	// 设置超时
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := c.aiService.CallDeepSeek(callCtx, req)
	if err != nil {
		return nil, fmt.Errorf("LLM调用失败: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("LLM返回错误: %s", resp.Error)
	}

	// 解析回复
	stance, confidence, reply := ParseResponse(resp.Content)

	return &AgentResponse{
		Content:    reply,
		Stance:     stance,
		Confidence: confidence,
		Model:      resp.Model,
	}, nil
}

// CompleteWithContext 带上下文的LLM调用
func (c *LLMClient) CompleteWithContext(ctx context.Context, persona *Persona, topic, phase string, round, maxRounds int, contextMessages string) (*AgentResponse, error) {
	prompt := BuildPrompt(persona, topic, phase, round, maxRounds, contextMessages)
	return c.Complete(ctx, persona.SystemPrompt, prompt)
}
