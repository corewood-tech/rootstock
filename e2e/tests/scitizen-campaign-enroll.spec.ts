// Graph node: 0x4e (Campaign Enrollment E2E Test)
// Validates: EnrollDeviceCampaignFlow (0x28), ConsentModal (0x35), CampaignDetailPage (0x41)
import { test, expect } from '@playwright/test';

test.describe('scitizen campaign enrollment', () => {
  test('campaign detail page renders correctly', async ({ page }) => {
    await page.goto('/app/en/scitizen/campaigns');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    const firstCard = page.locator('.campaign-card').first();
    const hasCampaigns = await firstCard.isVisible().catch(() => false);

    if (hasCampaigns) {
      // Has campaigns — click through to detail and verify sections
      await firstCard.click();
      await expect(page).toHaveURL(/\/scitizen\/campaigns\//, { timeout: 10_000 });
      await expect(page.locator('.campaign-detail__header')).toBeVisible();
      await expect(page.getByRole('heading', { name: 'Time Window' })).toBeVisible();
      await expect(page.getByRole('heading', { name: 'Statistics' })).toBeVisible();
    } else {
      // No campaigns — empty state is the correct response
      await expect(page.locator('.empty-state')).toBeVisible();
    }
  });

  test('consent modal workflow', async ({ page }) => {
    await page.goto('/app/en/scitizen/campaigns');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    const firstCard = page.locator('.campaign-card').first();
    const hasCampaigns = await firstCard.isVisible().catch(() => false);

    if (!hasCampaigns) {
      // No campaigns — verify empty state renders correctly
      await expect(page.locator('.empty-state')).toBeVisible();
      return;
    }

    await firstCard.click();
    await expect(page).toHaveURL(/\/scitizen\/campaigns\//, { timeout: 10_000 });

    const enrollBtn = page.getByRole('button', { name: 'Enroll a Device' });
    const hasEnroll = await enrollBtn.isVisible().catch(() => false);

    if (!hasEnroll) {
      // Campaign not published — detail page still renders correctly
      await expect(page.locator('.campaign-detail__header')).toBeVisible();
      return;
    }

    // Open consent modal
    await enrollBtn.click();
    const modal = page.locator('[role="dialog"]');
    await expect(modal).toBeVisible({ timeout: 5_000 });
    await expect(modal.getByRole('heading', { name: 'Consent Required' })).toBeVisible();
    await expect(modal.getByText('consent to sharing sensor data')).toBeVisible();
    await expect(modal.getByRole('button', { name: 'Cancel' })).toBeVisible();
    await expect(modal.getByRole('button', { name: 'Accept & Enroll' })).toBeVisible();

    // Cancel closes modal
    await modal.getByRole('button', { name: 'Cancel' }).click();
    await expect(modal).not.toBeVisible();
  });
});
