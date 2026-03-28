# NixOS 前端测试指南

## 环境陷阱与解决方案

### 1. Chromium 沙盒问题

NixOS 的 Chromium 沙盒默认很严，在 VM 里必须加 `--no-sandbox`：

```bash
# ❌ 会报错
chromium --headless http://localhost:3000

# ✅ 正确方式
chromium --headless=new --no-sandbox --disable-gpu http://localhost:3000
```

> ⚠️ `--no-sandbox` 仅限开发/测试环境使用

### 2. 截图全黑

无头 VM 缺显卡，解决方案：

```bash
# 方案A: 用新版 headless 模式
chromium --headless=new --no-sandbox --disable-gpu --screenshot=test.png http://localhost:3000

# 方案B: 用 xvfb-run 虚拟显示
nix-shell -p xvfb-run --run "xvfb-run chromium --no-sandbox --screenshot=test.png http://localhost:3000"
```

### 3. 中文字体缺失（显示方块□）

```nix
# configuration.nix
fonts.packages = with pkgs; [
  noto-fonts-cjk-sans
  noto-fonts-cjk-serif
];
```

或临时方案：
```bash
nix-shell -p noto-fonts-cjk-sans
```

### 4. 跨域问题

localhost 前端调 localhost 后端 API 时：

```bash
# Chromium 加 --disable-web-security
chromium --headless=new --no-sandbox --disable-web-security ...

# 或后端加 CORS 中间件（已配置）
```

## 运行 Playwright 测试

### 方式A: nix-shell（推荐）

```bash
nix-shell -p nodejs playwright-driver.browsers -p noto-fonts-cjk-sans --run "
  export PLAYWRIGHT_BROWSERS_PATH=${playwright-driver.browsers}
  npx playwright test
"
```

### 方式B: Python venv

```bash
python3 -m venv .venv
.venv/bin/pip install playwright
.venv/bin/playwright install chromium
FRONTEND_URL=http://localhost:3000 API_URL=http://localhost:4001 .venv/bin/python -m pytest
```

### 方式C: 直接用系统 Chromium（快速验证）

```bash
# 截图登录页
chromium --headless=new --no-sandbox --disable-gpu \
  --screenshot=login.png --window-size=1280,720 \
  --virtual-time-budget=8000 \
  http://localhost:3000

# 跑联调测试页
chromium --headless=new --no-sandbox --disable-gpu --disable-web-security \
  --screenshot=e2e.png --window-size=700,400 \
  --virtual-time-budget=10000 \
  http://localhost:3000/e2e_test.html
```

## 测试覆盖

| 测试 | 文件 | 方式 |
|------|------|------|
| 后端API + Git状态 | `scripts/e2e_test.py` | `python3 scripts/e2e_test.py` |
| 前端UI交互 | `tests/e2e.spec.js` | `npx playwright test` |
| 快速截图验证 | - | Chromium 命令行 |

## 已知限制

- WebGL/Three.js 在纯软件渲染下可能无法正常渲染3D模型
- Playwright 的 nix 包编译较慢，建议用 venv + pip 方式
- `--virtual-time-budget` 控制JS执行时间，复杂页面需要设大一些（>5000ms）
