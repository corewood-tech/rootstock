# ROOTSTOCK by Corewood

## Requirements Specification — Section 11: Maintainability and Support Requirements

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Maintainability requirements describe how easy the product must be to change, extend, and support over its lifetime.

> **Knowledge graph reference**: Requirements graph at `grapher/schema/rootstock_requirements.graphql`. Start with `GRAPH=reqs podman compose -f grapher/compose-grapher.yml up -d`.

---

### MT-001: Open Source Licensing (0x74)

**Priority:** Must | **Originator:** Corewood

**Description:** All source code, configuration, infrastructure-as-code, and documentation shall be published under an OSI-approved license. No runtime dependency requires a proprietary license.

**Rationale:** CON-001 mandates full open source. All code, config, IaC, and docs must be publicly available under OSI-approved license.

**Fit Criterion:** License audit of all runtime dependencies shows zero proprietary licenses. Project license is OSI-approved. All source repos are public. (Scale: boolean | Pass/Fail)

**Constrained by:** CON-001 — Open Source (0xd)
**Derived from:** BUC-01

---

### MT-002: Technology Traceability (0x76)

**Priority:** Must | **Originator:** Corewood

**Description:** Every technology component in the architecture shall be traceable to a specific requirement or constraint in this specification. An Architecture Decision Record (ADR) documents the rationale for each technology choice.

**Rationale:** CON-002 mandates first-principles design. Every technology choice must trace to a requirement — no component included because "it's what we usually use."

**Fit Criterion:** Every technology component has a corresponding ADR referencing requirements it satisfies. No component exists without a traced requirement. (Scale: boolean | Pass/Fail)

**Constrained by:** CON-002 — First Principles Design (0xb)
**Derived from:** BUC-01

---

### MT-003: Modular Architecture (0x78)

**Priority:** Must | **Originator:** Corewood

**Description:** The platform shall be structured with clean architectural boundaries. Dependencies point inward. Auth, observability, transport, business logic, and data access are separate components. Changing one does not require changing others.

**Rationale:** Volatility decomposition: components that change independently (auth, o11y, transport, business logic, data access) must be separated to prevent ripple effects.

**Fit Criterion:** Replacing the data access implementation (e.g., PostgreSQL to another store) requires changes only in the repo layer. No handler, flow, or op code changes. Verified by code review. (Scale: layers affected | Worst: 2 | Plan: 1 | Best: 1)

**Constrained by:** CON-002 — First Principles Design (0xb)
**Derived from:** BUC-01

---

*Next: [Section 12 — Security Requirements](./12_security.md)*
