import axios from 'axios';
import type { ApiResponse } from '../types/api';

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1';

export const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor - add JWT token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('asi_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor - handle errors & unwrap
api.interceptors.response.use(
  (response) => {
    // 后端统一返回 { code, msg, data }，自动解包 data
    if (response.data && typeof response.data === 'object' && 'data' in response.data) {
      return { ...response, data: response.data.data };
    }
    return response;
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('asi_token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Room API
export const roomApi = {
  getRooms: () => api.get<ApiResponse<unknown[]>>('/rooms'),
  getRoom: (id: string) => api.get<ApiResponse<unknown>>(`/rooms/${id}`),
  createRoom: (data: Record<string, unknown>) => api.post<ApiResponse<unknown>>('/rooms', data),
  startRoom: (id: string) => api.post<ApiResponse<unknown>>(`/rooms/${id}/start`),
  pauseRoom: (id: string) => api.post<ApiResponse<unknown>>(`/rooms/${id}/pause`),
  forkRoom: (id: string, data?: Record<string, unknown>) => 
    api.post<ApiResponse<unknown>>(`/rooms/${id}/fork`, data),
  deleteRoom: (id: string) => api.delete<ApiResponse<unknown>>(`/rooms/${id}`),
  getMessages: (id: string, params?: { round?: number; phase?: string }) => 
    api.get<ApiResponse<unknown[]>>(`/rooms/${id}/messages`, { params }),
  getAgentTemplates: () => api.get<ApiResponse<unknown[]>>('/agents/templates'),
  // 快照/回溯 API
  getSnapshots: (id: string) => api.get<ApiResponse<unknown>>(`/rooms/${id}/snapshots`),
  rollbackRoom: (id: string, snapshotId: string, reason?: string) => 
    api.post<ApiResponse<unknown>>(`/rooms/${id}/rollback`, { snapshot_id: snapshotId, reason }),
  forkFromSnapshot: (id: string, snapshotId: string, reason?: string) => 
    api.post<ApiResponse<unknown>>(`/rooms/${id}/fork`, { snapshot_id: snapshotId, reason }),
  createSnapshot: (id: string, reason?: string) => 
    api.post<ApiResponse<unknown>>(`/rooms/${id}/snapshots`, { reason }),
};

export default api;
