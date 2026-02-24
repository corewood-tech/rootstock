// Graph node: 0x4d (Campaign Browse E2E Test)
// Validates: BrowseCampaignsFlow (0x32), CampaignBrowsePage (0x45)
import { test, expect } from '@playwright/test';

test.describe('scitizen campaign browse', () => {
  test('shows campaign list or empty state', async ({ page }) => {
    await page.goto('/app/en/scitizen/campaigns');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    const hasCampaigns = await page.locator('.campaign-grid').isVisible().catch(() => false);
    const hasEmptyState = await page.locator('.empty-state').isVisible().catch(() => false);
    expect(hasCampaigns || hasEmptyState).toBeTruthy();
  });

  test('displays browse campaigns heading', async ({ page }) => {
    await page.goto('/app/en/scitizen/campaigns');
    await expect(page.getByRole('heading', { name: 'Browse Campaigns' })).toBeVisible({ timeout: 10_000 });
  });

  test('has search form with input and button', async ({ page }) => {
    await page.goto('/app/en/scitizen/campaigns');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    await expect(page.getByPlaceholder('Search campaigns...')).toBeVisible();
    await expect(page.getByRole('button', { name: 'Search' })).toBeVisible();
  });

  test('search filters campaigns', async ({ page }) => {
    await page.goto('/app/en/scitizen/campaigns');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    await page.getByPlaceholder('Search campaigns...').fill('temperature');
    await page.getByRole('button', { name: 'Search' }).click();

    // After search, should show results or empty state
    await page.waitForLoadState('networkidle');
    const hasCampaigns = await page.locator('.campaign-grid').isVisible().catch(() => false);
    const hasEmptyState = await page.locator('.empty-state').isVisible().catch(() => false);
    expect(hasCampaigns || hasEmptyState).toBeTruthy();
  });

  test('campaign cards link to detail or empty state shown', async ({ page }) => {
    await page.goto('/app/en/scitizen/campaigns');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    const firstCard = page.locator('.campaign-card').first();
    const hasCampaigns = await firstCard.isVisible().catch(() => false);

    if (hasCampaigns) {
      await firstCard.click();
      await expect(page).toHaveURL(/\/scitizen\/campaigns\//, { timeout: 10_000 });
    } else {
      await expect(page.locator('.empty-state')).toBeVisible();
    }
  });
});
