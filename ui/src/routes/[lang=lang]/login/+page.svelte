<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';
	import { login } from '$lib/auth/store';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);

	const lang = $derived($page.params.lang);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';
		loading = true;

		try {
			await login(email, password);
			await goto(`${base}/${lang}/researcher/`);
		} catch (err: any) {
			error = $t('auth.login_failed');
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex flex-col items-center justify-center min-h-screen px-6">
	<div class="w-full max-w-sm">
		<h1 class="text-2xl font-bold text-center mb-8" style="font-family: var(--font-display);">
			{$t('auth.login_title')}
		</h1>

		<form onsubmit={handleSubmit} class="flex flex-col gap-4">
			<div class="flex flex-col gap-1">
				<label for="email" class="text-sm text-cream-600">{$t('auth.email')}</label>
				<input
					id="email"
					type="email"
					bind:value={email}
					required
					autocomplete="email"
					class="px-4 py-3 bg-forest-900/50 border border-border-strong rounded text-cream-100 placeholder-cream-700 focus:outline-none focus:border-morpho-500 transition-colors"
				/>
			</div>

			<div class="flex flex-col gap-1">
				<label for="password" class="text-sm text-cream-600">{$t('auth.password')}</label>
				<input
					id="password"
					type="password"
					bind:value={password}
					required
					autocomplete="current-password"
					class="px-4 py-3 bg-forest-900/50 border border-border-strong rounded text-cream-100 placeholder-cream-700 focus:outline-none focus:border-morpho-500 transition-colors"
				/>
			</div>

			{#if error}
				<p class="text-dart-400 text-sm" role="alert">{error}</p>
			{/if}

			<button
				type="submit"
				disabled={loading}
				class="px-6 py-3 bg-morpho-500 hover:bg-morpho-400 disabled:opacity-50 text-cream-100 font-medium rounded transition-colors duration-200 mt-2"
			>
				{loading ? $t('common.loading') : $t('auth.login_button')}
			</button>
		</form>

		<p class="text-center text-sm text-cream-600 mt-6">
			{$t('auth.no_account')}
			<a href="{base}/{lang}/register" class="text-morpho-400 hover:text-morpho-300 underline">
				{$t('auth.register_link')}
			</a>
		</p>
	</div>
</div>
