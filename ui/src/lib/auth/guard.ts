import { get } from 'svelte/store';
import { goto } from '$app/navigation';
import { base } from '$app/paths';
import { isAuthenticated, authLoading, validateSession } from './store';

/** Validates session and redirects to login if unauthenticated. */
export async function requireAuth(lang: string): Promise<boolean> {
	const loading = get(authLoading);
	if (loading) {
		await validateSession();
	}

	if (!get(isAuthenticated)) {
		await goto(`${base}/${lang}/login`);
		return false;
	}
	return true;
}
