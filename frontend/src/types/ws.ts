import type { Message, Phase } from '../types/room';

export interface WSClientMessage {
  type: 'user_message' | 'command';
  payload: {
    content?: string;
    command?: 'pause' | 'resume' | 'next_phase' | 'fork';
  };
}

export interface WSServerMessage {
  type: 'message' | 'phase_change' | 'agent_state' | 'consensus_update' | 'system';
  payload: WSMessagePayload;
}

export type WSMessagePayload =
  | MessagePayload
  | PhaseChangePayload
  | AgentStatePayload
  | ConsensusUpdatePayload
  | SystemPayload;

export interface MessagePayload {
  message: Message;
}

export interface PhaseChangePayload {
  phase: Phase;
  round: number;
}

export interface AgentStatePayload {
  agentId: string;
  energy: number;
  confidence: number;
  stance: string;
}

export interface ConsensusUpdatePayload {
  topic: string;
  agreement: number;
}

export interface SystemPayload {
  event: 'room_started' | 'room_completed' | 'agent_joined' | 'room_paused' | 'room_resumed';
  data?: Record<string, unknown>;
}
