# 1. Executive Summary

## Project Overview

An internal **LLM Proxy Platform** that intermediates all application traffic to **Azure OpenAI** and **OpenAI** (via extensible provider adapters). It exposes **OpenAI-compatible** endpoints for developer speed while centralising **security, governance, cost controls, and observability**. It can be deployed in our tenant and in customer tenants, providing a commercial path with a strong security posture and transparent cost attribution.

**Primary outcomes**

- **Security & Governance:** Centralised authZ/authN, secret isolation, auditability.
- **Cost Management:** Budgets, quotas, model routing, granular usage attribution.
- **Oversight & Analytics:** Per-app/team reporting (tokens, cost, latency, error rates), alerting, export.
- **Developer Velocity:** Drop-in OpenAI API compatibility (baseURL + key).

## Business Context

Organisational AI use is expanding across teams and client environments. Direct SDK → provider patterns create key sprawl, fragmented governance, and opaque spend. A proxy consolidates control without sacrificing developer velocity.

**Value Proposition**

- **Reduce risk:** Eliminate raw provider keys in repos; enforce policy centrally.
- **Control spend:** Apply rate limits/quotas; route workloads to cost-appropriate deployments.
- **Operate at scale:** Consistent interface across environments; single pane of control.
- **Sellable asset:** Provider agnostic deployable package for client environments with private networking and clear governance.

---

# 2. Stakeholders & Benefits

| Stakeholder       | Needs                                              | Impact                                           |
| ----------------- | -------------------------------------------------- | ------------------------------------------------ |
| Developers        | OpenAI-compatible API, minimal changes             | Faster onboarding, less friction                 |
| IT / Security     | Centralised governance, authentication, audit logs | Reduced key sprawl, stronger compliance          |
| Product / Finance | Cost visibility and controls                       | Accurate show back/chargeback, predictable spend |
| Clients           | Deployable in their infrastructure                 | Data residency, control, trust                   |

**Key Benefits & Indicators**

- Centralised access control → **Unmanaged keys ↓** (target: ~80–100% reduction).
- Cost transparency → **100% usage attribution** by app/team/client.
- Developer velocity → **Time-to-first-call**: hours → minutes.
- Flexibility → **Faster rollouts** across Azure tenants via Terraform modules.

---

# 3. Technical Architecture

## 3.1 High-Level

**Hosting & Infra**

- **Azure Container Apps** (Proxy + Admin).
- **Redis** for rate limiting and caching (can be configured to not use).
- **Postgres DB**
- **Private Link** to Azure OpenAI; VNET integration; Private DNS.

**System Components**

- **Data Plane (Proxy):** Authenticates callers, enforces rate limits/quotas/policies, forwards to providers, logs usage.
- **Admin Plane (Go BFF + Nuxt UI):** Key lifecycle, routing rules, policies, dashboards, and **audit logs**.
- **Adapters:**
  - unified api surface rewriting and proxying requests to the correct provider under the hood, completely compatible with existing SDK’s and providers.

**Security Architecture**

- **Caller auth (MVP):** Proxy API keys (scoped, hashed at rest). **Nuxt:** Entra ID JWTs.
- **Upstream auth:** AOAI API keys via **secret references** today; **Managed Identity** on roadmap (no long-lived keys).
- **Audit:** All admin operations (who/what/when) captured and exportable.

## 3.2 Data & Secrets

**Data Entities**

| Entity        | Source   | Store             | Notes                                    |
| ------------- | -------- | ----------------- | ---------------------------------------- |
| API Keys      | Admin UI | Cosmos DB         | Hashed; scopes/roles                     |
| Rate Policies | Admin UI | Cosmos DB         | Config as data                           |
| Usage Metrics | Proxy    | Redis → Cosmos DB | Tokens, cost, latency rollups            |
| Routing Rules | Admin UI | Cosmos DB         | Client model → resource/deployment/alias |

**Secret Handling (aligned with implementation)**

- Store **only `secret_ref`** (e.g., `env:…`, `kv://vault/…`, `db://…`), never plaintext.
- Resolve at runtime via a pluggable **SecretResolver** with TTL cache (rotation-friendly).
- Roadmap: Managed Identity tokens for AOAI (`https://cognitiveservices.azure.com/.default`).

**Integrations**

- Proxy → Azure OpenAI: HTTPS per request (API key now; MI later).
- Admin → Entra ID: OIDC per login (roles: `llm.admin`, `llm.reader`).
- Proxy → Redis: token buckets, key lookups.
- Proxy → Cosmos DB: metadata on cache miss; periodic usage rollups.

---

# 4. Solution Components

## 4.1 Data Plane (Proxy)

- **OpenAI-compatible** endpoints (`/v1/chat/completions`, `/v1/embeddings`, …).
- Per-key **rate limits & quotas** (429 with `Retry-After`), **policy checks** (allowed models/endpoints, context/token caps).
- **Routing**: model → AOAI `{resource, deployment, api-version}`; optional aliasing for OpenAI.
- **Observability**: OpenTelemetry traces; metrics (RPS, latency P50/P95/P99, tokens in/out, cost estimates, error rates).

## 4.2 Admin Plane (API & UI)

- **Entra ID** authentication with role-based access.
- **Key provisioning & lifecycle** (create/disable/rotate; expiry).
- **Routing & policy management** with audit trail.
- **Dashboards & exports**: per-app/team/client usage, cost, model mix, error hotspots.

---

# 5. Acceptance Criteria (MVP)

**Functional**

- Applications call the proxy via standard OpenAI SDKs by changing **baseURL** and using a **proxy key**.
- Proxy correctly maps Azure routes (model → deployment) and supports **streaming**.
- Rate limits/quotas enforced with 429 + `Retry-After`.
- Admins can create/disable keys, edit routes, view usage and **audit logs**.

**Security**

- No plaintext upstream secrets in DB; only `secret_ref`.
- Caller `Authorization` stripped; upstream auth applied by proxy.
- Private networking to AOAI; egress allow-listed.
- Redaction policy applied to request/response logging; configurable retention.

**Non-Functional**

- Added latency ≤ **20ms P95** at 100 RPS (target).
- Availability **≥ 99.9%** (MVP); multi-zone Redis planned for higher SLOs.

**Deployability**

- Terraform modules provision full stack (VNET, Private Link, DNS, Redis, DB, Container Apps) in customer tenants.

---

# 6. Roadmap

- **MVP (Q1):** Proxy with API-key auth; Redis rate limits/quotas; usage logging; `azureopenai` & `openai` adapters; Admin UI (keys, routes, basic reports, audit); Terraform; Cosmos DB.
- **Phase 2:** Entra ID on data plane; **Managed Identity** upstream; Postgres adapter + migrations; richer dashboards/alerts.
- **Phase 3:** Advanced policies (downshift, prompt/context caps, anomaly detection); cost reporting; optional WAF/CDN.

---

## Appendix — Governance, Cost & Analytics

**Security & Governance**

- Single enforcement point for all LLM traffic.
- Central policy & role model; **all changes audited**.
- Secret isolation via **resolver**; no raw keys in repos.

**Cost Controls**

- **Budgets & caps** per app/team; alerts at 80/100%; automatic throttling or model downshift.
- **Attribution**: token → cost conversion per request; rollups by app/team/client; export to FinOps.

**Oversight & Reporting**

- Real-time dashboards (usage, cost, latency, error rates).
- “Top talkers” and error hotspots.
- Exports to SIEM/Power BI/Event Hub.
