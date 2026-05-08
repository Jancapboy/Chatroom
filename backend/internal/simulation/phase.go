package simulation

import (
	"sync"
)

// Phase 推演阶段定义
type Phase string

const (
	PhaseInfoGathering   Phase = "info_gathering"
	PhaseOpinion         Phase = "opinion_expression"
	PhaseDebate          Phase = "debate"
	PhaseConsensus       Phase = "consensus"
	PhaseDecision        Phase = "decision"
	PhaseSummary         Phase = "summary"
)

// PhaseController 阶段控制器
type PhaseController struct {
	currentPhase Phase
	phases       []Phase
	mu           sync.RWMutex
}

// NewPhaseController 创建阶段控制器
func NewPhaseController() *PhaseController {
	return &PhaseController{
		currentPhase: PhaseInfoGathering,
		phases: []Phase{
			PhaseInfoGathering,
			PhaseOpinion,
			PhaseDebate,
			PhaseConsensus,
			PhaseDecision,
			PhaseSummary,
		},
	}
}

// Enter 进入指定阶段
func (pc *PhaseController) Enter(phase Phase) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.currentPhase = phase
}

// Current 获取当前阶段
func (pc *PhaseController) Current() Phase {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return pc.currentPhase
}

// AllPhases 获取所有阶段列表
func (pc *PhaseController) AllPhases() []Phase {
	return pc.phases
}

// Next 获取下一个阶段
func (pc *PhaseController) Next() Phase {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	for i, p := range pc.phases {
		if p == pc.currentPhase && i+1 < len(pc.phases) {
			return pc.phases[i+1]
		}
	}
	return PhaseInfoGathering // 循环回到第一个
}

// PhaseName 获取阶段中文名
func PhaseName(phase Phase) string {
	switch phase {
	case PhaseInfoGathering:
		return "信息收集"
	case PhaseOpinion:
		return "观点表达"
	case PhaseDebate:
		return "辩论"
	case PhaseConsensus:
		return "共识"
	case PhaseDecision:
		return "决策"
	case PhaseSummary:
		return "总结"
	default:
		return string(phase)
	}
}

// PhaseDescription 获取阶段描述
func PhaseDescription(phase Phase) string {
	switch phase {
	case PhaseInfoGathering:
		return "各Agent陈述已知事实和数据"
	case PhaseOpinion:
		return "各Agent提出立场和方案"
	case PhaseDebate:
		return "Agent间交叉质疑和数据支撑"
	case PhaseConsensus:
		return "寻找共同点，推动达成一致"
	case PhaseDecision:
		return "关键议题投票或最终判断"
	case PhaseSummary:
		return "生成本轮推演报告"
	default:
		return ""
	}
}
