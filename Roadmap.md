# **1. Executive Summary**

## Project Overview

We are delivering an LLM Proxy Service that sits between internal developers/applications and Azure OpenAI instances. The proxy provides a secure, consistent, and scalable way for teams to use large language models without directly managing Azure resources.

## Business Context

The organization is expanding use of AI services across multiple teams and client environments. A proxy approach allows centralised governance, rate limiting per application, cost control, and future-proof flexibility, while still enabling developer speed and compatibility with existing SDKs.

**1.1 Key Objectives**

1. Provide a unified API gateway for Azure OpenAI access.
2. Ensure secure authentication and role-based control.
3. Support rate limiting, quotas, and usage tracking for accountability.
4. Remain compatible with existing OpenAI SDKs (drop-in replacement).
5. Deploy flexibly to multiple client environments with minimal overhead.
6. Enable future extension to managed identity / Entra-based authentication.

# **2. Business Context**

**2.1 Target Audience**

| **Stakeholder Group** | **Needs**                                          | **Impact**                                      |
| --------------------- | -------------------------------------------------- | ----------------------------------------------- |
| Developers            | Simple, drop-in compatible API for OpenAI SDKs     | Faster onboarding, less friction                |
| IT / Security         | Centralized governance, authentication, audit logs | Stronger compliance, minimized key sprawl       |
| Product Teams         | Reliable, scalable access to AI models             | Consistent performance & reduced downtime       |
| Clients               | Deployable solution in their infrastructure        | Confidence in security & control over resources |

**2.2 Business Benefits**

| **Benefit**                | **Metric**                       | **Expected Impact**                    |
| -------------------------- | -------------------------------- | -------------------------------------- |
| Centralized access control | Number of unmanaged keys reduced | Lower security risk                    |
| Cost transparency          | Per-tenant usage reports         | Improved cost allocation & forecasting |
| Developer velocity         | Time to first use of models      | Hours → Minutes                        |
| Flexibility                | Deployable across Azure tenants  | Faster client rollouts                 |

# **3. Technical Architecture**

**3.1 High-Level Architecture**

Infrastructure Overview

- Azure Container Apps (two roles: Proxy & Admin)
- Redis for rate limiting and caching
- Cosmos DB (initial store) with option for Postgres later
- Managed Identities for secure connection to Azure OpenAI

System Components

- Data Plane (Proxy): Handles SDK requests, applies authentication, rate limits, and forwards to AOAI.
- Admin Plane (API + UI): Provides secure UI (Nuxt + Go BFF) for managing keys, roles, policies, and reporting usage.
- Storage: Cosmos DB for metadata, Redis for performance-critical operations.

Security Architecture

- Data plane: API key authentication (MVP) with plan to add Entra ID JWTs.
- Admin plane: Entra ID login with app roles (llm.admin, llm.reader).
- Managed Identity replaces raw keys for upstream AOAI access.
- Audit logging for all admin actions.

**3.2 Detailed Architecture**

**3.2.1 Data Architecture**

| **Data Entity** | **Source** | **Destination**   | **Transform Logic**                   |
| --------------- | ---------- | ----------------- | ------------------------------------- |
| API Keys        | Admin UI   | Cosmos DB         | Hashing, role assignment              |
| Rate Policies   | Admin UI   | Cosmos DB         | Stored as configs                     |
| Usage Metrics   | Proxy      | Redis → Cosmos DB | Rollup & aggregation                  |
| Routing Rules   | Admin UI   | Cosmos DB         | Model alias → AOAI deployment mapping |

**3.2.2 Integrations**

| **Integration Point** | **Type**     | **Protocol**  | **Frequency**           | **Data Flow**                              |
| --------------------- | ------------ | ------------- | ----------------------- | ------------------------------------------ |
| Proxy → Azure OpenAI  | Service call | HTTPS         | Per request             | Forwarded requests (Managed Identity auth) |
| Admin → Entra ID      | Auth         | OAuth2 / OIDC | Per login               | User identity + roles                      |
| Proxy → Redis         | Cache        | TLS           | Per request             | Token buckets, key lookup                  |
| Proxy → Cosmos DB     | Store        | TLS           | On cache miss / rollups | Key metadata, usage                        |

**3.2.3 System Dependencies**

- Azure Container Apps
- Azure Cosmos DB (or Postgres future option)
- Azure Redis Cache
- Azure Entra ID (for admin auth)
- Azure OpenAI Services

# **4. Solution Components**

**4.1 Data Plane (Proxy)**

- Compatible with OpenAI SDKs.
- API key authentication (future: Entra ID).
- Rate limiting, quotas, and usage logging.
- Managed Identity to call Azure OpenAI.

**4.2 Admin Plane (API & UI)**

- Nuxt frontend served from Go API (BFF model).
- Entra ID authentication with role-based access.
- Key provisioning and lifecycle management.
- Usage dashboards and audit logs.

# **5. Acceptance Criteria**

Developers can call the proxy using standard OpenAI SDKs with only baseURL + proxy key.

- Proxy enforces per-key rate limits and quotas (429 with Retry-After).
- Responses (including streaming) are passed through correctly.
- Admins can log in via Entra, manage API keys, and view usage.
- No raw AOAI keys are exposed; Managed Identity is used.
- Audit logs capture all admin actions.
- Solution can be deployed in multiple client Azure environments via Terraform.

# **6. Roadmap**

Draft Roadmap

- MVP (Q1): API key–based proxy + admin UI with key mgmt, usage logs, audit. Cosmos + Redis.
- Phase 2: Add Entra ID authentication for data plane (JWT). Postgres adapter + migrations.
- Phase 3: Advanced routing policies, cost reporting, optional CDN/WAF in front.
- Phase 4: Multi-cloud portability (non-Azure clients).
