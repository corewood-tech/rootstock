import { test, expect } from '@playwright/test';
import { getVerificationLink, clearInbox } from './fixtures/maildev';

test.describe('login flow', () => {
  test('shows login form with required fields', async ({ page }) => {
    await page.goto('/app/en/login');
    await expect(page.getByRole('heading', { name: 'Log in to Rootstock' })).toBeVisible();
    await expect(page.getByLabel('Email')).toBeVisible();
    await expect(page.getByLabel('Password')).toBeVisible();
    await expect(page.getByRole('button', { name: 'Log in' })).toBeVisible();
  });

  test('shows error for invalid credentials', async ({ page }) => {
    await page.goto('/app/en/login');
    await page.getByLabel('Email').fill('nonexistent@test.com');
    await page.getByLabel('Password').fill('wrongpassword');
    await page.getByRole('button', { name: 'Log in' }).click();

    await expect(page.getByRole('alert')).toHaveText('Invalid email or password', {
      timeout: 10_000,
    });
  });

  test('links to register page', async ({ page }) => {
    await page.goto('/app/en/login');
    await page.getByRole('link', { name: 'Create one' }).click();
    await expect(page).toHaveURL(/\/en\/register/);
  });

  test('researcher login redirects to researcher dashboard', async ({ page, request }) => {
    await clearInbox(request);
    const user = {
      email: `login-res-${Date.now()}@rootstock.test`,
      password: 'TestPassword123!',
    };

    // Register as researcher
    await page.goto('/app/en/register');
    await page.locator('.role-option', { hasText: 'Researcher' }).click();
    await page.getByLabel('First name').fill('Login');
    await page.getByLabel('Last name').fill('Researcher');
    await page.getByLabel('Email').fill(user.email);
    await page.getByLabel('Password', { exact: true }).fill(user.password);
    await page.getByLabel('Confirm password').fill(user.password);
    await page.getByRole('button', { name: 'Create account' }).click();
    await expect(page.getByText('Check your email')).toBeVisible({ timeout: 15_000 });

    // Verify email
    const verifyLink = await getVerificationLink(request, user.email);
    await page.goto(verifyLink);
    await expect(page.getByText('Email verified')).toBeVisible({ timeout: 10_000 });

    // Login
    await page.goto('/app/en/login');
    await page.getByLabel('Email').fill(user.email);
    await page.getByLabel('Password').fill(user.password);
    await page.getByRole('button', { name: 'Log in' }).click();

    // Should redirect to researcher dashboard
    await expect(page).toHaveURL(/\/researcher/, { timeout: 15_000 });
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK');
  });

  test('scitizen login redirects to scitizen dashboard', async ({ page, request }) => {
    await clearInbox(request);
    const user = {
      email: `login-sci-${Date.now()}@rootstock.test`,
      password: 'TestPassword123!',
    };

    // Register as citizen scientist
    await page.goto('/app/en/register');
    await page.locator('.role-option', { hasText: 'Citizen Scientist' }).click();
    await page.getByLabel('First name').fill('Login');
    await page.getByLabel('Last name').fill('Scitizen');
    await page.getByLabel('Email').fill(user.email);
    await page.getByLabel('Password', { exact: true }).fill(user.password);
    await page.getByLabel('Confirm password').fill(user.password);
    await page.getByRole('button', { name: 'Create account' }).click();
    await expect(page.getByText('Check your email')).toBeVisible({ timeout: 15_000 });

    // Verify email
    const verifyLink = await getVerificationLink(request, user.email);
    await page.goto(verifyLink);
    await expect(page.getByText('Email verified')).toBeVisible({ timeout: 10_000 });

    // Login
    await page.goto('/app/en/login');
    await page.getByLabel('Email').fill(user.email);
    await page.getByLabel('Password').fill(user.password);
    await page.getByRole('button', { name: 'Log in' }).click();

    // Should redirect to scitizen dashboard
    await expect(page).toHaveURL(/\/scitizen/, { timeout: 15_000 });
  });
});
