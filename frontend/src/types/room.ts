export type Phase =
  | 'info_gathering'
  | 'opinion_expression'
  | 'debate'
  | 'consensus'
  | 'decision'
  | 'summary';

export type RoomStatus = 'preparing' | 'running' | 'paused' | 'completed' | 'archived';

export type AgentRole = 'architect' | 'risk_officer' | 'strategist' | 'analyst' | 'executor';

export type Stance = 'support' | 'oppose' | 'neutral';

export interface Agent {
  id: string;
  name: string;
  role: AgentRole;
  avatar: string;
  personality: string;
  expertise: string[];
  model: string;
  energy: number;        // 0-100
  confidence: number;    // 0-100
  stance: Stance;
  isActive: boolean;
  color: string;         // role color hex
}

export interface Message {
  id: string;
  roomId: string;
  senderId: string;
  senderType: 'agent' | 'user' | 'system';
  senderName: string;
  senderAvatar?: string;
  senderRole?: AgentRole;
  content: string;
  type: 'text' | 'decision' | 'consensus' | 'phase_change' | 'fork_notice' | 'system';
  phase: Phase;
  round: number;
  timestamp: number;
  metadata?: {
    confidence?: number;
    stance?: string;
    votes?: number;
  };
}

export interface ConsensusState {
  topic: string;
  agreement: number;      // 0-100
  breakdown: Record<string, number>; // per-agent stance
}

export interface Room {
  id: string;
  name: string;
  topic: string;
  description: string;
  status: RoomStatus;
  templateId?: string;
  agents: Agent[];
  currentPhase: Phase;
  currentRound: number;
  maxRounds: number;
  consensusScore: number;
  forkedFrom?: string;
  createdBy: string;
  createdAt: number;
  updatedAt: number;
  messageCount?: number;
}

export interface RoomTemplate {
  id: string;
  name: string;
  description: string;
  topic: string;
  defaultAgents: AgentRole[];
  maxRounds: number;
}

export interface User {
  id: string;
  name: string;
  email: string;
  avatar?: string;
  token?: string;
}
