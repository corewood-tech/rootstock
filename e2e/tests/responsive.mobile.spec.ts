import { test, expect } from '@playwright/test';

test.describe('responsive layout', () => {
  test('shows mobile nav toggle on small viewport', async ({ page, browserName }) => {
    // This test runs in the mobile-chrome project (Pixel 5 viewport)
    await page.goto('/app/en/researcher/');
    await expect(page.getByLabel('Open navigation')).toBeVisible({ timeout: 10_000 });
  });

  test('opens and closes mobile navigation drawer', async ({ page }) => {
    await page.goto('/app/en/researcher/');

    // Open drawer
    await page.getByLabel('Open navigation').click();
    const mobileNav = page.locator('.mobile-nav');
    await expect(page.getByLabel('Close navigation')).toBeVisible();
    await expect(mobileNav.getByRole('link', { name: 'Campaigns' })).toBeVisible();

    // Close drawer
    await page.getByLabel('Close navigation').click();
    await expect(page.getByLabel('Close navigation')).not.toBeVisible();
  });

  test('campaign wizard fits on mobile viewport', async ({ page }) => {
    await page.goto('/app/en/researcher/campaigns/new');
    const steps = page.locator('.wizard__steps');
    await expect(steps.getByText('Basics')).toBeVisible({ timeout: 10_000 });

    // Form fields should be visible and usable
    await expect(page.getByLabel('Start date')).toBeVisible();
    await expect(page.getByLabel('End date')).toBeVisible();
  });
});
