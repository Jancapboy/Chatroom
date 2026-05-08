// internal/service/ai.go
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// 公共的消息结构体
type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// API请求结构体
type AIRequest struct {
	Model    string      `json:"model"`
	Messages []AIMessage `json:"messages"`
	Stream   bool        `json:"stream"`
}

// API响应结构体
type AIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type UnifiedAIResponse struct {
	Success bool        `json:"success"`
	Content string      `json:"content"`
	Model   string      `json:"model,omitempty"`
	Usage   interface{} `json:"usage,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// 修改 AIService 结构体
type AIService struct {
	config AIConfig // 使用配置结构体
	client *http.Client
}

// AIConfig 结构体（与 setting.AISettingS 对应）
type AIConfig struct {
	APIKey  string
	Model   string
	BaseURL string
}

// 修改构造函数，接收配置
func NewAIService(config AIConfig) *AIService {
	// 设置默认值（双重保险）
	if config.Model == "" {
		config.Model = "deepseek-chat"
	}
	if config.BaseURL == "" {
		config.BaseURL = "https://api.deepseek.com"
	}

	return &AIService{
		config: config,
		client: &http.Client{},
	}
}

func (s *AIService) CallDeepSeek(ctx context.Context, req AIRequest) (*UnifiedAIResponse, error) {
	if s.config.APIKey == "" {
		return &UnifiedAIResponse{
			Success: false,
			Error:   "DeepSeek API key not configured",
		}, fmt.Errorf("DeepSeek API key not configured")
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return &UnifiedAIResponse{
			Success: false,
			Error:   "Failed to marshal request",
		}, err
	}

	// 使用配置中的BaseURL
	apiURL := s.config.BaseURL + "/v1/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return &UnifiedAIResponse{
			Success: false,
			Error:   "Failed to create request",
		}, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.config.APIKey)
	httpReq.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return &UnifiedAIResponse{
			Success: false,
			Error:   "API request failed",
		}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &UnifiedAIResponse{
			Success: false,
			Error:   "Failed to read response",
		}, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &UnifiedAIResponse{
			Success: false,
			Error:   fmt.Sprintf("API request failed with status %d", resp.StatusCode),
		}, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 解析DeepSeek的原始响应
	var deepseekResp AIResponse
	if err := json.Unmarshal(body, &deepseekResp); err != nil {
		return &UnifiedAIResponse{
			Success: false,
			Error:   "Failed to parse API response",
		}, fmt.Errorf("failed to parse response: %w", err)
	}

	// 转换为统一格式
	if len(deepseekResp.Choices) > 0 && deepseekResp.Choices[0].Message.Content != "" {
		return &UnifiedAIResponse{
			Success: true,
			Content: deepseekResp.Choices[0].Message.Content,
			Model:   deepseekResp.Model,
			Usage:   deepseekResp.Usage,
		}, nil
	}

	return &UnifiedAIResponse{
		Success: false,
		Error:   "No valid response from AI",
	}, fmt.Errorf("no valid response from AI")
}
