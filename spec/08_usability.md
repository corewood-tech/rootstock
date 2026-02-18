# ROOTSTOCK by Corewood

## Requirements Specification — Section 8: Usability and Humanity Requirements

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Usability requirements describe how easy the product must be to use, how quickly users must be able to learn it, and how it must accommodate different user capabilities.

> **Knowledge graph reference**: Requirements graph at `grapher/schema/rootstock_requirements.graphql`. Start with `GRAPH=reqs podman compose -f grapher/compose-grapher.yml up -d`.

---

### US-001: First-Use Device Enrollment Under 10 Minutes (0x5d)

**Priority:** Must | **Originator:** Scitizen

**Description:** A scitizen with no prior platform experience shall be able to register, enroll a Tier 1 device, and submit a first data reading within 10 minutes of landing on the platform.

**Rationale:** The first contribution is the retention inflection point — 82% retention post-first-contribution vs 39.7% without (FACT-006). Enrollment friction is the highest-leverage design problem.

**Fit Criterion:** 5 test users with no prior experience complete registration, Tier 1 device enrollment, and first reading submission. Median time is under 10 minutes. No user takes more than 15 minutes. (Scale: minutes | Worst: 15 | Plan: 10 | Best: 5)

**Derived from:** BUC-03, BUC-04 | **Cross-ref:** facts:0x33

---

### US-002: Campaign Creation Under 15 Minutes (0x5f)

**Priority:** Must | **Originator:** Researcher

**Description:** An authenticated researcher shall be able to create and publish a campaign with parameters, region, time window, and quality thresholds within 15 minutes.

**Rationale:** If campaign creation takes >30 minutes, researchers will not adopt the platform (ASM-003). Must be significantly faster than preparing a traditional data collection protocol.

**Fit Criterion:** 5 test researchers with defined data needs complete campaign creation and publication. Median time under 15 minutes. No researcher takes more than 30 minutes. (Scale: minutes | Worst: 30 | Plan: 15 | Best: 8)

**Derived from:** BUC-02

---

### US-003: Clear Feedback on Data Contribution (0x61)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall provide scitizens with clear, real-time feedback on their contributions: accepted readings count, active campaigns, contribution score, and recent badges. Device status and data flow indicators must be visible.

**Rationale:** Scitizens need clear feedback that their data matters. Without visible contribution value, engagement drops.

**Fit Criterion:** The scitizen dashboard shows accepted reading count, active campaign count, contribution score, and recent badges. Dashboard loads within 3 seconds. (Scale: seconds | Worst: 5 | Plan: 3 | Best: 1)

**Derived from:** BUC-10, BUC-03 | **Cross-ref:** facts:0x13

---

### US-004: Actionable Error Messages (0x63)

**Priority:** Must | **Originator:** Scitizen

**Description:** All error messages presented to users shall be specific and actionable. Device connection failures include reason and remediation steps. Enrollment failures include what to fix. Validation rejections include which parameter failed and why.

**Rationale:** Connection failures are surfaced as human-readable messages (Section 3a). Generic errors cause support burden and churn.

**Fit Criterion:** Every user-facing error message includes: what failed, why it failed, and what the user can do. No error message is a generic 500 or unformatted stack trace. (Scale: boolean | Pass/Fail)

**Derived from:** BUC-04, BUC-06

---

*Next: [Section 9 — Performance Requirements](./09_performance.md)*
