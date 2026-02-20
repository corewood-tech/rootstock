# ROOTSTOCK — Backend Development Roadmap

> Derived from the ops/flows/repos decomposition. Ordered by BUC dependency chain and implementation prerequisites.

---

## What Exists

| Layer | What's built | Status |
|-------|-------------|--------|
| **Config** | Multi-source loading (YAML, env, flags), CertConfig | Complete |
| **Globals** | O11y (OTel-backed), Events (DBOS-backed), Auth (OPA-backed) | Complete |
| **Repos** | AuthorizationRepo (OPA), ObservabilityRepo (OTel), EventsRepo (DBOS), SQL pool (pgx), IdentityRepo (Zitadel), CampaignRepo, DeviceRepo, ReadingRepo, ScoreRepo, CertRepo (in-process CA) | Complete |
| **Ops** | Org (5), Campaign (5), Device (8 + UpdateCertSerial), Cert (2), Reading (6), Score (4), Pure (ValidateReading, MatchEligibility) | Complete |
| **Flows** | OnboardInstitution (5 sub-flows), CreateCampaign, PublishCampaign, BrowseCampaigns, CampaignDashboard, RegisterDevice, RenewCert, GetDevice, RevokeDevice, ReinstateDevice, GetCACert, IngestReading, UpdateContributionScore, GetContribution | Complete |
| **Interceptors** | AuthorizationInterceptor, BinaryOnlyInterceptor, OTel auto-instrumentation | Complete |
| **Handlers** | HealthService, OrgService, CampaignService, ScoreService, DeviceService (Connect RPC), EnrollHandler (/enroll, /ca — HTTP) | Complete |
| **Proto** | HealthService, OrgService, CampaignService, ScoreService, DeviceService | Complete |
| **Server** | `server/rpc.go` (Connect RPC wiring), main.go decomposed | Complete |
| **Infra** | Postgres (app + Zitadel), Zitadel + Login v2, Caddy, OTel Collector, Prometheus, Tempo, Loki, Grafana, in-process CA | Running |
| **Migrations** | 5 migration pairs (campaigns, devices, readings, scores, unique constraints) | Complete |
| **Tests** | Unit tests across all repos, ops, flows (17 packages) | Complete |
| **Recently built** | MQTT pipeline (inline subscriptions, telemetry + renewal callbacks), EnrollInCampaign flow, MQTTOps (PushDeviceConfig), RefreshScitizenScoreFlow (score trigger after accepted reading), ExportData flow + PseudonymizeExport pure op, ExportCampaignData RPC, MQTT grace period (custom TLS verify + ACL restriction for expired certs), Offset pagination for QueryReadings, HMAC export config | Complete |
| **Phase 9** | SecurityResponse flow, NotificationRepo (log stub), NotificationOps, AdminService (SuspendByClass RPC) | Complete |

---

## Dependency Chain (from spec)

```
Phase 1 ──► Phase 2 ──► Phase 3 ──────────► Phase 4 ──► Phase 5
  Org     Campaign    Scitizen + Device    Enrollment    Ingestion
                                                            │
                                              Phase 6 ◄─────┤
                                              Export         │
                                                        Phase 7
                                                        Score
Phase 8: Cert Lifecycle (depends on Phase 3)
Phase 9: Security Response (independent, needs Phase 3 + Phase 5)
```

---

## Phase 0: Pipeline Foundation

**Goal:** Establish the flow/ops/repo pattern with working database schema so all subsequent phases build on proven plumbing.

### 0.1 — Database Migrations

Create migration tooling and base schema. Tables owned by each repo:

**CampaignRepo tables:**
```
campaigns (id, org_id, status, window_start, window_end, created_by, created_at)
campaign_parameters (id, campaign_id, name, unit, min_range, max_range, precision)
campaign_regions (id, campaign_id, geometry)  — PostGIS
campaign_eligibility (id, campaign_id, device_class, tier, required_sensors, firmware_min)
```

