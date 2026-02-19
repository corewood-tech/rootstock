# ROOTSTOCK by Corewood

## Requirements Specification — Section 20: User Documentation and Training

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> User documentation and training requirements define what documentation artifacts the platform must produce, who they serve, and what standards they must meet. Each artifact is sourced from external evidence and constrained by requirements captured in the user documentation knowledge graph (`grapher/schema/rootstock_userdocs.graphql`). Every claim traces to a graph node with an external citation.

> **Knowledge graph reference**: User documentation graph at `grapher/schema/rootstock_userdocs.graphql`. Start with `GRAPH=userdocs podman compose -f grapher/compose-grapher.yml up -d` and query at `http://localhost:18080`.

---

## 20a. Researcher Documentation

Researchers are the demand side of the platform — individuals who define data needs and analyze collected data on behalf of research institutions (0x1). Without documentation that enables self-service campaign creation, the platform risks researcher non-adoption (RISK-003, 0x2f). Two documentation artifacts serve this audience.

### DOC-001: Campaign Management Guide (0x39)

**Audience:** Researcher (0x1) | **Reading Level:** Professional | **Format:** Web-based, searchable, versioned with platform

**Description:** Step-by-step guide covering campaign creation, parameter definition, region configuration, time window setup, publication, quality monitoring, and data export. Includes worked examples for common campaign types (0x39).

**Rationale:** 75.3% of citizen science projects produce zero publications, with an average time to first publication of 9.15 years (RISK-003, Section 18). If researchers cannot create campaigns without support, they will not adopt the platform. The guide must enable complete self-service campaign creation.

**Fit Criterion:** A researcher with no prior Rootstock experience can create and publish a campaign by following the guide alone, without contacting support. Validated by usability testing — median completion time under 15 minutes (US-002, 0x2a).

**Constrained by:**
- US-002 — Campaign Creation Under 15 Minutes (0x2a): 5 test researchers with no prior experience complete campaign creation and publication with median time under 15 minutes
- Open Source Constraint (0x25): documentation must be openly accessible alongside the platform

**Sourced from:**
- Diataxis Framework (0x9): structures documentation into tutorials, how-to guides, reference, and explanation
- Write the Docs Minimum Viable Documentation (0xf): defines the minimum documentation set for an open-source project

**Mitigates:** RISK-003 — Researcher Non-Adoption (0x2f)

**Cross-ref:** FR-005, FR-006, FR-007, FR-008, FR-009, FR-010, FR-026

---

### DOC-002: Data Quality and Provenance Guide (0x3a)

**Audience:** Researcher (0x1) | **Reading Level:** Professional | **Format:** Web-based documentation with downloadable PDF

**Description:** Explanation of the validation pipeline, provenance metadata fields, data export formats, identity separation, and spatial resolution configuration. Describes FAIR principles alignment and how to cite Rootstock data in publications (0x3a).

**Rationale:** Researchers need to trust that exported data meets publication-grade standards. Institutional compliance offices require documentation of provenance and privacy controls before approving platform use. Without this guide, adoption is blocked by the 22 institutional barriers identified in RISK-004 (Section 18, 0x30).

**Fit Criterion:** A researcher can identify and explain every provenance field in an exported dataset by referring to this guide (0x3a).

**Constrained by:**
- CO-003 — Data Provenance for Publication (0x28): exported data includes provenance metadata sufficient for peer review, following W3C PROV or equivalent, supporting FAIR principles
- Open Source Constraint (0x25): documentation must be openly accessible

**Mitigates:** RISK-004 — Institutional Compliance Blocks Adoption (0x30)

**Cross-ref:** FR-022, FR-023, FR-026, FR-027, CO-003

---

## 20b. Scitizen Documentation

Scitizens are anyone willing to donate time, effort, and physical exertion to further scientific research by connecting personal IoT devices to active campaigns (0x2). Four documentation artifacts serve this audience. All scitizen-facing documentation is constrained by the 8th-grade reading level requirement (CU-003, 0x2c), multi-language support (0x2b), responsive design (0x26), and the 10-minute enrollment target (US-001, 0x1d).

