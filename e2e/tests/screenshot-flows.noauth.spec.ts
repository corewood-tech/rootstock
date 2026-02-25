import { test, expect } from '@playwright/test';
import { getVerificationLink, clearInbox } from './fixtures/maildev';

const SCREENSHOT_DIR = './screenshots';

test.describe('researcher registration flow screenshots', () => {
  const user = {
    email: `screenshot-res-${Date.now()}@rootstock.test`,
    password: 'TestPassword123!',
    givenName: 'Jane',
    familyName: 'Researcher',
  };

  test.beforeAll(async ({ request }) => {
    await clearInbox(request);
  });

  test('capture full researcher flow', async ({ page, request }) => {
    // 1. Register page — role selector visible
    await page.goto('/app/en/register', { waitUntil: 'networkidle' });
    await expect(page.getByLabel('First name')).toBeVisible({ timeout: 15_000 });
    await page.screenshot({ path: `${SCREENSHOT_DIR}/01-register-empty.png`, fullPage: true });

    // 2. Select researcher role
    await page.locator('.role-option', { hasText: 'Researcher' }).click();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/02-register-researcher-selected.png`, fullPage: true });

    // 3. Fill form
    await page.getByLabel('First name').fill(user.givenName);
    await page.getByLabel('Last name').fill(user.familyName);
    await page.getByLabel('Email').fill(user.email);
    await page.getByLabel('Password', { exact: true }).fill(user.password);
    await page.getByLabel('Confirm password').fill(user.password);
    await page.screenshot({ path: `${SCREENSHOT_DIR}/03-register-researcher-filled.png`, fullPage: true });

    // 4. Submit — check email screen
    await page.getByRole('button', { name: 'Create account' }).click();
    await expect(page.getByText('Check your email')).toBeVisible({ timeout: 15_000 });
    await page.screenshot({ path: `${SCREENSHOT_DIR}/04-check-email.png`, fullPage: true });

    // 5. Verify email
    const verifyLink = await getVerificationLink(request, user.email);
    await page.goto(verifyLink);
    await expect(page.getByText('Email verified')).toBeVisible({ timeout: 10_000 });
    await page.screenshot({ path: `${SCREENSHOT_DIR}/05-email-verified.png`, fullPage: true });

    // 6. Login
    await page.goto('/app/en/login');
    await page.screenshot({ path: `${SCREENSHOT_DIR}/06-login-empty.png`, fullPage: true });

    await page.getByLabel('Email').fill(user.email);
    await page.getByLabel('Password').fill(user.password);
    await page.screenshot({ path: `${SCREENSHOT_DIR}/07-login-filled.png`, fullPage: true });

    await page.getByRole('button', { name: 'Log in' }).click();
    await expect(page).toHaveURL(/\/researcher/, { timeout: 15_000 });
    await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });
    await page.screenshot({ path: `${SCREENSHOT_DIR}/08-researcher-dashboard.png`, fullPage: true });
  });
});

test.describe('citizen scientist registration flow screenshots', () => {
  const user = {
    email: `screenshot-sci-${Date.now()}@rootstock.test`,
    password: 'TestPassword123!',
    givenName: 'Alex',
    familyName: 'Scientist',
  };

  test.beforeAll(async ({ request }) => {
    await clearInbox(request);
  });

  test('capture full citizen scientist flow', async ({ page, request }) => {
    // 1. Select citizen scientist role
    await page.goto('/app/en/register', { waitUntil: 'networkidle' });
    await expect(page.getByLabel('First name')).toBeVisible({ timeout: 15_000 });
    await page.locator('.role-option', { hasText: 'Citizen Scientist' }).click();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/09-register-scitizen-selected.png`, fullPage: true });

    // 2. Fill form
    await page.getByLabel('First name').fill(user.givenName);
    await page.getByLabel('Last name').fill(user.familyName);
    await page.getByLabel('Email').fill(user.email);
    await page.getByLabel('Password', { exact: true }).fill(user.password);
    await page.getByLabel('Confirm password').fill(user.password);
    await page.screenshot({ path: `${SCREENSHOT_DIR}/10-register-scitizen-filled.png`, fullPage: true });

    // 3. Submit — check email screen
    await page.getByRole('button', { name: 'Create account' }).click();
    await expect(page.getByText('Check your email')).toBeVisible({ timeout: 15_000 });
    await page.screenshot({ path: `${SCREENSHOT_DIR}/11-scitizen-check-email.png`, fullPage: true });

    // 4. Verify email
    const verifyLink = await getVerificationLink(request, user.email);
    await page.goto(verifyLink);
    await expect(page.getByText('Email verified')).toBeVisible({ timeout: 10_000 });

    // 5. Login and verify scitizen redirect
    await page.goto('/app/en/login');
    await page.getByLabel('Email').fill(user.email);
    await page.getByLabel('Password').fill(user.password);
    await page.getByRole('button', { name: 'Log in' }).click();
    await expect(page).toHaveURL(/\/scitizen/, { timeout: 15_000 });
    await page.screenshot({ path: `${SCREENSHOT_DIR}/12-scitizen-dashboard.png`, fullPage: true });
  });
});
