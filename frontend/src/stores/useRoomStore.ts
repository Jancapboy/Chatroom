import { create } from 'zustand';
import type { Room, Agent, Message, Phase, ConsensusState } from '../types/room';
import { roomApi } from '../lib/api';

interface RoomState {
  currentRoom: Room | null;
  messages: Message[];
  agents: Agent[];
  consensus: ConsensusState | null;
  isConnected: boolean;
  isLoading: boolean;
  
  // Actions
  setRoom: (room: Room) => void;
  setMessages: (messages: Message[]) => void;
  addMessage: (message: Message) => void;
  updateAgent: (agent: Agent) => void;
  setConsensus: (consensus: ConsensusState) => void;
  setConnected: (connected: boolean) => void;
  setLoading: (loading: boolean) => void;
  fetchRoom: (roomId: string) => Promise<void>;
  fetchMessages: (roomId: string) => Promise<void>;
  updatePhase: (phase: Phase, round: number) => void;
}

function adaptAgent(raw: unknown): Agent {
  const a = raw as Record<string, unknown>;
  let expertise: string[] = [];
  try {
    expertise = JSON.parse(a.expertise as string);
  } catch { /* ignore */ }
  return {
    id: a.id as string,
    name: a.name as string,
    role: a.role as string as any,
    avatar: (a.avatar as string) || '',
    personality: a.personality as string,
    expertise,
    model: (a.model as string) || 'deepseek-chat',
    energy: (a.energy as number) || 100,
    confidence: (a.confidence as number) || 50,
    stance: (a.stance as string) as any,
    isActive: a.is_active !== false,
    color: getAgentColor(a.role as string),
  } as unknown as Agent;
}

function getAgentColor(role: string): string {
  const colors: Record<string, string> = {
    architect: '#3b82f6',
    risk_officer: '#ef4444',
    strategist: '#8b5cf6',
    analyst: '#10b981',
    executor: '#f59e0b',
  };
  return colors[role] || '#6b7280';
}

function adaptMessage(raw: unknown): Message {
  const m = raw as Record<string, unknown>;
  let metadata = {};
  try {
    metadata = JSON.parse((m.metadata as string) || '{}');
  } catch { /* ignore */ }
  return {
    id: m.id as string,
    roomId: m.room_id as string,
    senderId: m.sender_id as string,
    senderType: m.sender_type as 'agent' | 'user' | 'system',
    senderName: m.sender_name as string,
    senderAvatar: (m.sender_avatar as string) || undefined,
    content: m.content as string,
    type: m.msg_type as any,
    phase: (m.phase as string) as Phase || 'info_gathering',
    round: (m.round as number) || 1,
    timestamp: new Date(m.created_at as string).getTime(),
    metadata,
  };
}

export const useRoomStore = create<RoomState>((set) => ({
  currentRoom: null,
  messages: [],
  agents: [],
  consensus: null,
  isConnected: false,
  isLoading: false,

  setRoom: (room) => set({ currentRoom: room, agents: room.agents }),
  
  setMessages: (messages) => set({ messages }),
  
  addMessage: (message) => set((state) => ({ 
    messages: [...state.messages, message] 
  })),
  
  updateAgent: (agent) => set((state) => ({
    agents: state.agents.map((a) => (a.id === agent.id ? agent : a)),
    currentRoom: state.currentRoom ? {
      ...state.currentRoom,
      agents: state.currentRoom.agents.map((a) => (a.id === agent.id ? agent : a)),
    } : null,
  })),
  
  setConsensus: (consensus) => set({ consensus }),
  
  setConnected: (connected) => set({ isConnected: connected }),
  
  setLoading: (loading) => set({ isLoading: loading }),

  fetchRoom: async (roomId: string) => {
    set({ isLoading: true });
    try {
      const res = await roomApi.getRoom(roomId);
      const raw = res.data as unknown as Record<string, unknown>;
      if (!raw) {
        set({ isLoading: false });
        return;
      }
      const agents = (raw.agents as unknown[] || []).map(adaptAgent);
      const room: Room = {
        id: raw.id as string,
        name: raw.name as string,
        topic: (raw.topic as string) || '',
        description: (raw.description as string) || '',
        status: raw.status as any,
        templateId: (raw.template_id as string) || undefined,
        agents,
        currentPhase: (raw.current_phase as string) as Phase || 'info_gathering',
        currentRound: (raw.current_round as number) || 1,
        maxRounds: (raw.max_rounds as number) || 10,
        consensusScore: (raw.consensus_score as number) || 0,
        forkedFrom: (raw.forked_from as string) || undefined,
        createdBy: String(raw.created_by || ''),
        createdAt: new Date(raw.created_at as string).getTime(),
        updatedAt: new Date((raw.updated_at as string) || (raw.created_at as string)).getTime(),
      };
      set({ 
        currentRoom: room, 
        agents, 
        consensus: {
          topic: room.topic || room.name,
          agreement: room.consensusScore,
          breakdown: {},
        },
        isLoading: false,
      });
    } catch (err) {
      console.error('fetchRoom error:', err);
      set({ isLoading: false });
    }
  },

  fetchMessages: async (roomId: string) => {
    try {
      const res = await roomApi.getMessages(roomId);
      const rawList = (res.data as { list?: unknown[] }).list || [];
      const messages = rawList.map(adaptMessage);
      set({ messages });
    } catch (err) {
      console.error('fetchMessages error:', err);
    }
  },

  updatePhase: (phase, round) => set((state) => ({
    currentRoom: state.currentRoom ? {
      ...state.currentRoom,
      currentPhase: phase,
      currentRound: round,
    } : null,
  })),
}));
