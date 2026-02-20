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

**Constrained by:** CON-003 — No Shared Context Exists (0xc)
**Cross-ref:** scope:0x28, facts:0xe

---

### FR-002: Configure Organization Hierarchy (0x12)

**Priority:** Must | **Originator:** Research Institution

**Description:** The platform shall support nested sub-organizations within a tenant, allowing institutions to model departments, labs, and divisions as a hierarchy.

**Rationale:** Institutions have departments, labs, and divisions. The hierarchy must mirror real organizational structure for authorization scoping.

**Fit Criterion:** Sub-organizations can be nested to at least 5 levels. Each sub-organization inherits parent authorization rules unless explicitly overridden. (Scale: nesting levels | Worst: 3 | Plan: 5 | Best: 10)

**Constrained by:** CON-003 — No Shared Context Exists (0xc)
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

**Description:** The platform shall allow organization admins to invite researchers by email. Invited researchers authenticate via Zitadel and are associated with the inviting organization. The researcher must already have a registered account (FR-011) with `user_type` including `researcher`. The invitation creates the organization association, not the user account.

**Rationale:** Researchers must be invited to an organization to create campaigns. Registration (FR-011) and org association are separate concerns: a researcher registers once and may be invited to multiple organizations. Onboarding requires no integration with institutional systems (CON-003).

**Fit Criterion:** An invited researcher with an existing account can accept the invitation and access their organization within 2 minutes, with zero changes to institutional systems. A user without `researcher` in their `user_type` cannot be invited. (Scale: minutes | Worst: 5 | Plan: 2 | Best: 1)

**Constrained by:** CON-003 — No Shared Context Exists (0xc)
**Depends on:** FR-003, FR-011
**Cross-ref:** scope:0x28, facts:0xc

---

## 6b. Campaign Creation and Management (BUC-02, scope:0x29)

### FR-005: Create Campaign (0x16)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall allow an authenticated researcher with campaign creation permission to create a campaign specifying parameters, geographic region, time window, quality thresholds, and device eligibility criteria.

**Rationale:** Campaigns are the central organizing unit. A researcher must define what data is needed, where, when, and to what quality standard.

**Fit Criterion:** A campaign is created with all required fields (parameters, region, window, thresholds) and is queryable via API. Campaigns missing required fields are rejected with field-level validation errors. (Scale: boolean | Pass/Fail)

**Constrained by:** CON-002 — First Principles Design (0xb)
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

## 6c. User Registration (BUC-03, scope:0x2a)

### FR-011: User Account Registration (0x24)

**Priority:** Must | **Originator:** Scitizen, Researcher

**Description:** The platform shall allow open self-registration. Authentication is delegated to Zitadel (email/password or social login). On first authenticated access, the platform creates a local `app_users` record with a ULID primary key, the Zitadel user ID stored as `idp_id`, and a `user_type` indicating one or more roles: **scitizen**, **researcher**, or both. The user selects their type during registration. No institutional affiliation is required to register; researchers are associated with organizations separately via FR-004.

**Rationale:** Registration is the prerequisite for all platform participation. First contribution is the retention inflection point (FACT-006) — registration must be minimal friction. A local `app_users` table decouples platform identity from the IdP: the ULID PK is the canonical user identifier across all platform tables (`owner_id`, `created_by`, `scitizen_id`), while `idp_id` is a reference column linking to Zitadel. This follows CON-004 (ULID identifiers) and supports SEC-004 (identity separation from observation data). Users may be both researchers and scitizens — the roles are not mutually exclusive.

**Fit Criterion:** A user can complete registration (including type selection) and reach the campaign browse page in under 2 minutes from landing page. The `app_users` record is created with a ULID PK, `idp_id` referencing Zitadel, and `user_type` reflecting the user's selection. All platform tables reference users by ULID, not by IdP ID. (Scale: minutes | Worst: 5 | Plan: 2 | Best: 1)

**Constrained by:** CON-004 — ULID Identifiers
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

## 6d. Device Registration and Enrollment (BUC-04, scope:0x2b)

### FR-013: Generate Enrollment Code (0x2b)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall generate a short-lived (15-minute TTL), one-time-use, human-friendly enrollment code (6-8 alphanumeric characters, no ambiguous characters 0/O/1/l) when a scitizen initiates device registration.

**Rationale:** Enrollment code is the bootstrap trust anchor binding web session to physical device. Must be human-friendly, short-lived, and one-time use.

**Fit Criterion:** Enrollment codes expire after 15 minutes. Used codes are rejected on second use. Codes contain only unambiguous alphanumeric characters. (Scale: minutes | Pass/Fail)

