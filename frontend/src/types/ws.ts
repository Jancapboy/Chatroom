import type { Phase } from '../types/room';

export interface WSClientMessage {
  type: 'user_message' | 'command';
  payload: {
    content?: string;
    command?: 'pause' | 'resume' | 'next_phase' | 'fork';
  };
}

export interface WSServerMessage {
  type: 'message' | 'phase_change' | 'agent_state' | 'consensus_update' | 'system';
  payload: unknown;
}

// Backend actual payload shapes (snake_case)
export interface AgentMessagePayload {
  id: string;
  sender_id: string;
  sender_type: 'agent' | 'user' | 'system';
  sender_name: string;
  sender_avatar?: string;
  content: string;
  phase: string;
  round: number;
  confidence?: number;
  stance?: string;
}

export interface PhaseChangePayload {
  phase: Phase;
  round: number;
  phase_name?: string;
  description?: string;
}

export interface AgentStatePayload {
  agent_id: string;
  name: string;
  role: string;
  energy: number;
  confidence: number;
  stance: string;
}

export interface ConsensusUpdatePayload {
  topic: string;
  agreement: number;
  breakdown?: Record<string, unknown>;
}

export interface SystemPayload {
  event: 'room_started' | 'room_completed' | 'agent_joined' | 'room_paused' | 'room_resumed';
  message: string;
}
