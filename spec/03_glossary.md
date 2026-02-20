# ROOTSTOCK by Corewood

## Requirements Specification — Section 3: Naming Conventions and Terminology

> We acknowledge that this document uses material from the Volere Requirements Specification Template, copyright © 1995–2019 the Atlantic Systems Guild Limited.

> Names are very important. They invoke meanings that, if carefully defined, can save hours of explanations. Attention to names at this stage helps highlight misunderstandings. This glossary is used and extended throughout the project.

---

## 3a. Glossary

### Platform & Project

| Term | Definition |
|------|-----------|
| **Rootstock** | The open-source scientific data collection platform. The product being specified in this document. |
| **Corewood** | The organization building and maintaining Rootstock. Acts as platform steward and open-source project owner. |

### Stakeholders & Actors

| Term | Definition |
|------|-----------|
| **Research Institution** | An organization that employs or sponsors researchers. May be a university, government agency, private company, or independent lab. An institution may contain sub-institutions (e.g., a university contains departments, a company contains divisions). |
| **Researcher** | An individual who defines data needs and analyzes collected data on behalf of a research institution. The primary consumer of campaign data. |
| **Citizen Scientist** | A voluntary participant who contributes sensor data to campaigns using personal IoT devices. Not compensated directly — incentivized through gamification, recognition, and sweepstakes. See also: **Scitizen**. |
| **Scitizen** | Preferred term for Citizen Scientist within the Rootstock platform. A portmanteau of "science" and "citizen" that emphasizes the participant's active role: these are citizens doing science, not just volunteers. Scitizens have invested in equipment, have domain interest, and are functionally producing scientific measurements. The term reflects that their contribution is not charity — it is participation in a shared scientific enterprise, mediated by the platform. |
| **Oversight Body** | Any entity with governance, ethical, financial, or political authority over a research institution or researcher. Includes grant boards, ethics committees, IRBs, and funding agencies. |
| **Contributor** | General term for anyone providing data to the platform. Currently synonymous with Scitizen, but left as a separate term in case future data sources (e.g., automated stations, institutional sensors) are onboarded. |

### Core Domain Concepts

| Term | Definition |
|------|-----------|
| **Campaign** | A structured request for data, created by a researcher. A campaign defines what data is needed, where, when, to what quality standard, and under what constraints. Campaigns are the central organizing unit of the platform. |
| **Campaign Window** | The time period during which a campaign actively accepts data submissions. Has a start and end. |
| **Campaign Region** | The geographic boundary within which a campaign accepts data. May be a polygon, radius, administrative boundary, or unbounded. |
| **Campaign Parameters** | The set of measurable quantities a campaign requests (e.g., temperature, humidity, water pH, altitude, particulate matter). Each parameter includes acceptable ranges, units, and precision requirements. |
| **Campaign Enrollment** | The association between a device and a campaign, authorizing that device to submit data readings for that campaign. A device may be enrolled in multiple campaigns simultaneously. Enrollment is tracked in the device registry and enforced by OPA policy. |
| **Data Reading** | A single observation submitted by a contributor's device for a specific campaign. Contains a value, timestamp, geolocation, device identifier, and campaign reference. The atomic unit of collected data. |
| **Data Submission** | A batch of one or more data readings transmitted from a device to the platform in a single request. |
| **Observation** | Synonym for Data Reading. Preferred in research-facing contexts. |

### Data Quality & Trust

| Term | Definition |
|------|-----------|
| **Data Provenance** | The verifiable chain of origin for a data reading: which device produced it, where, when, under what firmware version, and through what transmission path it reached the platform. mTLS and the device registry together establish provenance — the certificate proves device identity, the registry provides context. |
| **Trust Boundary** | The point at which data enters the platform from an external source. All device data crosses a trust boundary and is treated as untrusted input until validated. The ingestion layer is the trust boundary. |
| **Quality Threshold** | A campaign-defined minimum standard that a data reading must meet to be accepted. May include precision, accuracy, sampling rate, or calibration requirements. |
| **Validation** | The process of checking a data reading against campaign-defined quality thresholds before acceptance. Validation is not silent — readings that fail validation are rejected with a reason. Includes schema validation, range checks, rate limiting, and anomaly flagging. |
| **Calibration** | The process by which a device's measurement accuracy is verified or adjusted against a known reference. Calibration status may affect whether a device is eligible for a campaign. |
| **Outlier** | A data reading that falls outside expected statistical bounds for a campaign. Outliers are flagged, not silently discarded. Disposition (accept, reject, quarantine) is determined by campaign policy. |
| **Quarantine** | A holding state for data readings or devices under review. Quarantined readings are not visible to researchers but are not deleted. Quarantined devices cannot submit new data until review is complete. |

