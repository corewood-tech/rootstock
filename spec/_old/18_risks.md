# ROOTSTOCK by Corewood

## Requirements Specification — Section 18: Risks

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Risks are potential problems that could threaten the success of the project. Each risk is traced to the assumption, fact, or requirement it derives from. Risks without mitigation strategies are escalated to open issues (Section 16).

---

## 18a. Supply-Side Risks

### RISK-001: Insufficient Scitizen Adoption

**Source:** ASM-001 — Consumer IoT Device Owners Will Contribute Data to Science If the Platform Exists

**Description:** The platform assumes a meaningful percentage of consumer sensor owners will enroll their devices. If enrollment rates are below 5% of contacted device owners, the supply side collapses and campaigns have no contributors.

**Probability:** Medium | **Impact:** Critical

**Mitigation:** Early access pilot with existing PurpleAir, Ecowitt, or Davis Instruments owners. Measure enrollment rate, first-contribution rate, and 30-day retention. If <5% enroll, revisit the value proposition before public launch.

**Affected requirements:** FR-011, FR-013, FR-014, FR-017

---

### RISK-002: Contributor Retention Failure

**Source:** ASM-002, ASM-005, FACT-006 (facts:0x4, facts:0x33)

**Description:** Even with automated IoT device contribution, scitizens may disengage — unplugging devices, changing WiFi passwords, or forgetting the platform. Sustained multi-year engagement in citizen science drops to 5.25% (CSMON-LIFE). If >20% of devices stop reporting within 60 days for non-technical reasons, the passive contribution assumption is challenged.

**Probability:** Medium | **Impact:** High

**Mitigation:** Track device uptime and data continuity across enrolled devices. Implement re-engagement campaigns for dormant contributors. The first contribution inflection point (82% retention post-first-contribution) means the enrollment-to-first-data path is the highest-leverage optimization target (US-001).

**Affected requirements:** FR-034, FR-035, FR-036, US-001

---

## 18b. Demand-Side Risks

### RISK-003: Researcher Non-Adoption

**Source:** ASM-003 — Researchers Will Create Campaigns on the Platform

**Description:** If campaign creation takes >30 minutes or the resulting data is not publication-grade, researchers will not invest time in the platform. Without campaigns, scitizens have nothing to contribute to.

**Probability:** Medium | **Impact:** Critical

**Mitigation:** Partner with 3–5 research groups during development. Have them define real campaigns. Validate that campaign creation is under 15 minutes (US-002) and that exported data meets publication standards (CO-003).

**Affected requirements:** FR-005, FR-006, FR-026, US-002, CO-003

---

### RISK-004: Institutional Compliance Blocks Adoption

**Source:** ASM-004 — Institutional Data Governance Will Not Block Adoption

**Description:** If each institution requires bespoke compliance work (per-institution IRB customization, data governance integration), onboarding becomes non-scalable. CON-003 mandates zero changes to institutional systems.

**Probability:** Medium | **Impact:** High

**Mitigation:** Engage 2–3 university IRBs during architecture phase. Present the consent model (CO-002), location privacy architecture (SEC-004), and data provenance chain (CO-003). Get pre-clearance or identify blockers before building.

**Affected requirements:** FR-001, FR-004, CO-001, CO-002, SEC-004

---

## 18c. Technical Risks

### RISK-005: Consumer Sensor Data Quality Insufficient

**Source:** ASM-007, OI-001 (0x92)

**Description:** Consumer sensor accuracy for weather and air quality is partially validated, but water quality and soil sensors have no systematic cross-validation against reference instruments. If consumer sensor error margins exceed researcher-defined thresholds for >50% of campaign types, the platform produces unusable data.

**Probability:** High (for water/soil) | **Impact:** Critical (for those categories)

**Mitigation:** Commission co-location studies for water quality (Atlas Scientific) and soil (DFRobot) sensors during pilot. Define minimum accuracy thresholds per parameter category. Launch with well-validated categories (weather, air quality) first.

