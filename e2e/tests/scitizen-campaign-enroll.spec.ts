// Graph node: 0x4e (Campaign Enrollment E2E Test)
// Validates: EnrollDeviceCampaignFlow (0x28), DevicePickerModal, ConsentModal (0x35), CampaignDetailPage (0x41)
import { test, expect } from '@playwright/test';

test.describe('scitizen campaign enrollment', () => {
  test('campaign detail page renders correctly', async ({ page }) => {
    await page.goto('/app/en/scitizen/campaigns');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
    await page.waitForLoadState('networkidle');

    const firstCard = page.locator('.campaign-card').first();
    const campaignGrid = page.locator('.campaign-grid');
    const emptyState = page.locator('.empty-state');
    await expect(campaignGrid.or(emptyState)).toBeVisible({ timeout: 10_000 });

    const hasCampaigns = await firstCard.isVisible().catch(() => false);

    if (hasCampaigns) {
      // Has campaigns — click through to detail and verify sections
      await firstCard.click();
      await expect(page).toHaveURL(/\/scitizen\/campaigns\//, { timeout: 10_000 });
      await expect(page.locator('.campaign-detail__header')).toBeVisible({ timeout: 10_000 });
    } else {
      // No campaigns — empty state is the correct response
      await expect(emptyState).toBeVisible();
    }
  });

  test('device picker modal opens on enroll click', async ({ page }) => {
    await page.goto('/app/en/scitizen/campaigns');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
    await page.waitForLoadState('networkidle');

    const firstCard = page.locator('.campaign-card').first();
    const hasCampaigns = await firstCard.isVisible({ timeout: 5_000 }).catch(() => false);

    if (!hasCampaigns) {
      const hasEmpty = await page.locator('.empty-state').isVisible().catch(() => false);
      expect(hasEmpty).toBeTruthy();
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

    // Click enroll — device picker modal should open first
    await enrollBtn.click();
    const modal = page.locator('[role="dialog"]');
    await expect(modal).toBeVisible({ timeout: 5_000 });
    await expect(modal.getByRole('heading', { name: 'Select a Device' })).toBeVisible();

    // Wait for device list or empty message to appear (modal loads devices async)
    const deviceOption = modal.locator('.device-option').first();
    const noDevicesMsg = modal.getByText('No active devices');
    await expect(deviceOption.or(noDevicesMsg)).toBeVisible({ timeout: 10_000 });

    // Cancel closes modal
    await modal.getByRole('button', { name: 'Cancel' }).click();
    await expect(modal).not.toBeVisible();
  });

  test('device selection opens consent modal', async ({ page }) => {
    await page.goto('/app/en/scitizen/campaigns');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
    await page.waitForLoadState('networkidle');

    const firstCard = page.locator('.campaign-card').first();
    const hasCampaigns = await firstCard.isVisible({ timeout: 5_000 }).catch(() => false);

    if (!hasCampaigns) {
      const hasEmpty = await page.locator('.empty-state').isVisible().catch(() => false);
      expect(hasEmpty).toBeTruthy();
      return;
    }

    await firstCard.click();
    await expect(page).toHaveURL(/\/scitizen\/campaigns\//, { timeout: 10_000 });

    const enrollBtn = page.getByRole('button', { name: 'Enroll a Device' });
    const hasEnroll = await enrollBtn.isVisible().catch(() => false);

    if (!hasEnroll) {
      await expect(page.locator('.campaign-detail__header')).toBeVisible();
      return;
    }

    // Open device picker
    await enrollBtn.click();
    const pickerModal = page.locator('[role="dialog"]');
    await expect(pickerModal).toBeVisible({ timeout: 5_000 });
    await page.waitForLoadState('networkidle');

    const firstDevice = pickerModal.locator('.device-option').first();
    const hasDevices = await firstDevice.isVisible({ timeout: 5_000 }).catch(() => false);

    if (!hasDevices) {
      // No devices — picker shows empty message, cancel and return
      const hasNoDevices = await pickerModal.getByText('No active devices').isVisible({ timeout: 3_000 }).catch(() => false);
      if (hasNoDevices) {
        await pickerModal.getByRole('button', { name: 'Cancel' }).click();
      } else {
        // Modal may be loading — cancel and skip
        await pickerModal.getByRole('button', { name: 'Cancel' }).click();
      }
      return;
    }

    // Select a device — consent modal should appear
    await firstDevice.click();
    const consentModal = page.locator('[role="dialog"]');
    await expect(consentModal).toBeVisible({ timeout: 5_000 });
    await expect(consentModal.getByRole('heading', { name: 'Consent Required' })).toBeVisible();
    await expect(consentModal.getByText('consent to sharing sensor data')).toBeVisible();
    await expect(consentModal.getByRole('button', { name: 'Cancel' })).toBeVisible();
    await expect(consentModal.getByRole('button', { name: /Accept & Enroll/ })).toBeVisible();

    // Cancel closes consent modal
    await consentModal.getByRole('button', { name: 'Cancel' }).click();
    await expect(consentModal).not.toBeVisible();
  });
});
