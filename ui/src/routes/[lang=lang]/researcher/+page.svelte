<script lang="ts">
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';
	import { campaignService } from '$lib/api/clients';
	import type { CampaignProto } from '$lib/api/gen/rootstock/v1/rootstock_pb';

	const lang = $derived($page.params.lang);

	let campaigns = $state<CampaignProto[]>([]);
	let loading = $state(true);
	let error = $state('');

	$effect(() => {
		loadCampaigns();
	});

	async function loadCampaigns() {
		loading = true;
		error = '';
		try {
			const resp = await campaignService.listCampaigns({});
			campaigns = resp.campaigns;
		} catch (e: any) {
			error = $t('campaign.load_error');
		} finally {
			loading = false;
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

	function statusClass(status: string): string {
		return `status-badge status-badge--${status}`;
	}
</script>

{#if loading}
	<div class="empty-state">
		<p class="text-secondary">{$t('common.loading')}</p>
	</div>
{:else if error}
	<div class="empty-state">
		<p class="form-error">{error}</p>
		<button onclick={loadCampaigns} class="btn btn--secondary mt-4">
			{$t('campaign.retry')}
		</button>
	</div>
{:else if campaigns.length === 0}
	<div class="empty-state">
		<h1 class="heading heading--xl">
			{$t('researcher.welcome')}
		</h1>
		<p class="empty-state__message">{$t('researcher.no_campaigns')}</p>
		<a href="{base}/{lang}/researcher/campaigns/new" class="btn btn--primary">
			{$t('campaign.create')}
		</a>
	</div>
{:else}
	<div class="campaign-list">
		<div class="campaign-list__header">
			<h1 class="heading heading--lg">{$t('campaign.list_title')}</h1>
			<a href="{base}/{lang}/researcher/campaigns/new" class="btn btn--primary">
				{$t('campaign.create')}
			</a>
		</div>

		{#each campaigns as campaign (campaign.id)}
			<a href="{base}/{lang}/researcher/campaigns/{campaign.id}" class="campaign-card campaign-card--link">
				<div class="campaign-card__top">
					<span class="campaign-card__id">{campaign.id.slice(0, 8)}</span>
					<span class={statusClass(campaign.status)}>{campaign.status}</span>
				</div>
				<div class="campaign-card__dates">
					{$t('campaign.window')}: {formatDate(campaign.windowStart)} — {formatDate(campaign.windowEnd)}
				</div>
				<div class="campaign-card__meta">
					<span>{$t('campaign.created')}: {formatDate(campaign.createdAt)}</span>
				</div>
			</a>
		{/each}
	</div>
{/if}
