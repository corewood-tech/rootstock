/** User types matching the DB CHECK constraint and proto enum. */
export const USER_TYPE = {
	RESEARCHER: 'researcher',
	SCITIZEN: 'scitizen',
	BOTH: 'both',
} as const;

export type UserType = typeof USER_TYPE[keyof typeof USER_TYPE];

/** Active role for a session — only the two concrete roles, never "both". */
export type ActiveRole = 'researcher' | 'scitizen';

/** Valid user types for registration. */
export type RegistrationRole = typeof USER_TYPE.RESEARCHER | typeof USER_TYPE.SCITIZEN;

/** Returns the dashboard base path segment for a given role. */
export function dashboardSegment(role: ActiveRole): string {
	return role === USER_TYPE.SCITIZEN ? 'scitizen' : 'researcher';
}

/** Returns the full dashboard path for a given role. */
export function getDashboardPath(lang: string, role: ActiveRole): string {
	return `/${lang}/${dashboardSegment(role)}/`;
}

/** Checks whether a user_type string represents a valid registration role. */
export function isRegistrationRole(value: string): value is RegistrationRole {
	return value === USER_TYPE.RESEARCHER || value === USER_TYPE.SCITIZEN;
}
