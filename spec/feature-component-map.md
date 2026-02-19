# ROOTSTOCK — Feature → Component → Handling Map

> Derived from functional requirements (Section 6), business use cases (Section 5), and non-functional constraints (Sections 7–12). Architecture emerges from the requirements, not the other way around (CON-002).

---

## Components (derived from requirement volatility)

Each component exists because it changes independently from the others. If two things always change together, they belong in the same component. If changing one shouldn't force changes in the other, they're separate.

| Component | Why it exists (volatility axis) | Changes when... |
|-----------|--------------------------------|-----------------|
| **UI** | Presentation changes independently from business logic | Look-and-feel, layout, user flows change; business rules don't |
| **RPC Handlers** | Protocol changes independently from business logic | Transport format, auth middleware, request shape change; domain logic doesn't |
| **Identity (Zitadel)** | Human authn changes independently from everything else | Auth provider, login flows, social login options change; nothing else does |
| **Auth (OPA)** | Authorization policy changes independently from code | Roles, permissions, ACLs change via Rego; no code deploys needed |
| **Campaign Service** | Campaign rules change independently from device management | Campaign parameters, regions, windows evolve; device handling doesn't care |
| **Device Registry** | Device lifecycle changes independently from campaigns | Enrollment, status, cert tracking evolve; campaign logic doesn't care |
| **Enrollment Service** | Enrollment protocol changes independently from registry storage | Enrollment code format, CSR handling, tier support evolve; registry schema doesn't care |
| **Certificate Authority** | Crypto changes independently from enrollment logic | Key algorithms, cert lifetime, HSM config change; enrollment flow doesn't care |
| **MQTT Broker** | Telemetry transport changes independently from validation | Broker software, QoS settings, topic structure change; validation pipeline doesn't care |
| **Auth Service (MQTT webhook)** | Broker auth format changes independently from policy logic | Webhook shape changes per broker vendor; OPA policy doesn't care |
| **Validation Pipeline** | Validation rules change independently from ingestion transport | Range checks, anomaly thresholds, rate limits change; MQTT/HTTP transport doesn't care |
| **Data Store** | Storage engine changes independently from business logic | Schema, indexes, partitioning change; domain logic uses repos |
| **Score Engine** | Gamification rules change independently from data collection | Score formula, badge criteria, sweepstakes rules change; ingestion doesn't care |
| **Notification** | Notification channels change independently from events that trigger them | Email/push/SMS change; security response logic doesn't care |
| **Proxy (Caddy)** | Routing changes independently from services behind it | URL structure, TLS termination, rate limiting change; services don't care |
| **O11y** | Observability vendor changes independently from everything | Traces/metrics/logs provider changes; no business code changes |

---

## Feature → Component → Handling

### Feature 1: Institutional Onboarding
**Requirements:** FR-001, FR-002, FR-003, FR-004 | **BUC-01** | All Must

| Step | Component | What it does |
|------|-----------|-------------|
| Admin navigates to org creation | **UI** | Org creation form, hierarchy editor, role assignment UI, member invite flow |
| Request hits platform | **Proxy** | Routes `/rootstock.v1.OrgService/*` to web server |
| Handler authenticates | **RPC Handler** → **Identity** | Validates Zitadel token, extracts user ID and org membership |
| Handler authorizes | **RPC Handler** → **Auth (OPA)** | Checks `org:create`, `org:manage_hierarchy`, `role:assign`, `member:invite` permissions |
| Create org / hierarchy / roles | **Campaign Service** (org ops) | Orchestrates: create tenant in Zitadel, create org record, configure hierarchy, define roles |
| Invite researcher | **Identity (Zitadel)** | Sends invite email, creates pending user, associates with org on accept |
| Persist org state | **Data Store** | Org, hierarchy, role, membership tables |
| Sync to OPA | **Auth (OPA)** | Bundle refresh picks up new org, roles, memberships within 30s |

**Cross-cutting:** O11y traces the full flow. Notification sends invite emails via Zitadel.

---

### Feature 2: Campaign Creation & Management
**Requirements:** FR-005, FR-006, FR-007, FR-008, FR-009, FR-010 | **BUC-02** | FR-010 is Should, rest Must

