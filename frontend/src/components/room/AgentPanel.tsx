import { motion } from 'framer-motion';
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
    <div className="h-full flex flex-col bg-bg-secondary border-r border-border-subtle">
      <div className="p-4 border-b border-border-subtle">
        <h2 className="text-sm font-semibold text-text-secondary uppercase tracking-wider">
          智能体状态
        </h2>
        <p className="text-xs text-text-muted mt-1">{agents.length} 个活跃Agent</p>
      </div>

      <div className="flex-1 overflow-y-auto p-3 space-y-3">
        {agents.map((agent, index) => {
          const Icon = roleIcons[agent.role] || Zap;
          return (
            <motion.div
              key={agent.id}
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: index * 0.1 }}
              className="bg-bg-elevated rounded-lg p-3 border border-border-subtle hover:border-border-glow transition-colors"
            >
              {/* Header */}
              <div className="flex items-center gap-3">
                <div
                  className="w-10 h-10 rounded-full flex items-center justify-center text-white text-sm font-bold shrink-0"
                  style={{ backgroundColor: agent.color }}
                >
                  {agent.name[0]}
                </div>
                <div className="min-w-0">
                  <div className="font-medium text-text-primary text-sm truncate">
                    {agent.name}
                  </div>
                  <div className="text-xs text-text-muted flex items-center gap-1">
                    <Icon size={12} />
                    {roleLabels[agent.role]}
                  </div>
                </div>
              </div>

              {/* Role color bar */}
              <div
                className="h-1 rounded-full mt-2"
                style={{ backgroundColor: agent.color, opacity: 0.6 }}
              />

              {/* Stats */}
              <div className="mt-3 space-y-2">
                {/* Energy */}
                <div>
                  <div className="flex justify-between text-xs mb-1">
                    <span className="text-text-muted">能量</span>
                    <span className="text-text-secondary">{agent.energy}%</span>
                  </div>
                  <div className="h-1.5 bg-bg-primary rounded-full overflow-hidden">
                    <motion.div
                      className="h-full rounded-full"
                      style={{ backgroundColor: agent.color }}
                      initial={{ width: 0 }}
                      animate={{ width: `${agent.energy}%` }}
                      transition={{ duration: 1, delay: index * 0.1 }}
                    />
                  </div>
                </div>

                {/* Confidence */}
                <div>
                  <div className="flex justify-between text-xs mb-1">
                    <span className="text-text-muted">置信度</span>
                    <span className="text-text-secondary">{agent.confidence}%</span>
                  </div>
                  <div className="h-1.5 bg-bg-primary rounded-full overflow-hidden">
                    <motion.div
                      className="h-full rounded-full bg-accent-cyan"
                      initial={{ width: 0 }}
                      animate={{ width: `${agent.confidence}%` }}
                      transition={{ duration: 1, delay: index * 0.1 + 0.2 }}
                    />
                  </div>
                </div>
              </div>

              {/* Stance badge */}
              <div className="mt-3 flex items-center justify-between">
                <span
                  className={`text-xs px-2 py-0.5 rounded-full text-white font-medium ${stanceColors[agent.stance]}`}
                >
                  {stanceLabels[agent.stance]}
                </span>
                {agent.isActive && (
                  <span className="w-2 h-2 rounded-full bg-agent-analyst animate-pulse" />
                )}
              </div>
            </motion.div>
          );
        })}
      </div>
    </div>
  );
}
