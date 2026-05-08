import { create } from 'zustand';
import type { Room, RoomStatus, Agent, Phase } from '../types/room';
import { roomApi } from '../lib/api';

interface LobbyState {
  rooms: Room[];
  filter: RoomStatus | 'all';
  isLoading: boolean;
  
  // Actions
  setFilter: (filter: RoomStatus | 'all') => void;
  fetchRooms: () => Promise<void>;
  addRoom: (room: Room) => void;
  updateRoom: (room: Room) => void;
  getFilteredRooms: () => Room[];
}

function adaptRoom(raw: unknown): Room {
  const r = raw as Record<string, unknown>;
  const agents = (r.agents as unknown[] || []).map((a: unknown) => {
    const ag = a as Record<string, unknown>;
    let expertise: string[] = [];
    try {
      expertise = JSON.parse(ag.expertise as string);
    } catch { /* ignore */ }
    return {
      id: ag.id as string,
      name: ag.name as string,
      role: ag.role as string as any,
      avatar: (ag.avatar as string) || '',
      personality: ag.personality as string,
      expertise,
      model: (ag.model as string) || 'deepseek-chat',
      energy: (ag.energy as number) || 100,
      confidence: (ag.confidence as number) || 50,
      stance: (ag.stance as string) as any,
      isActive: ag.is_active !== false,
      color: getAgentColor(ag.role as string),
    } as unknown as Agent;
  });

  return {
    id: r.id as string,
    name: r.name as string,
    topic: (r.topic as string) || '',
    description: (r.description as string) || '',
    status: r.status as RoomStatus,
    templateId: (r.template_id as string) || undefined,
    agents,
    currentPhase: (r.current_phase as string) as Phase || 'info_gathering',
    currentRound: (r.current_round as number) || 1,
    maxRounds: (r.max_rounds as number) || 10,
    consensusScore: (r.consensus_score as number) || 0,
    forkedFrom: (r.forked_from as string) || undefined,
    createdBy: String(r.created_by || ''),
    createdAt: new Date(r.created_at as string).getTime(),
    updatedAt: new Date((r.updated_at as string) || (r.created_at as string)).getTime(),
    messageCount: (r.message_count as number) || 0,
  };
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

export const useLobbyStore = create<LobbyState>((set, get) => ({
  rooms: [],
  filter: 'all',
  isLoading: false,

  setFilter: (filter) => set({ filter }),

  fetchRooms: async () => {
    set({ isLoading: true });
    try {
      const res = await roomApi.getRooms();
      const rawList = (res.data as { list?: unknown[] }).list || [];
      const rooms = rawList.map(adaptRoom);
      set({ rooms, isLoading: false });
    } catch (err) {
      console.error('fetchRooms error:', err);
      set({ isLoading: false });
    }
  },

  addRoom: (room) => {
    set((state) => ({ rooms: [room, ...state.rooms] }));
  },

  updateRoom: (room) => {
    set((state) => ({
      rooms: state.rooms.map((r) => (r.id === room.id ? room : r)),
    }));
  },

  getFilteredRooms: () => {
    const { rooms, filter } = get();
    if (filter === 'all') return rooms;
    return rooms.filter((r) => r.status === filter);
  },
}));
