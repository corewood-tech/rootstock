# ROOTSTOCK by Corewood

## Requirements Specification — Section 4: Relevant Facts and Assumptions

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Relevant facts are verifiable truths about the world that constrain, inform, or shape the product. Assumptions are beliefs about the world that the team holds to be true but cannot yet verify — they carry risk if wrong. Both categories are derived from research, not guessed at. Every entry in this section is traceable to a source.

> **Knowledge graph reference**: The research behind this section is captured in a persistent Dgraph knowledge graph (`grapher/schema/rootstock_facts.graphql`). Nodes are referenced by UID for traceability. Start Dgraph with `podman compose -f grapher/compose-grapher.yml up -d` and query at `http://localhost:18080`.

---

## 4a. Relevant Facts

### FACT-001: Traditional Field Data Collection Is Expensive and Inefficient

**Statement:** Traditional scientific field data collection requires grant-funded personnel, travel, equipment, and field time. Personnel costs alone consume 50–80% of grant budgets, leaving limited funds for data collection itself. The average annualized NSF research grant is $245,800 (FY2024), and ecology grants from the Division of Environmental Biology average roughly $833,000 over the award period.

**Evidence:**
- Personnel share of grant budgets: 50–80% (Science/AAAS)
- NSF annualized award FY2024: $245,800 median $185,600 (NSF by the Numbers)
- NSF DEB: $100M for ~120 awards in FY2024 (NSF DEB Solicitation)
- In the TRY global plant trait database, 70.2% of removed records were duplicates — indicating massive redundancy in manually collected data (Global Change Biology)

**Impact on Rootstock:** The platform must demonstrate a credible cost advantage. If consumer IoT devices can produce equivalent data quality for a fraction of the cost, the value proposition is defensible. The cost model in Section 1.2 (sweepstakes incentives with diminishing per-unit cost) must be validated against these baseline figures.

**Graph UIDs:** 0x1, 0x6, 0x20, 0x22, 0x25

---

### FACT-002: Geographic and Temporal Coverage Is Structurally Inadequate

**Statement:** Existing biodiversity and environmental datasets cover less than 7% of the Earth's surface at 5km resolution. Coverage is skewed toward wealthy countries and accessible areas. Tropical and subtropical regions — where biodiversity is greatest — are the least covered. Completeness is limited by distance to researchers, local funding, and participation in data-sharing networks.

**Evidence:**
- GBIF/OBIS surface coverage at 5km resolution: <7% (Nature Communications)
- Tropical biodiversity data gap is well-documented (Tropical Conservation Science)
- Australia lost its entire Long-Term Ecological Research Network in 2017 (PLOS Biology)

