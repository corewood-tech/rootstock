# ROOTSTOCK

## by Corewood

---

<!-- HERO -->

In horticulture, a rootstock refers to the established root of a fruit plant that can be grafted with limbs from other trees. Because the established root proves hearty and reliable, the branches grow from its steady supply of nutrients and solid grounding.

This is the missing rootstock for LLM-assisted code generation.

It's not a framework. It's not a specific technology.

It's the basic way of thinking about a software system that you can graft onto, extend, and evolve for your needs.

**Regrowth without uprooting.**

[Use Rootstock →](#architecture) · [Hire the people who built it →](https://corewood.io/schedule-meeting/)

---

<!-- SECTION: THE PROBLEM -->

## LLMs Write Code Fast. Making It Scale Is the Hard Part.

You've seen it. Maybe you've lived it.

A vibe-coded project starts fast, works well, then starts toppling under its own weight. Duplicated functions everywhere. Multiple dependencies managing the same job. The LLM gets confused, finds conflicting patterns, and stops helping.

Even disciplined approaches break down. A well-structured Python app with models, a repository pattern, and clean modules still hits the wall — because the context window explodes with every change.

This problem is structural, not skill-based. The architecture has to be designed for how LLMs actually work.

**Corewood spent ~11 months learning this the hard way** — building production software, yelling at the LLM, and refining the patterns that actually hold up. Rootstock is everything we learned, open-sourced.

---

<!-- SECTION: WHAT WE BUILT WITH THIS -->

## We Didn't Theorize This. We Shipped It.

Corewood built [LandScope](https://landscope.earth) with CEO [Mitch Rawlyk](https://mitch.earth) — a terrain intelligence platform used by land professionals across the world. Enterprise authentication. SSO/SAML integration. [OPA](https://www.openpolicyagent.org/)-driven permissions. [OpenTelemetry](https://opentelemetry.io/) from day one. 291,000+ acres mapped and counting.

We also built LLM inference engines, complex Postgres wire protocol interceptions, and a bunch of websites.

Then we asked: what if we open-sourced the architecture that made all of it possible?

> Most engineering decisions optimize for authoring software — trying to make writing code easier. Corewood optimizes for operations, because over time, you save money.

**Rootstock is the answer. Use it free. Or [hire us](https://corewood.io/schedule-meeting/) to build yours.**

---

<!-- SECTION: TWO LESSONS -->

## Two Lessons That Changed Everything

### 1. Manage Context Windows

Context windows are the bottleneck. LLMs effectively work at smaller scale, but as the application grows it buckles under its own weight. The LLMs get confused, find multiple patterns to follow, and ultimately fail to help your project grow.

Rootstock's architecture keeps every task small and self-contained. The LLM never needs to reason about the whole system.

### 2. Follow Strict Patterns

Every choice you give the LLM is a risk to the stability of the project. Don't give it more choices than absolutely necessary.

Rootstock uses identical file names, identical structure, identical conventions across every module. The LLM copies the nearest neighbor — and it works every time.

**Want to apply these lessons to your codebase?** [Talk to a Corewood engineer →](https://corewood.io/schedule-meeting/)

---

<!-- SECTION: THREE PRINCIPLES -->

## The Three Principles

### Use Performant Languages

LLM-generated code won't be perfect. Give it a runtime where imperfect code still performs.

In independent benchmarks, [Go](https://go.dev/) executes 3-20x faster than Python and uses 1.5-3x less memory across typical workloads. The performance gap is wide enough that even suboptimal Go code routinely outperforms well-written Python. That headroom is what makes LLM-assisted coding viable at scale.

### Definite, Reproducible Patterns

LLMs are pattern completion machines. Give them patterns that repeat identically — same file names, same structure, same conventions across every module. The LLM copies the nearest neighbor. The more consistent the pattern, the more reliable the output.

### Keep Context Small Per-Task

The context window is the hard limit. The architecture must ensure that any given task only requires a small slice of the codebase. One module, one layer, two type files. That's all the LLM needs to see.

---

<!-- SECTION: ARCHITECTURE -->

## The Architecture

Four layers. One direction. Every module follows the same shape.

### Handler → Auth → Flow → Ops → Repo

**Handler** — Decouples protocol and policy. Translates wire format to application types. Runs auth before anything else.

**Auth** — A sub-flow at the edge. Identity and authorization are resolved through globals that delegate to repos. Vendors are implementation details hidden behind interfaces. Unauthorized requests are dropped here — they never reach business logic. Auth-derived data (user ID, roles) is passed explicitly into the flow request. Clean contract. No side-channel dependencies.

**Flow** — Orchestrates. Calls one or more ops, converts between layer types. Thin by design.

**Ops** — Does the things. Business rules, validation, decisions. The fat layer.

**Repo** — Handles integration points. Any external boundary — database, identity provider, object store, observability vendor. Hides provider-specific details behind a clean interface. Swap the vendor, keep the contract.

### The Module Convention

Every module at every layer has exactly two type files:

- **receive_types** — what comes in
- **emit_types** — what comes out

That's the contract. That's the entire surface area. If you know the inputs and outputs, you know the module. The LLM knows the module.

### Directional Imports

Imports only go inward. Handler imports flow. Flow imports ops. Ops imports repo. Never the reverse.

The innermost layers are the most stable. The outermost layer absorbs all protocol changes without rippling inward. Circular dependencies are structurally impossible.

### Globals and Vendor Obfuscation

Cross-cutting concerns — observability, events, auth — live in `global/`. Each global is a thin singleton accessor that delegates to a repo. The repo wraps the vendor. The global knows nothing about the vendor. Swap [OpenTelemetry](https://opentelemetry.io/) for Datadog, swap [Zitadel](https://zitadel.com/) for Auth0 — one repo changes, everything else stays the same.

This is [volatility-based decomposition](https://www.informit.com/articles/article.aspx?p=2995357&seqNum=2): things that change independently are isolated behind separate boundaries. [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) governs the direction: dependencies always point inward.

**This is how Corewood structures every client engagement.** The architecture works whether you're building it yourself or [hiring us to build it for you](https://corewood.io/schedule-meeting/).

---

<!-- SECTION: THE STACK -->

## The Stack

| Layer | Choice | Why |
|-------|--------|-----|
| Backend | [Go](https://go.dev/) | Performance headroom, concurrency, ops-focused ecosystem, CGo for C interop |
| Frontend | [Svelte](https://svelte.dev/) + [SvelteKit](https://svelte.dev/docs/kit) (static bundle) | Compiles to static assets. No server runtime. |
| API Protocol | [ConnectRPC](https://connectrpc.com/) | Typed contracts from protos, codegen both sides, works over HTTP/1.1 |
| Database | [PostgreSQL](https://www.postgresql.org/) | Handles almost everything. Boring. Proven. |
| Object Store | [MinIO](https://min.io/) | S3-compatible, runs as a container, full dev/prod parity |
| Identity | [Zitadel](https://zitadel.com/) (SessionAPI only) | Isolated behind repo pattern — provider doesn't leak |
| Authorization | [OPA](https://www.openpolicyagent.org/) | Policies as data, evaluated in-process, recompiled on state changes |
| Observability | [OpenTelemetry](https://opentelemetry.io/) | Traces → [Tempo](https://grafana.com/oss/tempo/), metrics → [Prometheus](https://prometheus.io/), logs → [Loki](https://grafana.com/oss/loki/), dashboards → [Grafana](https://grafana.com/oss/grafana/) |
| Workflow Engine | [DBOS](https://dbos.dev/) | Durable execution for long-running processes |
| Reverse Proxy | [Caddy](https://caddyserver.com/) | Automatic TLS, config-driven routing |
| Containers | [Podman](https://podman.io/) Compose | Dev parity with production. Dependencies live in the container, not on your machine. |
| Config | [koanf](https://github.com/knadh/koanf) | Defaults → YAML → env vars → CLI flags. Resolved once at startup. |

Every choice reduces the number of things the LLM has to guess about.

**Corewood didn't just pick these tools — we've operated them in production.** [Let us help you pick yours →](https://corewood.io/schedule-meeting/)

---

<!-- SECTION: THE PROOF -->

## The Proof: A Citizen Science Platform That Actually Works

Rootstock isn't theoretical. We're building a complete citizen science IoT platform — called **Rootstock** — to demonstrate every piece of the architecture under real conditions. The full requirements specification follows the [Volere template](https://github.com/corewood-io/rootstock/tree/main/spec), traced to a [Dgraph](https://dgraph.io/) knowledge graph for auditability.

### The Problem

Traditional scientific field data collection is expensive and structurally inadequate. Personnel costs consume 50-80% of grant budgets. Existing biodiversity datasets cover less than 7% of Earth's surface at 5km resolution. 75% of citizen science projects produce zero peer-reviewed publications — because they collect data without a specific research question.

Meanwhile, 250,000+ personal weather stations (Weather Underground), 30,000+ air quality sensors (PurpleAir), and a growing market of water and soil probes are already deployed at massive scale. These are distributed mini-labs sitting idle. PurpleAir data has already been used in peer-reviewed wildfire research — the data quality is there. The coordination isn't.

No platform exists that lets researchers define multi-parameter campaigns and match them to citizen-owned IoT devices.

### The Platform

**Researchers** create data campaigns — structured requests specifying what sensor data they need, where, when, and to what quality standard. Campaign mechanics are borrowed directly from marketing: goal, audience, timeframe, measurement. The campaign model inverts the typical citizen science flow: instead of "collect data and hope researchers use it," researchers define what they need first.

**Scitizens** — our term for citizen scientists — discover campaigns relevant to their location and equipment, enroll their IoT devices, and contribute readings. Devices submit data automatically after enrollment, fundamentally changing the retention equation. The critical conversion is enrollment, not repeated contribution. Engagement is incentivized through gamification, recognition, and sweepstakes — not per-unit compensation. Prospect Theory confirms that lottery-style incentives are more motivating per dollar spent than fixed payments.

**Devices** are untrusted. Every sensor reading crosses a trust boundary. Over 50% of IoT devices have critical exploitable vulnerabilities. The platform authenticates via mTLS with a self-operated certificate authority, authorizes via [OPA](https://www.openpolicyagent.org/) against a device registry, and validates every reading against campaign-defined quality thresholds — schema, range, rate, geolocation, anomaly detection. Quarantined readings are never deleted.

**Data** is exported with full provenance metadata, quality metrics, and campaign context — meeting [FAIR principles](https://www.go-fair.org/fair-principles/) so the data is actually publishable. Contributor identity is separated from observation data. Raw contributor locations never appear in public-facing datasets — 4 spatiotemporal data points uniquely identify 95% of individuals.

### What This Exercises

This isn't a demo app. It's a platform with:

- **Multi-tenant hierarchical authorization** — institutions contain departments contain labs, roles are org-scoped, [OPA](https://www.openpolicyagent.org/) enforces everything
- **IoT device lifecycle management** — enrollment codes, CSR-based certificate issuance, 90-day cert rotation, automated renewal, bulk revocation on firmware vulnerabilities
- **mTLS device authentication** — two-tier CA hierarchy, device registry as the source of truth, no CRL/OCSP needed because the platform controls the relying party
- **Campaign-driven data ingestion** — schema validation, range checks, cross-device correlation, anomaly flagging
- **Cross-org participation** — scitizens have no institutional affiliation with the requesting researcher, which is the normal case
- **Privacy architecture** — GDPR-compliant consent model, granular per-campaign permissions, configurable spatial resolution, contributor identity never in exported datasets
- **Full observability from day one** — [OpenTelemetry](https://opentelemetry.io/) traces, metrics, and structured logs flowing through a collector to [Prometheus](https://prometheus.io/), [Tempo](https://grafana.com/oss/tempo/), [Loki](https://grafana.com/oss/loki/), and [Grafana](https://grafana.com/oss/grafana/) dashboards

Every architectural decision stress-tested by a genuinely complex problem.

**Not a todo app. Not a demo.** A platform that solves the last-mile problem between researchers and the data they need.

---

<!-- SECTION: DUAL CTA -->

## Two Ways Forward

### Use Rootstock

The architecture, the patterns, the stack — all open source under a BSD 3-Clause license. Clone it, learn from it, graft your own project onto it.

`make up` starts the full stack: web server with hot reload, [PostgreSQL](https://www.postgresql.org/), [Zitadel](https://zitadel.com/), [OPA](https://www.openpolicyagent.org/), [OpenTelemetry](https://opentelemetry.io/), [Prometheus](https://prometheus.io/), [Grafana](https://grafana.com/oss/grafana/), [Caddy](https://caddyserver.com/). Everything runs in containers. The only prerequisite is [Podman](https://podman.io/).

Rootstock gives you the foundation. You build the branches.

[GitHub: Rootstock →](https://github.com/corewood-io/rootstock)

### Hire Corewood

You have a real problem. You need software that handles user data, credit cards, enterprise auth, and doesn't fall apart at 2am.

We've built this before. We'll build it with you.

Corewood brings the same architecture, the same patterns, the same production-first mindset — applied to your problem, your domain, your timeline.

[Schedule a Consultation →](https://corewood.io/schedule-meeting/) · [Contact Us →](https://corewood.io/contact/)

---

<!-- SECTION: PHILOSOPHY -->

## The Corewood Philosophy

> Most engineering decisions optimize for authoring software — trying to make writing code easier. Corewood optimizes for operations, because over time, you save money.

The question that drives every decision: **What kind of system would I want to pay a team to operate?**

That's what Rootstock answers. And that's what Corewood builds.

---

<!-- FOOTER -->

[Corewood](https://corewood.io) · [LandScope](https://landscope.earth) · [Mitch Rawlyk](https://mitch.earth)

© 2026 Corewood LLC
