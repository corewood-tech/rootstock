# ROOTSTOCK by Corewood

## Requirements Specification — Section 7: Look and Feel Requirements

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Look and feel requirements describe the intended appearance and style of the product. They ensure the product conveys the right impression to its users.

> **Knowledge graph reference**: Requirements graph at `grapher/schema/rootstock_requirements.graphql`. Start with `GRAPH=reqs podman compose -f grapher/compose-grapher.yml up -d`.

---

### LF-001: Consistent Visual Identity (0x56)

**Priority:** Should | **Originator:** Scitizen

**Description:** The platform shall present a consistent visual identity across all interfaces: campaign browsing, device management, data dashboards, and profile pages. Typography, color palette, and layout patterns are uniform.

**Rationale:** A consistent visual identity builds trust and communicates that the platform is professionally maintained. Scitizens and researchers need confidence this is a credible scientific tool.

**Fit Criterion:** All pages use the same design system. A visual audit of the 5 primary user flows identifies zero deviations from the established component library. (Scale: deviations | Worst: 3 | Plan: 0 | Best: 0)

**Derived from:** BUC-03, BUC-02

---

### LF-002: Scientific Credibility in Presentation (0x58)

**Priority:** Should | **Originator:** Researcher

**Description:** Data-facing interfaces (campaign dashboards, export previews, quality metrics) shall present data with scientific conventions: SI units, appropriate significant figures, proper axis labels, and clear provenance indicators.

**Rationale:** Researchers at institutions will not adopt a platform that looks like a consumer app or a gamification toy. Data presentation must convey precision and rigor.

**Fit Criterion:** All data values display with correct SI units. Numeric precision matches campaign-defined significant figures. Charts include labeled axes with units. (Scale: boolean | Pass/Fail)

**Derived from:** BUC-02, BUC-07

---

### LF-003: Responsive Design (0x5a)

**Priority:** Must | **Originator:** Scitizen

**Description:** All platform interfaces shall render correctly on screens from 320px to 2560px width. Critical scitizen flows (registration, device enrollment, campaign browse) must be fully functional on mobile.

**Rationale:** Scitizens operate in the field. Device enrollment and campaign browsing must work on mobile devices. Researchers may use desktop or tablet.

**Fit Criterion:** All critical user flows complete without horizontal scrolling on a 320px viewport. No UI elements are clipped or overlapping on any tested viewport width (320, 768, 1024, 1440, 2560). (Scale: viewport widths | Worst: 3 | Plan: 5 | Best: 5)

**Derived from:** BUC-03, BUC-04

---

*Next: [Section 8 — Usability and Humanity Requirements](./08_usability.md)*