### DOC-003: Getting Started Guide (0x34)

**Audience:** Scitizen (0x2) | **Reading Level:** 8th grade (0x2c) | **Format:** In-app onboarding flow plus web-based guide, available in all supported languages

**Description:** Account registration, campaign browsing, enrollment code generation, Tier 1 direct enrollment, Tier 2 proxy enrollment, and first campaign enrollment. Written at 8th-grade reading level with visual step-by-step screenshots (0x34).

**Rationale:** IoT device onboarding is identified as the hardest UX step — involving repeated authentications and gateway processes that differ device to device (0x19). The 90-9-1 participation inequality rule predicts 90% of registrants will lurk (RISK-001, Section 18), but participants who make a first contribution retain at 82% (Section 18, 0x4 in analysis graph). This makes the registration-to-first-data path the highest-leverage documentation target. Zooniverse tutorial guidance recommends keeping onboarding tutorials short and focused (0x7), while Foldit demonstrated that scaffolded onboarding improves initial task completion (0x1b).

**Fit Criterion:** A scitizen with no prior experience completes registration, device enrollment, and first data submission by following the guide alone. Validated by usability testing — median completion time under 10 minutes, no user over 15 minutes (US-001, 0x1d).

**Constrained by:**
- US-001 — First-Use Device Enrollment Under 10 Minutes (0x1d): 5 test users complete registration through first reading with median under 10 minutes
- CU-003 — Inclusive Terminology and 8th Grade Reading Level (0x2c): no user-facing text uses "volunteer" or "amateur" for scitizens; all onboarding text tested at 8th-grade readability
- Multi-Language Interface Support (0x2b)
- Responsive Design (0x26): mobile-first layout

**Sourced from:**
- NIH 8th Grade Reading Level Recommendation (0xc)
- Zooniverse Tutorial Length Guidance (0x7)
- Foldit Scaffolded Onboarding (0x1b)
- IoT Onboarding Is Hardest UX Step (0x19)
- IoT Devices Need Physical and Digital Documentation (0x10)

**Mitigates:** RISK-001 — Insufficient Scitizen Adoption (0x2e)

**Cross-ref:** FR-011, FR-012, FR-013, FR-014, FR-015, FR-017

---

### DOC-004: Contribution and Recognition Guide (0x35)

**Audience:** Scitizen (0x2) | **Reading Level:** 8th grade (0x2c) | **Format:** In-app profile section plus web-based guide

**Description:** Contribution score calculation, badge descriptions and milestones, sweepstakes mechanics, and how scitizen data contributes to published research. Emphasizes scientific participation, not charity (0x35).

**Rationale:** Poor usability reduces volunteer motivation — when contributors cannot see the impact of their work, engagement drops (0x14). SciStarter's training modules with micro-credentials demonstrate that visible progress markers sustain participation (0x12). The guide must make the connection between individual data contributions and research outcomes explicit, reinforcing intrinsic motivation through informational feedback rather than controlling rewards.

**Fit Criterion:** A scitizen can explain what their contribution score means and how to earn badges by reading this guide (0x35).

**Constrained by:**
- US-003 — Clear Feedback on Data Contribution (0x24): scitizen dashboard shows accepted reading count, active campaign count, contribution score, and recent badges
- CU-003 — Inclusive Terminology and 8th Grade Reading Level (0x2c)

**Sourced from:**
- SciStarter Training Module with Micro-credentials (0x12)
- Poor Usability Reduces Volunteer Motivation (0x14)

**Mitigates:** RISK-001 — Insufficient Scitizen Adoption (0x2e)

**Cross-ref:** FR-034, FR-035, FR-036

---

### DOC-005: Device Troubleshooting Guide (0x3b)

**Audience:** Scitizen (0x2) | **Reading Level:** 8th grade (0x2c) | **Format:** Web-based searchable knowledge base; error messages link directly to relevant articles

