import { motion } from 'framer-motion';
import { TrendingUp, TrendingDown, Minus } from 'lucide-react';
import type { ConsensusState } from '../../types/room';

interface ConsensusGaugeProps {
  consensus: ConsensusState | null;
}

export function ConsensusGauge({ consensus }: ConsensusGaugeProps) {
  const agreement = consensus?.agreement ?? 0;
  
  // Calculate stroke for circular gauge
  const radius = 52;
  const circumference = 2 * Math.PI * radius;
  const strokeDashoffset = circumference - (agreement / 100) * circumference;

  const getStatusColor = (value: number) => {
    if (value >= 80) return '#10b981'; // green
    if (value >= 50) return '#f59e0b'; // orange
    return '#ef4444'; // red
  };

  const statusColor = getStatusColor(agreement);

  return (
    <div className="h-full flex flex-col bg-bg-secondary border-l border-border-subtle">
      <div className="p-4 border-b border-border-subtle">
        <h2 className="text-sm font-semibold text-text-secondary uppercase tracking-wider">
          共识仪表
        </h2>
      </div>

      <div className="flex-1 overflow-y-auto p-4 space-y-6">
        {/* Circular gauge */}
        <div className="flex flex-col items-center">
          <div className="relative w-36 h-36">
            <svg className="w-full h-full -rotate-90" viewBox="0 0 120 120">
              {/* Background circle */}
              <circle
                cx="60"
                cy="60"
                r={radius}
                fill="none"
                stroke="#1a1a25"
                strokeWidth="8"
              />
              {/* Progress circle */}
              <motion.circle
                cx="60"
                cy="60"
                r={radius}
                fill="none"
                stroke={statusColor}
                strokeWidth="8"
                strokeLinecap="round"
                strokeDasharray={circumference}
                initial={{ strokeDashoffset: circumference }}
                animate={{ strokeDashoffset }}
                transition={{ duration: 1.5, ease: 'easeOut' }}
              />
            </svg>
            {/* Center text */}
            <div className="absolute inset-0 flex flex-col items-center justify-center">
              <motion.span
                className="text-3xl font-bold"
                style={{ color: statusColor }}
                initial={{ opacity: 0, scale: 0.5 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.5, delay: 0.5 }}
              >
                {agreement}%
              </motion.span>
              <span className="text-xs text-text-muted mt-1">共识度</span>
            </div>
          </div>

          {/* Topic */}
          {consensus?.topic && (
            <p className="text-sm text-text-secondary mt-4 text-center px-2">
              {consensus.topic}
            </p>
          )}
        </div>

        {/* Breakdown */}
        {consensus?.breakdown && Object.keys(consensus.breakdown).length > 0 && (
          <div className="space-y-3">
            <h3 className="text-xs font-semibold text-text-muted uppercase tracking-wider">
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
                  <div className="h-1.5 bg-bg-primary rounded-full overflow-hidden">
                    <motion.div
                      className="h-full rounded-full"
                      style={{
                        backgroundColor: score >= 70 ? '#10b981' : score <= 40 ? '#ef4444' : '#f59e0b',
                      }}
                      initial={{ width: 0 }}
                      animate={{ width: `${score}%` }}
                      transition={{ duration: 1 }}
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
