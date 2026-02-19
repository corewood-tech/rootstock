ops# ROOTSTOCK — Ops, Flows, Repos Decomposition

> Derived from functional requirements. Transport (UI, RPC, proxy) is not architecture — it's plumbing. The architecture is: **Flows → Ops → Repos**.
>
> - **Flows** orchestrate ordering. They compose ops. They never call repos.
> - **Ops** execute business logic. Each op calls at most ONE repo. This is the clean architecture boundary.
> - **Repos** hide implementation details. Swapping Zitadel for Keycloak, step-ca for Vault PKI, Postgres for CockroachDB — only the repo changes.
>
> Constraint: each op calls **one** repo or **zero** (pure logic). If an op needs two data sources, it's two ops composed by a flow.

---

## Repos (9)

Each repo wraps exactly one implementation detail. The thing behind it can change without touching any op.

| Repo | Wraps | Changes when... |
|------|-------|-----------------|
| **IdentityRepo** | Zitadel | Identity provider changes, user/org CRUD API changes |
| **AuthRepo** | OPA | Policy engine changes, query format changes |
| **CampaignRepo** | Postgres (campaign tables) | Campaign schema evolves, query patterns change |
| **DeviceRepo** | Postgres (device registry) | Device schema evolves, enrollment model changes |
| **CertRepo** | step-ca | CA software changes, signing protocol changes, HSM config changes |
| **ReadingRepo** | Postgres/TimescaleDB (readings) | Storage engine changes, partitioning strategy changes |
| **ScoreRepo** | Postgres (scores/badges) | Gamification data model changes |
| **MQTTRepo** | EMQX API | Broker software changes, publish API changes |
| **NotificationRepo** | Email/push provider | Notification channel changes (email → push → SMS) |

---

## Ops (32 ops across 7 clusters + 2 pure logic)

### Org Ops → IdentityRepo

| # | Op | FR | What it does |
|---|----|----|-------------|
| 1 | CreateOrg | FR-001 | Create org tenant with unique ID, display name, initial admin |
| 2 | NestOrg | FR-002 | Create sub-org within parent (up to 5 levels) |
| 3 | DefineRole | FR-003 | Create org-scoped role with permissions |
| 4 | AssignRole | FR-003 | Assign role to user in specific org |
| 5 | InviteUser | FR-004 | Invite researcher by email, auto-associate on accept |

### Campaign Ops → CampaignRepo

| # | Op | FR | What it does |
|---|----|----|-------------|
| 6 | CreateCampaign | FR-005–008 | Create campaign with parameters (units, range, precision), region (polygon/radius/unbounded), window (start/end UTC), thresholds, eligibility |
| 7 | PublishCampaign | FR-009 | Set status draft → published, make discoverable |
| 8 | ListCampaigns | FR-009, FR-012 | Query with filters: geo proximity, sensor type, status |
| 9 | GetCampaignRules | FR-022 | Return validation criteria (ranges, region, window) for ingestion |
| 10 | GetCampaignEligibility | FR-019 | Return eligibility criteria (device class, tier, sensors, firmware) |

### Device Ops → DeviceRepo

| # | Op | FR | What it does |
|---|----|----|-------------|
| 11 | GenerateEnrollmentCode | FR-013 | 6-8 alphanum, no ambiguous chars, 15-min TTL, one-time use |
| 12 | RedeemEnrollmentCode | FR-013 | Validate code: exists? not expired? not used? Mark used |
| 13 | CreateDevice | FR-016 | Registry entry: device ID, owner, status, class, firmware, tier, sensors, cert serial |
| 14 | GetDevice | — | Read device from registry by ID |
| 15 | GetDeviceCapabilities | FR-019 | Return device class, tier, sensors, firmware for eligibility check |
| 16 | UpdateDeviceStatus | FR-030, FR-031, FR-033 | Change status: pending→active, active→suspended, suspended→active, any→revoked |
| 17 | QueryDevicesByClass | FR-031 | Batch query by device class + firmware version range |
| 18 | EnrollDeviceInCampaign | FR-017, FR-018 | Add campaign association to device record |

### Cert Ops → CertRepo