### Devices & Capability Tiers

| Term | Definition |
|------|-----------|
| **IoT Device** | Any consumer or prosumer electronic device capable of collecting sensor data and transmitting it to the platform. Includes smartphones, weather stations, water quality monitors, air quality sensors, GPS trackers, and similar equipment. |
| **Device Tier** | A classification of device capability that determines the enrollment method. Three tiers are defined; only Tier 1 and Tier 2 are in MVP scope. |
| **Tier 1 Device** | A device that can run an HTTPS client, has persistent storage, and has a UI or display. Examples: smartphones, Raspberry Pi, modern weather stations. Enrolls directly with the platform. |
| **Tier 2 Device** | A network-capable but constrained device with no rich UI and limited TLS capability. Examples: ESP32-based sensors, LoRa gateways. Enrolls via a companion app that proxies the enrollment flow. |
| **Tier 3 Device** | A very constrained or legacy device with no TLS capability. Examples: cheap 8-bit MCU sensors, BLE-only devices. Requires a gateway architecture. Out of MVP scope — separate project. |
| **Companion App** | A mobile application that acts as an enrollment proxy for Tier 2 devices. Communicates with the constrained device over BLE or local WiFi, requests a CSR from the device, submits the enrollment on its behalf, and pushes the issued certificate back to the device. The companion app never sees the device's private key. |
| **Device Class** | A category of devices grouped by manufacturer, model, and capability. Used for campaign eligibility restrictions, bulk operations (e.g., mass suspension of a vulnerable firmware version), and enrollment templates. |
| **Device Registry** | The central platform database (PostgreSQL) that stores all device identity, ownership, status, class, capability, enrollment history, and campaign associations. The single source of truth for device state. OPA syncs from this registry to make authorization decisions. |
| **Device Profile** | A record of a registered device's characteristics within the device registry: manufacturer, model, firmware version, sensor capabilities, calibration history, and reliability track record on the platform. |
| **Device Status** | The current operational state of a device on the platform. Possible values: `pending` (enrolled but not yet confirmed), `active` (operational), `suspended` (temporarily disabled, e.g., due to firmware vulnerability), `revoked` (permanently disabled). Status is checked by OPA on every action. |
| **Ownership Transfer** | The process of reassigning a device from one citizen scientist to another. Transfers update the registry binding — the device's certificate is unaffected. If the previous owner is uncooperative, the device is revoked and re-enrolled by the new owner. |

### Device Registration & Enrollment

| Term | Definition |
|------|-----------|
| **Device Registration** | The full process by which a citizen scientist associates a specific IoT device with their account, the device receives a certificate, and the device becomes authorized to communicate with the platform. Encompasses enrollment code generation, CSR submission, certificate issuance, and registry entry creation. |
| **Enrollment Code** | A short-lived (15-minute TTL), one-time-use, human-friendly token (6-8 alphanumeric characters, no ambiguous characters like 0/O or 1/l) that ties a pending device registration in the web UI to the physical device performing enrollment. The bootstrap trust anchor. |
| **Enrollment Service** | The platform component that handles device registration. Exposes two HTTP endpoints: `POST /enroll` (new device, authenticated by enrollment code) and `POST /renew` (existing device, authenticated by current mTLS cert). Validates enrollment codes, verifies CSRs, coordinates with the CA, and updates the device registry. |
| **CSR (Certificate Signing Request)** | A standard cryptographic request generated by a device containing its public key, submitted to the enrollment service. The device generates its own keypair — the private key never leaves the device. The CSR format is standard X.509; every TLS library can produce one. |
| **Direct Enrollment** | The enrollment path for Tier 1 devices. The device itself calls the enrollment service directly over HTTPS, presenting the enrollment code and a locally-generated CSR. |
| **Proxy Enrollment** | The enrollment path for Tier 2 devices. The companion app communicates with the constrained device over BLE/local WiFi, requests a CSR from the device (which generates its own keypair), submits the CSR to the enrollment service on the device's behalf, and pushes the returned certificate back to the device. |

### Certificate Authority & mTLS

