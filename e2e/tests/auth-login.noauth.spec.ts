import { test, expect } from '@playwright/test';

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
});
