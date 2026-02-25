import { test, expect } from '@playwright/test';
import { getVerificationLink, clearInbox } from './fixtures/maildev';

test.describe('registration flow', () => {
  const user = {
    email: `reg-${Date.now()}@rootstock.test`,
    password: 'TestPassword123!',
    givenName: 'Registration',
    familyName: 'Test',
  };

  test.beforeAll(async ({ request }) => {
    await clearInbox(request);
  });

  test('registers researcher, sends verification email, and verifies', async ({ page, request }) => {
    await page.goto('/app/en/register');
    await expect(page.getByRole('heading', { name: 'Create your account' })).toBeVisible();

    // Select researcher role
    await page.locator('.role-option', { hasText: 'Researcher' }).click();

    // Fill form using semantic locators
    await page.getByLabel('First name').fill(user.givenName);
    await page.getByLabel('Last name').fill(user.familyName);
    await page.getByLabel('Email').fill(user.email);
    await page.getByLabel('Password', { exact: true }).fill(user.password);
    await page.getByLabel('Confirm password').fill(user.password);

    // Submit
    await page.getByRole('button', { name: 'Create account' }).click();

    // Confirm "check your email" screen appears
    await expect(page.getByText('Check your email')).toBeVisible({ timeout: 15_000 });
    await expect(page.getByText('verification link')).toBeVisible();

    // Verify email via maildev
    const verifyLink = await getVerificationLink(request, user.email);
    expect(verifyLink).toBeTruthy();

    // Navigate to verification link
    await page.goto(verifyLink);
    await expect(page.getByText('Email verified')).toBeVisible({ timeout: 10_000 });

    // Verify we can now login
    await page.goto('/app/en/login');
    await page.getByLabel('Email').fill(user.email);
    await page.getByLabel('Password').fill(user.password);
    await page.getByRole('button', { name: 'Log in' }).click();

    // Should reach researcher dashboard
    await expect(page).toHaveURL(/\/researcher/, { timeout: 15_000 });
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
  });

  test('registers scitizen and redirects to scitizen dashboard', async ({ page, request }) => {
    await clearInbox(request);
    const scitizenUser = {
      email: `scitizen-${Date.now()}@rootstock.test`,
      password: 'TestPassword123!',
      givenName: 'Citizen',
      familyName: 'Scientist',
    };

    await page.goto('/app/en/register');

    // Select citizen scientist role
    await page.locator('.role-option', { hasText: 'Citizen Scientist' }).click();

    await page.getByLabel('First name').fill(scitizenUser.givenName);
    await page.getByLabel('Last name').fill(scitizenUser.familyName);
    await page.getByLabel('Email').fill(scitizenUser.email);
    await page.getByLabel('Password', { exact: true }).fill(scitizenUser.password);
    await page.getByLabel('Confirm password').fill(scitizenUser.password);

    await page.getByRole('button', { name: 'Create account' }).click();
    await expect(page.getByText('Check your email')).toBeVisible({ timeout: 15_000 });

    const verifyLink = await getVerificationLink(request, scitizenUser.email);
    await page.goto(verifyLink);
    await expect(page.getByText('Email verified')).toBeVisible({ timeout: 10_000 });

    await page.goto('/app/en/login');
    await page.getByLabel('Email').fill(scitizenUser.email);
    await page.getByLabel('Password').fill(scitizenUser.password);
    await page.getByRole('button', { name: 'Log in' }).click();

    // Should reach scitizen dashboard
    await expect(page).toHaveURL(/\/scitizen/, { timeout: 15_000 });
  });

  test('shows error when no role selected', async ({ page }) => {
    await page.goto('/app/en/register');

    await page.getByLabel('First name').fill('Test');
    await page.getByLabel('Last name').fill('User');
    await page.getByLabel('Email').fill('norole@test.com');
    await page.getByLabel('Password', { exact: true }).fill('Password123!');
    await page.getByLabel('Confirm password').fill('Password123!');

    await page.getByRole('button', { name: 'Create account' }).click();

    await expect(page.getByRole('alert')).toBeVisible();
  });

  test('shows error for mismatched passwords', async ({ page }) => {
    await page.goto('/app/en/register');

    await page.locator('.role-option', { hasText: 'Researcher' }).click();
    await page.getByLabel('First name').fill('Test');
    await page.getByLabel('Last name').fill('User');
    await page.getByLabel('Email').fill('mismatch@test.com');
    await page.getByLabel('Password', { exact: true }).fill('Password123!');
    await page.getByLabel('Confirm password').fill('Different456!');

    await page.getByRole('button', { name: 'Create account' }).click();

    await expect(page.getByRole('alert')).toHaveText('Passwords do not match');
  });

  test('links to login page', async ({ page }) => {
    await page.goto('/app/en/register');
    await page.getByRole('link', { name: 'Log in' }).click();
    await expect(page).toHaveURL(/\/en\/login/);
  });
});
