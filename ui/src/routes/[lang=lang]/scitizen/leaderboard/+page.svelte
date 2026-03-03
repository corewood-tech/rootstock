<script lang="ts">
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';
	import { scitizenService } from '$lib/api/clients';
	import { currentUser } from '$lib/auth/store';
	import type { LeaderboardEntryProto } from '$lib/api/gen/rootstock/v1/rootstock_pb';

	let entries = $state<LeaderboardEntryProto[]>([]);
	let total = $state(0);
	let requester = $state<LeaderboardEntryProto | null>(null);
	let loading = $state(true);
	let error = $state('');

	let timePeriod = $state('all');
	let currentOffset = $state(0);
	const pageSize = 20;

	$effect(() => {
		loadLeaderboard();
	});

	async function loadLeaderboard() {
		loading = true;
		error = '';
		try {
			const resp = await scitizenService.getLeaderboard({
				timePeriod: timePeriod === 'all' ? undefined : timePeriod,
				limit: pageSize,
				offset: currentOffset,
			});
			entries = resp.entries;
			total = resp.total;
			requester = resp.requester ?? null;
		} catch (e: any) {
			error = e.message || 'Failed to load leaderboard';
		} finally {
			loading = false;
		}
	}

	function changeTimePeriod(period: string) {
		timePeriod = period;
		currentOffset = 0;
		loadLeaderboard();
	}

	function nextPage() {
		currentOffset += pageSize;
		loadLeaderboard();
	}

	function prevPage() {
		currentOffset = Math.max(0, currentOffset - pageSize);
		loadLeaderboard();
	}
</script>

<div class="leaderboard">
	<header class="leaderboard__header">
		<h1 class="heading heading--lg">Leaderboard</h1>
		<div class="leaderboard__filters">
			<button
				class="btn btn--sm"
				class:btn--primary={timePeriod === 'all'}
				class:btn--ghost={timePeriod !== 'all'}
				onclick={() => changeTimePeriod('all')}
			>All Time</button>
			<button
				class="btn btn--sm"
				class:btn--primary={timePeriod === 'month'}
				class:btn--ghost={timePeriod !== 'month'}
				onclick={() => changeTimePeriod('month')}
			>Month</button>
			<button
				class="btn btn--sm"
				class:btn--primary={timePeriod === 'week'}
				class:btn--ghost={timePeriod !== 'week'}
				onclick={() => changeTimePeriod('week')}
			>Week</button>
		</div>
	</header>

	{#if loading}
		<div class="empty-state" role="status">
			<p class="text-secondary">{$t('common.loading')}</p>
		</div>
	{:else if error}
		<div class="empty-state">
			<p class="form-error">{error}</p>
		</div>
	{:else if entries.length === 0}
		<div class="empty-state">
			<p class="text-secondary">No scores yet. Start contributing to climb the ranks!</p>
		</div>
	{:else}
		{#if requester}
			<div class="leaderboard__your-rank">
				Your rank: <strong>#{requester.rank}</strong> &mdash; Score: <strong>{requester.score.toFixed(1)}</strong>
			</div>
		{/if}

		<table class="data-table">
			<thead>
				<tr>
					<th>#</th>
					<th>Scitizen</th>
					<th>Score</th>
					<th>Badges</th>
					<th>Campaigns</th>
				</tr>
			</thead>
			<tbody>
				{#each entries as entry}
					<tr class:leaderboard__highlight={entry.scitizenId === $currentUser?.id}>
						<td>{entry.rank}</td>
						<td><code>{entry.scitizenId.slice(0, 8)}</code></td>
						<td>{entry.score.toFixed(1)}</td>
						<td>{entry.badgeCount}</td>
						<td>{entry.campaignCount}</td>
					</tr>
				{/each}
			</tbody>
		</table>

		<div class="leaderboard__pagination">
			<button class="btn btn--ghost btn--sm" onclick={prevPage} disabled={currentOffset === 0}>
				Previous
			</button>
			<span class="text-secondary">
				{currentOffset + 1}–{Math.min(currentOffset + pageSize, total)} of {total}
			</span>
			<button class="btn btn--ghost btn--sm" onclick={nextPage} disabled={currentOffset + pageSize >= total}>
				Next
			</button>
		</div>
	{/if}
</div>

<style>
	.leaderboard__header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
		flex-wrap: wrap;
		gap: 1rem;
	}
	.leaderboard__filters {
		display: flex;
		gap: 0.5rem;
	}
	.leaderboard__your-rank {
		background: var(--color-surface, #1a1a2e);
		border: 1px solid var(--color-primary, #6366f1);
		border-radius: 0.75rem;
		padding: 1rem 1.5rem;
		margin-bottom: 1.5rem;
	}
	.data-table {
		width: 100%;
		border-collapse: collapse;
	}
	.data-table th,
	.data-table td {
		padding: 0.5rem 1rem;
		text-align: left;
		border-bottom: 1px solid var(--color-border, #333);
	}
	.data-table th {
		color: var(--color-text-secondary, #aaa);
		font-weight: 600;
		font-size: 0.875rem;
	}
	.data-table code {
		font-family: var(--font-mono, monospace);
		font-size: 0.875rem;
	}
	.leaderboard__highlight {
		background: var(--color-surface, #1a1a2e);
		border-left: 3px solid var(--color-primary, #6366f1);
	}
	.leaderboard__pagination {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 1rem;
		margin-top: 1.5rem;
	}
</style>
