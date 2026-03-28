// playwright.config.js
const { defineConfig } = require('@playwright/test');

module.exports = defineConfig({
  testDir: './tests',
  outputDir: './test-results',
  timeout: 30000,
  retries: 1,
  use: {
    baseURL: process.env.FRONTEND_URL || 'http://localhost:3000',
    headless: true,
    screenshot: 'only-on-failure',
    // NixOS VM 关键配置
    launchOptions: {
      args: [
        '--no-sandbox',           // NixOS VM 必须
        '--disable-gpu',          // 无显卡VM
        '--disable-web-security', // localhost跨域
      ],
      // NixOS 系统 Chromium 路径（如果不用 playwright 自带的）
      // executablePath: process.env.PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH || undefined,
    },
  },
  projects: [
    { name: 'chromium', use: { browserName: 'chromium' } },
  ],
});