**Description:** Common error messages and meanings, certificate renewal process, grace period behavior, device status explanations, and remediation steps for each failure mode (0x3b).

**Rationale:** IoT device maintenance is fragile — connectivity changes, certificate expiry, and firmware updates each produce failure modes that scitizens must resolve without support. ETSI EN 303 645, the first global cybersecurity standard for consumer IoT, explicitly requires clear maintenance guidance as part of its 33 cybersecurity requirements (0x18). IoT onboarding is the hardest UX step (0x19), and troubleshooting is effectively re-onboarding after failure. Every error message must lead to a resolution path, not a dead end.

**Fit Criterion:** Every user-facing error message has a corresponding troubleshooting article. Scitizens can resolve common issues — expired certificate, WiFi change, firmware update — without contacting support (0x3b).

**Constrained by:**
- US-004 — Actionable Error Messages (0x1e): every user-facing error includes what failed, why, and what to do next; no generic 500 errors
- SEC-006 — Connection Diagnostics (0x27): structured diagnostics for MQTT, enrollment, and auth failures
- CU-003 — Inclusive Terminology and 8th Grade Reading Level (0x2c)

**Sourced from:**
- ETSI EN 303 645 IoT Documentation Requirement (0x18)
- IoT Onboarding Is Hardest UX Step (0x19)

**Mitigates:** RISK-001 — Insufficient Scitizen Adoption (0x2e)

**Cross-ref:** FR-028, FR-029, US-004, SEC-006

---

### DOC-006: Privacy Policy and Consent Documentation (0x3e)

**Audience:** Scitizen (0x2), Institution Administrator (0x3) | **Reading Level:** 8th grade (0x2c) | **Format:** Web-based legal documentation, versioned

**Description:** GDPR rights explanation, consent process documentation, data portability format, erasure procedures, location data handling as PII, and IRB-compatible consent language per campaign (0x3e).

**Rationale:** This artifact serves two audiences with different needs. Scitizens need to understand their GDPR rights — erasure within 30 days, data portability, and consent withdrawal (CO-001, 0x21) — in plain language that passes 8th-grade readability. Institution administrators need consent language that passes IRB review, with per-campaign informed consent records that are timestamped, versioned, and exportable for audit (CO-002, 0x2d). IRBs require research consent forms at a 6th-to-8th grade reading level (0xe). GDPR compliance documentation carries its own cost burden — $5K–$50K for small organizations, up to $1.4M for mid-size (0x32) — making reusable, pre-approved consent templates a cost-saving investment.

**Fit Criterion:** A scitizen can locate and understand their GDPR rights (erasure, portability, consent withdrawal) within the privacy documentation. Consent language passes IRB review at 2–3 test universities (0x3e).

**Constrained by:**
- CO-001 — GDPR Compliance Architecture (0x21): right to erasure within 30 days, granular per-campaign consent, data portability, location data as PII
- CO-002 — IRB-Compatible Consent Model (0x2d): informed consent per-campaign with timestamp, version, scope; consent records exportable for IRB audit
- CU-003 — Inclusive Terminology and 8th Grade Reading Level (0x2c)

**Sourced from:**
- IRB 6th-8th Grade Reading Requirement (0xe)

**Mitigates:** RISK-004 — Institutional Compliance Blocks Adoption (0x30)

**Cost driver:** GDPR Compliance Documentation Cost (0x32)

**Cross-ref:** CO-001, CO-002, SEC-004, FR-027

---

## 20c. Institution Administrator Documentation

Institution administrators manage organizational tenants, researcher access, and compliance configuration (0x3). Two artifacts serve this audience: the Institution Onboarding Guide (0x33) and the shared Privacy Policy and Consent Documentation (DOC-006, 0x3e, documented in Section 20b).

### DOC-007: Institution Onboarding Guide (0x33)

**Audience:** Institution Administrator (0x3) | **Reading Level:** Professional | **Format:** Web-based documentation