| Term | Definition |
|------|-----------|
| **mTLS (Mutual TLS)** | The authentication mechanism for all device-to-platform communication. Both the device and the platform present certificates during the TLS handshake, establishing mutual identity verification. TLS does authentication only — authorization is handled by OPA. |
| **Rootstock CA** | The self-operated certificate authority that issues device certificates. Two-tier hierarchy: an offline root CA and online issuing CAs. Self-operated because external CAs (Let's Encrypt, AWS Private CA) either don't issue client certs at IoT scale or create cost/proprietary dependencies. |
| **Root CA** | The top of the certificate chain. Air-gapped or HSM-protected. 10-year certificate, EC P-384. Signs only issuing CA certificates. Touched rarely (approximately twice a year). |
| **Issuing CA** | An online CA signed by the root CA. Separate issuing CAs exist for devices and platform services. 2-year certificate, HSM-backed. Signs device and service certificates on request. |
| **Device Certificate** | An X.509 certificate issued to a device during enrollment. Minimal content: `CN=<device-id>`, issuer, serial, validity window, key usage. No SANs, no OUs, no device metadata in the cert — all mutable metadata lives in the device registry. 90-day lifetime. |
| **Certificate Renewal** | The process by which a device with a valid certificate obtains a new certificate before expiration. The device presents its current cert via mTLS to `POST /renew` with a new CSR. Automated, triggered at day 60 of a 90-day cert. |
| **Certificate Revocation** | Handled by the device registry and OPA, not by traditional CRL or OCSP mechanisms. When a device is revoked, its status is updated in the registry. OPA denies all subsequent actions on the next policy evaluation cycle. The platform controls the relying party, so traditional revocation infrastructure is unnecessary. |
| **Grace Period** | A 7-day window after certificate expiration during which a device can still reach the renewal endpoint (and only the renewal endpoint) to obtain a new certificate. Devices expired beyond the grace period must re-enroll. |
| **HSM (Hardware Security Module)** | Hardware that protects CA private keys. SoftHSM for development, YubiHSM 2 or cloud HSM for production. Required for issuing CA keys — not optional at scale. Accessed via PKCS#11 interface. |
| **In-process CA** | Certificate signing implemented directly in Go using `crypto/x509` stdlib, wrapped by CertRepo. No external CA software (step-ca, Vault PKI). The CA key is exclusively owned by a single goroutine via channel-based concurrency. Enrollment orchestration (code validation, CSR verification, registry integration) is custom platform code. Replaces the earlier step-ca plan. |

### Authentication & Authorization

| Term | Definition |
|------|-----------|
| **Authn/Authz Split** | The architectural principle that authentication (proving identity) and authorization (granting access) are separate concerns handled by separate systems. TLS handles authn (device presents cert, platform verifies CA signature, extracts device ID). OPA handles authz (device ID is checked against registry for status, campaign enrollment, topic ACLs). |
| **OPA (Open Policy Agent)** | The policy engine that makes all authorization decisions on the platform. Evaluates Rego policies against device registry data. Queried by the authorization interceptor (Connect RPC) and in-process MQTT auth hook. Data synced from device registry via bundles (≤30-second refresh). |
| **Rego** | The policy language used by OPA. All authorization rules (device status checks, topic ACLs, campaign enrollment verification) are expressed as Rego policies. Testable and auditable. |
| **OPA Bundle** | A periodically-refreshed data package that syncs device registry state to OPA for policy evaluation. Refresh interval determines revocation latency — a 30-second bundle refresh means a revoked device is denied within 30 seconds. |
| **Tenant** | An entity with an isolated data and authorization boundary on the platform. Research institutions are tenants. |
| **Organization** | A structural entity that can contain users, sub-organizations, and campaigns. Organizations nest — a university is an organization containing department organizations. |
| **Organization Hierarchy** | The tree of nested organizations. Authorization rules may inherit down the hierarchy (a university admin can see department campaigns) or be scoped to a specific level. |
| **Role** | A named set of permissions assigned to a user within the context of an organization. Roles are not global — a user may hold different roles in different organizations. |
| **Permission** | A specific authorized action (e.g., create campaign, view readings, manage devices, invite members). Permissions are granted via roles, never directly. |
| **Membership** | The association between a user and an organization, along with the role(s) they hold in that organization. A user may have memberships in multiple organizations. |
| **Cross-Org Participation** | The ability for a citizen scientist to contribute data to campaigns belonging to organizations they are not a member of. This is the normal case — citizen scientists typically have no institutional affiliation with the requesting researcher. |

### Communication Protocols

| Term | Definition |
|------|-----------|
| **MQTT 5.0** | The primary device-to-platform communication protocol. Chosen for IoT suitability: low overhead, intermittent connectivity support (QoS levels, session resumption, last will), tiny packet size, broad client library availability (ESP32 through smartphones). MQTT 5.0 adds user properties (metadata), shared subscriptions (load balancing), and flow control. All MQTT connections use mTLS. |
| **MQTT Broker** | The server component that receives, routes, and delivers MQTT messages. Embedded in-process using Mochi MQTT (Go, `github.com/mochi-mqtt/server/v2`). Same rationale as in-process CA — reduces infrastructure dependency. Auth handled by an in-process Go hook (mTLS cert verification + topic ACL), not an external webhook. InlineClient enables server-side publish (config push) and subscribe (telemetry consumption). |
| **MQTT Topic** | The hierarchical address to which devices publish data and from which the platform consumes it. Structure: `rootstock/{device-id}/data/{campaign-id}` for sensor data, `rootstock/{device-id}/status` for health, `rootstock/{device-id}/config` for platform-pushed configuration, `rootstock/{device-id}/cmd` for commands. |
| **Topic ACL** | Access control on MQTT topics, enforced by the in-process MQTT auth hook (`OnACLCheck`). A device can only publish to topics under its own device ID (`rootstock/{its-own-device-id}/*`). The inline client has unrestricted access. Prevents a compromised device from injecting data as another device. |
| **HTTP/2 Gateway** | The secondary data ingestion path for Tier 1 devices that prefer request/response communication (smartphones, Raspberry Pis). Also serves all human-facing APIs, enrollment endpoints, and service-to-service communication. All connections use mTLS. |
| **QoS (Quality of Service)** | MQTT delivery guarantee levels. QoS 0: at most once (fire and forget). QoS 1: at least once (acknowledged). QoS 2: exactly once (four-part handshake). Campaign requirements may mandate minimum QoS levels. |
| **Retained Message** | An MQTT feature where the broker stores the last message on a topic and delivers it to new subscribers immediately. Used for device configuration — a reconnecting device receives its current config without waiting for a push. |
| **Last Will and Testament (LWT)** | An MQTT feature where a device registers a message to be published if it disconnects unexpectedly. Used for device status monitoring — the platform knows immediately when a device drops. |

### Ingestion & Data Pipeline

| Term | Definition |
|------|-----------|
| **Ingestion** | The platform process that receives device data (via MQTT or HTTP/2), validates it against campaign parameters, associates it with a campaign, and persists it. The ingestion layer is the trust boundary — all data entering here is untrusted until validated. |
| **Input Validation** | The first stage of ingestion. Checks performed on every incoming data submission: schema validation (correct fields and types), range checks (values within physical possibility), rate limiting (per device, to detect malfunctioning or compromised devices), and anomaly flagging (geolocation outside campaign region, timestamps outside campaign window). Failures are rejected with explicit reasons. |
| **Telemetry** | The raw data stream transmitted by an IoT device, before validation or association with a campaign. What the device sends. |
| **Data Routing** | The process of associating validated telemetry with one or more campaigns based on the topic structure (`rootstock/{device-id}/data/{campaign-id}`) and the device's campaign enrollments. |

### Failure Modes & Diagnostics

| Term | Definition |
|------|-----------|
| **Clock Drift** | The tendency of cheap IoT device clocks to lose accuracy over time. Affects TLS handshake (cert validation requires accurate time) and data provenance (timestamps may be wrong). Mitigated by NTP on boot, relaxed `notBefore` validation (48-hour tolerance), and per-device drift monitoring. |
| **Connection Failure Diagnostics** | Structured log entries produced on every failed device connection attempt. Includes device IP, presented cert (if any), failure reason, and timestamp. Failure reasons are specific and actionable: "Certificate expired on {date}", "Device {id} status: suspended", "Device {id} not enrolled in campaign {campaign}". Surfaced to citizen scientists as human-readable messages on their dashboard. |
| **Bulk Revocation** | The process of suspending or revoking all devices of a specific device class in response to a discovered vulnerability. Registry batch update → OPA denies on next bundle refresh → mass notification to affected citizen scientists. |

### Incentives & Engagement

| Term | Definition |
|------|-----------|
| **Gamification** | The system of recognition, achievements, and incentives that encourages ongoing citizen scientist participation. |
| **Sweepstakes** | A reward mechanism in which active contributors are entered into drawings for prizes or experiences, rather than receiving per-unit compensation. The primary cost-control mechanism for the incentive model. |
| **Recognition** | Public or profile-visible acknowledgment of a citizen scientist's contributions (e.g., contribution counts, badges, leaderboards, campaign acknowledgments in publications). |
| **Contribution Score** | A computed metric reflecting a citizen scientist's overall participation: volume, quality, consistency, and diversity of contributions. Used for recognition and sweepstakes eligibility. |

### Platform Operations

| Term | Definition |
|------|-----------|
| **Campaign Management** | The set of capabilities for researchers to create, configure, publish, monitor, and close campaigns. The CMS-like element of the platform. |
| **Data Export** | The process by which a researcher extracts collected campaign data for analysis in external tools. Export formats and access controls are campaign-scoped. |
| **Platform Health** | The observable operational state of the platform: ingestion throughput, validation error rates, API latency, storage consumption, cert expiration rates, OPA denial rates, and system availability. |

---

## 3b. Data Dictionary

> The data dictionary will be populated as functional requirements are defined. It will specify the data inputs and outputs for each atomic requirement — field names, types, ranges, units, and constraints.

### Data Reading

| Data Element | Type | Description | Constraints |
|-------------|------|-------------|-------------|
| `reading.value` | numeric | The measured value of a single observation | Units and range defined per campaign parameter |
| `reading.timestamp` | datetime (UTC) | When the measurement was taken by the device | Must be within campaign window |
| `reading.geolocation` | lat/lon + altitude | Where the measurement was taken | Must be within campaign region if defined |
| `reading.device_id` | identifier | Reference to the device that produced the reading | Must match a registered, active device |
| `reading.campaign_id` | identifier | Reference to the campaign this reading is submitted against | Must reference an active campaign; device must be enrolled |

### Campaign

| Data Element | Type | Description | Constraints |
|-------------|------|-------------|-------------|
| `campaign.window_start` | datetime (UTC) | When the campaign begins accepting readings | Required |
| `campaign.window_end` | datetime (UTC) | When the campaign stops accepting readings | Required; must be after window_start |
| `campaign.region` | geometry | Geographic boundary for the campaign | Optional; unbounded if not set |

### Device Registry

| Data Element | Type | Description | Constraints |
|-------------|------|-------------|-------------|
| `device.id` | identifier | Unique platform identifier for the device | Generated at enrollment; matches certificate CN |
| `device.owner_id` | identifier | Reference to the citizen scientist who owns this device | Required; updated on ownership transfer |
| `device.status` | enum | Current operational state | One of: `pending`, `active`, `suspended`, `revoked` |
| `device.class` | string | Device class identifier (manufacturer + model) | Used for campaign eligibility and bulk operations |
| `device.firmware_version` | semver string | Current firmware version reported by device | Used for provenance, eligibility, and bulk vulnerability response |
| `device.tier` | enum | Device capability tier | One of: `tier_1`, `tier_2` |
| `device.campaigns` | identifier[] | List of campaigns this device is enrolled in | Enforced by OPA on every publish |

### Device Certificate

| Data Element | Type | Description | Constraints |
|-------------|------|-------------|-------------|
| `cert.serial` | string | Unique certificate serial number | Issued by CA; tracked in registry |
| `cert.subject_cn` | string | Certificate common name | Must equal `device.id` |
| `cert.not_before` | datetime (UTC) | Certificate validity start | Set at issuance |
| `cert.not_after` | datetime (UTC) | Certificate validity end | 90 days from issuance |
| `cert.issued_at` | datetime (UTC) | When the certificate was issued | Tracked for audit |

### Enrollment

| Data Element | Type | Description | Constraints |
|-------------|------|-------------|-------------|
| `enrollment_code.value` | string | Human-friendly one-time token | 6-8 alphanumeric, no ambiguous chars (0/O, 1/l) |
| `enrollment_code.ttl` | duration | Time-to-live from generation | 15 minutes |
| `enrollment_code.device_id` | identifier | The pending device this code is bound to | One code per pending device |
| `enrollment_code.used` | boolean | Whether the code has been consumed | One-time use; rejected after first use |

---

*Next: [Section 4 — Relevant Facts and Assumptions](./04-relevant-facts.md)*