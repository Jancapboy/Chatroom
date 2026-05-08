import { useEffect, useState, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { ArrowLeft, Menu, X } from 'lucide-react';
import { AgentPanel } from '../components/room/AgentPanel';
import { ChatStream } from '../components/room/ChatStream';
import { ConsensusGauge } from '../components/room/ConsensusGauge';
import { PhaseIndicator } from '../components/room/PhaseIndicator';
import { RoundBadge } from '../components/room/RoundBadge';
import { Timeline } from '../components/room/Timeline';
import { RoomControls } from '../components/room/RoomControls';
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

  const [leftPanelOpen, setLeftPanelOpen] = useState(true);
  const [rightPanelOpen, setRightPanelOpen] = useState(true);
  const [isMobile, setIsMobile] = useState(false);

  // Handle responsive layout
  useEffect(() => {
    const checkMobile = () => {
      const mobile = window.innerWidth < 768;
      setIsMobile(mobile);
      if (mobile) {
        setLeftPanelOpen(false);
        setRightPanelOpen(false);
      } else {
        setLeftPanelOpen(true);
        setRightPanelOpen(true);
      }
    };
    checkMobile();
    window.addEventListener('resize', checkMobile);
    return () => window.removeEventListener('resize', checkMobile);
  }, []);

  // Fetch room + messages on mount
  useEffect(() => {
    if (roomId) {
      fetchRoom(roomId);
      fetchMessages(roomId);
    }
  }, [roomId, fetchRoom, fetchMessages]);

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
          breakdown: {}, // TODO: map breakdown from backend
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
  }, [lastMessage, addMessage, updateAgent, setConsensus, updatePhase, agents, currentRoom, roomId]);

  const handleSendMessage = useCallback((content: string) => {
    sendUserMessage(content);
    // Optimistically add to local state
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
      // Refresh room state to get updated status
      fetchRoom(roomId);
    } catch (err) {
      console.error('start room error:', err);
      alert('启动推演失败: ' + String(err));
    }
  }, [roomId, fetchRoom]);

  const handlePause = useCallback(() => sendCommand('pause'), [sendCommand]);
  const handleResume = useCallback(() => sendCommand('resume'), [sendCommand]);
  const handleFork = useCallback(() => sendCommand('fork'), [sendCommand]);
  const handleNextPhase = useCallback(() => sendCommand('next_phase'), [sendCommand]);

  if (!currentRoom) {
    return (
      <div className="min-h-screen bg-bg-primary flex items-center justify-center text-text-muted">
        <div className="animate-pulse">加载中...</div>
      </div>
    );
  }

  return (
    <div className="h-screen bg-bg-primary flex flex-col overflow-hidden">
      {/* Top bar */}
      <div className="bg-bg-secondary border-b border-border-subtle px-4 py-3 flex items-center justify-between shrink-0">
        <div className="flex items-center gap-3">
          <button
            onClick={() => onNavigate('lobby')}
            className="p-2 hover:bg-bg-elevated rounded-md transition-colors"
          >
            <ArrowLeft size={18} className="text-text-secondary" />
          </button>

          <div>
            <h1 className="text-base font-semibold text-text-primary">{currentRoom.name}</h1>
            <div className="flex items-center gap-2 mt-0.5">
              <RoundBadge round={currentRoom.currentRound} maxRounds={currentRoom.maxRounds} />
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

        {/* Mobile toggles */}
        {isMobile && (
          <div className="flex gap-2">
            <button
              onClick={() => setLeftPanelOpen(!leftPanelOpen)}
              className="p-2 hover:bg-bg-elevated rounded-md transition-colors"
            >
              {leftPanelOpen ? <X size={18} /> : <Menu size={18} />}
            </button>
          </div>
        )}
      </div>

      <PhaseIndicator
        currentPhase={currentRoom.currentPhase}
        currentRound={currentRoom.currentRound}
        maxRounds={currentRoom.maxRounds}
      />

      {/* Main content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Left panel - AgentPanel */}
        <AnimatePresence>
          {leftPanelOpen && (
            <motion.div
              initial={{ width: 0, opacity: 0 }}
              animate={{ width: 280, opacity: 1 }}
              exit={{ width: 0, opacity: 0 }}
              transition={{ duration: 0.2 }}
              className="shrink-0 overflow-hidden"
            >
              <AgentPanel agents={agents} />
            </motion.div>
          )}
        </AnimatePresence>

        {/* Center - ChatStream */}
        <div className="flex-1 min-w-0 flex flex-col">
          <ChatStream messages={messages} />
          <RoomControls
            status={currentRoom.status}
            onStart={handleStart}
            onPause={handlePause}
            onResume={handleResume}
            onFork={handleFork}
            onNextPhase={handleNextPhase}
            onSendMessage={handleSendMessage}
          />
        </div>

        {/* Right panel - Consensus + Timeline */}
        <AnimatePresence>
          {rightPanelOpen && (
            <motion.div
              initial={{ width: 0, opacity: 0 }}
              animate={{ width: 260, opacity: 1 }}
              exit={{ width: 0, opacity: 0 }}
              transition={{ duration: 0.2 }}
              className="shrink-0 overflow-hidden flex flex-col"
            >
              <ConsensusGauge consensus={consensus} />
              <div className="p-3 border-t border-border-subtle">
                <Timeline
                  currentRound={currentRoom.currentRound}
                  maxRounds={currentRoom.maxRounds}
                />
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}
