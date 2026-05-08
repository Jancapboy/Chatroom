import { Users, MessageCircle, ArrowRight, Clock, CheckCircle, PauseCircle, PlayCircle } from 'lucide-react';
import type { Room } from '../../types/room';

const statusConfig = {
  preparing: { icon: Clock, color: 'text-accent-warn', bg: 'bg-accent-warn/10', label: '准备中' },
  running: { icon: PlayCircle, color: 'text-agent-analyst', bg: 'bg-agent-analyst/10', label: '推演中' },
  paused: { icon: PauseCircle, color: 'text-accent-warn', bg: 'bg-accent-warn/10', label: '已暂停' },
  completed: { icon: CheckCircle, color: 'text-accent-cyan', bg: 'bg-accent-cyan/10', label: '已完成' },
  archived: { icon: Clock, color: 'text-text-muted', bg: 'bg-text-muted/10', label: '已归档' },
};

interface RoomCardProps {
  room: Room;
  onEnter: (roomId: string) => void;
}

export function RoomCard({ room, onEnter }: RoomCardProps) {
  const status = statusConfig[room.status];
  const StatusIcon = status.icon;

  return (
    <div
      className="bg-bg-secondary rounded-lg border border-border-subtle overflow-hidden hover:border-border-glow hover:-translate-y-1 hover:shadow-[0_8px_30px_rgba(0,212,170,0.08)] transition-all duration-200 cursor-pointer"
      onClick={() => onEnter(room.id)}
    >
      {/* Header */}
      <div className="p-4 pb-3">
        <div className="flex items-start justify-between mb-2">
          <h3 className="text-base font-semibold text-text-primary line-clamp-1">
            {room.name}
          </h3>
          <span
            className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-medium ${status.bg} ${status.color}`}
          >
            <StatusIcon size={12} />
            {status.label}
          </span>
        </div>
        <p className="text-sm text-text-secondary line-clamp-2 mb-3">
          {room.description}
        </p>

        {/* Topic tag */}
        <div className="inline-block px-2 py-1 bg-bg-elevated rounded text-xs text-text-muted border border-border-subtle">
          {room.topic}
        </div>
      </div>

      {/* Footer */}
      <div className="px-4 py-3 bg-bg-elevated/50 border-t border-border-subtle flex items-center justify-between">
        <div className="flex items-center gap-3">
          {/* Agent avatars */}
          <div className="flex -space-x-2">
            {room.agents.slice(0, 4).map((agent) => (
              <div
                key={agent.id}
                className="w-7 h-7 rounded-full border-2 border-bg-secondary flex items-center justify-center text-[10px] font-bold text-white"
                style={{ backgroundColor: agent.color }}
                title={agent.name}
              >
                {agent.name[0]}
              </div>
            ))}
            {room.agents.length > 4 && (
              <div className="w-7 h-7 rounded-full border-2 border-bg-secondary bg-bg-card flex items-center justify-center text-[10px] text-text-muted">
                +{room.agents.length - 4}
              </div>
            )}
          </div>

          {/* Stats */}
          <div className="flex items-center gap-2 text-xs text-text-muted">
            <span className="flex items-center gap-1">
              <Users size={12} />
              {room.agents.length}
            </span>
            {room.messageCount !== undefined && (
              <span className="flex items-center gap-1">
                <MessageCircle size={12} />
                {room.messageCount}
              </span>
            )}
          </div>
        </div>

        <div className="flex items-center gap-1 text-sm text-accent-cyan font-medium hover:translate-x-1 transition-transform">
          进入
          <ArrowRight size={16} />
        </div>
      </div>
    </div>
  );
}