**DeviceRepo tables:**
```
devices (id, owner_id, status, class, firmware_version, tier, sensors, cert_serial, created_at)
enrollment_codes (code, device_id, expires_at, used)
device_campaigns (device_id, campaign_id, enrolled_at)
```

**ReadingRepo tables:**
```
readings (id, device_id, campaign_id, value, timestamp, geolocation, firmware_version, cert_serial, ingested_at, status, quarantine_reason)
```

**ScoreRepo tables:**
```
scores (scitizen_id, volume, quality_rate, consistency, diversity, total, updated_at)
badges (id, scitizen_id, badge_type, awarded_at)
sweepstakes_entries (id, scitizen_id, entries, milestone_trigger, granted_at)
```

**Decision:** Migration tool. Options: goose, golang-migrate, or raw SQL files executed by an init container. Pick one.

### 0.2 — Business Repo Interfaces

Define interfaces for the 4 Postgres-backed business repos. Each follows the established pattern: interface in `interface.go`, types in `receive_types.go` / `emit_types.go`, channel-based implementation in `repository.go`.

```
repo/
├── authorization/     ← exists
├── events/            ← exists
├── observability/     ← exists
├── sql/connect/       ← exists
├── campaign/          ← NEW
│   ├── interface.go
│   ├── receive_types.go
│   ├── emit_types.go
│   └── repository.go
├── device/            ← NEW
│   ├── interface.go
│   ├── receive_types.go
│   ├── emit_types.go
│   └── repository.go
├── reading/           ← NEW
│   ├── interface.go
│   ├── receive_types.go
│   ├── emit_types.go
│   └── repository.go
└── score/             ← NEW
    ├── interface.go
    ├── receive_types.go
    ├── emit_types.go
    └── repository.go
```

Each repo gets the pgx pool injected at construction. All mutable state managed via goroutine + channels (no mutexes).

### 0.3 — Ops and Flows Directory Structure

```
ops/
├── org/           ← ops 1-5, calls IdentityRepo
├── campaign/      ← ops 6-10, calls CampaignRepo
├── device/        ← ops 11-18, calls DeviceRepo
├── cert/          ← op 19, calls CertRepo (Phase 3)
├── reading/       ← ops 20-25, calls ReadingRepo
├── score/         ← ops 26-29, calls ScoreRepo
└── pure/          ← ops 32-33, no repo (ValidateReading, MatchEligibility)

flows/
├── onboard_institution.go
├── create_campaign.go
├── publish_campaign.go
├── browse_campaigns.go
├── register_device.go
├── enroll_in_campaign.go
├── ingest_reading.go
├── export_data.go
├── renew_cert.go
├── revoke_device.go
├── security_response.go
├── reinstate_device.go
├── update_score.go
└── campaign_dashboard.go
```

### 0.4 — Wire Identity into Auth Interceptor

The `AuthorizationInterceptor` has a TODO for `session_user_id`. Zitadel is running. Wire:
1. Extract bearer token from request headers
2. Introspect/validate against Zitadel (or verify JWT locally)
3. Pass `session_user_id` + org membership into `AuthzInput`
4. Update OPA policy to use `session_user_id` and org context for real RBAC

This unblocks all authenticated flows.

---

## Phase 1: Org + Campaign (Researcher Side)

**Goal:** A researcher can create an org, invite members, create and publish a campaign, and browse campaigns.

**BUCs:** BUC-01, BUC-02 | **FRs:** FR-001–010 (9 Must, 1 Should)

### 1.1 — Proto: OrgService + CampaignService

```protobuf
service OrgService {
  rpc CreateOrg(CreateOrgRequest) returns (CreateOrgResponse);
  rpc NestOrg(NestOrgRequest) returns (NestOrgResponse);
  rpc DefineRole(DefineRoleRequest) returns (DefineRoleResponse);
  rpc AssignRole(AssignRoleRequest) returns (AssignRoleResponse);
  rpc InviteUser(InviteUserRequest) returns (InviteUserResponse);
}

service CampaignService {
  rpc CreateCampaign(CreateCampaignRequest) returns (CreateCampaignResponse);
  rpc PublishCampaign(PublishCampaignRequest) returns (PublishCampaignResponse);
  rpc ListCampaigns(ListCampaignsRequest) returns (ListCampaignsResponse);
  rpc GetCampaignDashboard(GetCampaignDashboardRequest) returns (GetCampaignDashboardResponse);
}
```