| # | Op | FR | What it does |
|---|----|----|-------------|
| 19 | IssueCert | FR-014, FR-015, FR-028 | Submit CSR to CA, return signed X.509. CN=device-id, 90-day lifetime |

### Reading Ops → ReadingRepo

| # | Op | FR | What it does |
|---|----|----|-------------|
| 20 | PersistReading | FR-023 | Write reading with full provenance (device ID, timestamp, geo, firmware, cert serial, campaign ID, ingested_at) |
| 21 | QuarantineReading | FR-025 | Flag statistical outlier, mark quarantined, preserve for review |
| 22 | QueryReadings | FR-026 | Read campaign data with provenance for export/dashboard |
| 23 | PseudonymizeExport | FR-027 | Hash device IDs, strip PII, apply spatial resolution policy |
| 24 | QuarantineByWindow | FR-032 | Batch-flag readings from affected devices during vulnerability window |
| 25 | GetCampaignQuality | FR-010 | Aggregated metrics: accepted/rejected counts, rejection reasons, geo/temporal coverage |

### Score Ops → ScoreRepo

| # | Op | FR | What it does |
|---|----|----|-------------|
| 26 | UpdateScore | FR-034 | Recompute: volume, quality rate, consistency, campaign diversity |
| 27 | CheckMilestones | FR-035 | Evaluate badge criteria against current score |
| 28 | AwardBadge | FR-035 | Grant badge: first contribution, 100/1K readings, campaign complete, 30-day streak |
| 29 | GrantSweepstakes | FR-036 | Add entries at score milestones (variable-ratio schedule) |

### Infra Ops (thin wrappers)

| # | Op | Repo | FR | What it does |
|---|----|----- |----|-------------|
| 30 | PushDeviceConfig | MQTTRepo | FR-017 | Retained message on `rootstock/{device-id}/config` |
| 31 | Notify | NotificationRepo | FR-031 | Send notification to scitizen(s) |

### Pure Logic (no repo)

| # | Op | FR | What it does |
|---|----|----|-------------|
| 32 | ValidateReading | FR-022 | Pure function: (reading, rules) → valid/invalid + reason. Schema, range, rate, geo, timestamp |
| 33 | MatchEligibility | FR-019 | Pure function: (device capabilities, campaign criteria) → eligible/reason |

---

## Flows (14)

Each flow composes ops. Flows never call repos directly. A flow can use ops from ANY cluster — no vertical grouping.

### Flow 1: OnboardInstitution
**BUC-01** | FR-001–004

```
CreateOrg → NestOrg → DefineRole → AssignRole → InviteUser
```

| Op | Cluster | Repo |
|----|---------|------|
| CreateOrg | Org | IdentityRepo |
| NestOrg | Org | IdentityRepo |
| DefineRole | Org | IdentityRepo |
| AssignRole | Org | IdentityRepo |
| InviteUser | Org | IdentityRepo |

Single-cluster flow. All identity provider operations.

---

### Flow 2: CreateCampaign
**BUC-02** | FR-005–008

```
CreateCampaign
```

| Op | Cluster | Repo |
|----|---------|------|
| CreateCampaign | Campaign | CampaignRepo |

Single op. Campaign definition (params, region, window, thresholds, eligibility) is one atomic write.

---

### Flow 3: PublishCampaign
**BUC-02** | FR-009

```
PublishCampaign
```

| Op | Cluster | Repo |
|----|---------|------|
| PublishCampaign | Campaign | CampaignRepo |

Single op. Separate from creation because publishing has different authorization and is a distinct business event.

---

### Flow 4: BrowseCampaigns
**BUC-03** | FR-009, FR-012

```
ListCampaigns
```

| Op | Cluster | Repo |
|----|---------|------|
| ListCampaigns | Campaign | CampaignRepo |

Single op. Filtered query (geo, sensor type, status).

---

### Flow 5: RegisterDevice
**BUC-04** | FR-013–016

```
RedeemEnrollmentCode → IssueCert → CreateDevice → UpdateDeviceStatus(pending→active)
```

