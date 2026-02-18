# ROOTSTOCK by Corewood

## Requirements Specification — Section 10: Operational and Environmental Requirements

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Operational requirements describe the environment in which the product will operate, including partner systems, deployment constraints, and expected operating conditions.

> **Knowledge graph reference**: Requirements graph at `grapher/schema/rootstock_requirements.graphql`. Start with `GRAPH=reqs podman compose -f grapher/compose-grapher.yml up -d`.

---

### OP-001: Container-Based Deployment (0x72)

**Priority:** Must | **Originator:** Corewood

**Description:** All platform components shall be deployable as OCI-compliant containers. A compose file shall bring up the full stack for local development.

**Rationale:** Open-source platform must be deployable by any organization. Container-based deployment ensures reproducibility and portability.

**Fit Criterion:** The full platform stack starts with a single compose command. All containers pass health checks within 60 seconds. No host-specific dependencies beyond a container runtime. (Scale: seconds | Worst: 120 | Plan: 60 | Best: 30)

**Constrained by:** CON-001 — Open Source (0xd)
**Derived from:** BUC-01

---

### OP-002: Dual Protocol Ingestion (0x6e)

**Priority:** Must | **Originator:** Scitizen

**Description:** The platform shall accept data ingestion via both MQTT 5.0 and HTTP/2 Gateway. Both paths use mTLS and feed into the same validation and persistence pipeline.

**Rationale:** Different device tiers suit different protocols. MQTT for constrained devices; HTTP/2 for smartphones and Raspberry Pis. Both must be supported.

**Fit Criterion:** The same reading submitted via MQTT and HTTP/2 produces identical persisted results. Both paths enforce mTLS and OPA authorization. (Scale: boolean | Pass/Fail)

**Derived from:** BUC-06, BUC-04

---

### OP-003: Graceful Degradation on Device Failure (0x70)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall detect device disconnections via MQTT LWT and update device status accordingly. Reconnecting devices resume data submission without re-enrollment. Campaign data continues flowing from other enrolled devices.

**Rationale:** IoT devices are inherently unreliable (intermittent connectivity, power loss, environmental interference). The platform must degrade gracefully.

**Fit Criterion:** A device that disconnects and reconnects resumes data submission without re-enrollment or manual intervention. LWT triggers status update within 60 seconds of unexpected disconnect. (Scale: seconds | Worst: 120 | Plan: 60 | Best: 10)

**Derived from:** BUC-06, BUC-08

---

*Next: [Section 11 — Maintainability and Support Requirements](./11_maintainability.md)*
