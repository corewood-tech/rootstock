<script lang="ts">
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';
	import { isAuthenticated, currentUser, authLoading, logout } from '$lib/auth/store';
	import { requireAuth } from '$lib/auth/guard';
	import { goto } from '$app/navigation';

	let { children } = $props();

	const lang = $derived($page.params.lang);

	let guardChecked = $state(false);

	$effect(() => {
		requireAuth(lang).then((ok) => {
			guardChecked = ok;
		});
	});

	async function handleLogout() {
		await logout();
		await goto(`${base}/${lang}/`);
	}
</script>

{#if $authLoading || !guardChecked}
	<div class="flex items-center justify-center min-h-screen">
		<p class="text-cream-600">{$t('common.loading')}</p>
	</div>
{:else}
	<div class="min-h-screen flex flex-col">
		<header class="flex items-center justify-between px-6 py-4 border-b border-border">
			<div class="flex items-center gap-6">
				<a href="{base}/{lang}/researcher/" class="flex items-center gap-2">
					<img src="{base}/corewood_symbol_transparent_ON-DARK.png" alt="Rootstock" class="h-8 w-auto" />
					<span class="text-sm font-bold tracking-wider" style="font-family: var(--font-display);">ROOTSTOCK</span>
				</a>
				<nav class="flex items-center gap-4">
					<span class="text-sm text-cream-700 cursor-not-allowed">{$t('nav.campaigns')}</span>
				</nav>
			</div>
			<div class="flex items-center gap-4">
				{#if $currentUser}
					<span class="text-sm text-cream-600">{$currentUser.userType}</span>
				{/if}
				<button
					onclick={handleLogout}
					class="text-sm text-cream-600 hover:text-cream-100 transition-colors"
				>
					{$t('nav.logout')}
				</button>
			</div>
		</header>

		<div class="flex-grow">
			{@render children()}
		</div>
	</div>
{/if}
