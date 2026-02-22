import { test, expect } from '@playwright/test';

test.describe('authenticated dashboard', () => {
  test('shows researcher dashboard when authenticated', async ({ page }) => {
    await page.goto('/app/en/researcher/');
    // Dashboard renders either campaign list or empty welcome state
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
  });

  test('shows navigation with campaigns and logout', async ({ page }) => {
    await page.goto('/app/en/researcher/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    // On mobile, nav items are behind the hamburger menu
    const hamburger = page.locator('.mobile-nav-toggle');
    if (await hamburger.isVisible()) {
      await hamburger.click();
      const mobileNav = page.locator('.mobile-nav');
      await expect(mobileNav).toBeVisible();
      await expect(mobileNav.getByRole('link', { name: 'Campaigns' })).toBeVisible();
      await expect(mobileNav.getByRole('button', { name: 'Log out' })).toBeVisible();
    } else {
      await expect(page.getByRole('link', { name: 'Campaigns' })).toBeVisible();
      await expect(page.getByRole('button', { name: 'Log out' })).toBeVisible();
    }
  });
});
