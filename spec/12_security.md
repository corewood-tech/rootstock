# ROOTSTOCK by Corewood

## Requirements Specification — Section 12: Security Requirements

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Security requirements describe how the product must protect itself, its data, and its users from unauthorized access, modification, or disclosure.

> **Knowledge graph reference**: Requirements graph at `grapher/schema/rootstock_requirements.graphql`. Start with `GRAPH=reqs podman compose -f grapher/compose-grapher.yml up -d`.

---

### SEC-001: mTLS for All Device Communication (0x83)

**Priority:** Must | **Originator:** Researcher

**Description:** All device-to-platform communication (MQTT and HTTP/2) shall use mutual TLS. Both the device and the platform present certificates. Connections without valid client certificates are rejected at the TLS layer.

**Rationale:** Every device connection must be mutually authenticated. Without mTLS, data provenance is unverifiable and spoofing is trivial.

**Fit Criterion:** No device connection succeeds without a valid client certificate signed by the platform CA. Tested with: no cert, expired cert, self-signed cert, wrong CA cert. All rejected. (Scale: boolean | Pass/Fail)

**Derived from:** BUC-06, BUC-08 | **Cross-ref:** facts:0x5

---

### SEC-002: Private Key Never Leaves Device (0x84)

**Priority:** Must | **Originator:** Researcher

**Description:** During enrollment, the device shall generate its own keypair locally. Only the CSR (containing the public key) is transmitted to the enrollment service. The private key is never transmitted over any network.

**Rationale:** If private keys leave the device, the entire identity model is compromised. The device generates its own keypair; only the CSR (containing the public key) is transmitted.

**Fit Criterion:** Network capture during enrollment (direct and proxy) shows no private key material. CSR contains public key only. Companion app never receives private key. (Scale: boolean | Pass/Fail)

**Derived from:** BUC-04, BUC-08

---

### SEC-003: RBAC via OPA Policies (0x7e)

**Priority:** Must | **Originator:** Research Institution

**Description:** The platform shall enforce role-based access control where all authorization decisions are made by OPA. Policies are expressed in Rego and are testable and auditable. Roles are scoped to organizations.

**Rationale:** Permissions are granted via roles, never directly. Roles are organization-scoped. All authorization decisions are made by OPA evaluating Rego policies.

**Fit Criterion:** A user without campaign_create permission is denied campaign creation. A user with the role in org A cannot create campaigns in org B. All policies have Rego unit tests with >90% coverage. (Scale: percentage | Worst: 80 | Plan: 90 | Best: 100)

**Derived from:** BUC-01, BUC-06

---

### SEC-004: Data Privacy and Identity Separation (0x80)

**Priority:** Must | **Originator:** Oversight Body

**Description:** The platform shall architecturally separate contributor identity from observation data. Consent is granular (per-campaign, per-use), versioned, and auditable. Spatial resolution is configurable per campaign. Raw contributor locations never appear in public-facing datasets.

**Rationale:** 4 spatiotemporal data points uniquely identify 95% of individuals (FACT-009). GDPR treats location data as PII. Contributor identity must be architecturally separable from observation data.

**Fit Criterion:** Database schema physically separates identity tables from observation tables. Join requires explicit authorization. Exported datasets contain zero PII. Consent records are versioned and auditable. (Scale: boolean | Pass/Fail)

**Derived from:** BUC-07 | **Cross-ref:** facts:0x30

---

### SEC-005: HSM Protection for CA Keys (0x81)

**Priority:** Must | **Originator:** Researcher

**Description:** Issuing CA private keys shall be stored in an HSM (SoftHSM for development, YubiHSM 2 or cloud HSM for production). Root CA is air-gapped or HSM-protected. Access via PKCS#11 interface.

**Rationale:** CA private keys are the root of trust for the entire device identity system. HSM protection is required for issuing CA keys — not optional at scale.

**Fit Criterion:** Issuing CA private key is not extractable from the HSM. Certificate signing operations use PKCS#11 API. Root CA key is stored offline. (Scale: boolean | Pass/Fail)

**Derived from:** BUC-08

---

### SEC-006: Connection Failure Diagnostics (0x7a)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall produce structured log entries on every failed device connection attempt. Logs include device IP, presented cert, failure reason, and timestamp. Failure reasons are specific and surfaced to scitizens as human-readable messages.

**Rationale:** Every failed connection must produce structured, actionable diagnostics. Specific failure reasons enable self-service troubleshooting.

**Fit Criterion:** Every connection failure produces a structured log with: device IP, cert info (if any), specific reason (expired, suspended, not enrolled, wrong CA), and timestamp. Reasons are human-readable. (Scale: boolean | Pass/Fail)

**Derived from:** BUC-06, BUC-04

---

*Next: [Section 13 — Cultural Requirements](./13_cultural.md)*
