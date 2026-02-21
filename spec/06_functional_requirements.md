# ROOTSTOCK by Corewood

## Requirements Specification — Section 6: Functional Requirements

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Functional requirements describe what the product must do — the actions, behaviors, and computations it performs. Each requirement is derived from one or more business use cases (Section 5) and includes a fit criterion that makes it testable.

> **Knowledge graph reference**: The requirements model behind this section is captured in a persistent Dgraph knowledge graph (`grapher/schema/rootstock_requirements.graphql`). Nodes are referenced by UID for traceability. Start the requirements graph with `GRAPH=reqs podman compose -f grapher/compose-grapher.yml up -d` and query at `http://localhost:18080`.

---

## 6a. Institutional Onboarding (BUC-01, scope:0x28)

### FR-001: Create Organization Tenant (0x10)

**Priority:** Must | **Originator:** Research Institution

**Description:** The platform shall allow a research institution to create an organizational tenant with a unique identifier, display name, and initial administrator.

**Rationale:** Institutions need an isolated organizational boundary to manage researchers, campaigns, and data. Multi-tenancy is foundational to platform adoption.

**Fit Criterion:** A new organization tenant is created with a unique identifier, the requesting user is assigned the admin role, and the tenant is queryable via API within 5 seconds of creation. (Scale: seconds | Worst: 10 | Plan: 5 | Best: 2)

**Constrained by:** CON-003 — No Shared Context Exists
**Cross-ref:** scope:0x28, facts:0xe

---

### FR-002: Configure Organization Hierarchy (0x12)

**Priority:** Must | **Originator:** Research Institution

**Description:** The platform shall support nested sub-organizations within a tenant, allowing institutions to model departments, labs, and divisions as a hierarchy.

**Rationale:** Institutions have departments, labs, and divisions. The hierarchy must mirror real organizational structure for authorization scoping.

**Fit Criterion:** Sub-organizations can be nested to at least 5 levels. Each sub-organization inherits parent authorization rules unless explicitly overridden. (Scale: nesting levels | Worst: 3 | Plan: 5 | Best: 10)

**Constrained by:** CON-003 — No Shared Context Exists
**Depends on:** FR-001
**Cross-ref:** scope:0x28

---

### FR-003: Define and Assign Roles (0x14)

**Priority:** Must | **Originator:** Research Institution

**Description:** The platform shall support organization-scoped roles with assignable permissions. Users may hold different roles in different organizations.

**Rationale:** Roles scope permissions to organizations. A user may hold different roles in different organizations. Permissions are granted via roles, never directly.

**Fit Criterion:** A user assigned a role in one organization does not inherit that role in sibling organizations. Role changes take effect on the next API request. (Scale: boolean | Pass/Fail)

**Depends on:** FR-001
**Cross-ref:** scope:0x28

---

### FR-004: Invite and Onboard Researchers (0xe)

**Priority:** Must | **Originator:** Research Institution

**Description:** The platform shall allow organization admins to invite researchers by email. Invited researchers authenticate via Zitadel and are automatically associated with the inviting organization.

**Rationale:** Researchers must be invited to an organization to create campaigns. Identity is delegated to Zitadel. Onboarding requires no integration with institutional systems (CON-003).

**Fit Criterion:** An invited researcher can authenticate and access their organization within 2 minutes of accepting the invitation, with zero changes to institutional systems. (Scale: minutes | Worst: 5 | Plan: 2 | Best: 1)

**Constrained by:** CON-003 — No Shared Context Exists
**Depends on:** FR-003
**Cross-ref:** scope:0x28, facts:0xc

---

### FR-046: Role Hierarchy and Inheritance (0x2723)

**Priority:** Must | **Originator:** Research Institution

**Description:** The platform shall support role inheritance along the organization hierarchy. A user with a role at a parent organization inherits that role in all child sub-organizations unless explicitly overridden at the child level. Override applies downward from the override point.

**Rationale:** FR-002 establishes nested sub-organizations and FR-003 establishes org-scoped roles, but neither specifies how roles interact with nesting. Without inheritance rules, a university admin would need manual role assignment in every department.

**Fit Criterion:** A user with admin role at a parent organization can perform admin actions in all child sub-organizations. A child sub-organization with an explicit role override for that user applies the override role, not the inherited role. Override does not affect sibling organizations. (Scale: boolean | Pass/Fail)

**Depends on:** FR-002, FR-003

---

### FR-047: Permission Model Definition (0x2725)

**Priority:** Must | **Originator:** Research Institution

**Description:** The platform shall define a discrete set of permissions that roles grant: campaign.create, campaign.publish, campaign.edit, campaign.cancel, device.suspend, device.revoke, data.export, user.invite, org.manage, role.assign. Permissions are granted via roles, never directly to users. New permissions may be added without schema changes.

**Rationale:** FR-003 says roles have assignable permissions but does not enumerate them. Without a defined permission model, OPA policies cannot be written and tested. The permission set must be explicit and extensible.

**Fit Criterion:** Every API endpoint is gated by at least one permission. A user without the required permission receives a 403 response. Adding a new permission to the system does not require a database migration. (Scale: boolean | Pass/Fail)

**Depends on:** FR-003

---

### FR-048: Organization Member Management (0x2727)

**Priority:** Must | **Originator:** Research Institution

**Description:** The platform shall allow organization admins to list members, remove members, and transfer organization ownership. Removing a member revokes their role in that organization and all child sub-organizations. Campaigns created by a removed member are retained under the organization, not deleted.

**Rationale:** Organizations need to manage their member roster. Researchers leave institutions. Campaigns are organizational assets that must survive member departure.

**Fit Criterion:** A removed member cannot access organization resources. Campaigns created by the removed member remain accessible to remaining org members. Organization ownership transfer requires confirmation from both parties. (Scale: boolean | Pass/Fail)

**Depends on:** FR-003, FR-004

---

### FR-049: Audit Trail for Role Changes (0x272b)

**Priority:** Should | **Originator:** Research Institution

**Description:** The platform shall maintain an immutable audit trail for all role and permission changes: who changed what role, for whom, in which organization, and when. Audit entries include the actor, target user, old role, new role, organization, and timestamp.

**Rationale:** Research institutions have compliance and governance requirements. Role changes affect who can access data and manage campaigns. An audit trail provides accountability and supports compliance reporting.

**Fit Criterion:** Every role assignment, removal, or change produces an audit entry with actor, target, old role, new role, organization, and timestamp. Audit entries cannot be modified or deleted. Audit log is queryable by organization, user, and time range. (Scale: boolean | Pass/Fail)

**Depends on:** FR-003

---

### FR-086: Organization Dashboard (0x4e52)

**Priority:** Should | **Originator:** Research Institution

**Description:** The platform shall provide organization administrators with an organization dashboard showing: member count and recent activity, campaigns by status (draft, published, active, completed, cancelled), aggregate data volume across all org campaigns, total enrolled devices contributing to org campaigns, and sub-organization summary if hierarchy exists.

**Rationale:** BUC-01 establishes organizational tenants with hierarchy. FR-048 provides member management. But no requirement gives the org admin an operational overview. An institution needs to see what their researchers are doing, how much data their campaigns are collecting, and whether the platform is delivering value.

**Fit Criterion:** Organization dashboard displays member count, campaign count by status, aggregate data volume, and enrolled device count. If sub-organizations exist, each sub-org shows its own summary. Dashboard data refreshed within 15 minutes. Accessible only to users with org admin role. (Scale: minutes | Worst: 30 | Plan: 15 | Best: 5)

**Depends on:** FR-001, FR-048

---

### FR-096: Invitation Management (0x4e66)

**Priority:** Must | **Originator:** Research Institution

**Description:** The platform shall allow organization admins to view pending invitations, resend invitation emails, and revoke unaccepted invitations. Invitations expire after a configurable period (default 7 days). Expired invitations are automatically cleaned up. Invited users who have not yet registered see the invitation context upon registration.

**Rationale:** FR-004 covers sending invitations but not managing them. An org admin who sent an invitation to the wrong email, or whose invitee did not receive the email, needs to resend or revoke. Without invitation management, the admin has no visibility into pending onboarding.

**Fit Criterion:** Org admins can list pending invitations with status (pending, accepted, expired, revoked). Resend triggers a new invitation email. Revoke invalidates the invitation immediately. Invitations expire after the configured period. Expired invitations are not visible to the invitee. (Scale: boolean | Pass/Fail)

**Depends on:** FR-004

---

### FR-099: Organization Activity Audit Trail (0x4e6c)

**Priority:** Should | **Originator:** Research Institution

**Description:** The platform shall provide organization administrators with an org-scoped audit trail: member joins and departures, role changes, campaign creation and state changes, data exports, and invitation activity within their organization. This is the org-level view of the system audit log (FR-077) filtered to the organizations scope.

**Rationale:** FR-049 covers role change audit only. FR-077 covers the full system audit log visible to platform admins. But org admins need visibility into activity within their organization for institutional compliance, data governance, and operational oversight — without needing platform admin access.

**Fit Criterion:** Org admins can view audit trail filtered to their organization scope. Trail includes member joins/departures, role changes, campaign state changes, data exports, and invitation activity. Entries include actor, action, target, timestamp. Accessible only to org admin role. Filterable by action type and date range. (Scale: boolean | Pass/Fail)

**Depends on:** FR-048, FR-077

---

## 6b. Campaign Creation and Management (BUC-02, scope:0x29)

### FR-005: Create Campaign (0x16)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall allow an authenticated researcher with campaign creation permission to create a campaign specifying parameters, geographic region, time window, quality thresholds, and device eligibility criteria.

**Rationale:** Campaigns are the central organizing unit. A researcher must define what data is needed, where, when, and to what quality standard.

**Fit Criterion:** A campaign is created with all required fields (parameters, region, window, thresholds) and is queryable via API. Campaigns missing required fields are rejected with field-level validation errors. (Scale: boolean | Pass/Fail)

**Constrained by:** CON-002 — First Principles Design
**Depends on:** FR-004
**Cross-ref:** scope:0x29, facts:0x9

---

### FR-006: Define Campaign Parameters (0x1c)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall allow researchers to define one or more measurable parameters per campaign, each with units, acceptable range, and precision requirements.

**Rationale:** Campaign parameters must be explicitly defined — no open-ended collection. Each parameter includes acceptable ranges, units, and precision requirements.

**Fit Criterion:** Each campaign parameter has a defined unit, minimum range, maximum range, and precision. Readings outside the defined range are rejected during ingestion. (Scale: boolean | Pass/Fail)

**Depends on:** FR-005
**Cross-ref:** scope:0x29

---

### FR-007: Define Campaign Region (0x1e)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall allow researchers to define a geographic boundary for a campaign as a polygon, radius, or administrative boundary. Campaigns may also be unbounded.

**Rationale:** Geographic boundaries determine which data is relevant. Readings outside the region are rejected.

**Fit Criterion:** Readings with geolocation inside the campaign region are accepted. Readings outside are rejected with reason indicating out-of-region. Unbounded campaigns accept readings from any location. (Scale: boolean | Pass/Fail)

**Depends on:** FR-005
**Cross-ref:** scope:0x29

---

### FR-008: Define Campaign Time Window (0x17)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall require each campaign to have a start and end time (UTC). Readings with timestamps outside the window are rejected.

**Rationale:** Campaigns have defined start and end times. Readings outside the window are rejected.

**Fit Criterion:** Readings timestamped before campaign start or after campaign end are rejected. The platform pushes campaign configuration to enrolled devices when the window opens. (Scale: boolean | Pass/Fail)

**Depends on:** FR-005
**Cross-ref:** scope:0x29

---

### FR-009: Publish and Discover Campaigns (0x18)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall allow researchers to publish campaigns. Published campaigns are discoverable by all authenticated scitizens regardless of organizational affiliation.

**Rationale:** Campaigns must be discoverable by scitizens. Cross-org participation is the normal case — scitizens typically have no institutional affiliation with the researcher.

**Fit Criterion:** A published campaign is visible in the campaign listing API to any authenticated scitizen within 30 seconds of publication. Unpublished campaigns are not visible. (Scale: seconds | Worst: 60 | Plan: 30 | Best: 5)

**Depends on:** FR-005
**Cross-ref:** scope:0x29, facts:0x9

---

### FR-010: Monitor Campaign Data Quality (0x19)

**Priority:** Should | **Originator:** Researcher

**Description:** The platform shall provide researchers with a campaign dashboard showing accepted/rejected reading counts, quality metrics, geographic distribution, and temporal coverage.

**Rationale:** Researchers need visibility into data quality during collection to adjust campaigns if needed.

**Fit Criterion:** Campaign dashboard data is refreshed within 5 minutes of the latest reading submission. Dashboard shows accepted count, rejected count, rejection reasons, and geographic coverage. (Scale: minutes | Worst: 15 | Plan: 5 | Best: 1)

**Depends on:** FR-005
**Cross-ref:** scope:0x29

---

### FR-054: Campaign Editing Constraints (0x2738)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall allow researchers to edit published campaigns with the following constraints: descriptive fields (name, description) are always mutable. Structural fields (parameters, region, time window, eligibility criteria) are mutable only before the campaign window opens. Once the window is open, structural changes are prohibited to protect data consistency for already-enrolled devices. If eligibility criteria change pre-window, enrolled devices that no longer meet criteria are disenrolled with notification.

**Rationale:** Researchers discover errors or need to adjust campaigns. But structural changes to a live campaign would invalidate enrolled devices and in-flight data. The editing model must balance researcher flexibility with data integrity.

**Fit Criterion:** A researcher can edit campaign name and description at any time. A researcher cannot modify parameters, region, window, or eligibility after window opens. Eligibility changes pre-window disenroll ineligible devices with notification. Editing a closed campaign is prohibited. (Scale: boolean | Pass/Fail)

**Depends on:** FR-005

---

### FR-055: Campaign Cancellation and Archival (0x273a)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall allow researchers to cancel an active campaign. Cancellation closes the campaign window immediately, disenrolls all devices, sends cancellation notifications to enrolled scitizens, and preserves all collected data as read-only. Cancelled campaigns are archived and remain queryable for data export. Scitizen contribution scores for the cancelled campaign are preserved.

**Rationale:** Research priorities change. Funding is revoked. Campaigns may need to be cancelled mid-collection. Collected data is still valuable. Scitizens who contributed should retain credit.

**Fit Criterion:** A cancelled campaign stops accepting readings immediately. All enrolled devices are disenrolled. Scitizens receive cancellation notification. Collected data remains exportable. Contribution scores for the campaign are preserved. Campaign status shows cancelled. (Scale: boolean | Pass/Fail)

**Depends on:** FR-005

---

### FR-056: Campaign Duplication (0x2734)

**Priority:** Could | **Originator:** Researcher

**Description:** The platform shall allow researchers to duplicate an existing campaign, creating a new draft campaign with the same parameters, region, eligibility criteria, and quality thresholds but a new time window. The duplicate is independent: editing it does not affect the original. Campaign data is not copied.

**Rationale:** Researchers often repeat campaigns in different seasons, regions, or time windows. Duplication reduces setup friction for repeat experiments and encourages longitudinal study design.

**Fit Criterion:** A duplicated campaign creates a new draft with identical parameters, region, and eligibility as the source. The duplicate has a new unique ID and requires a new time window before publication. Editing the duplicate does not affect the original. (Scale: boolean | Pass/Fail)

**Depends on:** FR-005

---

### FR-057: Campaign Collaboration (0x2736)

**Priority:** Should | **Originator:** Researcher

