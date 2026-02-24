// Graph node: 0x3c (ScitizenState)
// Scitizen state module: dashboard data, onboarding state. Context-based (SSR-safe).
import { getContext, setContext } from 'svelte';
import { scitizenService } from '$lib/api/clients';
import type {
	GetDashboardResponse,
	GetOnboardingStateResponse,
} from '$lib/api/gen/rootstock/v1/rootstock_pb';

const SCITIZEN_CTX = Symbol('scitizen');

export interface ScitizenState {
	dashboard: GetDashboardResponse | null;
	onboarding: GetOnboardingStateResponse | null;
	loading: boolean;
	error: string;
	loadDashboard(): Promise<void>;
	loadOnboarding(): Promise<void>;
}

export function createScitizenState(): ScitizenState {
	let dashboard = $state<GetDashboardResponse | null>(null);
	let onboarding = $state<GetOnboardingStateResponse | null>(null);
	let loading = $state(false);
	let error = $state('');

	async function loadDashboard() {
		loading = true;
		error = '';
		try {
			dashboard = await scitizenService.getDashboard({});
		} catch (e: any) {
			error = e.message || 'Failed to load dashboard';
		} finally {
			loading = false;
		}
	}

	async function loadOnboarding() {
		try {
			onboarding = await scitizenService.getOnboardingState({});
		} catch (e: any) {
			error = e.message || 'Failed to load onboarding state';
		}
	}

	const state: ScitizenState = {
		get dashboard() { return dashboard; },
		get onboarding() { return onboarding; },
		get loading() { return loading; },
		get error() { return error; },
		loadDashboard,
		loadOnboarding,
	};

	setContext(SCITIZEN_CTX, state);
	return state;
}

export function getScitizenState(): ScitizenState {
	return getContext<ScitizenState>(SCITIZEN_CTX);
}
