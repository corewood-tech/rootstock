<script lang="ts">
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';
	import { verifyEmail } from '$lib/auth/store';

	const lang = $derived($page.params.lang);

	let status = $state<'verifying' | 'success' | 'error'>('verifying');
	let errorMessage = $state('');

	$effect(() => {
		const url = new URL(window.location.href);
		const userId = url.searchParams.get('userId');
		const code = url.searchParams.get('code');

		if (!userId || !code) {
			status = 'error';
			errorMessage = $t('auth.verify_missing_params');
			return;
		}

		verifyEmail(userId, code)
			.then((verified) => {
				status = verified ? 'success' : 'error';
				if (!verified) {
					errorMessage = $t('auth.verify_failed');
				}
			})
			.catch(() => {
				status = 'error';
				errorMessage = $t('auth.verify_failed');
			});
	});
</script>

<div class="verify-email">
	{#if status === 'verifying'}
		<p class="text-secondary">{$t('common.loading')}</p>
	{:else if status === 'success'}
		<h1 class="heading heading--lg">{$t('auth.verify_success_title')}</h1>
		<p class="verify-email__message">{$t('auth.verify_success_message')}</p>
		<a href="{base}/{lang}/login" class="btn btn--primary">
			{$t('auth.login_button')}
		</a>
	{:else}
		<h1 class="heading heading--lg">{$t('auth.verify_error_title')}</h1>
		<p class="verify-email__message">{errorMessage}</p>
		<a href="{base}/{lang}/register" class="btn btn--secondary">
			{$t('auth.register_button')}
		</a>
	{/if}
</div>
