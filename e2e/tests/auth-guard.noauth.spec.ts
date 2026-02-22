import { test, expect } from '@playwright/test';

test.describe('auth guard', () => {
  test('redirects unauthenticated users from researcher dashboard to login', async ({ page }) => {
    await page.goto('/app/en/researcher/');
    await expect(page).toHaveURL(/\/en\/login/, { timeout: 10_000 });
  });
});
