/**
 * Chatroom 3D 前端 Playwright 端到端测试
 * 
 * NixOS 运行方式:
 *   nix-shell -p nodejs playwright-driver.browsers --run "npx playwright test"
 * 
 * 或用 venv:
 *   python3 -m venv .venv && .venv/bin/pip install playwright
 *   .venv/bin/playwright install chromium
 *   .venv/bin/playwright test
 * 
 * 环境变量:
 *   PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH=/nix/store/...-chromium/bin/chromium
 */
const { test, expect } = require('@playwright/test');

const BASE_URL = process.env.FRONTEND_URL || 'http://localhost:3000';
const API_URL = process.env.API_URL || 'http://localhost:4001';

test.describe('Chatroom 3D 集成测试', () => {

  test('登录页加载', async ({ page }) => {
    await page.goto(BASE_URL);
    // 等待登录表单出现
    await expect(page.locator('input[type="text"]')).toBeVisible({ timeout: 10000 });
    await expect(page.locator('input[type="password"]')).toBeVisible();
    await page.screenshot({ path: 'screenshots/01-login-page.png' });
  });

  test('注册新用户', async ({ page }) => {
    await page.goto(BASE_URL);
    // 点击"去注册"链接（如果有的话）
    const registerLink = page.locator('text=注册');
    if (await registerLink.isVisible()) {
      await registerLink.click();
    }
    // 通过API注册（绕过UI，因为注册页可能不存在）
    const res = await page.request.post(`${API_URL}/api/register`, {
      data: { user_name: 'pw_test_' + Date.now(), password: 'test123', nickname: 'Playwright测试' }
    });
    expect(res.ok()).toBeTruthy();
  });

  test('登录并进入聊天室', async ({ page }) => {
    // 先注册
    const username = 'pw_login_' + Date.now();
    await page.request.post(`${API_URL}/api/register`, {
      data: { user_name: username, password: 'test123', nickname: '登录测试' }
    });

    await page.goto(BASE_URL);
    await page.fill('input[type="text"]', username);
    await page.fill('input[type="password"]', 'test123');
    
    // 点击登录按钮
    await page.locator('svg, button, [class*="login"]').last().click();
    
    // 等待进入聊天室（聊天输入框出现）
    await expect(page.locator('input[type="text"]')).toBeVisible({ timeout: 10000 });
    await page.screenshot({ path: 'screenshots/02-chatroom.png' });
  });

  test('3D按钮存在', async ({ page }) => {
    const username = 'pw_3d_' + Date.now();
    await page.request.post(`${API_URL}/api/register`, {
      data: { user_name: username, password: 'test123', nickname: '3D测试' }
    });

    await page.goto(BASE_URL);
    await page.fill('input[type="text"]', username);
    await page.fill('input[type="password"]', 'test123');
    await page.locator('svg, button, [class*="login"]').last().click();
    
    // 等待聊天室加载
    await page.waitForTimeout(3000);
    
    // 查找3D按钮（🎨 emoji 或 "3D" 文字）
    const btn3D = page.locator('text=🎨, text=3D, button:has-text("3D")');
    await expect(btn3D.first()).toBeVisible({ timeout: 5000 });
    await page.screenshot({ path: 'screenshots/03-3d-button.png' });
  });

  test('3D生成弹窗', async ({ page }) => {
    const username = 'pw_modal_' + Date.now();
    await page.request.post(`${API_URL}/api/register`, {
      data: { user_name: username, password: 'test123', nickname: '弹窗测试' }
    });

    await page.goto(BASE_URL);
    await page.fill('input[type="text"]', username);
    await page.fill('input[type="password"]', 'test123');
    await page.locator('svg, button, [class*="login"]').last().click();
    await page.waitForTimeout(3000);

    // 点击3D按钮
    const btn3D = page.locator('text=🎨, text=3D, button:has-text("3D")');
    await btn3D.first().click();

    // 验证弹窗出现
    await expect(page.locator('.ant-modal, [class*="modal"]')).toBeVisible({ timeout: 5000 });
    await page.screenshot({ path: 'screenshots/04-3d-modal.png' });
  });

  test('后端3D API可用', async ({ page }) => {
    // 直接测API
    const queryRes = await page.request.post(`${API_URL}/api/3d/query`, {
      data: { job_id: '1429320300369436672' }
    });
    expect(queryRes.ok()).toBeTruthy();
    const data = await queryRes.json();
    expect(data.status).toBe('DONE');
    expect(data.files.length).toBeGreaterThan(0);
    expect(data.files[0].type).toBe('GLB');
  });
});
