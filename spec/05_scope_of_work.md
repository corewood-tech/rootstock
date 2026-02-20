# ROOTSTOCK by Corewood

## Requirements Specification — Section 5: Scope of the Work

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Section 5 bridges "why are we doing this?" (Sections 1–4) to "what should the product do?" (Section 6+). It covers: the current situation, the work context (adjacent systems and data flows), the business event list, and business use cases.

> **Knowledge graph reference**: The scope model behind this section is captured in a persistent Dgraph knowledge graph (`grapher/schema/rootstock_scope.graphql`). Nodes are referenced by UID for traceability. Start the scope graph with `GRAPH=scope podman compose -f grapher/compose-grapher.yml up -d` and query at `http://localhost:18080`.

---

## 5a. The Current Situation

There is no current system (CON-003). The work that Rootstock aims to support — researchers requesting field data and citizen scientists providing it — currently happens across disconnected processes with no shared platform, protocol, or data format. Four fragmented processes define the status quo:

### Manual Field Data Collection (scope:0x1)

Researchers hire field teams or travel personally to collect sensor readings at specific locations. Personnel costs consume 50–80% of grant budgets (the average annualized NSF research grant is $245,800 in FY2024). Coverage is limited to areas researchers can physically reach, temporal coverage is limited to field visit windows, and collection cannot scale without proportional budget increases. In the TRY global plant trait database, 70.2% of removed records were duplicates — indicating massive redundancy in manually collected data. (facts:0x1, facts:0x6, facts:0x7)

### Siloed Institutional Data Management (scope:0x2)

Each research institution manages collected data independently using internal tools, formats, and storage. No shared platform, protocol, or data format exists across institutions. Researchers at different institutions cannot combine or compare datasets without manual reconciliation. Data silos prevent large-scale analysis and cause duplicated collection effort across institutions studying similar phenomena in overlapping regions. (facts:0x3)

### Observation-Based Citizen Science (scope:0x3)

Existing citizen science platforms — iNaturalist (300M+ observations), eBird (2.1B observations), Zooniverse (2.7M volunteers) — rely on human observers actively photographing, identifying, or classifying data. None provides a general-purpose platform for connecting researchers' structured data needs to citizen-owned IoT devices across multiple sensor types. 75% of citizen science projects produce zero peer-reviewed publications. Sustained multi-year engagement drops to 5.25%. Contributions follow a power law: ~10% of participants produce ~90% of data. (facts:0x2, facts:0x4, facts:0x32, facts:0x35)

### Uncoordinated Consumer Sensor Networks (scope:0x4)

250,000+ personal weather stations (Weather Underground) and 30,000+ air quality sensors (PurpleAir) are already deployed at scale. These devices generate continuous data but feed into vendor-specific aggregation platforms with no research campaign interface. PurpleAir data has been used in peer-reviewed wildfire research, proving the data is publication-grade — but access is ad-hoc, not systematic. No mechanism exists for researchers to request specific data from specific locations, and no cross-parameter correlation is possible. Water quality and soil monitoring have no consumer network equivalent at all. (facts:0x8, facts:0x2e, facts:0x27, facts:0x28)

---

## 5b. The Work Context

### Adjacent Systems

The following people, organizations, and automated systems interact with Rootstock. Each is an external entity that sends data to or receives data from the product.

| UID | System | Type | Description |
|-----|--------|------|-------------|
| 0x5 | Researcher | Person | Defines data needs, creates campaigns, monitors quality, exports results. (facts:0xc) |
| 0x6 | Scitizen | Person | Owns consumer IoT devices, enrolls them in campaigns, contributes sensor data. (facts:0xd) |
| 0x7 | Research Institution | Organization | Employs researchers. Creates organizational tenants. Has compliance and governance requirements. (facts:0xe) |
| 0x8 | Oversight Body | Regulatory | IRBs, ethics committees, grant boards, funding agencies. Imposes data governance and consent requirements. (facts:0xf) |
| 0x9 | IoT Device | Automated | Consumer sensor devices (weather stations, air quality monitors, water/soil probes). Tier 1 (direct) and Tier 2 (proxy). All data crosses a trust boundary. (facts:0x8, facts:0x5) |
| 0xa | Rootstock CA | Automated | Self-operated certificate authority. Two-tier hierarchy. Issues and renews device certificates for mTLS. |
| 0xb | Identity Provider (Zitadel) | Automated | External identity provider for human users. Handles authentication, user management, org hierarchy. |
| 0xc | Policy Engine (OPA) | Automated | Makes all authorization decisions. Evaluates Rego policies against device registry data. 30-second bundle refresh. |
| 0xd | MQTT Broker | Automated | Routes IoT device telemetry. Embedded in-process (Mochi MQTT). All connections use strict mTLS (`RequireAndVerifyClientCert`). Auth handled by in-process Go hook — mTLS cert verification + topic ACL enforcement. Port 8883. |

