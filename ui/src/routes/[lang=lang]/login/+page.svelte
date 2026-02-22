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

<div class="page-centered">
	<div class="form-card">
		<h1 class="heading heading--lg heading--center">
			{$t('auth.login_title')}
		</h1>

		<form onsubmit={handleSubmit} class="form-stack">
			<div class="field">
				<label for="email" class="field__label">{$t('auth.email')}</label>
				<input
					id="email"
					type="email"
					bind:value={email}
					required
					autocomplete="email"
					class="field__input"
				/>
			</div>

			<div class="field">
				<label for="password" class="field__label">{$t('auth.password')}</label>
				<input
					id="password"
					type="password"
					bind:value={password}
					required
					autocomplete="current-password"
					class="field__input"
				/>
			</div>

			{#if error}
				<p class="form-error" role="alert">{error}</p>
			{/if}

			<button
				type="submit"
				disabled={loading}
				class="btn btn--primary mt-2"
			>
				{loading ? $t('common.loading') : $t('auth.login_button')}
			</button>
		</form>

		<p class="form-footer">
			{$t('auth.no_account')}
			<a href="{base}/{lang}/register" class="form-footer__link">
				{$t('auth.register_link')}
			</a>
		</p>
	</div>
</div>
