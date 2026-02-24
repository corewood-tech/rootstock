// Graph node: 0x4b (Scitizen Registration E2E Test)
// Validates: ScitzenRegistrationFlow (0x2f), ScitizenRegisterPage (0x40)
import { test, expect } from '@playwright/test';
import { getVerificationLink, clearInbox } from './fixtures/maildev';

test.describe('scitizen registration flow', () => {
  const user = {
    email: `scitizen-${Date.now()}@rootstock.test`,
    password: 'TestPassword123!',
    givenName: 'Citizen',
    familyName: 'Scientist',
  };

  test.beforeAll(async ({ request }) => {
    await clearInbox(request);
  });

  test('registers scitizen, verifies email, and reaches dashboard', async ({ page, request }) => {
    // Navigate to scitizen registration page
    await page.goto('/app/en/register-scitizen');
    await expect(page.getByRole('heading', { name: 'Join as Citizen Scientist' })).toBeVisible({ timeout: 15_000 });

    // Fill form
    await page.getByLabel('First Name').fill(user.givenName);
    await page.getByLabel('Last Name').fill(user.familyName);
    await page.getByLabel('Email').fill(user.email);
    await page.getByLabel('Password').fill(user.password);
    await page.getByLabel(/Terms of Service/).check();

    // Submit
    await page.getByRole('button', { name: 'Register' }).click();

    // Confirm success state
    await expect(page.getByText('Check your email')).toBeVisible({ timeout: 15_000 });
    await expect(page.getByText(user.email)).toBeVisible();

    // Verify email via maildev
    const verifyLink = await getVerificationLink(request, user.email);
    expect(verifyLink).toBeTruthy();

    await page.goto(verifyLink);
    await expect(page.getByText('Email verified')).toBeVisible({ timeout: 10_000 });

    // Login
    await page.goto('/app/en/login');
    await page.getByLabel('Email').fill(user.email);
    await page.getByLabel('Password').fill(user.password);
    await page.getByRole('button', { name: 'Log in' }).click();

    // Should reach an authenticated area (researcher or scitizen, depends on user_type routing)
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 15_000 });
  });

  test('shows error when ToS not accepted', async ({ page }) => {
    await page.goto('/app/en/register-scitizen');
    await expect(page.getByRole('heading', { name: 'Join as Citizen Scientist' })).toBeVisible({ timeout: 15_000 });

    await page.getByLabel('First Name').fill('Test');
    await page.getByLabel('Last Name').fill('User');
    await page.getByLabel('Email').fill('tos-test@rootstock.test');
    await page.getByLabel('Password').fill('TestPassword123!');
    // Do NOT check ToS

    await page.getByRole('button', { name: 'Register' }).click();

    await expect(page.getByRole('alert')).toHaveText('You must accept the Terms of Service');
  });

  test('links to login page', async ({ page }) => {
    await page.goto('/app/en/register-scitizen');
    await page.getByRole('link', { name: 'Log in' }).click();
    await expect(page).toHaveURL(/\/en\/login/);
  });
});
