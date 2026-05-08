import { useState, useEffect } from 'react';
import { ArrowLeft, Sparkles, Users, Clock, Check } from 'lucide-react';
import { roomApi } from '../lib/api';
import type { AgentRole } from '../types/room';

interface Template {
  id: string;
  name: string;
  description: string;
  defaultAgents: string[];
  maxRounds: number;
}

interface AgentTemplate {
  id: string;
  name: string;
  role: string;
  color: string;
  expertise: string[];
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

interface CreateRoomProps {
  onNavigate: (page: string) => void;
}

export function CreateRoom({ onNavigate }: CreateRoomProps) {
  const [name, setName] = useState('');
  const [topic, setTopic] = useState('');
  const [description, setDescription] = useState('');
  const [selectedTemplate, setSelectedTemplate] = useState('default');
  const [selectedAgents, setSelectedAgents] = useState<AgentRole[]>(['architect', 'risk_officer', 'strategist']);
  const [maxRounds, setMaxRounds] = useState(8);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [templates, setTemplates] = useState<Template[]>([]);
  const [agentTemplates, setAgentTemplates] = useState<AgentTemplate[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const res = await roomApi.getAgentTemplates();
        const tpls = (res.data as { templates?: unknown[] }).templates || [];
        const agents = tpls.map((t: unknown) => {
          const raw = t as Record<string, unknown>;
          let expertise: string[] = [];
          try { expertise = JSON.parse(raw.expertise as string); } catch { /* ignore */ }
          return {
            id: raw.id as string,
            name: raw.name as string,
            role: raw.role as string,
            color: getAgentColor(raw.role as string),
            expertise,
          };
        });
        setAgentTemplates(agents);
        setTemplates([
          {
            id: 'default',
            name: '默认推演',
            description: '通用多Agent推演模板',
            defaultAgents: ['architect', 'risk_officer', 'strategist'],
            maxRounds: 8,
          },
          {
            id: 'product_review',
            name: '产品评审',
            description: '技术方案评审与决策',
            defaultAgents: ['architect', 'analyst', 'executor'],
            maxRounds: 6,
          },
          {
            id: 'crisis',
            name: '危机应对',
            description: '突发事件应急推演',
            defaultAgents: ['risk_officer', 'strategist', 'executor'],
            maxRounds: 4,
          },
        ]);
        if (agents.length > 0) {
          setSelectedAgents([agents[0].role as AgentRole, agents[1]?.role as AgentRole].filter(Boolean));
        }
      } catch (err) {
        console.error('load templates error:', err);
      } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  const toggleAgent = (role: AgentRole) => {
    setSelectedAgents((prev) =>
      prev.includes(role) ? prev.filter((r) => r !== role) : [...prev, role]
    );
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || selectedAgents.length < 2) return;

    setIsSubmitting(true);
    try {
      const agentIds = agentTemplates
        .filter((a) => selectedAgents.includes(a.role as AgentRole))
        .map((a) => a.id);
      await roomApi.createRoom({
        name,
        topic,
        description,
        max_rounds: maxRounds,
        agent_ids: agentIds,
      });
      onNavigate('lobby');
    } catch (err) {
      console.error('create room error:', err);
      alert('创建房间失败: ' + String(err));
    } finally {
      setIsSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-bg-primary flex items-center justify-center">
        <p className="text-text-muted">加载中...</p>
      </div>
    );
  }

  const currentTemplate = templates.find((t) => t.id === selectedTemplate) || templates[0];
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  void currentTemplate;

