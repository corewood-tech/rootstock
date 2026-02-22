import { test, expect } from '@playwright/test';

test.describe('campaign dashboard', () => {
  test('shows campaign list or empty state', async ({ page }) => {
    await page.goto('/app/en/researcher/');
    // Wait for the dashboard to render (either state)
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
    const hasCampaigns = await page.locator('.campaign-list').isVisible().catch(() => false);
    const hasEmptyState = await page.locator('.empty-state').isVisible().catch(() => false);
    expect(hasCampaigns || hasEmptyState).toBeTruthy();
  });

  test('navigates to campaign creation form', async ({ page }) => {
    await page.goto('/app/en/researcher/');
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
    await page.getByRole('link', { name: 'New campaign' }).click();
    await expect(page).toHaveURL(/\/researcher\/campaigns\/new/, { timeout: 10_000 });
  });
});

test.describe('campaign creation wizard', () => {
  test('shows step indicator with 5 steps', async ({ page }) => {
    await page.goto('/app/en/researcher/campaigns/new');
    const steps = page.locator('.wizard__steps');
    await expect(steps).toBeVisible({ timeout: 10_000 });
    await expect(steps.getByText('Basics')).toBeVisible();
    await expect(steps.getByText('Parameters')).toBeVisible();
    await expect(steps.getByText('Regions')).toBeVisible();
    await expect(steps.getByText('Eligibility')).toBeVisible();
    await expect(steps.getByText('Review')).toBeVisible();
  });

  test('navigates through wizard steps', async ({ page }) => {
    await page.goto('/app/en/researcher/campaigns/new');
    await expect(page.locator('.wizard__steps')).toBeVisible({ timeout: 10_000 });

    // Step 1: Basics
    await expect(page.getByLabel('Start date')).toBeVisible();
    await page.getByRole('button', { name: 'Next' }).click();

    // Step 2: Parameters
    await expect(page.getByLabel('Parameter name')).toBeVisible();
    await page.getByRole('button', { name: 'Next' }).click();

    // Step 3: Regions
    await expect(page.getByLabel('GeoJSON')).toBeVisible();
    await page.getByRole('button', { name: 'Next' }).click();

    // Step 4: Eligibility
    await expect(page.getByLabel('Device class')).toBeVisible();
    await page.getByRole('button', { name: 'Next' }).click();

    // Step 5: Review
    await expect(page.getByRole('button', { name: 'Create campaign' })).toBeVisible();
  });

  test('can go back through steps', async ({ page }) => {
    await page.goto('/app/en/researcher/campaigns/new');
    await expect(page.locator('.wizard__steps')).toBeVisible({ timeout: 10_000 });

    await page.getByRole('button', { name: 'Next' }).click();
    await expect(page.getByLabel('Parameter name')).toBeVisible();

    await page.getByRole('button', { name: 'Back' }).click();
    await expect(page.getByLabel('Start date')).toBeVisible();
  });

  test('adds and removes parameters', async ({ page }) => {
    await page.goto('/app/en/researcher/campaigns/new');
    await expect(page.locator('.wizard__steps')).toBeVisible({ timeout: 10_000 });
    await page.getByRole('button', { name: 'Next' }).click();

    // Add a parameter
    await page.getByRole('button', { name: 'Add parameter' }).click();
    const paramInputs = page.getByLabel('Parameter name');
    await expect(paramInputs).toHaveCount(2);

    // Remove a parameter
    const removeButtons = page.getByRole('button', { name: 'Remove' });
    await removeButtons.first().click();
    await expect(page.getByLabel('Parameter name')).toHaveCount(1);
  });

  test('creates campaign with basic data', async ({ page }) => {
    await page.goto('/app/en/researcher/campaigns/new');
    await expect(page.locator('.wizard__steps')).toBeVisible({ timeout: 10_000 });

    // Step 1: Set dates
    await page.getByLabel('Start date').fill('2026-03-01');
    await page.getByLabel('End date').fill('2026-06-01');
    await page.getByRole('button', { name: 'Next' }).click();

    // Step 2: Add parameter
    await page.getByLabel('Parameter name').fill('Temperature');
    await page.getByLabel('Unit').fill('°C');
    await page.getByRole('button', { name: 'Next' }).click();

    // Step 3: Skip regions
    await page.getByRole('button', { name: 'Next' }).click();

    // Step 4: Skip eligibility
    await page.getByRole('button', { name: 'Next' }).click();

    // Step 5: Review & submit
    await expect(page.getByText('Temperature')).toBeVisible();
    await page.getByRole('button', { name: 'Create campaign' }).click();

    // Should redirect to dashboard
    await expect(page).toHaveURL(/\/researcher/, { timeout: 15_000 });
  });
});
