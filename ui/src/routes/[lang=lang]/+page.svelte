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

			// Decode protobuf field 1 (string): tag 0x0a, varint length, UTF-8 bytes
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

<div class="flex flex-col items-center justify-center min-h-screen gap-12 px-6">
	<div class="flex flex-col items-center gap-6 w-4/5">
		<img src="{base}/rootstock.svg" alt="Rootstock" class="w-full" />
		<img src="{base}/corewood_symbol_transparent_ON-DARK.png" alt="Corewood" class="h-20 w-auto" />
		<div class="flex flex-col items-center" style="font-family: 'Roca Two', system-ui, sans-serif;">
			<span class="text-2xl font-bold tracking-widest">ROOTSTOCK</span>
			<span class="text-sm font-light tracking-wide">By COREWOOD</span>
		</div>
	</div>

	<div class="flex flex-col items-center gap-4">
		<div class="flex gap-4">
			<a
				href="{base}/{lang}/login"
				class="px-6 py-3 bg-morpho-500 hover:bg-morpho-400 text-cream-100 font-medium rounded transition-colors duration-200"
			>
				{$t('home.login')}
			</a>
			<a
				href="{base}/{lang}/register"
				class="px-6 py-3 border border-morpho-500 hover:bg-morpho-500/10 text-morpho-400 font-medium rounded transition-colors duration-200"
			>
				{$t('home.register')}
			</a>
		</div>

		<button
			onclick={checkHealth}
			disabled={loading}
			class="px-4 py-2 text-sm text-cream-600 hover:text-cream-100 transition-colors duration-200 mt-4"
			aria-live="polite"
		>
			{loading ? 'Checking...' : 'Health Check'}
		</button>

		{#if status}
			<p class="text-forest-300 font-mono text-sm" role="status">{status}</p>
		{/if}

		{#if error}
			<p class="text-dart-400 font-mono text-sm" role="alert">{error}</p>
		{/if}
	</div>
</div>
