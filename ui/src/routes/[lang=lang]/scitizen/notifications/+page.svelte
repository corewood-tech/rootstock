<script lang="ts">
	// Graph node: 0x3b (NotificationPage)
	// Implements: FR-114 (Notification History), FR-051 (Notification System)
	import { t } from '$lib/i18n';
	import { scitizenService } from '$lib/api/clients';
	import { getNotificationState } from '$lib/state/notifications.svelte';
	import type { NotificationProto } from '$lib/api/gen/rootstock/v1/rootstock_pb';

	const notifState = getNotificationState();

	let typeFilter = $state('');

	$effect(() => {
		notifState.load(typeFilter || undefined);
	});

	function formatDate(dateStr?: string): string {
		if (!dateStr) return '—';
		try { return new Date(dateStr).toLocaleString(); } catch { return dateStr; }
	}
</script>

<div class="notifications-page">
	<div class="notifications-page__header">
		<h1 class="heading heading--lg">Notifications</h1>
		<span class="unread-count">{notifState.unreadCount} unread</span>
	</div>

	<div class="filter-bar">
		<label for="type-filter" class="visually-hidden">Filter by type</label>
		<select id="type-filter" bind:value={typeFilter} class="input">
			<option value="">All types</option>
			<option value="enrollment">Enrollment</option>
			<option value="badge">Badge</option>
			<option value="campaign">Campaign</option>
			<option value="device">Device</option>
		</select>
	</div>

	{#if notifState.loading}
		<div class="empty-state" role="status">
			<p class="text-secondary">{$t('common.loading')}</p>
		</div>
	{:else if notifState.notifications.length === 0}
		<div class="empty-state">
			<p>No notifications.</p>
		</div>
	{:else}
		<ul class="notification-list" role="list">
			{#each notifState.notifications as notif (notif.id)}
				<li class="notification-item" class:notification-item--unread={!notif.read}>
					<div class="notification-item__header">
						<span class="notification-item__type">{notif.type}</span>
						<time class="notification-item__time">{formatDate(notif.createdAt)}</time>
					</div>
					<p class="notification-item__message">{notif.message}</p>
					{#if notif.resourceLink}
						<a href={notif.resourceLink} class="notification-item__link">View details</a>
					{/if}
				</li>
			{/each}
		</ul>
	{/if}
</div>

<style>
	.notifications-page__header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
	}
	.unread-count {
		font-size: 0.875rem;
		color: var(--color-primary, #60a5fa);
	}
	.filter-bar { margin-bottom: 1.5rem; }
	.notification-list { list-style: none; padding: 0; }
	.notification-item {
		padding: 1rem;
		border-bottom: 1px solid var(--color-border, #333);
	}
	.notification-item--unread {
		border-left: 3px solid var(--color-primary, #60a5fa);
		background: var(--color-surface, #1a1a2e);
	}
	.notification-item__header {
		display: flex;
		justify-content: space-between;
		margin-bottom: 0.25rem;
	}
	.notification-item__type {
		font-size: 0.75rem;
		text-transform: uppercase;
		color: var(--color-text-secondary, #9ca3af);
	}
	.notification-item__time {
		font-size: 0.75rem;
		color: var(--color-text-secondary, #9ca3af);
	}
	.notification-item__link {
		font-size: 0.875rem;
		color: var(--color-primary, #60a5fa);
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