**Description:** The platform shall allow campaign owners to add collaborators (other researchers within the same organization) with defined permissions: view-only, data-export, or co-manage (full edit/cancel/announce). Collaborator access is scoped to the specific campaign, not inherited from organization role.

**Rationale:** Research is collaborative. Multiple researchers often co-manage a campaign. Campaign-level permissions allow fine-grained control beyond organization roles.

**Fit Criterion:** A campaign collaborator with view-only permission can see campaign data but cannot edit, cancel, or send announcements. A co-manage collaborator can perform all campaign operations. Removing a collaborator revokes their campaign-specific access immediately. (Scale: boolean | Pass/Fail)

**Depends on:** FR-003, FR-005

---

### FR-082: Campaign Detail View (0x4e4a)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall provide a campaign detail view accessible to all authenticated users showing: full campaign description, required parameters with units and ranges, geographic region visualization, time window, device eligibility criteria, current enrollment count, and campaign progress (percentage of data collection target met, if defined). The detail view is the entry point for campaign enrollment.

**Rationale:** FR-012 provides a campaign list with summary. FR-009 makes campaigns discoverable. But neither specifies the detail view that a scitizen uses to evaluate whether to enroll. Without a detail view, the scitizen cannot assess eligibility, understand parameters, or make an informed enrollment decision.

**Fit Criterion:** Campaign detail view displays full description, parameters with units and ranges, region visualization, time window, eligibility criteria, and enrollment count. View loads within 3 seconds. Enrollment action is accessible from the detail view. Ineligible devices are indicated with reason before enrollment attempt. (Scale: seconds | Worst: 10 | Plan: 3 | Best: 1)

**Depends on:** FR-009, FR-012

---

### FR-103: Public Campaign Showcase (0x4e74)

**Priority:** Should | **Originator:** Corewood

**Description:** The platform shall provide a public-facing page (no authentication required) showcasing active campaigns, aggregate platform statistics (total campaigns, total contributors, total readings), and a call-to-action for registration. This is the platform landing page and primary recruitment surface for new scitizens.

**Rationale:** Section 1.2 Goal states community growth is organic — rewarded citizen scientists recruit others through word-of-mouth. Word-of-mouth requires a shareable URL that shows what the platform offers. FR-009 makes campaigns discoverable to authenticated users but unauthenticated visitors see nothing. The landing page is the top of the recruitment funnel.

**Fit Criterion:** The landing page is accessible without authentication. It displays active campaign highlights, aggregate platform statistics, and registration call-to-action. Campaign detail links from the showcase redirect to the login flow and return to the campaign detail after authentication. Page loads within 3 seconds. (Scale: seconds | Worst: 5 | Plan: 3 | Best: 1)


---

### FR-105: Campaign Lifecycle State Machine (0x4e78)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall enforce a defined campaign lifecycle with states: draft (created but not published, editable), published (visible to scitizens, enrollable, not yet collecting), active (campaign window open, collecting data), completed (window closed, data retained), and cancelled (terminated early per FR-055). State transitions are: draft to published, published to active (automatic at window open), active to completed (automatic at window close), any non-completed state to cancelled. Only valid transitions are permitted.

**Rationale:** FR-005 creates a campaign. FR-009 publishes it. FR-008 defines the time window. FR-055 covers cancellation. But no FR defines the full state machine or makes clear that draft campaigns are fully editable while published campaigns have restricted editing (FR-054). The state machine is the backbone of campaign management — every other campaign FR depends on knowing what state a campaign is in.

**Fit Criterion:** Campaigns exist in one of five states: draft, published, active, completed, cancelled. Only valid state transitions are permitted. Invalid transitions are rejected with reason. Draft-to-published requires all required fields. Published-to-active occurs automatically at window open. Active-to-completed occurs automatically at window close. Campaign state is queryable via API. (Scale: boolean | Pass/Fail)

**Depends on:** FR-005, FR-008, FR-009

---

## 6c. User Registration and Account Management (BUC-03, scope:0x2a)

### FR-011: Scitizen Account Registration (0x24)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall allow open registration for scitizens. Identity is delegated to Zitadel. Registration requires only email and password (or social login). No institutional affiliation required.

**Rationale:** Registration is the prerequisite for device enrollment and campaign participation. First contribution is the retention inflection point (FACT-006) — registration must be minimal friction.

**Fit Criterion:** A scitizen can complete registration and reach the campaign browse page in under 2 minutes from landing page. (Scale: minutes | Worst: 5 | Plan: 2 | Best: 1)

**Cross-ref:** scope:0x2a, facts:0xd

---

### FR-012: Browse Campaigns (0x22)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall allow authenticated scitizens to browse published campaigns, filtered by geographic proximity, sensor type compatibility, and campaign status.

**Rationale:** Scitizens need to discover campaigns relevant to their location and devices before enrolling.

**Fit Criterion:** Campaign listing returns results filtered by location and device type within 2 seconds. Results include campaign summary, region, window, and required parameters. (Scale: seconds | Worst: 5 | Plan: 2 | Best: 0.5)

**Depends on:** FR-011
**Cross-ref:** scope:0x2a

---

### FR-037: Password Reset Delegation (0x2712)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall provide a password reset entry point that redirects the user to the Zitadel password reset flow and handles the return redirect. The platform does not store or process credentials but must own the UI entry point and post-reset navigation.

**Rationale:** Users will forget passwords. The platform must provide a discoverable entry point even though the actual reset is handled by Zitadel. Without a platform-owned entry point, users have no way to initiate reset from the Rootstock UI.

**Fit Criterion:** A user who clicks forgot password is redirected to Zitadel reset flow and returned to the Rootstock login page upon completion. The full round-trip completes without manual URL entry. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011

---

### FR-038: Account Deactivation (0x2714)

**Priority:** Should | **Originator:** Scitizen

**Description:** The platform shall allow a user to temporarily deactivate their account. Deactivation sets the app_users status to deactivated, suspends all owned devices, pauses campaign enrollments, and prevents login. Deactivation is reversible: reactivation restores previous device and enrollment state.

**Rationale:** Users may need to take a break from the platform without losing their devices, contribution history, or campaign enrollments. Deactivation provides a reversible pause distinct from permanent deletion.

**Fit Criterion:** A deactivated user cannot login. All owned devices show status suspended. Campaign enrollments are paused. Upon reactivation, devices return to their pre-deactivation status and campaign enrollments resume. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011

---

### FR-039: Account Deletion and Right to Erasure (0x2711)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall allow a user to permanently delete their account. Deletion revokes all owned devices, removes campaign enrollments, deletes the app_users record (including idp_id reference), and anonymizes contributor references in persisted readings. Submitted readings are preserved with pseudonymized provenance but the link to the deleted user is permanently severed. The platform notifies Zitadel to delete the corresponding identity.

**Rationale:** GDPR Article 17 right to erasure and general user trust require permanent deletion capability. The platform stores only idp_id and application state. Deletion removes the idp_id link and application state. Submitted readings are scientific data that must be preserved, but contributor identity must be unrecoverable post-deletion.

**Fit Criterion:** After account deletion: no app_users record exists for the deleted user, all owned devices are revoked, no campaign enrollments reference the deleted user, persisted readings retain pseudonymized provenance but no query can resolve the deleted user identity. Deletion completes within 24 hours of request. (Scale: hours | Worst: 48 | Plan: 24 | Best: 1)

**Depends on:** FR-011

---

### FR-042: Profile Visibility Settings (0x2719)

**Priority:** Could | **Originator:** Scitizen

**Description:** The platform shall allow users to control their profile visibility: public (contribution score, badges, and campaign participation visible to all authenticated users), organization-scoped (visible only to members of shared organizations), or private (visible only to self). Profile visibility applies to application state only; PII from Zitadel is never exposed to other users without explicit consent.

**Rationale:** Some scitizens want public recognition; others prefer privacy. Profile visibility must be user-controlled. Since PII lives in Zitadel and the app DB stores only application state, visibility controls apply to contribution data, badges, and participation history.

**Fit Criterion:** A user with private profile is not visible in any user listing or campaign contributor list. A user with public profile shows contribution score and badges to any authenticated user. Default visibility is private. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011

---

### FR-043: Profile Data Resolution from IdP (0x271b)

**Priority:** Must | **Originator:** Scitizen, Researcher

**Description:** The platform shall resolve user display name and profile picture from Zitadel on demand when rendering user-facing contexts (profile page, campaign contributor lists, org member lists). PII is never cached in the app DB. The platform requests PII from Zitadel at render time using the idp_id reference.

**Rationale:** The app DB stores only idp_id and application state, not PII. Display name and profile picture are personal data managed by Zitadel. The platform must fetch this data at render time to display user-facing contexts without violating the data privacy principle.

**Fit Criterion:** No PII (display name, email, profile picture) is stored in the app DB. User-facing displays resolve PII from Zitadel via idp_id. If Zitadel is unreachable, the platform displays a placeholder, not cached PII. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011

---

### FR-050: User Type Modification (0x272c)

**Priority:** Should | **Originator:** Scitizen, Researcher

**Description:** The platform shall allow a user to add or remove user_type roles (scitizen, researcher) after initial registration. Adding researcher type does not automatically grant institutional affiliation. Removing scitizen type does not affect owned devices or contribution history. Removing researcher type does not delete campaigns created while a researcher.

**Rationale:** FR-011 allows type selection at registration, but user needs evolve. A scitizen may later become a researcher, or a researcher may want to also contribute as a scitizen. Types are additive roles, not mutually exclusive states.

**Fit Criterion:** A user can add researcher to their user_type and subsequently access campaign creation (after org invitation). A user can remove researcher from their user_type; existing campaigns are unaffected. A user must retain at least one user_type. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011

---

### FR-078: Email Verification (0x4e42)

**Priority:** Must | **Originator:** Corewood

**Description:** The platform shall require email verification before an account is fully activated. Upon registration, a verification link is sent to the provided email address. The link expires after 24 hours. Unverified accounts can authenticate but cannot enroll devices, join campaigns, or create campaigns. Re-sending the verification email is allowed with rate limiting.

**Rationale:** FR-011 does not specify email verification. Without it, spam registrations pollute the user base, notification delivery (FR-051) fails silently, and account recovery (FR-037) is unreliable. Email verification is a standard platform requirement and an OWASP recommendation.

**Fit Criterion:** Registration triggers a verification email. Unverified accounts cannot enroll devices or join/create campaigns. Verification link expires after 24 hours. Re-send is rate-limited to 3 per hour. Verified accounts gain full access immediately upon verification. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011

---

### FR-097: Password Change (0x4e69)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall allow authenticated users to change their password from account settings. The user must provide their current password and the new password. The new password is validated against the identity providers password policy. On successful change, all other active sessions are invalidated and the user is notified via email.

**Rationale:** FR-037 covers password reset (forgot password, unauthenticated). Password change is distinct — the user knows their current password and wants to update it. Session invalidation on password change is an OWASP requirement to prevent continued use of a potentially compromised session.

**Fit Criterion:** Authenticated users can change password from account settings by providing current and new password. Invalid current password is rejected. New password must pass identity provider policy. On success, all other sessions are invalidated. Email notification is sent confirming the change. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011, FR-079

---

## 6d. Device Registration and Enrollment (BUC-04, scope:0x2b)

### FR-013: Generate Enrollment Code (0x2b)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall generate a short-lived (15-minute TTL), one-time-use, human-friendly enrollment code (6-8 alphanumeric characters, no ambiguous characters 0/O/1/l) when a scitizen initiates device registration.

**Rationale:** Enrollment code is the bootstrap trust anchor binding web session to physical device. Must be human-friendly, short-lived, and one-time use.

**Fit Criterion:** Enrollment codes expire after 15 minutes. Used codes are rejected on second use. Codes contain only unambiguous alphanumeric characters. (Scale: minutes | Pass/Fail)

**Depends on:** FR-011
**Cross-ref:** scope:0x2b

---

### FR-014: Direct Device Enrollment (Tier 1) (0x2d)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall accept device enrollment via POST /enroll where the device presents an enrollment code and a locally-generated CSR over HTTPS. The enrollment service validates the code, verifies the CSR, coordinates with the CA, and returns the issued certificate.

**Rationale:** Tier 1 devices (smartphones, Raspberry Pi, modern weather stations) can run HTTPS clients and generate keypairs directly.

**Fit Criterion:** A Tier 1 device presenting a valid enrollment code and CSR receives an X.509 certificate and is registered in the device registry with status active within 10 seconds. (Scale: seconds | Worst: 30 | Plan: 10 | Best: 3)

**Depends on:** FR-013
**Cross-ref:** scope:0x2b, facts:0x8

---

### FR-015: Proxy Device Enrollment (Tier 2) (0x27)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall support proxy enrollment where a companion app communicates with a Tier 2 device over BLE/local WiFi, requests a CSR from the device, submits the enrollment on behalf of the device, and pushes the issued certificate back.

**Rationale:** Tier 2 devices (ESP32, LoRa gateways) have no rich UI and limited TLS. Companion app proxies the enrollment flow. Private key never leaves device.

**Fit Criterion:** A Tier 2 device enrolled via companion app receives a certificate and is registered with status active. The device private key is never transmitted to the companion app or platform. (Scale: boolean | Pass/Fail)

**Depends on:** FR-013
**Cross-ref:** scope:0x2b, facts:0x8

---

### FR-016: Device Registry Entry (0x29)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall create a device registry entry upon enrollment containing: device ID (matching certificate CN), owner ID, status, device class, firmware version, tier, sensor capabilities, and certificate serial.

**Rationale:** The device registry is the single source of truth for device state. All mutable metadata lives here, not in the certificate.

**Fit Criterion:** Every enrolled device has a registry entry with all required fields. Device ID matches certificate CN. Registry is queryable by device ID, owner ID, status, and device class. (Scale: boolean | Pass/Fail)

**Depends on:** FR-014
**Cross-ref:** scope:0x2b

---

### FR-058: Device Deregistration (0x273b)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall allow a scitizen to deregister a device they own. Deregistration revokes the device certificate, removes all campaign enrollments for that device, and sets the device registry status to revoked. Submitted readings from the deregistered device are preserved with provenance intact. The device must re-enroll to rejoin the platform.

**Rationale:** Scitizens sell, discard, or retire devices. A deregistration path ensures clean removal from the platform. Submitted data is scientific output that must persist regardless of device lifecycle.

**Fit Criterion:** A deregistered device has status revoked in the registry, zero campaign enrollments, and a revoked certificate. OPA denies all actions from the device within 30 seconds. Previously submitted readings remain queryable with full provenance. (Scale: seconds | Worst: 60 | Plan: 30 | Best: 5)

**Depends on:** FR-016

---

### FR-059: Device Ownership Transfer (0x273d)

**Priority:** Should | **Originator:** Scitizen

**Description:** The platform shall support transferring device ownership from one scitizen to another. Transfer updates the owner_id in the device registry. Campaign enrollments are preserved if the new owner confirms them. The device certificate remains valid (CN is device-id, not owner-id). If the original owner is uncooperative, the device is revoked and must be re-enrolled by the new owner.

**Rationale:** Devices change hands (gifts, sales, hand-me-downs). The glossary defines ownership transfer. Certificate is tied to device-id, not owner, so transfer does not require re-enrollment in the cooperative case.

**Fit Criterion:** After cooperative transfer, device owner_id reflects the new owner. Campaign enrollments are preserved pending new owner confirmation. After uncooperative transfer (revoke + re-enroll), the device has a new certificate and the new owner as owner_id. Previous readings retain the original owner provenance. (Scale: boolean | Pass/Fail)

**Depends on:** FR-016

---

