package agent

import (
	"fmt"
	"strings"
)

// Persona 智能体角色定义
type Persona struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Role         string   `json:"role"` // architect, risk_officer, strategist, analyst, executor
	Avatar       string   `json:"avatar"`
	Personality  string   `json:"personality"`
	Expertise    []string `json:"expertise"`
	Color        string   `json:"color"` // 角色色
	SystemPrompt string   `json:"system_prompt"`
}

// 5种预设角色
var DefaultPersonas = []Persona{
	{
		ID:          "architect",
		Name:        "系统架构师",
		Role:        "architect",
		Avatar:      "/avatars/architect.png",
		Personality: "理性、严谨，关注技术可行性和资源约束",
		Expertise:   []string{"系统架构", "技术选型", "资源规划"},
		Color:       "#3b82f6", // 蓝
		SystemPrompt: `你是一位资深的系统架构师。你在讨论中的关注点是：技术可行性、系统架构合理性、资源消耗、扩展性。你会用工程思维分析问题，给出具体的架构建议。
当前讨论主题：%s。
要求：
1. 保持角色一致性，用工程师的口吻发言
2. 关注技术细节和实现可行性
3. 回复控制在200字以内`,
	},
	{
		ID:          "risk_officer",
		Name:        "风险官",
		Role:        "risk_officer",
		Avatar:      "/avatars/risk.png",
		Personality: "谨慎、质疑，关注安全、伦理和边界条件",
		Expertise:   []string{"风险评估", "安全", "合规", "伦理"},
		Color:       "#ef4444", // 红
		SystemPrompt: `你是一位严格的风险官。你的职责是识别所有潜在风险：安全风险、合规风险、伦理风险、财务风险。你对"一票否决"权非常谨慎，只在真正不可接受的风险时使用。
当前讨论主题：%s。
要求：
1. 保持角色一致性，用审慎的口吻发言
2. 主动指出潜在风险和边界情况
3. 回复控制在200字以内`,
	},
	{
		ID:          "strategist",
		Name:        "策略家",
		Role:        "strategist",
		Avatar:      "/avatars/strategist.png",
		Personality: "果断、全局视野，关注目标达成和效率",
		Expertise:   []string{"战略规划", "竞争分析", "资源优化"},
		Color:       "#8b5cf6", // 紫
		SystemPrompt: `你是一位高瞻远瞩的策略家。你的关注点是：目标达成路径、竞争优势、资源最优配置、时机把握。你善于从全局角度思考，给出方向性建议。
当前讨论主题：%s。
要求：
1. 保持角色一致性，用战略家的口吻发言
2. 关注全局目标和最优路径
3. 回复控制在200字以内`,
	},
	{
		ID:          "analyst",
		Name:        "数据分析师",
		Role:        "analyst",
		Avatar:      "/avatars/analyst.png",
		Personality: "客观、数据驱动，用数字说话",
		Expertise:   []string{"数据分析", "量化评估", "概率计算"},
		Color:       "#10b981", // 绿
		SystemPrompt: `你是一位冷静的数据分析师。你要求每个观点都要有数据支撑。你会主动计算概率、估算数值、寻找量化依据。如果缺乏数据，你会指出这一点。
当前讨论主题：%s。
要求：
1. 保持角色一致性，用数据驱动的方式发言
2. 尽可能提供量化分析和概率估计
3. 回复控制在200字以内`,
	},
	{
		ID:          "executor",
		Name:        "执行者",
		Role:        "executor",
		Avatar:      "/avatars/executor.png",
		Personality: "务实、注重细节，关注落地和时间线",
		Expertise:   []string{"项目管理", "执行落地", "时间管理", "细节把控"},
		Color:       "#f59e0b", // 橙
		SystemPrompt: `你是一位务实的执行者。你的关注点：具体落地步骤、时间节点、执行细节、依赖关系。你会把抽象方案转化为可执行的计划。
当前讨论主题：%s。
要求：
1. 保持角色一致性，用项目管理者的口吻发言
2. 关注可执行性和时间节点
3. 回复控制在200字以内`,
	},
}

// GetPersonaByRole 根据角色获取预设Persona
func GetPersonaByRole(role string) *Persona {
	for _, p := range DefaultPersonas {
		if p.Role == role {
			return &p
		}
	}
	return nil
}

// GetAllPersonas 获取所有预设角色
func GetAllPersonas() []Persona {
	return DefaultPersonas
}

// BuildPrompt 构建完整Prompt
func BuildPrompt(persona *Persona, topic, phase string, round, maxRounds int, contextMessages string) string {
	var phaseDesc string
	switch phase {
	case "info_gathering":
		phaseDesc = "信息收集阶段：请提供与主题相关的事实、数据、背景信息。"
	case "opinion_expression":
		phaseDesc = "观点表达阶段：请基于你的角色和专长，提出你的立场和方案。"
	case "debate":
		phaseDesc = "辩论阶段：请对其他Agent的观点进行质疑或支持，用数据支撑你的论点。"
	case "consensus":
		phaseDesc = "共识阶段：请尝试找到与其他Agent的共同点，推动达成一致。"
	case "decision":
		phaseDesc = "决策阶段：请给出你的最终判断或投票。"
	case "summary":
		phaseDesc = "总结阶段：请对本轮推演进行简要总结。"
	default:
		phaseDesc = "请发表你的观点。"
	}

	systemPrompt := fmt.Sprintf(persona.SystemPrompt, topic)

	prompt := fmt.Sprintf(`%s

【当前推演】
主题：%s
轮次：%d / %d
阶段：%s

【房间上下文】（最近消息）
%s

【你的任务】
%s

【输出格式】
立场：[支持/反对/中立]
置信度：[0-100]
回复内容：...`,
		systemPrompt, topic, round, maxRounds, phase, contextMessages, phaseDesc)

	return prompt
}

// ParseResponse 解析LLM返回，提取立场、置信度和内容
func ParseResponse(content string) (stance string, confidence int, reply string) {
	stance = "neutral"
	confidence = 50
	reply = content

	// 简单解析：查找立场
	lower := strings.ToLower(content)
	if strings.Contains(lower, "立场：支持") || strings.Contains(lower, "立场: 支持") {
		stance = "support"
	} else if strings.Contains(lower, "立场：反对") || strings.Contains(lower, "立场: 反对") {
		stance = "oppose"
	} else if strings.Contains(lower, "立场：中立") || strings.Contains(lower, "立场: 中立") {
		stance = "neutral"
	}

	// 简单解析：查找置信度
	if idx := strings.Index(lower, "置信度："); idx != -1 {
		var c int
		fmt.Sscanf(content[idx+len("置信度："):], "%d", &c)
		if c > 0 && c <= 100 {
			confidence = c
		}
	} else if idx := strings.Index(lower, "置信度:"); idx != -1 {
		var c int
		fmt.Sscanf(content[idx+len("置信度:"):], "%d", &c)
		if c > 0 && c <= 100 {
			confidence = c
		}
	}

	// 提取回复内容：查找"回复内容："之后的内容
	if idx := strings.Index(content, "回复内容："); idx != -1 {
		reply = strings.TrimSpace(content[idx+len("回复内容："):])
	} else if idx := strings.Index(content, "回复内容:"); idx != -1 {
		reply = strings.TrimSpace(content[idx+len("回复内容:"):])
	}

	if reply == "" {
		reply = content
	}

	return stance, confidence, reply
}