### 1.2 — IdentityRepo (wraps Zitadel)

New repo at `repo/identity/`. Wraps Zitadel's gRPC Management API:
- CreateOrg → Zitadel Org API
- CreateSubOrg → Zitadel nested Org
- CreateRole → Zitadel Project Role
- AssignRole → Zitadel User Grant
- InviteUser → Zitadel User Import / Invite

### 1.3 — Org Ops (5 ops)

Each op calls IdentityRepo. Thin business logic layer:
- `CreateOrg` — validate name, call repo, return org ID
- `NestOrg` — validate parent exists, nesting depth ≤5, call repo
- `DefineRole` — validate permissions are known, call repo
- `AssignRole` — validate user exists, role exists, call repo
- `InviteUser` — validate email, call repo

### 1.4 — Campaign Ops (5 ops)

Each op calls CampaignRepo:
- `CreateCampaign` — validate params/region/window/thresholds, persist
- `PublishCampaign` — check status is draft, set to published
- `ListCampaigns` — query with filters (geo proximity via PostGIS, sensor type, status)
- `GetCampaignRules` — return validation criteria (used later by IngestReading)
- `GetCampaignEligibility` — return eligibility criteria (used later by EnrollInCampaign)

### 1.5 — CampaignRepo Implementation

Postgres queries for campaign CRUD. PostGIS for geographic queries. Uses pgx pool from `sql/connect`.

### 1.6 — Flows

- `OnboardInstitution` — composes: CreateOrg → NestOrg → DefineRole → AssignRole → InviteUser
- `CreateCampaign` — composes: CreateCampaign (single op)
- `PublishCampaign` — composes: PublishCampaign (single op)
- `BrowseCampaigns` — composes: ListCampaigns (single op)
- `CampaignDashboard` — composes: GetCampaignQuality (reading op, deferred to Phase 5)

### 1.7 — Handlers

- `OrgServiceHandler` — maps RPC → OnboardInstitution flow
- `CampaignServiceHandler` — maps RPC → campaign flows

### 1.8 — OPA Policy Update

Add org-scoped RBAC rules:
- `org:create` — platform-level (Corewood admin or self-service)
- `org:manage` — org admin only
- `campaign:create` — researcher role in org
- `campaign:publish` — researcher role in org

### 1.9 — E2E Test

Playwright: researcher logs in via Zitadel → creates campaign → publishes → verifies in listing.

**Deliverable:** Researcher can create org, invite members, create and publish campaigns. Campaigns are browsable.

---

## Phase 2: Scitizen Registration + Campaign Browsing

**Goal:** A scitizen can register and browse published campaigns.

**BUCs:** BUC-03 | **FRs:** FR-011, FR-012 (2 Must)

### 2.1 — Zitadel Configuration

- Enable open registration (no org required)
- Configure scitizen role (no campaign creation, can browse + enroll)
- Social login providers (optional, defer to later)

### 2.2 — OPA Policy Update

- Scitizen role can call `ListCampaigns`
- Scitizen role cannot call `CreateCampaign`, `PublishCampaign`, org management

### 2.3 — E2E Test

Scitizen registers → browses campaigns → sees published campaigns from Phase 1.

**Deliverable:** Scitizens can register and browse. Smallest possible phase — enables Phase 3.

---

## Phase 3: Device Registration

**Goal:** A scitizen can register a Tier 1 device and get a certificate.

**BUCs:** BUC-04 | **FRs:** FR-013–016 (4 Must)

### 3.1 — Infrastructure: In-Process CA

**Decision (implemented):** In-process CA using Go's `crypto/x509` stdlib instead of step-ca. The cert repo wraps the crypto implementation the same way the identity repo wraps Zitadel. No external container dependency. `make ca-init` generates dev CA + server leaf cert.

