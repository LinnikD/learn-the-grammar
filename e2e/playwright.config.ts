import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './tests',

  retries: process.env.CI ? 1 : 0,

  use: {
    baseURL: 'http://localhost',
    trace: 'on-first-retry',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
})