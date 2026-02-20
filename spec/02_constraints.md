# ROOTSTOCK by Corewood

## Requirements Specification — Section 2: Constraints

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Constraints are global — they apply to the entire product and restrict how the problem can be solved. Each constraint carries a description, rationale, and fit criterion. The Volere template warns against false constraints: only include what is genuinely non-negotiable.

---

## 2a. Solution Constraints

### CON-001: Open Source

**Description:** The product must be fully open source. All platform components, tooling, and documentation must be publicly available under an OSI-approved license.

**Rationale:** Rootstock exists to enable citizen science. Trust, transparency, and community contribution are foundational to that mission. Proprietary components would undermine adoption by research institutions who need to audit what touches their data, and would prevent the community contribution model that scales the platform beyond Corewood's capacity.

**Fit Criterion:** All source code, configuration, infrastructure-as-code, and documentation are published in a public repository under an OSI-approved license. No runtime dependency requires a proprietary license to operate.

---

### CON-002: First Principles Design

**Description:** The product must be designed from first principles. There are no preconceived technology mandates, framework requirements, or inherited architectural decisions. Solution choices are driven by the problem, not by prior preference.

**Rationale:** The problem space — high-volume IoT data ingestion, multi-tenant hierarchical authorization, campaign-driven data collection — is specific enough that off-the-shelf assumptions are likely to be wrong. Premature technology commitments create false constraints that compound over time. The architecture must emerge from the requirements, not precede them.

**Fit Criterion:** Every technology choice in the architecture is traceable to a specific requirement or constraint in this document. No component is included because "it's what we usually use."

---

### CON-004: ULID Identifiers

**Description:** All platform-generated identifiers shall be ULIDs (Universally Unique Lexicographically Sortable Identifiers). External identifiers (e.g., Zitadel user IDs, certificate serials) are stored as-is but referenced via foreign key to a ULID primary key.

**Rationale:** ULIDs combine the uniqueness of UUIDs with lexicographic sortability (monotonic within the same millisecond), making them suitable for primary keys in time-ordered tables without index fragmentation. They are 128-bit compatible with UUID storage but encode as 26-character Crockford Base32 strings, making them human-readable in logs and URLs. Application-generated IDs (vs. database-generated `gen_random_uuid()`) keep ID creation in the ops layer where business logic lives, consistent with the clean architecture boundary.

**Fit Criterion:** All primary keys for platform-created entities are ULIDs generated in application code (`github.com/oklog/ulid/v2`). No table uses `gen_random_uuid()` for primary keys. External identifiers (IdP user IDs, certificate serials) are stored as TEXT in reference columns, not as primary keys. (Scale: boolean | Pass/Fail)

---

## 2b. Implementation Environment of the Current System

### CON-003: No Shared Context Exists

**Description:** There is no current system. The work that Rootstock aims to support — researchers requesting field data and citizen scientists providing it — currently happens across disconnected groups of institutions with no shared platform, protocol, or data format. Each institution operates independently with its own tools, processes, and data silos.

**Rationale:** This is a greenfield product, but the absence of a shared system is itself a constraint. It means there is no existing data to migrate, no incumbent workflow to preserve, and no integration contract to honor — but it also means there is no shared vocabulary, no agreed-upon data standards, and no existing trust infrastructure between institutions. The product must create all of these from scratch.

**Fit Criterion:** The product does not assume any pre-existing shared infrastructure, identity system, or data format across participating institutions. Onboarding a new institution requires zero changes to that institution's existing internal systems.

---

## 2c. Partner or Collaborative Applications

> To be determined. No partner application integrations have been mandated at this time. Candidates will likely emerge as the product scope is defined — for example, institutional identity providers (SSO/SAML), existing research data repositories, or IoT device management platforms. These should be captured here as they become known constraints rather than assumed prematurely.

---

## 2d. Off-the-Shelf Software

> No off-the-shelf software has been mandated. Per CON-002, technology choices will be evaluated against requirements as they are defined. This section will be updated if and when specific OTS components are selected and become constraints on the rest of the system.

---

## 2e. Anticipated Workplace Environment

> To be determined. The product serves two distinct operational environments that will need to be characterized:
>
> - **Researcher environment:** Likely institutional (university, lab, office). Web-based access assumed but not yet mandated.
> - **Citizen scientist environment:** Field conditions. Mobile connectivity, varying device capabilities, intermittent network access. This environment is inherently unpredictable and will generate significant operational and design constraints once characterized.
> - **IoT device environment:** See Section 1.4 (Non-Human Actors). Device volatility — hardware fragmentation, firmware instability, connectivity unreliability — constitutes an environmental constraint on the ingestion layer.

---

## 2f. Schedule Constraints

> No schedule constraints have been mandated at this time.

---

## 2g. Budget Constraints

> No explicit budget constraints have been mandated. However, the project goals (Section 1.2) establish that the platform must drive data acquisition costs down significantly. This creates an implicit constraint: the platform's own operational cost must remain low enough that it does not become the cost problem it was designed to solve.

---

## 2h. Enterprise Constraints

> No enterprise constraints have been mandated. As an open-source project built by Corewood, there are no inherited organizational policies, mandated development processes, or corporate technology standards that constrain the solution. Governance model for the open-source project itself is not yet defined and may introduce constraints when established.

---

*Next: [Section 3 — Naming Conventions and Terminology](./03-naming-conventions.md)*