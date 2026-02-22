import { defineConfig, devices } from '@playwright/test';

const STORAGE_STATE = './tests/.auth/user.json';

export default defineConfig({
  testDir: './tests',
  timeout: 60_000,
  fullyParallel: false,
  workers: 1,
  reporter: 'list',
  use: {
    baseURL: 'http://caddy:9999',
    headless: true,
    screenshot: 'only-on-failure',
  },
  projects: [
    {
      name: 'setup',
      testMatch: /.*\.setup\.ts/,
    },
    {
      name: 'chromium',
      testIgnore: [/.*\.noauth\.spec\.ts/, /.*\.mobile\.spec\.ts/],
      use: {
        browserName: 'chromium',
        storageState: STORAGE_STATE,
      },
      dependencies: ['setup'],
    },
    {
      name: 'mobile-chrome',
      testIgnore: /.*\.noauth\.spec\.ts/,
      use: {
        ...devices['Pixel 5'],
        storageState: STORAGE_STATE,
      },
      dependencies: ['setup'],
    },
    {
      name: 'no-auth',
      testMatch: /.*\.noauth\.spec\.ts/,
      use: { browserName: 'chromium' },
    },
  ],
});
