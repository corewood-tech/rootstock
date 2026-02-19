# ROOTSTOCK by Corewood

## Requirements Specification — Section 19: Costs

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> This section identifies the cost factors that the project must account for. Rootstock is an open-source reference project — there is no commercial budget to specify. Instead, this section documents the cost categories that affect architectural decisions, operational sustainability, and the platform's core value proposition.

---

## 19a. Development Costs

### Open Source Development Model

Rootstock is built by Corewood as an open-source project (CON-001). No schedule or budget constraints have been mandated (Sections 2f, 2g). Development costs are borne by Corewood with the expectation that community contribution scales the effort beyond Corewood's capacity.

**Cost driver:** Every technology choice must be traceable to a requirement (CON-002, MT-002). This prevents gratuitous complexity and keeps the dependency surface small, which directly controls development and maintenance cost.

**Cost constraint:** No runtime dependency may require a proprietary license (MT-001). This eliminates licensing costs for the platform and for every institution that deploys it.

---

## 19b. Infrastructure Costs

### Self-Operated Certificate Authority

The Rootstock CA is self-operated because external CAs (Let's Encrypt, AWS Private CA) either do not issue client certificates at IoT scale or create cost/proprietary dependencies. Self-operation shifts the cost from per-certificate fees to HSM hardware and operational overhead.

| Component | Development | Production |
|-----------|------------|------------|
| Root CA key storage | SoftHSM (free) | Air-gapped HSM or YubiHSM 2 (~$650) |
| Issuing CA key storage | SoftHSM (free) | YubiHSM 2 (~$650) or cloud HSM |
| CA software (step-ca) | Open source | Open source |

**Cost implication for SEC-005:** HSM hardware is a fixed cost. Per-certificate cost is near zero after initial setup.

### Compute and Storage

The platform must sustain 10,000 readings/second (PE-001) with 99.9% availability (PE-004). Infrastructure costs scale with device count and data volume, not with user count. Key scaling dimensions:

| Dimension | Cost Driver | Requirement |
|-----------|------------|-------------|
| Ingestion throughput | MQTT broker + validation pipeline compute | PE-001, PE-002 |
| Storage | Time-series data volume per campaign | FR-023 |
| OPA policy evaluation | Memory for registry bundle at scale | PE-005 |
| Certificate operations | CA signing throughput | FR-014, FR-028 |

**Cost constraint from Section 1.2:** The platform's own operational cost must remain low enough that it does not become the cost problem it was designed to solve. If the platform costs more to run than the field data collection it replaces, the value proposition fails.

---

## 19c. Scitizen Incentive Costs

### Sweepstakes Model

The cost model in Section 1.2 uses sweepstakes entries, experiences, and recognition rather than per-unit compensation. This is the primary cost-control mechanism for the incentive model.

**Why this works (FACT-008):** Prospect Theory demonstrates that people overweight small probabilities, making lottery-style incentives psychologically more motivating per dollar spent than fixed per-unit payments. Variable-ratio reinforcement produces the strongest sustained engagement patterns.

**Cost scaling:** As more scitizens join, per-unit data cost approaches zero while coverage expands. The incentive budget is fixed (sweepstakes prize pool), not proportional to data volume. This is the fundamental economic advantage over traditional field data collection (FACT-001: personnel costs consume 50–80% of grant budgets).

**Legal cost caveat (OI-002):** Sweepstakes administration may require legal review and potentially a third-party administration service, adding operational cost.

---

## 19d. Compliance Costs

### GDPR and IRB

GDPR compliance (CO-001) and IRB compatibility (CO-002) are architectural requirements, not add-on services. The cost is in design and implementation, not in ongoing fees. Key cost items:

- **Privacy architecture:** Architectural separation of identity from observation data (SEC-004). This is a design cost, not an operational cost.
- **Consent management:** Versioned, auditable consent records per campaign. Storage cost scales with enrollment count.
- **Right to erasure:** Data deletion pipeline must be implemented and tested. Operational cost per erasure request.
- **Audit trail:** Immutable consent and authorization logs. Storage cost scales with platform activity.

**Cost avoidance:** The open-source, auditable nature of the platform (CON-001) supports institutional trust and reduces the compliance burden for each adopting institution. Without this, each institution would need to audit a proprietary black box.

---

## 19e. Cost Summary

| Category | Nature | Scaling |
|----------|--------|---------|
| Development | Corewood + community | Fixed effort, community scales |
| HSM hardware | One-time | Fixed per deployment |
| Compute/storage | Ongoing | Scales with device count and data volume |
| Incentives | Ongoing | Fixed budget, not proportional to data |
| Compliance | Design-time | One-time architecture cost |
| Sweepstakes admin | Ongoing | Per-jurisdiction legal review |

**Key insight:** The platform's cost model inverts the traditional data collection cost structure. Traditional field collection has variable cost proportional to data volume (more data = more personnel). Rootstock has fixed infrastructure cost with marginal cost approaching zero per additional reading.

---

*Next: [Section 20 — User Documentation and Training](./20_user_documentation.md)*
