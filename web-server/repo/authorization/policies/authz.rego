package authz

import rego.v1

# Graph node: 0x51 (ABAC Policy Engine)
# ABAC policy: subject attributes + resource URN + action. Default deny.
# URN format: urn:rootstock:<resource-type>:<resource-id>

# Public endpoints — accessible without authentication
public_methods := {
	"/rootstock.v1.HealthService/Check",
	"/rootstock.v1.UserService/Login",
	"/rootstock.v1.UserService/RegisterResearcher",
	"/rootstock.v1.UserService/VerifyEmail",
	"/rootstock.v1.ScitizenService/RegisterScitizen",
}

# Scitizen-accessible endpoints — any authenticated scitizen or researcher
scitizen_methods := {
	"/rootstock.v1.ScitizenService/GetDashboard",
	"/rootstock.v1.ScitizenService/BrowsePublishedCampaigns",
	"/rootstock.v1.ScitizenService/GetCampaignDetail",
	"/rootstock.v1.ScitizenService/SearchCampaigns",
	"/rootstock.v1.ScitizenService/EnrollDevice",
	"/rootstock.v1.ScitizenService/WithdrawEnrollment",
	"/rootstock.v1.ScitizenService/GetDevices",
	"/rootstock.v1.ScitizenService/GetDeviceDetail",
	"/rootstock.v1.ScitizenService/GetNotifications",
	"/rootstock.v1.ScitizenService/GetContributions",
	"/rootstock.v1.ScitizenService/GetOnboardingState",
	"/rootstock.v1.ScoreService/GetContribution",
}

# Notification endpoints — any authenticated user
notification_methods := {
	"/rootstock.v1.NotificationService/ListNotifications",
	"/rootstock.v1.NotificationService/MarkRead",
	"/rootstock.v1.NotificationService/GetPreferences",
	"/rootstock.v1.NotificationService/UpdatePreferences",
}

# Main decision object — default deny
default decision := {
	"allow": false,
	"reason": "denied",
}

# Rule: Public endpoints (allow without auth)
decision := {
	"allow": true,
	"reason": "public_endpoint",
} if {
	input.method in public_methods
}

# Rule: Scitizen-accessible endpoints (any authenticated user)
decision := {
	"allow": true,
	"reason": "scitizen_access",
} if {
	input.session_user_id != ""
	input.method in scitizen_methods
}

# Rule: Notification endpoints (any authenticated user)
decision := {
	"allow": true,
	"reason": "notification_access",
} if {
	input.session_user_id != ""
	input.method in notification_methods
}

# Rule: Authenticated user accessing non-public, non-scitizen endpoint
# (researcher/admin endpoints — campaigns, orgs, devices, admin)
decision := {
	"allow": true,
	"reason": "authenticated",
} if {
	input.session_user_id != ""
	not input.method in public_methods
	not input.method in scitizen_methods
	not input.method in notification_methods
}