**Depends on:** FR-011
**Cross-ref:** scope:0x2b

---

### FR-014: Direct Device Enrollment — Tier 1 (0x2d)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall accept device enrollment via POST /enroll where the device presents an enrollment code and a locally-generated CSR over HTTPS. The enrollment service validates the code, verifies the CSR, coordinates with the CA, and returns the issued certificate.

**Rationale:** Tier 1 devices (smartphones, Raspberry Pi, modern weather stations) can run HTTPS clients and generate keypairs directly.

**Fit Criterion:** A Tier 1 device presenting a valid enrollment code and CSR receives an X.509 certificate and is registered in the device registry with status active within 10 seconds. (Scale: seconds | Worst: 30 | Plan: 10 | Best: 3)

**Depends on:** FR-013
**Cross-ref:** scope:0x2b, facts:0x8

---

### FR-015: Proxy Device Enrollment — Tier 2 (0x27)

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

**Description:** The platform shall enforce topic ACLs via OPA: a device may only publish to rootstock/{its-own-device-id}/* topics. Attempts to publish to another device's topic are denied and logged.

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

## 6k. Requirement Summary

| Req ID | Name | BUC | Priority | Originator |
|--------|------|-----|----------|-----------|
| FR-001 | Create Organization Tenant | BUC-01 | Must | Research Institution |
| FR-002 | Configure Organization Hierarchy | BUC-01 | Must | Research Institution |
| FR-003 | Define and Assign Roles | BUC-01 | Must | Research Institution |
| FR-004 | Invite and Onboard Researchers | BUC-01 | Must | Research Institution |
| FR-005 | Create Campaign | BUC-02 | Must | Researcher |
| FR-006 | Define Campaign Parameters | BUC-02 | Must | Researcher |
| FR-007 | Define Campaign Region | BUC-02 | Must | Researcher |
| FR-008 | Define Campaign Time Window | BUC-02 | Must | Researcher |
| FR-009 | Publish and Discover Campaigns | BUC-02 | Must | Researcher |
| FR-010 | Monitor Campaign Data Quality | BUC-02 | Should | Researcher |
| FR-011 | User Account Registration | BUC-03 | Must | Scitizen, Researcher |
| FR-012 | Browse Campaigns | BUC-03 | Must | Scitizen |
| FR-013 | Generate Enrollment Code | BUC-04 | Must | Scitizen |
| FR-014 | Direct Device Enrollment (Tier 1) | BUC-04 | Must | Scitizen |
| FR-015 | Proxy Device Enrollment (Tier 2) | BUC-04 | Must | Scitizen |
| FR-016 | Device Registry Entry | BUC-04 | Must | Scitizen |
| FR-017 | Enroll Device in Campaign | BUC-05 | Must | Scitizen |
| FR-018 | Multi-Campaign Device Enrollment | BUC-05 | Must | Scitizen |
| FR-019 | Device Eligibility Check | BUC-05 | Must | Researcher |
| FR-020 | Authenticate Device via mTLS | BUC-06 | Must | Researcher |
| FR-021 | Authorize Device Actions via OPA | BUC-06 | Must | Researcher |
| FR-022 | Validate Sensor Readings | BUC-06 | Must | Researcher |
| FR-023 | Persist Valid Readings with Provenance | BUC-06 | Must | Researcher |
| FR-024 | Topic ACL Enforcement | BUC-06 | Must | Researcher |
| FR-025 | Anomaly Flagging | BUC-06 | Should | Researcher |
| FR-026 | Export Campaign Data | BUC-07 | Must | Researcher |
| FR-027 | Separate Contributor Identity from Data | BUC-07 | Must | Oversight Body |
| FR-028 | Automated Certificate Renewal | BUC-08 | Must | Scitizen |
| FR-029 | Grace Period for Expired Certificates | BUC-08 | Must | Scitizen |
| FR-030 | Certificate Revocation via Registry | BUC-08 | Must | Researcher |
| FR-031 | Bulk Device Suspension by Class | BUC-09 | Must | Researcher |
| FR-032 | Flag Data from Vulnerable Window | BUC-09 | Should | Researcher |
| FR-033 | Device Reinstatement After Patch | BUC-09 | Should | Scitizen |
| FR-034 | Compute Contribution Score | BUC-10 | Must | Scitizen |
| FR-035 | Award Badges and Recognition | BUC-10 | Should | Scitizen |
| FR-036 | Manage Sweepstakes Entries | BUC-10 | Should | Scitizen |

**MoSCoW Summary:** 29 Must, 7 Should, 0 Could, 0 Won't

---

*Next: [Section 7 — Look and Feel Requirements](./07_look_and_feel.md)*
