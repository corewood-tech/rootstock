# ROOTSTOCK by Corewood

## Requirements Specification — Section 14: Compliance Requirements

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Compliance requirements describe the legal, regulatory, and standards obligations the product must satisfy.

> **Knowledge graph reference**: Requirements graph at `grapher/schema/rootstock_requirements.graphql`. Start with `GRAPH=reqs podman compose -f grapher/compose-grapher.yml up -d`.

---

### CO-001: GDPR Compliance Architecture (0x89)

**Priority:** Must | **Originator:** Oversight Body

**Description:** The platform shall implement GDPR-compliant data handling: right to erasure, granular consent (per-campaign, per-use), data portability, and privacy by design. Location data is treated as PII throughout the architecture.

**Rationale:** GDPR treats location data as PII and applies to all persons in the EEA. The platform collects geolocation data from devices owned by individuals. GDPR compliance is non-negotiable for institutional adoption.

**Fit Criterion:** A scitizen can request erasure of all their data. The request is fulfilled within 30 days. Consent records are versioned and auditable. Data portability export is available in machine-readable format. (Scale: days | Worst: 30 | Plan: 14 | Best: 7)

**Derived from:** BUC-07, BUC-01 | **Cross-ref:** facts:0x30

---

### CO-002: IRB-Compatible Consent Model (0x8b)

**Priority:** Must | **Originator:** Oversight Body

**Description:** The platform shall implement an informed consent model that satisfies IRB requirements. Consent is obtained per-campaign, records are versioned, and consent status is queryable by institution for audit purposes.

**Rationale:** US and international IRBs require informed consent for research involving human participants. GDPR-compliant consent satisfies most IRB requirements, but the model must be auditable.

**Fit Criterion:** Every campaign enrollment records explicit consent with timestamp, version, and scope. Consent records are exportable for IRB audit. Consent withdrawal stops further data collection from the device for that campaign. (Scale: boolean | Pass/Fail)

**Derived from:** BUC-01 | **Cross-ref:** facts:0x31

---

### CO-003: Data Provenance for Publication (0x91)

**Priority:** Should | **Originator:** Researcher

**Description:** Exported campaign data shall include provenance metadata sufficient for peer review: collection methodology, quality thresholds applied, device calibration status, validation pipeline version, and consent scope. Data format supports FAIR principles.

**Rationale:** FAIR principles and TRUST principles define scientific data management standards. Data must be findable, accessible, interoperable, and reusable to be publication-grade.

**Fit Criterion:** Exported data includes provenance metadata sufficient for a reviewer to trace any reading back to its source device, collection time, quality checks applied, and consent scope. Metadata follows W3C PROV or equivalent. (Scale: boolean | Pass/Fail)

**Derived from:** BUC-07 | **Cross-ref:** facts:0x34

---

*Next: [Section 15 — Pricing](./15_pricing.md)*