| Step | Component | What it does |
|------|-----------|-------------|
| Researcher defines campaign | **UI** | Campaign form: parameters (units, range, precision), region (polygon/radius/admin boundary/unbounded), time window (start/end UTC), quality thresholds, device eligibility |
| Request hits platform | **Proxy** | Routes `/rootstock.v1.CampaignService/*` |
| Handler authenticates + authorizes | **RPC Handler** → **Identity** → **Auth (OPA)** | Validates token, checks `campaign:create` in researcher's org |
| Validate campaign definition | **Campaign Service** | Field-level validation: all required fields present, window_end > window_start, region geometry valid, parameter ranges make physical sense |
| Persist campaign | **Data Store** | Campaign table with parameters, region (PostGIS geometry), window, thresholds, eligibility criteria, status (draft/published) |
| Publish campaign | **Campaign Service** | Sets status to published; becomes discoverable by scitizens |
| Campaign discovery | **Campaign Service** | Listing API with filters: geo proximity, sensor type compatibility, campaign status |
| Quality dashboard | **Campaign Service** + **Data Store** | Aggregates accepted/rejected counts, rejection reasons, geo distribution, temporal coverage; refreshed within 5 min of latest reading |

**Cross-cutting:** O11y. Campaign window open/close are temporal events (BE-08, BE-09) — campaign service handles scheduling.

---

### Feature 3: Scitizen Registration & Campaign Browsing
**Requirements:** FR-011, FR-012 | **BUC-03** | Both Must

| Step | Component | What it does |
|------|-----------|-------------|
| Scitizen registers | **UI** → **Identity (Zitadel)** | Open registration, email + password or social login. Identity fully delegated. Target: registration → campaign browse < 2 min |
| Browse campaigns | **UI** | Campaign listing filtered by location, device compatibility, status |
| Campaign listing API | **RPC Handler** → **Campaign Service** → **Data Store** | Returns published campaigns with summary, region, window, required parameters. Response < 2s |

**Cross-cutting:** No org membership required for scitizens. Cross-org participation is the normal case.

---

### Feature 4: Device Registration & Enrollment
**Requirements:** FR-013, FR-014, FR-015, FR-016 | **BUC-04** | All Must

| Step | Component | What it does |
|------|-----------|-------------|
| Scitizen initiates enrollment | **UI** | Generates enrollment code request |
| Generate enrollment code | **Enrollment Service** | 6-8 alphanumeric, no ambiguous chars, 15-min TTL, one-time use. Persisted with pending device ID |
| **Tier 1 (direct):** Device calls `/enroll` | **Enrollment Service** | Device presents enrollment code + locally-generated CSR over HTTPS |
| **Tier 2 (proxy):** Companion app mediates | **Enrollment Service** | Companion app gets CSR from device over BLE/WiFi, submits to `/enroll`. Private key never leaves device |
| Validate code + CSR | **Enrollment Service** | Code valid? Not expired? Not used? CSR well-formed? |
| Issue certificate | **Certificate Authority** (step-ca) | Signs CSR. CN = device-id. 90-day lifetime. No metadata in cert |
| Return cert to device | **Enrollment Service** | Direct: HTTPS response. Proxy: companion app pushes cert back to device |
| Create registry entry | **Device Registry** → **Data Store** | device ID, owner ID, status (pending → active), device class, firmware version, tier, sensor capabilities, cert serial |
| Sync to OPA | **Auth (OPA)** | Bundle refresh picks up new device within 30s |

**Cross-cutting:** SEC-002 (private key never leaves device). O11y traces enrollment flow.

---

### Feature 5: Campaign Enrollment
**Requirements:** FR-017, FR-018, FR-019 | **BUC-05** | All Must

| Step | Component | What it does |
|------|-----------|-------------|
| Scitizen selects campaign for device | **UI** | Shows eligible campaigns for selected device |
| Eligibility check | **Campaign Service** | Checks device class, tier, sensor capabilities, firmware version against campaign criteria. Rejects ineligible with reason |
| Enroll device in campaign | **Device Registry** → **Data Store** | Adds campaign association to device record |
| OPA picks up enrollment | **Auth (OPA)** | Bundle refresh (≤30s) — device now allowed to publish to campaign topic |
| Push campaign config to device | **MQTT Broker** | Retained message on `rootstock/{device-id}/config` with campaign parameters (sampling rate, etc.) |
| Multi-campaign support | **Device Registry** | Device can be enrolled in N campaigns simultaneously. Topic routing associates readings with correct campaign |