**Impact on Rootstock:** This is the scale problem that justifies distributed collection. Consumer devices are already deployed in residential areas globally. A network of 250,000+ weather stations (Weather Underground's current network) dwarfs any professional field team. Rootstock's coverage depends on scitizen geographic distribution, not research institution budgets.

**Graph UIDs:** 0x7, 0x21

---

### FACT-003: Consumer IoT Devices Are Functionally Scientific Instruments

**Statement:** The consumer environmental monitoring market is large and growing. Weather stations ($1.4B market, 8.1% CAGR to $2.8B by 2033), air quality monitors (PurpleAir: 30,000+ sensors, $200–300/unit), and soil/water sensors are already deployed at scale. Weather Underground alone aggregates data from 250,000+ personal weather stations. These devices produce measurements comparable to — and in some cases cross-validated against — professional reference stations.

**Evidence:**
- Wireless home weather station market: $1.4B → $2.8B by 2033 (Spherical Insights)
- Environmental sensor market: $1.7B → $3.33B by 2032 (SNS Insider)
- PurpleAir: 30,000+ deployed sensors, data used in peer-reviewed wildfire research (ACS ES&T Letters)
- Weather Underground: 250,000+ personal weather stations (Weather Underground)
- EU Commission: citizen air quality sensors improved by co-location quality control (EC Environment News, 2025)

**Impact on Rootstock:** The supply side exists. The hardware is deployed, powered, connected, and generating data. The gap is not equipment — it is coordination and motive. Rootstock does not need to manufacture or distribute sensors; it needs to connect existing sensors to research demand.

**Graph UIDs:** 0x8, 0x2e, 0x27, 0x28

---

### FACT-004: No Platform Connects Device-Owning Citizens to Researchers

**Statement:** Existing citizen science platforms are either observation-based (iNaturalist, eBird — humans photograph/report), classification-based (Zooniverse — humans label data), or single-parameter (Safecast — radiation only). None provides a general-purpose platform for connecting researchers' structured data needs to citizen-owned IoT devices across multiple sensor types. This is the gap Rootstock fills.

**Evidence:**
- iNaturalist: 300M+ observations, 3.3M observers — but photography-based, not sensor-based
- eBird: 2.1B observations — bird-only, checklist-based
- Zooniverse: 2.7M volunteers, 450+ papers — classification, not collection
- Safecast: 180M+ radiation measurements — closest model but single-parameter, custom hardware
- PurpleAir / Weather Underground: large networks but not research platforms
- No platform exists that lets researchers define multi-parameter campaigns and match them to citizen-owned IoT devices

**Impact on Rootstock:** The competitive landscape confirms the opportunity. Rootstock is not competing with iNaturalist or eBird — those platforms serve different data modalities. The closest analog is Safecast, but Safecast is locked to one parameter. Rootstock must be parameter-agnostic and campaign-driven.

**Graph UIDs:** 0x2, 0x29, 0x2a, 0x2b, 0x26, 0x19

---

### FACT-005: Citizen Science Produces Significant Economic Value

**Statement:** The in-kind value of citizen science volunteers across 388 surveyed projects (1.3M volunteers) is estimated at $0.7–2.5 billion annually. Citizen science involvement produces a 46–77% cost reduction compared to contracted data collection. eBird alone logged 302 million new observations in 2025 from 618,000 active contributors.

**Evidence:**
- Volunteer in-kind value: $0.7–2.5B/year (Phys.org / Silvertown 2022)
- Cost reduction vs. contractors: 46–77% (Phys.org)
- eBird 2025: 302M observations, 618K active eBirders (Cornell Chronicle)
- iNaturalist: 4,000+ research papers (BioScience/Oxford)

**Impact on Rootstock:** The economics work at scale. The challenge is not whether citizen science is valuable — it demonstrably is — but whether the platform can attract and retain enough contributors to reach useful density in target geographies.

**Graph UIDs:** 0x23, 0x24, 0xa

---

### FACT-006: Retention Is the Critical Problem, and the First Contribution Is the Inflection Point

**Statement:** Citizen science platforms suffer high attrition. Project FeederWatch data shows 82% year-over-year retention for contributors who submitted at least one observation, but only 39.7% retention for those who registered but never submitted. CSMON-LIFE found sustained multi-year engagement drops to 5.25%. Contributions follow a power law: ~10% of participants produce ~90% of data.

**Evidence:**
- FeederWatch: 82% retention post-first-contribution, 39.7% without (BioScience/Oxford)
- CSMON-LIFE: 70% nominal retention, 5.25% sustained (MDPI/Sustainability)
- Power-law distribution: ~10% produce ~90% (Sauermann & Franzoni, PNAS 2015)
- Average app Day-1 retention: ~25%, Day-30: ~5–10% (Mixpanel)

**Impact on Rootstock:** Getting a scitizen from registration to first successful data submission is the highest-leverage design problem. IoT devices that submit data automatically after enrollment fundamentally change this equation — the device does the contributing, reducing activation energy to near zero after setup. This is Rootstock's structural advantage over observation-based platforms like iNaturalist.

**Graph UIDs:** 0x4, 0x33, 0x13

---

### FACT-007: Marketing Campaign Mechanics Are Directly Transferable

**Statement:** Marketing campaigns and research campaigns are structurally identical: defined goal, target audience, timeframe, creative assets, incentives, measurement criteria, and iteration. The most successful citizen science events already use campaign mechanics — eBird's Global Big Day (150,000+ participants, 200+ countries) and iNaturalist's City Nature Challenge (76,000+ participants, 2.2M observations, 500+ cities) are marketing campaigns applied to data collection.

**Evidence:**
- eBird Global Big Day 2024: 150,000+ participants across 200+ countries (eBird.org)
- City Nature Challenge 2024: 76,000+ participants, 2.2M observations, 500+ cities (citynaturechallenge.org)
- CMS campaign structure maps directly to research campaigns: objective → data need, audience → geographic/device segment, scheduling → campaign windows, analytics → data quality dashboards
- Gamification increases participation 20–100% across studies (Morschheuser et al, 2017)

**Impact on Rootstock:** The platform's CMS component is not a metaphor — it is a literal campaign management system with researchers as campaign creators and scitizens as the audience. The marketing automation playbook (segmentation, lifecycle automation, A/B testing, re-engagement) is directly applicable.

**Graph UIDs:** 0x9, 0xb, 0x14, 0x12, 0x2d

---

### FACT-008: Sweepstakes Incentives Are Behaviorally Optimal for This Use Case

**Statement:** Prospect Theory (Kahneman & Tversky, 1979) demonstrates that people overweight small probabilities, making lottery-style incentives psychologically more motivating per dollar spent than fixed per-unit payments. Volpp et al. (2008) confirmed this in a randomized controlled trial: lottery-based incentives outperformed per-unit payment for sustained behavior change. Variable-ratio reinforcement schedules (Skinner) produce the strongest sustained engagement patterns.

**Evidence:**
- Prospect Theory: probability overweighting (Econometrica, 1979)
- Lottery vs. fixed payment RCT: lottery group showed greater sustained behavior (JAMA, 2008)
- Variable-ratio reinforcement: highest response rates of all reinforcement schedules (Skinner, foundational behavioral psychology)
- Deci & Ryan Self-Determination Theory: informational rewards (recognition, badges) enhance intrinsic motivation, while controlling rewards can undermine it

**Impact on Rootstock:** The sweepstakes model in Section 1.2 is not just cheaper — it is more effective per dollar at sustaining engagement. The key risk is that extrinsic rewards (even sweepstakes) can crowd out intrinsic motivation if poorly designed. Recognition and visible research impact must be the primary layer; sweepstakes are the supplement, not the core.

**Graph UIDs:** 0x2c, 0x11, 0x10

---

### FACT-009: Location Data Privacy Is a Fundamental Platform Tension

**Statement:** Just 4 spatiotemporal data points uniquely identify 95% of individuals in a dataset of 1.5 million people. Simple GPS truncation is not effective at preventing re-identification. GDPR treats location data as PII, applies to all persons in the EEA regardless of citizenship, and mandates erasure rights that the platform must support. IRBs cannot waive informed consent under GDPR.

**Evidence:**
- 4 points → 95% re-identification (de Montjoye et al., Nature Scientific Reports; PMC 7961185)
- GDPR: location data is PII; right to erasure is absolute (Multiple university IRB guidance)
- IRB cannot waive consent under GDPR (LSU GDPR Consent Guide)
- Future of Privacy Forum: geolocation data is inherently high-risk (FPF)

**Impact on Rootstock:** This creates a fundamental architectural requirement. Contributor identity must be separable from observation data. The consent model must be granular (per-campaign, per-use), versioned, and auditable. Spatial resolution must be configurable per campaign, with researcher-justified precision requests subject to ethics review. Raw contributor locations must never appear in public-facing datasets. This is non-negotiable for any institution with GDPR exposure.

**Graph UIDs:** 0x30, 0x31

---

### FACT-010: IoT Security Cannot Be Trusted at the Device Level

**Statement:** Over 50% of IoT devices have critical exploitable vulnerabilities. 60% of IoT breaches originate from unpatched firmware. 17% of devices contain hardcoded credentials. The average device risk score is 8.98/10 (Forescout, 2025). 540,531 MQTT broker vulnerabilities were found worldwide via Shodan. Consumer device security is a baseline assumption of compromise, not a guarantee.

**Evidence:**
- 50%+ critical vulnerabilities (Forescout 2025 report)
- 60% breaches from unpatched firmware (JumpCloud 2025)
- 17% hardcoded credentials (industry surveys)
- 540,531 MQTT vulnerabilities on Shodan (Wavestone RiskInsight, 2024)
- OWASP IoT Top 10: hardcoded passwords, insecure network services, insecure ecosystem interfaces

**Impact on Rootstock:** The trust boundary in Section 1.4 is confirmed as architecturally correct. The platform must assume every device is potentially compromised. Server-side validation, anomaly detection, cross-device correlation, and cryptographic attestation of data origin are requirements, not nice-to-haves. The mTLS architecture and OPA-enforced authorization already address this — this fact validates those choices.

**Graph UIDs:** 0x5, 0x2f, 0x18

---

### FACT-011: Relevant Standards Exist and Should Be Evaluated

**Statement:** The OGC SensorThings API provides an open, geospatial-enabled standard for IoT data interconnection, built on REST+JSON+MQTT and based on ISO 19156 (Observations and Measurements). The FAIR Principles (Findable, Accessible, Interoperable, Reusable) and TRUST Principles define scientific data management and repository standards. The W3C PROV model defines data provenance interchange. All are potentially relevant to Rootstock.

**Evidence:**
- OGC SensorThings API: REST+MQTT, ISO 19156 data model, FROST-Server reference implementation (OGC)
- FAIR Principles (Nature, 2016)
- TRUST Principles for repositories (Scientific Data, 2020)
- W3C PROV model for provenance interchange (W3C, 2013)
- Citizen Science: Theory and Practice: citizen science projects often lack open access, interoperability, and sustainable infrastructure

**Impact on Rootstock:** Per CON-002 (First Principles Design), these standards should be evaluated on their merits against Rootstock's requirements, not adopted reflexively. SensorThings is the closest existing standard — its data model and MQTT extension align well. FAIR/TRUST compliance would make Rootstock data publishable in major journals. W3C PROV is implementation-ready for provenance chains.

**Graph UIDs:** 0x34

---

### FACT-012: Water Quality and Soil Monitoring Have No Consumer Network Equivalent

**Statement:** No widely deployed consumer water quality or soil monitoring network exists equivalent to PurpleAir (air quality, 30,000+ sensors) or Weather Underground (weather, 250,000+ stations). IoT-compatible water quality and soil sensors exist — Atlas Scientific and DFRobot sell consumer-grade probes with digital interfaces — but no open research network aggregates their data. This is a concrete market gap and an early differentiation opportunity for Rootstock.

**Evidence:**
- No equivalent of PurpleAir or Weather Underground exists for water quality or soil monitoring (market survey, 2025)
- Atlas Scientific: IoT-compatible water quality probes (pH, dissolved oxygen, conductivity, ORP) with UART/I2C interfaces
- DFRobot: consumer soil moisture, pH, and EC sensors with Arduino/ESP32 compatibility
- EPA and USGS water quality monitoring relies on fixed stations with limited geographic coverage
- Soil monitoring is almost entirely institutional (NRCS, ISRIC) with no citizen participation model

**Impact on Rootstock:** Water quality and soil monitoring represent underserved sensor categories where Rootstock could establish a network before a specialized competitor does. These categories also align well with the campaign model — researchers studying watershed health or soil degradation need spatially distributed, temporally continuous data that consumer sensors can provide. However, consumer sensor accuracy for water chemistry and soil parameters is less validated than for weather or air quality, making data quality thresholds especially important for these campaign types.

**Graph UIDs:** 0x35, 0x8

---

### FACT-013: Most Citizen Science Projects Fail to Produce Peer-Reviewed Publications

**Statement:** A meta-analysis of 895 citizen science projects spanning 1890–2018 found that 75% produced zero peer-reviewed publications. Of the 221 projects that did publish, they produced 2,075 papers total — but 5 projects accounted for nearly half of all publications. The average time from project launch to first publication was 9.15 years.

**Evidence:**
- Meta-analysis: 895 projects, 75% with zero publications (ResearchGate)
- 5 projects produced ~50% of all citizen science publications
- Average time to first publication: 9.15 years
- Successful projects (eBird, iNaturalist, Galaxy Zoo) share common traits: structured data collection, clear research questions, institutional backing, and low-friction contribution

**Impact on Rootstock:** This is the failure mode Rootstock's campaign model is designed to prevent. The 75% failure rate is driven by projects that collect data without a specific research question — the data never connects to a publishable hypothesis. By requiring researchers to define campaigns with explicit parameters, regions, time windows, and quality thresholds, Rootstock structures the collection around the research need from the start. The campaign model inverts the typical citizen science flow: instead of "collect data and hope researchers use it," it starts with "researcher defines what they need" and matches contributors to that need.

**Graph UIDs:** 0x32, 0x14, 0x12

---

## 4b. Assumptions

> Each assumption carries risk if wrong. Assumptions should be validated as early as possible and converted to facts or discarded.

### ASM-001: Consumer IoT Device Owners Will Contribute Data to Science If the Platform Exists

**Statement:** We assume that a meaningful percentage of consumer weather station, air quality monitor, and environmental sensor owners will enroll their devices in research campaigns if the process is low-friction, the purpose is clear, and recognition is provided.

**Risk if wrong:** The platform has no contributors. The supply side collapses.

**Validation path:** Early access pilot with existing PurpleAir, Ecowitt, or Davis Instruments owners. Measure enrollment rate, first-contribution rate, and 30-day retention. If <5% of contacted device owners enroll, revisit the value proposition.

**Supporting evidence:** 250,000+ weather stations on Weather Underground suggest willingness to share data. 14% of Millennial/Gen Z adults have participated in citizen science in the past year (Pew Research). But contributing sensor data passively is different from actively classifying images.

---

### ASM-002: Sweepstakes and Recognition Are Sufficient Incentive (No Per-Unit Payment Needed)

**Statement:** We assume that lottery-based incentives, gamification, and visible recognition are sufficient to drive sustained participation without direct per-contribution compensation.

**Risk if wrong:** Contributor retention drops below viable levels. The cost model in Section 1.2 fails, and the platform must introduce per-unit payments that change the economic model.

**Validation path:** A/B test incentive structures during pilot: sweepstakes-only vs. sweepstakes-plus-micro-recognition vs. small per-unit payment. Measure retention at 30, 60, and 90 days.

**Supporting evidence:** Prospect Theory supports lottery efficacy. FeederWatch retention data shows intrinsic motivation matters more than extrinsic rewards. But Rootstock asks for passive device enrollment (lower effort) — the incentive bar may be different from active observation platforms.

---

### ASM-003: Researchers Will Create Campaigns on the Platform

**Statement:** We assume that researchers at institutions will invest time in defining campaigns (parameters, regions, time windows, quality thresholds) on Rootstock, and that they will find the resulting data useful for publication.

**Risk if wrong:** The platform has no demand side. Without campaigns, scitizens have nothing to contribute to.

**Validation path:** Partner with 3–5 research groups during development. Have them define real campaigns. If campaign creation takes >30 minutes or the resulting data is not publication-grade, the researcher UX or data quality pipeline needs work.

**Supporting evidence:** 75% of citizen science projects never produce a paper. The bottleneck is not collection but the connection between collected data and research questions. A structured campaign model — where the researcher defines exactly what they need — may improve this ratio.

---

### ASM-004: Institutional Data Governance Will Not Block Adoption

**Statement:** We assume that research institutions' IRB, ethics, and data governance requirements can be satisfied by Rootstock's consent model, data provenance, and privacy architecture without requiring per-institution customization that fragments the platform.

**Risk if wrong:** Each institution requires bespoke compliance work, making onboarding non-scalable. The platform cannot be "zero changes to the institution's existing systems" (CON-003) if compliance requires deep integration.

**Validation path:** Engage 2–3 university IRBs during architecture phase. Present the consent model, location privacy architecture, and data provenance chain. Get pre-clearance or identify blockers.

**Supporting evidence:** GDPR consent requirements are well-documented and satisfiable with granular, versioned consent. US IRB requirements vary by institution but generally accept informed consent for voluntary data contribution. The open-source, auditable nature of the platform (CON-001) supports institutional trust.

---

### ASM-005: Automated IoT Devices Change the Retention Equation

**Statement:** We assume that because IoT devices submit data automatically after enrollment (unlike platforms where humans must manually observe and report), the retention problem is fundamentally different. A weather station enrolled in a campaign contributes data 24/7 without ongoing human action. The critical conversion is enrollment, not repeated contribution.

**Risk if wrong:** Device owners still disengage — they unplug devices, change WiFi passwords, let subscriptions lapse, or simply forget the platform exists. Passive contribution does not guarantee passive retention.

**Validation path:** Track device uptime and data continuity across enrolled devices. If >20% of devices stop reporting within 60 days for non-technical reasons (not firmware failure or hardware fault), the assumption is challenged.

**Supporting evidence:** No direct precedent exists for this model at scale. Safecast is the closest analog (continuous sensor data from citizen-owned devices), but its community is self-selected and highly motivated. Weather Underground's 250K station network suggests long-term device operation is common among enthusiasts.

---

### ASM-006: The Market Will Accept a Multi-Parameter, Protocol-Agnostic Platform

**Statement:** We assume that a platform supporting arbitrary sensor types, data parameters, and communication protocols (MQTT, HTTP/2) across device tiers will find adoption, despite the complexity this introduces. We assume this generality is a strength, not a liability.

**Risk if wrong:** The platform is too generic to be excellent at any one thing. Specialized platforms (PurpleAir for air quality, Weather Underground for weather) remain preferred because they are purpose-built.

**Validation path:** Launch with 2–3 well-defined sensor categories (weather, air quality, water quality). Measure whether researchers and scitizens prefer the unified platform or continue using specialized networks. If cross-parameter campaigns (e.g., correlating air quality with weather data) generate unique research value, the generality justifies itself.

---

### ASM-007: Consumer Sensor Accuracy Is Sufficient for Peer-Reviewed Research Across Target Parameters

**Statement:** We assume that consumer-grade environmental sensors (weather stations, air quality monitors, water quality probes) produce data of sufficient accuracy and precision to support peer-reviewed research when validated against campaign-defined quality thresholds. We assume this holds across weather parameters (well-validated), air quality (partially validated), and water/soil parameters (largely unvalidated).

**Risk if wrong:** The platform collects high volumes of data that researchers cannot publish. The value proposition collapses — cost savings are meaningless if the data is not usable. Researchers stop creating campaigns.

**Validation path:** For each sensor category, identify peer-reviewed studies that cross-validate consumer devices against reference instruments. Weather: extensive validation exists (Weather Underground co-location studies). Air quality: PurpleAir PM2.5 data has been validated in wildfire research (ACS ES&T Letters). Water quality: limited validation — commission a co-location study with Atlas Scientific probes vs. reference instruments during pilot. If consumer sensor error margins exceed researcher-defined thresholds for >50% of campaign types, the platform must either restrict sensor eligibility or introduce calibration requirements.

**Supporting evidence:** PurpleAir data is used in peer-reviewed wildfire research. EU Commission found citizen air quality sensors improved with co-location quality control. Weather Underground data is referenced in meteorological studies. But no systematic cross-validation exists for consumer water quality or soil sensors in research contexts.

**Graph UIDs:** 0x8, 0x35, 0x27, 0x2e

---

*Next: [Section 5 — The Scope of the Work](./05-scope-of-work.md)*
