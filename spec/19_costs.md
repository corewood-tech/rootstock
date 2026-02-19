# ROOTSTOCK by Corewood

## Requirements Specification — Section 19: Costs

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Costs are the resources — financial, operational, and organizational — required to build and run the product. Each cost category is sourced from external evidence captured in the analysis knowledge graph (`grapher/schema/rootstock_analysis.graphql`). Every claim traces to a graph node with an external citation.

> **Knowledge graph reference**: Analysis graph at `grapher/schema/rootstock_analysis.graphql`. Start with `GRAPH=analysis podman compose -f grapher/compose-grapher.yml up -d` and query at `http://localhost:18080`.

---

## 19a. Development Costs

### COST-001: Development — Open Source Model (0x22)

**Category:** Development | **Nature:** Ongoing

**Description:** Corewood bears the initial development cost. Every technology choice must trace to a requirement (CON-002) and no proprietary runtime dependencies are permitted (CON-001). Community contribution is expected to scale development capacity beyond Corewood's team, but this model has well-documented sustainability challenges.

**Evidence-Based Scaling:**
- 60% of open source maintainers are unpaid (0x11, OpenSSF, 2025).
- The OSU Open Source Lab sustains 160+ projects on approximately $250,000/year, providing a reference for lean infrastructure costs (0x11, Linux Foundation/OSU OSL, 2025).
- Security audits cost $5,000–$15,000 per penetration test (0x11, OpenSSF, 2025).

The open source model shifts cost from licensing to community infrastructure, governance, and security auditing. The 60% unpaid maintainer statistic underscores that community contribution is not free — it requires investment in contributor experience, documentation, and governance (see OI-005).

**Affected Requirements:** MT-001, MT-002
**Affected Constraints:** CON-001 (Open Source), CON-002 (First Principles Design)