### FR-060: Device Metadata Update (0x273f)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall update device registry metadata (firmware version, sensor capabilities) when reported by the device. Firmware version is reported on each MQTT connect or HTTP request via client metadata. Sensor capability changes (e.g., after OTA update adding a new sensor type) are reported by the device and validated against the known device class before registry update.

**Rationale:** Devices receive OTA firmware updates outside platform control. The registry must stay current for eligibility checks, security response (bulk suspension by firmware version), and provenance tracking. Stale metadata leads to incorrect authorization decisions.

**Fit Criterion:** When a device reports a new firmware version, the registry is updated within the same connection session. Sensor capability changes are reflected in the registry after validation. Eligibility checks for campaign enrollment use the current registry metadata, not stale data. (Scale: boolean | Pass/Fail)

**Depends on:** FR-016

---

### FR-061: Device Health Monitoring (0x2741)

**Priority:** Should | **Originator:** Scitizen

**Description:** The platform shall track device health via MQTT Last Will and Testament (LWT) messages and periodic heartbeat on the rootstock/{device-id}/status topic. When a device disconnects unexpectedly, the LWT message updates the device last-seen timestamp and sets a disconnected flag. The scitizen device management dashboard shows device online/offline status and last-seen time.

**Rationale:** Scitizens need to know if their devices are functioning. The MQTT LWT feature (defined in glossary) provides immediate disconnect detection. Without health monitoring, a silently offline device produces no data and the scitizen has no visibility into the problem.

**Fit Criterion:** A device that disconnects unexpectedly triggers an LWT update within 30 seconds. The device management dashboard shows the device as offline with last-seen timestamp. A device that reconnects clears the offline flag within one heartbeat interval. (Scale: seconds | Worst: 60 | Plan: 30 | Best: 5)

**Depends on:** FR-016

---

### FR-109: Device Certificate Expiry Warning (0x4e80)

**Priority:** Must | **Originator:** Gap analysis — FR-028 handles automated renewal but no FR specifies the scitizen-facing notification when a certificate approaches expiry or renewal fails

**Description:** The platform shall notify device owners when a device certificate is approaching its expiry date and when automated renewal succeeds or fails. If renewal fails, the notification shall include instructions for manual resolution.

**Rationale:** Certificate expiry that goes unnoticed results in silent device disconnection and data loss. Automated renewal (FR-028) handles the happy path, but scitizens need visibility into the process — especially when it fails — to maintain device uptime and data contribution continuity.

**Fit Criterion:** Device owners receive a notification at least 7 days before certificate expiry. If automated renewal fails, a follow-up notification with resolution steps is sent within 1 hour of failure. Notification delivery is logged and verifiable. (Scale: boolean | Pass/Fail)

**Depends on:** FR-028, FR-051

---

### FR-118: Device Firmware Compatibility Check (0x4e92)

**Priority:** Should | **Originator:** Gap analysis — FR-019 checks eligibility by sensor type and region but no FR covers firmware version compatibility between campaign requirements and device capabilities

**Description:** Campaigns shall be able to specify minimum firmware version requirements for enrolled devices. The eligibility check (FR-019) shall include firmware version comparison. Devices that do not meet the firmware requirement shall be informed of the incompatibility and, where possible, directed to update their firmware.

**Rationale:** IoT device firmware determines data format, sensor capabilities, and communication protocols. A campaign expecting a specific data schema cannot accept readings from devices running incompatible firmware. FR-060 (device metadata update) tracks firmware version; this requirement uses that information at enrollment time.

**Fit Criterion:** Campaign parameters include an optional minimum firmware version field. When set, the enrollment eligibility check compares the device firmware version from the registry against the campaign requirement. Incompatible devices receive a clear message identifying the required version and their current version. (Scale: boolean | Pass/Fail)

**Depends on:** FR-016, FR-060

---

### FR-122: Companion Device Setup Flow (0x4e9b)

**Priority:** Must | **Originator:** Gap analysis — FR-015 defines proxy enrollment but no FR specifies the guided setup experience for pairing a companion device (phone) with a sensor and verifying data flow

**Description:** The platform shall provide a guided setup flow for Tier 2 (proxy) device enrollment. The flow shall walk the scitizen through pairing their companion device with the sensor, verifying the connection, confirming sensor readings are received, and completing enrollment. The flow shall handle common failure modes (pairing failure, no readings received) with clear recovery instructions.

**Rationale:** Proxy enrollment (FR-015) is the most complex enrollment path — it involves three entities (platform, phone, sensor) instead of two. Non-technical scitizens in field conditions need step-by-step guidance. US-001 requires first-use enrollment under 10 minutes; achieving that for Tier 2 requires a well-designed guided flow.

**Fit Criterion:** The proxy enrollment flow presents numbered steps guiding the scitizen through pairing, connection verification, and reading confirmation. Each step has success and failure states with recovery instructions. A scitizen can complete the flow in under 10 minutes (per US-001). The flow detects common failure modes and provides specific guidance for each. (Scale: boolean | Pass/Fail)

**Depends on:** FR-015

---

## 6e. Campaign Enrollment (BUC-05, scope:0x2c)

### FR-017: Enroll Device in Campaign (0x32)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall allow a scitizen to enroll a registered, active device in a published campaign whose window has not closed, provided the device meets campaign eligibility criteria (device class, tier, sensor capabilities).

**Rationale:** Campaign enrollment links devices to data needs. A device may be enrolled in multiple campaigns. Cross-org participation is the normal case.

**Fit Criterion:** After enrollment, the device registry reflects the campaign association. OPA allows publish to the campaign topic on the next bundle refresh (within 30 seconds). Device receives campaign configuration. (Scale: seconds | Worst: 60 | Plan: 30 | Best: 5)

**Depends on:** FR-016
**Cross-ref:** scope:0x2c

---

### FR-018: Multi-Campaign Device Enrollment (0x2e)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall support simultaneous enrollment of a single device in multiple active campaigns. Data routing uses the topic structure to associate readings with the correct campaign.

**Rationale:** A device may be enrolled in multiple campaigns simultaneously. This maximizes data utility and scitizen contribution value.

**Fit Criterion:** A device enrolled in N campaigns can publish to all N campaign topics. Readings are correctly routed to each campaign independently. (Scale: campaigns | Pass/Fail)

**Depends on:** FR-017
**Cross-ref:** scope:0x2c

---

### FR-019: Device Eligibility Check (0x30)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall check device eligibility against campaign criteria (device class, tier, sensor capabilities, firmware version) at enrollment time and reject ineligible devices with a specific reason.

**Rationale:** Researchers define which device types and capabilities are acceptable for their campaign. Eligibility is checked at enrollment time.

**Fit Criterion:** An ineligible device (wrong class, tier, or missing required sensor) is rejected at enrollment with a human-readable reason. Eligible devices are enrolled. (Scale: boolean | Pass/Fail)

**Depends on:** FR-017
**Cross-ref:** scope:0x2c

---

### FR-062: Campaign Configuration Push (0x4e21)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall push campaign configuration (sampling parameters, quality thresholds, geographic bounds) to all enrolled devices when a campaign window opens (BE-08). Configuration is delivered via MQTT retained messages on rootstock/{device-id}/config and via HTTP/2 for devices using the gateway path.

**Rationale:** BE-08 (Campaign Window Opens) produces Data Flow 0x13 (Device Configuration). BUC-05 postcondition states device receives campaign configuration. FR-008 mentions config push only in its fit criterion — a fit criterion is a test condition, not a requirement. The platform behavior of delivering configuration at window open must be an explicit functional requirement.

**Fit Criterion:** When a campaign window opens, all enrolled devices receive campaign configuration within 60 seconds. Configuration is available as a retained MQTT message on each device config topic. Devices connecting after window open receive the retained config immediately. HTTP/2 gateway devices can poll for config. (Scale: seconds | Worst: 120 | Plan: 60 | Best: 10)

**Depends on:** FR-008, FR-017

---

### FR-063: Campaign Enrollment Withdrawal (0x4e23)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall allow a scitizen to withdraw a device from a specific campaign without deregistering the device from the platform. Withdrawal updates the device registry to remove the campaign association, OPA revokes publish authorization for that campaign topic on the next bundle refresh, and the consent record for that enrollment is marked as withdrawn with a timestamp.

**Rationale:** CO-002 fit criterion states consent withdrawal stops further data collection from the device for that campaign. CO-001 requires granular per-campaign consent. FR-058 covers full device deregistration but no FR covers partial withdrawal from a single campaign. A scitizen contributing to campaigns A and B must be able to stop contributing to A while continuing B.

**Fit Criterion:** After withdrawal, OPA denies publish to the withdrawn campaign topic within 30 seconds. The device remains enrolled in other campaigns. The consent record shows withdrawal timestamp. Previously submitted readings for the withdrawn campaign are preserved. The device registry no longer lists the withdrawn campaign for that device. (Scale: seconds | Worst: 60 | Plan: 30 | Best: 5)

**Depends on:** FR-017

---

### FR-064: Consent Capture at Campaign Enrollment (0x4e26)

**Priority:** Must | **Originator:** Oversight Body

**Description:** The platform shall present campaign-specific consent terms to the scitizen at the time of campaign enrollment and record explicit consent with timestamp, consent version, and scope before completing enrollment. Consent records are stored in the app database and are exportable for IRB audit. Enrollment does not proceed without recorded consent.

**Rationale:** CO-002 requires every campaign enrollment to record explicit consent with timestamp, version, and scope. CO-001 requires granular per-campaign consent. FR-017 covers enrollment mechanics but does not specify the consent capture step. The functional mechanism for capturing consent at enrollment time is a platform behavior that must be explicitly specified.

**Fit Criterion:** Every completed campaign enrollment has a consent record with timestamp, consent version, and scope. Enrollment attempted without consent acceptance is rejected. Consent records are exportable in a machine-readable format for IRB audit. Consent version is traceable to the campaign consent terms active at enrollment time. (Scale: boolean | Pass/Fail)

**Depends on:** FR-017

---

### FR-092: Campaign Enrollment Funnel (0x4e5f)

**Priority:** Should | **Originator:** Researcher

**Description:** The platform shall show researchers an enrollment funnel for each campaign: campaign page views, eligible device count, enrolled device count, actively contributing device count, and devices that have submitted at least one valid reading. This allows researchers to diagnose whether a campaign is failing at discovery, eligibility, enrollment, or activation.

**Rationale:** No existing FR tracks the scitizen journey from campaign discovery to active contribution. A campaign with many views but few enrollments signals a usability or eligibility problem. A campaign with many enrollments but few submissions signals a device configuration or motivation problem. The funnel makes this diagnosable.

**Fit Criterion:** Campaign dashboard displays funnel: page views, eligible devices, enrolled devices, actively contributing devices, devices with at least one accepted reading. Each stage count refreshed within 15 minutes. (Scale: minutes | Worst: 30 | Plan: 15 | Best: 5)

**Depends on:** FR-010, FR-017

---

### FR-095: Scitizen Campaign Progress View (0x4e64)

**Priority:** Should | **Originator:** Scitizen

**Description:** The platform shall show enrolled scitizens the progress of campaigns they are contributing to: how much data has been collected overall (anonymized), how close the campaign is to its data collection target (if defined), time remaining in the campaign window, and the scitizens own contribution relative to the campaign total.

**Rationale:** BUC-10 and the glossary emphasize that scitizens need to see that their data matters. Aggregate campaign progress — not just personal stats — shows the scitizen that they are part of something larger. Self-Determination Theory (referenced in BUC-10 business rules) holds that perceived autonomy, competence, and relatedness drive intrinsic motivation. Seeing campaign progress feeds relatedness.

**Fit Criterion:** Enrolled scitizens see campaign progress: total readings collected (anonymized count), progress toward target (percentage if target defined), time remaining, and their own contribution count. Progress data refreshed within 15 minutes. No individual contributor data is visible to other scitizens. (Scale: minutes | Worst: 30 | Plan: 15 | Best: 5)

**Depends on:** FR-040

---

### FR-116: Campaign Enrollment Approval (0x4e8e)

**Priority:** Should | **Originator:** Gap analysis — FR-017 covers enrollment and FR-019 covers eligibility check but no FR allows researchers to manually approve or reject enrollment requests for campaigns requiring human vetting

**Description:** The platform shall support an optional enrollment approval workflow for campaigns. When enabled by the researcher, enrollment requests enter a pending state and require researcher approval before the device begins collecting data. Researchers shall see pending enrollments in their campaign management view with device details sufficient to make an approval decision.

**Rationale:** FR-019 automates eligibility checks against sensor type and geographic region. But some campaigns have qualitative criteria that cannot be automated — device placement quality, environmental suitability, or capacity limits. Manual approval gives researchers control over campaign composition when needed.

**Fit Criterion:** Campaigns can be configured to require enrollment approval. When enabled, enrollment requests enter a pending state visible to the researcher. Researchers can approve or reject with an optional reason. Scitizens are notified of the decision. Approved devices begin data collection; rejected devices do not. (Scale: boolean | Pass/Fail)

**Depends on:** FR-016, FR-017

---

## 6f. Data Ingestion and Validation (BUC-06, scope:0x2d)

### FR-020: Authenticate Device via mTLS (0x34)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall require mTLS for all device connections (MQTT and HTTP/2). The device presents its certificate; the platform verifies against the CA chain and extracts the device ID from the certificate CN.

**Rationale:** All device-to-platform communication must be mutually authenticated. TLS does authentication — device presents cert, platform verifies CA signature, extracts device ID.

**Fit Criterion:** Connections without a valid client certificate are rejected at the TLS layer. Connections with expired, revoked, or untrusted certificates are rejected with structured diagnostic logs. (Scale: boolean | Pass/Fail)

**Cross-ref:** scope:0x2d, facts:0x5

---

### FR-021: Authorize Device Actions via OPA (0x36)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall query OPA for authorization on every device action (MQTT connect, publish, subscribe; HTTP request). OPA evaluates device status (active/suspended/revoked), campaign enrollment, and topic ACLs against the device registry.

**Rationale:** Authorization is separate from authentication. OPA checks device status, campaign enrollment, and topic ACLs on every action.

**Fit Criterion:** A suspended device is denied all actions. A revoked device is denied all actions. A device publishing to a topic for a campaign it is not enrolled in is denied. Denial reasons are logged. (Scale: boolean | Pass/Fail)

**Depends on:** FR-020
**Cross-ref:** scope:0x2d

---

### FR-022: Validate Sensor Readings (0x3b)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall validate every incoming reading against: schema (correct fields and types), parameter range (within campaign-defined bounds), rate limits (per device), geolocation (within campaign region), and timestamp (within campaign window). Rejected readings return explicit reasons.

**Rationale:** All device data crosses a trust boundary and is treated as untrusted (FACT-010). Validation includes schema, range, rate limiting, geolocation, and timestamp checks.

**Fit Criterion:** Readings failing any validation check are rejected with a specific reason (schema error, out of range, rate exceeded, out of region, out of window). Valid readings are accepted and persisted with full provenance. (Scale: boolean | Pass/Fail)

**Depends on:** FR-021
**Cross-ref:** scope:0x2d, facts:0x2f

---

### FR-023: Persist Valid Readings with Provenance (0x38)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall persist every valid reading with full provenance metadata: device ID, timestamp, geolocation, firmware version, certificate serial, campaign ID, and ingestion timestamp.

**Rationale:** Data provenance is required for publication-grade data. Each reading must be traceable to the device, location, time, firmware version, and transmission path.

**Fit Criterion:** Every persisted reading has all provenance fields populated. No reading is persisted without device ID, timestamp, geolocation, firmware version, and campaign ID. (Scale: boolean | Pass/Fail)