### 3.2 — CertRepo

New repo at `repo/cert/`. In-process X.509 CA:
- `IssueCert(ctx, IssueCertInput) → IssuedCert` — parse CSR, validate key, sign with CA, return cert PEM
- `GetCACert(ctx) → CACert` — return CA cert PEM for device bootstrapping
- Channel-based concurrency; CA key exclusively owned by manage() goroutine

### 3.3 — DeviceRepo Implementation

Postgres queries for device registry CRUD:
- enrollment_codes table: generate, validate, mark used
- devices table: create, read, update status, query by class

### 3.4 — Device Ops (8 ops)

- `GenerateEnrollmentCode` — create code, store with TTL
- `RedeemEnrollmentCode` — validate + mark used
- `CreateDevice` — insert registry entry
- `GetDevice` — read by ID
- `GetDeviceCapabilities` — return class, tier, sensors
- `UpdateDeviceStatus` — status transition (the most reused op)
- `QueryDevicesByClass` — batch query
- `EnrollDeviceInCampaign` — add campaign association (used in Phase 4)

### 3.5 — Cert Ops (1 op)

- `IssueCert` — call CertRepo with CSR, return signed cert

### 3.6 — RegisterDevice Flow

```
RedeemEnrollmentCode → IssueCert → CreateDevice → UpdateDeviceStatus(pending→active)
```

### 3.7 — Proto: DeviceService (admin ops) + HTTP Enrollment

```protobuf
service DeviceService {
  rpc GetDevice(GetDeviceRequest) returns (GetDeviceResponse);
  rpc RevokeDevice(RevokeDeviceRequest) returns (RevokeDeviceResponse);
  rpc ReinstateDevice(ReinstateDeviceRequest) returns (ReinstateDeviceResponse);
}
```

Enrollment is NOT a Connect RPC call — it's a plain HTTP POST on the same port (8080). Devices don't have JWT tokens. Auth is via enrollment code in the request body.

- `POST /enroll` — enrollment code + CSR → cert PEM. No mTLS required.
- `GET /ca` — public, returns CA cert PEM for device bootstrapping.

### 3.8 — Handler: DeviceServiceHandler + EnrollHandler

- `DeviceServiceHandler` — Connect RPC handler for admin ops (GetDevice, Revoke, Reinstate). JWT auth.
- `EnrollHandler` — plain HTTP handler for `/enroll` and `/ca`. Mounted on the same mux as Connect RPC.

### 3.9 — E2E Test

Scitizen generates enrollment code → simulated device calls `POST /enroll` with code + CSR → receives cert → device appears in registry as active.

**Deliverable:** Devices can be registered with real certificates. Foundation for ingestion.

---

## Phase 4: Campaign Enrollment

**Goal:** A scitizen can enroll a registered device in a published campaign.

**BUCs:** BUC-05 | **FRs:** FR-017–019 (3 Must)

### 4.1 — Infrastructure: Embedded MQTT Broker (Mochi MQTT)

