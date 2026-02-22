import { test, expect } from '@playwright/test';

const PROMETHEUS_BASE = 'http://prometheus:9090/prometheus';
const METRIC_QUERY = 'span_metrics_calls_total{service_name="rootstock"}';

test('health check produces ok status and increments span metrics', async ({ page, request }) => {
  // 1. Query Prometheus for baseline metric value
  const baselineResponse = await request.get(`${PROMETHEUS_BASE}/api/v1/query`, {
    params: { query: METRIC_QUERY },
    timeout: 5_000,
  });
  expect(baselineResponse.ok(), `Prometheus returned ${baselineResponse.status()}`).toBeTruthy();

  const baselineData = await baselineResponse.json();
  const baselineResults = baselineData.data?.result ?? [];
  const baseline = baselineResults.reduce(
    (sum: number, r: { value: [number, string] }) => sum + parseFloat(r.value[1]),
    0,
  );

  // 2. Navigate to the app and click the Health Check button multiple times
  await page.goto('/app/en/', { timeout: 10_000, waitUntil: 'networkidle' });
  for (let i = 0; i < 3; i++) {
    await page.getByRole('button', { name: 'Health Check' }).click({ timeout: 10_000 });
    await expect(page.getByRole('status')).toHaveText('ok', { timeout: 10_000 });
  }

  // 4. Poll Prometheus until span_metrics_calls_total increments above baseline
  await expect(async () => {
    const response = await request.get(`${PROMETHEUS_BASE}/api/v1/query`, {
      params: { query: METRIC_QUERY },
      timeout: 5_000,
    });
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    const results = data.data?.result ?? [];
    const current = results.reduce(
      (sum: number, r: { value: [number, string] }) => sum + parseFloat(r.value[1]),
      0,
    );

    expect(current).toBeGreaterThan(baseline);
  }).toPass({
    intervals: [2_000],
    timeout: 30_000,
  });
});
