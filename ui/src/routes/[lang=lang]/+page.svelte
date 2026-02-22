<script lang="ts">
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';

	let status = $state('');
	let error = $state('');
	let loading = $state(false);

	const lang = $derived($page.params.lang);

	async function checkHealth() {
		loading = true;
		status = '';
		error = '';

		try {
			const res = await fetch('/rootstock.v1.HealthService/Check', {
				method: 'POST',
				headers: { 'Content-Type': 'application/proto' },
				body: new Uint8Array(0)
			});

			if (!res.ok) {
				throw new Error(`${res.status} ${res.statusText}`);
			}

			const buf = new Uint8Array(await res.arrayBuffer());

			if (buf.length >= 2 && buf[0] === 0x0a) {
				const len = buf[1];
				status = new TextDecoder().decode(buf.slice(2, 2 + len));
			} else {
				status = 'OK (empty response)';
			}
		} catch (e: any) {
			error = e.message;
		} finally {
			loading = false;
		}
	}
</script>

<div class="landing">
	<div class="landing__brand">
		<img src="{base}/rootstock.svg" alt="Rootstock" class="landing__logo" />
		<img src="{base}/corewood_symbol_transparent_ON-DARK.png" alt="Corewood" class="landing__symbol" />
		<div class="landing__wordmark">
			<span class="landing__title">ROOTSTOCK</span>
			<span class="landing__subtitle">By COREWOOD</span>
		</div>
	</div>

	<div class="landing__actions">
		<div class="landing__cta-group">
			<a href="{base}/{lang}/login" class="btn btn--primary">
				{$t('home.login')}
			</a>
			<a href="{base}/{lang}/register" class="btn btn--secondary">
				{$t('home.register')}
			</a>
		</div>

		<button
			onclick={checkHealth}
			disabled={loading}
			class="btn btn--ghost mt-4"
			aria-live="polite"
		>
			{loading ? 'Checking...' : 'Health Check'}
		</button>

		{#if status}
			<p class="status-text status-text--success" role="status">{status}</p>
		{/if}

		{#if error}
			<p class="status-text status-text--error" role="alert">{error}</p>
		{/if}
	</div>
</div>
