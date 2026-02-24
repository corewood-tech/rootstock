<script lang="ts">
	// Graph node: 0x44 (DeviceDetailPage)
	// Implements: FR-041 (Device Dashboard), FR-093 (Connection History)
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';
	import { scitizenService } from '$lib/api/clients';
	import type { GetDeviceDetailResponse } from '$lib/api/gen/rootstock/v1/rootstock_pb';

	const deviceId = $derived($page.params.id);

	let detail = $state<GetDeviceDetailResponse | null>(null);
	let loading = $state(true);
	let error = $state('');

	$effect(() => {
		loadDetail();
	});

	async function loadDetail() {
		loading = true;
		try {
			detail = await scitizenService.getDeviceDetail({ deviceId });
		} catch (e: any) {
			error = e.message || 'Failed to load device';
		} finally {
			loading = false;
		}
	}

	function formatDate(dateStr?: string): string {
		if (!dateStr) return '—';
		try { return new Date(dateStr).toLocaleString(); } catch { return dateStr; }
	}
</script>

<div class="device-detail">
	{#if loading}
		<div class="empty-state" role="status">
			<p class="text-secondary">{$t('common.loading')}</p>
		</div>
	{:else if error}
		<div class="empty-state">
			<p class="form-error">{error}</p>
		</div>
	{:else if detail?.device}
		{@const d = detail.device}
		<h1 class="heading heading--lg">Device {d.id.slice(0, 8)}</h1>

		<section class="detail-section" aria-label="Device info">
			<div class="info-grid">
				<div><strong>Status:</strong> <span class="status-badge status-badge--{d.status}">{d.status}</span></div>
				<div><strong>Class:</strong> {d.class}</div>
				<div><strong>Tier:</strong> {d.tier}</div>
				<div><strong>Firmware:</strong> {d.firmwareVersion}</div>
				<div><strong>Created:</strong> {formatDate(d.createdAt)}</div>
				{#if d.sensors.length > 0}
					<div><strong>Sensors:</strong> {d.sensors.join(', ')}</div>
				{/if}
			</div>
		</section>

		{#if detail.enrollments.length > 0}
			<section class="detail-section" aria-label="Enrollments">
				<h2 class="heading heading--md">Campaign Enrollments</h2>
				{#each detail.enrollments as enrollment (enrollment.id)}
					<div class="enrollment-row">
						<span>Campaign: {enrollment.campaignId.slice(0, 8)}</span>
						<span class="status-badge status-badge--{enrollment.status}">{enrollment.status}</span>
						<span>{formatDate(enrollment.enrolledAt)}</span>
					</div>
				{/each}
			</section>
		{/if}

		{#if detail.connectionHistory.length > 0}
			<section class="detail-section" aria-label="Connection history">
				<h2 class="heading heading--md">Connection History</h2>
				<table class="data-table" role="table">
					<thead>
						<tr><th>Event</th><th>Time</th><th>Reason</th></tr>
					</thead>
					<tbody>
						{#each detail.connectionHistory as event}
							<tr>
								<td>{event.eventType}</td>
								<td>{formatDate(event.timestamp)}</td>
								<td>{event.reason ?? '—'}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</section>
		{/if}
	{/if}
</div>

<style>
	.detail-section { margin-bottom: 2rem; }
	.info-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
		gap: 0.75rem;
	}
	.enrollment-row {
		display: flex;
		gap: 1rem;
		align-items: center;
		padding: 0.5rem 0;
		border-bottom: 1px solid var(--color-border, #333);
	}
	.data-table { width: 100%; border-collapse: collapse; }
	.data-table th, .data-table td {
		padding: 0.5rem 1rem;
		text-align: left;
		border-bottom: 1px solid var(--color-border, #333);
	}
</style>
