package simulation

import (
	"github.com/Jancapboy/Chatroom/backend/internal/model"
)

// Skill Agent能力技能接口
type Skill interface {
	// Name 技能名称
	Name() string

	// Description 技能描述
	Description() string

	// CanHandle 判断该技能是否能处理当前话题
	CanHandle(topic string) bool

	// EnhancePrompt 增强System Prompt，注入技能相关的指令
	EnhancePrompt(basePrompt string, context map[string]interface{}) string

	// AfterDebate 辩论后提取经验，返回结构化经验数据
	AfterDebate(agentID string, messages []model.Message, consensus float64) map[string]interface{}
}

// BaseSkill 基础技能（所有Agent默认拥有）
type BaseSkill struct{}

func (s *BaseSkill) Name() string { return "debate" }
func (s *BaseSkill) Description() string { return "基础辩论能力，参与多Agent讨论" }
func (s *BaseSkill) CanHandle(topic string) bool { return true } // 所有话题都能参与
func (s *BaseSkill) EnhancePrompt(basePrompt string, context map[string]interface{}) string {
	return basePrompt + "\n你具备基础辩论能力，能够清晰地表达观点、回应他人。"
}
func (s *BaseSkill) AfterDebate(agentID string, messages []model.Message, consensus float64) map[string]interface{} {
	return nil // 基础技能不提取经验
}

// CalculationSkill 数据计算技能
type CalculationSkill struct{}

func (s *CalculationSkill) Name() string { return "calculation" }
func (s *CalculationSkill) Description() string { return "数据计算、概率评估、量化分析" }
func (s *CalculationSkill) CanHandle(topic string) bool {
	// 话题包含数字、成本、价格、概率等关键词时激活
	calcKeywords := []string{"成本", "价格", "概率", "预算", "ROI", "收益", "数量", "规模"}
	for _, kw := range calcKeywords {
		if contains(topic, kw) {
			return true
		}
	}
	return false
}
func (s *CalculationSkill) EnhancePrompt(basePrompt string, context map[string]interface{}) string {
	return basePrompt + "\n你具备数据计算和量化分析能力。在讨论中，你应该主动进行数值估算、概率计算、成本效益分析。如果缺乏数据，请明确指出并请求补充。"
}
func (s *CalculationSkill) AfterDebate(agentID string, messages []model.Message, consensus float64) map[string]interface{} {
	// 提取数值相关经验
	return map[string]interface{}{
		"skill": "calculation",
		"note":  "在本次讨论中使用了量化分析",
	}
}

// RiskAssessmentSkill 风险评估技能
type RiskAssessmentSkill struct{}

func (s *RiskAssessmentSkill) Name() string { return "risk_assessment" }
func (s *RiskAssessmentSkill) Description() string { return "风险评估、合规检查、安全分析" }
func (s *RiskAssessmentSkill) CanHandle(topic string) bool {
	riskKeywords := []string{"风险", "安全", "合规", "漏洞", "事故", "危机", "故障"}
	for _, kw := range riskKeywords {
		if contains(topic, kw) {
			return true
		}
	}
	return false
}
func (s *RiskAssessmentSkill) EnhancePrompt(basePrompt string, context map[string]interface{}) string {
	return basePrompt + "\n你具备专业的风险评估能力。你应该从多个维度识别风险：技术风险、安全风险、合规风险、财务风险、运营风险。对于不可接受的风险，你有权提出否决。"
}
func (s *RiskAssessmentSkill) AfterDebate(agentID string, messages []model.Message, consensus float64) map[string]interface{} {
	return map[string]interface{}{
		"skill": "risk_assessment",
		"note":  "在本次讨论中识别了关键风险点",
	}
}

// CostAnalysisSkill 成本分析技能
type CostAnalysisSkill struct{}

func (s *CostAnalysisSkill) Name() string { return "cost_analysis" }
func (s *CostAnalysisSkill) Description() string { return "成本分析、ROI计算、资源优化" }
func (s *CostAnalysisSkill) CanHandle(topic string) bool {
	costKeywords := []string{"成本", "预算", "费用", "ROI", "投入", "产出", "资源", "人力"}
	for _, kw := range costKeywords {
		if contains(topic, kw) {
			return true
		}
	}
	return false
}
func (s *CostAnalysisSkill) EnhancePrompt(basePrompt string, context map[string]interface{}) string {
	return basePrompt + "\n你具备成本分析和资源优化能力。你应该计算TCO（总拥有成本）、ROI（投资回报率），并建议资源最优配置方案。"
}
func (s *CostAnalysisSkill) AfterDebate(agentID string, messages []model.Message, consensus float64) map[string]interface{} {
	return map[string]interface{}{
		"skill": "cost_analysis",
		"note":  "在本次讨论中进行了成本效益分析",
	}
}

// SkillRegistry 技能注册表
var SkillRegistry = map[string]Skill{
	"debate":           &BaseSkill{},
	"calculation":      &CalculationSkill{},
	"risk_assessment":  &RiskAssessmentSkill{},
	"cost_analysis":    &CostAnalysisSkill{},
}

// GetSkillsForTopic 根据话题获取推荐的技能列表
func GetSkillsForTopic(topic string) []string {
	var skills []string
	for name, skill := range SkillRegistry {
		if skill.CanHandle(topic) {
			skills = append(skills, name)
		}
	}
	return skills
}

// contains 字符串包含判断
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSub(s, substr))
}

func containsSub(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
