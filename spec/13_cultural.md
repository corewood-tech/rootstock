# ROOTSTOCK by Corewood

## Requirements Specification — Section 13: Cultural Requirements

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Cultural requirements describe how the product must accommodate different cultural, political, and social contexts.

> **Knowledge graph reference**: Requirements graph at `grapher/schema/rootstock_requirements.graphql`. Start with `GRAPH=reqs podman compose -f grapher/compose-grapher.yml up -d`.

---

### CU-001: Multi-Language Interface Support (0x8e)

**Priority:** Should | **Originator:** Scitizen

**Description:** The platform shall support internationalization (i18n). All user-facing strings shall be externalizable. Initial deployment supports English with architecture ready for additional languages.

**Rationale:** The platform targets global participation. Geographic coverage gaps are in tropical/subtropical regions (FACT-002) where English may not be primary. i18n enables global reach.

**Fit Criterion:** All user-facing strings are stored in locale files, not hardcoded. Adding a new language requires only a new locale file, no code changes. Verified by adding a test locale. (Scale: boolean | Pass/Fail)

**Derived from:** BUC-03, BUC-10 | **Cross-ref:** facts:0x7

---

### CU-002: Multi-Institutional Collaboration (0x8f)

**Priority:** Must | **Originator:** Research Institution

**Description:** The platform shall support multiple research institutions operating independently on the same platform. Each institution maintains its own hierarchy, roles, and governance. Cross-institution campaign participation by scitizens does not require institutional affiliation.

**Rationale:** Institutions have different governance models, compliance requirements, and organizational structures. The platform must support heterogeneous multi-tenant collaboration without requiring institutional standardization.

**Fit Criterion:** Two institutions with different role structures operate independently. A scitizen contributes to campaigns from both without belonging to either. Data isolation between institutions is maintained. (Scale: boolean | Pass/Fail)

**Constrained by:** CON-003 — No Shared Context Exists (0xc)
**Derived from:** BUC-01

---

### CU-003: Inclusive Terminology (0x88)

**Priority:** Should | **Originator:** Corewood

**Description:** The platform shall use inclusive language. The term Scitizen (not volunteer or amateur) reflects active participation in science. Platform copy avoids jargon that excludes non-expert participants.

**Rationale:** Scitizen as a term reflects participant agency (Section 3a). Platform language should be inclusive and avoid gatekeeping scientific participation.

**Fit Criterion:** No user-facing text uses the terms "volunteer" or "amateur" to describe scitizens. All onboarding text tested for readability at an 8th-grade reading level. (Scale: boolean | Pass/Fail)

**Derived from:** BUC-03, BUC-10

---

*Next: [Section 14 — Compliance Requirements](./14_compliance.md)*
