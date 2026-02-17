package authz

import rego.v1

# Public endpoints — accessible without authentication
public_methods := {
	"/rootstock.v1.HealthService/Check",
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

# Rule: Authenticated user accessing non-public endpoint
decision := {
	"allow": true,
	"reason": "authenticated",
} if {
	input.session_user_id != ""
	not input.method in public_methods
}
