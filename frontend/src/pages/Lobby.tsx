import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { Plus, Search, Filter, Flame, Clock, CheckCircle, PlayCircle } from 'lucide-react';
import { RoomCard } from '../components/lobby/RoomCard';
import { useLobbyStore } from '../stores/useLobbyStore';
import type { RoomStatus } from '../types/room';

const filters: { key: RoomStatus | 'all'; label: string; icon: typeof Clock }[] = [
  { key: 'all', label: '全部', icon: Filter },
  { key: 'running', label: '进行中', icon: PlayCircle },
  { key: 'preparing', label: '准备中', icon: Clock },
  { key: 'completed', label: '已完成', icon: CheckCircle },
];

interface LobbyProps {
  onNavigate: (page: string, params?: Record<string, string>) => void;
}

export function Lobby({ onNavigate }: LobbyProps) {
  const { filter, setFilter, getFilteredRooms } = useLobbyStore();
  const [searchQuery, setSearchQuery] = useState('');
  const filtered = getFilteredRooms().filter(
    (r) =>
      r.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      r.topic.toLowerCase().includes(searchQuery.toLowerCase())
  );

  useEffect(() => {
    useLobbyStore.getState().fetchRooms();
  }, []);

  return (
    <div className="min-h-screen bg-bg-primary">
      {/* Header */}
      <header className="border-b border-border-subtle bg-bg-secondary/80 backdrop-blur sticky top-0 z-10">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-lg bg-accent-cyan/10 flex items-center justify-center">
              <Flame size={20} className="text-accent-cyan" />
            </div>
            <div>
              <h1 className="text-lg font-bold text-text-primary">ASI Chatroom</h1>
              <p className="text-xs text-text-muted">多智能体推演平台</p>
            </div>
          </div>

          <motion.button
            whileTap={{ scale: 0.95 }}
            onClick={() => onNavigate('create')}
            className="flex items-center gap-2 px-4 py-2 bg-accent-cyan text-bg-primary rounded-md text-sm font-semibold hover:bg-accent-cyan-dim transition-colors"
          >
            <Plus size={16} />
            创建房间
          </motion.button>
        </div>
      </header>

      <div className="max-w-7xl mx-auto px-4 py-6">
        {/* Search + Filters */}
        <div className="flex flex-col sm:flex-row gap-4 mb-6">
          <div className="relative flex-1">
            <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="搜索房间..."
              className="w-full bg-bg-secondary border border-border-subtle rounded-md pl-10 pr-4 py-2.5 text-sm text-text-primary placeholder:text-text-muted focus:outline-none focus:border-accent-cyan transition-colors"
            />
          </div>

          <div className="flex gap-2">
            {filters.map(({ key, label, icon: Icon }) => (
              <button
                key={key}
                onClick={() => setFilter(key)}
                className={`flex items-center gap-1.5 px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                  filter === key
                    ? 'bg-accent-cyan/10 text-accent-cyan border border-accent-cyan/30'
                    : 'bg-bg-secondary text-text-secondary border border-border-subtle hover:bg-bg-elevated'
                }`}
              >
                <Icon size={14} />
                {label}
              </button>
            ))}
          </div>
        </div>

        {/* Room grid */}
        {filtered.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {filtered.map((room, index) => (
              <motion.div
                key={room.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.05 }}
              >
                <RoomCard
                  room={room}
                  onEnter={(id) => onNavigate('room', { roomId: id })}
                />
              </motion.div>
            ))}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center py-20 text-text-muted">
            <Search size={48} className="mb-4 opacity-30" />
            <p className="text-lg font-medium">暂无房间</p>
            <p className="text-sm mt-1">创建一个房间开始推演吧</p>
          </div>
        )}
      </div>
    </div>
  );
}
