<script lang="ts">
	import '../../app.css';
	import { base } from '$app/paths';
	import { t } from '$lib/i18n';
	import { validateSession } from '$lib/auth/store';

	let { children, data } = $props();

	const currentLang = $derived(data.lang);

	$effect(() => {
		document.documentElement.lang = currentLang;
	});

	$effect(() => {
		validateSession();
	});
</script>

<svelte:head>
	{@html `<style>
		@font-face {
			font-family: 'Roca Two';
			src: url('${base}/fonts/roca-two-bold.woff2') format('woff2');
			font-weight: 700; font-style: normal; font-display: swap;
		}
		@font-face {
			font-family: 'Roca Two';
			src: url('${base}/fonts/roca-two-thin.woff2') format('woff2');
			font-weight: 300; font-style: normal; font-display: swap;
		}
		@font-face {
			font-family: 'HK Grotesk Pro';
			src: url('${base}/fonts/hk-grotesk-pro-regular.woff2') format('woff2');
			font-weight: 400; font-style: normal; font-display: swap;
		}
		@font-face {
			font-family: 'HK Grotesk Pro';
			src: url('${base}/fonts/hk-grotesk-pro-medium.woff2') format('woff2');
			font-weight: 500; font-style: normal; font-display: swap;
		}
		@font-face {
			font-family: 'HK Grotesk Pro';
			src: url('${base}/fonts/hk-grotesk-pro-bold.woff2') format('woff2');
			font-weight: 700; font-style: normal; font-display: swap;
		}
	</style>`}
</svelte:head>

<div class="page-shell">
	<a href="#main-content" class="sr-only">
		{$t('nav.skip_to_content')}
	</a>

	<main id="main-content" class="page-content">
		{@render children()}
	</main>
</div>
