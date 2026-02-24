<script lang="ts">
	// Graph node: 0x3a (ScitizenDashboardPage)
	// Implements: US-003 (Clear Feedback), FR-040 (Contribution Dashboard)
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';
	import { getScitizenState } from '$lib/state/scitizen.svelte';

	const lang = $derived($page.params.lang);
	const state = getScitizenState();

	$effect(() => {
		state.loadDashboard();
		state.loadOnboarding();
	});

	function formatDate(dateStr?: string): string {
		if (!dateStr) return '—';
		try { return new Date(dateStr).toLocaleDateString(); } catch { return dateStr; }
	}
</script>

<div class="dashboard">
	{#if state.loading}
		<div class="empty-state" role="status">
			<p class="text-secondary">{$t('common.loading')}</p>
		</div>
	{:else if state.error}
		<div class="empty-state">
			<p class="form-error">{state.error}</p>
			<button onclick={() => state.loadDashboard()} class="btn btn--secondary mt-4">Retry</button>
		</div>
	{:else if state.dashboard}
		{@const d = state.dashboard}

		{#if state.onboarding?.state && (!state.onboarding.state.deviceRegistered || !state.onboarding.state.campaignEnrolled)}
			<section class="onboarding-banner" aria-label="Onboarding progress">
				<h2 class="heading heading--md">Get Started</h2>
				<ul class="onboarding-steps">
					<li class:completed={state.onboarding.state.tosAccepted}>Accept Terms of Service</li>
					<li class:completed={state.onboarding.state.deviceRegistered}>Register a device</li>
					<li class:completed={state.onboarding.state.campaignEnrolled}>Enroll in a campaign</li>
					<li class:completed={state.onboarding.state.firstReadingSubmitted}>Submit first reading</li>
				</ul>
			</section>
		{/if}

		<section class="stats-grid" aria-label="Dashboard statistics">
			<div class="stat-card">
				<span class="stat-card__value">{d.activeEnrollments}</span>
				<span class="stat-card__label">Active Campaigns</span>
			</div>
			<div class="stat-card">
				<span class="stat-card__value">{d.totalReadings}</span>
				<span class="stat-card__label">Total Readings</span>
			</div>
			<div class="stat-card">
				<span class="stat-card__value">{d.acceptedReadings}</span>
				<span class="stat-card__label">Accepted</span>
			</div>
			<div class="stat-card">
				<span class="stat-card__value">{d.contributionScore.toFixed(1)}</span>
				<span class="stat-card__label">Score</span>
			</div>
		</section>

		{#if d.badges.length > 0}
			<section aria-label="Badges">
				<h2 class="heading heading--md">Badges</h2>
				<div class="badge-grid">
					{#each d.badges as badge (badge.id)}
						<div class="badge-card">
							<span class="badge-card__type">{badge.badgeType}</span>
							<span class="badge-card__date">{formatDate(badge.awardedAt)}</span>
						</div>
					{/each}
				</div>
			</section>
		{/if}

		{#if d.enrollments.length > 0}
			<section aria-label="Active enrollments">
				<h2 class="heading heading--md">Enrollments</h2>
				<div class="enrollment-list">
					{#each d.enrollments as enrollment (enrollment.id)}
						<div class="enrollment-card">
							<span class="enrollment-card__campaign">Campaign: {enrollment.campaignId.slice(0, 8)}</span>
							<span class="enrollment-card__device">Device: {enrollment.deviceId.slice(0, 8)}</span>
							<span class="status-badge status-badge--{enrollment.status}">{enrollment.status}</span>
						</div>
					{/each}
				</div>
			</section>
		{:else}
			<div class="empty-state">
				<p>No active enrollments.</p>
				<a href="{base}/{lang}/scitizen/campaigns" class="btn btn--primary">Browse Campaigns</a>
			</div>
		{/if}
	{/if}
</div>

<style>
	.stats-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
		gap: 1rem;
		margin-bottom: 2rem;
	}
	.stat-card {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 1.5rem;
		border-radius: 0.5rem;
		background: var(--color-surface, #1a1a2e);
		border: 1px solid var(--color-border, #333);
	}
	.stat-card__value {
		font-size: 2rem;
		font-weight: 700;
		color: var(--color-primary, #60a5fa);
	}
	.stat-card__label {
		font-size: 0.875rem;
		color: var(--color-text-secondary, #9ca3af);
		margin-top: 0.25rem;
	}
	.onboarding-banner {
		padding: 1.5rem;
		border-radius: 0.5rem;
		background: var(--color-surface, #1a1a2e);
		border: 1px solid var(--color-primary, #60a5fa);
		margin-bottom: 2rem;
	}
	.onboarding-steps {
		list-style: none;
		padding: 0;
		margin: 1rem 0 0;
	}
	.onboarding-steps li {
		padding: 0.5rem 0;
		color: var(--color-text-secondary, #9ca3af);
	}
	.onboarding-steps li.completed {
		color: var(--color-success, #22c55e);
		text-decoration: line-through;
	}
	.badge-grid {
		display: flex;
		gap: 0.75rem;
		flex-wrap: wrap;
		margin-bottom: 2rem;
	}
	.badge-card {
		display: flex;
		flex-direction: column;
		padding: 0.75rem 1rem;
		border-radius: 0.5rem;
		background: var(--color-surface, #1a1a2e);
		border: 1px solid var(--color-border, #333);
	}
	.enrollment-card {
		display: flex;
		gap: 1rem;
		align-items: center;
		padding: 0.75rem 1rem;
		border-radius: 0.5rem;
		background: var(--color-surface, #1a1a2e);
		border: 1px solid var(--color-border, #333);
		margin-bottom: 0.5rem;
	}
</style>