**Description:** Organization creation, hierarchy configuration, role definition, researcher invitation, and compliance configuration. Emphasizes the zero-integration requirement (0x33).

**Rationale:** No shared context exists across participating institutions — no common system, platform, protocol, or data format (CON-003, 0x29). Onboarding must require zero changes to institutional systems. A systematic review identified 22 barriers to institutional adoption of open data platforms (RISK-004, Section 18). If each institution requires bespoke negotiation, onboarding becomes non-scalable. This guide must demonstrate that an administrator can self-service the entire onboarding path — from tenant creation through researcher invitation — without contacting Corewood support.

**Fit Criterion:** An institution administrator can onboard their organization and invite researchers without contacting Corewood support (0x33).

**Constrained by:**
- CON-003 — No Shared Context Exists (0x29): onboarding requires zero changes to institutional systems
- CO-002 — IRB-Compatible Consent Model (0x2d): consent configuration must be documented for IRB review
- Open Source Constraint (0x25)

**Mitigates:** RISK-004 — Institutional Compliance Blocks Adoption (0x30)

**Cross-ref:** FR-001, FR-002, FR-003, FR-004, CON-003

**See also:** DOC-006 (Privacy Policy and Consent Documentation) also serves this audience.

---

## 20d. Platform Operator Documentation

Platform operators deploy and maintain Rootstock instances — responsible for infrastructure, security, and incident response (0x4). Two artifacts serve this audience.

### DOC-008: Deployment Guide (0x3c)

**Audience:** Platform Operator (0x4) | **Reading Level:** Technical | **Format:** README and documentation in the source repository

**Description:** Container-based deployment, environment configuration, CA setup, OPA policy configuration, MQTT broker integration, database setup, and observability stack configuration (0x3c).

**Rationale:** Incomplete or outdated documentation is observed by 93% of open-source respondents (0x6, GitHub Open Source Survey 2017), and incomplete or confusing documentation is the top complaint about open-source software (0x17, Google Open Source Blog). For an open-source platform, the deployment guide is the first-contact artifact — if an operator cannot deploy the stack from the public repository alone, adoption fails. The guide must cover the full path from clone to running stack with zero proprietary tools or undocumented steps.

**Fit Criterion:** An operator can deploy the full Rootstock stack from the public repository using only the deployment guide and a container runtime. No proprietary tools or undocumented steps (0x3c).

**Constrained by:**
- OP-001 — Container-Based Deployment (0x1f): all components deployable as OCI containers, single compose command starts full stack, health checks within 60 seconds
- MT-001 — Open Source Licensing (0x20): all source, configuration, infrastructure-as-code, and documentation under OSI-approved license; zero proprietary runtime dependencies
- Open Source Constraint (0x25)

**Sourced from:**
- 93% Observe Incomplete Documentation (0x6, GitHub Open Source Survey 2017)
- Documentation Is Top OSS Complaint (0x17, Google Open Source Blog)
- Write the Docs Minimum Viable Documentation (0xf)

**Cross-ref:** OP-001, SEC-005, CON-001

---

### DOC-009: Security Operations Runbook (0x3d)

**Audience:** Platform Operator (0x4) | **Reading Level:** Technical | **Format:** Operational runbook in the source repository

**Description:** Bulk device suspension procedures, vulnerability window data flagging, device reinstatement, certificate lifecycle monitoring, OPA policy management, and HSM key ceremony procedures (0x3d).

**Rationale:** 57% of IoT devices are vulnerable to medium- or high-severity attacks, and 60% of security breaches originate from unpatched firmware (RISK-007, Section 18). When a firmware vulnerability is discovered, the operator must suspend affected devices, flag data from the vulnerability window, and reinstate devices after patching — all under time pressure. Without a runbook, operators make ad-hoc decisions that increase blast radius and recovery time.

**Fit Criterion:** An operator can execute a bulk device suspension in response to a firmware vulnerability by following the runbook, without ad-hoc decision-making (0x3d).