| Op | Cluster | Repo |
|----|---------|------|
| RedeemEnrollmentCode | Device | DeviceRepo |
| IssueCert | Cert | CertRepo |
| CreateDevice | Device | DeviceRepo |
| UpdateDeviceStatus | Device | DeviceRepo |

Cross-cluster: Device + Cert. Enrollment code generated by a prior step (UI request → GenerateEnrollmentCode). Tier 1 vs Tier 2 is a transport concern — same flow, different caller.

---

### Flow 6: EnrollInCampaign
**BUC-05** | FR-017–019

```
GetDeviceCapabilities + GetCampaignEligibility → MatchEligibility → EnrollDeviceInCampaign → PushDeviceConfig
```

| Op | Cluster | Repo |
|----|---------|------|
| GetDeviceCapabilities | Device | DeviceRepo |
| GetCampaignEligibility | Campaign | CampaignRepo |
| MatchEligibility | Pure | — |
| EnrollDeviceInCampaign | Device | DeviceRepo |
| PushDeviceConfig | Infra | MQTTRepo |

Cross-cluster: Device + Campaign + Infra + Pure. This is the flow that ties the researcher side (campaign criteria) to the scitizen side (device capabilities). The pure MatchEligibility op means neither side knows about the other.

---

### Flow 7: IngestReading
**BUC-06** | FR-020–025

```
GetCampaignRules → ValidateReading → PersistReading → [QuarantineReading]
```

| Op | Cluster | Repo |
|----|---------|------|
| GetCampaignRules | Campaign | CampaignRepo |
| ValidateReading | Pure | — |
| PersistReading | Reading | ReadingRepo |
| QuarantineReading | Reading | ReadingRepo |

Cross-cluster: Campaign + Pure + Reading. Hot path (10K/sec, P99 < 200ms). Campaign rules can be cached — GetCampaignRules doesn't need to hit the DB on every reading. ValidateReading is pure logic: no I/O, no repo, deterministic.

**Note:** mTLS authentication and OPA authorization happen at the handler/transport layer (global auth), before this flow is entered. The flow assumes the caller is already authenticated and authorized.

---

### Flow 8: ExportData
**BUC-07** | FR-026–027

```
QueryReadings → PseudonymizeExport
```

| Op | Cluster | Repo |
|----|---------|------|
| QueryReadings | Reading | ReadingRepo |
| PseudonymizeExport | Reading | ReadingRepo |

Single-cluster. Both ops hit ReadingRepo. PseudonymizeExport applies hashing and spatial resolution before output.

---

### Flow 9: RenewCert
**BUC-08** | FR-028–029

```
IssueCert → UpdateDeviceStatus (update cert serial + expiry)
```

| Op | Cluster | Repo |
|----|---------|------|
| IssueCert | Cert | CertRepo |
| UpdateDeviceStatus | Device | DeviceRepo |

Cross-cluster: Cert + Device. Same IssueCert op reused from RegisterDevice. Grace period logic (expired ≤7 days → allow renewal; >7 days → reject) lives in the handler or flow guard, not in the op.

---

### Flow 10: RevokeDevice
**BUC-08** | FR-030

```
UpdateDeviceStatus(→revoked)
```

| Op | Cluster | Repo |
|----|---------|------|
| UpdateDeviceStatus | Device | DeviceRepo |

Single op. OPA picks up status change on next bundle refresh (≤30s).

---

### Flow 11: SecurityResponse
**BUC-09** | FR-031–032

```
QueryDevicesByClass → UpdateDeviceStatus(bulk→suspended) → QuarantineByWindow → Notify
```

| Op | Cluster | Repo |
|----|---------|------|
| QueryDevicesByClass | Device | DeviceRepo |
| UpdateDeviceStatus | Device | DeviceRepo |
| QuarantineByWindow | Reading | ReadingRepo |
| Notify | Infra | NotificationRepo |

Cross-cluster: Device + Reading + Infra. UpdateDeviceStatus is the same op used everywhere — here it runs in batch mode.

---

### Flow 12: ReinstateDevice
**BUC-09** | FR-033

```
UpdateDeviceStatus(suspended→active)
```

| Op | Cluster | Repo |
|----|---------|------|
| UpdateDeviceStatus | Device | DeviceRepo |

Single op. Campaign enrollments preserved.

