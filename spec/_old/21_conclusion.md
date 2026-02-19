# ROOTSTOCK by Corewood

## Requirements Specification — Conclusion

---

## What This Specification Defines

This document specifies Rootstock — an open-source scientific data collection platform that connects researchers with scitizens through consumer IoT devices. Twenty-one sections trace the project from motivation through requirements to risks and documentation.

**Sections 1–4** establish why Rootstock exists. Traditional field data collection is expensive and coverage-limited (FACT-001, FACT-002). Consumer IoT devices are already deployed at massive scale but disconnected from research needs (FACT-003, FACT-004). No platform bridges this gap (FACT-004). The economic model works: citizen science produces $0.7–2.5B in annual value with 46–77% cost reduction over contracted collection (FACT-005).

**Section 5** defines the scope: 9 adjacent systems, 12 business events, and 10 business use cases covering the full lifecycle from institutional onboarding through data export and scitizen recognition.

**Sections 6–14** specify 66 requirements — 36 functional and 30 non-functional — each with a testable fit criterion. 52 are Must priority, 14 are Should. Every requirement traces to at least one BUC. Every requirement has a fit criterion.

**Sections 15 and 17** are explicitly not applicable. This is a reference project, not a commercial product or an evaluation of alternatives.

**Section 16** documents 6 open issues that must be resolved before the affected requirements can be fully implemented.

**Sections 18–20** address risks (9 identified, each with mitigation), costs (inverted cost model where marginal data cost approaches zero), and the documentation needed to make the platform self-service for both researchers and scitizens.

---

## Three Constraints Shape Everything

Every architectural and design decision in Rootstock flows from three constraints:

**CON-001 — Open Source.** All code, config, and documentation are public under an OSI-approved license. No proprietary runtime dependencies. This is not just a licensing choice — it is the trust foundation. Research institutions need to audit what touches their data. Community contribution scales the project beyond Corewood's capacity.

**CON-002 — First Principles.** No preconceived technology mandates. Every component is traceable to a requirement. This prevents the accumulated complexity that kills open-source projects — if a dependency cannot justify its existence against a specific requirement, it does not belong.

**CON-003 — No Shared Context.** There is no current system, no shared vocabulary, no agreed-upon data standards, and no existing trust infrastructure between institutions. Rootstock creates all of these from scratch. Onboarding a new institution requires zero changes to that institution's existing systems.

---

## The Core Bet

Rootstock makes one fundamental bet: **that the gap between research data needs and consumer IoT device capability is a coordination problem, not a technology problem.**

The devices exist. 250,000+ weather stations, 30,000+ air quality monitors, and a growing number of water and soil sensors are already deployed, powered, connected, and generating data. The research demand exists — field data collection consumes 50–80% of grant budgets and covers less than 7% of the Earth's surface.

The missing piece is a platform that lets researchers define exactly what data they need and matches scitizens with the right devices to that need. Rootstock is that platform.

The campaign model inverts the typical citizen science flow: instead of "collect data and hope researchers use it" (the model that produces zero publications 75% of the time), Rootstock starts with "researcher defines what they need" and matches contributors to that need. Every data reading enters the platform with a purpose.

---

## What Comes Next

This specification is the input to architecture and implementation. The requirements graph (`grapher/schema/rootstock_requirements.graphql`) contains all 151 nodes — requirements, fit criteria, BUC references, constraints, and open issues — as a queryable knowledge base for the implementation phase.

The implementation order follows the BUC dependency chain:

```
BUC-01 Institutional Onboarding
  → BUC-02 Campaign Creation
    → BUC-05 Campaign Enrollment ← BUC-04 Device Registration
      → BUC-06 Data Ingestion
        → BUC-07 Data Export
        → BUC-10 Scitizen Recognition

BUC-03 Scitizen Registration
  → BUC-04 Device Registration
    → BUC-08 Certificate Lifecycle

BUC-09 Device Security Response (independent)
```

Three open issues must be resolved early:
1. **OI-001** — Consumer sensor accuracy validation (co-location studies for water/soil)
2. **OI-004** — Clock drift tolerance strategy
3. **OI-006** — Spatial resolution privacy tradeoff

Everything else is specified, traced, and testable.

---

*Rootstock — because the best science grows from strong roots and citizen ground truth.*