**Depends on:** FR-022
**Cross-ref:** scope:0x2d

---

### FR-024: Topic ACL Enforcement (0x3d)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall enforce topic ACLs via OPA: a device may only publish to rootstock/{its-own-device-id}/* topics. Attempts to publish to another devices topic are denied and logged.

**Rationale:** A device can only publish to topics under its own device ID. This prevents a compromised device from injecting data as another device.

**Fit Criterion:** A device attempting to publish to rootstock/{other-device-id}/data/{campaign} is denied. The denial is logged with the device ID, attempted topic, and timestamp. (Scale: boolean | Pass/Fail)

**Depends on:** FR-021
**Cross-ref:** scope:0x2d, facts:0x2f

---

### FR-025: Anomaly Flagging (0x3f)

**Priority:** Should | **Originator:** Researcher

**Description:** The platform shall flag readings that pass validation but fall outside expected statistical bounds for the campaign. Flagged readings are quarantined for review, not deleted. Campaign policy determines disposition.

**Rationale:** Outliers should be flagged, not silently discarded. Disposition is determined by campaign policy.

**Fit Criterion:** Readings outside 3 standard deviations of the campaign rolling average are flagged and quarantined. Quarantined readings are queryable separately. No quarantined reading is deleted without explicit researcher action. (Scale: standard deviations | Pass/Fail)

**Depends on:** FR-022
**Cross-ref:** scope:0x2d

---

### FR-065: Quarantine Review Workflow (0x4e27)

**Priority:** Should | **Originator:** Researcher

**Description:** The platform shall provide researchers with a quarantine review interface for dispositioning flagged readings. Researchers can accept (move to validated), reject (mark as discarded with reason), or escalate (request additional review) quarantined readings individually or in bulk. Discarded readings are soft-deleted and remain in audit logs.

**Rationale:** FR-025 (anomaly flagging) and FR-032 (vulnerable window flagging) both state that quarantined readings cannot be deleted without explicit researcher action. This creates an implied requirement for a review mechanism but no FR specifies the affirmative capability. BUC-09 may quarantine thousands of readings from an entire device class — bulk disposition is essential.

**Fit Criterion:** Researchers can list quarantined readings filtered by campaign, device, and quarantine reason. Individual and bulk accept/reject operations are supported. Rejected readings are soft-deleted with reason and remain in audit logs. Accepted readings move to validated status and appear in campaign data. All disposition actions are logged with researcher ID and timestamp. (Scale: boolean | Pass/Fail)

**Depends on:** FR-025, FR-032

---

### FR-123: Campaign Data Quality Thresholds (0x4e9c)

**Priority:** Must | **Originator:** Gap analysis — FR-010 monitors data quality and FR-084 provides alerting but no FR specifies user-configurable quality thresholds that trigger automated responses

**Description:** Researchers shall be able to define data quality thresholds for their campaigns — anomaly rate, missing reading rate, device dropout rate, and minimum reading frequency. When a threshold is breached, the platform shall trigger an alert (per FR-084) and optionally take automated action such as pausing enrollment or flagging the campaign for review.

**Rationale:** FR-010 provides data quality monitoring as a view. But passive monitoring requires the researcher to check regularly. Active threshold-based alerting turns monitoring into a responsive system — the researcher is notified when quality degrades, enabling timely intervention before data integrity is compromised.

**Fit Criterion:** Campaign configuration includes settable thresholds for at least: anomaly rate, missing reading rate, and device dropout rate. Threshold breaches generate an alert within 15 minutes. Researchers can configure whether threshold breach triggers notification only or notification plus automated action. Threshold values are editable after campaign publication. (Scale: boolean | Pass/Fail)

**Depends on:** FR-010, FR-084

---

### FR-124: Device Location Verification (0x4ea4)

**Priority:** Should | **Originator:** Gap analysis — FR-007 defines campaign regions and FR-019 checks eligibility but no FR specifies how device location claims are verified against actual placement

**Description:** The platform shall support device location verification mechanisms. When a campaign requires geographic placement, the platform shall accept location evidence from the device (GPS coordinates, network-derived location) and compare it against the claimed deployment location. Significant discrepancies shall flag the device for review without automatically rejecting it.

**Rationale:** Campaign data quality depends on devices being where they claim to be. A temperature sensor enrolled in a coastal monitoring campaign but placed inland produces misleading data. Location verification is not foolproof but provides a baseline integrity check that supports data provenance (CO-003).

**Fit Criterion:** Devices enrolled in geographically-scoped campaigns report location data. The platform compares reported location against the campaign region boundary. Devices outside the boundary are flagged for researcher review. Location checks occur at enrollment and periodically during the campaign. Flagged devices are not automatically rejected. (Scale: boolean | Pass/Fail)

**Depends on:** FR-007, FR-019

---

### FR-125: Reading Rejection Feedback (0x4ea2)

**Priority:** Must | **Originator:** Gap analysis — FR-022 validates readings and FR-025 flags anomalies but no FR specifies what feedback the scitizen receives when their readings are rejected or flagged

**Description:** When a device reading fails validation (FR-022) or is flagged as anomalous (FR-025), the platform shall provide feedback to the device owner. Feedback shall indicate which readings were affected, the reason for rejection or flagging, and any action the scitizen can take (reposition sensor, check calibration, restart device). Feedback shall not disclose validation algorithms or thresholds.

**Rationale:** US-003 requires clear feedback on data contribution. Without rejection feedback, scitizens submit data into a void — they cannot tell if their contributions are valuable or if their device has a problem. Actionable feedback enables scitizens to self-correct, improving data quality without researcher intervention.

**Fit Criterion:** Scitizens can see the validation status of their submitted readings (accepted, rejected, flagged). Rejected and flagged readings show a human-readable reason and suggested corrective action. Feedback is available within 5 minutes of reading submission. Validation algorithm details and thresholds are not disclosed. (Scale: boolean | Pass/Fail)

**Depends on:** FR-022, FR-025

---

## 6g. Data Export and Analysis (BUC-07, scope:0x2e)

### FR-026: Export Campaign Data (0x43)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall allow researchers with export permission to export validated campaign data including readings, provenance metadata, quality metrics, and campaign context.

**Rationale:** Researchers need to extract collected data for analysis in external tools. Export must include provenance metadata for publication.

**Fit Criterion:** Exported data includes all validated readings with provenance fields. Export completes within 60 seconds for campaigns with up to 1 million readings. (Scale: seconds | Worst: 120 | Plan: 60 | Best: 15)

**Depends on:** FR-023
**Cross-ref:** scope:0x2e

---

### FR-027: Separate Contributor Identity from Data (0x41)

**Priority:** Must | **Originator:** Oversight Body

**Description:** The platform shall ensure that exported campaign data does not contain contributor identity. Device IDs in exports are pseudonymized. Raw contributor locations never appear in public-facing datasets.

**Rationale:** 4 spatiotemporal data points uniquely identify 95% of individuals (FACT-009). Contributor identity must be separable from observation data.

**Fit Criterion:** No exported dataset contains scitizen names, emails, or unhashed device IDs. Pseudonymized device IDs cannot be reversed without platform access. Spatial resolution in exports is configurable per campaign. (Scale: boolean | Pass/Fail)

**Depends on:** FR-026
**Cross-ref:** scope:0x2e, facts:0x30

---

### FR-083: Researcher In-Platform Data Explorer (0x4e4c)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall provide researchers with an in-platform data exploration interface for their campaigns: browse validated readings with filtering by time range, parameter, device class, and geographic sub-region; view summary statistics (count, mean, median, min, max, standard deviation) for filtered results; and preview data before export. All device identifiers are pseudonymized. Location precision is limited to the campaign-configured spatial resolution.

**Rationale:** FR-026 provides bulk export but no in-platform browsing. A researcher evaluating data quality, identifying patterns, or deciding what to export needs to filter and aggregate without downloading the full dataset. Data explorer is the read path for campaign data — export is the download path. Privacy constraints (FR-027, SEC-004) apply identically to in-platform views.

**Fit Criterion:** Researchers can browse validated readings filtered by time range, parameter, device class, and geographic sub-region. Summary statistics are displayed for filtered results. Device IDs are pseudonymized. Location precision limited to campaign spatial resolution. Query results return within 5 seconds for up to 100,000 readings. (Scale: seconds | Worst: 15 | Plan: 5 | Best: 2)

**Depends on:** FR-026, FR-027

---

### FR-087: Programmatic API Access for Researchers (0x4e54)

**Priority:** Should | **Originator:** Researcher

**Description:** The platform shall provide authenticated API endpoints for researchers to programmatically query their campaign data: list campaigns, query readings with filtering (time range, parameter, device class, region), retrieve summary statistics, and trigger data exports. API access respects the same permission model, pseudonymization, and spatial resolution constraints as the UI. API authentication uses bearer tokens issued by the identity provider.

**Rationale:** Researchers use external analysis tools (R, Python, MATLAB, Jupyter). UI-only data access forces manual export-download-import workflows. Programmatic access enables automated pipelines, reproducible analysis, and integration with institutional data infrastructure.

**Fit Criterion:** Authenticated researchers can query campaign data via API with the same filtering capabilities as the UI data explorer. API responses use a documented schema. Pseudonymization and spatial resolution constraints are enforced identically to UI access. API documentation is publicly available. Rate limiting prevents abuse. (Scale: boolean | Pass/Fail)

**Depends on:** FR-026, FR-027

---

### FR-090: Data Retention Policy (0x4e5a)

**Priority:** Must | **Originator:** Corewood

**Description:** The platform shall enforce configurable data retention policies: campaign data is retained for a researcher-specified period after campaign completion (minimum 1 year, default 5 years). After the retention period, data is archived or deleted per institutional policy. Account deletion (FR-039) anonymizes the contributor association but preserves the scientific data. Audit logs are retained for the legally required period. Retention policies are visible to researchers at campaign creation time.

**Rationale:** No FR specifies what happens to data over time. Scientific data has long-term value — peer review, replication studies, meta-analyses may occur years after collection. But GDPR (CO-001) requires that data not be retained indefinitely without purpose. The retention policy balances scientific value against data minimization.

**Fit Criterion:** Researchers set a retention period at campaign creation (minimum 1 year). Data is available for the full retention period after campaign completion. After retention expiry, data is archived or deleted per policy. Account deletion anonymizes contributor association without deleting readings. Retention policy is displayed at campaign creation and in campaign detail. (Scale: boolean | Pass/Fail)

**Depends on:** FR-026, FR-039

---

### FR-106: Campaign Data Visualization (0x4e7a)

**Priority:** Should | **Originator:** Researcher

**Description:** The platform shall provide researchers with data visualization for campaign readings: time series plots for each parameter, distribution histograms, and scatter plots for multi-parameter campaigns. Visualizations respect the same pseudonymization and spatial resolution constraints as the data explorer (FR-083). Visualizations are interactive — researchers can zoom, filter, and select subsets.

**Rationale:** FR-083 provides tabular data exploration. But researchers think in charts — a time series shows trends, seasonality, and gaps that a table cannot. Distribution histograms reveal outlier patterns. Scatter plots show parameter correlations. Visualization is how researchers evaluate data quality before committing to export and analysis.

**Fit Criterion:** Researchers can view time series, histograms, and scatter plots for campaign data. Visualizations render within 10 seconds for up to 100,000 readings. Interactive zoom and filter are supported. Pseudonymization and spatial resolution constraints are enforced. Visualizations can be exported as images. (Scale: seconds | Worst: 30 | Plan: 10 | Best: 3)

**Depends on:** FR-083

---

### FR-107: Campaign Completion Summary (0x4e7c)

**Priority:** Should | **Originator:** Researcher

**Description:** The platform shall generate a campaign completion summary when a campaign transitions to completed state: total readings collected, acceptance rate, rejection breakdown by reason, geographic coverage achieved, temporal coverage achieved, contributing device count, contributing scitizen count (anonymized), anomaly count, and data quality assessment relative to campaign thresholds. The summary is downloadable as a report.

**Rationale:** When a campaign completes, the researcher needs a comprehensive summary before deciding to export. The summary also serves as documentation for grant reports, IRB reporting, and publication methodology sections. No FR specifies what happens at campaign completion from the researchers perspective.

**Fit Criterion:** Campaign completion automatically generates a summary report. Report includes all specified metrics. Report is available within 1 hour of campaign window close. Report is downloadable. Contributing scitizen count is anonymized — no individual identification. (Scale: hours | Worst: 4 | Plan: 1 | Best: 0.25)

**Depends on:** FR-010, FR-105

---

### FR-120: Cross-Campaign Data Comparison (0x4e96)

**Priority:** Should | **Originator:** Gap analysis — FR-083 provides in-platform data exploration within a single campaign but no FR covers comparison of data across multiple campaigns

**Description:** The platform shall allow researchers to select multiple campaigns and compare their data side by side. Comparison views shall support overlaying time series from different campaigns, comparing aggregate statistics, and identifying differences in data quality metrics across campaigns.

**Rationale:** Scientific research frequently involves comparing datasets across conditions — same study in different regions, same region in different seasons, or control vs. experimental campaigns. Without cross-campaign comparison, researchers must export data from each campaign and compare externally, reducing the value of the in-platform analytics.

**Fit Criterion:** Researchers can select two or more campaigns they have access to and view a comparison dashboard. The comparison supports overlaid time series charts, side-by-side aggregate statistics, and data quality metric comparison. Campaigns with incompatible parameter types display a clear incompatibility warning. (Scale: boolean | Pass/Fail)

**Depends on:** FR-083

---

## 6h. Certificate Lifecycle Management (BUC-08, scope:0x2f)

### FR-028: Automated Certificate Renewal (0x48)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall support automated certificate renewal via POST /renew where the device presents its current valid certificate via mTLS and a new CSR. Renewal is triggered at day 60 of 90-day cert lifetime.

**Rationale:** Certificates have 90-day lifetime. Automated renewal at day 60 prevents service disruption. Device presents current cert via mTLS + new CSR.

**Fit Criterion:** A device with a valid certificate obtains a new certificate via /renew within 10 seconds. The new certificate has a fresh 90-day validity window. The device registry is updated with the new certificate serial. (Scale: seconds | Worst: 30 | Plan: 10 | Best: 3)

**Cross-ref:** scope:0x2f

---

### FR-029: Grace Period for Expired Certificates (0x44)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall allow devices with certificates expired for up to 7 days to reach the /renew endpoint only. All other endpoints are denied. Devices expired beyond 7 days must re-enroll.

**Rationale:** Devices in the field may miss the renewal window. A 7-day grace period allows recovery without full re-enrollment.

**Fit Criterion:** A device with cert expired 5 days ago can reach /renew and obtain a new cert. A device with cert expired 8 days ago is denied at /renew and must re-enroll. During grace period, only /renew is accessible. (Scale: days | Pass/Fail)

**Depends on:** FR-028
**Cross-ref:** scope:0x2f

---

### FR-030: Certificate Revocation via Registry (0x46)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall revoke device access by updating the device status in the registry to revoked. OPA denies all subsequent actions on the next bundle refresh cycle.

**Rationale:** The platform controls the relying party, so traditional CRL/OCSP is unnecessary. Registry status update + OPA denial achieves revocation within 30 seconds.

**Fit Criterion:** A revoked device is denied all actions within 30 seconds of registry status update. No CRL or OCSP infrastructure is required. (Scale: seconds | Worst: 60 | Plan: 30 | Best: 5)

**Cross-ref:** scope:0x2f

---

## 6i. Device Security Response (BUC-09, scope:0x30)

