import { test as setup, expect } from '@playwright/test';
import { getVerificationLink } from './fixtures/maildev';

const STORAGE_STATE = './tests/.auth/user.json';

const TEST_USER = {
  email: `test-${Date.now()}@rootstock.test`,
  password: 'TestPassword123!',
  givenName: 'Test',
  familyName: 'Researcher',
};

setup('register, verify email, and login', async ({ page, request }) => {
  // 1. Clear maildev inbox
  await request.delete('http://caddy:9999/maildev/email/all');

  // 2. Warm up vite (first load triggers dep optimization + page reload)
  await page.goto('/app/en/', { waitUntil: 'networkidle', timeout: 30_000 });

  // 3. Register via UI — select researcher role
  await page.goto('/app/en/register', { waitUntil: 'networkidle' });
  await expect(page.getByLabel('First name')).toBeVisible({ timeout: 15_000 });

  await page.locator('.role-option', { hasText: 'Researcher' }).click();
  await page.getByLabel('First name').fill(TEST_USER.givenName);
  await page.getByLabel('Last name').fill(TEST_USER.familyName);
  await page.getByLabel('Email').fill(TEST_USER.email);
  await page.getByLabel('Password', { exact: true }).fill(TEST_USER.password);
  await page.getByLabel('Confirm password').fill(TEST_USER.password);
  await page.getByRole('button', { name: 'Create account' }).click();

  // 4. Wait for "check your email" confirmation
  await expect(page.getByText('Check your email')).toBeVisible({ timeout: 15_000 });

  // 5. Get verification link from maildev
  const verifyLink = await getVerificationLink(request, TEST_USER.email);

  // 6. Navigate to verification link
  await page.goto(verifyLink);
  await expect(page.getByText('Email verified')).toBeVisible({ timeout: 10_000 });

  // 7. Login
  await page.goto('/app/en/login');
  await page.getByLabel('Email').fill(TEST_USER.email);
  await page.getByLabel('Password').fill(TEST_USER.password);
  await page.getByRole('button', { name: 'Log in' }).click();

  // 8. Wait for researcher dashboard (researcher role selected at registration)
  await expect(page).toHaveURL(/\/researcher/, { timeout: 15_000 });
  await expect(page.locator('.app-header__brand-name')).toHaveText('ROOTSTOCK', { timeout: 10_000 });

  // 9. Save storage state
  await page.context().storageState({ path: STORAGE_STATE });
});