**Constrained by:**
- MT-001 — Open Source Licensing (0x20)
- Open Source Constraint (0x25)

**Cross-ref:** FR-028, FR-029, FR-030, FR-031, FR-032, FR-033

---

## 20e. Developer and Integrator Documentation

Developers and integrators build integrations, extend the platform, or contribute to the open-source project — including device firmware developers integrating mTLS enrollment and MQTT data submission (0x5). Three artifacts serve this audience.

### DOC-010: API Reference (0x36)

**Audience:** Developer / Integrator (0x5) | **Reading Level:** Technical | **Format:** Auto-generated from protobuf definitions, hosted alongside the platform

**Description:** All public API endpoints with request/response schemas, authentication requirements, error codes, and rate limits. Generated from protobuf definitions via Buf Schema Registry (0x36).

**Rationale:** Connect RPC separates package-level API reference from narrative documentation and is cURL-friendly unlike vanilla gRPC (0x11). The Buf Schema Registry automatically generates documentation on every protobuf module push, providing live documentation for every commit with syntax highlighting, definitions, and searchable UI (0x1c). This means the API reference is never stale — it regenerates on every release from the same protobuf definitions that generate the server and client code.

**Fit Criterion:** Every public API endpoint is documented with at least one example request and response. Documentation is regenerated on every release (0x36).

**Constrained by:**
- MT-001 — Open Source Licensing (0x20)
- Open Source Constraint (0x25)

**Sourced from:**
- Connect RPC Documentation Pattern (0x11)
- Buf Schema Registry Auto-Generated Docs (0x1c)

**Cross-ref:** All public endpoints

---

### DOC-011: Open Source Governance Documentation (0x37)

**Audience:** Developer / Integrator (0x5) | **Reading Level:** Technical | **Format:** Repository root files (CONTRIBUTING.md, SECURITY.md) plus docs site

**Description:** CONTRIBUTING.md with contribution guidelines, SECURITY.md with security disclosure process, release process documentation, and decision-making authority model (0x37).

**Rationale:** 93% of open-source respondents observe incomplete or outdated documentation (0x6), and 60% of maintainers are unpaid — community contribution requires explicit investment in contributor experience (0x31). The open source governance model is currently unresolved (OI-005, 0x23): contribution guidelines, release process, decision-making authority, and security disclosure process are undefined and must be documented before public launch. Without these artifacts, potential contributors face an invisible barrier that suppresses the community growth the project depends on.

**Fit Criterion:** A potential contributor can find contribution guidelines, submit a PR following the documented process, and report a security vulnerability through the documented disclosure channel (0x37).

**Constrained by:**
- MT-001 — Open Source Licensing (0x20): all documentation under OSI-approved license
- OI-005 — Open Source Governance Model, Unresolved (0x23): must be resolved before public launch
- Open Source Constraint (0x25)

**Sourced from:**
- 93% Observe Incomplete Documentation (0x6, GitHub Open Source Survey 2017)
- Documentation Is Top OSS Complaint (0x17, Google Open Source Blog)

**Cost driver:** Development Cost — Documentation as Investment (0x31)

**Cross-ref:** OI-005, MT-001, CON-001

---

### DOC-012: Device Integration Guide (0x38)

**Audience:** Developer / Integrator (0x5) | **Reading Level:** Technical | **Format:** Technical documentation in repository, with protocol diagrams

**Description:** mTLS enrollment protocol (Tier 1 direct, Tier 2 proxy), MQTT topic structure, data submission format, dual-protocol support (MQTT and HTTP/2), and certificate lifecycle from the device perspective (0x38).

**Rationale:** Unlike purely digital products, IoT devices require users to learn physical and digital use simultaneously — documentation must address both physical device setup and protocol integration (0x10). ETSI EN 303 645 requires clear maintenance guidance as part of consumer IoT cybersecurity standards (0x18). Firmware developers integrating new device classes need a complete protocol reference: enrollment, authentication, topic ACLs, and data submission. Undocumented protocol steps are the equivalent of an unverifiable trust boundary.

