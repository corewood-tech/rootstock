<script lang="ts">
	// Graph node: 0x42 (BrowseCampaignsPage)
	// Implements: FR-009 (Discover Campaigns), FR-012 (Browse Campaigns), FR-088 (Search)
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';
	import { scitizenService } from '$lib/api/clients';
	import type { CampaignSummaryProto } from '$lib/api/gen/rootstock/v1/rootstock_pb';

	const lang = $derived($page.params.lang);

	let campaigns = $state<CampaignSummaryProto[]>([]);
	let total = $state(0);
	let loading = $state(true);
	let error = $state('');
	let searchQuery = $state('');
	let currentOffset = $state(0);
	const pageSize = 20;

	$effect(() => {
		loadCampaigns();
	});

	async function loadCampaigns() {
		loading = true;
		error = '';
		try {
			if (searchQuery.trim()) {
				const resp = await scitizenService.searchCampaigns({
					query: searchQuery,
					limit: pageSize,
					offset: currentOffset,
				});
				campaigns = resp.campaigns;
				total = resp.total;
			} else {
				const resp = await scitizenService.browsePublishedCampaigns({
					limit: pageSize,
					offset: currentOffset,
				});
				campaigns = resp.campaigns;
				total = resp.total;
			}
		} catch (e: any) {
			error = e.message || 'Failed to load campaigns';
		} finally {
			loading = false;
		}
	}

	function handleSearch() {
		currentOffset = 0;
		loadCampaigns();
	}

	function formatDate(dateStr?: string): string {
		if (!dateStr) return '—';
		try { return new Date(dateStr).toLocaleDateString(); } catch { return dateStr; }
	}
</script>

<div class="campaigns-page">
	<div class="campaigns-page__header">
		<h1 class="heading heading--lg">Browse Campaigns</h1>
		<form class="search-form" onsubmit={(e) => { e.preventDefault(); handleSearch(); }} role="search">
			<label for="campaign-search" class="visually-hidden">Search campaigns</label>
			<input
				id="campaign-search"
				type="search"
				bind:value={searchQuery}
				placeholder="Search campaigns..."
				class="input"
			/>
			<button type="submit" class="btn btn--secondary">Search</button>
		</form>
	</div>

	{#if loading}
		<div class="empty-state" role="status">
			<p class="text-secondary">{$t('common.loading')}</p>
		</div>
	{:else if error}
		<div class="empty-state">
			<p class="form-error">{error}</p>
			<button onclick={loadCampaigns} class="btn btn--secondary mt-4">Retry</button>
		</div>
	{:else if campaigns.length === 0}
		<div class="empty-state">
			<p>No campaigns found.</p>
		</div>
	{:else}
		<div class="campaign-grid" role="list">
			{#each campaigns as campaign (campaign.id)}
				<a
					href="{base}/{lang}/scitizen/campaigns/{campaign.id}"
					class="campaign-card"
					role="listitem"
				>
					<div class="campaign-card__top">
						<span class="campaign-card__id">{campaign.id.slice(0, 8)}</span>
						<span class="status-badge status-badge--{campaign.status}">{campaign.status}</span>
					</div>
					<div class="campaign-card__dates">
						{formatDate(campaign.windowStart)} — {formatDate(campaign.windowEnd)}
					</div>
					<div class="campaign-card__meta">
						<span>{campaign.enrollmentCount} enrolled</span>
						{#if campaign.requiredSensors.length > 0}
							<span>Sensors: {campaign.requiredSensors.join(', ')}</span>
						{/if}
					</div>
				</a>
			{/each}
		</div>

		{#if total > pageSize}
			<nav class="pagination" aria-label="Campaign pagination">
				<button
					onclick={() => { currentOffset = Math.max(0, currentOffset - pageSize); loadCampaigns(); }}
					disabled={currentOffset === 0}
					class="btn btn--ghost"
				>Previous</button>
				<span class="pagination__info">
					{currentOffset + 1}–{Math.min(currentOffset + pageSize, total)} of {total}
				</span>
				<button
					onclick={() => { currentOffset += pageSize; loadCampaigns(); }}
					disabled={currentOffset + pageSize >= total}
					class="btn btn--ghost"
				>Next</button>
			</nav>
		{/if}
	{/if}
</div>

<style>
	.campaigns-page__header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
		flex-wrap: wrap;
		gap: 1rem;
	}
	.search-form {
		display: flex;
		gap: 0.5rem;
	}
	.campaign-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 1rem;
	}
	.campaign-card {
		display: block;
		padding: 1.25rem;
		border-radius: 0.5rem;
		background: var(--color-surface, #1a1a2e);
		border: 1px solid var(--color-border, #333);
		text-decoration: none;
		color: inherit;
		transition: border-color 0.2s;
	}
	.campaign-card:hover, .campaign-card:focus-visible {
		border-color: var(--color-primary, #60a5fa);
		outline: none;
	}
	.campaign-card__top {
		display: flex;
		justify-content: space-between;
		margin-bottom: 0.5rem;
	}
	.campaign-card__meta {
		display: flex;
		gap: 1rem;
		font-size: 0.875rem;
		color: var(--color-text-secondary, #9ca3af);
	}
	.pagination {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 1rem;
		margin-top: 2rem;
	}
	.visually-hidden {
		position: absolute;
		width: 1px;
		height: 1px;
		padding: 0;
		margin: -1px;
		overflow: hidden;
		clip: rect(0, 0, 0, 0);
		white-space: nowrap;
		border: 0;
	}
</style>
