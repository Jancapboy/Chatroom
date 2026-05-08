import { useState } from 'react';
import { Clock, ChevronDown, ChevronRight, Camera, GitBranch } from 'lucide-react';

interface TimelineProps {
  currentRound: number;
  maxRounds: number;
  onJumpToRound?: (round: number) => void;
  onViewSnapshot?: (round: number) => void;
  onForkFromRound?: (round: number) => void;
}

export function Timeline({ currentRound, maxRounds, onJumpToRound, onViewSnapshot, onForkFromRound }: TimelineProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const [menuRound, setMenuRound] = useState<number | null>(null);

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

      <div
        className={`overflow-hidden transition-all duration-200 ${isExpanded ? 'max-h-96' : 'max-h-0'}`}
      >
        <div className="px-4 pb-3 space-y-1">
          {rounds.map((round) => {
            const isCurrent = round === currentRound;
            const isPast = round < currentRound;
            const isFuture = round > currentRound;
            const isMenuOpen = menuRound === round;

            return (
              <div key={round} className="relative">
                <button
                  onClick={() => {
                    if (isPast || isCurrent) {
                      onJumpToRound?.(round);
                      setMenuRound(isMenuOpen ? null : round);
                    }
                  }}
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

                {isMenuOpen && (isPast || isCurrent) && (
                  <div className="ml-6 mt-1 flex gap-1 transition-opacity duration-150">
                    <button
                      onClick={() => {
                        onViewSnapshot?.(round);
                        setMenuRound(null);
                      }}
                      className="flex items-center gap-1 px-2 py-1 bg-bg-elevated text-text-secondary rounded text-[11px] hover:bg-bg-card border border-border-subtle transition-colors"
                    >
                      <Camera size={10} />
                      查看快照
                    </button>
                    <button
                      onClick={() => {
                        onForkFromRound?.(round);
                        setMenuRound(null);
                      }}
                      className="flex items-center gap-1 px-2 py-1 bg-accent-purple/10 text-accent-purple rounded text-[11px] hover:bg-accent-purple/20 transition-colors"
                    >
                      <GitBranch size={10} />
                      从此处分支
                    </button>
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