---

### Feature 6: Data Ingestion & Validation
**Requirements:** FR-020, FR-021, FR-022, FR-023, FR-024, FR-025 | **BUC-06** | FR-025 is Should, rest Must

| Step | Component | What it does |
|------|-----------|-------------|
| Device connects | **MQTT Broker** or **HTTP/2 Gateway** | mTLS handshake. Platform verifies device cert against CA chain. Extracts device ID from CN |
| Authenticate (mTLS) | TLS layer (**MQTT Broker** / **Proxy**) | Reject: no cert, expired cert, self-signed, wrong CA. Structured diagnostic log on failure |
| Authorize (OPA) | **Auth Service** → **Auth (OPA)** | On every MQTT connect/publish/subscribe or HTTP request: check device status (active?), campaign enrollment, topic ACL |
| Topic ACL enforcement | **Auth (OPA)** | Device can only publish to `rootstock/{its-own-device-id}/*`. Cross-device publishing denied + logged |
| Receive telemetry | **MQTT Broker** / **HTTP/2 Gateway** | Raw reading: value, timestamp, geolocation, device ID, campaign ID |
| Validate reading | **Validation Pipeline** | Schema (correct fields/types), parameter range (within campaign bounds), rate limit (per device), geolocation (within campaign region), timestamp (within campaign window). Reject with specific reason |
| Anomaly flagging | **Validation Pipeline** | Readings outside 3σ of campaign rolling average → quarantined for review, not deleted |
| Persist with provenance | **Data Store** | Reading + device ID, timestamp, geolocation, firmware version, cert serial, campaign ID, ingestion timestamp |
| Update contribution | **Score Engine** | Accepted reading triggers score update (async) |

**Performance:** 10K readings/sec sustained (PE-001), P99 < 200ms (PE-002), OPA < 5ms P99 (PE-005).

---

### Feature 7: Data Export & Analysis
**Requirements:** FR-026, FR-027 | **BUC-07** | Both Must

| Step | Component | What it does |
|------|-----------|-------------|
| Researcher requests export | **UI** | Export configuration: format, filters, date range |
| Handler authenticates + authorizes | **RPC Handler** → **Identity** → **Auth (OPA)** | Checks `campaign:export` permission for this campaign in researcher's org |
| Query validated readings | **Data Store** | Readings with full provenance metadata, quality metrics, campaign context |
| Pseudonymize device IDs | **Campaign Service** (export ops) | Hash device IDs. No scitizen names, emails, or unhashed device IDs in output. Spatial resolution per campaign policy |
| Package export | **Campaign Service** | Readings + provenance + quality metrics. FAIR-compliant metadata |
| Deliver | **RPC Handler** | Stream or download. < 60s for up to 1M readings |

**Cross-cutting:** SEC-004 (identity separation). FACT-009 (4 spatiotemporal points = 95% re-identification).

---

### Feature 8: Certificate Lifecycle Management
**Requirements:** FR-028, FR-029, FR-030 | **BUC-08** | All Must

| Step | Component | What it does |
|------|-----------|-------------|
| Renewal trigger (day 60 of 90) | Device-initiated | Device presents current cert via mTLS + new CSR to `POST /renew` |
| Validate renewal | **Enrollment Service** | Cert valid? Cert matches device in registry? CSR well-formed? |
| Issue new cert | **Certificate Authority** | Signs new CSR. Fresh 90-day window |
| Update registry | **Device Registry** → **Data Store** | New cert serial, new expiry |
| Grace period (days 90-97) | **Enrollment Service** | Expired ≤7 days: `/renew` only. Expired >7 days: denied, must re-enroll |
| Revocation | **Device Registry** → **Auth (OPA)** | Set device status to `revoked` in registry. OPA denies all actions within 30s (bundle refresh). No CRL/OCSP needed |

---

### Feature 9: Device Security Response
**Requirements:** FR-031, FR-032, FR-033 | **BUC-09** | FR-031 Must, FR-032/033 Should