**Fit Criterion:** A firmware developer can implement the enrollment flow and data submission protocol for a new device class using only this guide and the API reference. No undocumented protocol steps (0x38).

**Constrained by:**
- Open Source Constraint (0x25)

**Sourced from:**
- IoT Devices Need Physical and Digital Documentation (0x10)
- Connect RPC Documentation Pattern (0x11)
- ETSI EN 303 645 IoT Documentation Requirement (0x18)

**Cross-ref:** FR-014, FR-015, FR-020, FR-024

---

## 20f. Training

Training requirements are shaped by three converging findings from citizen science research:

1. **A one-size-fits-all approach to training fails** because it does not draw on individuals' strengths and motivations. Online crowdsourced projects show exceptionally high turnover rates (0x13, West & Pateman, *Citizen Science: Theory and Practice*).
2. **Temporary volunteers have lower classification accuracy** than intermittent or persistent volunteers. More time on the platform correlates with more accurate contributions (0x8, Zooniverse study, *Biological Conservation*).
3. **Scaffolded training with micro-credentials works.** SciStarter offers self-guided online training with badge and micro-credential systems, developed with NLM, ASU, NCSU, Moore Foundation, and IMLS support (0x12).

These findings have direct implications for Rootstock: the platform collects sensor data (not human classifications), so "accuracy" maps to data quality — correct device calibration, proper placement, consistent uptime. Training must address these physical-world skills, not just software navigation.

### Training Strategy

**Role-differentiated training.** Rather than a single training path, documentation and onboarding flows are tailored per audience (0x13):
- **Scitizens** receive scaffolded onboarding (DOC-003) with progressive complexity — registration first, then device enrollment, then campaign participation. The 10-minute target (US-001, 0x1d) applies to the initial path; deeper training unfolds through the contribution guide (DOC-004) and troubleshooting guide (DOC-005).
- **Researchers** receive task-oriented campaign management documentation (DOC-001) with worked examples. The 15-minute target (US-002, 0x2a) constrains the campaign creation path.
- **Operators and developers** receive reference-oriented documentation (DOC-008, DOC-009, DOC-010, DOC-012) optimized for lookup, not linear reading.

**Progressive skill building.** The badge and micro-credential system (FR-035) serves double duty: it recognizes contribution milestones while implicitly training scitizens through increasingly sophisticated participation. SciStarter's model demonstrates that self-guided training with visible credentials sustains engagement (0x12). Zooniverse's finding that persistent volunteers produce more accurate data (0x8) reinforces that the platform should invest in retention-oriented training rather than front-loading complexity.

**Cross-ref:** US-001, US-002, FR-035, DOC-003, DOC-004

---

## 20g. Documentation Standards

All Rootstock documentation must conform to the following standards, derived from external evidence on readability, documentation structure, and citizen science documentation patterns.

### Readability

Scitizen-facing documentation must meet an **8th-grade reading level** (CU-003, 0x2c). This target is supported by converging evidence from three domains:

- **NIH Plain Language Guidelines** (0xc): NIH recommends patient and public materials be written at or below the 8th-grade reading level under the Plain Writing Act of 2010.
- **WCAG 2.1 SC 3.1.5** (0xd): Text above lower secondary education level (grades 7–9) must have supplemental content or a simplified version available (AAA criterion).
- **IRB consent requirements** (0xe): Most IRBs now recommend or require that consent forms and participant-facing materials meet a 6th-to-8th grade reading level.

For unavoidable technical terms, the **AGU Plain Language Summary Guidelines** (0x16) recommend limiting sentences to 15–20 words, avoiding discipline-specific jargon, and including pronunciation guides for complex terms.

### Documentation Structure: Diataxis Framework

All documentation follows the **Diataxis framework** (0x9), which organizes documentation into four types:

| Type | Orientation | Rootstock Application |
|------|------------|----------------------|
| **Tutorials** | Learning-oriented | Getting Started Guide (DOC-003), Institution Onboarding Guide (DOC-007) |
| **How-To Guides** | Goal-oriented | Campaign Management Guide (DOC-001), Device Troubleshooting Guide (DOC-005) |
| **Reference** | Information-oriented | API Reference (DOC-010), Device Integration Guide (DOC-012) |
| **Explanation** | Understanding-oriented | Data Quality and Provenance Guide (DOC-002), Contribution and Recognition Guide (DOC-004) |

Diataxis is adopted by Canonical/Ubuntu, Python, and Django (0x9). Applying it prevents the common failure of mixing tutorial content with reference material, which degrades both.

### Citizen Science Documentation Patterns

Existing citizen science platforms provide models for documentation design:

- **Zooniverse** (0x7): Recommends tutorials of 6 steps or fewer — "the longer your Tutorial, the less likely it is that volunteers will read the whole thing." This constrains the Getting Started Guide (DOC-003) to a brief, focused path.
- **eBird** (0xa): Provides a Help Center with Getting Started, Rules and Best Practices, Community Guidelines, and Review Process sections. Offers a free eBird Essentials online course as a step-by-step guide. This pattern maps to Rootstock's combination of DOC-003 and DOC-004.
- **iNaturalist** (0x1a): Provides Getting Started guide, video tutorials, interactive field guides, community forum, Help Center, and FAQ. Runs an Ambassador Program requiring 100 observations and 100 identifications for peer onboarding — a model for community-driven training at scale.
- **Foldit** (0x1b): Uses scaffolded onboarding through intro puzzles that progressively unlock tools. Community wiki translated into 6 languages. This progressive-complexity model informs the Rootstock training strategy (Section 20f).

### Minimum Viable Documentation

Write the Docs defines minimum viable documentation as answering three questions: why the project exists, how to install it, and how to use it (0xf). "If people cannot figure out how to install or use your code, they will not use it." For Rootstock, this maps to:

1. **Why**: Covered by project README and Section 1 of this specification
2. **How to install**: Deployment Guide (DOC-008)
3. **How to use**: Getting Started Guide (DOC-003) for scitizens, Campaign Management Guide (DOC-001) for researchers

---

## 20h. Documentation Artifact Summary

| DOC ID | Name | UID | Audience | Format | Reading Level | Key Constraints | Mitigates |
|--------|------|-----|----------|--------|--------------|-----------------|-----------|
| DOC-001 | Campaign Management Guide | 0x39 | Researcher | Web-based, versioned | Professional | US-002 (15 min) | RISK-003 |
| DOC-002 | Data Quality and Provenance Guide | 0x3a | Researcher | Web + PDF | Professional | CO-003 (provenance) | RISK-004 |
| DOC-003 | Getting Started Guide | 0x34 | Scitizen | In-app + web, i18n | 8th grade | US-001 (10 min), CU-003 | RISK-001 |
| DOC-004 | Contribution and Recognition Guide | 0x35 | Scitizen | In-app + web | 8th grade | US-003, CU-003 | RISK-001 |
| DOC-005 | Device Troubleshooting Guide | 0x3b | Scitizen | Searchable KB | 8th grade | US-004, SEC-006, CU-003 | RISK-001 |
| DOC-006 | Privacy Policy and Consent | 0x3e | Scitizen, Institution Admin | Web, versioned | 8th grade | CO-001, CO-002, CU-003 | RISK-004 |
| DOC-007 | Institution Onboarding Guide | 0x33 | Institution Admin | Web-based | Professional | CON-003, CO-002 | RISK-004 |
| DOC-008 | Deployment Guide | 0x3c | Platform Operator | Repo README + docs | Technical | OP-001, MT-001 | — |
| DOC-009 | Security Operations Runbook | 0x3d | Platform Operator | Repo runbook | Technical | MT-001 | — |
| DOC-010 | API Reference | 0x36 | Developer | Auto-generated from protobuf | Technical | MT-001 | — |
| DOC-011 | OSS Governance Documentation | 0x37 | Developer | Repo root files + docs | Technical | MT-001, OI-005 | — |
| DOC-012 | Device Integration Guide | 0x38 | Developer | Repo docs + diagrams | Technical | — | — |

