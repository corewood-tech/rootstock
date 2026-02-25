package authz

import rego.v1

# ABAC policy: principal attributes + method. Default deny.
# input.principal.role: "researcher" | "scitizen" | "both" | ""
# input.session_user_id: Zitadel user ID or ""
# input.method: Connect RPC procedure path

# --- Role sets ---

# Roles that grant researcher access
researcher_roles := {"researcher", "both"}

# Roles that grant scitizen access
scitizen_roles := {"scitizen", "both"}

# --- Endpoint sets ---

# Public endpoints — accessible without authentication
public_methods := {
	"/rootstock.v1.HealthService/Check",
	"/rootstock.v1.UserService/Login",
	"/rootstock.v1.UserService/RegisterResearcher",
	"/rootstock.v1.UserService/VerifyEmail",
	"/rootstock.v1.ScitizenService/RegisterScitizen",
}

# Endpoints accessible to any authenticated user regardless of role
any_authed_methods := {
	"/rootstock.v1.UserService/GetMe",
	"/rootstock.v1.UserService/RegisterUser",
	"/rootstock.v1.UserService/Logout",
	"/rootstock.v1.UserService/UpdateUserType",
	"/rootstock.v1.NotificationService/ListNotifications",
	"/rootstock.v1.NotificationService/MarkRead",
	"/rootstock.v1.NotificationService/GetPreferences",
	"/rootstock.v1.NotificationService/UpdatePreferences",
}

# Scitizen endpoints — require scitizen or both role
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

# Researcher endpoints — require researcher or both role
researcher_methods := {
	"/rootstock.v1.CampaignService/CreateCampaign",
	"/rootstock.v1.CampaignService/PublishCampaign",
	"/rootstock.v1.CampaignService/ListCampaigns",
	"/rootstock.v1.CampaignService/GetCampaignDashboard",
	"/rootstock.v1.CampaignService/ExportCampaignData",
	"/rootstock.v1.OrgService/CreateOrg",
	"/rootstock.v1.OrgService/NestOrg",
	"/rootstock.v1.OrgService/DefineRole",
	"/rootstock.v1.OrgService/AssignRole",
	"/rootstock.v1.OrgService/InviteUser",
	"/rootstock.v1.DeviceService/GetDevice",
	"/rootstock.v1.DeviceService/RevokeDevice",
	"/rootstock.v1.DeviceService/ReinstateDevice",
	"/rootstock.v1.DeviceService/EnrollInCampaign",
	"/rootstock.v1.AdminService/SuspendByClass",
}

# --- Decision rules ---

# Default deny
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

# Rule: Any authenticated user endpoints
decision := {
	"allow": true,
	"reason": "any_authenticated",
} if {
	input.session_user_id != ""
	input.method in any_authed_methods
}

# Rule: Scitizen endpoints — principal must have scitizen role
decision := {
	"allow": true,
	"reason": "scitizen_access",
} if {
	input.session_user_id != ""
	input.principal.role in scitizen_roles
	input.method in scitizen_methods
}

# Rule: Researcher endpoints — principal must have researcher role
decision := {
	"allow": true,
	"reason": "researcher_access",
} if {
	input.session_user_id != ""
	input.principal.role in researcher_roles
	input.method in researcher_methods
}