### Data Flows

| UID | Flow | Direction | From/To | Data Elements |
|-----|------|-----------|---------|---------------|
| 0xe | Campaign Definition | Inbound | Researcher → Rootstock | Campaign parameters, geographic region, time window, quality thresholds, device eligibility |
| 0xf | Collected Data Export | Outbound | Rootstock → Researcher | Validated readings, data provenance, quality metrics |
| 0x10 | Device Enrollment | Inbound | Scitizen + Device → Rootstock | Enrollment code, CSR, device profile |
| 0x11 | Device Certificate | Outbound | Rootstock CA → Device | X.509 certificate |
| 0x12 | Sensor Telemetry | Inbound | Device → Rootstock | Sensor readings, timestamp, geolocation, device ID, campaign ID |
| 0x13 | Device Configuration | Outbound | Rootstock → Device | Campaign config, sampling parameters |
| 0x14 | Authorization Decision | Bidirectional | Rootstock ↔ OPA | Policy query, allow/deny decision, device registry bundle |
| 0x15 | Identity Token | Inbound | Zitadel → Rootstock | Identity token, user ID, org membership, roles |
| 0x16 | Compliance Requirements | Inbound | Oversight Body → Rootstock | Ethics requirements, consent model, data governance rules |
| 0x17 | Scitizen Recognition | Outbound | Rootstock → Scitizen | Contribution score, badges, sweepstakes entries, campaign acknowledgments |
| 0x18 | Institution Onboarding | Inbound | Research Institution → Rootstock | Org structure, member invitations, role assignments |

---

## 5c. Business Event List

| Number | Event | Type | Originates From | Input/Output |
|--------|-------|------|-----------------|--------------|
| BE-01 | Institution Requests Platform Access (0x1c) | External, non-temporal | Research Institution (0x7) | Input: Institution Onboarding (0x18) |
| BE-02 | Researcher Creates Campaign (0x1d) | External, non-temporal | Researcher (0x5) | Input: Campaign Definition (0xe) |
| BE-03 | Scitizen Registers Account (0x1e) | External, non-temporal | Scitizen (0x6) | Input: Identity Token (0x15) |
| BE-04 | Scitizen Enrolls Device (0x1f) | External, non-temporal | Scitizen (0x6) | Input: Device Enrollment (0x10); Output: Device Certificate (0x11) |
| BE-05 | Device Submits Sensor Data (0x20) | External, non-temporal | IoT Device (0x9) | Input: Sensor Telemetry (0x12) |
| BE-06 | Researcher Exports Campaign Data (0x21) | External, non-temporal | Researcher (0x5) | Output: Collected Data Export (0xf) |
| BE-07 | Device Certificate Approaches Expiration (0x22) | External, temporal | IoT Device (0x9) | — |
| BE-08 | Campaign Window Opens (0x23) | External, temporal | — | Output: Device Configuration (0x13) |
| BE-09 | Campaign Window Closes (0x24) | External, temporal | — | — |
| BE-10 | Device Firmware Vulnerability Discovered (0x25) | External, non-temporal | IoT Device (0x9) | — |
| BE-11 | Scitizen Enrolls Device in Campaign (0x26) | External, non-temporal | Scitizen (0x6) | — |
| BE-12 | Oversight Body Updates Compliance Requirements (0x27) | External, non-temporal | Oversight Body (0x8) | Input: Compliance Requirements (0x16) |

---

## 5d. Business Use Cases

### BUC-01: Institutional Onboarding (scope:0x28)

**Summary:** A research institution establishes its presence on the platform by creating an organizational tenant, configuring hierarchy (departments, labs), defining roles, and inviting researchers. Zero changes to the institution's existing internal systems.

**Triggered by:** BE-01 — Institution Requests Platform Access (0x1c)

**Preconditions:** Institution has decided to use Rootstock. At least one admin user exists or can be created via Zitadel.

**Postconditions:** Organization tenant exists. Hierarchy configured. Researchers invited and able to authenticate. Roles and permissions enforced by OPA.

**Business rules:** Onboarding requires no integration with institutional systems (CON-003). Organization hierarchy supports nesting. Roles are organization-scoped, not global.

**Constrained by:** CON-003 — No Shared Context Exists (0x1b)
**Supported by:** CON-001 — Open Source (0x19)

---

### BUC-02: Campaign Creation and Management (scope:0x29)

