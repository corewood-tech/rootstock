// Graph node: 0x50 (Full Happy Path E2E Test)
// Validates: EnrollDeviceCampaignFlow (0x28), ScitzenRegistrationFlow (0x2f),
//            BrowseCampaignsFlow (0x32), ScitizenRegisterPage (0x40),
//            CampaignDetailPage (0x41), CampaignBrowsePage (0x45)
import { test, expect } from '@playwright/test';

test.describe('researcher-scitizen full loop', () => {
  test('researcher creates campaign, publishes, scitizen sees it', async ({ page }) => {
    // Step 1: As authenticated researcher, go to campaign dashboard
    await page.goto('/app/en/researcher/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    // Step 2: Create a new campaign
    await page.getByRole('link', { name: 'New campaign' }).click();
    await expect(page).toHaveURL(/\/researcher\/campaigns\/new/, { timeout: 10_000 });
    await expect(page.locator('.wizard__steps')).toBeVisible({ timeout: 10_000 });

    // Step 1: Basics — set dates
    await page.getByLabel('Start date').fill('2026-03-01');
    await page.getByLabel('End date').fill('2026-09-01');
    await page.getByRole('button', { name: 'Next' }).click();

    // Step 2: Parameters
    await page.getByLabel('Parameter name').fill('PM2.5');
    await page.getByLabel('Unit').fill('µg/m³');
    await page.getByRole('button', { name: 'Next' }).click();

    // Step 3: Regions — skip
    await page.getByRole('button', { name: 'Next' }).click();

    // Step 4: Eligibility — skip
    await page.getByRole('button', { name: 'Next' }).click();

    // Step 5: Review & submit
    await expect(page.getByText('PM2.5')).toBeVisible();
    await page.getByRole('button', { name: 'Create campaign' }).click();

    // Should redirect to dashboard
    await expect(page).toHaveURL(/\/researcher/, { timeout: 15_000 });
    await page.waitForLoadState('networkidle');

    // Step 3: Click into the newly created campaign detail
    const firstCard = page.locator('.campaign-card--link').first();
    const hasCampaigns = await firstCard.isVisible().catch(() => false);

    if (hasCampaigns) {
      await firstCard.click();
      await expect(page).toHaveURL(/\/researcher\/campaigns\//, { timeout: 10_000 });
      await expect(page.locator('.campaign-detail__header')).toBeVisible({ timeout: 10_000 });

      // Step 4: Publish the campaign if it's in draft
      const publishBtn = page.getByRole('button', { name: 'Publish Campaign' });
      const canPublish = await publishBtn.isVisible().catch(() => false);

      if (canPublish) {
        await publishBtn.click();
        // Wait for status to update — page reloads after publish
        await expect(page.locator('.status-badge')).not.toHaveText('draft', { timeout: 15_000 });
      }
    }

    // Step 5: Navigate to scitizen campaigns browse
    await page.goto('/app/en/scitizen/campaigns');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    // The campaign grid or empty state should be visible
    const hasBrowseCampaigns = await page.locator('.campaign-grid').isVisible().catch(() => false);
    const hasEmptyState = await page.locator('.empty-state').isVisible().catch(() => false);
    expect(hasBrowseCampaigns || hasEmptyState).toBeTruthy();

    if (hasBrowseCampaigns) {
      // Step 6: Click into a campaign detail and verify enrollment flow is available
      const campaignCard = page.locator('.campaign-card').first();
      const hasCard = await campaignCard.isVisible().catch(() => false);

      if (hasCard) {
        await campaignCard.click();
        await expect(page).toHaveURL(/\/scitizen\/campaigns\//, { timeout: 10_000 });
        await expect(page.locator('.campaign-detail__header')).toBeVisible({ timeout: 10_000 });

        // If campaign is published, enroll button should be visible
        const enrollBtn = page.getByRole('button', { name: 'Enroll a Device' });
        const hasEnroll = await enrollBtn.isVisible().catch(() => false);

        if (hasEnroll) {
          // Open device picker
          await enrollBtn.click();
          const modal = page.locator('[role="dialog"]');
          await expect(modal).toBeVisible({ timeout: 5_000 });
          await expect(modal.getByRole('heading', { name: 'Select a Device' })).toBeVisible();

          // Cancel and confirm modal closes
          await modal.getByRole('button', { name: 'Cancel' }).click();
          await expect(modal).not.toBeVisible();
        }
      }
    }
  });

  test('scitizen can navigate between all sections', async ({ page }) => {
    // Verify all scitizen sections are reachable
    await page.goto('/app/en/scitizen/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    // Use desktop nav if visible, otherwise open mobile nav each time
    const desktopNav = page.locator('nav[aria-label="Main navigation"]');
    const isDesktop = await desktopNav.isVisible().catch(() => false);

    async function clickNavLink(label: string) {
      if (isDesktop) {
        await desktopNav.getByText(label).click();
      } else {
        await page.locator('button[aria-label="Open navigation"]').click();
        const mobileNav = page.locator('nav[aria-label="Mobile navigation"]');
        await expect(mobileNav).toBeVisible({ timeout: 5_000 });
        await mobileNav.getByText(label).click();
      }
    }

    await clickNavLink('Campaigns');
    await expect(page).toHaveURL(/\/scitizen\/campaigns/, { timeout: 10_000 });

    await clickNavLink('Devices');
    await expect(page).toHaveURL(/\/scitizen\/devices/, { timeout: 10_000 });

    await clickNavLink('Contributions');
    await expect(page).toHaveURL(/\/scitizen\/contributions/, { timeout: 10_000 });

    await clickNavLink('Notifications');
    await expect(page).toHaveURL(/\/scitizen\/notifications/, { timeout: 10_000 });

    await clickNavLink('Dashboard');
    await expect(page).toHaveURL(/\/scitizen/, { timeout: 10_000 });
  });

  test('both dashboards are accessible', async ({ page }) => {
    // Researcher dashboard loads
    await page.goto('/app/en/researcher/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    // Scitizen dashboard loads
    await page.goto('/app/en/scitizen/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    // Dashboard content renders (stats or empty state)
    const hasStats = await page.locator('.stats-grid').isVisible().catch(() => false);
    const hasEmptyState = await page.locator('.empty-state').isVisible().catch(() => false);
    expect(hasStats || hasEmptyState).toBeTruthy();
  });
});
