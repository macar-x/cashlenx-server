# CashLenX Roadmap (Versioned Todo)

## Versioning Policy
- Current version: `v0.1.0`
- Increase minor by `+0.1` for each completed big feature
- Optionally use patch `+0.0.1` for bugfix releases

## Tags
- `#api` API endpoints and contracts
- `#security` authentication, authorization, secrets, rate limits
- `#docs` documentation and developer experience
- `#devops` CI/CD, build, release
- `#observability` logs, metrics, tracing, pprof
- `#data` migrations, backups, schemas
- `#performance` pagination, caching, efficiency
- `#dx` configuration and local setup

## v0.1.0
- [x] CLI commands and base REST API #api #dx
- [x] Cash flow and category CRUD endpoints #api
- [x] Persistence abstraction; MySQL and MongoDB backends #data
- [x] Import/export MVP (Excel) and CLI commands #api #data #cli
- [x] Docker Compose baseline for self-hosted deployments #devops #cloud
- [x] Middleware: CORS and request logging #observability
- [x] Health and version endpoints #api
- [x] Basic docs: README, API, CLI, deployment #docs
- [x] Unit tests in validation, cache, and errors #api

## v0.2.0 — Big Feature: API Contract & OpenAPI
- [ ] Finalize OpenAPI coverage for all endpoints #api #docs
- [ ] Auto-validate requests/responses against schema in dev/test #api
- [ ] Publish HTML docs artifact from OpenAPI #docs #devops
- [ ] Introduce consistent response wrapper `{data,error,meta}` #api
- [ ] Centralize error types and mapping #api
- [ ] Add pagination and filtering to listing endpoints #api #performance
- [ ] Increase unit test coverage in services and mappers #api
- [ ] Align toolchain to Go 1.21 across `go.mod`, Docker, local builds #devops #dx
- [ ] Set up CI (build, test, lint, Docker image) #devops

## v0.3.0 — Big Feature: User Management & Authentication
- [ ] Optional JWT auth middleware protecting mutating endpoints #security #api
- [ ] Minimal user model and env toggles (`AUTH_ENABLED`, `JWT_SECRET`) #security #dx
- [ ] Role-less single-user default; document flows #docs
- [ ] OIDC login support with local user records #security #api
- [ ] Per-user data isolation across storage backends #security #data
- [ ] Admin endpoints for user lifecycle (disable/export/delete) #api #security

## v0.4.0 — Big Feature: Observability
- [ ] Request IDs propagation and structured logging #observability
- [ ] `/metrics` endpoint with Prometheus counters/histograms #observability #devops
- [ ] Enable `pprof` in development #observability

## v0.5.0 — Big Feature: Import/Export Refinements
- [ ] CSV format alongside Excel; unify parsers #api #data
- [ ] Bulk import validation with partial success reporting #api
- [ ] Schema versioning for exports #data #docs
- [ ] User-scoped exports/imports with ownership checks #security #data

## v0.6.0 — Big Feature: Migration Tooling
- [ ] Introduce MySQL migration tooling and track schema changes #data #devops
- [ ] Validate MongoDB indexes at startup and apply scripts #data
- [ ] Backup/restore CLI with progress and validation #data #devops
- [ ] Integration tests via Docker Compose for MongoDB/MySQL #data #devops

## v0.7.0 — Big Feature: Statistics & Insights
- [ ] Summary endpoints per user, per category, per period #api #stats
- [ ] Trend and aggregation endpoints (weekly/monthly) #api #stats
- [ ] Export stats reports (CSV/Excel) #api #stats #data

## v0.8.0 — Big Feature: Performance & Caching
- [ ] Extend category cache; add invalidation on writes #performance
- [ ] Optional read-through cache for recent queries #performance
- [ ] Benchmarks for summaries and mapper queries #performance #devops

## v0.9.0 — Big Feature: Cloud & Self-Hosted
- [ ] Docker Compose profiles for single-tenant and multi-tenant #devops #cloud
- [ ] Helm chart draft for cloud deployments (optional) #devops #cloud
- [ ] Secure defaults for production (CORS, rate limits, secrets) #security #devops

## v1.0.0 — Big Feature: DevOps & Releases
- [ ] GitHub Actions: release pipeline with tagged binaries and images #devops
- [ ] Module caching; reproducible builds #devops
- [ ] Start `CHANGELOG.md` and sync displayed version with code #docs #devops

## Version Sources
- `model/version.go:4` defines the canonical version constant
- `cmd/version.go:12` prints version in CLI `cashlenx-server version`
- `controller/server.go:91` returns version from `GET /api/version`

## Notes
- Each completed big feature increases the minor version by `+0.1`
- When a version’s todo list is fully checked, tag and release that version
- Keep roadmap synchronized with CLI `version` output and `CHANGELOG.md`
