# ROOTSTOCK by Corewood

## Requirements Specification — Section 20: User Documentation and Training

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> This section describes the documentation and training materials that must accompany the product. Rootstock serves two distinct user populations with different needs: researchers (data consumers) and scitizens (data producers). Documentation must serve both without assuming technical expertise from scitizens or dumbing down for researchers.

---

## 20a. Researcher Documentation

### Campaign Management Guide

**Audience:** Researchers creating and managing data campaigns.

**Content:** Step-by-step guide covering campaign creation (FR-005), parameter definition (FR-006), region configuration (FR-007), time window setup (FR-008), publication (FR-009), quality monitoring (FR-010), and data export (FR-026). Includes worked examples for common campaign types (weather monitoring, air quality, water quality).

**Format:** Web-based documentation, searchable, versioned with the platform.

**Fit Criterion:** A researcher with no prior Rootstock experience can create and publish a campaign by following the guide alone, without contacting support. Validated by usability testing (US-002: under 15 minutes).

---

### Data Quality and Provenance Guide

**Audience:** Researchers preparing data for peer review.

**Content:** Explanation of the validation pipeline (FR-022), provenance metadata fields (FR-023), data export formats (FR-026), identity separation (FR-027), and spatial resolution configuration. Describes how Rootstock data aligns with FAIR principles (CO-003) and how to cite Rootstock data in publications.

**Format:** Web-based documentation with downloadable PDF for offline reference.

**Fit Criterion:** A researcher can identify and explain every provenance field in an exported dataset by referring to this guide.

---

### Institution Onboarding Guide

**Audience:** Institution administrators setting up organizational tenants.

**Content:** Organization creation (FR-001), hierarchy configuration (FR-002), role definition (FR-003), researcher invitation (FR-004), and compliance configuration. Emphasizes zero integration requirement (CON-003).

**Format:** Web-based documentation.

**Fit Criterion:** An institution administrator can onboard their organization and invite researchers without contacting Corewood support.

---

## 20b. Scitizen Documentation

### Getting Started Guide

**Audience:** Scitizens registering and enrolling their first device.

**Content:** Account registration (FR-011), campaign browsing (FR-012), enrollment code generation (FR-013), Tier 1 direct enrollment (FR-014), Tier 2 proxy enrollment (FR-015), and first campaign enrollment (FR-017). Written at an 8th-grade reading level (CU-003). Visual step-by-step with screenshots.

**Format:** In-app onboarding flow plus web-based guide. Available in all supported languages (CU-001).

**Fit Criterion:** A scitizen with no prior experience completes registration, device enrollment, and first data submission by following the guide. Validated by usability testing (US-001: under 10 minutes).

---

### Device Troubleshooting Guide

**Audience:** Scitizens experiencing device connectivity or enrollment issues.

**Content:** Common error messages and their meanings (US-004, SEC-006), certificate renewal process (FR-028), grace period behavior (FR-029), device status explanations (pending, active, suspended, revoked), and remediation steps for each failure mode.

**Format:** Web-based searchable knowledge base. Error messages link directly to relevant troubleshooting articles.

**Fit Criterion:** Every user-facing error message has a corresponding troubleshooting article. Scitizens can resolve common issues (expired cert, WiFi change, firmware update) without contacting support.

---

### Contribution and Recognition Guide

**Audience:** Active scitizens wanting to understand their impact.

**Content:** Contribution score calculation (FR-034), badge descriptions and milestones (FR-035), sweepstakes mechanics (FR-036), and how scitizen data contributes to published research. Emphasizes that contributions are scientific participation, not charity (CU-003).

**Format:** In-app profile section plus web-based guide.

**Fit Criterion:** A scitizen can explain what their contribution score means and how to earn badges by reading this guide.

---

## 20c. Operator Documentation

### Deployment Guide

**Audience:** Organizations deploying their own Rootstock instance.

**Content:** Container-based deployment (OP-001), environment configuration, CA setup (SEC-005), OPA policy configuration (SEC-003), MQTT broker integration, database setup, and observability stack configuration.

**Format:** README and documentation in the source repository (CON-001).

**Fit Criterion:** An operator can deploy the full Rootstock stack from the public repository using only the deployment guide and a container runtime. No proprietary tools or undocumented steps.

---

### Security Operations Guide

**Audience:** Platform operators managing device security.

**Content:** Bulk device suspension procedures (FR-031), vulnerability window data flagging (FR-032), device reinstatement (FR-033), certificate lifecycle monitoring (FR-028, FR-029, FR-030), OPA policy management, and HSM key ceremony procedures.

**Format:** Operational runbook in the source repository.

**Fit Criterion:** An operator can execute a bulk device suspension in response to a firmware vulnerability by following the runbook, without ad-hoc decision-making.

---

## 20d. API Documentation

### API Reference

**Audience:** Developers building integrations or extending the platform.

**Content:** All public API endpoints with request/response schemas, authentication requirements, error codes, and rate limits. Generated from source code annotations to stay in sync with the implementation.

**Format:** Auto-generated API documentation (e.g., from protobuf/OpenAPI definitions). Hosted alongside the platform.

**Fit Criterion:** Every public API endpoint is documented with at least one example request and response. Documentation is regenerated on every release.

---

## 20e. Training

### Formal Training

No formal training programs are planned for initial release. The platform must be usable by its target audiences (researchers and scitizens) through documentation and in-app guidance alone. This is validated by the usability requirements:

- **US-001:** Scitizen first-use under 10 minutes
- **US-002:** Campaign creation under 15 minutes
- **US-004:** Actionable error messages eliminate the need for support-mediated troubleshooting

If usability testing reveals that documentation alone is insufficient, training materials will be developed as a follow-on effort.

---

## Documentation Summary

| Document | Audience | Format | Key Requirements |
|----------|----------|--------|-----------------|
| Campaign Management Guide | Researcher | Web | FR-005–FR-010, FR-026 |
| Data Quality Guide | Researcher | Web + PDF | FR-022, FR-023, FR-026, CO-003 |
| Institution Onboarding Guide | Admin | Web | FR-001–FR-004 |
| Getting Started Guide | Scitizen | In-app + Web | FR-011–FR-017, US-001 |
| Device Troubleshooting | Scitizen | Knowledge base | US-004, SEC-006, FR-028–FR-029 |
| Contribution Guide | Scitizen | In-app + Web | FR-034–FR-036 |
| Deployment Guide | Operator | Repository | OP-001, SEC-005 |
| Security Operations | Operator | Runbook | FR-031–FR-033, FR-028–FR-030 |
| API Reference | Developer | Auto-generated | All public endpoints |

---

*Next: [Conclusion](./21_conclusion.md)*
