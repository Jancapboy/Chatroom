import { create } from 'zustand';
import type { User } from '../types/room';

interface AuthState {
  user: User | null;
  isLoggedIn: boolean;
  isLoading: boolean;
  
  // Actions
  login: (email: string, password: string) => Promise<boolean>;
  logout: () => void;
  setUser: (user: User) => void;
}

// Mock login for now - will connect to real API later
const MOCK_USER: User = {
  id: 'user-1',
  name: 'Observer',
  email: 'observer@chatroom.asi',
  avatar: '/avatars/user.png',
  token: 'mock-jwt-token-12345',
};

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isLoggedIn: false,
  isLoading: false,

  login: async (email: string, password: string) => {
    // TODO: Connect to real API
    // For now, accept any non-empty credentials
    if (email && password) {
      set({ user: MOCK_USER, isLoggedIn: true, isLoading: false });
      // Persist token
      localStorage.setItem('asi_token', MOCK_USER.token || '');
      return true;
    }
    return false;
  },

  logout: () => {
    localStorage.removeItem('asi_token');
    set({ user: null, isLoggedIn: false });
  },

  setUser: (user: User) => {
    set({ user, isLoggedIn: true });
  },
}));

// Check for existing token on init
const token = localStorage.getItem('asi_token');
if (token) {
  useAuthStore.setState({ user: { ...MOCK_USER, token }, isLoggedIn: true });
}
