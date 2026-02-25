import { writable, derived } from 'svelte/store';
import { userService } from '$lib/api/clients';
import type { UserProto } from '$lib/api/gen/rootstock/v1/rootstock_pb';
import { USER_TYPE, getDashboardPath, type ActiveRole, type RegistrationRole } from './roles';

const SESSION_ID_KEY = 'rootstock_session_id';
const SESSION_TOKEN_KEY = 'rootstock_session_token';
const ACTIVE_ROLE_KEY = 'rootstock_active_role';

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

/** Active role for the current session. For "both" users, this determines which view they see. */
export const activeRole = writable<ActiveRole | null>(null);

/** Re-export for convenience. */
export { getDashboardPath, USER_TYPE, type ActiveRole, type RegistrationRole };

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

/** Resolve the active role from user_type. For "both", check sessionStorage or default to researcher. */
function resolveActiveRole(userType: string): ActiveRole {
	if (userType === USER_TYPE.BOTH) {
		const stored = sessionStorage.getItem(ACTIVE_ROLE_KEY);
		if (stored === USER_TYPE.RESEARCHER || stored === USER_TYPE.SCITIZEN) {
			return stored;
		}
		return USER_TYPE.RESEARCHER;
	}
	if (userType === USER_TYPE.SCITIZEN) return USER_TYPE.SCITIZEN;
	return USER_TYPE.RESEARCHER;
}

/** Set the active role for this session. Only meaningful for "both" users. */
export function setActiveRole(role: ActiveRole): void {
	sessionStorage.setItem(ACTIVE_ROLE_KEY, role);
	activeRole.set(role);
}

export async function login(email: string, password: string): Promise<void> {
	const resp = await userService.login({ email, password });
	saveTokens(resp.sessionId, resp.sessionToken);
	const user = resp.user ?? null;
	authState.set({
		tokens: { sessionId: resp.sessionId, sessionToken: resp.sessionToken },
		user,
		loading: false,
	});
	if (user) {
		const role = resolveActiveRole(user.userType);
		setActiveRole(role);
	}
}

export async function logout(): Promise<void> {
	try {
		await userService.logout({});
	} finally {
		clearTokens();
		sessionStorage.removeItem(ACTIVE_ROLE_KEY);
		activeRole.set(null);
		authState.set({ tokens: null, user: null, loading: false });
	}
}

/** Register a user. Returns userId — user must verify email before login. */
export async function register(
	email: string,
	password: string,
	givenName: string,
	familyName: string,
	userType: RegistrationRole,
): Promise<{ userId: string; emailVerificationSent: boolean }> {
	const resp = await userService.registerResearcher({
		email,
		password,
		givenName,
		familyName,
		userType,
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

/** Update user type via backend. */
export async function updateUserType(newType: string): Promise<void> {
	const resp = await userService.updateUserType({ userType: newType });
	if (resp.user) {
		authState.update((s) => ({ ...s, user: resp.user ?? s.user }));
		const role = resolveActiveRole(resp.user.userType);
		setActiveRole(role);
	}
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
		const user = resp.user ?? null;
		authState.set({
			tokens,
			user,
			loading: false,
		});
		if (user) {
			const role = resolveActiveRole(user.userType);
			activeRole.set(role);
			if (!sessionStorage.getItem(ACTIVE_ROLE_KEY)) {
				sessionStorage.setItem(ACTIVE_ROLE_KEY, role);
			}
		}
	} catch {
		clearTokens();
		authState.set({ tokens: null, user: null, loading: false });
	}
}
