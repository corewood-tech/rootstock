import { derived, writable } from 'svelte/store';

export type Locale = 'en';

export const locale = writable<Locale>('en');
export const locales: Locale[] = ['en'];

interface Translations {
	[key: string]: any;
}

const translations = writable<Translations>({});

export async function loadTranslations(lang: Locale) {
	const translation = await import(`./locales/${lang}.json`);
	translations.set(translation.default);
}

function getNestedTranslation(obj: any, path: string): string {
	return path.split('.').reduce((current, key) => current?.[key], obj) || path;
}

export const t = derived(
	translations,
	($translations) => (key: string) => {
		return getNestedTranslation($translations, key);
	}
);