### FR-031: Bulk Device Suspension by Class (0x4a)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall support bulk suspension of all devices matching a specified device class and firmware version range. Registry batch update, OPA denial on next bundle refresh, mass notification to affected scitizens.

**Rationale:** 50%+ of IoT devices have critical exploitable vulnerabilities (FACT-010). When a firmware vulnerability is discovered, all affected devices must be suspended immediately.

**Fit Criterion:** All devices matching the specified class and firmware version are set to suspended status within 60 seconds. OPA denies their requests on next bundle refresh. Affected scitizens receive notification. (Scale: seconds | Worst: 120 | Plan: 60 | Best: 30)

**Cross-ref:** scope:0x30, facts:0x2f

---

### FR-032: Flag Data from Vulnerable Window (0x4c)

**Priority:** Should | **Originator:** Researcher

**Description:** The platform shall flag all readings submitted by affected devices during the vulnerability window for researcher review. Flagged data is quarantined, not deleted.

**Rationale:** Data submitted by affected devices during the vulnerable window may be compromised and should be flagged for researcher review.

**Fit Criterion:** All readings from affected devices between vulnerability introduction and suspension are flagged and queryable as quarantined. No flagged readings are deleted without explicit researcher action. (Scale: boolean | Pass/Fail)

**Depends on:** FR-031
**Cross-ref:** scope:0x30

---

### FR-033: Device Reinstatement After Patch (0x4e)

**Priority:** Should | **Originator:** Scitizen

**Description:** The platform shall support reinstating suspended devices when they report a firmware version at or above the patched version. Reinstatement restores active status and campaign enrollments.

**Rationale:** Suspension is reversible. Devices that update to a patched firmware version should be reinstatable without full re-enrollment.

**Fit Criterion:** A suspended device reporting patched firmware is reinstated to active status. Campaign enrollments are preserved. OPA allows requests on next bundle refresh. (Scale: boolean | Pass/Fail)

**Depends on:** FR-031
**Cross-ref:** scope:0x30

---

## 6j. Scitizen Recognition and Incentives (BUC-10, scope:0x31)

### FR-034: Compute Contribution Score (0x53)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall compute a contribution score for each scitizen based on volume of accepted readings, data quality rate, consistency of submissions, and diversity of campaigns contributed to.

**Rationale:** Contribution scores reflect participation: volume, quality, consistency, and diversity. They drive recognition and sweepstakes eligibility.

**Fit Criterion:** Contribution score is updated within 15 minutes of a new reading acceptance. Score reflects volume, quality, consistency, and diversity dimensions. Score is visible on the scitizen profile. (Scale: minutes | Worst: 30 | Plan: 15 | Best: 5)

**Cross-ref:** scope:0x31, facts:0x2c

---

### FR-035: Award Badges and Recognition (0x55)

**Priority:** Should | **Originator:** Scitizen

**Description:** The platform shall award badges for milestone achievements (first contribution, campaign completion, quality streaks, geographic diversity). Badges are visible on scitizen profiles and campaign acknowledgments.

**Rationale:** Recognition must not crowd out intrinsic motivation (Self-Determination Theory). Informational rewards (badges, acknowledgments) enhance intrinsic motivation.

**Fit Criterion:** Badges are awarded automatically when milestone criteria are met. At minimum: first contribution, 100 readings, 1000 readings, campaign completion, and 30-day consistency streak. (Scale: badge types | Worst: 3 | Plan: 5 | Best: 10)

**Depends on:** FR-034
**Cross-ref:** scope:0x31, facts:0x11

---

### FR-036: Manage Sweepstakes Entries (0x51)

**Priority:** Should | **Originator:** Scitizen

**Description:** The platform shall grant sweepstakes entries based on contribution activity. Entry accumulation follows a variable-ratio schedule tied to contribution score milestones.

**Rationale:** Prospect Theory: people overweight small probabilities, making lottery-style incentives more motivating per dollar than fixed payments (FACT-008). Variable-ratio reinforcement produces strongest engagement.

**Fit Criterion:** Sweepstakes entries are granted at defined contribution score milestones. Entry count is visible to the scitizen. Entries are auditable and tamper-evident. (Scale: boolean | Pass/Fail)

**Depends on:** FR-034
**Cross-ref:** scope:0x31, facts:0x2c

---

### FR-121: Scitizen Leaderboard (0x4e98)

**Priority:** Should | **Originator:** Gap analysis — FR-034 computes scores and FR-035 awards badges but no FR specifies a public-facing leaderboard for comparing contributions and motivating participation

**Description:** The platform shall provide a leaderboard view showing scitizen contribution rankings. Leaderboards shall be filterable by campaign, organization, time period, and overall. Only scitizens who have opted into public profile visibility (FR-042) shall appear on leaderboards. Leaderboards shall show contribution scores, badge counts, and campaign participation.

**Rationale:** BUC-10 covers scitizen recognition and incentives. Contribution scores (FR-034) and badges (FR-035) are awarded but have no public visibility surface. Citizen science participation is driven partly by social recognition — leaderboards create healthy competition and community visibility.

**Fit Criterion:** A leaderboard view displays ranked scitizens with contribution scores and badge counts. Only users with public profile visibility appear. Leaderboards support filtering by campaign, organization, and time period. Rankings update within 24 hours of new contribution score calculations. (Scale: boolean | Pass/Fail)

**Depends on:** FR-034, FR-042

---

### FR-126: Sweepstakes Entry Visibility (0x4ea6)

**Priority:** Must | **Originator:** Gap analysis — FR-036 manages sweepstakes entries but no FR specifies the scitizen-facing view of their entries, eligibility status, and results

**Description:** Scitizens shall be able to view their sweepstakes entries, including which campaigns qualified them, how many entries they have accumulated, the status of each sweepstakes drawing, and results of past drawings. The view shall clearly show entry criteria and how contributions translate to entries.

**Rationale:** Sweepstakes are an incentive mechanism (BUC-10). Without visibility into entries and results, scitizens cannot see the connection between their contributions and the incentive. Transparency in sweepstakes mechanics builds trust and reinforces the contribution-reward loop that drives participation.

**Fit Criterion:** Scitizens can view a list of their sweepstakes entries with qualifying campaign, entry count, and drawing status. Past drawing results are visible. The entry criteria and contribution-to-entry conversion are explained. The view updates within 24 hours of new entry accumulation. (Scale: boolean | Pass/Fail)

**Depends on:** FR-036

---

### FR-127: Badge Display and Sharing (0x4ea8)

**Priority:** Should | **Originator:** Gap analysis — FR-035 awards badges but no FR specifies where badges are displayed or how scitizens can share their achievements

**Description:** Earned badges shall be displayed on the scitizen profile (visible per FR-042 visibility settings), on the contribution dashboard (FR-040), and in leaderboard views (FR-121). Scitizens shall be able to generate a shareable link or image of their badge collection for use outside the platform.

**Rationale:** Badges are recognition tokens (BUC-10). Recognition has no value if it is not visible. Displaying badges on profiles creates social proof within the community. External sharing extends the platforms reach by letting scitizens showcase participation to peers, employers, and academic networks.

**Fit Criterion:** Earned badges appear on the scitizen profile page and contribution dashboard. Each badge shows its name, description, award date, and qualifying campaign. Scitizens can generate a shareable URL that displays their badge collection to unauthenticated viewers (respecting profile visibility settings). (Scale: boolean | Pass/Fail)

**Depends on:** FR-035, FR-040

---

## 6k. Authentication and Session Security

### FR-044: Login Redirect Flow (0x271f)

**Priority:** Must | **Originator:** Scitizen, Researcher

**Description:** The platform shall redirect unauthenticated users who access a protected page to the Zitadel login flow, preserving the originally requested URL. After successful authentication, the user is returned to the originally requested page, not a generic landing page.

**Rationale:** Users arrive at the platform via deep links (shared campaign URLs, email notification links, bookmarks). Losing the target page after authentication creates friction and confusion.

**Fit Criterion:** An unauthenticated user accessing a protected page is redirected to Zitadel login and returned to the exact originally requested URL after authentication. No protected page returns a 200 to an unauthenticated request. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011

---

### FR-045: Token Refresh and Session Continuity (0x2721)

**Priority:** Must | **Originator:** Scitizen, Researcher

**Description:** The platform shall silently refresh expired access tokens using the Zitadel refresh token grant without requiring the user to re-authenticate. If the refresh token is also expired or revoked, the user is redirected to the login flow. Active sessions persist across browser tabs.

**Rationale:** Forcing re-login on token expiry disrupts workflow, especially for researchers monitoring campaigns. Silent refresh keeps sessions alive as long as the refresh token is valid.

**Fit Criterion:** A user with an expired access token but valid refresh token continues using the platform without re-authenticating. Token refresh completes in under 2 seconds. A user with an expired refresh token is redirected to login. (Scale: seconds | Worst: 5 | Plan: 2 | Best: 0.5)

**Depends on:** FR-011

---

### FR-069: Multi-Factor Authentication Policy (0x4e30)

**Priority:** Must | **Originator:** Research Institution

**Description:** The platform shall define and enforce MFA policies per user type and permission level via Zitadel. Organization admins and users with campaign management, data export, or member management permissions shall require MFA. Scitizens without elevated permissions may use single-factor authentication. Organizations may override the default policy to require MFA for all members. Supported second factors include TOTP authenticator apps and WebAuthn/passkeys.

**Rationale:** Zitadel handles MFA mechanics, but the platform defines the policy. Researchers with access to campaign data, export permissions, or member management represent a higher risk profile than scitizens browsing campaigns. Institutional compliance requirements (IRB, data governance) may mandate MFA for all researchers. WebAuthn/passkeys are the modern standard replacing SMS-based 2FA.

**Fit Criterion:** Users with campaign management, data export, or member management permissions cannot complete login without a second factor. Organizations can enable MFA-required for all members. TOTP and WebAuthn are supported as second factors. Scitizens without elevated permissions can log in with single factor unless org policy overrides. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011, FR-047

---

### FR-070: Session Management (0x4e32)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall enforce session inactivity timeout, allow users to view their active sessions (device, location, last activity), and revoke individual sessions. Sessions expire after a configurable inactivity period (default 30 minutes for researchers, 24 hours for scitizens). Users can log out of all sessions from any active session. Session state is managed via Zitadel tokens with platform-side inactivity tracking.

**Rationale:** FR-045 covers token refresh but not session lifecycle. A user who loses a device or suspects compromise needs to see and revoke active sessions. Different inactivity timeouts for researchers (higher-risk, manage campaign data) vs scitizens (lower-risk, contribute data) reflect the risk profile difference.

**Fit Criterion:** Sessions expire after configured inactivity period. Users can list active sessions showing device type and last activity timestamp. Users can revoke individual sessions or all sessions. Revoked sessions are denied on next request. Session list updates within 60 seconds of new login or revocation. (Scale: seconds | Worst: 120 | Plan: 60 | Best: 10)

**Depends on:** FR-011, FR-045

---

### FR-071: Account Lockout and Brute Force Protection (0x4e34)

**Priority:** Must | **Originator:** Corewood

**Description:** The platform shall configure Zitadel to enforce account lockout after repeated failed authentication attempts. After N consecutive failures (configurable, default 5), the account is temporarily locked for a configurable duration (default 15 minutes). The user is notified of the lockout via email. Failed attempt counts reset after successful authentication or lockout expiry.

**Rationale:** Standard security requirement. No existing FR addresses authentication failure handling. Open registration (BUC-03) means any email can be targeted for brute force. Lockout protects both the user and the platform.

**Fit Criterion:** After 5 consecutive failed login attempts, the account is locked for 15 minutes. The user receives an email notification of the lockout. Login attempts during lockout are rejected with a message indicating temporary lock. Failed count resets after successful login or lockout expiry. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011

---

### FR-079: Reauthentication for Sensitive Operations (0x4e44)

**Priority:** Must | **Originator:** Corewood

**Description:** The platform shall require reauthentication before executing sensitive operations: account deletion, data export, role or permission changes, organization ownership transfer, device bulk operations, and campaign deletion. Reauthentication means the user must re-enter their password or complete a second-factor challenge, even if they have an active session.

**Rationale:** OWASP ASVS requires reauthentication for sensitive operations. An active session that was compromised (session hijacking, unattended workstation) should not grant access to destructive or data-exfiltration operations without confirming the actor is the legitimate user.

**Fit Criterion:** All specified sensitive operations prompt for reauthentication before execution. The reauthentication challenge is the same factor(s) required for initial login. A successful reauthentication is valid for 5 minutes for the specific operation requested. Failure to reauthenticate blocks the operation. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011

---

### FR-101: Security Event Notification (0x4e70)

**Priority:** Must | **Originator:** Corewood

**Description:** The platform shall notify users of security-relevant events on their account: successful login from a new device or location, password change, email change, MFA enrollment or removal, session revocation, account lockout, and account deactivation or deletion request. Notifications are delivered via email to the address on file with the identity provider, regardless of in-app notification preferences.

**Rationale:** OWASP ASVS 2.2.3 requires secure notifications after authentication changes. ASVS 2.5.5 requires notification of authentication factor changes. These are security notifications, not feature notifications — they cannot be disabled by user preference because their purpose is to alert the user to potentially unauthorized account activity.

**Fit Criterion:** All specified security events trigger an email notification to the user. Notifications are sent within 5 minutes of the event. Security notifications cannot be disabled through notification preferences. Each notification includes the event type, timestamp, and device/location context where available. (Scale: minutes | Worst: 10 | Plan: 5 | Best: 1)

**Depends on:** FR-011, FR-051

---

### FR-102: Generic Error Responses for Authentication (0x4e72)

**Priority:** Must | **Originator:** Corewood

**Description:** The platform shall use generic error messages for authentication failures that do not reveal whether an account exists. Login failure, password reset for unknown email, and registration for existing email all return identical response timing and messaging. This prevents user enumeration attacks.

**Rationale:** OWASP Authentication Cheat Sheet explicitly requires generic messaging to prevent user enumeration. Example: password reset says If that email address is in our database, we will send you an email to reset your password regardless of whether the email exists. Open registration (BUC-03) makes the platform a target for credential stuffing; user enumeration makes it worse.

**Fit Criterion:** Login failure for nonexistent and existent accounts produces identical error messages and response times (within 100ms variance). Password reset for unknown email produces the same response as for known email. Registration for existing email does not confirm the email is taken. No authentication endpoint reveals account existence through response content, status code, or timing. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011, FR-037

---

## 6l. Platform Administration

### FR-072: Platform Administrator Role (0x4e36)

**Priority:** Must | **Originator:** Corewood

**Description:** The platform shall define a platform administrator role with cross-organization visibility and management capabilities. Platform admins can view all organizations, campaigns, users, and devices. Platform admin is distinct from organization admin — it is scoped to platform operations, not to any single organization. Platform admin actions are logged in the system audit trail.

**Rationale:** FR-003 defines organization-scoped roles. No requirement covers platform-scoped administration. Corewood as platform owner (Section 1.3) requires operational visibility across all tenants for incident response, abuse detection, and platform health management.

**Fit Criterion:** A platform administrator can view all organizations, campaigns, users, and devices across tenants. Platform admin role cannot be assigned by organization admins — only by existing platform admins. All platform admin actions are logged with admin identity and timestamp. (Scale: boolean | Pass/Fail)

**Depends on:** FR-003

---

### FR-073: Platform Health Dashboard (0x4e38)

**Priority:** Must | **Originator:** Corewood

**Description:** The platform shall provide platform administrators with a health dashboard showing: total registered users (by type), active campaigns (by status), registered devices (by status and class), data ingestion throughput, validation error rates, authentication failure rates, policy engine denial rates, certificate expiration distribution, and storage consumption. Metrics are real-time or near-real-time.

