import { create } from 'zustand';
import type { Room, Agent, Message, Phase, ConsensusState } from '../types/room';
import { mockRoomDetail, mockMessages } from '../lib/mockData';

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
  updatePhase: (phase: Phase, round: number) => void;
}

export const useRoomStore = create<RoomState>((set) => ({
  currentRoom: mockRoomDetail,
  messages: mockMessages,
  agents: mockRoomDetail?.agents || [],
  consensus: {
    topic: '火星殖民计划可行性',
    agreement: 65,
    breakdown: {
      'agent-1': 85, // architect - support
      'agent-2': 30, // risk - oppose
      'agent-3': 70, // strategist - support
      'agent-4': 60, // analyst - neutral
      'agent-5': 55, // executor - neutral
    },
  },
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
      // TODO: Connect to real API
      // const res = await api.get(`/rooms/${roomId}`);
      // set({ currentRoom: res.data.data });
      
      // Use mock data for now
      const room = { ...mockRoomDetail, id: roomId };
      set({ 
        currentRoom: room, 
        agents: room.agents, 
        messages: mockMessages,
        isLoading: false,
      });
    } catch {
      set({ isLoading: false });
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
