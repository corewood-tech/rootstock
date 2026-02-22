<script lang="ts">
	import { t } from '$lib/i18n';
</script>

<div class="flex flex-col items-center justify-center min-h-[60vh] px-6">
	<h1 class="text-3xl font-bold mb-4" style="font-family: var(--font-display);">
		{$t('researcher.welcome')}
	</h1>
	<p class="text-cream-600 mb-6">{$t('researcher.no_campaigns')}</p>
	<button
		disabled
		class="px-6 py-3 bg-morpho-500/50 text-cream-100/50 font-medium rounded cursor-not-allowed"
	>
		{$t('researcher.create_first')}
	</button>
</div>
