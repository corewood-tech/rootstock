# ROOTSTOCK by Corewood

## Requirements Specification — Section 21: Conclusion

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995--2019 the Atlantic Systems Guild Limited.

---

## 21a. What Rootstock Is

Rootstock is an open-source platform that connects researchers who need field data with citizen scientists who own consumer IoT devices capable of producing it. The core mechanism is the **campaign**: a researcher defines what data they need, where, when, and to what quality standard, and the platform matches that request to enrolled devices operated by scitizens in the field.

This specification defines the full system across 20 Volere sections: the problem it solves, the constraints it operates under, the facts and assumptions it rests on, the functional and non-functional requirements it must satisfy, and the risks, costs, and documentation obligations it must manage.

---

## 21b. The Problem Being Solved

Traditional scientific field data collection consumes 50--80% of grant budgets on personnel alone. Existing biodiversity and environmental datasets cover less than 7% of the Earth's surface at 5km resolution. Meanwhile, 250,000+ personal weather stations, 30,000+ air quality sensors, and a growing base of water and soil monitors are already deployed, powered, connected, and generating data that no research platform systematically captures.

Existing citizen science platforms are either observation-based (iNaturalist, eBird), classification-based (Zooniverse), or single-parameter (Safecast). None connects researchers' structured data needs to citizen-owned IoT devices across multiple sensor types. This is the gap.

---

## 21c. What Emerged from the Requirements Process

### The campaign model inverts the failure mode

75% of citizen science projects produce zero peer-reviewed publications. The common pattern is: collect data, then hope a researcher uses it. Rootstock inverts this by starting with the researcher's defined need --- parameters, region, time window, quality thresholds --- and matching contributors to that need. The campaign is the organizing unit, not the data.

### Automated devices change the retention equation

On observation-based platforms, every data point requires active human effort. Retention drops to 5.25% for sustained multi-year engagement. IoT devices submit data automatically after enrollment. The critical conversion is enrollment, not repeated contribution. The 82% post-first-contribution retention rate makes the enrollment-to-first-data path the highest-leverage design problem --- hence the US-001 requirement of under 10 minutes from landing to first reading.

### Security is architectural, not supplementary

Over 50% of IoT devices have critical exploitable vulnerabilities. Every device connection crosses a trust boundary. The specification treats this as a foundational constraint, not an afterthought: mTLS authentication (SEC-001), private keys that never leave the device (SEC-002), OPA-enforced authorization on every action (SEC-003), a self-operated CA with HSM-protected keys (SEC-005), and layered validation from schema checks through anomaly detection (FR-022--FR-025). The bulk suspension capability (FR-031) and vulnerability window flagging (FR-032) are operational necessities, not edge cases.

### Privacy and science are in tension

Four spatiotemporal data points uniquely identify 95% of individuals. Researchers want high spatial resolution for scientific accuracy. Privacy requires spatial degradation. The specification addresses this through architectural separation of contributor identity from observation data (SEC-004), configurable spatial resolution per campaign, and GDPR-compliant consent architecture (CO-001, CO-002). This tension is not fully resolved --- OI-006 remains open --- but the architecture ensures the tradeoff is explicit and per-campaign, not a platform-wide compromise.

### The cost model is the product thesis

Traditional field data collection has variable cost proportional to data volume. Rootstock inverts this: infrastructure costs are fixed, incentive costs are fixed (sweepstakes prize pool), and the marginal cost per additional reading approaches zero as more scitizens contribute. The implicit budget constraint (Section 2g) requires that the platform's own operational cost not become the cost problem it was designed to solve. Self-hosted infrastructure over managed services (COST-003), hardware HSM over cloud HSM (COST-002), and sweepstakes over per-unit payment (COST-004) all reflect this constraint.

---

## 21d. Constraints That Shape Everything

Three constraints govern the entire system:

1. **CON-001 --- Open Source.** All components, tooling, and documentation are publicly available under an OSI-approved license. No proprietary runtime dependencies. This is not a preference --- it is foundational to institutional trust, community contribution, and auditability.

