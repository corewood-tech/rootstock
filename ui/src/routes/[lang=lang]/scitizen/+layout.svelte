<script lang="ts">
	// Graph node: 0x36 (ScitizenLayout)
	// Implements: FR-091 (Accessibility), LF-003 (Responsive Design)
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';
	import { isAuthenticated, currentUser, authLoading, logout } from '$lib/auth/store';
	import { requireAuth } from '$lib/auth/guard';
	import { goto } from '$app/navigation';
	import { createScitizenState } from '$lib/state/scitizen.svelte';
	import { createNotificationState } from '$lib/state/notifications.svelte';

	let { children } = $props();

	const lang = $derived($page.params.lang);

	let guardChecked = $state(false);
	let mobileNavOpen = $state(false);

	const scitizenState = createScitizenState();
	const notificationState = createNotificationState();

	$effect(() => {
		requireAuth(lang).then((ok) => {
			guardChecked = ok;
			if (ok) {
				notificationState.load();
			}
		});
	});

	async function handleLogout() {
		await logout();
		await goto(`${base}/${lang}/`);
	}

	const navItems = $derived([
		{ href: `${base}/${lang}/scitizen/`, label: 'Dashboard', icon: '📊' },
		{ href: `${base}/${lang}/scitizen/campaigns`, label: 'Campaigns', icon: '🔬' },
		{ href: `${base}/${lang}/scitizen/devices`, label: 'Devices', icon: '📡' },
		{ href: `${base}/${lang}/scitizen/contributions`, label: 'Contributions', icon: '📈' },
		{ href: `${base}/${lang}/scitizen/notifications`, label: 'Notifications', icon: '🔔' },
	]);
</script>

{#if $authLoading || !guardChecked}
	<div class="loading-screen" role="status" aria-label="Loading">
		<p>{$t('common.loading')}</p>
	</div>
{:else}
	<div class="page-shell">
		<header class="app-header" role="banner">
			<div class="app-header__left">
				<a href="{base}/{lang}/scitizen/" class="app-header__brand">
					<img src="{base}/corewood_symbol_transparent_ON-DARK.png" alt="Rootstock" class="app-header__brand-logo" />
					<span class="app-header__brand-name">ROOTSTOCK</span>
				</a>
				<nav class="app-header__nav" aria-label="Main navigation">
					{#each navItems as item}
						<a
							href={item.href}
							class="app-header__nav-item"
							aria-current={$page.url.pathname === item.href ? 'page' : undefined}
						>
							{item.label}
							{#if item.label === 'Notifications' && notificationState.unreadCount > 0}
								<span class="notification-badge" aria-label="{notificationState.unreadCount} unread">
									{notificationState.unreadCount}
								</span>
							{/if}
						</a>
					{/each}
				</nav>
			</div>

			<button
				class="mobile-nav-toggle"
				onclick={() => mobileNavOpen = true}
				aria-label="Open navigation"
				aria-expanded={mobileNavOpen}
			>
				&#9776;
			</button>

			<div class="app-header__right">
				{#if $currentUser}
					<span class="app-header__user-type">{$currentUser.userType}</span>
				{/if}
				<button onclick={handleLogout} class="btn btn--ghost">
					{$t('nav.logout')}
				</button>
			</div>
		</header>

		{#if mobileNavOpen}
			<nav class="mobile-nav mobile-nav--open" aria-label="Mobile navigation">
				<button class="mobile-nav__close" onclick={() => mobileNavOpen = false} aria-label="Close navigation">
					&#10005;
				</button>
				{#each navItems as item}
					<a href={item.href} class="mobile-nav__link" onclick={() => mobileNavOpen = false}>
						{item.label}
					</a>
				{/each}
				<button class="mobile-nav__link btn--ghost" onclick={() => { mobileNavOpen = false; handleLogout(); }}>
					{$t('nav.logout')}
				</button>
			</nav>
		{/if}

		<main class="page-content" role="main">
			{@render children()}
		</main>
	</div>
{/if}

<style>
	.notification-badge {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		min-width: 1.25rem;
		height: 1.25rem;
		padding: 0 0.25rem;
		border-radius: 9999px;
		background: var(--color-error, #ef4444);
		color: white;
		font-size: 0.75rem;
		font-weight: 600;
		margin-left: 0.25rem;
	}
</style>
