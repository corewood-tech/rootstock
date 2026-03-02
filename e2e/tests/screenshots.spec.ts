import { test, expect } from '@playwright/test';
import { MockDevice } from './fixtures/mock-device';
import * as db from './fixtures/db';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const SHOTS = join(__dirname, 'screenshots');

let device: MockDevice | null = null;

test.afterAll(async () => {
  if (device) await device.disconnect();
  await db.cleanup();
});

test.describe.serial('capture screenshots', () => {
  test.setTimeout(120_000);

  let campaignID: string;
  let userID: string;

  test('01 — registration page', async ({ browser }) => {
    // Fresh context without auth to see the registration form
    const ctx = await browser.newContext();
    const page = await ctx.newPage();

    await page.goto('/app/en/register', { waitUntil: 'networkidle' });
    await expect(page.getByLabel('First name')).toBeVisible({ timeout: 15_000 });

    // Fill out the form for a realistic screenshot
    await page.locator('.role-option', { hasText: 'Researcher' }).click();
    await page.getByLabel('First name').fill('Jane');
    await page.getByLabel('Last name').fill('Mendoza');
    await page.getByLabel('Email').fill('jane.mendoza@example.org');
    await page.getByLabel('Password', { exact: true }).fill('SecurePass99!');
    await page.getByLabel('Confirm password').fill('SecurePass99!');

    await page.screenshot({ path: `${SHOTS}/01-registration.png`, fullPage: true });
    await ctx.close();
  });

  test('02 — creating a campaign (wizard)', async ({ page }) => {
    await page.goto('/app/en/researcher/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
    await page.waitForLoadState('networkidle');

    await page.getByRole('link', { name: 'New campaign' }).first().click();
    await expect(page.locator('.wizard__steps')).toBeVisible({ timeout: 10_000 });

    // Step 1: Basics
    await page.getByLabel('Start date').fill('2026-04-01');
    await page.getByLabel('End date').fill('2026-10-01');
    await page.screenshot({ path: `${SHOTS}/02a-campaign-wizard-basics.png`, fullPage: true });
    await page.getByRole('button', { name: 'Next' }).click({ force: true });

    // Step 2: Parameters
    await expect(page.getByLabel('Parameter name')).toBeVisible();
    await page.getByLabel('Parameter name').fill('PM2.5');
    await page.getByLabel('Unit').fill('µg/m³');
    await page.screenshot({ path: `${SHOTS}/02b-campaign-wizard-parameters.png`, fullPage: true });
    await page.getByRole('button', { name: 'Next' }).click({ force: true });

    // Step 3: Regions
    await expect(page.getByLabel('GeoJSON')).toBeVisible();
    await page.screenshot({ path: `${SHOTS}/02c-campaign-wizard-regions.png`, fullPage: true });
    await page.getByRole('button', { name: 'Next' }).click({ force: true });

    // Step 4: Eligibility
    await expect(page.getByLabel('Device class')).toBeVisible();
    await page.screenshot({ path: `${SHOTS}/02d-campaign-wizard-eligibility.png`, fullPage: true });
    await page.getByRole('button', { name: 'Next' }).click({ force: true });

    // Step 5: Review
    await expect(page.getByRole('button', { name: 'Create campaign' })).toBeVisible();
    await page.screenshot({ path: `${SHOTS}/02e-campaign-wizard-review.png`, fullPage: true });

    // Actually create the campaign so we have data for later screenshots
    await page.getByRole('button', { name: 'Create campaign' }).click();
    await expect(page).toHaveURL(/\/researcher/, { timeout: 15_000 });
    await page.waitForLoadState('networkidle');

    // Click into campaign detail and publish it
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

    // Publish
    const publishBtn = page.getByRole('button', { name: 'Publish Campaign' });
    await expect(publishBtn).toBeVisible({ timeout: 5_000 });
    await publishBtn.click();
    await expect(page.locator('.status-badge')).not.toHaveText('draft', { timeout: 15_000 });
  });

  test('03 — seed device, enroll, and submit reading', async ({ page, request }) => {
    // Get the authenticated user's internal ID
    userID = await db.getCampaignCreator(campaignID);

    // Create mock device (seeds DB + enrolls via HTTP)
    device = await MockDevice.create(request, userID);

    // Seed enrollment linking device → campaign → scitizen
    await db.seedEnrollment(device.deviceId, campaignID, userID);

    // Promote to 'both' role for scitizen access
    await db.promoteUserToBoth(userID);
    await db.seedScitizenProfile(userID);

    // Submit a reading via MQTT
    await device.publishReading(request, campaignID, 23.5);

    // Wait for server to process
    await page.waitForTimeout(2_000);
  });

  test('04 — accepting a campaign (enrollment flow)', async ({ page }) => {
    // Navigate to scitizen campaign browse
    await page.goto('/app/en/scitizen/campaigns');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
    await page.waitForLoadState('networkidle');

    // Wait for content to render
    const campaignGrid = page.locator('.campaign-grid');
    const emptyState = page.locator('.empty-state');
    await expect(campaignGrid.or(emptyState)).toBeVisible({ timeout: 10_000 });

    await page.screenshot({ path: `${SHOTS}/03a-campaign-browse.png`, fullPage: true });

    // Click into first campaign
    const firstCard = page.locator('.campaign-card').first();
    const hasCampaigns = await firstCard.isVisible({ timeout: 5_000 }).catch(() => false);

    if (!hasCampaigns) return;

    await firstCard.click();
    await expect(page).toHaveURL(/\/scitizen\/campaigns\//, { timeout: 10_000 });
    await expect(page.locator('.campaign-detail__header')).toBeVisible({ timeout: 10_000 });
    await page.screenshot({ path: `${SHOTS}/03b-campaign-detail.png`, fullPage: true });

    // Open enrollment modal
    const enrollBtn = page.getByRole('button', { name: 'Enroll a Device' });
    const hasEnroll = await enrollBtn.isVisible().catch(() => false);

    if (hasEnroll) {
      await enrollBtn.click();
      const modal = page.locator('[role="dialog"]');
      await expect(modal).toBeVisible({ timeout: 5_000 });

      // Wait for device list or empty message
      const deviceOption = modal.locator('.device-option').first();
      const noDevicesMsg = modal.getByText('No active devices');
      await expect(deviceOption.or(noDevicesMsg)).toBeVisible({ timeout: 10_000 });

      await page.screenshot({ path: `${SHOTS}/03c-device-picker-modal.png`, fullPage: true });

      // If a device is available, select it to show consent modal
      const hasDevice = await deviceOption.isVisible().catch(() => false);
      if (hasDevice) {
        await deviceOption.click();
        await expect(modal.getByRole('heading', { name: 'Consent Required' })).toBeVisible({ timeout: 5_000 });
        await page.screenshot({ path: `${SHOTS}/03d-consent-modal.png`, fullPage: true });
        await modal.getByRole('button', { name: 'Cancel' }).click();
      } else {
        await modal.getByRole('button', { name: 'Cancel' }).click();
      }
    }
  });

  test('05 — device management', async ({ page }) => {
    await page.goto('/app/en/scitizen/devices');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
    await page.waitForLoadState('networkidle');

    const deviceGrid = page.locator('.device-grid');
    const emptyState = page.locator('.empty-state');
    await expect(deviceGrid.or(emptyState)).toBeVisible({ timeout: 10_000 });

    await page.screenshot({ path: `${SHOTS}/04a-device-list.png`, fullPage: true });

    // Click into device detail if available
    const firstCard = page.locator('.device-card').first();
    const hasDevices = await firstCard.isVisible({ timeout: 3_000 }).catch(() => false);

    if (hasDevices) {
      await firstCard.click();
      await expect(page).toHaveURL(/\/scitizen\/devices\//, { timeout: 10_000 });
      await expect(page.locator('.info-grid')).toBeVisible({ timeout: 10_000 });
      await page.screenshot({ path: `${SHOTS}/04b-device-detail.png`, fullPage: true });
    }
  });

  test('06 — receiving data (dashboard + contributions)', async ({ page }) => {
    // Scitizen dashboard with stats
    await page.goto('/app/en/scitizen/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
    await page.waitForLoadState('networkidle');

    // Poll for stats to appear
    await expect(async () => {
      await page.reload();
      await expect(page.locator('.stats-grid')).toBeVisible({ timeout: 5_000 });
    }).toPass({ intervals: [2_000, 3_000, 5_000], timeout: 30_000 });

    await page.screenshot({ path: `${SHOTS}/05a-scitizen-dashboard.png`, fullPage: true });

    // Contributions page
    await page.goto('/app/en/scitizen/contributions');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
    await page.waitForLoadState('networkidle');

    await page.screenshot({ path: `${SHOTS}/05b-contributions.png`, fullPage: true });
  });
});
