# ROOTSTOCK by Corewood

## Requirements Specification — Section 16: Open Issues

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Open issues are unresolved decisions that affect or block one or more requirements. Each issue has a description, impact assessment, and proposed resolution path. Issues are tracked in the requirements knowledge graph as `OpenIssue` nodes.

> **Knowledge graph reference**: Requirements graph at `grapher/schema/rootstock_requirements.graphql`. Open issues are `OpenIssue` type with `blocks` edges to affected requirements.

---

### OI-001: Consumer Sensor Accuracy Validation (0x92)

**Description:** ASM-007 assumes consumer sensors produce publication-grade data. Weather and air quality are partially validated (PurpleAir PM2.5 data used in peer-reviewed wildfire research; Weather Underground data cross-validated against reference stations), but water quality and soil sensors have no systematic cross-validation against reference instruments. A co-location study is needed before campaigns in these categories.

**Impact:** If consumer sensor accuracy is insufficient, campaigns in water/soil categories produce unusable data. Researchers stop creating these campaigns.

**Resolution:** Commission co-location studies for water quality (Atlas Scientific) and soil (DFRobot) sensors during pilot. Define minimum accuracy thresholds per parameter category.

**Blocks:** FR-022 (Validate Sensor Readings), FR-006 (Define Campaign Parameters)

---

### OI-002: Sweepstakes Legal Compliance (0x93)

**Description:** Sweepstakes and lottery-based incentives are subject to jurisdiction-specific regulations. "No purchase necessary" requirements, age restrictions, geographic exclusions, and official rules vary by country and US state.

**Impact:** Non-compliance could expose Corewood or participating institutions to legal liability. Sweepstakes may need to be restricted by geography.

**Resolution:** Legal review of sweepstakes model across target jurisdictions before launch. Consider partnering with a sweepstakes administration service.

**Blocks:** FR-036 (Manage Sweepstakes Entries)

---

### OI-003: Tier 3 Device Architecture (0x94)

**Description:** Tier 3 devices (cheap 8-bit MCU sensors, BLE-only devices) require a gateway architecture. This is explicitly out of MVP scope but the architecture should not preclude it.

**Impact:** If the enrollment and ingestion architecture is not gateway-aware, Tier 3 support becomes a redesign rather than an extension.

**Resolution:** Ensure enrollment and ingestion APIs are designed with extensibility for gateway-mediated device enrollment. Separate project for Tier 3 gateway.

**Blocks:** FR-014 (Direct Device Enrollment), FR-015 (Proxy Device Enrollment)

---

### OI-004: Clock Drift Tolerance Strategy (0x95)

**Description:** Cheap IoT device clocks drift. TLS handshake requires accurate time for cert validation. Data timestamps may be wrong. Current mitigation (NTP on boot, 48-hour notBefore tolerance) is assumed but not validated at scale.

**Impact:** Excessive clock drift causes: TLS failures (cert appears not yet valid), data provenance issues (wrong timestamps), campaign window enforcement errors.

**Resolution:** Define acceptable drift tolerance per campaign type. Implement server-side timestamp reconciliation. Monitor per-device drift rates.

**Blocks:** FR-022 (Validate Sensor Readings), FR-020 (Authenticate Device via mTLS)

---

### OI-005: Open Source Governance Model (0x96)

**Description:** The project is open source (CON-001) but no governance model has been defined. Questions: contribution guidelines, release process, decision-making authority, security disclosure process.

**Impact:** Without governance, external contributions may be blocked, security issues may not be disclosed responsibly, and the project lacks legitimacy for institutional adoption.

**Resolution:** Define governance model (BDFL, steering committee, or foundation) before public launch. Publish CONTRIBUTING.md, SECURITY.md, and release process documentation.

**Blocks:** None directly — governance is an organizational issue, not a functional blocker.

---

### OI-006: Spatial Resolution Privacy Tradeoff (0x97)

**Description:** Researchers want high spatial resolution for scientific accuracy. Privacy requires spatial degradation. The optimal balance is campaign-dependent and may require ethics review for each campaign.

**Impact:** Too much degradation makes data scientifically useless. Too little degradation enables re-identification.

**Resolution:** Define spatial resolution tiers (e.g., 100m, 1km, 10km). Require researcher justification for resolutions finer than 1km. Subject fine-resolution requests to automated privacy risk assessment.

**Blocks:** FR-027 (Separate Contributor Identity from Data), SEC-004 (Data Privacy and Identity Separation)

---

*Next: [Section 17 — Off-the-Shelf Solutions](./17_off_the_shelf.md)*
