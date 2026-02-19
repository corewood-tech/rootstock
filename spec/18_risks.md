# ROOTSTOCK by Corewood

## Requirements Specification — Section 18: Risks

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Risks are conditions or events that, if they occur, would negatively affect the project's ability to meet its goals. Each risk is sourced from external evidence captured in the analysis knowledge graph (`grapher/schema/rootstock_analysis.graphql`). Every claim traces to a graph node with an external citation.

> **Knowledge graph reference**: Analysis graph at `grapher/schema/rootstock_analysis.graphql`. Start with `GRAPH=analysis podman compose -f grapher/compose-grapher.yml up -d` and query at `http://localhost:18080`.

---

## 18a. Supply-Side Risks

These risks affect the platform's ability to attract and retain the citizen scientists who produce data.

### RISK-001: Insufficient Scitizen Adoption (0x19)

**Probability:** Medium | **Impact:** Critical | **Category:** Supply-side

**Description:** The platform assumes consumer sensor owners will enroll devices and contribute data. However, the 90-9-1 participation inequality rule — documented across citizen science projects — predicts that 90% of registrants will lurk, 9% will contribute occasionally, and only 1% will become bulk contributors (0x2, Recruiting and Retaining Participants in Citizen Science, *Citizen Science: Theory and Practice*, 2017). If enrollment rates fall below 5% of contacted device owners, the supply side collapses.

The retention data offers a counterpoint: participants who make a first contribution retain at 82%, compared to 39.7% for those who do not (0x4, Participant retention in a continental-scale citizen science project, *BioScience*, 2023). This makes first-contribution optimization the highest-leverage intervention.

**Mitigation:** Early access pilot with existing PurpleAir, Ecowitt, and Davis device owners. Measure enrollment rate, first-contribution rate, and 30-day retention. If fewer than 5% enroll, revisit value proposition before public launch.

**Affected Requirements:** FR-011 (Scitizen Account Registration), FR-013 (Generate Enrollment Code), FR-014 (Direct Device Enrollment), FR-017 (Enroll Device in Campaign)