**Rationale:** The glossary defines Platform Health as an observable operational state. No functional requirement specifies how platform operators access this information. Without a health dashboard, Corewood cannot detect system degradation, capacity issues, or security incidents until users report them.

**Fit Criterion:** Platform health dashboard displays all specified metrics. Metrics refresh within 60 seconds. Dashboard is accessible only to platform administrators. Historical trend data is available for at least 90 days. (Scale: seconds | Worst: 120 | Plan: 60 | Best: 15)

**Depends on:** FR-072

---

### FR-074: Platform User Management (0x4e3a)

**Priority:** Must | **Originator:** Corewood

**Description:** The platform shall allow platform administrators to search users by identifier, email, user type, or organization membership; view user profile, activity summary, and organization memberships; and take moderation actions: suspend (temporary, reversible) or ban (permanent) a user account. Suspended users cannot authenticate. Banned users cannot re-register. Moderation actions are logged in the system audit trail and trigger notification to the affected user.

**Rationale:** FR-048 covers org-scoped member management. No requirement addresses platform-wide user management for abuse response, support, or compliance enforcement. Open registration (BUC-03) means bad actors can register. The platform needs tools to respond.

**Fit Criterion:** Platform admins can search users by identifier, email, type, or organization. User detail view shows profile, activity summary, and org memberships. Suspend and ban actions take effect on next authentication attempt. Moderation actions are logged and trigger user notification. Banned email addresses cannot re-register. (Scale: boolean | Pass/Fail)

**Depends on:** FR-072

---

### FR-075: Platform Campaign Oversight (0x4e3c)

**Priority:** Must | **Originator:** Corewood

**Description:** The platform shall provide platform administrators with a cross-organization campaign view: all campaigns by status (draft, published, active, completed, cancelled), organization, data volume, enrolled device count, and quality metrics. Platform admins can flag or suspend campaigns that violate platform policies.

**Rationale:** No requirement provides cross-tenant campaign visibility. A campaign could collect inappropriate data, violate terms of service, or consume disproportionate platform resources. Platform operators need visibility and the ability to intervene.

**Fit Criterion:** Platform admins can list all campaigns across organizations, filtered by status, organization, and date range. Campaign detail shows data volume, device count, and quality metrics. Platform admins can suspend a campaign, which stops data ingestion and notifies the campaign owner. Oversight actions are logged. (Scale: boolean | Pass/Fail)

**Depends on:** FR-010, FR-072

---

### FR-076: Platform Device Fleet Overview (0x4e3f)

**Priority:** Must | **Originator:** Corewood

**Description:** The platform shall provide platform administrators with a fleet-wide device view: device counts by status, class, tier, and firmware version; certificate expiry distribution; devices approaching renewal window; devices in grace period. Platform admins can filter by device class and firmware version to assess exposure during security response (BUC-09).

**Rationale:** FR-031 (bulk suspension) assumes the ability to identify affected devices by class and firmware version. No requirement provides the fleet-level view needed to make that assessment. Without fleet visibility, security response is blind.

**Fit Criterion:** Platform admins can view device fleet aggregated by status, class, tier, and firmware version. Certificate expiry distribution is displayed. Filtering by class and firmware version returns matching devices within 5 seconds for fleets up to 1 million devices. View accessible only to platform administrators. (Scale: seconds | Worst: 15 | Plan: 5 | Best: 2)

**Depends on:** FR-016, FR-072

---

### FR-077: System Audit Log (0x4e41)

**Priority:** Must | **Originator:** Corewood

**Description:** The platform shall maintain an immutable audit log of all significant platform events: user registration, authentication successes and failures, role assignments, campaign creation and state changes, device registration and status changes, data exports, moderation actions, and policy engine denials. Each log entry includes actor identity, action, target, timestamp, and result. The audit log is queryable by platform administrators with filtering by actor, action type, target, and time range.

**Rationale:** FR-049 covers role change audit only. OWASP ASVS requires logging of authentication events and access control failures. Institutional compliance (IRB, data governance) may require audit trails for data access. Without a comprehensive audit log, incident investigation and compliance reporting are impossible.

**Fit Criterion:** All specified event types produce audit log entries. Each entry contains actor, action, target, timestamp, and result. Log entries are immutable after creation. Platform admins can query the log filtered by actor, action type, target, and time range. Query results return within 10 seconds for up to 30 days of log data. Log retention is at least 1 year. (Scale: seconds | Worst: 30 | Plan: 10 | Best: 3)

**Depends on:** FR-072

---

### FR-098: Content Moderation for Campaign Descriptions (0x4e6a)

**Priority:** Should | **Originator:** Corewood

**Description:** The platform shall allow platform administrators to review and moderate user-generated content: campaign descriptions, campaign announcements, and organization display names. Content flagged by users (via FR-094) or automated filters is queued for review. Platform admins can approve, edit, or remove flagged content. Removed content triggers notification to the content author.

**Rationale:** Researchers create campaign descriptions (FR-005) and announcements (FR-053) visible to all scitizens. Organizations set display names (FR-001). Any user-generated content visible to others requires moderation capability. Without it, the platform has no recourse against misleading campaigns, inappropriate announcements, or offensive organization names.

**Fit Criterion:** Flagged content appears in a moderation queue accessible to platform admins. Admins can approve, edit, or remove content. Removed content is replaced with a notice. Content authors are notified of moderation actions. User-reported content (FR-094) feeds into the moderation queue. (Scale: boolean | Pass/Fail)

**Depends on:** FR-072

---

### FR-100: Human-Facing API Rate Limiting (0x4e6e)

**Priority:** Must | **Originator:** Corewood

**Description:** The platform shall enforce rate limiting on all human-facing API endpoints: per-user and per-IP limits for authentication, registration, search, data queries, and export requests. Rate limits are configurable by endpoint category. Exceeded limits return a standard rate-limit response with retry-after indication. Device-facing ingestion endpoints have separate rate limiting defined by FR-022.

**Rationale:** OWASP ASVS 2.2.1 requires rate limiting on authentication endpoints. But rate limiting applies to all human-facing endpoints — not just auth. Without rate limits, a single user or automated scraper can degrade the platform for everyone. FR-022 rate-limits device ingestion per-device; this covers the human-facing surface.

**Fit Criterion:** All human-facing API endpoints enforce rate limits. Rate-limited responses include HTTP 429 status and Retry-After header. Rate limits are configurable per endpoint category. Per-user and per-IP limits are enforced independently. Legitimate usage patterns do not trigger rate limits under normal operation. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011

---

### FR-115: Bulk Device Operations (0x4e8c)

**Priority:** Should | **Originator:** Gap analysis — FR-031 covers bulk suspension for security response but no FR covers general-purpose bulk operations on devices for campaign management or scitizen convenience

**Description:** The platform shall support bulk operations on devices: select multiple devices and apply an action (enroll in campaign, withdraw from campaign, deregister, update metadata) in a single operation. Bulk operations shall provide a preview of affected devices before execution and a summary of results after completion.

**Rationale:** Campaigns with hundreds of enrolled devices and scitizens with multiple devices need efficient management. Performing actions one device at a time does not scale. FR-031 established the pattern for bulk suspension in security contexts; the same pattern applies to routine device management.

**Fit Criterion:** Users can select multiple devices from a list view and apply a supported action. A confirmation screen shows the count and list of affected devices before execution. After execution, a result summary shows successes and failures. Failed operations identify the specific devices and reasons. (Scale: boolean | Pass/Fail)

**Depends on:** FR-041

---

### FR-117: Platform Announcements (0x4e90)

**Priority:** Must | **Originator:** Gap analysis — FR-053 covers campaign announcements from researchers but no FR covers platform-wide announcements from administrators to all users

**Description:** Platform administrators shall be able to create and publish announcements visible to all users or targeted user segments (all researchers, all scitizens, specific organizations). Announcements shall support scheduling for future publication and automatic expiration. Active announcements shall be displayed prominently in the platform UI.

**Rationale:** Platform operations require communicating maintenance windows, policy changes, feature releases, and security advisories to users. Without a platform-level announcement system, administrators must rely on external channels that may not reach all users.

**Fit Criterion:** Platform administrators can create announcements with a title, body, target audience, publish date, and expiration date. Active announcements appear in a banner or notification area visible to targeted users. Expired announcements are automatically hidden. Users can dismiss announcements and they do not reappear. (Scale: boolean | Pass/Fail)

**Depends on:** FR-072

---

### FR-119: Administrative Data Export (0x4e94)

**Priority:** Must | **Originator:** Gap analysis — FR-026 covers campaign data export but no FR covers export of administrative data (member lists, device inventories, audit logs) for compliance and reporting

**Description:** Organization managers and platform administrators shall be able to export administrative data including member lists, device inventories, campaign metadata, and audit logs in structured formats. Exports shall respect access control — users can only export data they are authorized to view. Exported data shall not include identity provider PII unless explicitly authorized by the data subject.

**Rationale:** Compliance requirements (CO-001 GDPR, CO-002 IRB) may require organizations to produce records of platform activity. Institutional reporting and migration needs require structured exports beyond campaign data. Without administrative export, organizations are locked into the platform with no portability of their operational data.

**Fit Criterion:** Organization managers can export member lists, device inventories, and organization audit logs. Platform administrators can export platform-level audit logs and aggregate statistics. All exports produce structured files (CSV or JSON). No export includes PII from the identity provider without explicit per-user consent. (Scale: boolean | Pass/Fail)

**Depends on:** FR-048, FR-077

---

## 6m. Notifications

### FR-051: Notification Delivery System (0x272f)

**Priority:** Must | **Originator:** Scitizen, Researcher

**Description:** The platform shall deliver notifications to users via two channels: in-app (visible in the platform UI) and email (sent to the address on file in Zitadel). Notification types include: device suspended, device certificate expiring, campaign starting, campaign ending, campaign cancelled, badge awarded, invitation received, account status change, and security response alert. Each notification includes a type, timestamp, message body, and optional action link.

**Rationale:** FR-031 mentions affected scitizens receive notification but no notification system is defined. Multiple existing requirements implicitly depend on notification capability (device suspension, badge award, campaign lifecycle). Without a defined delivery system, these promises are unimplementable.

**Fit Criterion:** Notifications are delivered in-app within 60 seconds of the triggering event. Email notifications are sent within 5 minutes. Each notification type defined in the requirement produces both in-app and email delivery. Notifications include type, timestamp, message, and action link where applicable. (Scale: seconds | Worst: 120 | Plan: 60 | Best: 10)

**Depends on:** FR-011

---

### FR-052: Notification Preferences (0x2731)

**Priority:** Should | **Originator:** Scitizen

**Description:** The platform shall allow users to configure notification preferences per channel (in-app, email) and per notification type. Users may disable email notifications for specific types while retaining in-app delivery. Security-critical notifications (device suspension, security response) cannot be disabled.

**Rationale:** Users have different notification tolerances. Some want email for everything; others only want critical alerts. Providing per-type, per-channel control reduces notification fatigue and churn.

**Fit Criterion:** A user who disables email for badge-awarded notifications does not receive email for that type but still sees in-app notification. A user cannot disable security-critical notification types. Preference changes take effect immediately. (Scale: boolean | Pass/Fail)

**Depends on:** FR-051

---

### FR-053: Campaign Announcements (0x272d)

**Priority:** Should | **Originator:** Researcher

**Description:** The platform shall allow researchers with campaign management permission to send announcements to all scitizens enrolled in their campaign. Announcements are delivered via the notification system (in-app and email per user preferences). Announcements include a subject, body text, and optional action link.

**Rationale:** Researchers need to communicate with their campaign participants: updates on campaign progress, changes to collection parameters, gratitude for contributions, or instructions for device recalibration. Without campaign announcements, researchers have no direct communication channel to their contributors.

**Fit Criterion:** A researcher can send an announcement to all enrolled scitizens. The announcement is delivered via the notification system respecting each recipients preferences. Only users with campaign management permission for that campaign can send announcements. (Scale: boolean | Pass/Fail)

**Depends on:** FR-005, FR-051

---

### FR-084: Researcher Campaign Alerting (0x4e4e)

**Priority:** Should | **Originator:** Researcher

**Description:** The platform shall allow researchers to configure alert rules on campaign metrics: data submission rate drops below a threshold, rejection rate exceeds a threshold, temporal gap exceeds a duration, enrolled device count drops below a threshold, or anomaly rate exceeds a threshold. Triggered alerts are delivered via the notification system. Alerts can be acknowledged, snoozed, or disabled.

**Rationale:** FR-010 provides passive monitoring. Researchers running multi-week campaigns cannot continuously watch a dashboard. Alerting converts passive monitoring into active quality management — the researcher is notified when intervention is needed rather than discovering problems during periodic review.

**Fit Criterion:** Researchers can create alert rules with configurable thresholds on specified metrics. Alert conditions are evaluated within 15 minutes of metric update. Triggered alerts are delivered via the notification system. Alerts include the metric value that triggered the condition and a link to the relevant dashboard view. (Scale: minutes | Worst: 30 | Plan: 15 | Best: 5)

**Depends on:** FR-010, FR-051

---

### FR-114: Notification History View (0x4e8a)

**Priority:** Must | **Originator:** Gap analysis — FR-051 delivers notifications and FR-052 manages preferences but no FR specifies where users view their notification history

**Description:** The platform shall provide a notification history view where users can see all past notifications, mark them as read or unread, and filter by notification type. Unread notification count shall be visible from any page. Notifications shall link to the relevant resource (campaign, device, badge) when applicable.

**Rationale:** Notifications are ephemeral without a persistent history view. Users who miss a notification or want to revisit one have no recovery path. A notification inbox is the standard pattern for making the notification system usable beyond real-time delivery.

**Fit Criterion:** Users can access a notification history view showing all notifications received in the last 90 days. Each notification shows timestamp, type, summary, and a link to the related resource. Unread count is visible on every page. Users can mark individual or all notifications as read. (Scale: boolean | Pass/Fail)

**Depends on:** FR-051

---

## 6n. User Experience and Accessibility

### FR-040: Scitizen Contribution Dashboard (0x271d)

**Priority:** Should | **Originator:** Scitizen

**Description:** The platform shall provide each scitizen with a dashboard showing their contribution history: readings submitted per campaign, acceptance/rejection rates, contribution score, badges earned, active campaign enrollments, and sweepstakes entries. All data displayed is application state from the app DB.

**Rationale:** Scitizens need visibility into their contributions to stay engaged. Clear feedback that their data matters is a key need (Section 1.4). The contribution dashboard surfaces application state only, no PII.

**Fit Criterion:** The dashboard displays per-campaign reading counts, acceptance rate, current contribution score, earned badges, and active enrollments. Data is refreshed within 15 minutes of the latest reading submission. (Scale: minutes | Worst: 30 | Plan: 15 | Best: 5)

**Depends on:** FR-011, FR-034

---

### FR-041: Device Management Dashboard (0x2717)

**Priority:** Should | **Originator:** Scitizen

**Description:** The platform shall provide each scitizen with a device management view showing all registered devices, their current status (active, suspended, revoked, pending), certificate expiry date, firmware version, campaign enrollments per device, and last-seen timestamp.

**Rationale:** Scitizens manage physical devices in the field. They need a single view of device health, status, and enrollment to troubleshoot issues and manage their fleet.

**Fit Criterion:** The device management view lists all devices owned by the scitizen with status, certificate expiry, firmware version, campaign enrollments, and last-seen timestamp. Device status changes are reflected within 60 seconds. (Scale: seconds | Worst: 120 | Plan: 60 | Best: 10)

