# ROOTSTOCK by Corewood

## Requirements Specification — Section 9: Performance Requirements

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Performance requirements describe the speed, capacity, reliability, and scalability the product must achieve. Each has measurable fit criteria.

> **Knowledge graph reference**: Requirements graph at `grapher/schema/rootstock_requirements.graphql`. Start with `GRAPH=reqs podman compose -f grapher/compose-grapher.yml up -d`.

---

### PE-001: Data Ingestion Throughput (0x6b)

**Priority:** Must | **Originator:** Researcher

**Description:** The ingestion pipeline shall sustain a minimum of 10,000 validated readings per second under steady-state load.

**Rationale:** A network of 250K+ devices submitting readings at varying intervals requires sustained ingestion capacity. Peak load during campaign window opens.

**Fit Criterion:** Load test sustaining 10,000 readings/second for 10 minutes with zero data loss. All readings validated and persisted. (Scale: readings/second | Worst: 5,000 | Plan: 10,000 | Best: 50,000)

**Derived from:** BUC-06

---

### PE-002: Ingestion Latency (0x65)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform shall process a single reading (receive, validate, persist) in under 200ms at the 99th percentile under steady-state load.

**Rationale:** Low ingestion latency ensures near-real-time data availability for campaigns. Excessive latency degrades researcher experience and delays quality feedback.

**Fit Criterion:** P99 ingestion latency (receive to persist) is under 200ms during sustained 10K readings/second load. P50 is under 50ms. (Scale: milliseconds | Worst: 500 | Plan: 200 | Best: 50)

**Derived from:** BUC-06

---

### PE-003: API Response Time (0x67)

**Priority:** Must | **Originator:** Researcher

**Description:** All human-facing API endpoints shall respond within 500ms at the 95th percentile under normal load.

**Rationale:** Human-facing API calls (campaign listing, dashboard queries, device management) must be responsive for usability.

**Fit Criterion:** P95 API response time for human-facing endpoints is under 500ms. P50 is under 200ms. Measured under 100 concurrent users. (Scale: milliseconds | Worst: 1,000 | Plan: 500 | Best: 100)

**Derived from:** BUC-02, BUC-03

---

### PE-004: Platform Availability (0x6c)

**Priority:** Must | **Originator:** Researcher

**Description:** The platform ingestion pipeline shall maintain 99.9% availability measured monthly. Planned maintenance windows are excluded if announced 48 hours in advance.

**Rationale:** IoT devices operate continuously. Downtime causes data loss during campaign windows that cannot be recovered.

**Fit Criterion:** Monthly uptime is at least 99.9% (maximum 43.8 minutes downtime per month). Uptime measured as ingestion endpoint responding to health checks. (Scale: percentage | Worst: 99.5 | Plan: 99.9 | Best: 99.99)

**Derived from:** BUC-06

---

### PE-005: OPA Authorization Latency (0x6a)

**Priority:** Must | **Originator:** Researcher

**Description:** OPA authorization decisions shall complete in under 5ms at the 99th percentile. Bundle refresh interval shall not exceed 30 seconds.

**Rationale:** OPA is in the hot path for every device action. Authorization latency directly impacts ingestion throughput.

**Fit Criterion:** P99 OPA decision latency under 5ms with 100K device registry entries loaded. Bundle refresh completes in under 30 seconds. (Scale: milliseconds | Worst: 10 | Plan: 5 | Best: 1)

**Derived from:** BUC-06, BUC-09

---

*Next: [Section 10 — Operational and Environmental Requirements](./10_operational.md)*
