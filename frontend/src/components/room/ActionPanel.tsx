import { X, CheckCircle, ArrowRight, FileText, GitBranch } from 'lucide-react';

export interface ActionDraft {
  id: string;
  title: string;
  description: string;
}

interface ActionPanelProps {
  conclusion: string;
  actions: ActionDraft[];
  onExecute: (actionId: string) => void;
  onSkip: () => void;
  onClose: () => void;
}

export function ActionPanel({ conclusion, actions, onExecute, onSkip, onClose }: ActionPanelProps) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/60" onClick={onClose} />

      {/* Panel */}
      <div className="relative w-full max-w-lg bg-bg-secondary border border-border-subtle rounded-xl shadow-2xl overflow-hidden flex flex-col max-h-[85vh]">
        {/* Header */}
        <div className="px-5 py-4 border-b border-border-subtle flex items-center justify-between shrink-0">
          <div className="flex items-center gap-2">
            <CheckCircle size={18} className="text-accent-cyan" />
            <h2 className="text-base font-semibold text-text-primary">推演结果</h2>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 hover:bg-bg-elevated rounded-md transition-colors"
          >
            <X size={16} className="text-text-muted" />
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-5 space-y-5">
          {/* Conclusion */}
          <div>
            <h3 className="text-xs font-semibold text-text-muted uppercase tracking-wider mb-2">
              达成的结论
            </h3>
            <div className="bg-bg-elevated rounded-lg p-4 border border-border-subtle">
              <p className="text-sm text-text-primary leading-relaxed">{conclusion}</p>
            </div>
          </div>

          {/* Actions */}
          {actions.length > 0 && (
            <div>
              <h3 className="text-xs font-semibold text-text-muted uppercase tracking-wider mb-2">
                建议的操作
              </h3>
              <div className="space-y-2">
                {actions.map((action) => (
                  <button
                    key={action.id}
                    onClick={() => onExecute(action.id)}
                    className="w-full flex items-center gap-3 p-3 bg-bg-elevated rounded-lg border border-border-subtle hover:border-accent-cyan/40 hover:bg-bg-card transition-all text-left active:scale-[0.98]"
                  >
                    <div className="w-8 h-8 rounded-full bg-accent-cyan/10 flex items-center justify-center shrink-0">
                      {action.id === '2' ? (
                        <GitBranch size={14} className="text-accent-cyan" />
                      ) : action.id === '3' ? (
                        <FileText size={14} className="text-accent-cyan" />
                      ) : (
                        <CheckCircle size={14} className="text-accent-cyan" />
                      )}
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="text-sm font-medium text-text-primary">{action.title}</div>
                      <div className="text-xs text-text-muted mt-0.5">{action.description}</div>
                    </div>
                    <ArrowRight size={14} className="text-text-muted shrink-0" />
                  </button>
                ))}
              </div>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="px-5 py-3 border-t border-border-subtle flex justify-end shrink-0">
          <button
            onClick={onSkip}
            className="px-4 py-2 text-sm text-text-muted hover:text-text-primary transition-colors"
          >
            关闭
          </button>
        </div>
      </div>
    </div>
  );
}
