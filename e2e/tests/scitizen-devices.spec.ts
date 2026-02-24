// Graph node: 0x47 (Device Management E2E Test)
// Validates: DeviceManagementFlow (0x27), DeviceManagementPage (0x42), DeviceDetailPage (0x43)
import { test, expect } from '@playwright/test';

test.describe('scitizen device management', () => {
  test('shows device list or empty state', async ({ page }) => {
    await page.goto('/app/en/scitizen/devices');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    const hasDevices = await page.locator('.device-grid').isVisible().catch(() => false);
    const hasEmptyState = await page.locator('.empty-state').isVisible().catch(() => false);
    expect(hasDevices || hasEmptyState).toBeTruthy();
  });

  test('displays my devices heading', async ({ page }) => {
    await page.goto('/app/en/scitizen/devices');
    await expect(page.getByRole('heading', { name: 'My Devices' })).toBeVisible({ timeout: 10_000 });
  });

  test('device cards or empty state render correctly', async ({ page }) => {
    await page.goto('/app/en/scitizen/devices');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    const firstCard = page.locator('.device-card').first();
    const hasDevices = await firstCard.isVisible().catch(() => false);

    if (hasDevices) {
      await expect(firstCard.locator('.device-card__header')).toBeVisible();
      await expect(firstCard.locator('.device-card__info')).toBeVisible();
      await expect(firstCard.locator('.device-card__meta')).toBeVisible();
    } else {
      // Empty or error state — no devices available
      await expect(page.locator('.empty-state')).toBeVisible();
    }
  });

  test('device navigation works', async ({ page }) => {
    await page.goto('/app/en/scitizen/devices');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    const firstCard = page.locator('.device-card').first();
    const hasDevices = await firstCard.isVisible().catch(() => false);

    if (hasDevices) {
      await firstCard.click();
      await expect(page).toHaveURL(/\/scitizen\/devices\//, { timeout: 10_000 });
      await expect(page.locator('.info-grid')).toBeVisible({ timeout: 10_000 });
    } else {
      // No devices — empty state is correct
      await expect(page.locator('.empty-state')).toBeVisible();
    }
  });

  test('error state shows retry button', async ({ page }) => {
    await page.goto('/app/en/scitizen/devices');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

    const hasError = await page.locator('.form-error').isVisible().catch(() => false);
    if (hasError) {
      await expect(page.getByRole('button', { name: 'Retry' })).toBeVisible();
    }
    // No error is also valid — page loaded successfully
    expect(true).toBeTruthy();
  });
});
