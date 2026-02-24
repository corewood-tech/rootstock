<script lang="ts">
	// Graph node: 0x43 (DeviceListPage)
	// Implements: FR-041 (Device Management Dashboard)
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';
	import { scitizenService } from '$lib/api/clients';
	import type { DeviceSummaryProto } from '$lib/api/gen/rootstock/v1/rootstock_pb';

	const lang = $derived($page.params.lang);

	let devices = $state<DeviceSummaryProto[]>([]);
	let loading = $state(true);
	let error = $state('');

	$effect(() => {
		loadDevices();
	});

	async function loadDevices() {
		loading = true;
		error = '';
		try {
			const resp = await scitizenService.getDevices({});
			devices = resp.devices;
		} catch (e: any) {
			error = e.message || 'Failed to load devices';
		} finally {
			loading = false;
		}
	}

	function formatDate(dateStr?: string): string {
		if (!dateStr) return 'Never';
		try { return new Date(dateStr).toLocaleString(); } catch { return dateStr; }
	}
</script>

<div class="devices-page">
	<h1 class="heading heading--lg">My Devices</h1>

	{#if loading}
		<div class="empty-state" role="status">
			<p class="text-secondary">{$t('common.loading')}</p>
		</div>
	{:else if error}
		<div class="empty-state">
			<p class="form-error">{error}</p>
			<button onclick={loadDevices} class="btn btn--secondary mt-4">Retry</button>
		</div>
	{:else if devices.length === 0}
		<div class="empty-state">
			<p>No devices registered yet.</p>
		</div>
	{:else}
		<div class="device-grid" role="list">
			{#each devices as device (device.id)}
				<a href="{base}/{lang}/scitizen/devices/{device.id}" class="device-card" role="listitem">
					<div class="device-card__header">
						<span class="device-card__id">{device.id.slice(0, 8)}</span>
						<span class="status-badge status-badge--{device.status}">{device.status}</span>
					</div>
					<div class="device-card__info">
						<span>Class: {device.class}</span>
						<span>Tier: {device.tier}</span>
						<span>FW: {device.firmwareVersion}</span>
					</div>
					<div class="device-card__meta">
						<span>{device.activeEnrollments} active campaigns</span>
						<span>Last seen: {formatDate(device.lastSeen)}</span>
					</div>
					{#if device.sensors.length > 0}
						<div class="device-card__sensors">
							{#each device.sensors as sensor}
								<span class="sensor-tag">{sensor}</span>
							{/each}
						</div>
					{/if}
				</a>
			{/each}
		</div>
	{/if}
</div>

<style>
	.device-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
		gap: 1rem;
		margin-top: 1.5rem;
	}
	.device-card {
		display: block;
		padding: 1.25rem;
		border-radius: 0.5rem;
		background: var(--color-surface, #1a1a2e);
		border: 1px solid var(--color-border, #333);
		text-decoration: none;
		color: inherit;
		transition: border-color 0.2s;
	}
	.device-card:hover, .device-card:focus-visible {
		border-color: var(--color-primary, #60a5fa);
	}
	.device-card__header {
		display: flex;
		justify-content: space-between;
		margin-bottom: 0.75rem;
	}
	.device-card__info {
		display: flex;
		gap: 1rem;
		font-size: 0.875rem;
		margin-bottom: 0.5rem;
	}
	.device-card__meta {
		font-size: 0.8rem;
		color: var(--color-text-secondary, #9ca3af);
	}
	.device-card__sensors {
		display: flex;
		gap: 0.25rem;
		flex-wrap: wrap;
		margin-top: 0.5rem;
	}
	.sensor-tag {
		display: inline-block;
		padding: 0.125rem 0.5rem;
		border-radius: 9999px;
		background: var(--color-primary-muted, #1e3a5f);
		font-size: 0.75rem;
	}
</style>
