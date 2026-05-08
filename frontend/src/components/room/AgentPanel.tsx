import { Zap, Target, Shield, Brain, BarChart3, Hammer } from 'lucide-react';
import type { Agent } from '../../types/room';

const roleIcons = {
  architect: Brain,
  risk_officer: Shield,
  strategist: Target,
  analyst: BarChart3,
  executor: Hammer,
};

const roleLabels: Record<string, string> = {
  architect: '架构师',
  risk_officer: '风险官',
  strategist: '策略家',
  analyst: '分析师',
  executor: '执行者',
};

const stanceLabels: Record<string, string> = {
  support: '支持',
  oppose: '反对',
  neutral: '中立',
};

const stanceColors: Record<string, string> = {
  support: 'bg-agent-analyst',
  oppose: 'bg-agent-risk',
  neutral: 'bg-agent-executor',
};

interface AgentPanelProps {
  agents: Agent[];
}

export function AgentPanel({ agents }: AgentPanelProps) {
  return (
    <div className="flex flex-col bg-bg-secondary rounded-lg border border-border-subtle overflow-hidden">
      <div className="p-3 border-b border-border-subtle">
        <h2 className="text-xs font-semibold text-text-secondary uppercase tracking-wider">
          智能体状态
        </h2>
        <p className="text-[11px] text-text-muted mt-0.5">{agents.length} 个活跃Agent</p>
      </div>

      <div className="p-2 space-y-2 max-h-64 overflow-y-auto">
        {agents.map((agent) => {
          const Icon = roleIcons[agent.role] || Zap;
          return (
            <div
              key={agent.id}
              className="bg-bg-elevated rounded-lg p-2.5 border border-border-subtle hover:border-border-glow transition-colors"
            >
              <div className="flex items-center gap-2.5">
                <div
                  className="w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0"
                  style={{ backgroundColor: agent.color }}
                >
                  {agent.name[0]}
                </div>
                <div className="min-w-0 flex-1">
                  <div className="font-medium text-text-primary text-xs truncate">
                    {agent.name}
                  </div>
                  <div className="text-[10px] text-text-muted flex items-center gap-1">
                    <Icon size={10} />
                    {roleLabels[agent.role]}
                  </div>
                </div>
                {agent.isActive && (
                  <span className="w-1.5 h-1.5 rounded-full bg-agent-analyst animate-pulse shrink-0" />
                )}
              </div>

              {/* Stats */}
              <div className="mt-2 space-y-1.5">
                <div>
                  <div className="flex justify-between text-[10px] mb-0.5">
                    <span className="text-text-muted">能量</span>
                    <span className="text-text-secondary">{agent.energy}%</span>
                  </div>
                  <div className="h-1 bg-bg-primary rounded-full overflow-hidden">
                    <div
                      className="h-full rounded-full transition-all duration-500"
                      style={{ width: `${agent.energy}%`, backgroundColor: agent.color }}
                    />
                  </div>
                </div>
                <div>
                  <div className="flex justify-between text-[10px] mb-0.5">
                    <span className="text-text-muted">置信度</span>
                    <span className="text-text-secondary">{agent.confidence}%</span>
                  </div>
                  <div className="h-1 bg-bg-primary rounded-full overflow-hidden">
                    <div
                      className="h-full rounded-full bg-accent-cyan transition-all duration-500"
                      style={{ width: `${agent.confidence}%` }}
                    />
                  </div>
                </div>
              </div>

              <div className="mt-2">
                <span
                  className={`text-[10px] px-1.5 py-0.5 rounded-full text-white font-medium ${stanceColors[agent.stance]}`}
                >
                  {stanceLabels[agent.stance]}
                </span>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
