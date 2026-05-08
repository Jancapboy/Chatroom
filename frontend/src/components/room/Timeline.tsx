import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Clock, ChevronDown, ChevronRight } from 'lucide-react';

interface TimelineProps {
  currentRound: number;
  maxRounds: number;
  onJumpToRound?: (round: number) => void;
}

export function Timeline({ currentRound, maxRounds, onJumpToRound }: TimelineProps) {
  const [isExpanded, setIsExpanded] = useState(false);

  const rounds = Array.from({ length: maxRounds }, (_, i) => i + 1);

  return (
    <div className="bg-bg-secondary border border-border-subtle rounded-lg overflow-hidden">
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className="w-full flex items-center justify-between px-4 py-3 hover:bg-bg-elevated transition-colors"
      >
        <div className="flex items-center gap-2">
          <Clock size={16} className="text-text-muted" />
          <span className="text-sm font-medium text-text-secondary">时间线 / 回溯</span>
        </div>
        {isExpanded ? (
          <ChevronDown size={16} className="text-text-muted" />
        ) : (
          <ChevronRight size={16} className="text-text-muted" />
        )}
      </button>

      <AnimatePresence>
        {isExpanded && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="overflow-hidden"
          >
            <div className="px-4 pb-3 space-y-1">
              {rounds.map((round) => {
                const isCurrent = round === currentRound;
                const isPast = round < currentRound;
                const isFuture = round > currentRound;

                return (
                  <button
                    key={round}
                    onClick={() => onJumpToRound?.(round)}
                    disabled={isFuture}
                    className={`w-full flex items-center gap-3 px-3 py-2 rounded-md text-sm transition-colors ${
                      isCurrent
                        ? 'bg-accent-cyan/10 text-accent-cyan border border-accent-cyan/30'
                        : isPast
                          ? 'text-text-secondary hover:bg-bg-elevated'
                          : 'text-text-muted cursor-not-allowed opacity-50'
                    }`}
                  >
                    <span
                      className={`w-2 h-2 rounded-full shrink-0 ${
                        isCurrent
                          ? 'bg-accent-cyan'
                          : isPast
                            ? 'bg-text-muted'
                            : 'bg-border-subtle'
                      }`}
                    />
                    <span>第 {round} 轮</span>
                    {isCurrent && (
                      <span className="ml-auto text-[10px] px-1.5 py-0.5 rounded bg-accent-cyan/20 text-accent-cyan">
                        当前
                      </span>
                    )}
                  </button>
                );
              })}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
