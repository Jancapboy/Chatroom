import { motion } from 'framer-motion';
import { Pause, Play, GitBranch, Forward, MessageSquare, Rocket } from 'lucide-react';
import { useState } from 'react';
import type { RoomStatus } from '../../types/room';

interface RoomControlsProps {
  status: RoomStatus;
  onStart?: () => void;
  onPause?: () => void;
  onResume?: () => void;
  onFork?: () => void;
  onNextPhase?: () => void;
  onSendMessage?: (content: string) => void;
}

export function RoomControls({
  status,
  onStart,
  onPause,
  onResume,
  onFork,
  onNextPhase,
  onSendMessage,
}: RoomControlsProps) {
  const [inputValue, setInputValue] = useState('');
  const isPreparing = status === 'preparing';
  const isRunning = status === 'running';
  const isPaused = status === 'paused';

  const handleSend = () => {
    if (inputValue.trim() && onSendMessage) {
      onSendMessage(inputValue.trim());
      setInputValue('');
    }
  };

  return (
    <div className="bg-bg-secondary border-t border-border-subtle p-4 space-y-3">
      {/* Control buttons */}
      <div className="flex items-center gap-2 flex-wrap">
        {isPreparing && (
          <motion.button
            whileTap={{ scale: 0.95 }}
            onClick={onStart}
            className="flex items-center gap-2 px-4 py-2 bg-accent-cyan text-bg-primary rounded-md text-sm font-medium hover:bg-accent-cyan-dim transition-colors"
          >
            <Rocket size={16} />
            开始推演
          </motion.button>
        )}

        {isRunning && (
          <motion.button
            whileTap={{ scale: 0.95 }}
            onClick={onPause}
            className="flex items-center gap-2 px-4 py-2 bg-accent-warn/10 text-accent-warn rounded-md text-sm font-medium hover:bg-accent-warn/20 transition-colors"
          >
            <Pause size={16} />
            暂停
          </motion.button>
        )}

        {isPaused && (
          <motion.button
            whileTap={{ scale: 0.95 }}
            onClick={onResume}
            className="flex items-center gap-2 px-4 py-2 bg-agent-analyst/10 text-agent-analyst rounded-md text-sm font-medium hover:bg-agent-analyst/20 transition-colors"
          >
            <Play size={16} />
            继续
          </motion.button>
        )}

        {(isRunning || isPaused) && (
          <motion.button
            whileTap={{ scale: 0.95 }}
            onClick={onNextPhase}
            className="flex items-center gap-2 px-4 py-2 bg-bg-elevated text-text-secondary rounded-md text-sm font-medium hover:bg-bg-card transition-colors border border-border-subtle"
          >
            <Forward size={16} />
            下一阶段
          </motion.button>
        )}

        <motion.button
          whileTap={{ scale: 0.95 }}
          onClick={onFork}
          className="flex items-center gap-2 px-4 py-2 bg-accent-purple/10 text-accent-purple rounded-md text-sm font-medium hover:bg-accent-purple/20 transition-colors"
        >
          <GitBranch size={16} />
          Fork
        </motion.button>
      </div>

      {/* Chat input */}
      <div className="flex gap-2">
        <input
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleSend()}
          placeholder="插话 / 提问 / 指令..."
          className="flex-1 bg-bg-primary border border-border-subtle rounded-md px-3 py-2 text-sm text-text-primary placeholder:text-text-muted focus:outline-none focus:border-accent-cyan transition-colors"
        />
        <motion.button
          whileTap={{ scale: 0.95 }}
          onClick={handleSend}
          disabled={!inputValue.trim()}
          className="px-4 py-2 bg-accent-cyan text-bg-primary rounded-md text-sm font-medium hover:bg-accent-cyan-dim transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <MessageSquare size={16} />
        </motion.button>
      </div>
    </div>
  );
}