  return (
    <div className="min-h-screen bg-bg-primary">
      {/* Header */}
      <div className="bg-bg-secondary border-b border-border-subtle px-4 py-4">
        <div className="max-w-3xl mx-auto flex items-center gap-3">
          <button
            onClick={() => onNavigate('lobby')}
            className="p-2 hover:bg-bg-elevated rounded-md transition-colors"
          >
            <ArrowLeft size={18} className="text-text-secondary" />
          </button>
          <h1 className="text-lg font-semibold text-text-primary">创建推演房间</h1>
        </div>
      </div>

      <div className="max-w-3xl mx-auto px-4 py-8">
        <form onSubmit={handleSubmit} className="space-y-8">
          {/* Basic info */}
          <section className="space-y-4">
            <h2 className="text-sm font-semibold text-text-secondary uppercase tracking-wider flex items-center gap-2">
              <Sparkles size={16} className="text-accent-cyan" />
              基本信息
            </h2>

            <div className="space-y-3">
              <div>
                <label className="block text-sm text-text-secondary mb-1.5">房间名称</label>
                <input
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="例如：火星殖民计划推演"
                  required
                  className="w-full bg-bg-secondary border border-border-subtle rounded-md px-3 py-2.5 text-sm text-text-primary placeholder:text-text-muted focus:outline-none focus:border-accent-cyan transition-colors"
                />
              </div>

              <div>
                <label className="block text-sm text-text-secondary mb-1.5">推演主题</label>
                <input
                  type="text"
                  value={topic}
                  onChange={(e) => setTopic(e.target.value)}
                  placeholder="例如：火星殖民计划可行性分析"
                  className="w-full bg-bg-secondary border border-border-subtle rounded-md px-3 py-2.5 text-sm text-text-primary placeholder:text-text-muted focus:outline-none focus:border-accent-cyan transition-colors"
                />
              </div>

              <div>
                <label className="block text-sm text-text-secondary mb-1.5">描述</label>
                <textarea
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  placeholder="描述推演的目标和背景..."
                  rows={3}
                  className="w-full bg-bg-secondary border border-border-subtle rounded-md px-3 py-2.5 text-sm text-text-primary placeholder:text-text-muted focus:outline-none focus:border-accent-cyan transition-colors resize-none"
                />
              </div>
            </div>
          </section>

          {/* Template selection */}
          <section className="space-y-4">
            <h2 className="text-sm font-semibold text-text-secondary uppercase tracking-wider flex items-center gap-2">
              <Sparkles size={16} className="text-accent-purple" />
              选择模板
            </h2>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              {templates.map((tpl) => (
                <button
                  key={tpl.id}
                  type="button"
                  onClick={() => {
                    setSelectedTemplate(tpl.id);
                    setMaxRounds(tpl.maxRounds);
                    setSelectedAgents(tpl.defaultAgents as AgentRole[]);
                  }}
                  className={`p-4 rounded-lg border text-left transition-all ${
                    selectedTemplate === tpl.id
                      ? 'border-accent-cyan bg-accent-cyan/5'
                      : 'border-border-subtle bg-bg-secondary hover:border-border-glow'
                  }`}
                >
                  <div className="font-medium text-text-primary text-sm">{tpl.name}</div>
                  <p className="text-xs text-text-muted mt-1">{tpl.description}</p>
                  <div className="flex items-center gap-3 mt-2 text-xs text-text-muted">
                    <span className="flex items-center gap-1">
                      <Users size={12} />
                      {tpl.defaultAgents.length} 个Agent
                    </span>
                    <span className="flex items-center gap-1">
                      <Clock size={12} />
                      {tpl.maxRounds} 轮
                    </span>
                  </div>
                </button>
              ))}
            </div>
          </section>

          {/* Agent configuration */}
          <section className="space-y-4">
            <h2 className="text-sm font-semibold text-text-secondary uppercase tracking-wider flex items-center gap-2">
              <Users size={16} className="text-agent-strategist" />
              Agent 配置
            </h2>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
              {agentTemplates.map((agent) => {
                const isSelected = selectedAgents.includes(agent.role as AgentRole);
                return (
                  <button
                    key={agent.id}
                    type="button"
                    onClick={() => toggleAgent(agent.role as AgentRole)}
                    className={`flex items-center gap-3 p-3 rounded-lg border transition-all ${
                      isSelected
                        ? 'border-accent-cyan bg-accent-cyan/5'
                        : 'border-border-subtle bg-bg-secondary opacity-60 hover:opacity-100'
                    }`}
                  >
                    <div
                      className="w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0"
                      style={{ backgroundColor: agent.color }}
                    >
                      {isSelected && <Check size={14} />}
                    </div>
                    <div className="text-left">
                      <div className="text-sm font-medium text-text-primary">{agent.name}</div>
                      <div className="text-xs text-text-muted">{agent.expertise.join(' · ')}</div>
                    </div>
                  </button>
                );
              })}
            </div>

            {selectedAgents.length < 2 && (
              <p className="text-xs text-accent-danger">至少需要选择 2 个Agent</p>
            )}
          </section>

          {/* Max rounds */}
          <section className="space-y-3">
            <h2 className="text-sm font-semibold text-text-secondary uppercase tracking-wider flex items-center gap-2">
              <Clock size={16} className="text-agent-executor" />
              推演设置
            </h2>
            <div>
              <label className="block text-sm text-text-secondary mb-1.5">
                最大轮次: <span className="text-accent-cyan font-medium">{maxRounds}</span>
              </label>
              <input
                type="range"
                min={3}
                max={20}
                value={maxRounds}
                onChange={(e) => setMaxRounds(Number(e.target.value))}
                className="w-full accent-accent-cyan"
              />
              <div className="flex justify-between text-xs text-text-muted mt-1">
                <span>3轮</span>
                <span>20轮</span>
              </div>
            </div>
          </section>

          {/* Submit */}
          <div className="pt-4">
            <button
              type="submit"
              disabled={isSubmitting || !name.trim() || selectedAgents.length < 2}
              className="w-full py-3 bg-accent-cyan text-bg-primary rounded-md font-semibold text-sm hover:bg-accent-cyan-dim transition-colors disabled:opacity-50 disabled:cursor-not-allowed active:scale-[0.98]"
            >
              {isSubmitting ? '创建中...' : '创建房间'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