### Reading Level Distribution

| Reading Level | Artifact Count | Audiences |
|--------------|---------------|-----------|
| 8th grade | 4 (DOC-003, DOC-004, DOC-005, DOC-006) | Scitizen, Institution Administrator |
| Professional | 3 (DOC-001, DOC-002, DOC-007) | Researcher, Institution Administrator |
| Technical | 5 (DOC-008, DOC-009, DOC-010, DOC-011, DOC-012) | Platform Operator, Developer |

### Cost Drivers

Two cost nodes constrain documentation investment:

- **Development Cost — Documentation as Investment** (0x31): 60% of OSS maintainers are unpaid. Community contribution requires investment in contributor experience and documentation. Documentation is explicitly identified as a cost factor for sustaining open source.
- **GDPR Compliance Documentation Cost** (0x32): Legal documentation and compliance guides required for adopting institutions. GDPR compliance costs $5K–$50K for small organizations, up to $1.4M for mid-size.

---

## 20i. Traceability

Every claim in this section traces to a graph node in the user documentation knowledge graph (`GRAPH=userdocs`). Each graph node carries an external source citation. The traceability chain is:

```
Claim in spec → Graph UID → External source (publication, standard, or platform)
```

### Evidence Summary

| UID | Name | Domain | Source |
|-----|------|--------|--------|
| 0x6 | 93% Observe Incomplete Documentation | OSS | GitHub Open Source Survey 2017 |
| 0x7 | Zooniverse Tutorial Length Guidance | Citizen science | Zooniverse Project Builder Help |
| 0x8 | Temporary Volunteers Lower Accuracy | Citizen science | *Biological Conservation* |
| 0x9 | Diataxis Framework | Documentation | Daniele Procida |
| 0xa | eBird Documentation and Training | Citizen science | eBird Help Center |
| 0xb | GLOBE Observer In-App Guided Documentation | Citizen science | GLOBE Observer |
| 0xc | NIH 8th Grade Reading Level | Readability | NIH Plain Language Guidelines |
| 0xd | WCAG 2.1 Lower Secondary Reading Level | Readability | W3C WCAG 2.1 |
| 0xe | IRB 6th-8th Grade Reading Requirement | Readability | PMC — Using Plain Language in Research |
| 0xf | Write the Docs Minimum Viable Documentation | OSS | Write the Docs |
| 0x10 | IoT Physical and Digital Documentation | IoT UX | IoT For All |
| 0x11 | Connect RPC Documentation Pattern | API docs | Connect RPC official docs |
| 0x12 | SciStarter Training Module | Citizen science | SciStarter |
| 0x13 | One Size Fits All Training Fails | Citizen science | West & Pateman, *CS: Theory and Practice* |
| 0x14 | Poor Usability Reduces Volunteer Motivation | Citizen science | Robinson et al., *JEPM* Vol 64 |
| 0x15 | Four Barriers to Volunteer Retention | Citizen science | *JEPM* Vol 66 No 1 |
| 0x16 | AGU Plain Language Summary Guidelines | Readability | AGU Toolkit |
| 0x17 | Documentation Is Top OSS Complaint | OSS | Google Open Source Blog |
| 0x18 | ETSI EN 303 645 IoT Documentation | Standards | ETSI EN 303 645 V3.1.3 |
| 0x19 | IoT Onboarding Is Hardest UX Step | IoT UX | IoT For All |
| 0x1a | iNaturalist Documentation Types | Citizen science | iNaturalist Help |
| 0x1b | Foldit Scaffolded Onboarding | Citizen science | PMC — Online Citizen Science Games |
| 0x1c | Buf Schema Registry Auto-Generated Docs | API docs | Buf Docs |

---

*Previous: [Section 19 — Costs](./19_costs.md)*