**Depends on:** FR-011, FR-016

---

### FR-080: Terms of Service and Privacy Policy Acceptance (0x4e46)

**Priority:** Must | **Originator:** Corewood

**Description:** The platform shall require users to accept the terms of service and privacy policy at registration. Acceptance is recorded with the document version, timestamp, and user identity. When terms are updated, users are prompted to re-accept on next login. Users who do not accept updated terms cannot continue using the platform until acceptance is recorded.

**Rationale:** Legal baseline for any platform. CO-001 (GDPR) requires informed consent for data processing. Without versioned acceptance records, the platform has no legal basis for data handling and cannot demonstrate compliance.

**Fit Criterion:** Registration requires ToS and privacy policy acceptance before account activation. Acceptance records include document version, timestamp, and user identity. Updated terms trigger re-acceptance prompt on next login. Users who decline updated terms are blocked from platform use until accepted. Acceptance records are retained for the lifetime of the account plus the legally required retention period. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011

---

### FR-081: Scitizen Reading History (0x4e49)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall allow scitizens to view their own submitted readings per campaign and per device: individual reading values with timestamps, acceptance or rejection status, and rejection reasons where applicable. Readings are the scitizens own data — displayed without pseudonymization since they are the data owner. Filterable by campaign, device, time range, and acceptance status.

**Rationale:** FR-040 shows aggregate contribution metrics but no individual reading detail. FR-022 rejects readings with explicit reasons, but those reasons are only useful if the scitizen can see them. A scitizen whose readings are being rejected needs to diagnose whether the device is miscalibrated, geolocating incorrectly, or timestamping with drift. Without per-reading feedback, the scitizen cannot self-correct and will churn.

**Fit Criterion:** Scitizens can view submitted readings filtered by campaign, device, time range, and acceptance status. Each reading shows value, timestamp, acceptance/rejection status, and rejection reason if rejected. Paginated to 100 readings per page. History available within 5 minutes of submission. (Scale: minutes | Worst: 15 | Plan: 5 | Best: 1)

**Depends on:** FR-023, FR-040

---

### FR-085: Scitizen Data Portability (0x4e50)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall allow scitizens to export their own data in a machine-readable format: all submitted readings (with acceptance status), device registry entries, campaign enrollment history, contribution scores, badges earned, and consent records. This is the scitizens personal data export, distinct from researcher campaign data export (FR-026).

**Rationale:** CO-001 requires GDPR-compliant data portability. The fit criterion for CO-001 states data portability export is available in machine-readable format. No functional requirement specifies the mechanism. This is a user right, not a researcher feature.

**Fit Criterion:** Scitizens can request a personal data export from their account settings. Export includes all submitted readings, device entries, enrollment history, scores, badges, and consent records. Export is delivered in a machine-readable format within 24 hours of request. Export is available for download for 7 days. (Scale: hours | Worst: 48 | Plan: 24 | Best: 1)

**Depends on:** FR-011

---

### FR-088: Campaign Search (0x4e56)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall provide full-text search across published campaigns: campaign name, description, parameter types, and geographic region names. Search results are ranked by relevance and can be further filtered by the structured criteria in FR-012 (proximity, sensor type, status). Search is available to all authenticated users.

**Rationale:** FR-012 provides structured filtering (proximity, sensor type, status). But a scitizen may not know the structured terms — they search by natural language description of what they want to contribute to. Search is the unstructured discovery path; browse is the structured filtering path. Both are needed.

**Fit Criterion:** Full-text search returns relevant campaigns within 2 seconds. Search matches against campaign name, description, parameter types, and region names. Results are ranked by relevance. Search results can be further filtered by structured criteria (proximity, sensor type, status). Empty results show helpful suggestions. (Scale: seconds | Worst: 5 | Plan: 2 | Best: 0.5)

**Depends on:** FR-009, FR-012

---

### FR-089: User Activity Log (0x4e58)

**Priority:** Should | **Originator:** Scitizen

**Description:** The platform shall provide each user with a personal activity log showing their own significant actions: login events, device registrations and status changes, campaign enrollments and withdrawals, data exports, profile and settings changes, and role assignments received. The log is read-only and accessible from account settings.

**Rationale:** Standard platform transparency feature. Users need to verify their own account activity for security (detecting unauthorized access) and for personal recordkeeping. This is the user-facing subset of the system audit log (FR-077) — users see only their own actions.

**Fit Criterion:** Users can view their activity log from account settings. Log shows login events, device actions, campaign actions, export actions, profile changes, and role changes. Each entry includes action, timestamp, and result. Log entries appear within 5 minutes of the action. Log is paginated and sortable by date. (Scale: minutes | Worst: 15 | Plan: 5 | Best: 1)

**Depends on:** FR-011

---

### FR-091: Accessibility Compliance (0x4e5c)

**Priority:** Must | **Originator:** Corewood

**Description:** All user-facing platform interfaces shall conform to WCAG 2.1 Level AA. This includes keyboard navigation, screen reader compatibility, sufficient color contrast, text alternatives for non-text content, and accessible form controls. Geographic visualizations shall provide non-visual alternatives (data tables).

**Rationale:** Standard web accessibility requirement. Open-source platforms serving a global user base must be accessible. Many research institutions have institutional accessibility mandates. Geographic visualizations (FR-068) and dashboards require specific attention to non-visual alternatives.

**Fit Criterion:** All user-facing pages pass WCAG 2.1 Level AA automated checks. Keyboard-only navigation can reach all interactive elements. Screen readers can access all content and controls. Color contrast ratios meet AA minimums (4.5:1 for normal text, 3:1 for large text). Geographic visualizations have data table alternatives. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011

---

### FR-093: Device Connection History (0x4e60)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall provide scitizens with a connection history per device: successful connections, failed connection attempts with reasons (certificate expired, device suspended, policy denial), disconnection events, and last-will-and-testament triggers. The glossary defines Connection Failure Diagnostics as surfaced to citizen scientists as human-readable messages on their dashboard.

**Rationale:** The glossary explicitly states that connection failure diagnostics are surfaced to citizen scientists on their dashboard. FR-041 shows device status but not connection history. SEC-006 defines structured diagnostic logs but does not specify the scitizen-facing view. Without connection history, a scitizen whose device silently stops connecting has no way to diagnose the problem.

**Fit Criterion:** Device management view includes connection history showing successful connections, failures with human-readable reasons, disconnections, and last-will triggers. History retains at least 30 days of events. Events appear within 5 minutes of occurrence. Failure reasons are actionable (e.g., Certificate expired on date - renew your device). (Scale: minutes | Worst: 15 | Plan: 5 | Best: 1)

**Depends on:** FR-041

---

### FR-094: User Feedback and Issue Reporting (0x4e62)

**Priority:** Should | **Originator:** Scitizen

**Description:** The platform shall provide a feedback and issue reporting mechanism accessible from any page. Users can report bugs, request features, or flag content issues. Reports include a category, description, and optional screenshot. Reports are visible to platform administrators and routed to the appropriate team. Users receive acknowledgment and can track report status.

**Rationale:** Standard platform feature. An open-source platform serving citizen scientists needs a low-friction feedback channel. Without it, issues go unreported or surface only as user churn. The open-source community model (CON-001) depends on user feedback for improvement.

**Fit Criterion:** Feedback form is accessible from every authenticated page. Reports include category, description, and optional attachment. Submitted reports are visible in the platform admin interface. Users receive acknowledgment and can view their report status. Reports are queryable by category and status. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011

---

### FR-104: Scitizen Onboarding Flow (0x4e76)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall guide newly registered scitizens through a first-time onboarding flow: welcome message, prompt to register their first device, prompt to browse and enroll in their first campaign, and explanation of how contribution scoring and recognition work. The onboarding flow is skippable but resurfaces key incomplete steps until the scitizen has registered at least one device and enrolled in at least one campaign.

**Rationale:** BUC-03 business rules state: minimize friction at registration — the first contribution is the critical retention inflection point (FACT-006). 82% retention post-first-contribution vs 39.7% without. The onboarding flow bridges the gap between account creation and first contribution. Without it, new scitizens land on an empty dashboard with no guidance.

**Fit Criterion:** New scitizens see an onboarding flow after first login. Flow guides toward device registration and campaign enrollment. Flow is skippable. Incomplete onboarding steps are surfaced on subsequent logins until device registration and campaign enrollment are complete. Onboarding completion is tracked per user. (Scale: boolean | Pass/Fail)

**Depends on:** FR-011

---

### FR-108: Researcher Onboarding Flow (0x4e7e)

**Priority:** Must | **Originator:** Gap analysis — BUC-01 covers institutional onboarding but no FR specifies the guided experience for individual researchers joining an organization

**Description:** When a researcher accepts an invitation to join an organization, the platform shall present a guided onboarding flow that walks them through completing their profile in the identity provider, understanding their assigned role and permissions, and accessing the campaign management tools available to them.

**Rationale:** Institutional onboarding (FR-001) creates the organization. Role assignment (FR-003) grants permissions. But the individual researcher experience between invitation acceptance and productive use is unspecified. Without guided onboarding, researchers must discover capabilities by trial and error, increasing time-to-value and support burden.

**Fit Criterion:** A newly invited researcher who accepts an invitation is presented with a step-by-step onboarding flow. The flow completes when the researcher has a populated profile, understands their role, and has accessed at least one campaign management view. Completion rate and drop-off points are measurable. (Scale: boolean | Pass/Fail)

**Depends on:** FR-003, FR-048

---

### FR-110: Timezone Display Normalization (0x4e82)

**Priority:** Must | **Originator:** Gap analysis — all timestamps stored as UTC but no FR specifies how users see dates and times in their local timezone

**Description:** The platform shall store all timestamps in UTC and display them to users in their configured or detected timezone. Users shall be able to set a preferred timezone in their profile. Campaign time windows, reading timestamps, device last-seen times, and all other temporal displays shall reflect the users timezone preference.

**Rationale:** Citizen science campaigns span multiple timezones. Researchers in one timezone set campaign windows that scitizens in other timezones must interpret correctly. Displaying raw UTC to non-technical users causes confusion about when campaigns start and end, when readings were taken, and when devices last reported.

**Fit Criterion:** All user-facing timestamps display in the users configured timezone. A user who changes their timezone preference sees all timestamps update immediately. Campaign time windows display correctly for users in different timezones viewing the same campaign. (Scale: boolean | Pass/Fail)

**Depends on:** FR-045

---

### FR-111: Error State Pages (0x4e84)

**Priority:** Must | **Originator:** Gap analysis — no FR specifies user-facing behavior for error conditions (not found, server error, maintenance, expired session)

**Description:** The platform shall present user-friendly error pages for common error states: resource not found, server error, scheduled maintenance, and expired or invalid sessions. Error pages shall not expose internal system details. Each error page shall provide navigation back to a known good state.

**Rationale:** Unhandled error states result in raw framework error pages that confuse users, expose implementation details, and provide no recovery path. Defining error state behavior as a requirement ensures consistent UX during failures and prevents information leakage.

**Fit Criterion:** All HTTP error responses (4xx, 5xx) render a branded error page with a human-readable message and a link to return to the dashboard or home page. No error page exposes stack traces, internal paths, or framework identifiers. Maintenance mode displays an estimated return time when available. (Scale: boolean | Pass/Fail)


---

### FR-112: List Pagination and Sorting (0x4e86)

**Priority:** Must | **Originator:** Gap analysis — multiple FRs produce list views (campaigns, devices, readings, members) but no FR specifies pagination, sorting, or filtering behavior

**Description:** All list views in the platform shall support pagination, sorting by relevant columns, and filtering by key attributes. The platform shall provide consistent pagination controls across all list views. Default sort order shall be most-recent-first for time-ordered lists and alphabetical for name-ordered lists.

**Rationale:** A researcher managing dozens of campaigns, a scitizen with multiple devices, an organization with hundreds of members — all produce lists that exceed a single screen. Without pagination and sorting, users cannot efficiently find what they need. Consistent behavior across list views reduces cognitive load.

**Fit Criterion:** Every list view with more than 25 items displays pagination controls. Users can sort by at least two columns per list view. Pagination state persists when the user navigates away and returns within the same session. Filter criteria are clearable with a single action. (Scale: boolean | Pass/Fail)


---

### FR-113: In-App Help and Documentation (0x4e88)

**Priority:** Must | **Originator:** Gap analysis — no FR specifies how users access help, tooltips, or documentation within the platform

**Description:** The platform shall provide contextual help accessible from every major view. Help content shall include tooltips for technical terms, guided walkthroughs for complex workflows (device enrollment, campaign creation), and links to external documentation. Help content shall be role-appropriate — scitizens see different help than researchers.

**Rationale:** Citizen science platforms serve non-technical users by definition. Scitizens enrolling devices or interpreting campaign instructions need inline guidance. Without in-app help, support burden falls entirely on human channels, which does not scale with an open-source community platform.

**Fit Criterion:** Every major workflow (device enrollment, campaign creation, campaign enrollment) has at least one guided walkthrough. Technical terms in the UI have tooltip definitions. Help content is accessible within two clicks from any view. Help content differs based on user role. (Scale: boolean | Pass/Fail)


---

## 6o. Campaign Monitoring and Quality

### FR-066: Campaign Device Performance Breakdown (0x4e2a)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall provide researchers with a per-device breakdown within each campaign: acceptance rate, rejection rate by reason, reading volume over time, last submission timestamp, and device health status. Device identity is pseudonymized — researchers see device class, tier, and performance metrics but not owner identity or raw device ID.

**Rationale:** FR-010 provides campaign-level aggregate metrics. But a researcher diagnosing data quality issues needs to identify whether problems are systemic (all devices) or isolated (one malfunctioning device class). Per-device granularity is essential for quality management. Device pseudonymization required per FR-027 and SEC-004.

**Fit Criterion:** Campaign dashboard includes a device-level view showing pseudonymized device ID, device class, acceptance rate, rejection reasons, submission volume timeline, and last-seen timestamp. No researcher can see the scitizen owner or raw device ID. Data refreshed within 5 minutes of latest submission. (Scale: minutes | Worst: 15 | Plan: 5 | Best: 1)

**Depends on:** FR-010, FR-027

---

### FR-067: Campaign Temporal Coverage Analysis (0x4e2c)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall display temporal coverage for each campaign: reading submission rate over time (hourly/daily/weekly aggregations), identification of temporal gaps (periods with zero or below-threshold submissions), and comparison of actual collection rate against campaign target rate if defined.

**Rationale:** FR-010 mentions temporal coverage but does not specify gap detection or rate analysis. A researcher needs to know if data collection has stalled, if certain hours of day are under-represented, or if the campaign is on track to meet its data volume goals. Temporal gaps are a data quality signal that may require intervention (announcements to scitizens, adjusted parameters).

**Fit Criterion:** Campaign dashboard shows submission rate at hourly, daily, and weekly granularity. Periods with zero submissions exceeding 1 hour are flagged as gaps. If the campaign defines a target submission rate, actual vs target comparison is displayed. Data refreshed within 5 minutes. (Scale: minutes | Worst: 15 | Plan: 5 | Best: 1)

**Depends on:** FR-010

---

### FR-068: Campaign Geographic Coverage Visualization (0x4e2f)

**Priority:** Should | **Originator:** Researcher

**Description:** The platform shall display geographic coverage for each campaign: a visualization showing where readings have been submitted relative to the campaign region, identification of under-covered sub-regions, and submission density. Location data in this view is aggregated to the spatial resolution configured for the campaign — individual reading locations are not displayed at full precision.

