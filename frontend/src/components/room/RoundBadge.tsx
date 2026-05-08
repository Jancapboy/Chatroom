import { RotateCcw } from 'lucide-react';

interface RoundBadgeProps {
  round: number;
  maxRounds: number;
}

export function RoundBadge({ round, maxRounds }: RoundBadgeProps) {
  return (
    <div className="inline-flex items-center gap-2 px-3 py-1.5 bg-bg-elevated rounded-full border border-border-subtle">
      <RotateCcw size={14} className="text-accent-cyan" />
      <span className="text-sm font-semibold text-accent-cyan">
        第 {round} / {maxRounds} 轮
      </span>
    </div>
  );
}
