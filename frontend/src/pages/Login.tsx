import { useState } from 'react';
import { LogIn, Eye, EyeOff, Flame } from 'lucide-react';
import { useAuthStore } from '../stores/useAuthStore';

interface LoginProps {
  onNavigate: (page: string) => void;
}

export function Login({ onNavigate }: LoginProps) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const login = useAuthStore((state) => state.login);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    const success = await login(email, password);
    setIsLoading(false);

    if (success) {
      onNavigate('lobby');
    } else {
      setError('登录失败，请检查邮箱和密码');
    }
  };

  return (
    <div className="min-h-screen bg-bg-primary flex items-center justify-center px-4">
      <div className="w-full max-w-md animate-fadeIn">
        {/* Logo */}
        <div className="text-center mb-8">
          <div className="w-16 h-16 mx-auto rounded-xl bg-accent-cyan/10 flex items-center justify-center mb-4">
            <Flame size={32} className="text-accent-cyan" />
          </div>
          <h1 className="text-2xl font-bold text-text-primary">ASI Chatroom</h1>
          <p className="text-sm text-text-muted mt-2">多智能体推演平台</p>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm text-text-secondary mb-1.5">邮箱</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="your@email.com"
              required
              className="w-full bg-bg-secondary border border-border-subtle rounded-md px-3 py-2.5 text-sm text-text-primary placeholder:text-text-muted focus:outline-none focus:border-accent-cyan transition-colors"
            />
          </div>

          <div>
            <label className="block text-sm text-text-secondary mb-1.5">密码</label>
            <div className="relative">
              <input
                type={showPassword ? 'text' : 'password'}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="输入密码"
                required
                className="w-full bg-bg-secondary border border-border-subtle rounded-md px-3 py-2.5 text-sm text-text-primary placeholder:text-text-muted focus:outline-none focus:border-accent-cyan transition-colors pr-10"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-secondary"
              >
                {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
              </button>
            </div>
          </div>

          {error && (
            <div className="text-sm text-accent-danger bg-accent-danger/10 px-3 py-2 rounded-md">
              {error}
            </div>
          )}

          <button
            type="submit"
            disabled={isLoading}
            className="w-full py-3 bg-accent-cyan text-bg-primary rounded-md font-semibold text-sm hover:bg-accent-cyan-dim transition-colors disabled:opacity-50 flex items-center justify-center gap-2 active:scale-[0.98]"
          >
            <LogIn size={16} />
            {isLoading ? '登录中...' : '登录'}
          </button>
        </form>

        {/* Guest access */}
        <div className="mt-6 text-center">
          <button
            onClick={() => onNavigate('lobby')}
            className="text-sm text-text-muted hover:text-accent-cyan transition-colors"
          >
            以游客身份浏览大厅 →
          </button>
        </div>
      </div>
    </div>
  );
}
