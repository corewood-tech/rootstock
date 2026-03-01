// Graph node: 0x4e (Campaign Enrollment E2E Test)
// Validates: EnrollDeviceCampaignFlow (0x28), DevicePickerModal, ConsentModal (0x35), CampaignDetailPage (0x41)
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

  test('device picker modal opens on enroll click', async ({ page }) => {
    await page.goto('/app/en/scitizen/campaigns');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    const firstCard = page.locator('.campaign-card').first();
    const hasCampaigns = await firstCard.isVisible().catch(() => false);

    if (!hasCampaigns) {
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

    // Click enroll — device picker modal should open first
    await enrollBtn.click();
    const modal = page.locator('[role="dialog"]');
    await expect(modal).toBeVisible({ timeout: 5_000 });
    await expect(modal.getByRole('heading', { name: 'Select a Device' })).toBeVisible();

    // Should show device list or "No active devices" message
    const hasDevices = await modal.locator('.device-option').first().isVisible().catch(() => false);
    const hasNoDevices = await modal.getByText('No active devices').isVisible().catch(() => false);
    expect(hasDevices || hasNoDevices).toBeTruthy();

    // Cancel closes modal
    await modal.getByRole('button', { name: 'Cancel' }).click();
    await expect(modal).not.toBeVisible();
  });

  test('device selection opens consent modal', async ({ page }) => {
    await page.goto('/app/en/scitizen/campaigns');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    const firstCard = page.locator('.campaign-card').first();
    const hasCampaigns = await firstCard.isVisible().catch(() => false);

    if (!hasCampaigns) {
      await expect(page.locator('.empty-state')).toBeVisible();
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

    const firstDevice = pickerModal.locator('.device-option').first();
    const hasDevices = await firstDevice.isVisible().catch(() => false);

    if (!hasDevices) {
      // No devices — picker shows empty message, cancel and return
      await expect(pickerModal.getByText('No active devices')).toBeVisible();
      await pickerModal.getByRole('button', { name: 'Cancel' }).click();
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
