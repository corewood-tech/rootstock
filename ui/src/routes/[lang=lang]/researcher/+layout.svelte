<script lang="ts">
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';
	import { currentUser, authLoading, logout, activeRole, setActiveRole } from '$lib/auth/store';
	import { requireAuth } from '$lib/auth/guard';
	import { goto } from '$app/navigation';

	let { children } = $props();

	const lang = $derived($page.params.lang);
	const isBothUser = $derived($currentUser?.userType === 'both');

	let guardChecked = $state(false);
	let mobileNavOpen = $state(false);

	$effect(() => {
		requireAuth(lang).then((ok) => {
			guardChecked = ok;
		});
	});

	async function handleLogout() {
		await logout();
		await goto(`${base}/${lang}/`);
	}

	async function switchToScitizen() {
		setActiveRole('scitizen');
		await goto(`${base}/${lang}/scitizen/`);
	}
</script>

{#if $authLoading || !guardChecked}
	<div class="loading-screen">
		<p>{$t('common.loading')}</p>
	</div>
{:else}
	<div class="page-shell">
		<header class="app-header">
			<div class="app-header__left">
				<a href="{base}/{lang}/researcher/" class="app-header__brand">
					<img src="{base}/corewood_symbol_transparent_ON-DARK.png" alt="Rootstock" class="app-header__brand-logo" />
					<span class="app-header__brand-name">ROOTSTOCK</span>
				</a>
				<nav class="app-header__nav">
					<a href="{base}/{lang}/researcher/" class="app-header__nav-item app-header__nav-item--active">{$t('nav.campaigns')}</a>
				</nav>
			</div>

			<button
				class="mobile-nav-toggle"
				onclick={() => mobileNavOpen = true}
				aria-label="Open navigation"
			>
				&#9776;
			</button>

			<div class="app-header__right">
				{#if isBothUser}
					<button onclick={switchToScitizen} class="btn btn--ghost btn--sm" title="Switch to Citizen Scientist view">
						{$t('nav.switch_to_scitizen')}
					</button>
				{/if}
				{#if $currentUser}
					<span class="app-header__user-type">{$activeRole ?? $currentUser.userType}</span>
				{/if}
				<button onclick={handleLogout} class="btn btn--ghost">
					{$t('nav.logout')}
				</button>
			</div>
		</header>

		{#if mobileNavOpen}
			<div class="mobile-nav mobile-nav--open">
				<button class="mobile-nav__close" onclick={() => mobileNavOpen = false} aria-label="Close navigation">
					&#10005;
				</button>
				<a href="{base}/{lang}/researcher/" class="mobile-nav__link" onclick={() => mobileNavOpen = false}>
					{$t('nav.campaigns')}
				</a>
				{#if isBothUser}
					<button class="mobile-nav__link btn--ghost" onclick={() => { mobileNavOpen = false; switchToScitizen(); }}>
						{$t('nav.switch_to_scitizen')}
					</button>
				{/if}
				<button class="mobile-nav__link btn--ghost" onclick={() => { mobileNavOpen = false; handleLogout(); }}>
					{$t('nav.logout')}
				</button>
			</div>
		{/if}

		<div class="page-content">
			{@render children()}
		</div>
	</div>
{/if}
