# ROOTSTOCK — Backend Development Roadmap

> Derived from the ops/flows/repos decomposition. Ordered by BUC dependency chain and implementation prerequisites.

---

## What Exists

| Layer | What's built | Status |
|-------|-------------|--------|
| **Config** | Multi-source loading (YAML, env, flags) | Complete |
| **Globals** | O11y (OTel-backed), Events (DBOS-backed), Auth (OPA-backed) | Complete |
| **Repos** | AuthorizationRepo (OPA), ObservabilityRepo (OTel), EventsRepo (DBOS), SQL pool (pgx) | Complete |
| **Interceptors** | AuthorizationInterceptor, BinaryOnlyInterceptor, OTel auto-instrumentation | Complete |
| **Handlers** | HealthService.Check | Complete |
| **Proto** | HealthService only | Scaffolded |
| **Infra** | Postgres (app + Zitadel), Zitadel + Login v2, Caddy, OTel Collector, Prometheus, Tempo, Loki, Grafana | Running |
| **Flows** | — | Not started |
| **Ops** | — | Not started |
| **Business repos** | — | Not started |
| **Migrations** | — | Not started |

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

### 3.1 — Infrastructure: step-ca

New compose file `compose-ca.yml` (added to `COMPOSE_FILES`):
- step-ca container with SoftHSM (dev)
- Root CA + Issuing CA hierarchy
- ACME or custom enrollment endpoint

### 3.2 — CertRepo

New repo at `repo/cert/`. Wraps step-ca's signing API:
- `IssueCert(ctx, CSR) → Certificate` — submit CSR, get signed cert back

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

### 3.7 — Proto: EnrollmentService

```protobuf
service EnrollmentService {
  rpc GenerateCode(GenerateCodeRequest) returns (GenerateCodeResponse);
  rpc Enroll(EnrollRequest) returns (EnrollResponse);          // Tier 1 + Tier 2
  rpc Renew(RenewRequest) returns (RenewResponse);             // Phase 8
}
```

### 3.8 — Handler: EnrollmentServiceHandler

Maps `Enroll` RPC → RegisterDevice flow. `GenerateCode` is a separate RPC called from the UI before device enrollment begins.

### 3.9 — E2E Test

Scitizen generates enrollment code → simulated device calls `/enroll` with code + CSR → receives cert → device appears in registry as active.

**Deliverable:** Devices can be registered with real certificates. Foundation for ingestion.

---

## Phase 4: Campaign Enrollment

**Goal:** A scitizen can enroll a registered device in a published campaign.

**BUCs:** BUC-05 | **FRs:** FR-017–019 (3 Must)

### 4.1 — Infrastructure: MQTT Broker

New compose file `compose-mqtt.yml` (added to `COMPOSE_FILES`):
- EMQX container with mTLS enabled
- Auth webhook pointing to web-server
- Topic ACL structure: `rootstock/{device-id}/data/{campaign-id}`

### 4.2 — MQTTRepo

New repo at `repo/mqtt/`. Wraps EMQX management API:
- `PushDeviceConfig(ctx, deviceID, config)` — publish retained message to `rootstock/{device-id}/config`

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

### 4.6 — OPA Policy: Device Authorization

- Device status check (active → allow, suspended/revoked → deny)
- Campaign enrollment check (device enrolled in campaign → allow publish to topic)
- Topic ACL (device can only publish to `rootstock/{own-device-id}/*`)

### 4.7 — Auth Service (MQTT Webhook)

Thin glue between EMQX auth webhook format and OPA query format. Can be a separate small Go binary or an endpoint on the web-server.

### 4.8 — E2E Test

Scitizen enrolls device in campaign → device appears in campaign's enrolled devices → config retained message published.

**Deliverable:** Devices enrolled in campaigns. MQTT infrastructure ready. OPA enforcing device authorization.

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

### 5.5 — MQTT Consumer

Subscribes to `rootstock/+/data/+` topics. For each message:
1. Extract device ID and campaign ID from topic
2. Deserialize reading (protobuf)
3. Call IngestReading flow
4. ACK or NACK

This is a new component — either part of web-server or a separate consumer process. Decision: separate process keeps the hot path isolated from the RPC server.

### 5.6 — HTTP/2 Ingestion Endpoint (Dual Protocol)

```protobuf
service IngestionService {
  rpc SubmitReading(SubmitReadingRequest) returns (SubmitReadingResponse);
}
```

Same IngestReading flow behind both MQTT and HTTP/2 paths (OP-002).

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
IssueCert → UpdateDeviceStatus (update cert serial + expiry)
```

Both ops already exist from Phase 3. New flow, no new ops.

### 8.2 — RevokeDevice Flow

```
UpdateDeviceStatus(→revoked)
```

Single op, already exists. OPA denies within 30s.

### 8.3 — Grace Period Logic

In the EnrollmentServiceHandler (or a flow guard):
- Expired ≤7 days → allow Renew RPC only
- Expired >7 days → reject, must re-enroll

### 8.4 — E2E Test

Device with expiring cert → calls `/renew` → gets new cert. Revoked device → denied all actions.

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
| **3** | 04 | CertRepo, DeviceRepo (impl) | 9 | 1 | EnrollmentService | step-ca |
| **4** | 05 | MQTTRepo | 1 pure | 1 | DeviceService | EMQX |
| **5** | 06 | ReadingRepo (impl) | 7 | 1 | IngestionService | MQTT consumer |
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