**Rationale:** FR-010 mentions geographic distribution but does not specify visualization, under-coverage detection, or spatial resolution controls. A researcher needs to see whether their campaign region is evenly covered or has blind spots requiring targeted scitizen recruitment. Aggregation to campaign spatial resolution respects SEC-004 and FR-027 privacy constraints.

**Fit Criterion:** Campaign dashboard shows a geographic visualization with reading density within the campaign region. Sub-regions with zero readings are visually distinguishable from covered areas. Location data is aggregated to the campaign-configured spatial resolution. No individual reading geolocation is displayed at full precision. (Scale: boolean | Pass/Fail)

**Depends on:** FR-007, FR-010

---

## 6p. Requirement Summary

| Req ID | Name | Section | Priority | Originator |
|--------|------|---------|----------|-----------|
| FR-001 | Create Organization Tenant | 6a (BUC-01) | Must | Research Institution |
| FR-002 | Configure Organization Hierarchy | 6a (BUC-01) | Must | Research Institution |
| FR-003 | Define and Assign Roles | 6a (BUC-01) | Must | Research Institution |
| FR-004 | Invite and Onboard Researchers | 6a (BUC-01) | Must | Research Institution |
| FR-005 | Create Campaign | 6b (BUC-02) | Must | Researcher |
| FR-006 | Define Campaign Parameters | 6b (BUC-02) | Must | Researcher |
| FR-007 | Define Campaign Region | 6b (BUC-02) | Must | Researcher |
| FR-008 | Define Campaign Time Window | 6b (BUC-02) | Must | Researcher |
| FR-009 | Publish and Discover Campaigns | 6b (BUC-02) | Must | Researcher |
| FR-010 | Monitor Campaign Data Quality | 6b (BUC-02) | Should | Researcher |
| FR-011 | Scitizen Account Registration | 6c (BUC-03) | Must | Scitizen |
| FR-012 | Browse Campaigns | 6c (BUC-03) | Must | Scitizen |
| FR-013 | Generate Enrollment Code | 6d (BUC-04) | Must | Scitizen |
| FR-014 | Direct Device Enrollment (Tier 1) | 6d (BUC-04) | Must | Scitizen |
| FR-015 | Proxy Device Enrollment (Tier 2) | 6d (BUC-04) | Must | Scitizen |
| FR-016 | Device Registry Entry | 6d (BUC-04) | Must | Scitizen |
| FR-017 | Enroll Device in Campaign | 6e (BUC-05) | Must | Scitizen |
| FR-018 | Multi-Campaign Device Enrollment | 6e (BUC-05) | Must | Scitizen |
| FR-019 | Device Eligibility Check | 6e (BUC-05) | Must | Researcher |
| FR-020 | Authenticate Device via mTLS | 6f (BUC-06) | Must | Researcher |
| FR-021 | Authorize Device Actions via OPA | 6f (BUC-06) | Must | Researcher |
| FR-022 | Validate Sensor Readings | 6f (BUC-06) | Must | Researcher |
| FR-023 | Persist Valid Readings with Provenance | 6f (BUC-06) | Must | Researcher |
| FR-024 | Topic ACL Enforcement | 6f (BUC-06) | Must | Researcher |
| FR-025 | Anomaly Flagging | 6f (BUC-06) | Should | Researcher |
| FR-026 | Export Campaign Data | 6g (BUC-07) | Must | Researcher |
| FR-027 | Separate Contributor Identity from Data | 6g (BUC-07) | Must | Oversight Body |
| FR-028 | Automated Certificate Renewal | 6h (BUC-08) | Must | Scitizen |
| FR-029 | Grace Period for Expired Certificates | 6h (BUC-08) | Must | Scitizen |
| FR-030 | Certificate Revocation via Registry | 6h (BUC-08) | Must | Researcher |
| FR-031 | Bulk Device Suspension by Class | 6i (BUC-09) | Must | Researcher |
| FR-032 | Flag Data from Vulnerable Window | 6i (BUC-09) | Should | Researcher |
| FR-033 | Device Reinstatement After Patch | 6i (BUC-09) | Should | Scitizen |
| FR-034 | Compute Contribution Score | 6j (BUC-10) | Must | Scitizen |
| FR-035 | Award Badges and Recognition | 6j (BUC-10) | Should | Scitizen |
| FR-036 | Manage Sweepstakes Entries | 6j (BUC-10) | Should | Scitizen |
| FR-037 | Password Reset Delegation | 6c (BUC-03) | Must | Scitizen |
| FR-038 | Account Deactivation | 6c (BUC-03) | Should | Scitizen |
| FR-039 | Account Deletion and Right to Erasure | 6c (BUC-03) | Must | Scitizen |
| FR-040 | Scitizen Contribution Dashboard | 6n (Cross-cutting) | Should | Scitizen |
| FR-041 | Device Management Dashboard | 6n (Cross-cutting) | Should | Scitizen |
| FR-042 | Profile Visibility Settings | 6c (BUC-03) | Could | Scitizen |
| FR-043 | Profile Data Resolution from IdP | 6c (BUC-03) | Must | Scitizen, Researcher |
| FR-044 | Login Redirect Flow | 6k (Cross-cutting) | Must | Scitizen, Researcher |
| FR-045 | Token Refresh and Session Continuity | 6k (Cross-cutting) | Must | Scitizen, Researcher |
| FR-046 | Role Hierarchy and Inheritance | 6a (BUC-01) | Must | Research Institution |
| FR-047 | Permission Model Definition | 6a (BUC-01) | Must | Research Institution |
| FR-048 | Organization Member Management | 6a (BUC-01) | Must | Research Institution |
| FR-049 | Audit Trail for Role Changes | 6a (BUC-01) | Should | Research Institution |
| FR-050 | User Type Modification | 6c (BUC-03) | Should | Scitizen, Researcher |
| FR-051 | Notification Delivery System | 6m (Cross-cutting) | Must | Scitizen, Researcher |
| FR-052 | Notification Preferences | 6m (Cross-cutting) | Should | Scitizen |
| FR-053 | Campaign Announcements | 6m (Cross-cutting) | Should | Researcher |
| FR-054 | Campaign Editing Constraints | 6b (BUC-02) | Must | Researcher |
| FR-055 | Campaign Cancellation and Archival | 6b (BUC-02) | Must | Researcher |
| FR-056 | Campaign Duplication | 6b (BUC-02) | Could | Researcher |
| FR-057 | Campaign Collaboration | 6b (BUC-02) | Should | Researcher |
| FR-058 | Device Deregistration | 6d (BUC-04) | Must | Scitizen |
| FR-059 | Device Ownership Transfer | 6d (BUC-04) | Should | Scitizen |
| FR-060 | Device Metadata Update | 6d (BUC-04) | Must | Scitizen |
| FR-061 | Device Health Monitoring | 6d (BUC-04) | Should | Scitizen |
| FR-062 | Campaign Configuration Push | 6e (BUC-05) | Must | Researcher |
| FR-063 | Campaign Enrollment Withdrawal | 6e (BUC-05) | Must | Scitizen |
| FR-064 | Consent Capture at Campaign Enrollment | 6e (BUC-05) | Must | Oversight Body |
| FR-065 | Quarantine Review Workflow | 6f (BUC-06) | Should | Researcher |
| FR-066 | Campaign Device Performance Breakdown | 6o (Cross-cutting) | Must | Researcher |
| FR-067 | Campaign Temporal Coverage Analysis | 6o (Cross-cutting) | Must | Researcher |
| FR-068 | Campaign Geographic Coverage Visualization | 6o (Cross-cutting) | Should | Researcher |
| FR-069 | Multi-Factor Authentication Policy | 6k (Cross-cutting) | Must | Research Institution |
| FR-070 | Session Management | 6k (Cross-cutting) | Must | Scitizen |
| FR-071 | Account Lockout and Brute Force Protection | 6k (Cross-cutting) | Must | Corewood |
| FR-072 | Platform Administrator Role | 6l (Cross-cutting) | Must | Corewood |
| FR-073 | Platform Health Dashboard | 6l (Cross-cutting) | Must | Corewood |
| FR-074 | Platform User Management | 6l (Cross-cutting) | Must | Corewood |
| FR-075 | Platform Campaign Oversight | 6l (Cross-cutting) | Must | Corewood |
| FR-076 | Platform Device Fleet Overview | 6l (Cross-cutting) | Must | Corewood |
| FR-077 | System Audit Log | 6l (Cross-cutting) | Must | Corewood |
| FR-078 | Email Verification | 6c (BUC-03) | Must | Corewood |
| FR-079 | Reauthentication for Sensitive Operations | 6k (Cross-cutting) | Must | Corewood |
| FR-080 | Terms of Service and Privacy Policy Acceptance | 6n (Cross-cutting) | Must | Corewood |
| FR-081 | Scitizen Reading History | 6n (Cross-cutting) | Must | Scitizen |
| FR-082 | Campaign Detail View | 6b (BUC-02) | Must | Scitizen |
| FR-083 | Researcher In-Platform Data Explorer | 6g (BUC-07) | Must | Researcher |
| FR-084 | Researcher Campaign Alerting | 6m (Cross-cutting) | Should | Researcher |
| FR-085 | Scitizen Data Portability | 6n (Cross-cutting) | Must | Scitizen |
| FR-086 | Organization Dashboard | 6a (BUC-01) | Should | Research Institution |
| FR-087 | Programmatic API Access for Researchers | 6g (BUC-07) | Should | Researcher |
| FR-088 | Campaign Search | 6n (Cross-cutting) | Must | Scitizen |
| FR-089 | User Activity Log | 6n (Cross-cutting) | Should | Scitizen |
| FR-090 | Data Retention Policy | 6g (BUC-07) | Must | Corewood |
| FR-091 | Accessibility Compliance | 6n (Cross-cutting) | Must | Corewood |
| FR-092 | Campaign Enrollment Funnel | 6e (BUC-05) | Should | Researcher |
| FR-093 | Device Connection History | 6n (Cross-cutting) | Must | Scitizen |
| FR-094 | User Feedback and Issue Reporting | 6n (Cross-cutting) | Should | Scitizen |
| FR-095 | Scitizen Campaign Progress View | 6e (BUC-05) | Should | Scitizen |
| FR-096 | Invitation Management | 6a (BUC-01) | Must | Research Institution |
| FR-097 | Password Change | 6c (BUC-03) | Must | Scitizen |
| FR-098 | Content Moderation for Campaign Descriptions | 6l (Cross-cutting) | Should | Corewood |
| FR-099 | Organization Activity Audit Trail | 6a (BUC-01) | Should | Research Institution |
| FR-100 | Human-Facing API Rate Limiting | 6l (Cross-cutting) | Must | Corewood |
| FR-101 | Security Event Notification | 6k (Cross-cutting) | Must | Corewood |
| FR-102 | Generic Error Responses for Authentication | 6k (Cross-cutting) | Must | Corewood |
| FR-103 | Public Campaign Showcase | 6b (BUC-02) | Should | Corewood |
| FR-104 | Scitizen Onboarding Flow | 6n (Cross-cutting) | Must | Scitizen |
| FR-105 | Campaign Lifecycle State Machine | 6b (BUC-02) | Must | Researcher |
| FR-106 | Campaign Data Visualization | 6g (BUC-07) | Should | Researcher |
| FR-107 | Campaign Completion Summary | 6g (BUC-07) | Should | Researcher |
| FR-108 | Researcher Onboarding Flow | 6n (Cross-cutting) | Must | Gap analysis — BUC-01 covers institutional onboarding but no FR specifies the guided experience for individual researchers joining an organization |
| FR-109 | Device Certificate Expiry Warning | 6d (BUC-04) | Must | Gap analysis — FR-028 handles automated renewal but no FR specifies the scitizen-facing notification when a certificate approaches expiry or renewal fails |
| FR-110 | Timezone Display Normalization | 6n (Cross-cutting) | Must | Gap analysis — all timestamps stored as UTC but no FR specifies how users see dates and times in their local timezone |
| FR-111 | Error State Pages | 6n (Cross-cutting) | Must | Gap analysis — no FR specifies user-facing behavior for error conditions (not found, server error, maintenance, expired session) |
| FR-112 | List Pagination and Sorting | 6n (Cross-cutting) | Must | Gap analysis — multiple FRs produce list views (campaigns, devices, readings, members) but no FR specifies pagination, sorting, or filtering behavior |
| FR-113 | In-App Help and Documentation | 6n (Cross-cutting) | Must | Gap analysis — no FR specifies how users access help, tooltips, or documentation within the platform |
| FR-114 | Notification History View | 6m (Cross-cutting) | Must | Gap analysis — FR-051 delivers notifications and FR-052 manages preferences but no FR specifies where users view their notification history |
| FR-115 | Bulk Device Operations | 6l (Cross-cutting) | Should | Gap analysis — FR-031 covers bulk suspension for security response but no FR covers general-purpose bulk operations on devices for campaign management or scitizen convenience |
| FR-116 | Campaign Enrollment Approval | 6e (BUC-05) | Should | Gap analysis — FR-017 covers enrollment and FR-019 covers eligibility check but no FR allows researchers to manually approve or reject enrollment requests for campaigns requiring human vetting |
| FR-117 | Platform Announcements | 6l (Cross-cutting) | Must | Gap analysis — FR-053 covers campaign announcements from researchers but no FR covers platform-wide announcements from administrators to all users |
| FR-118 | Device Firmware Compatibility Check | 6d (BUC-04) | Should | Gap analysis — FR-019 checks eligibility by sensor type and region but no FR covers firmware version compatibility between campaign requirements and device capabilities |
| FR-119 | Administrative Data Export | 6l (Cross-cutting) | Must | Gap analysis — FR-026 covers campaign data export but no FR covers export of administrative data (member lists, device inventories, audit logs) for compliance and reporting |
| FR-120 | Cross-Campaign Data Comparison | 6g (BUC-07) | Should | Gap analysis — FR-083 provides in-platform data exploration within a single campaign but no FR covers comparison of data across multiple campaigns |
| FR-121 | Scitizen Leaderboard | 6j (BUC-10) | Should | Gap analysis — FR-034 computes scores and FR-035 awards badges but no FR specifies a public-facing leaderboard for comparing contributions and motivating participation |
| FR-122 | Companion Device Setup Flow | 6d (BUC-04) | Must | Gap analysis — FR-015 defines proxy enrollment but no FR specifies the guided setup experience for pairing a companion device (phone) with a sensor and verifying data flow |
| FR-123 | Campaign Data Quality Thresholds | 6f (BUC-06) | Must | Gap analysis — FR-010 monitors data quality and FR-084 provides alerting but no FR specifies user-configurable quality thresholds that trigger automated responses |
| FR-124 | Device Location Verification | 6f (BUC-06) | Should | Gap analysis — FR-007 defines campaign regions and FR-019 checks eligibility but no FR specifies how device location claims are verified against actual placement |
| FR-125 | Reading Rejection Feedback | 6f (BUC-06) | Must | Gap analysis — FR-022 validates readings and FR-025 flags anomalies but no FR specifies what feedback the scitizen receives when their readings are rejected or flagged |
| FR-126 | Sweepstakes Entry Visibility | 6j (BUC-10) | Must | Gap analysis — FR-036 manages sweepstakes entries but no FR specifies the scitizen-facing view of their entries, eligibility status, and results |
| FR-127 | Badge Display and Sharing | 6j (BUC-10) | Should | Gap analysis — FR-035 awards badges but no FR specifies where badges are displayed or how scitizens can share their achievements |

**MoSCoW Summary:** 88 Must, 37 Should, 2 Could, 0 Won't

---

*Next: [Section 7 — Look and Feel Requirements](./07_look_and_feel.md)*