| Step | Component | What it does |
|------|-----------|-------------|
| Vulnerability discovered | External trigger (BE-10) | Specific device class + firmware version range identified |
| Bulk suspension | **Device Registry** → **Data Store** | Batch update: all matching devices → status `suspended`. < 60s |
| OPA denial | **Auth (OPA)** | Bundle refresh → suspended devices denied on next request |
| Mass notification | **Notification** | Alert affected scitizens with remediation instructions |
| Flag vulnerable-window data | **Validation Pipeline** → **Data Store** | Readings from affected devices during vulnerability window → quarantined for researcher review |
| Reinstatement after patch | **Device Registry** | Device reports patched firmware → status restored to `active`. Campaign enrollments preserved |

---

### Feature 10: Recognition & Incentives
**Requirements:** FR-034, FR-035, FR-036 | **BUC-10** | FR-034 Must, FR-035/036 Should

| Step | Component | What it does |
|------|-----------|-------------|
| Reading accepted | **Validation Pipeline** (event) | Triggers async score update |
| Compute contribution score | **Score Engine** | Volume, quality rate, consistency, campaign diversity. Updated within 15 min |
| Check badge milestones | **Score Engine** | First contribution, 100/1000 readings, campaign completion, 30-day streak, geographic diversity |
| Award badges | **Score Engine** → **Data Store** | Persist badge, visible on profile |
| Grant sweepstakes entries | **Score Engine** → **Data Store** | Variable-ratio schedule tied to score milestones. Auditable, tamper-evident |
| Dashboard display | **UI** | Contribution score, badges, sweepstakes entries, active campaigns, accepted reading count |

**Cross-cutting:** Self-Determination Theory — recognition must not crowd out intrinsic motivation. Variable-ratio reinforcement (FACT-008).

---

## Component Interaction Summary

```
                                    ┌─────────────┐
                                    │   Identity   │
                                    │  (Zitadel)   │
                                    └──────┬───────┘
                                           │ authn
┌────────┐    ┌───────┐    ┌───────────────┼───────────────────┐
│   UI   │───▸│ Proxy │───▸│          RPC Handlers             │
│(Svelte)│    │(Caddy)│    │  (Connect RPC, proto-only)        │
└────────┘    └───┬───┘    └───────────┬───────────────────────┘
                  │                    │
                  │              ┌─────┴──────┐
                  │              │ Auth (OPA) │◂──── OPA Bundle ◂── Device Registry
                  │              └─────┬──────┘
                  │                    │ authz
                  │              ┌─────┴──────────────────────────────────┐
                  │              │              Flows                      │
                  │              │  (orchestrate across ops, no repo calls)│
                  │              └─────┬──────────────────────────────────┘
                  │                    │
                  │    ┌───────────────┼───────────────────────┐
                  │    │               │                       │
                  │    ▼               ▼                       ▼
                  │ Campaign Ops   Device Ops            Score Ops
                  │    │               │                       │
                  │    ▼               ▼                       ▼
                  │ Campaign Repo  Device Repo           Score Repo
                  │    │               │                       │
                  │    └───────────────┼───────────────────────┘
                  │                    │
                  │              ┌─────▼──────┐
                  │              │ Data Store  │
                  │              │ (Postgres)  │
                  │              └────────────┘
                  │
           ┌──────┴──────────────────────────────────┐
           │        Device Data Path                  │
           │                                          │
  ┌────────▼─────┐     ┌──────────────┐     ┌────────┴──────┐
  │ MQTT Broker  │     │ Auth Service  │     │ HTTP/2 Gateway│
  │   (EMQX)    │────▸│ (webhook→OPA) │◂────│               │
  └──────┬───────┘     └──────────────┘     └───────┬───────┘
         │                                          │
         └────────────────┬─────────────────────────┘
                          │
                    ┌─────▼──────────┐
                    │  Validation    │
                    │  Pipeline      │
                    └─────┬──────────┘
                          │
                    ┌─────▼──────┐
                    │ Data Store  │
                    └────────────┘

  ┌──────────────┐          ┌──────────────┐
  │ Enrollment   │─────────▸│     CA       │
  │  Service     │          │  (step-ca)   │
  │ /enroll      │          └──────────────┘
  │ /renew       │
  └──────┬───────┘
         │
  ┌──────▼───────┐
  │Device Registry│
  └──────────────┘
```

