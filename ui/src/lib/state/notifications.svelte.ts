// Graph node: 0x39 (NotificationState)
// Notification state: unread count, notification list. Context-based.
import { getContext, setContext } from 'svelte';
import { scitizenService } from '$lib/api/clients';
import type {
	NotificationProto,
} from '$lib/api/gen/rootstock/v1/rootstock_pb';

const NOTIFICATION_CTX = Symbol('notifications');

export interface NotificationState {
	notifications: NotificationProto[];
	unreadCount: number;
	total: number;
	loading: boolean;
	load(typeFilter?: string): Promise<void>;
}

export function createNotificationState(): NotificationState {
	let notifications = $state<NotificationProto[]>([]);
	let unreadCount = $state(0);
	let total = $state(0);
	let loading = $state(false);

	async function load(typeFilter?: string) {
		loading = true;
		try {
			const resp = await scitizenService.getNotifications({
				typeFilter,
				limit: 50,
				offset: 0,
			});
			notifications = resp.notifications;
			unreadCount = resp.unreadCount;
			total = resp.total;
		} catch {
			// silently fail for badge
		} finally {
			loading = false;
		}
	}

	const state: NotificationState = {
		get notifications() { return notifications; },
		get unreadCount() { return unreadCount; },
		get total() { return total; },
		get loading() { return loading; },
		load,
	};

	setContext(NOTIFICATION_CTX, state);
	return state;
}

export function getNotificationState(): NotificationState {
	return getContext<NotificationState>(NOTIFICATION_CTX);
}
