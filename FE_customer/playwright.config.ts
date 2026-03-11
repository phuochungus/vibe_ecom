import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  outputDir: '../tests/.playwright-customer',
  reporter: [
    ['list'],
    ['../tests/playwright-evidence-reporter.mjs', { appName: 'fe-customer', evidenceDir: '../tests' }],
  ],
  use: {
    baseURL: 'http://127.0.0.1:3001',
    trace: 'on-first-retry',
    video: 'on',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
})