---

## Requirement Coverage Matrix

| Feature | FR IDs | Components Touched | Priority Mix |
|---------|--------|-------------------|-------------|
| 1. Institutional Onboarding | FR-001–004 | UI, Proxy, RPC Handler, Identity, Auth, Data Store | 4 Must |
| 2. Campaign Management | FR-005–010 | UI, Proxy, RPC Handler, Identity, Auth, Campaign Service, Data Store | 5 Must, 1 Should |
| 3. Scitizen Registration | FR-011–012 | UI, Identity, RPC Handler, Campaign Service, Data Store | 2 Must |
| 4. Device Registration | FR-013–016 | UI, Enrollment Service, CA, Device Registry, Data Store, Auth | 4 Must |
| 5. Campaign Enrollment | FR-017–019 | UI, RPC Handler, Campaign Service, Device Registry, Auth, MQTT Broker | 3 Must |
| 6. Data Ingestion | FR-020–025 | MQTT Broker, HTTP/2 Gateway, Auth Service, Auth, Validation Pipeline, Data Store | 5 Must, 1 Should |
| 7. Data Export | FR-026–027 | UI, RPC Handler, Auth, Campaign Service, Data Store | 2 Must |
| 8. Certificate Lifecycle | FR-028–030 | Enrollment Service, CA, Device Registry, Auth | 3 Must |
| 9. Security Response | FR-031–033 | Device Registry, Auth, Notification, Validation Pipeline, Data Store | 1 Must, 2 Should |
| 10. Recognition | FR-034–036 | Score Engine, UI, Data Store | 1 Must, 2 Should |

**Totals:** 16 components, 36 functional requirements, 10 features, 29 Must / 7 Should

---

## Non-Functional Requirement Mapping

| NFR | Components Affected | Constraint |
|-----|-------------------|-----------|
| PE-001 (10K readings/sec) | Validation Pipeline, MQTT Broker, Data Store | Ingestion path must be horizontally scalable |
| PE-002 (P99 < 200ms ingestion) | Validation Pipeline, Data Store | No synchronous external calls in hot path |
| PE-003 (P95 < 500ms API) | RPC Handlers, Data Store | Query optimization, connection pooling |
| PE-004 (99.9% availability) | All ingestion-path components | Health checks, graceful degradation |
| PE-005 (OPA < 5ms P99) | Auth (OPA) | In-memory policy evaluation, efficient bundle format |
| LF-001 (consistent visual identity) | UI | Shared design system |
| LF-002 (scientific credibility) | UI | SI units, sig figs, proper chart labels |
| LF-003 (responsive 320–2560px) | UI | Mobile-first for scitizen flows |
| US-001 (enrollment < 10 min) | UI, Enrollment Service | Minimal steps, clear feedback |
| US-002 (campaign creation < 15 min) | UI, Campaign Service | Guided form, sensible defaults |
| US-003 (contribution feedback) | UI, Score Engine | Dashboard loads < 3s |
| US-004 (actionable errors) | All user-facing components | Structured error responses |
| SEC-001 (mTLS everywhere) | MQTT Broker, HTTP/2 Gateway, Enrollment Service | TLS termination config |
| SEC-002 (private key on device) | Enrollment Service, Companion App | CSR-only protocol |
| SEC-003 (RBAC via OPA) | Auth (OPA) | Rego policies with >90% test coverage |
| SEC-004 (identity separation) | Data Store, Campaign Service (export) | Schema-level separation, pseudonymization |
| SEC-005 (HSM for CA keys) | Certificate Authority | PKCS#11, SoftHSM (dev) / YubiHSM (prod) |
| SEC-006 (connection diagnostics) | MQTT Broker, Auth Service, Enrollment Service | Structured logs with specific failure reasons |
| OP-001 (container deployment) | All | OCI images, compose for full stack |
| OP-002 (dual protocol) | MQTT Broker, HTTP/2 Gateway | Same validation pipeline behind both |
| OP-003 (graceful degradation) | MQTT Broker, Device Registry | LWT, auto-reconnect, no re-enrollment |
