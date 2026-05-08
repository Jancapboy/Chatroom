import type { Phase } from '../../types/room';

const phases: { key: Phase; label: string }[] = [
  { key: 'info_gathering', label: '信息收集' },
  { key: 'opinion_expression', label: '观点表达' },
  { key: 'debate', label: '辩论' },
  { key: 'consensus', label: '共识' },
  { key: 'decision', label: '决策' },
  { key: 'summary', label: '总结' },
];

interface PhaseIndicatorProps {
  currentPhase: Phase;
  currentRound: number;
  maxRounds: number;
}

export function PhaseIndicator({ currentPhase, currentRound, maxRounds }: PhaseIndicatorProps) {
  const currentIndex = phases.findIndex((p) => p.key === currentPhase);

  return (
    <div className="bg-bg-secondary border-b border-border-subtle px-4 py-3">
      <div className="flex items-center gap-2 mb-2">
        <span className="text-xs text-text-muted">第 {currentRound} / {maxRounds} 轮</span>
        <span className="text-xs text-accent-cyan font-medium">
          {phases[currentIndex]?.label}
        </span>
      </div>
      <div className="flex gap-1">
        {phases.map((phase, index) => {
          const isActive = index === currentIndex;
          const isCompleted = index < currentIndex;

          return (
            <div
              key={phase.key}
              className={`flex-1 h-2 rounded-full ${
                isActive
                  ? 'bg-accent-cyan'
                  : isCompleted
                    ? 'bg-accent-cyan/40'
                    : 'bg-bg-elevated'
              }`}
            >
              <div className="relative group">
                <div className="absolute -top-7 left-1/2 -translate-x-1/2 opacity-0 group-hover:opacity-100 transition-opacity bg-bg-elevated text-text-secondary text-[10px] px-2 py-0.5 rounded whitespace-nowrap border border-border-subtle">
                  {phase.label}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
