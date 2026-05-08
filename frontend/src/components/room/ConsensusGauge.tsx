import { TrendingUp, TrendingDown, Minus } from 'lucide-react';
import type { ConsensusState } from '../../types/room';

interface ConsensusGaugeProps {
  consensus: ConsensusState | null;
}

export function ConsensusGauge({ consensus }: ConsensusGaugeProps) {
  const agreement = consensus?.agreement ?? 0;

  const radius = 52;
  const circumference = 2 * Math.PI * radius;
  const strokeDashoffset = circumference - (agreement / 100) * circumference;

  const getStatusColor = (value: number) => {
    if (value >= 80) return '#10b981';
    if (value >= 50) return '#f59e0b';
    return '#ef4444';
  };

  const statusColor = getStatusColor(agreement);

  return (
    <div className="flex flex-col bg-bg-secondary rounded-lg border border-border-subtle overflow-hidden">
      <div className="p-3 border-b border-border-subtle">
        <h2 className="text-xs font-semibold text-text-secondary uppercase tracking-wider">
          共识仪表
        </h2>
      </div>

      <div className="p-3 space-y-4">
        {/* Circular gauge */}
        <div className="flex flex-col items-center">
          <div className="relative w-28 h-28">
            <svg className="w-full h-full -rotate-90" viewBox="0 0 120 120">
              <circle
                cx="60"
                cy="60"
                r={radius}
                fill="none"
                stroke="#1a1a25"
                strokeWidth="8"
              />
              <circle
                cx="60"
                cy="60"
                r={radius}
                fill="none"
                stroke={statusColor}
                strokeWidth="8"
                strokeLinecap="round"
                strokeDasharray={circumference}
                strokeDashoffset={strokeDashoffset}
                className="transition-all duration-1000 ease-out"
              />
            </svg>
            <div className="absolute inset-0 flex flex-col items-center justify-center">
              <span className="text-2xl font-bold" style={{ color: statusColor }}>
                {agreement}%
              </span>
              <span className="text-[10px] text-text-muted mt-0.5">共识度</span>
            </div>
          </div>

          {consensus?.topic && (
            <p className="text-xs text-text-secondary mt-2 text-center px-2">
              {consensus.topic}
            </p>
          )}
        </div>

        {/* Breakdown */}
        {consensus?.breakdown && Object.keys(consensus.breakdown).length > 0 && (
          <div className="space-y-2">
            <h3 className="text-[10px] font-semibold text-text-muted uppercase tracking-wider">
              各Agent立场
            </h3>
            {Object.entries(consensus.breakdown).map(([agentId, score]) => {
              const getIcon = () => {
                if (score >= 70) return <TrendingUp size={14} className="text-agent-analyst" />;
                if (score <= 40) return <TrendingDown size={14} className="text-agent-risk" />;
                return <Minus size={14} className="text-agent-executor" />;
              };
              return (
                <div key={agentId} className="space-y-1">
                  <div className="flex justify-between text-xs">
                    <span className="text-text-secondary">{agentId}</span>
                    <div className="flex items-center gap-1">
                      {getIcon()}
                      <span className="text-text-muted">{score}%</span>
                    </div>
                  </div>
                  <div className="h-1 bg-bg-primary rounded-full overflow-hidden">
                    <div
                      className="h-full rounded-full transition-all duration-500"
                      style={{
                        width: `${score}%`,
                        backgroundColor: score >= 70 ? '#10b981' : score <= 40 ? '#ef4444' : '#f59e0b',
                      }}
                    />
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
