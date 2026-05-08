package simulation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Jancapboy/Chatroom/backend/internal/agent"
	"github.com/Jancapboy/Chatroom/backend/internal/model"
)

// LearningEngine Agent学习引擎
type LearningEngine struct {
	llmClient *agent.LLMClient
}

// NewLearningEngine 创建学习引擎
func NewLearningEngine(llmClient *agent.LLMClient) *LearningEngine {
	return &LearningEngine{llmClient: llmClient}
}

// ExtractExperience 从一轮辩论中提取经验
func (le *LearningEngine) ExtractExperience(
	ctx context.Context,
	agentID string,
	agentName string,
	agentRole string,
	topic string,
	messages []model.Message,
	consensus float64,
) (*model.AgentMemory, error) {
	// 构建提示词
	prompt := le.buildExtractionPrompt(agentName, agentRole, topic, messages, consensus)

	// 调用LLM提取经验
	resp, err := le.llmClient.Complete(ctx, "你是一位经验提取专家。分析以下Agent在辩论中的表现，提取关键经验教训。", prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM提取经验失败: %w", err)
	}

	// 解析结果
	experience := strings.TrimSpace(resp.Content)
	if experience == "" || experience == "无" || len(experience) < 10 {
		return nil, nil // 没有值得记录的经验
	}

	// 判断结果好坏
	outcome := "neutral"
	if consensus >= 70 {
		outcome = "success"
	} else if consensus < 30 {
		outcome = "failure"
	}

	// 创建记忆
	memory := &model.AgentMemory{
		ID:         fmt.Sprintf("mem-%d", time.Now().UnixNano()),
		AgentID:    agentID,
		Topic:      extractTopicTag(topic),
		Experience: experience,
		Outcome:    outcome,
		Confidence: int(consensus),
		CreatedAt:  time.Now(),
	}

	return memory, nil
}

// GenerateLearningPrompt 生成带有学习记忆的System Prompt
func (le *LearningEngine) GenerateLearningPrompt(
	ctx context.Context,
	agentID string,
	basePrompt string,
	currentTopic string,
	getMemories func(agentID, topic string, limit int) ([]model.AgentMemory, error),
) (string, error) {
	// 获取相关记忆
	memories, err := getMemories(agentID, extractTopicTag(currentTopic), 3)
	if err != nil {
		return basePrompt, nil // 没记忆也不报错
	}

	if len(memories) == 0 {
		return basePrompt, nil
	}

	// 构建记忆上下文
	var sb strings.Builder
	sb.WriteString(basePrompt)
	sb.WriteString("\n\n【你的历史经验】\n")
	for i, m := range memories {
		sb.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, m.Outcome, m.Experience))
	}
	sb.WriteString("\n请在本次讨论中参考以上经验，避免重复过去的错误，发扬有效的策略。")

	return sb.String(), nil
}

// buildExtractionPrompt 构建经验提取提示词
func (le *LearningEngine) buildExtractionPrompt(
	agentName, agentRole, topic string,
	messages []model.Message,
	consensus float64,
) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("话题: %s\n", topic))
	sb.WriteString(fmt.Sprintf("Agent: %s（角色: %s）\n", agentName, agentRole))
	sb.WriteString(fmt.Sprintf("本轮共识度: %.1f%%\n\n", consensus))
	sb.WriteString("该Agent在本轮的发言:\n")

	for _, msg := range messages {
		if msg.SenderName == agentName {
			sb.WriteString(fmt.Sprintf("- [%s] %s\n", msg.Phase, msg.Content))
		}
	}

	sb.WriteString("\n其他Agent对该Agent观点的反馈:\n")
	for _, msg := range messages {
		if msg.SenderName != agentName && msg.SenderType == "agent" {
			// 简单判断是否是回应
			if containsStr(msg.Content, agentName) || containsStr(msg.Content, "同意") || containsStr(msg.Content, "反对") {
				sb.WriteString(fmt.Sprintf("- %s: %s\n", msg.SenderName, msg.Content))
			}
		}
	}

	sb.WriteString("\n请提取该Agent在本轮讨论中的关键经验教训（1-3条，每条30字以内）：\n")
	sb.WriteString("- 格式: \"在[场景]中，[策略]导致了[结果]，下次应该[调整]\"\n")
	sb.WriteString("- 如果没有值得记录的经验，回复\"无\"\n")
	sb.WriteString("- 只输出经验内容，不要解释")

	return sb.String()
}

// extractTopicTag 提取话题标签（简化）
func extractTopicTag(topic string) string {
	if len(topic) <= 20 {
		return topic
	}
	return topic[:20]
}

// containsStr 字符串包含判断
func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || strings.Contains(s, substr))
}