---

### Flow 13: UpdateContributionScore
**BUC-10** | FR-034–036

```
UpdateScore → CheckMilestones → AwardBadge → GrantSweepstakes
```

| Op | Cluster | Repo |
|----|---------|------|
| UpdateScore | Score | ScoreRepo |
| CheckMilestones | Score | ScoreRepo |
| AwardBadge | Score | ScoreRepo |
| GrantSweepstakes | Score | ScoreRepo |

Single-cluster. Triggered async after reading acceptance. Score input (reading stats) passed in by the caller, not fetched by the op — keeps it within ScoreRepo.

---

### Flow 14: CampaignDashboard
**BUC-02** | FR-010

```
GetCampaignQuality
```

| Op | Cluster | Repo |
|----|---------|------|
| GetCampaignQuality | Reading | ReadingRepo |

Single op. Aggregated from readings, not from campaign table.

---

## Op Reuse Matrix

Which ops appear in multiple flows? This drives architecture decisions — highly reused ops justify their existence as separate units.

| Op | Flows that use it | Reuse count |
|----|-------------------|-------------|
| **UpdateDeviceStatus** | RegisterDevice, RenewCert, RevokeDevice, SecurityResponse, ReinstateDevice | **5** |
| **IssueCert** | RegisterDevice, RenewCert | **2** |
| **ListCampaigns** | BrowseCampaigns (+ UI campaign selection in EnrollInCampaign) | **2** |
| **GetCampaignRules** | IngestReading (+ potentially CampaignDashboard) | **1–2** |
| All other ops | 1 flow each | **1** |

**UpdateDeviceStatus** is the most reused op in the system. It appears in 5 flows across 3 BUCs. It's the workhorse of device lifecycle management. This validates it as a standalone op rather than being inlined into each flow.

**IssueCert** appears in both initial enrollment and renewal — same crypto operation, different flow context.

---

## Grouping Analysis

### Question: Should any op clusters merge?

**Candidates for merge:**

| Merge | Pro | Con | Verdict |
|-------|-----|-----|---------|
| Cert → Device | Only 1 op. IssueCert always called in device flows | CertRepo (step-ca) changes independently from DeviceRepo (Postgres). Merging means CA software swap touches device code | **No** — different volatility |
| MQTT → Device | PushDeviceConfig only used in device enrollment flows | MQTTRepo (EMQX API) changes independently. Broker swap touches device code | **No** — different volatility |
| Notification → Score | Notifications triggered by engagement events | Notification channel changes (email→push) shouldn't touch score logic | **No** — different volatility |
| Org → Campaign | Both serve researcher persona | Adding org hierarchy levels doesn't affect campaign parameters. Different change reasons | **No** — different volatility |
| Reading → Campaign | IngestReading needs campaign rules | Campaign rules are READ by readings, not owned. CampaignRepo owns rules. ReadingRepo owns readings. Clean boundary | **No** — different data ownership |

**Verdict: Keep all clusters separate.** The 1-op clusters (Cert, MQTT, Notification) are thin by design — they wrap a single external dependency with a single operation. That's correct. A repo with one op is still a legitimate boundary if the thing it wraps changes independently.

---

## Data Map

Each repo owns specific tables. No table is shared across repos.

### CampaignRepo (Postgres)

| Table | Key Fields | Notes |
|-------|-----------|-------|
| campaigns | id, org_id, status (draft/published/closed), window_start, window_end, created_by | Core campaign record |
| campaign_parameters | campaign_id, name, unit, min_range, max_range, precision | One row per parameter per campaign |
| campaign_regions | campaign_id, geometry (PostGIS) | Polygon/radius/admin boundary. Nullable = unbounded |
| campaign_eligibility | campaign_id, device_class, tier, required_sensors[], firmware_min | What devices can participate |

### DeviceRepo (Postgres)

| Table | Key Fields | Notes |
|-------|-----------|-------|
| devices | id (=cert CN), owner_id, status, class, firmware_version, tier, sensors[], cert_serial | Device registry — single source of truth |
| enrollment_codes | code, device_id, expires_at, used | 15-min TTL, one-time use |
| device_campaigns | device_id, campaign_id, enrolled_at | Many-to-many: device ↔ campaign |

