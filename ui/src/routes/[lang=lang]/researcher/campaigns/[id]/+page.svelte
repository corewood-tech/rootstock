<script lang="ts">
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';
	import { campaignService } from '$lib/api/clients';
	import type { CampaignProto, GetCampaignDashboardResponse } from '$lib/api/gen/rootstock/v1/rootstock_pb';

	const lang = $derived($page.params.lang);
	const campaignId = $derived($page.params.id);

	let campaign = $state<CampaignProto | null>(null);
	let dashboard = $state<GetCampaignDashboardResponse | null>(null);
	let loading = $state(true);
	let error = $state('');
	let publishing = $state(false);
	let publishError = $state('');

	$effect(() => {
		loadCampaign();
	});

	async function loadCampaign() {
		loading = true;
		error = '';
		try {
			const resp = await campaignService.listCampaigns({});
			campaign = resp.campaigns.find((c) => c.id === campaignId) ?? null;
			if (!campaign) {
				error = 'Campaign not found';
				return;
			}
			if (campaign.status !== 'draft') {
				await loadDashboard();
			}
		} catch (e: any) {
			error = e.message || $t('campaign.load_error');
		} finally {
			loading = false;
		}
	}

	async function loadDashboard() {
		try {
			dashboard = await campaignService.getCampaignDashboard({ campaignId });
		} catch {
			// Dashboard may not be available yet — not a fatal error
		}
	}

	async function handlePublish() {
		publishing = true;
		publishError = '';
		try {
			await campaignService.publishCampaign({ campaignId });
			await loadCampaign();
		} catch (e: any) {
			publishError = e.message || 'Publish failed';
		} finally {
			publishing = false;
		}
	}

	function formatDate(dateStr?: string): string {
		if (!dateStr) return '—';
		try {
			return new Date(dateStr).toLocaleDateString();
		} catch {
			return dateStr;
		}
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
			<a href="{base}/{lang}/researcher" class="btn btn--secondary mt-4">
				{$t('common.back')}
			</a>
		</div>
	{:else if campaign}
		<header class="campaign-detail__header">
			<div>
				<a href="{base}/{lang}/researcher" class="back-link">&larr; {$t('campaign.list_title')}</a>
				<h1 class="heading heading--lg">Campaign {campaign.id.slice(0, 8)}</h1>
			</div>
			<span class="status-badge status-badge--{campaign.status}">{campaign.status}</span>
		</header>

		<section class="campaign-detail__section" aria-label="Campaign window">
			<h2 class="heading heading--md">{$t('campaign.window')}</h2>
			<p>{formatDate(campaign.windowStart)} — {formatDate(campaign.windowEnd)}</p>
		</section>

		<section class="campaign-detail__section" aria-label="Campaign metadata">
			<h2 class="heading heading--md">Details</h2>
			<dl class="detail-list">
				<dt>Created</dt>
				<dd>{formatDate(campaign.createdAt)}</dd>
				<dt>Organization</dt>
				<dd>{campaign.orgId || '—'}</dd>
			</dl>
		</section>

		{#if campaign.status === 'draft'}
			<section class="campaign-detail__section">
				<button onclick={handlePublish} class="btn btn--primary" disabled={publishing}>
					{publishing ? 'Publishing...' : 'Publish Campaign'}
				</button>
				{#if publishError}
					<p class="form-error mt-2">{publishError}</p>
				{/if}
			</section>
		{/if}

		{#if dashboard}
			<section class="campaign-detail__section" aria-label="Dashboard metrics">
				<h2 class="heading heading--md">Dashboard</h2>
				<div class="stats-row">
					<div class="stat-card">
						<span class="stat-card__value">{dashboard.acceptedCount}</span>
						<span class="stat-card__label">Accepted Readings</span>
					</div>
					<div class="stat-card">
						<span class="stat-card__value">{dashboard.quarantineCount}</span>
						<span class="stat-card__label">Quarantined Readings</span>
					</div>
				</div>
			</section>
		{/if}
	{/if}
</div>

<style>
	.campaign-detail__header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 2rem;
	}
	.back-link {
		display: inline-block;
		margin-bottom: 0.5rem;
		color: var(--color-text-secondary, #aaa);
		text-decoration: none;
	}
	.back-link:hover {
		color: var(--color-text, #fff);
	}
	.campaign-detail__section {
		margin-bottom: 2rem;
	}
	.detail-list {
		display: grid;
		grid-template-columns: auto 1fr;
		gap: 0.5rem 1.5rem;
	}
	.detail-list dt {
		color: var(--color-text-secondary, #aaa);
	}
	.stats-row {
		display: flex;
		gap: 1.5rem;
	}
	.stat-card {
		background: var(--color-surface, #1a1a2e);
		border: 1px solid var(--color-border, #333);
		border-radius: 0.75rem;
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		align-items: center;
		min-width: 160px;
	}
	.stat-card__value {
		font-size: 2rem;
		font-weight: 700;
	}
	.stat-card__label {
		color: var(--color-text-secondary, #aaa);
		font-size: 0.875rem;
		margin-top: 0.25rem;
	}
</style>
