import { create } from 'zustand';
import type { Room, RoomStatus } from '../types/room';
import { mockRooms } from '../lib/mockData';

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

export const useLobbyStore = create<LobbyState>((set, get) => ({
  rooms: mockRooms,
  filter: 'all',
  isLoading: false,

  setFilter: (filter) => set({ filter }),

  fetchRooms: async () => {
    set({ isLoading: true });
    try {
      // TODO: Connect to real API
      // const res = await api.get('/rooms');
      // set({ rooms: res.data.data });
      
      // For now, use mock data
      set({ rooms: mockRooms, isLoading: false });
    } catch {
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
