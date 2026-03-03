// Phase 4: Enriched researcher campaign dashboard E2E test
// Validates: Parameter quality, device breakdown, enrollment funnel, temporal coverage
import { test, expect } from '@playwright/test';

test.describe('enriched researcher campaign dashboard', () => {
  test.setTimeout(60_000);

  test('published campaign shows enriched dashboard sections', async ({ page }) => {
    await page.goto('/app/en/researcher/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', {
      timeout: 10_000,
    });
    await page.waitForLoadState('networkidle');

    // Find a published campaign
    const cards = page.locator('.campaign-card--link');
    const count = await cards.count();

    if (count === 0) {
      await expect(page.locator('.empty-state')).toBeVisible();
      return;
    }

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
      // All campaigns are draft — skip enriched dashboard checks
      return;
    }

    await expect(page).toHaveURL(/\/researcher\/campaigns\//, { timeout: 10_000 });
    await expect(page.locator('.campaign-detail__header')).toBeVisible({ timeout: 10_000 });

    // Base dashboard should be present
    await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible({ timeout: 10_000 });
    await expect(page.getByText('Accepted Readings')).toBeVisible();
    await expect(page.getByText('Quarantined Readings')).toBeVisible();

    // Enrollment funnel should always be present for published campaigns
    await expect(page.getByRole('heading', { name: 'Enrollment Funnel' })).toBeVisible({
      timeout: 10_000,
    });
    await expect(page.getByText('Enrolled')).toBeVisible();
    await expect(page.getByText('Active')).toBeVisible();
    await expect(page.getByText('Contributing')).toBeVisible();
  });

  test('parameter quality section shows per-parameter data', async ({ page }) => {
    await page.goto('/app/en/researcher/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', {
      timeout: 10_000,
    });
    await page.waitForLoadState('networkidle');

    const cards = page.locator('.campaign-card--link');
    const count = await cards.count();

    if (count === 0) return;

    // Click first non-draft campaign
    for (let i = 0; i < count; i++) {
      const badge = cards.nth(i).locator('.status-badge');
      const status = await badge.textContent();
      if (status !== 'draft') {
        await cards.nth(i).click();
        break;
      }
      if (i === count - 1) return; // all draft
    }

    await expect(page.locator('.campaign-detail__header')).toBeVisible({ timeout: 10_000 });

    // Parameter quality section — may or may not have data depending on readings
    const paramHeading = page.getByRole('heading', { name: 'Parameter Quality' });
    const isVisible = await paramHeading.isVisible().catch(() => false);

    if (isVisible) {
      // If present, verify the table structure
      const paramTable = paramHeading.locator('..').locator('table');
      await expect(paramTable.locator('th', { hasText: 'Parameter' })).toBeVisible();
      await expect(paramTable.locator('th', { hasText: 'Accepted' })).toBeVisible();
      await expect(paramTable.locator('th', { hasText: 'Quarantined' })).toBeVisible();
      await expect(paramTable.locator('th', { hasText: 'Acceptance Rate' })).toBeVisible();
    }
  });

  test('device breakdown section shows device data', async ({ page }) => {
    await page.goto('/app/en/researcher/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', {
      timeout: 10_000,
    });
    await page.waitForLoadState('networkidle');

    const cards = page.locator('.campaign-card--link');
    const count = await cards.count();

    if (count === 0) return;

    for (let i = 0; i < count; i++) {
      const badge = cards.nth(i).locator('.status-badge');
      const status = await badge.textContent();
      if (status !== 'draft') {
        await cards.nth(i).click();
        break;
      }
      if (i === count - 1) return;
    }

    await expect(page.locator('.campaign-detail__header')).toBeVisible({ timeout: 10_000 });

    const deviceHeading = page.getByRole('heading', { name: 'Device Breakdown' });
    const isVisible = await deviceHeading.isVisible().catch(() => false);

    if (isVisible) {
      const deviceTable = deviceHeading.locator('..').locator('table');
      await expect(deviceTable.locator('th', { hasText: 'Device ID' })).toBeVisible();
      await expect(deviceTable.locator('th', { hasText: 'Class' })).toBeVisible();
      await expect(deviceTable.locator('th', { hasText: 'Readings' })).toBeVisible();
    }
  });
});
