<script lang="ts">
	// Graph node: 0x37 (ContributionPage)
	// Implements: FR-034 (Contribution Score)
	import { t } from '$lib/i18n';
	import { scitizenService } from '$lib/api/clients';
	import type { GetContributionsResponse } from '$lib/api/gen/rootstock/v1/rootstock_pb';

	let data = $state<GetContributionsResponse | null>(null);
	let loading = $state(true);
	let error = $state('');

	$effect(() => {
		loadContributions();
	});

	async function loadContributions() {
		loading = true;
		try {
			data = await scitizenService.getContributions({});
		} catch (e: any) {
			error = e.message || 'Failed to load contributions';
		} finally {
			loading = false;
		}
	}

	function formatDate(dateStr?: string): string {
		if (!dateStr) return '—';
		try { return new Date(dateStr).toLocaleDateString(); } catch { return dateStr; }
	}
</script>

<div class="contributions-page">
	<h1 class="heading heading--lg">My Contributions</h1>

	{#if loading}
		<div class="empty-state" role="status">
			<p class="text-secondary">{$t('common.loading')}</p>
		</div>
	{:else if error}
		<div class="empty-state">
			<p class="form-error">{error}</p>
		</div>
	{:else if data}
		<section class="score-section" aria-label="Contribution score">
			<div class="score-card">
				<span class="score-card__value">{data.contributionScore.toFixed(1)}</span>
				<span class="score-card__label">Contribution Score</span>
			</div>
		</section>

		{#if data.badges.length > 0}
			<section aria-label="Badges">
				<h2 class="heading heading--md">Badges</h2>
				<div class="badge-list">
					{#each data.badges as badge (badge.id)}
						<div class="badge-item">
							<span>{badge.badgeType}</span>
							<span class="text-secondary">{formatDate(badge.awardedAt)}</span>
						</div>
					{/each}
				</div>
			</section>
		{/if}

		{#if data.histories.length > 0}
			<section aria-label="Reading history">
				<h2 class="heading heading--md">Reading History</h2>
				<table class="data-table" role="table">
					<thead>
						<tr><th>Device</th><th>Campaign</th><th>Total</th><th>Accepted</th><th>Rejected</th></tr>
					</thead>
					<tbody>
						{#each data.histories as h}
							<tr>
								<td>{h.deviceId.slice(0, 8)}</td>
								<td>{h.campaignId.slice(0, 8)}</td>
								<td>{h.totalReadings}</td>
								<td>{h.accepted}</td>
								<td>{h.rejected}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</section>
		{:else}
			<div class="empty-state">
				<p>No readings submitted yet.</p>
			</div>
		{/if}
	{/if}
</div>

<style>
	.score-section { margin-bottom: 2rem; }
	.score-card {
		display: inline-flex;
		flex-direction: column;
		align-items: center;
		padding: 2rem 3rem;
		border-radius: 0.75rem;
		background: var(--color-surface, #1a1a2e);
		border: 1px solid var(--color-primary, #60a5fa);
	}
	.score-card__value {
		font-size: 3rem;
		font-weight: 700;
		color: var(--color-primary, #60a5fa);
	}
	.score-card__label {
		font-size: 0.875rem;
		color: var(--color-text-secondary, #9ca3af);
	}
	.badge-list { display: flex; gap: 0.75rem; flex-wrap: wrap; margin-bottom: 2rem; }
	.badge-item {
		display: flex;
		flex-direction: column;
		padding: 0.75rem 1rem;
		border-radius: 0.5rem;
		background: var(--color-surface, #1a1a2e);
		border: 1px solid var(--color-border, #333);
	}
	.data-table { width: 100%; border-collapse: collapse; }
	.data-table th, .data-table td {
		padding: 0.5rem 1rem;
		text-align: left;
		border-bottom: 1px solid var(--color-border, #333);
	}
</style>