**Summary:** A researcher creates a structured data campaign specifying what data is needed (parameters), where (region), when (time window), to what quality standard (thresholds), and from which devices (eligibility). The campaign is published and discoverable by scitizens.

**Triggered by:** BE-02 — Researcher Creates Campaign (0x1d)

**Preconditions:** Researcher is authenticated via Zitadel. Researcher belongs to an onboarded institution with campaign creation permission.

**Postconditions:** Campaign is published with defined parameters, region, window, and quality thresholds. Campaign is discoverable. Scitizens can enroll devices.

**Business rules:** Campaign parameters must be explicitly defined — no open-ended collection. Every technology choice traceable to a requirement (CON-002). Campaign structure mirrors marketing campaign mechanics: goal, audience, timeframe, measurement. (facts:0x9, facts:0x14, facts:0x12)

**Constrained by:** CON-002 — First Principles Design (0x1a)
**Depends on:** BUC-01 — Institutional Onboarding (0x28)

---

### BUC-03: Scitizen Registration (scope:0x2a)

**Summary:** A citizen scientist creates an account on the platform. Identity managed by Zitadel. Registration is the prerequisite for device enrollment and campaign participation.

**Triggered by:** BE-03 — Scitizen Registers Account (0x1e)

**Preconditions:** None. Open registration.

**Postconditions:** Scitizen has an authenticated identity. Can browse campaigns and initiate device enrollment.

**Business rules:** Minimize friction at registration — the first contribution is the critical retention inflection point (FACT-006, facts:0x33, facts:0x13). Identity delegated to Zitadel.

---

### BUC-04: Device Registration and Enrollment (scope:0x2b)

**Summary:** A scitizen registers a consumer IoT device on the platform. Enrollment code binds web session to physical device. Device generates keypair and CSR. Enrollment service validates, CA issues certificate. Device entered in registry. Supports Tier 1 (direct) and Tier 2 (proxy via companion app) enrollment paths.

**Triggered by:** BE-04 — Scitizen Enrolls Device (0x1f)

**Preconditions:** Scitizen is authenticated. Device is Tier 1 or Tier 2. Device can generate a keypair and CSR.

**Postconditions:** Device has a valid X.509 certificate. Device is registered in the device registry with status pending → active. Device can authenticate via mTLS.

**Business rules:** Enrollment code: 6–8 alphanumeric, no ambiguous chars (0/O, 1/l), 15-minute TTL, one-time use. Private key never leaves device. Certificate: CN=device-id, 90-day lifetime, no metadata in cert — all mutable state in registry.

**Depends on:** BUC-03 — Scitizen Registration (0x2a)

---

### BUC-05: Campaign Enrollment (scope:0x2c)

**Summary:** A scitizen discovers a campaign relevant to their location and device capabilities and enrolls a registered device. Campaign enrollment is tracked in the device registry and enforced by OPA on every publish. Cross-org participation is the normal case — scitizens typically have no institutional affiliation with the requesting researcher.

**Triggered by:** BE-11 — Scitizen Enrolls Device in Campaign (0x26)

**Preconditions:** Scitizen has a registered, active device. Campaign is published and its window has not closed. Device meets campaign eligibility criteria.

**Postconditions:** Device is enrolled in campaign. Registry updated. OPA allows publish to campaign topic on next bundle refresh. Device receives campaign configuration.

**Business rules:** A device may be enrolled in multiple campaigns simultaneously. Cross-org participation does not require institutional affiliation. Eligibility checked against device class, tier, and sensor capabilities.

**Depends on:** BUC-04 — Device Registration and Enrollment (0x2b), BUC-02 — Campaign Creation and Management (0x29)

---

### BUC-06: Data Ingestion and Validation (scope:0x2d)

**Summary:** An enrolled device transmits sensor telemetry. The platform authenticates (mTLS), authorizes (OPA checks device status, campaign enrollment, topic ACL), validates (schema, range, rate, anomaly), routes to campaign, and persists. Rejected readings return explicit reasons.

**Triggered by:** BE-05 — Device Submits Sensor Data (0x20)

**Preconditions:** Device has valid certificate. Device is enrolled in at least one active campaign. Campaign window is open.

**Postconditions:** Valid readings persisted with full provenance. Invalid readings rejected with reason. Anomalies flagged for review. Contribution score updated.

**Business rules:** All device data crosses a trust boundary and is treated as untrusted (FACT-010, facts:0x5, facts:0x2f). Validation includes schema, range, rate limiting, geolocation, and timestamp checks. Quarantined readings are not deleted. Topic ACL: device can only publish under its own device-id.

**Depends on:** BUC-05 — Campaign Enrollment (0x2c)

