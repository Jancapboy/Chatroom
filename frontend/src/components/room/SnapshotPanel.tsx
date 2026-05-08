import { useState, useEffect } from 'react';
import { Camera, GitBranch, RotateCcw, X, ChevronDown, ChevronRight, AlertTriangle } from 'lucide-react';
import { roomApi } from '../../lib/api';

interface Snapshot {
  id: string;
  room_id: string;
  round: number;
  phase: string;
  consensus_score: number;
  trigger_reason: string;
  created_at: string;
}

interface SnapshotPanelProps {
  roomId: string;
  onFork?: (snapshotId: string) => void;
  onRollback?: (snapshotId: string) => void;
}

export function SnapshotPanel({ roomId, onFork, onRollback }: SnapshotPanelProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const [snapshots, setSnapshots] = useState<Snapshot[]>([]);
  const [loading, setLoading] = useState(false);
  const [activeMenu, setActiveMenu] = useState<string | null>(null);

  const fetchSnapshots = async () => {
    if (!roomId) return;
    setLoading(true);
    try {
      const res = await roomApi.getSnapshots(roomId);
      const data = (res.data as { list?: Snapshot[] })?.list || [];
      setSnapshots(data);
    } catch (err) {
      console.error('fetchSnapshots error:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isExpanded) {
      fetchSnapshots();
    }
  }, [isExpanded, roomId]);

  const handleCreateSnapshot = async () => {
    try {
      await roomApi.createSnapshot(roomId, '用户手动创建');
      fetchSnapshots();
    } catch (err) {
      console.error('createSnapshot error:', err);
      alert('创建检查点失败');
    }
  };

  const formatTime = (ts: string) => {
    const d = new Date(ts);
    return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', second: '2-digit' });
  };

  return (
    <div className="bg-bg-secondary border border-border-subtle rounded-lg overflow-hidden">
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className="w-full flex items-center justify-between px-4 py-3 hover:bg-bg-elevated transition-colors"
      >
        <div className="flex items-center gap-2">
          <Camera size={16} className="text-text-muted" />
          <span className="text-sm font-medium text-text-secondary">检查点 / 回溯</span>
          {snapshots.length > 0 && (
            <span className="text-[10px] px-1.5 py-0.5 rounded-full bg-accent-cyan/10 text-accent-cyan">
              {snapshots.length}
            </span>
          )}
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
        <div className="px-4 pb-3 space-y-2">
          <button
            onClick={handleCreateSnapshot}
            className="w-full flex items-center justify-center gap-2 px-3 py-2 bg-accent-cyan/10 text-accent-cyan rounded-md text-xs font-medium hover:bg-accent-cyan/20 transition-colors border border-accent-cyan/20"
          >
            <Camera size={14} />
            创建检查点
          </button>

          {loading && (
            <div className="text-xs text-text-muted text-center py-2">加载中...</div>
          )}

          {!loading && snapshots.length === 0 && (
            <div className="text-xs text-text-muted text-center py-2">暂无检查点</div>
          )}

          {snapshots.map((snap) => (
            <div key={snap.id} className="relative">
              <button
                onClick={() => setActiveMenu(activeMenu === snap.id ? null : snap.id)}
                className={`w-full flex items-center gap-2 px-3 py-2 rounded-md text-xs transition-colors ${
                  snap.trigger_reason.includes('共识度') || snap.trigger_reason.includes('风险官')
                    ? 'bg-agent-risk/5 hover:bg-agent-risk/10 border border-agent-risk/20'
                    : 'bg-bg-elevated hover:bg-bg-card border border-border-subtle'
                }`}
              >
                {(snap.trigger_reason.includes('共识度') || snap.trigger_reason.includes('风险官')) && (
                  <AlertTriangle size={12} className="text-agent-risk shrink-0" />
                )}
                <div className="flex-1 text-left">
                  <div className="flex items-center gap-1.5">
                    <span className="font-medium">第 {snap.round} 轮</span>
                    <span className="text-text-muted">· {snap.phase}</span>
                  </div>
                  <div className="text-text-muted mt-0.5 truncate">
                    {snap.trigger_reason}
                  </div>
                </div>
                <span className="text-[10px] text-text-muted shrink-0">
                  {formatTime(snap.created_at)}
                </span>
              </button>

              {activeMenu === snap.id && (
                <div className="mt-1 flex gap-1 px-1 transition-opacity duration-150">
                  <button
                    onClick={() => {
                      onRollback?.(snap.id);
                      setActiveMenu(null);
                    }}
                    className="flex-1 flex items-center justify-center gap-1 px-2 py-1.5 bg-accent-warn/10 text-accent-warn rounded text-[11px] hover:bg-accent-warn/20 transition-colors"
                  >
                    <RotateCcw size={12} />
                    回滚
                  </button>
                  <button
                    onClick={() => {
                      onFork?.(snap.id);
                      setActiveMenu(null);
                    }}
                    className="flex-1 flex items-center justify-center gap-1 px-2 py-1.5 bg-accent-purple/10 text-accent-purple rounded text-[11px] hover:bg-accent-purple/20 transition-colors"
                  >
                    <GitBranch size={12} />
                    分支
                  </button>
                  <button
                    onClick={() => setActiveMenu(null)}
                    className="p-1.5 text-text-muted hover:text-text-secondary transition-colors"
                  >
                    <X size={12} />
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
