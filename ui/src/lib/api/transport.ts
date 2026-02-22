import { createConnectTransport } from '@connectrpc/connect-web';
import type { Interceptor } from '@connectrpc/connect';
import { getSessionTokens } from '$lib/auth/store';

const authInterceptor: Interceptor = (next) => async (req) => {
	const tokens = getSessionTokens();
	if (tokens) {
		req.header.set('Authorization', `Bearer ${tokens.sessionId}|${tokens.sessionToken}`);
	}
	return next(req);
};

export const transport = createConnectTransport({
	baseUrl: '/',
	useBinaryFormat: true,
	interceptors: [authInterceptor],
});
