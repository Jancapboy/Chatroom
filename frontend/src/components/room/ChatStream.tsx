import { motion, AnimatePresence } from 'framer-motion';
import { Bot, User, Info } from 'lucide-react';
import type { Message, AgentRole } from '../../types/room';

const roleColors: Record<AgentRole, string> = {
  architect: '#3b82f6',
  risk_officer: '#ef4444',
  strategist: '#8b5cf6',
  analyst: '#10b981',
  executor: '#f59e0b',
};

const phaseLabels: Record<string, string> = {
  info_gathering: '信息收集',
  opinion_expression: '观点表达',
  debate: '辩论',
  consensus: '共识',
  decision: '决策',
  summary: '总结',
};

interface ChatStreamProps {
  messages: Message[];
}

export function ChatStream({ messages }: ChatStreamProps) {
  // Group messages by round+phase
  const groups: { key: string; round: number; phase: string; items: Message[] }[] = [];
  let currentGroup: { key: string; round: number; phase: string; items: Message[] } | null = null;

  for (const msg of messages) {
    const key = `${msg.round}-${msg.phase}`;
    if (!currentGroup || currentGroup.key !== key) {
      if (currentGroup) groups.push(currentGroup);
      currentGroup = { key, round: msg.round, phase: msg.phase, items: [msg] };
    } else {
      currentGroup.items.push(msg);
    }
  }
  if (currentGroup) groups.push(currentGroup);

  return (
    <div className="h-full flex flex-col bg-bg-primary">
      {/* Messages area */}
      <div className="flex-1 overflow-y-auto p-4 space-y-1">
        <AnimatePresence>
          {groups.map((group) => (
            <div key={group.key} className="mb-4">
              {/* Phase/Round divider */}
              <div className="flex items-center gap-3 mb-3">
                <div className="h-px flex-1 bg-border-subtle" />
                <span className="text-xs text-text-muted font-medium">
                  第 {group.round} 轮 · {phaseLabels[group.phase]}
                </span>
                <div className="h-px flex-1 bg-border-subtle" />
              </div>

              {group.items.map((msg) => (
                <MessageBubble key={msg.id} message={msg} />
              ))}
            </div>
          ))}
        </AnimatePresence>
      </div>
    </div>
  );
}

function MessageBubble({ message }: { message: Message }) {
  const isAgent = message.senderType === 'agent';
  const isUser = message.senderType === 'user';
  const isSystem = message.senderType === 'system';

  if (isSystem || message.type === 'phase_change' || message.type === 'consensus') {
    return (
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        className="flex justify-center my-2"
      >
        <div className="flex items-center gap-2 px-4 py-2 bg-bg-elevated rounded-full border border-border-subtle">
          <Info size={14} className="text-text-muted" />
          <span className="text-xs text-text-muted">{message.content}</span>
        </div>
      </motion.div>
    );
  }

  const color = isAgent && message.senderRole
    ? roleColors[message.senderRole]
    : isUser
      ? '#00d4aa'
      : '#8888a0';

  return (
    <motion.div
      initial={{ opacity: 0, y: 12, scale: 0.98 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      transition={{ duration: 0.3, ease: 'easeOut' }}
      className={`flex gap-3 mb-3 ${isUser ? 'flex-row-reverse' : 'flex-row'}`}
    >
      {/* Avatar */}
      <div
        className="w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0"
        style={{ backgroundColor: color }}
      >
        {isAgent ? <Bot size={16} /> : isUser ? <User size={16} /> : '?'}
      </div>

      {/* Content */}
      <div className={`max-w-[75%] ${isUser ? 'items-end' : 'items-start'} flex flex-col`}>
        <div className="flex items-center gap-2 mb-1">
          <span className="text-xs text-text-muted font-medium">
            {message.senderName}
          </span>
          <span className="text-[10px] text-text-muted">
            {new Date(message.timestamp).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })}
          </span>
        </div>

        <div
          className="px-3 py-2.5 rounded-lg text-sm leading-relaxed relative"
          style={{
            backgroundColor: isUser ? 'rgba(0, 212, 170, 0.1)' : 'var(--color-bg-elevated)',
            borderLeft: isAgent ? `3px solid ${color}` : isUser ? '3px solid #00d4aa' : undefined,
          }}
        >
          <p className="text-text-primary">{message.content}</p>

          {/* Metadata badges */}
          {message.metadata && (
            <div className="flex gap-2 mt-2">
              {message.metadata.confidence !== undefined && (
                <span className="text-[10px] px-1.5 py-0.5 rounded bg-bg-primary text-text-muted">
                  置信 {message.metadata.confidence}%
                </span>
              )}
              {message.metadata.stance && (
                <span className="text-[10px] px-1.5 py-0.5 rounded bg-bg-primary text-text-muted">
                  {message.metadata.stance === 'support' ? '支持' : message.metadata.stance === 'oppose' ? '反对' : '中立'}
                </span>
              )}
            </div>
          )}
        </div>
      </div>
    </motion.div>
  );
}