---

### BUC-07: Data Export and Analysis (scope:0x2e)

**Summary:** A researcher exports validated campaign data for analysis in external tools. Export includes readings with provenance metadata, quality metrics, and campaign context. Contributor identity is separated from observation data.

**Triggered by:** BE-06 — Researcher Exports Campaign Data (0x21)

**Preconditions:** Campaign has collected validated readings. Researcher has export permission for the campaign.

**Postconditions:** Data exported in requested format with full provenance. Contributor identity not present in exported dataset. Data meets FAIR principles for findability, accessibility, interoperability, and reusability.

**Business rules:** Raw contributor locations never appear in public-facing datasets — 4 spatiotemporal data points uniquely identify 95% of individuals (FACT-009, facts:0x30). Export formats and access controls are campaign-scoped. Spatial resolution configurable per campaign.

**Depends on:** BUC-06 — Data Ingestion and Validation (0x2d)

---

### BUC-08: Certificate Lifecycle Management (scope:0x2f)

**Summary:** The platform manages the lifecycle of device certificates: automated renewal at day 60 of 90-day lifetime, 7-day grace period for expired certs (renewal endpoint only), re-enrollment required beyond grace period. Revocation handled by registry status + OPA denial, not CRL/OCSP.

**Triggered by:** BE-07 — Device Certificate Approaches Expiration (0x22)

**Preconditions:** Device has been enrolled and has an issued certificate.

**Postconditions:** Certificate renewed or device re-enrolled. Registry updated. Expired-beyond-grace devices cannot communicate until re-enrolled.

**Business rules:** Renewal: device presents current cert via mTLS + new CSR. Grace period: 7 days post-expiry, renewal endpoint only. Revocation: registry status update, OPA denies within 30 seconds (bundle refresh). The platform controls the relying party, so traditional CRL/OCSP infrastructure is unnecessary.

**Depends on:** BUC-04 — Device Registration and Enrollment (0x2b)

---

### BUC-09: Device Security Response (scope:0x30)

**Summary:** The platform responds to a discovered firmware vulnerability by bulk suspending or revoking all devices of the affected device class. Registry batch update, OPA denial on next bundle refresh, mass notification to affected scitizens.

**Triggered by:** BE-10 — Device Firmware Vulnerability Discovered (0x25)

**Preconditions:** Vulnerability identified in a specific device class and firmware version. Affected devices are in the registry.

**Postconditions:** All affected devices suspended or revoked. OPA denying their requests. Scitizens notified with remediation instructions. Data from affected devices during vulnerable window flagged for review.

**Business rules:** 50%+ of IoT devices have critical exploitable vulnerabilities (FACT-010, facts:0x2f). Platform must assume every device is potentially compromised. Bulk operations by device class. Suspension is reversible; revocation is permanent.

---

### BUC-10: Scitizen Recognition and Incentives (scope:0x31)

**Summary:** The platform tracks scitizen contributions, computes contribution scores, awards badges and recognition, and manages sweepstakes entries. Recognition and visible research impact are the primary engagement layer; sweepstakes supplement but do not replace intrinsic motivation.

**Preconditions:** Scitizen has contributed data via enrolled devices.

**Postconditions:** Contribution score updated. Badges awarded where earned. Sweepstakes entries granted. Recognition visible on profile and dashboard.

**Business rules:** Sweepstakes are behaviorally optimal per Prospect Theory — people overweight small probabilities (FACT-008, facts:0x2c). Recognition must not crowd out intrinsic motivation (Self-Determination Theory). The first contribution is the retention inflection point — 82% retention post-first-contribution vs. 39.7% without (FACT-006, facts:0x33). Variable-ratio reinforcement produces the strongest sustained engagement patterns (facts:0x11).

**Depends on:** BUC-06 — Data Ingestion and Validation (0x2d)

---

## 5e. BUC Dependency Order

```
BUC-01: Institutional Onboarding
  └─► BUC-02: Campaign Creation and Management
        └─► BUC-05: Campaign Enrollment ◄── BUC-04: Device Registration and Enrollment
              └─► BUC-06: Data Ingestion and Validation
                    ├─► BUC-07: Data Export and Analysis
                    └─► BUC-10: Scitizen Recognition and Incentives

BUC-03: Scitizen Registration
  └─► BUC-04: Device Registration and Enrollment
        ├─► BUC-05: Campaign Enrollment (above)
        └─► BUC-08: Certificate Lifecycle Management

BUC-09: Device Security Response (independent — triggered by external discovery)
```

---

*Next: [Section 6 — Functional Requirements](./06-functional-requirements.md)*
