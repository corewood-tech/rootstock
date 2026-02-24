<script lang="ts">
	// Graph node: 0x41 (CampaignDetailPage)
	// Implements: FR-082 (Campaign Detail View)
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';
	import { scitizenService } from '$lib/api/clients';
	import type { GetCampaignDetailResponse } from '$lib/api/gen/rootstock/v1/rootstock_pb';

	const campaignId = $derived($page.params.id);

	let detail = $state<GetCampaignDetailResponse | null>(null);
	let loading = $state(true);
	let error = $state('');
	let enrolling = $state(false);
	let enrollResult = $state('');
	let showConsent = $state(false);

	$effect(() => {
		loadDetail();
	});

	async function loadDetail() {
		loading = true;
		error = '';
		try {
			detail = await scitizenService.getCampaignDetail({ campaignId });
		} catch (e: any) {
			error = e.message || 'Failed to load campaign';
		} finally {
			loading = false;
		}
	}

	async function handleEnroll(deviceId: string) {
		enrolling = true;
		enrollResult = '';
		try {
			const resp = await scitizenService.enrollDevice({
				deviceId,
				campaignId,
				consent: { version: '1.0', scope: 'data_collection' },
			});
			enrollResult = resp.enrolled ? 'Enrolled successfully!' : resp.reason;
			showConsent = false;
		} catch (e: any) {
			enrollResult = e.message || 'Enrollment failed';
		} finally {
			enrolling = false;
		}
	}

	function formatDate(dateStr?: string): string {
		if (!dateStr) return '—';
		try { return new Date(dateStr).toLocaleDateString(); } catch { return dateStr; }
	}
</script>

<div class="campaign-detail">
	{#if loading}
		<div class="empty-state" role="status">
			<p class="text-secondary">{$t('common.loading')}</p>
		</div>
	{:else if error}
		<div class="empty-state">
			<p class="form-error">{error}</p>
		</div>
	{:else if detail}
		<header class="campaign-detail__header">
			<h1 class="heading heading--lg">Campaign {detail.campaignId.slice(0, 8)}</h1>
			<span class="status-badge status-badge--{detail.status}">{detail.status}</span>
		</header>

		<section class="campaign-detail__section" aria-label="Campaign window">
			<h2 class="heading heading--md">Time Window</h2>
			<p>{formatDate(detail.windowStart)} — {formatDate(detail.windowEnd)}</p>
		</section>

		<section class="campaign-detail__section" aria-label="Campaign statistics">
			<h2 class="heading heading--md">Statistics</h2>
			<div class="stats-row">
				<span>{detail.enrollmentCount} enrolled devices</span>
				<span>{detail.progressPercent.toFixed(1)}% progress</span>
			</div>
		</section>

		{#if detail.parameters.length > 0}
			<section class="campaign-detail__section" aria-label="Parameters">
				<h2 class="heading heading--md">Parameters</h2>
				<table class="data-table" role="table">
					<thead>
						<tr><th>Name</th><th>Unit</th><th>Range</th></tr>
					</thead>
					<tbody>
						{#each detail.parameters as param}
							<tr>
								<td>{param.name}</td>
								<td>{param.unit}</td>
								<td>{param.minRange ?? '—'} – {param.maxRange ?? '—'}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</section>
		{/if}

		{#if detail.eligibility.length > 0}
			<section class="campaign-detail__section" aria-label="Eligibility">
				<h2 class="heading heading--md">Eligibility</h2>
				{#each detail.eligibility as elig}
					<div class="elig-card">
						<span>Class: {elig.deviceClass}</span>
						<span>Min Tier: {elig.tier}</span>
						{#if elig.requiredSensors.length > 0}
							<span>Sensors: {elig.requiredSensors.join(', ')}</span>
						{/if}
					</div>
				{/each}
			</section>
		{/if}

		{#if detail.status === 'published'}
			<section class="campaign-detail__section">
				<button onclick={() => showConsent = true} class="btn btn--primary" disabled={enrolling}>
					{enrolling ? 'Enrolling...' : 'Enroll a Device'}
				</button>
				{#if enrollResult}
					<p class="mt-2">{enrollResult}</p>
				{/if}
			</section>
		{/if}

		{#if showConsent}
			<div class="modal-overlay" role="dialog" aria-modal="true" aria-label="Consent">
				<div class="modal">
					<h2 class="heading heading--md">Consent Required</h2>
					<p>By enrolling, you consent to sharing sensor data from your device for this campaign.</p>
					<div class="modal__actions">
						<button onclick={() => showConsent = false} class="btn btn--ghost">Cancel</button>
						<button onclick={() => handleEnroll('default-device')} class="btn btn--primary">Accept & Enroll</button>
					</div>
				</div>
			</div>
		{/if}
	{/if}
</div>

<style>
	.campaign-detail__header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
	}
	.campaign-detail__section {
		margin-bottom: 2rem;
	}
	.stats-row {
		display: flex;
		gap: 2rem;
	}
	.data-table {
		width: 100%;
		border-collapse: collapse;
	}
	.data-table th, .data-table td {
		padding: 0.5rem 1rem;
		text-align: left;
		border-bottom: 1px solid var(--color-border, #333);
	}
	.elig-card {
		display: flex;
		gap: 1rem;
		padding: 0.75rem;
		background: var(--color-surface, #1a1a2e);
		border-radius: 0.5rem;
		margin-bottom: 0.5rem;
	}
	.modal-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.7);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 50;
	}
	.modal {
		background: var(--color-bg, #0f0f23);
		border: 1px solid var(--color-border, #333);
		border-radius: 0.75rem;
		padding: 2rem;
		max-width: 480px;
		width: 90%;
	}
	.modal__actions {
		display: flex;
		gap: 1rem;
		justify-content: flex-end;
		margin-top: 1.5rem;
	}
</style>