**Evidence:**
- 0x11 — 60% unpaid maintainers; $250K/yr for 160+ projects; audits $5K–$15K ([source](https://openssf.org/blog/2025/09/23/open-infrastructure-is-not-free-a-joint-statement-on-sustainable-stewardship/), 2025)

---

## 19b. Infrastructure Costs

### COST-002: Self-Operated Certificate Authority (0x23)

**Category:** Infrastructure | **Nature:** One-time + operational

**Description:** The platform requires a self-operated CA because external certificate authorities do not issue client certificates at IoT scale, and using one would create a proprietary dependency that violates CON-001. This shifts cost from per-certificate fees to HSM hardware and operational overhead.

**Evidence-Based Pricing:**
- **Development:** SoftHSM is free and sufficient for local development and testing (0x13, Smallstep docs).
- **Production hardware:** YubiHSM 2 costs approximately $650 per unit — a fixed, one-time cost (0x12, Yubico Store, 2026).
- **Cloud alternative:** AWS CloudHSM costs $876–$1,825/month; Azure and Google Cloud KMS are similarly priced (0x12, AWS/Azure/Google Cloud pricing, 2026).
- **CA software:** step-ca is free under Apache 2.0 (0x13, Smallstep GitHub).
- **Data gap:** No IoT-specific operational cost case studies exist for self-operated CAs (0x13). Operational cost (key ceremony, monitoring, rotation procedures) must be estimated during pilot.

The cost decision is between a ~$650 one-time hardware cost (YubiHSM 2) versus $10,500–$21,900/year for cloud HSM. For a project with CON-001 constraints and no revenue model, hardware HSM is the clear choice for production.

**Affected Requirements:** SEC-005 (HSM Protection for CA Keys), FR-014 (Direct Device Enrollment), FR-028 (Automated Certificate Renewal)

**Evidence:**
- 0x12 — YubiHSM $650; cloud HSM $876–$1,825/mo ([source](https://aws.amazon.com/cloudhsm/pricing/), 2026)
- 0x13 — step-ca free (Apache 2.0); no IoT operational cost data ([source](https://smallstep.com/docs/step-ca/), 2026)

---

### COST-003: Compute and MQTT Infrastructure (0x24)

**Category:** Infrastructure | **Nature:** Ongoing

**Description:** The platform must sustain 10,000 validated readings per second (PE-001) with 99.9% availability (PE-004). Key cost drivers are the MQTT broker, validation pipeline compute, time-series storage, and OPA policy evaluation.

**Evidence-Based Pricing:**
- **Managed MQTT:** At 10,000 messages/second, managed MQTT services (AWS IoT Core, HiveMQ Cloud) cost approximately $26,000/month (0x14, EMQX Open MQTT Benchmarking Comparison, 2023).
- **Self-hosted MQTT:** Self-hosted brokers reduce cost to the underlying VM/container infrastructure. Mosquitto is single-threaded with a maximum throughput of approximately 37,000 messages/second. EMQX supports 100 million connections but requires horizontal scaling (0x14, EMQX benchmarks, 2023).
- **Per-device costs:** Industry benchmarks range from $0.01 to $5.00 per device per month depending on the pricing model — messaging at $1 per million messages, storage at $0.30/GB-month (0x17, Monetizely IoT Cost Breakdown, 2024; AWS IoT Core Pricing).

The $26,000/month managed versus VM-only self-hosted differential is the most significant infrastructure cost decision. Given CON-001 (open source) and the implicit budget constraint (Section 2g — platform cost must not become the problem it was designed to solve), self-hosted MQTT is the expected path.

**Affected Requirements:** PE-001 (Data Ingestion Throughput), PE-002 (Ingestion Latency), PE-004 (Platform Availability), PE-005 (OPA Authorization Latency), FR-023 (Persist Valid Readings with Provenance)

**Evidence:**
- 0x14 — Managed MQTT ~$26K/mo at 10K msg/s; self-hosted = VM cost ([source](https://www.emqx.com/en/blog/open-mqtt-benchmarking-comparison-mqtt-brokers-in-2023), 2023)
- 0x17 — Per-device $0.01–$5/mo; $1/1M messages ([source](https://aws.amazon.com/iot-core/pricing/), 2024)

---

## 19c. Incentive Costs

### COST-004: Scitizen Incentives — Sweepstakes Model (0x25)

**Category:** Incentive | **Nature:** Ongoing

**Description:** The cost model (Section 1.2) specifies sweepstakes entries, experiences, and recognition rather than per-unit compensation. This is grounded in Prospect Theory: people overweight small probabilities, making lottery-style rewards psychologically disproportionate to their actual cost. However, the evidence is more nuanced than the theory suggests.

**Evidence-Based Scaling:**
- Lottery-style incentives reduce initial uptake by 5.5% compared to guaranteed payments, but produce the highest engagement persistence over time (0x10, PMC 2011; BMC Public Health, 2019). Health behavior studies used expected values of $1.40–$2.80/day for lottery incentives.
- The traditional cost baseline: citizen science data collection has been valued at $667 million to $2.5 billion annually across 388 biodiversity projects, with per-observation costs of 37–300 EUR when setup costs are included (0x15, Theobald et al., *Biological Conservation*, 2015).

The fixed prize pool model means incentive costs do not scale with data volume. As more scitizens join, per-unit data cost approaches zero — this is the core economic advantage. However, the 5.5% uptake reduction suggests a hybrid approach (small guaranteed recognition plus occasional lottery) may optimize both adoption and retention.

**Affected Requirements:** FR-034 (Compute Contribution Score), FR-035 (Award Badges and Recognition), FR-036 (Manage Sweepstakes Entries)

**Evidence:**
- 0x10 — Lottery uptake 5.5% lower; highest engagement persistence ([source](https://pmc.ncbi.nlm.nih.gov/articles/PMC3207198/), 2019)
- 0x15 — $667M–$2.5B/yr citizen science value; 37–300 EUR/observation ([source](https://www.sciencedirect.com/science/article/pii/S0006320714004029), 2015)

---

## 19d. Compliance Costs

### COST-005: GDPR and Compliance (0x26)

**Category:** Compliance | **Nature:** Design-time + operational

**Description:** GDPR compliance (CO-001) and IRB compatibility (CO-002) are architectural requirements, not aftermarket additions. The cost is primarily in design and implementation: privacy separation architecture, consent management, right-to-erasure pipeline, and audit trail storage. Open-source auditability (CON-001) reduces the compliance burden for each adopting institution.

**Evidence-Based Pricing:**
- GDPR compliance costs range from $5,000–$50,000 for small organizations, to $1.4 million for mid-size organizations (0x16, ItsASAP GDPR Cost Guide, 2024; Usercentrics, 2025).
- A Data Protection Officer costs $40,000–$150,000/year in-house, or $60–$720/month outsourced (0x16).
- GDPR caused a 20% increase in the average cost of data across industries (0x16, MIT Sloan, 2024).
- **Data gap:** No OSS-specific GDPR compliance cost data exists (0x16). The standard enterprise figures likely overestimate costs for an open-source project with no revenue model, but underestimate the volunteer engineering effort required.

For Rootstock, the primary compliance cost is engineering time to build privacy-by-design architecture. Ongoing costs are audit trail storage (scales with platform activity) and legal review of consent models across jurisdictions.

**Affected Requirements:** CO-001 (GDPR Compliance Architecture), CO-002 (IRB-Compatible Consent Model), SEC-004 (Data Privacy and Identity Separation)
**Affected Constraints:** CON-001 (Open Source)

**Evidence:**
- 0x16 — GDPR costs $5K–$50K small, $1.4M mid-size; DPO $40K–$150K/yr; no OSS data ([source](https://www.itsasap.com/blog/cost-gdpr-compliance), 2024)

---

## 19e. Economic Context

### COST-006: Citizen Science Economic Value — The Inverted Cost Model (0x21)

**Category:** Economics | **Nature:** Reference

**Description:** This is not a platform cost but the economic context that justifies the platform's existence. Traditional field data collection consumes 50–80% of research grant budgets. Rootstock inverts this model: infrastructure costs are fixed, and the marginal cost per additional reading approaches zero as more scitizens contribute.

**Evidence-Based Context:**
- Citizen science produces an estimated $667 million to $2.5 billion in annual value across 388 biodiversity monitoring projects alone (0x15, Theobald et al., *Biological Conservation*, 2015).
- However, per-observation costs range from 37 to 300 EUR when setup, training, and quality assurance costs are included (0x15).
- The platform inverts the cost structure: traditional field collection has variable cost proportional to data volume (more data = proportionally more cost), while the platform has fixed infrastructure cost with marginal cost approaching zero per additional reading (more data = lower unit cost).

The implicit budget constraint (Section 2g) requires that the platform's own operational cost remains low enough that it does not become the cost problem it was designed to solve. This means total infrastructure + incentive cost must remain well below the traditional per-observation cost of 37–300 EUR multiplied by the data volume the platform enables.

**Affected Requirements:** None directly — this is the economic rationale, not a cost to manage.

**Evidence:**
- 0x15 — $667M–$2.5B/yr citizen science value; 37–300 EUR/observation ([source](https://www.sciencedirect.com/science/article/pii/S0006320714004029), 2015)

---

## 19f. Cost Summary

| Cost ID | Name | Category | Nature | Key Pricing Data |
|---------|------|----------|--------|-----------------|
| COST-001 | Development (OSS Model) | Development | Ongoing | 60% maintainers unpaid; audits $5K–$15K (0x11) |
| COST-002 | Self-Operated CA | Infrastructure | One-time + operational | YubiHSM $650 vs cloud $876–$1,825/mo (0x12) |
| COST-003 | Compute and MQTT | Infrastructure | Ongoing | Managed $26K/mo vs self-hosted VM-only (0x14) |
| COST-004 | Scitizen Incentives | Incentive | Ongoing | Fixed prize pool; per-unit cost → zero at scale (0x10) |
| COST-005 | GDPR and Compliance | Compliance | Design-time + operational | $5K–$50K small org; no OSS data (0x16) |
| COST-006 | Economic Context | Economics | Reference | $667M–$2.5B/yr citizen science value (0x15) |

### Cost Scaling Model

| Cost | Scales With | Direction |
|------|------------|-----------|
| COST-001 | Developer time, community size | Sublinear — community contribution offsets Corewood effort |
| COST-002 | Fixed (hardware) or monthly (cloud) | Constant — per-cert cost near zero |
| COST-003 | Device count, message volume | Linear with device count; self-hosted reduces slope |
| COST-004 | Prize pool budget (fixed) | Constant — per-unit data cost → zero as supply grows |
| COST-005 | Platform activity (audit trail) | Sublinear — architecture cost is one-time, storage scales |

### Data Gaps

- **0x13** — No IoT-specific operational cost data for self-operated CAs. Pilot must measure key ceremony, monitoring, and rotation overhead.
- **0x16** — No OSS-specific GDPR compliance cost data. Enterprise figures likely overestimate financial cost but underestimate volunteer engineering effort.

---

*Next: [Section 20 — User Documentation and Training](./20_user_documentation.md)*