**Evidence:**
- 0x2 — 90-9-1 participation inequality rule ([source](https://theoryandpractice.citizenscienceassociation.org/articles/10.5334/cstp.8), 2017)
- 0x4 — 82% retention after first contribution ([source](https://doi.org/10.1093/biosci/biad041), 2023)

---

### RISK-002: Contributor Retention Failure (0x1a)

**Probability:** Medium | **Impact:** High | **Category:** Supply-side

**Description:** Research on the CSMON-LIFE project found that 72% of citizen science volunteers are active for only one month, with long-term retention dropping to 5.25% (0x1, Volunteers Recruitment, Retention, and Performance during CSMON-LIFE, *Sustainability*, 2021). The 90-9-1 rule further compounds this: even among those who register, the vast majority will not contribute sustained data (0x2, 2017).

The critical inflection point is the first contribution: 82% of participants who make a first contribution are retained, versus 39.7% who do not (0x4, *BioScience*, 2023). For an IoT-based platform, "disengagement" may manifest as unplugging devices, changing WiFi passwords, or simply forgetting — even when no active effort is required.

**Mitigation:** Optimize the enrollment-to-first-data path (US-001, under 10 minutes). Track device uptime and data continuity as retention proxies. Implement re-engagement campaigns for dormant contributors. The 82% post-first-contribution retention rate makes first-data the highest-leverage optimization target.

**Affected Requirements:** FR-034 (Compute Contribution Score), FR-035 (Award Badges and Recognition), FR-036 (Manage Sweepstakes Entries), US-001

**Evidence:**
- 0x1 — 72% active only 1 month; 5.25% long-term retention ([source](https://www.mdpi.com/2071-1050/13/19/11110), 2021)
- 0x2 — 90-9-1 participation inequality rule ([source](https://theoryandpractice.citizenscienceassociation.org/articles/10.5334/cstp.8), 2017)
- 0x4 — 82% retention post-first-contribution ([source](https://doi.org/10.1093/biosci/biad041), 2023)

---

## 18b. Demand-Side Risks

These risks affect the platform's ability to attract researchers who create the campaigns that give scitizens something to contribute to.

### RISK-003: Researcher Non-Adoption (0x1e)

**Probability:** Medium | **Impact:** Critical | **Category:** Demand-side

**Description:** An analysis of citizen science research found that 75.3% of citizen science projects produce zero publications, with an average time to first publication of 9.15 years (0x3, Follett & Strezov, An Analysis of Citizen Science Based Research, *PLOS ONE*, 2015). If campaign creation takes more than 30 minutes or the exported data does not meet publication-grade standards, researchers will not invest time in the platform. Without campaigns, scitizens have nothing to contribute to — the demand side and supply side are coupled.

**Mitigation:** Partner with 3–5 research groups during development. Have them define real campaigns. Validate that campaign creation completes in under 15 minutes (US-002) and that exported data meets publication standards (CO-003).

**Affected Requirements:** FR-005 (Create Campaign), FR-006 (Define Campaign Parameters), FR-026 (Export Campaign Data), US-002, CO-003

**Evidence:**
- 0x3 — 75.3% zero publications; 9.15 years average to first publication ([source](https://www.mdpi.com/2071-1050/15/5/4577), 2015)

---

### RISK-004: Institutional Compliance Blocks Adoption (0x1b)

**Probability:** Medium | **Impact:** High | **Category:** Demand-side

**Description:** A systematic review using the TOE (Technology-Organization-Environment) framework identified 22 distinct barriers to institutional adoption of open data platforms (0xc, ScienceDirect, 2023). The number of university IRBs with existing open data policies is unknown. Researchers perceive IRB and ethics offices as integral but burdensome. If each institution requires bespoke compliance negotiation, onboarding becomes non-scalable.

**Mitigation:** Engage 2–3 university IRBs during the architecture phase. Present the consent model (CO-002), location privacy architecture (SEC-004), and provenance chain (CO-003). Obtain pre-clearance or identify blockers before building.

**Affected Requirements:** FR-001 (Create Organization Tenant), FR-004 (Invite and Onboard Researchers), CO-001 (GDPR Compliance Architecture), CO-002 (IRB-Compatible Consent Model), SEC-004 (Data Privacy and Identity Separation)

**Evidence:**
- 0xc — 22 barriers to institutional adoption; no quantitative adoption data ([source](https://www.sciencedirect.com/science/article/pii/S2543925123000232), 2023)

---

## 18c. Technical Risks

These risks arise from the technical characteristics of the platform's operating environment — consumer IoT devices, sensor accuracy, and clock synchronization.

### RISK-005: Consumer Sensor Data Quality Insufficient (0x1f)

**Probability:** High | **Impact:** Critical | **Category:** Technical

**Description:** Sensor accuracy varies dramatically by category and conditions:

- **Air quality:** PurpleAir PM2.5 sensors show R²=0.43 against reference instruments when uncorrected, improving to R²=0.999 with EPA correction factors. However, during dust events sensors underestimate by 5–6x, and only 5 of 15 sensors met EPA performance targets during wildfire smoke (0x5, Intercomparison of PurpleAir Sensor Performance, *Sensors*, 2022; EPA correction evaluation, *AMT*, 2023).
- **Water quality:** Atlas Scientific sensors achieve R²=0.97 in laboratory conditions, but error ranges vary from -0.33% to 33.77% across manufacturers in multi-sensor comparisons (0x6, Low-Cost Water Quality Sensors systematic review, *Sensors*, 2023).
- **Soil moisture:** Low-cost sensors achieve RMSE of 1.87–3.78% VWC depending on ADC and calibration (0x7, Automated Low-Cost Soil Moisture Sensors, *Sensors*, 2023).
- **Data gap:** No peer-reviewed co-location study exists for Ecowitt or Davis weather stations against reference instruments (0x8, research gap identified during systematic search, 2026).

**Mitigation:** Commission co-location studies for water quality (Atlas Scientific) and soil (DFRobot) sensors during pilot. Define minimum accuracy thresholds per parameter category. Launch with well-validated categories (weather, air quality) first. Require correction factors (like the EPA PurpleAir equation) as part of the validation pipeline.

**Affected Requirements:** FR-006 (Define Campaign Parameters), FR-022 (Validate Sensor Readings), PE-001 (Data Ingestion Throughput)
**Related Open Issue:** OI-001 (Consumer Sensor Accuracy Validation)

**Evidence:**
- 0x5 — PurpleAir PM2.5 R²=0.43 raw, 0.999 corrected ([source](https://pmc.ncbi.nlm.nih.gov/articles/PMC9002513/), 2022)
- 0x6 — Atlas Scientific R²=0.97 lab, error range to 33.77% ([source](https://pmc.ncbi.nlm.nih.gov/articles/PMC10181703/), 2023)
- 0x7 — Soil moisture RMSE 1.87–3.78% VWC ([source](https://pmc.ncbi.nlm.nih.gov/articles/PMC10007478/), 2023)
- 0x8 — No Ecowitt/Davis co-location study (research gap, 2026)

---

### RISK-006: Clock Drift Causes Systematic Data Issues (0x20)

**Probability:** Medium | **Impact:** High | **Category:** Technical

**Description:** The ESP32 internal RTC drifts up to 8 minutes per day without calibration, and even with calibration drifts approximately 8 seconds per day. Devices without battery-backed RTC reset to epoch (January 1, 1970) on power loss, which breaks all TLS connections because the certificate appears "not yet valid." Azure IoT Edge has documented clock drift even in containers using host NTP. NTP itself requires UDP port 123, which some networks block (0xa, Climbers.net ESP32 clock study; AWS IoT Blog; Azure IoT Edge GitHub #7389).

Clock drift impacts three areas: TLS handshake failures (certificate validation requires accurate time), data provenance (wrong timestamps on readings), and campaign window enforcement (readings rejected or accepted incorrectly).

**Mitigation:** Define acceptable drift tolerance per campaign type. Implement server-side timestamp reconciliation. Monitor per-device drift rates and alert on systematic drift. Apply a 48-hour `notBefore` tolerance on certificates to absorb reasonable drift.

**Affected Requirements:** FR-020 (Authenticate Device via mTLS), FR-022 (Validate Sensor Readings), FR-028 (Automated Certificate Renewal), SEC-001 (mTLS for All Device Communication)
**Related Open Issue:** OI-004 (Clock Drift Tolerance Strategy)

**Evidence:**
- 0xa — ESP32 drifts 8 min/day uncalibrated; epoch reset on power loss ([source](https://climbers.net/sbc/esp32-accurate-clock-sleep-ntp/))

---

### RISK-007: IoT Device Compromise at Scale (0x1c)

**Probability:** Medium | **Impact:** Critical | **Category:** Technical

**Description:** The IoT security landscape is hostile:

- 57% of IoT devices are vulnerable to medium- or high-severity attacks, with 98% of IoT traffic unencrypted, across a study of 1.2 million devices (0xd, Unit 42 IoT Threat Report, Palo Alto Networks, 2020).
- 60% of IoT security breaches originate from unpatched firmware, with an 88% increase in hardware vulnerabilities year-over-year (0xe, JumpCloud/Forescout, 2025).
- The Mirai botnet demonstrated mass compromise using just 62 default credentials, infecting hundreds of thousands of devices and generating 1 Tbit/sec DDoS attacks (0xf, Inside the Infamous Mirai IoT Botnet, Cloudflare, 2023).

No documented case of scientific data poisoning via compromised IoT devices exists, but once a device is compromised, its data integrity is unverifiable. The platform must assume every device is potentially compromised.

**Mitigation:** The architecture treats every device connection as crossing a trust boundary. Defenses are layered: mTLS authentication (FR-020), OPA authorization on every action (FR-021), reading validation (FR-022), topic ACL enforcement (FR-024), anomaly flagging (FR-025), bulk suspension by device class (FR-031), and vulnerability window data flagging (FR-032).

**Affected Requirements:** FR-020 (Authenticate Device via mTLS), FR-021 (Authorize Device Actions via OPA), FR-022 (Validate Sensor Readings), FR-024 (Topic ACL Enforcement), FR-025 (Anomaly Flagging), FR-031 (Bulk Device Suspension by Class), SEC-001 (mTLS for All Device Communication)

**Evidence:**
- 0xd — 57% vulnerable, 98% unencrypted, 1.2M devices ([source](https://unit42.paloaltonetworks.com/iot-threat-report-2020/), 2020)
- 0xe — 60% breaches from firmware, 88% increase in HW vulns ([source](https://jumpcloud.com/blog/iot-security-risks-stats-and-trends-to-know-in-2025), 2025)
- 0xf — Mirai: 62 credentials, hundreds of thousands of devices, 1 Tbit/sec ([source](https://blog.cloudflare.com/inside-mirai-the-infamous-iot-botnet-a-retrospective-analysis/), 2023)

---

## 18d. Legal and Regulatory Risks

These risks arise from the legal and regulatory environment in which the platform operates.

### RISK-008: Location Privacy Re-Identification (0x1d)

**Probability:** Medium | **Impact:** Critical | **Category:** Legal

**Description:** Research by de Montjoye et al. demonstrated that just 4 spatiotemporal data points uniquely identify 95% of individuals in a dataset of 1.5 million people. At country scale (60 million), 93% remain re-identifiable. GPS truncation to 4 decimal places yields areas of approximately 11 x 9.5 meters — still too precise for anonymity. The authors conclude that truncation is not a robust anonymization strategy; differential privacy is recommended instead (0x9, Unique in the Crowd, *Scientific Reports*, 2013; follow-up in *Nature Communications*, 2021).

GDPR treats location data as PII. If the platform collects device geolocation without adequate privacy controls, it creates legal liability for both Corewood and participating institutions.

**Mitigation:** Architectural separation of contributor identity from observation data (SEC-004). Configurable spatial resolution per campaign. Automated privacy risk assessment for requests at sub-1km resolution. No raw contributor locations in exported datasets (FR-027).

**Affected Requirements:** FR-027 (Separate Contributor Identity from Data), SEC-004 (Data Privacy and Identity Separation), CO-001 (GDPR Compliance Architecture)
**Related Open Issue:** OI-006 (Spatial Resolution Privacy Tradeoff)

**Evidence:**
- 0x9 — 4 spatiotemporal points identify 95% of individuals; truncation insufficient ([source](https://www.nature.com/articles/srep01376), 2013)

---

### RISK-009: Sweepstakes Regulatory Exposure (0x18)

**Probability:** Medium | **Impact:** Medium | **Category:** Legal

**Description:** Under U.S. law, prize + chance + consideration = illegal lottery. States including New York and Florida require registration and bonding when prize pools exceed $5,000. FTC fines range from $5,000 to $50,000 per violation; the Publishers Clearing House settlement was $18.5 million (0xb, Olshan Frome Wolosky LLP; Klein Moynihan Turco; FTC enforcement actions, 2024).

Additionally, evidence on lottery-style incentives is mixed: lottery incentives reduce uptake by 5.5% compared to guaranteed payments, though they produce the highest engagement persistence (0x10, PMC 2011; BMC Public Health, 2019). A hybrid model — small guaranteed recognition plus occasional lottery — may be needed.

**Mitigation:** Legal review of sweepstakes model across target jurisdictions before launch. Partner with a sweepstakes administration service. Design the incentive system modularly so sweepstakes can be restricted by geography without affecting the rest of the recognition pipeline. Consider a hybrid approach: small guaranteed rewards plus occasional lottery.

**Affected Requirements:** FR-036 (Manage Sweepstakes Entries), FR-034 (Compute Contribution Score), FR-035 (Award Badges and Recognition)
**Related Open Issue:** OI-002 (Sweepstakes Legal Compliance)

**Evidence:**
- 0xb — Prize+chance+consideration = illegal lottery; FTC fines $5K–$50K ([source](https://www.olshanlaw.com/sweepstakes-law-basics), 2024)
- 0x10 — Lottery uptake 5.5% lower than guaranteed; highest engagement persistence ([source](https://pmc.ncbi.nlm.nih.gov/articles/PMC3207198/), 2019)

---

## 18e. Risk Summary

| Risk ID | Name | Category | Probability | Impact | Key Evidence |
|---------|------|----------|-------------|--------|-------------|
| RISK-001 | Insufficient Scitizen Adoption | Supply-side | Medium | Critical | 90-9-1 rule (0x2); 82% first-contribution retention (0x4) |
| RISK-002 | Contributor Retention Failure | Supply-side | Medium | High | 72% active 1 month only (0x1); 5.25% long-term (0x1) |
| RISK-003 | Researcher Non-Adoption | Demand-side | Medium | Critical | 75.3% zero publications (0x3) |
| RISK-004 | Institutional Compliance Blocks | Demand-side | Medium | High | 22 adoption barriers (0xc) |
| RISK-005 | Consumer Sensor Quality | Technical | High | Critical | R²=0.43 raw PM2.5 (0x5); no co-location studies (0x8) |
| RISK-006 | Clock Drift | Technical | Medium | High | 8 min/day ESP32 drift (0xa) |
| RISK-007 | IoT Device Compromise | Technical | Medium | Critical | 57% vulnerable (0xd); 60% firmware breaches (0xe) |
| RISK-008 | Location Re-Identification | Legal | Medium | Critical | 4 points = 95% identified (0x9) |
| RISK-009 | Sweepstakes Regulatory | Legal | Medium | Medium | FTC fines $5K–$50K (0xb) |

### Risk Matrix

|  | **Medium Impact** | **High Impact** | **Critical Impact** |
|---|---|---|---|
| **High Probability** | | | RISK-005 |
| **Medium Probability** | RISK-009 | RISK-002, RISK-004, RISK-006 | RISK-001, RISK-003, RISK-007, RISK-008 |

### Data Gaps

The following evidence nodes explicitly document gaps that require validation work:

- **0x8** — No peer-reviewed co-location study exists for Ecowitt or Davis weather stations. This gap directly affects RISK-005 and OI-001.
- **0xc** — No quantitative adoption data for institutional open data platforms. The 22 barriers are qualitative (TOE framework). This limits the probability estimate for RISK-004.

---

*Next: [Section 19 — Costs](./19_costs.md)*
