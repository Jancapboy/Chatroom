import { useEffect, useState, useCallback } from 'react';
import { ArrowLeft, Menu, X } from 'lucide-react';
import { AgentPanel } from '../components/room/AgentPanel';
import { ChatStream } from '../components/room/ChatStream';
import { ConsensusGauge } from '../components/room/ConsensusGauge';
import { RoomControls } from '../components/room/RoomControls';
import { SnapshotPanel } from '../components/room/SnapshotPanel';
import { Timeline } from '../components/room/Timeline';
import { ActionPanel } from '../components/room/ActionPanel';
import { useRoomStore } from '../stores/useRoomStore';
import { useWebSocket } from '../hooks/useWebSocket';
import { roomApi } from '../lib/api';
import type { Message, Phase } from '../types/room';
import type { AgentMessagePayload, AgentStatePayload, ConsensusUpdatePayload, PhaseChangePayload } from '../types/ws';

interface RoomPageProps {
  onNavigate: (page: string) => void;
}

export function RoomPage({ onNavigate }: RoomPageProps) {
  const roomId = window.location.pathname.split('/room/')[1] || '';
  const { currentRoom, messages, agents, consensus, fetchRoom, fetchMessages, addMessage, updateAgent, setConsensus, updatePhase } = useRoomStore();
  const { isConnected, lastMessage, sendUserMessage, sendCommand } = useWebSocket(roomId);

  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [showActionPanel, setShowActionPanel] = useState(false);

  // Fetch room + messages on mount
  useEffect(() => {
    if (roomId) {
      fetchRoom(roomId);
      fetchMessages(roomId);
    }
  }, [roomId, fetchRoom, fetchMessages]);

  // Show action panel when completed
  useEffect(() => {
    if (currentRoom?.status === 'completed') {
      setShowActionPanel(true);
    }
  }, [currentRoom?.status]);

  // Handle WS messages
  useEffect(() => {
    if (!lastMessage) return;

    switch (lastMessage.type) {
      case 'message': {
        const p = lastMessage.payload as AgentMessagePayload;
        const msg: Message = {
          id: p.id,
          roomId,
          senderId: p.sender_id,
          senderType: p.sender_type,
          senderName: p.sender_name,
          senderAvatar: p.sender_avatar,
          content: p.content,
          type: 'text',
          phase: p.phase as Phase,
          round: p.round,
          timestamp: Date.now(),
          metadata: {
            confidence: p.confidence,
            stance: p.stance,
          },
        };
        addMessage(msg);
        break;
      }
      case 'agent_state': {
        const p = lastMessage.payload as AgentStatePayload;
        const agent = agents.find((a) => a.id === p.agent_id);
        if (agent) {
          updateAgent({
            ...agent,
            energy: p.energy,
            confidence: p.confidence,
            stance: p.stance as 'support' | 'oppose' | 'neutral',
          });
        }
        break;
      }
      case 'consensus_update': {
        const p = lastMessage.payload as ConsensusUpdatePayload;
        setConsensus({
          topic: p.topic,
          agreement: Math.round(p.agreement),
          breakdown: {},
        });
        break;
      }
      case 'phase_change': {
        const p = lastMessage.payload as PhaseChangePayload;
        updatePhase(p.phase, p.round);
        break;
      }
      case 'system': {
        break;
      }
    }
  }, [lastMessage, addMessage, updateAgent, setConsensus, updatePhase, agents, roomId]);

  const handleSendMessage = useCallback((content: string) => {
    sendUserMessage(content);
    addMessage({
      id: `user-${Date.now()}`,
      roomId,
      senderId: 'user-1',
      senderType: 'user',
      senderName: 'Observer',
      content,
      type: 'text',
      phase: currentRoom?.currentPhase || 'debate',
      round: currentRoom?.currentRound || 1,
      timestamp: Date.now(),
    });
  }, [sendUserMessage, addMessage, roomId, currentRoom]);

  const handleStart = useCallback(async () => {
    if (!roomId) return;
    try {
      await roomApi.startRoom(roomId);
      fetchRoom(roomId);
    } catch (err) {
      console.error('start room error:', err);
      alert('启动推演失败: ' + String(err));
    }
  }, [roomId, fetchRoom]);

  const handlePause = useCallback(() => sendCommand('pause'), [sendCommand]);
  const handleResume = useCallback(() => sendCommand('resume'), [sendCommand]);

  const handleCreateSnapshot = useCallback(async () => {
    if (!roomId) return;
    try {
      await roomApi.createSnapshot(roomId, '用户手动创建');
      alert('检查点已创建');
    } catch (err) {
      console.error('create snapshot error:', err);
      alert('创建检查点失败');
    }
  }, [roomId]);

  const handleRollback = useCallback(async (snapshotId: string) => {
    if (!roomId) return;
    const reason = window.prompt('回滚原因（可选）：') || '';
    try {
      const res = await roomApi.rollbackRoom(roomId, snapshotId, reason);
      const newRoom = res.data as { id?: string };
      if (newRoom?.id) {
        alert(`已回滚到新房间，ID: ${newRoom.id}`);
        onNavigate('lobby');
      }
    } catch (err) {
      console.error('rollback error:', err);
      alert('回滚失败');
    }
  }, [roomId, onNavigate]);

  const handleForkFromSnapshot = useCallback(async (snapshotId: string) => {
    if (!roomId) return;
    const reason = window.prompt('分支原因（可选）：') || '';
    try {
      const res = await roomApi.forkFromSnapshot(roomId, snapshotId, reason);
      const newRoom = res.data as { id?: string };
      if (newRoom?.id) {
        alert(`分支创建成功，新房间ID: ${newRoom.id}`);
        onNavigate('lobby');
      }
    } catch (err) {
      console.error('fork error:', err);
      alert('分支创建失败');
    }
  }, [roomId, onNavigate]);

  if (!currentRoom) {
    return (
      <div className="min-h-screen bg-bg-primary flex items-center justify-center text-text-muted">
        <div className="animate-pulse">加载中...</div>
      </div>
    );
  }

  return (
    <div className="h-screen bg-bg-primary flex flex-col overflow-hidden">
      {/* Top bar — minimal */}
      <div className="bg-bg-secondary border-b border-border-subtle px-4 py-3 flex items-center justify-between shrink-0 z-10">
        <div className="flex items-center gap-3 min-w-0">
          <button
            onClick={() => onNavigate('lobby')}
            className="p-2 hover:bg-bg-elevated rounded-md transition-colors shrink-0"
          >
            <ArrowLeft size={18} className="text-text-secondary" />
          </button>

          <div className="min-w-0">
            <h1 className="text-base font-semibold text-text-primary truncate">{currentRoom.name}</h1>
            <div className="flex items-center gap-2 mt-0.5">
              <span className="text-xs text-text-muted">
                第 {currentRoom.currentRound} / {currentRoom.maxRounds} 轮
              </span>
              <span
                className={`text-xs px-2 py-0.5 rounded-full ${
                  isConnected
                    ? 'bg-agent-analyst/10 text-agent-analyst'
                    : 'bg-agent-risk/10 text-agent-risk'
                }`}
              >
                {isConnected ? '已连接' : '未连接'}
              </span>
            </div>
          </div>
        </div>

        {/* Hamburger — opens sidebar */}
        <button
          onClick={() => setSidebarOpen(true)}
          className="p-2 hover:bg-bg-elevated rounded-md transition-colors shrink-0"
          aria-label="展开侧边栏"
        >
          <Menu size={18} className="text-text-secondary" />
        </button>
      </div>

      {/* Main content — chat takes all space */}
      <div className="flex-1 min-h-0 flex flex-col">
        <ChatStream messages={messages} />
      </div>

      {/* Bottom controls */}
      <div className="shrink-0 z-10">
        <RoomControls
          status={currentRoom.status}
          onStart={handleStart}
          onPause={handlePause}
          onResume={handleResume}
          onSendMessage={handleSendMessage}
          onCreateSnapshot={handleCreateSnapshot}
          onShowResults={() => setShowActionPanel(true)}
        />
      </div>

      {/* Sidebar overlay */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50"
          onClick={() => setSidebarOpen(false)}
        />
      )}
      <aside
        className={`fixed top-0 right-0 h-full w-80 max-w-[85vw] bg-bg-secondary border-l border-border-subtle z-50 flex flex-col overflow-hidden transition-transform duration-300 ease-out ${
          sidebarOpen ? 'translate-x-0' : 'translate-x-full'
        }`}
      >
        {/* Sidebar header */}
        <div className="px-4 py-3 border-b border-border-subtle flex items-center justify-between shrink-0">
          <h2 className="text-sm font-semibold text-text-secondary">高级功能</h2>
          <button
            onClick={() => setSidebarOpen(false)}
            className="p-1.5 hover:bg-bg-elevated rounded-md transition-colors"
          >
            <X size={16} className="text-text-muted" />
          </button>
        </div>

        {/* Sidebar content */}
        <div className="flex-1 overflow-y-auto">
          <div className="p-3 space-y-3">
            <AgentPanel agents={agents} />
            <ConsensusGauge consensus={consensus} />
            <Timeline
              currentRound={currentRoom.currentRound}
              maxRounds={currentRoom.maxRounds}
              onViewSnapshot={(round) => {
                alert(`查看第 ${round} 轮快照（可在检查点面板查看详情）`);
              }}
              onForkFromRound={(round) => {
                alert(`从第 ${round} 轮创建分支（请使用检查点面板操作）`);
              }}
            />
            <SnapshotPanel
              roomId={roomId}
              onFork={handleForkFromSnapshot}
              onRollback={handleRollback}
            />
          </div>
        </div>
      </aside>

      {/* Action panel — shown when completed */}
      {showActionPanel && (
        <ActionPanel
          conclusion={consensus?.topic || currentRoom.topic || '推演已完成'}
          actions={[
            { id: '1', title: '采纳结论', description: '将推演结论应用到实际决策中' },
            { id: '2', title: '创建分支', description: '基于当前状态创建新的推演分支' },
            { id: '3', title: '导出报告', description: '生成完整的推演过程报告' },
          ]}
          onExecute={(id) => {
            if (id === '2') {
              sendCommand('fork');
            }
            setShowActionPanel(false);
          }}
          onSkip={() => setShowActionPanel(false)}
          onClose={() => setShowActionPanel(false)}
        />
      )}
    </div>
  );
}