**Decision:** Embed the MQTT broker in-process using [Mochi MQTT](https://github.com/mochi-mqtt/server/v2) (`github.com/mochi-mqtt/server/v2`). Same rationale as in-process CA — no external container dependency. No compose-mqtt.yml. No auth webhook service.

- MQTT listener on port 8883 with strict mTLS (`RequireAndVerifyClientCert`)
- Auth hook is a Go struct implementing `mqtt.Hook` — extracts device ID from mTLS cert CN, enforces topic ACLs in-process
- InlineClient enabled for server-side publish (config push) and subscribe (telemetry consumption)
- Topic structure: `rootstock/{device-id}/data/{campaign-id}`, `rootstock/{device-id}/config`, `rootstock/{device-id}/renew`, `rootstock/{device-id}/cert`

Two ports total for the platform:
- **Port 8080**: Connect RPC (JWT) + `/enroll` (enrollment code) + `/ca` (public)
- **Port 8883**: MQTT (mTLS) — all post-enrollment device traffic

### 4.2 — MQTTRepo

New repo at `repo/mqtt/`. Wraps embedded Mochi server's inline client:
- `PushDeviceConfig(ctx, deviceID, config)` — `server.Publish("rootstock/{device-id}/config", payload, true, 1)` (retained message)
- `PublishToDevice(ctx, topic, payload)` — generic device publish (e.g., renewal cert response)
- Channel-based concurrency (same pattern as all repos)

### 4.3 — Pure Logic: MatchEligibility

`ops/pure/match_eligibility.go`:
- Input: device capabilities + campaign eligibility criteria
- Output: eligible (bool) + reason (string)
- No repo, no I/O, deterministic, unit-testable

### 4.4 — EnrollInCampaign Flow

```
GetDeviceCapabilities + GetCampaignEligibility → MatchEligibility → EnrollDeviceInCampaign → PushDeviceConfig
```

Crosses 3 clusters + pure logic. This is the flow that bridges researcher ↔ scitizen.

### 4.5 — Proto: Add to CampaignService or DeviceService

```protobuf
service DeviceService {
  rpc EnrollInCampaign(EnrollInCampaignRequest) returns (EnrollInCampaignResponse);
  rpc GetDevice(GetDeviceRequest) returns (GetDeviceResponse);
}
```

### 4.6 — MQTT Auth Hook (replaces OPA webhook + Auth Service)

In-process Go hook implementing `mqtt.Hook`:
- `OnConnectAuthenticate` — extracts device ID from mTLS cert CN, verifies against CA cert pool
- `OnACLCheck` — enforces topic ACLs: device can only publish/subscribe to `rootstock/{own-device-id}/*`
- Device status check (active → allow, suspended/revoked → deny) via device ops lookup
- Campaign enrollment check (device enrolled in campaign → allow publish to campaign topic)

No separate auth service. No webhook. The hook is a Go struct with direct access to device ops.

### 4.7 — E2E Test

Scitizen enrolls device in campaign → device appears in campaign's enrolled devices → config retained message published via embedded broker.

**Deliverable:** Devices enrolled in campaigns. MQTT broker embedded and operational. Auth enforced in-process.

---

## Phase 5: Data Ingestion

**Goal:** An enrolled device can publish readings. Readings are validated and persisted with full provenance.

**BUCs:** BUC-06 | **FRs:** FR-020–025 (5 Must, 1 Should)

### 5.1 — ReadingRepo Implementation

Postgres (or TimescaleDB) queries:
- `PersistReading` — insert with all provenance fields
- `QuarantineReading` — update status to quarantined
- `QueryReadings` — campaign-scoped query with filters
- `PseudonymizeExport` — hash device IDs in query results
- `QuarantineByWindow` — batch update by device + time range
- `GetCampaignQuality` — aggregated metrics

### 5.2 — Pure Logic: ValidateReading

`ops/pure/validate_reading.go`:
- Input: reading + campaign rules
- Checks: schema, parameter range, rate limit, geolocation (PostGIS point-in-polygon), timestamp within window
- Output: valid (bool) + rejection reasons[]
- No repo, no I/O, deterministic

### 5.3 — Reading Ops (6 ops)

Each calls ReadingRepo:
- `PersistReading`, `QuarantineReading`, `QueryReadings`, `PseudonymizeExport`, `QuarantineByWindow`, `GetCampaignQuality`

### 5.4 — IngestReading Flow

```
GetCampaignRules → ValidateReading(pure) → PersistReading → [QuarantineReading]
```

Hot path. GetCampaignRules should be cached (campaign rules don't change mid-window). ValidateReading is pure — no I/O.

### 5.5 — MQTT Inline Subscription (replaces separate consumer)

**Decision:** No separate MQTT consumer process. The embedded Mochi broker's InlineClient subscribes to `rootstock/+/data/+`. The callback runs in-process:

```go
server.Subscribe("rootstock/+/data/+", 1, func(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
    // Extract device ID and campaign ID from topic path segments
    // Deserialize reading from pk.Payload (protobuf)
    // Call IngestReading flow directly
})
```

This is the hot path. The callback must be non-blocking for sustained 10K msg/sec. Heavy work (DB persist, score update) runs in goroutines.

### 5.6 — HTTP/2 Ingestion Endpoint (Dual Protocol)

```protobuf
service IngestionService {
  rpc SubmitReading(SubmitReadingRequest) returns (SubmitReadingResponse);
}
```

Same IngestReading flow behind both MQTT and HTTP/2 paths (OP-002). HTTP/2 path uses Connect RPC on port 8080 with device cert auth (or JWT — TBD).

### 5.7 — E2E Test

Device publishes reading via MQTT → reading persisted with provenance → Prometheus metric increments → campaign quality shows accepted count.

**Deliverable:** Data flows end-to-end. Validated readings with full provenance. Hot path operational.

---

## Phase 6: Data Export

**Goal:** A researcher can export validated campaign data with pseudonymized device IDs.

**BUCs:** BUC-07 | **FRs:** FR-026–027 (2 Must)

### 6.1 — ExportData Flow

```
QueryReadings → PseudonymizeExport
```

Single-cluster (ReadingRepo). Ops already built in Phase 5.

### 6.2 — Proto: Add to CampaignService

```protobuf
rpc ExportCampaignData(ExportRequest) returns (stream ExportChunk);
```

Server-streaming for large exports (up to 1M readings, < 60s).

### 6.3 — CampaignDashboard Flow

```
GetCampaignQuality
```

Also already built in Phase 5. Wire handler.

### 6.4 — E2E Test

Researcher exports campaign data → exported data has pseudonymized device IDs → no PII in output.

**Deliverable:** Researchers can extract and analyze collected data.

---

## Phase 7: Recognition & Incentives

**Goal:** Contribution scores computed, badges awarded, sweepstakes entries granted.

**BUCs:** BUC-10 | **FRs:** FR-034–036 (1 Must, 2 Should)

### 7.1 — ScoreRepo Implementation

Postgres queries for scores, badges, sweepstakes.

### 7.2 — Score Ops (4 ops)

- `UpdateScore` — recompute from reading stats
- `CheckMilestones` — evaluate badge criteria
- `AwardBadge` — persist badge
- `GrantSweepstakes` — add entries

### 7.3 — UpdateContributionScore Flow

```
UpdateScore → CheckMilestones → AwardBadge → GrantSweepstakes
```

Triggered async after reading acceptance in IngestReading (event/channel from Phase 5).

### 7.4 — Proto: ScoreService

```protobuf
service ScoreService {
  rpc GetContribution(GetContributionRequest) returns (GetContributionResponse);
}
```

Read-only for the scitizen dashboard.

### 7.5 — E2E Test

Device submits readings → score updates within 15 min → badges awarded at milestones.

**Deliverable:** Gamification loop operational.

---

## Phase 8: Certificate Lifecycle

**Goal:** Devices can renew certs. Devices can be revoked. Grace period enforced.

**BUCs:** BUC-08 | **FRs:** FR-028–030 (3 Must)

### 8.1 — RenewCert Flow

```
IssueCert → UpdateCertSerial
```

Both ops already exist from Phase 3. New flow, no new ops. **Renewal happens over MQTT**, not HTTP:
1. Device publishes CSR to `rootstock/{device-id}/renew` (already mTLS-authenticated)
2. MQTT inline subscription calls RenewCert flow
3. Server publishes new cert PEM to `rootstock/{device-id}/cert`
4. Device stores new cert, reconnects

### 8.2 — RevokeDevice Flow

```
UpdateDeviceStatus(→revoked)
```

Single op, already exists. Auth hook denies revoked devices on next MQTT connect/action.

### 8.3 — Grace Period Logic

In the MQTT auth hook's `OnConnectAuthenticate`:
- Expired ≤7 days → allow connection, restrict to `rootstock/{device-id}/renew` topic only
- Expired >7 days → reject connection, must re-enroll via HTTP `/enroll`

### 8.4 — E2E Test

Device with expiring cert → publishes to `rootstock/{id}/renew` → receives new cert on `rootstock/{id}/cert`. Revoked device → MQTT connection rejected.

**Deliverable:** Full certificate lifecycle. No new ops needed — pure flow composition.

---

## Phase 9: Security Response

**Goal:** Bulk suspend vulnerable devices, flag data, notify scitizens, reinstate after patch.

**BUCs:** BUC-09 | **FRs:** FR-031–033 (1 Must, 2 Should)

### 9.1 — NotificationRepo

New repo at `repo/notification/`. Initial implementation: log-based (print notification). Real provider (email/push) deferred.

### 9.2 — SecurityResponse Flow

```
QueryDevicesByClass → UpdateDeviceStatus(bulk→suspended) → QuarantineByWindow → Notify
```

All ops except Notify already exist. Cross-cluster: Device + Reading + Notification.

### 9.3 — ReinstateDevice Flow

```
UpdateDeviceStatus(suspended→active)
```

Single op, already exists.

### 9.4 — Proto: AdminService

```protobuf
service AdminService {
  rpc SuspendByClass(SuspendByClassRequest) returns (SuspendByClassResponse);
  rpc ReinstateDevice(ReinstateDeviceRequest) returns (ReinstateDeviceResponse);
}
```

### 9.5 — E2E Test

Bulk suspend by class → affected devices denied → data flagged → reinstate after firmware update.

**Deliverable:** Security response operational. Completes all BUCs.

---

## Phase Summary

| Phase | BUC | New Repos | New Ops | New Flows | New Proto Services | Infra Added |
|-------|-----|-----------|---------|-----------|-------------------|-------------|
| **0** | — | CampaignRepo, DeviceRepo, ReadingRepo, ScoreRepo (interfaces) | — | — | — | Migrations |
| **1** | 01, 02 | IdentityRepo, CampaignRepo (impl) | 10 | 4 | OrgService, CampaignService | — |
| **2** | 03 | — | — | — | — | Zitadel config |
| **3** | 04 | CertRepo (in-process CA), DeviceRepo (impl) | 9 | 1 | DeviceService + HTTP /enroll | — (in-process CA) |
| **4** | 05 | MQTTRepo (embedded Mochi) | 1 pure | 1 | (extend DeviceService) | — (embedded broker) |
| **5** | 06 | ReadingRepo (impl) | 7 | 1 | IngestionService | MQTT inline subscription |
| **6** | 07 | — | — | 2 | (extend CampaignService) | — |
| **7** | 10 | ScoreRepo (impl) | 4 | 1 | ScoreService | — |
| **8** | 08 | — | — | 2 | (extend EnrollmentService) | — |
| **9** | 09 | NotificationRepo | 1 | 2 | AdminService | — |

**Cumulative: 9 repos, 33 ops, 14 flows, 7 proto services**

---

## Op Reuse Across Phases

Ops built early get reused later — this validates the decomposition:

| Op | Built in | Reused in |
|----|----------|-----------|
| `UpdateDeviceStatus` | Phase 3 | Phase 8 (renew, revoke), Phase 9 (suspend, reinstate) |
| `IssueCert` | Phase 3 | Phase 8 (renewal) |
| `GetCampaignRules` | Phase 1 | Phase 5 (ingestion validation) |
| `GetCampaignEligibility` | Phase 1 | Phase 4 (enrollment) |
| `GetDeviceCapabilities` | Phase 3 | Phase 4 (enrollment) |
| `QueryReadings` | Phase 5 | Phase 6 (export) |
| `PseudonymizeExport` | Phase 5 | Phase 6 (export) |
| `GetCampaignQuality` | Phase 5 | Phase 6 (dashboard) |
| `QuarantineByWindow` | Phase 5 | Phase 9 (security response) |
| `QueryDevicesByClass` | Phase 3 | Phase 9 (security response) |

**10 of 33 ops are reused across phases.** No customization needed — same op, different flow.