**Affected requirements:** FR-006, FR-022, PE-001

---

### RISK-006: Clock Drift Causes Systematic Data Issues

**Source:** OI-004 (0x95)

**Description:** Cheap IoT device clocks drift, causing TLS handshake failures (cert appears not yet valid), wrong timestamps on data readings, and campaign window enforcement errors. The 48-hour notBefore tolerance and NTP-on-boot mitigation are assumed but not validated at scale.

**Probability:** Medium | **Impact:** High

**Mitigation:** Define acceptable drift tolerance per campaign type. Implement server-side timestamp reconciliation. Monitor per-device drift rates and alert on systematic drift.

**Affected requirements:** FR-020, FR-022, FR-028, SEC-001

---

### RISK-007: IoT Device Compromise at Scale

**Source:** FACT-010 (facts:0x2f, facts:0x5)

**Description:** Over 50% of IoT devices have critical exploitable vulnerabilities. 60% of breaches originate from unpatched firmware. A coordinated attack using compromised consumer devices could inject fabricated data at scale, undermining data integrity across multiple campaigns.

**Probability:** Medium | **Impact:** Critical

**Mitigation:** The architecture already assumes every device is potentially compromised (trust boundary at ingestion). FR-022 (validation), FR-024 (topic ACL), FR-031 (bulk suspension), and FR-032 (vulnerability window flagging) collectively address this. Cross-device correlation and anomaly detection (FR-025) provide an additional defense layer.

**Affected requirements:** FR-020, FR-021, FR-022, FR-024, FR-025, FR-031, SEC-001

---

## 18d. Legal and Regulatory Risks

### RISK-008: Location Privacy Re-Identification

**Source:** FACT-009 (facts:0x30), OI-006 (0x97)

**Description:** 4 spatiotemporal data points uniquely identify 95% of individuals. Simple GPS truncation is insufficient. If the platform's spatial degradation is too weak, exported datasets enable contributor re-identification, violating GDPR and exposing institutions to liability.

**Probability:** Medium | **Impact:** Critical

**Mitigation:** Architectural separation of contributor identity from observation data (SEC-004). Configurable spatial resolution per campaign with researcher justification for fine resolution. Automated privacy risk assessment for sub-1km resolution requests. No raw contributor locations in public-facing datasets (FR-027).

**Affected requirements:** FR-027, SEC-004, CO-001

---

### RISK-009: Sweepstakes Regulatory Exposure

**Source:** OI-002 (0x93)

**Description:** Sweepstakes and lottery-based incentives are subject to jurisdiction-specific regulations that vary by country and US state. Non-compliance could expose Corewood or participating institutions to legal liability.

**Probability:** Medium | **Impact:** Medium

**Mitigation:** Legal review of sweepstakes model across target jurisdictions before launch. Consider partnering with a sweepstakes administration service. Design the incentive system to be modular so sweepstakes can be restricted by geography without removing other recognition mechanisms (FR-034, FR-035).

**Affected requirements:** FR-036

---

## Risk Summary

| Risk | Probability | Impact | Primary Mitigation |
|------|------------|--------|-------------------|
| RISK-001 Scitizen Adoption | Medium | Critical | Pilot with existing device owners |
| RISK-002 Retention Failure | Medium | High | First-contribution optimization |
| RISK-003 Researcher Non-Adoption | Medium | Critical | Partner with research groups |
| RISK-004 Compliance Blocks | Medium | High | Pre-clearance with IRBs |
| RISK-005 Sensor Data Quality | High | Critical | Co-location validation studies |
| RISK-006 Clock Drift | Medium | High | Server-side reconciliation |
| RISK-007 Device Compromise | Medium | Critical | Trust boundary + validation |
| RISK-008 Re-Identification | Medium | Critical | Architectural identity separation |
| RISK-009 Sweepstakes Regulation | Medium | Medium | Legal review + modular design |

---

*Next: [Section 19 — Costs](./19_costs.md)*
