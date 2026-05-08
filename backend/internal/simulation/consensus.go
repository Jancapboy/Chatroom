package simulation

import (
	"fmt"
	"math"
	"strings"

	"github.com/Jancapboy/Chatroom/backend/internal/model"
)

// ConsensusEngine 共识度计算引擎
type ConsensusEngine struct{}

// NewConsensusEngine 创建共识引擎
func NewConsensusEngine() *ConsensusEngine {
	return &ConsensusEngine{}
}

// Calculate 计算共识度（0-100）
func (ce *ConsensusEngine) Calculate(messages []model.Message, agents []model.RoomAgent) float64 {
	if len(agents) == 0 {
		return 0
	}

	// 1. 提取所有stance标记
	stances := make(map[string]string) // agentID -> stance
	confidenceMap := make(map[string]int)

	for _, msg := range messages {
		if msg.SenderType == "agent" {
			// 尝试从metadata解析
			stance, confidence := parseMetadata(msg.Metadata)
			if stance != "" {
				stances[msg.SenderID] = stance
			}
			if confidence > 0 {
				confidenceMap[msg.SenderID] = confidence
			}
		}
	}

	// 2. 统计立场
	support := 0
	oppose := 0
	neutral := 0
	unknown := 0

	for _, agent := range agents {
		if !agent.IsActive {
			continue
		}
		stance, ok := stances[agent.ID]
		if !ok {
			unknown++
			continue
		}
		switch stance {
		case "support":
			support++
		case "oppose":
			oppose++
		case "neutral":
			neutral++
		default:
			unknown++
		}
	}

	totalActive := support + oppose + neutral + unknown
	if totalActive == 0 {
		return 0
	}

	// 3. 共识度算法：
	// - 支持 = 1.0 分
	// - 中立 = 0.5 分
	// - 反对 = 0.0 分
	// - 未知 = 0.3 分（未表态视为低共识）
	score := float64(support) + float64(neutral)*0.5 + float64(unknown)*0.3
	consensus := score / float64(totalActive) * 100

	// 4. 置信度加权：如果高置信度Agent占多数，略微提升共识度
	avgConfidence := calculateAvgConfidence(confidenceMap, agents)
	consensus = consensus * (0.8 + avgConfidence/500) // 0.8 ~ 1.0 的加权系数

	return math.Min(consensus, 100)
}

// CalculateBreakdown 计算各Agent的立场分布
func (ce *ConsensusEngine) CalculateBreakdown(messages []model.Message, agents []model.RoomAgent) map[string]interface{} {
	stances := make(map[string]string)
	confidenceMap := make(map[string]int)

	for _, msg := range messages {
		if msg.SenderType == "agent" {
			stance, confidence := parseMetadata(msg.Metadata)
			if stance != "" {
				stances[msg.SenderID] = stance
			}
			if confidence > 0 {
				confidenceMap[msg.SenderID] = confidence
			}
		}
	}

	breakdown := make(map[string]interface{})
	for _, agent := range agents {
		if !agent.IsActive {
			continue
		}
		stance := stances[agent.ID]
		if stance == "" {
			stance = "unknown"
		}
		breakdown[agent.Name] = map[string]interface{}{
			"stance":     stance,
			"confidence": confidenceMap[agent.ID],
		}
	}

	return breakdown
}

// parseMetadata 解析metadata JSON字符串
func parseMetadata(metadata string) (string, int) {
	if metadata == "" {
		return "", 0
	}
	var stance string
	var confidence int
	fmt.Sscanf(metadata, `{"confidence":%d,"stance":"%s"}`, &confidence, &stance)
	// 去除可能的尾部引号
	stance = strings.TrimSuffix(stance, `"`)
	return stance, confidence
}

// calculateAvgConfidence 计算平均置信度
func calculateAvgConfidence(confidences map[string]int, agents []model.RoomAgent) float64 {
	if len(confidences) == 0 {
		return 50 // 默认中等置信度
	}

	sum := 0
	count := 0
	for _, agent := range agents {
		if !agent.IsActive {
			continue
		}
		if c, ok := confidences[agent.ID]; ok {
			sum += c
			count++
		}
	}

	if count == 0 {
		return 50
	}
	return float64(sum) / float64(count)
}
