// Graph node: 0x60 (Full Happy-Path E2E: Device + Data Loop)
// Validates: RegisterDeviceFlow, IngestReadingFlow, RefreshScitizenScoreFlow,
//            ScitzenDashboardFlow, campaign creation/publish UI
import { test, expect } from '@playwright/test';
import { MockDevice } from './fixtures/mock-device';
import * as db from './fixtures/db';

let device: MockDevice | null = null;

test.afterAll(async () => {
  if (device) await device.disconnect();
  await db.cleanup();
});

test.describe('full happy-path: campaign → device → reading → dashboard', () => {
  test.setTimeout(120_000);

  let campaignID: string;
  let userID: string;

  test('end-to-end data loop', async ({ page, request }) => {
    // ─── Phase 1: Researcher creates + publishes campaign ───

    await page.goto('/app/en/researcher/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', {
      timeout: 10_000,
    });

    // Create campaign via wizard
    await page.getByRole('link', { name: 'New campaign' }).click();
    await expect(page).toHaveURL(/\/researcher\/campaigns\/new/, { timeout: 10_000 });
    await expect(page.locator('.wizard__steps')).toBeVisible({ timeout: 10_000 });

    // Step 1: Basics — dates
    await page.getByLabel('Start date').fill('2026-03-01');
    await page.getByLabel('End date').fill('2026-09-01');
    // force: true needed — on mobile viewports the .repeater overlay intercepts clicks
    await page.getByRole('button', { name: 'Next' }).click({ force: true });

    // Step 2: Parameters — PM2.5
    await page.getByLabel('Parameter name').fill('PM2.5');
    await page.getByLabel('Unit').fill('µg/m³');
    await page.getByRole('button', { name: 'Next' }).click({ force: true });

    // Step 3: Regions — skip
    await page.getByRole('button', { name: 'Next' }).click({ force: true });

    // Step 4: Eligibility — skip
    await page.getByRole('button', { name: 'Next' }).click({ force: true });

    // Step 5: Review & submit
    await expect(page.getByText('PM2.5')).toBeVisible();
    await page.getByRole('button', { name: 'Create campaign' }).click();

    // Should redirect to researcher dashboard
    await expect(page).toHaveURL(/\/researcher/, { timeout: 15_000 });
    await page.waitForLoadState('networkidle');

    // Click into the newly created campaign
    const firstCard = page.locator('.campaign-card--link').first();
    await expect(firstCard).toBeVisible({ timeout: 10_000 });
    await firstCard.click();
    await expect(page).toHaveURL(/\/researcher\/campaigns\//, { timeout: 10_000 });
    await expect(page.locator('.campaign-detail__header')).toBeVisible({ timeout: 10_000 });

    // Extract campaign ID from URL
    const url = page.url();
    const match = url.match(/\/researcher\/campaigns\/([^/]+)/);
    expect(match).toBeTruthy();
    campaignID = match![1];

    // Publish the campaign
    const publishBtn = page.getByRole('button', { name: 'Publish Campaign' });
    await expect(publishBtn).toBeVisible({ timeout: 5_000 });
    await publishBtn.click();
    await expect(page.locator('.status-badge')).not.toHaveText('draft', { timeout: 15_000 });

    // ─── Phase 2: Seed device + enroll ───

    // Get the authenticated user's internal ID via the campaign they just created
    userID = await db.getCampaignCreator(campaignID);

    // Create mock device (seeds DB + enrolls via HTTP)
    device = await MockDevice.create(request, userID);

    // Seed campaign enrollment linking device → campaign → scitizen
    await db.seedEnrollment(device.deviceId, campaignID, userID);

    // Promote researcher to 'both' role so they can access scitizen dashboard
    await db.promoteUserToBoth(userID);
    await db.seedScitizenProfile(userID);

    // ─── Phase 3: Submit reading via MQTT ───

    await device.publishReading(request, campaignID, 23.5);

    // Brief wait for server to process the reading inline
    await page.waitForTimeout(2_000);

    // ─── Phase 4: Verify on scitizen dashboard ───

    await page.goto('/app/en/scitizen/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', {
      timeout: 10_000,
    });

    // Poll for stats to appear — the dashboard should show updated data
    await expect(async () => {
      await page.reload();
      const statsGrid = page.locator('.stats-grid');
      await expect(statsGrid).toBeVisible({ timeout: 5_000 });

      // Verify Total Readings value is > 0
      const readingsValue = statsGrid
        .locator('.stat-card', { has: page.getByText('Total Readings') })
        .locator('.stat-card__value');
      const value = await readingsValue.textContent();
      expect(Number(value)).toBeGreaterThan(0);
    }).toPass({ intervals: [2_000, 3_000, 5_000], timeout: 30_000 });
  });
});
