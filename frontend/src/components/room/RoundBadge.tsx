import { motion } from 'framer-motion';
import { RotateCcw } from 'lucide-react';

interface RoundBadgeProps {
  round: number;
  maxRounds: number;
}

export function RoundBadge({ round, maxRounds }: RoundBadgeProps) {
  return (
    <motion.div
      key={round}
      initial={{ scale: 0.8, opacity: 0 }}
      animate={{ scale: 1, opacity: 1 }}
      className="inline-flex items-center gap-2 px-3 py-1.5 bg-bg-elevated rounded-full border border-border-subtle"
    >
      <RotateCcw size={14} className="text-accent-cyan" />
      <span className="text-sm font-semibold text-accent-cyan">
        第 {round} / {maxRounds} 轮
      </span>
    </motion.div>
  );
}
