import { writable, derived } from 'svelte/store';
import { userService } from '$lib/api/clients';
import type { UserProto } from '$lib/api/gen/rootstock/v1/rootstock_pb';

const SESSION_ID_KEY = 'rootstock_session_id';
const SESSION_TOKEN_KEY = 'rootstock_session_token';

interface AuthState {
	tokens: { sessionId: string; sessionToken: string } | null;
	user: UserProto | null;
	loading: boolean;
}

const authState = writable<AuthState>({
	tokens: null,
	user: null,
	loading: true,
});

export const isAuthenticated = derived(authState, ($s) => $s.tokens !== null && $s.user !== null);
export const currentUser = derived(authState, ($s) => $s.user);
export const authLoading = derived(authState, ($s) => $s.loading);

/** Synchronous read for the transport interceptor. */
export function getSessionTokens(): { sessionId: string; sessionToken: string } | null {
	const sessionId = localStorage.getItem(SESSION_ID_KEY);
	const sessionToken = localStorage.getItem(SESSION_TOKEN_KEY);
	if (sessionId && sessionToken) {
		return { sessionId, sessionToken };
	}
	return null;
}

function saveTokens(sessionId: string, sessionToken: string) {
	localStorage.setItem(SESSION_ID_KEY, sessionId);
	localStorage.setItem(SESSION_TOKEN_KEY, sessionToken);
}

function clearTokens() {
	localStorage.removeItem(SESSION_ID_KEY);
	localStorage.removeItem(SESSION_TOKEN_KEY);
}

export async function login(email: string, password: string): Promise<void> {
	const resp = await userService.login({ email, password });
	saveTokens(resp.sessionId, resp.sessionToken);
	authState.set({
		tokens: { sessionId: resp.sessionId, sessionToken: resp.sessionToken },
		user: resp.user ?? null,
		loading: false,
	});
}

export async function logout(): Promise<void> {
	try {
		await userService.logout({});
	} finally {
		clearTokens();
		authState.set({ tokens: null, user: null, loading: false });
	}
}

/** Register a researcher. Returns userId — user must verify email before login. */
export async function registerResearcher(
	email: string,
	password: string,
	givenName: string,
	familyName: string,
): Promise<{ userId: string; emailVerificationSent: boolean }> {
	const resp = await userService.registerResearcher({
		email,
		password,
		givenName,
		familyName,
	});
	return {
		userId: resp.userId,
		emailVerificationSent: resp.emailVerificationSent,
	};
}

/** Verify email using code from verification link. */
export async function verifyEmail(
	userId: string,
	verificationCode: string,
): Promise<boolean> {
	const resp = await userService.verifyEmail({
		userId,
		verificationCode,
	});
	return resp.verified;
}

/** Validate existing session from localStorage. Call on app mount. */
export async function validateSession(): Promise<void> {
	const tokens = getSessionTokens();
	if (!tokens) {
		authState.set({ tokens: null, user: null, loading: false });
		return;
	}

	try {
		const resp = await userService.getMe({});
		authState.set({
			tokens,
			user: resp.user ?? null,
			loading: false,
		});
	} catch {
		clearTokens();
		authState.set({ tokens: null, user: null, loading: false });
	}
}
