import { useState } from 'react';
import { Pause, Play, MessageSquare, Rocket, Camera, CheckCircle } from 'lucide-react';
import type { RoomStatus } from '../../types/room';

interface RoomControlsProps {
  status: RoomStatus;
  onStart?: () => void;
  onPause?: () => void;
  onResume?: () => void;
  onSendMessage?: (content: string) => void;
  onCreateSnapshot?: () => void;
  onShowResults?: () => void;
}

export function RoomControls({
  status,
  onStart,
  onPause,
  onResume,
  onSendMessage,
  onCreateSnapshot,
  onShowResults,
}: RoomControlsProps) {
  const [inputValue, setInputValue] = useState('');
  const isPreparing = status === 'preparing';
  const isRunning = status === 'running';
  const isPaused = status === 'paused';
  const isCompleted = status === 'completed';

  const handleSend = () => {
    if (inputValue.trim() && onSendMessage) {
      onSendMessage(inputValue.trim());
      setInputValue('');
    }
  };

  return (
    <div className="bg-bg-secondary border-t border-border-subtle p-3 space-y-2">
      {/* Status row */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          {isPreparing && (
            <button
              onClick={onStart}
              className="flex items-center gap-2 px-4 py-2 bg-accent-cyan text-bg-primary rounded-md text-sm font-medium hover:bg-accent-cyan-dim transition-colors active:scale-95"
            >
              <Rocket size={16} />
              开始推演
            </button>
          )}

          {isRunning && (
            <button
              onClick={onPause}
              className="flex items-center gap-2 px-4 py-2 bg-accent-warn/10 text-accent-warn rounded-md text-sm font-medium hover:bg-accent-warn/20 transition-colors active:scale-95"
            >
              <Pause size={16} />
              暂停
            </button>
          )}

          {isPaused && (
            <button
              onClick={onResume}
              className="flex items-center gap-2 px-4 py-2 bg-agent-analyst/10 text-agent-analyst rounded-md text-sm font-medium hover:bg-agent-analyst/20 transition-colors active:scale-95"
            >
              <Play size={16} />
              继续
            </button>
          )}

          {isCompleted && (
            <button
              onClick={onShowResults}
              className="flex items-center gap-2 px-4 py-2 bg-accent-cyan/10 text-accent-cyan rounded-md text-sm font-medium hover:bg-accent-cyan/20 transition-colors active:scale-95"
            >
              <CheckCircle size={16} />
              推演完成，查看结果
            </button>
          )}

          {(isRunning || isPaused) && onCreateSnapshot && (
            <button
              onClick={onCreateSnapshot}
              className="flex items-center gap-2 px-3 py-2 bg-bg-elevated text-text-secondary rounded-md text-sm font-medium hover:bg-bg-card transition-colors border border-border-subtle active:scale-95"
            >
              <Camera size={16} />
              检查点
            </button>
          )}
        </div>

        {/* Status hint */}
        {isRunning && (
          <span className="text-xs text-text-muted animate-pulse">推演中...</span>
        )}
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
        <button
          onClick={handleSend}
          disabled={!inputValue.trim()}
          className="px-4 py-2 bg-accent-cyan text-bg-primary rounded-md text-sm font-medium hover:bg-accent-cyan-dim transition-colors disabled:opacity-50 disabled:cursor-not-allowed active:scale-95 flex items-center justify-center"
        >
          <MessageSquare size={16} />
        </button>
      </div>
    </div>
  );
}