2. **CON-002 --- First Principles Design.** Every technology choice traces to a requirement. No component is included because "it's what we usually use." This constraint forced the specification to justify every architectural decision against the problem it solves (MT-002).

3. **CON-003 --- No Shared Context Exists.** There is no current system, no shared vocabulary, no agreed-upon data standards, and no existing trust infrastructure between institutions. The platform creates all of these from scratch. Onboarding requires zero changes to an institution's existing internal systems.

---

## 21e. What Remains Open

Six open issues require resolution before or during pilot:

| Issue | Core Question | Blocking |
|-------|--------------|----------|
| OI-001 | Are consumer water/soil sensors accurate enough for publication? | FR-022, FR-006 |
| OI-002 | Can the sweepstakes model comply across jurisdictions? | FR-036 |
| OI-003 | Can the architecture extend to Tier 3 (gateway) devices? | FR-014, FR-015 |
| OI-004 | How much clock drift is tolerable per campaign type? | FR-022, FR-020 |
| OI-005 | What governance model sustains the open-source project? | None directly |
| OI-006 | Where is the spatial resolution / privacy boundary per campaign? | FR-027, SEC-004 |

Two evidence gaps are documented in the risk and cost sections:
- No peer-reviewed co-location study exists for Ecowitt or Davis weather stations (RISK-005, 0x8).
- No OSS-specific GDPR compliance cost data exists (COST-005, 0x16).

Seven assumptions (ASM-001 through ASM-007) carry explicit validation paths. The most consequential is ASM-005 --- that automated IoT devices fundamentally change the retention equation. No direct precedent exists for this model at scale. The pilot must measure device uptime and data continuity to validate or falsify this assumption.

---

## 21f. What This Specification Does Not Cover

This specification defines *what* the product must do and *why*. It does not prescribe *how* --- the architecture, technology stack, and implementation emerge from these requirements per CON-002. Architectural decisions are documented separately in ADRs (MT-002).

The specification also does not cover:
- **Governance model** for the open-source project (OI-005).
- **Tier 3 device architecture** --- gateway-mediated devices are explicitly out of MVP scope (OI-003).
- **Partner integrations** --- no partner applications have been mandated (Section 2c).
- **Schedule** --- no timeline constraints exist (Section 2f).

---

## 21g. Requirement Count Summary

| Section | Category | Count |
|---------|----------|-------|
| 6 | Functional Requirements | 36 (FR-001 through FR-036) |
| 7 | Look and Feel | 3 (LF-001 through LF-003) |
| 8 | Usability | 4 (US-001 through US-004) |
| 9 | Performance | 5 (PE-001 through PE-005) |
| 10 | Operational | 3 (OP-001 through OP-003) |
| 11 | Maintainability | 3 (MT-001 through MT-003) |
| 12 | Security | 6 (SEC-001 through SEC-006) |
| 13 | Cultural | 3 (CU-001 through CU-003) |
| 14 | Compliance | 3 (CO-001 through CO-003) |
| | **Constraints** | 3 (CON-001 through CON-003) |
| | **Open Issues** | 6 (OI-001 through OI-006) |
| | **Risks** | 9 (RISK-001 through RISK-009) |
| | **Costs** | 6 (COST-001 through COST-006) |
| | **Documentation Artifacts** | 12 |
| | **Facts** | 13 (FACT-001 through FACT-013) |
| | **Assumptions** | 7 (ASM-001 through ASM-007) |

---

## 21h. The Core Bet

Rootstock rests on a single structural bet: that the combination of researcher-defined campaigns, automated IoT data collection, and sweepstakes-based incentives can produce publication-grade scientific data at a fraction of traditional field collection costs.

The specification has documented the evidence for this bet (Sections 4--5), defined the requirements to execute it (Sections 6--14), identified what could go wrong (Section 18), quantified what it costs (Section 19), and specified how to explain it to its users (Section 20).

What remains is to build it, pilot it, and find out.

---

*This concludes the Rootstock Requirements Specification.*
