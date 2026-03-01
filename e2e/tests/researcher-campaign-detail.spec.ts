// Graph node: 0x52 (Researcher Campaign Detail E2E Test)
// Validates: CampaignDetailPage (researcher), PublishCampaignFlow, CampaignDashboardFlow
import { test, expect } from '@playwright/test';

test.describe('researcher campaign detail', () => {
  test('campaign cards link to detail page', async ({ page }) => {
    await page.goto('/app/en/researcher/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
    await page.waitForLoadState('networkidle');

    const firstCard = page.locator('.campaign-card--link').first();
    const hasCampaigns = await firstCard.isVisible().catch(() => false);

    if (!hasCampaigns) {
      await expect(page.locator('.empty-state')).toBeVisible();
      return;
    }

    await firstCard.click();
    await expect(page).toHaveURL(/\/researcher\/campaigns\/[a-f0-9-]+/, { timeout: 10_000 });
  });

  test('campaign detail page renders header and sections', async ({ page }) => {
    await page.goto('/app/en/researcher/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
    await page.waitForLoadState('networkidle');

    const firstCard = page.locator('.campaign-card--link').first();
    const hasCampaigns = await firstCard.isVisible().catch(() => false);

    if (!hasCampaigns) {
      await expect(page.locator('.empty-state')).toBeVisible();
      return;
    }

    await firstCard.click();
    await expect(page).toHaveURL(/\/researcher\/campaigns\//, { timeout: 10_000 });

    // Header renders with campaign ID and status badge
    await expect(page.locator('.campaign-detail__header')).toBeVisible({ timeout: 10_000 });
    await expect(page.locator('.status-badge')).toBeVisible();

    // Back link to campaign list
    await expect(page.locator('.back-link')).toBeVisible();

    // Time window section
    await expect(page.getByRole('heading', { name: /window/i })).toBeVisible();

    // Details section
    await expect(page.getByRole('heading', { name: 'Details' })).toBeVisible();
  });

  test('draft campaign shows publish button', async ({ page }) => {
    await page.goto('/app/en/researcher/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
    await page.waitForLoadState('networkidle');

    const firstCard = page.locator('.campaign-card--link').first();
    const hasCampaigns = await firstCard.isVisible().catch(() => false);

    if (!hasCampaigns) {
      await expect(page.locator('.empty-state')).toBeVisible();
      return;
    }

    await firstCard.click();
    await expect(page).toHaveURL(/\/researcher\/campaigns\//, { timeout: 10_000 });
    await expect(page.locator('.campaign-detail__header')).toBeVisible({ timeout: 10_000 });

    const statusBadge = page.locator('.status-badge');
    const statusText = await statusBadge.textContent();

    if (statusText === 'draft') {
      // Draft campaigns have publish button
      await expect(page.getByRole('button', { name: 'Publish Campaign' })).toBeVisible();
    } else {
      // Non-draft campaigns should not show publish button
      await expect(page.getByRole('button', { name: 'Publish Campaign' })).not.toBeVisible();
    }
  });

  test('published campaign shows dashboard metrics', async ({ page }) => {
    await page.goto('/app/en/researcher/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
    await page.waitForLoadState('networkidle');

    // Find a non-draft campaign if one exists
    const cards = page.locator('.campaign-card--link');
    const count = await cards.count();

    if (count === 0) {
      await expect(page.locator('.empty-state')).toBeVisible();
      return;
    }

    // Try to find a published/active campaign
    let foundPublished = false;
    for (let i = 0; i < count; i++) {
      const badge = cards.nth(i).locator('.status-badge');
      const status = await badge.textContent();
      if (status !== 'draft') {
        await cards.nth(i).click();
        foundPublished = true;
        break;
      }
    }

    if (!foundPublished) {
      // All campaigns are draft — that's fine, just verify detail page works
      await cards.first().click();
      await expect(page.locator('.campaign-detail__header')).toBeVisible({ timeout: 10_000 });
      return;
    }

    await expect(page).toHaveURL(/\/researcher\/campaigns\//, { timeout: 10_000 });
    await expect(page.locator('.campaign-detail__header')).toBeVisible({ timeout: 10_000 });

    // Dashboard section should be visible for non-draft campaigns
    await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible({ timeout: 10_000 });
    await expect(page.locator('.stat-card')).toHaveCount(2);
    await expect(page.getByText('Accepted Readings')).toBeVisible();
    await expect(page.getByText('Quarantined Readings')).toBeVisible();
  });

  test('back link navigates to campaign list', async ({ page }) => {
    await page.goto('/app/en/researcher/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
    await page.waitForLoadState('networkidle');

    const firstCard = page.locator('.campaign-card--link').first();
    const hasCampaigns = await firstCard.isVisible().catch(() => false);

    if (!hasCampaigns) {
      await expect(page.locator('.empty-state')).toBeVisible();
      return;
    }

    await firstCard.click();
    await expect(page).toHaveURL(/\/researcher\/campaigns\//, { timeout: 10_000 });

    await page.locator('.back-link').click();
    await expect(page).toHaveURL(/\/researcher\/?$/, { timeout: 10_000 });
  });
});