### ReadingRepo (Postgres/TimescaleDB)

| Table | Key Fields | Notes |
|-------|-----------|-------|
| readings | id, device_id, campaign_id, value, timestamp, geolocation, firmware_version, cert_serial, ingested_at, status (accepted/quarantined), quarantine_reason | Provenance-complete. Append-only |

### ScoreRepo (Postgres)

| Table | Key Fields | Notes |
|-------|-----------|-------|
| scores | scitizen_id, volume, quality_rate, consistency, diversity, total, updated_at | Recomputed periodically |
| badges | scitizen_id, badge_type, awarded_at | Immutable once awarded |
| sweepstakes_entries | scitizen_id, entries, milestone_trigger, granted_at | Auditable, tamper-evident |

### IdentityRepo (Zitadel)

No local tables — org hierarchy, users, roles, grants all managed by Zitadel. IdentityRepo wraps Zitadel's gRPC/REST API.

### CertRepo (step-ca)

No local tables — cert issuance handled by step-ca. CertRepo wraps step-ca's signing API via PKCS#11/HSM.

### AuthRepo (OPA)

No local tables — policy evaluation is in-memory. AuthRepo wraps OPA's REST query API. Bundle data synced FROM DeviceRepo.

---

## Visual Components (UI views derived from flows)

UI views map 1:1 to flows. Each view calls exactly one flow through the handler.

### Researcher Views

| View | Flow | Key Interactions |
|------|------|-----------------|
| Org Management | OnboardInstitution | Create org, nest sub-orgs, define roles, invite members |
| Campaign Creator | CreateCampaign | Form: parameters, region (map picker), window (date range), thresholds, eligibility |
| Campaign List | BrowseCampaigns | Table with status, window, region, enrolled device count |
| Campaign Dashboard | CampaignDashboard | Accepted/rejected counts, quality metrics, geo distribution, temporal coverage |
| Data Export | ExportData | Format selection, date range, filters → download |

### Scitizen Views

| View | Flow | Key Interactions |
|------|------|-----------------|
| Campaign Browse | BrowseCampaigns | Filtered by location + device compatibility |
| Device Enrollment | RegisterDevice | Generate code → enter on device → device gets cert → active |
| Campaign Enrollment | EnrollInCampaign | Select device, select campaign, eligibility check, enroll |
| Contribution Dashboard | (reads from ScoreRepo) | Score, badges, sweepstakes entries, active campaigns, reading counts |

### Admin Views

| View | Flow | Key Interactions |
|------|------|-----------------|
| Security Response | SecurityResponse | Select device class + firmware range → bulk suspend → notify |
| Device Management | RevokeDevice, ReinstateDevice | Status changes, individual or batch |

---

## Flow Complexity Summary

| Flow | Ops | Clusters crossed | Repos touched | Complexity |
|------|-----|------------------|---------------|-----------|
| OnboardInstitution | 5 | 1 | 1 | Low |
| CreateCampaign | 1 | 1 | 1 | Trivial |
| PublishCampaign | 1 | 1 | 1 | Trivial |
| BrowseCampaigns | 1 | 1 | 1 | Trivial |
| RegisterDevice | 4 | 2 | 2 | Medium |
| EnrollInCampaign | 5 | 3 + pure | 3 | **High** |
| IngestReading | 3–4 | 2 + pure | 2 | Medium (hot path) |
| ExportData | 2 | 1 | 1 | Low |
| RenewCert | 2 | 2 | 2 | Low |
| RevokeDevice | 1 | 1 | 1 | Trivial |
| SecurityResponse | 4 | 3 | 3 | **High** |
| ReinstateDevice | 1 | 1 | 1 | Trivial |
| UpdateContributionScore | 4 | 1 | 1 | Low |
| CampaignDashboard | 1 | 1 | 1 | Trivial |

**Two high-complexity flows:** EnrollInCampaign (bridges researcher world to scitizen world) and SecurityResponse (coordinates across device, reading, and notification concerns). These are the flows that justify the architecture — without flows as an orchestration layer, these would force cross-cluster coupling.
