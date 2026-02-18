import type { LayoutLoad } from './$types';
import { loadTranslations, locale, type Locale } from '$lib/i18n';

export const ssr = false;
export const prerender = false;

export const load: LayoutLoad = async ({ params }) => {
	const lang = params.lang as Locale;
	locale.set(lang);
	await loadTranslations(lang);

	return {
		lang
	};
};
