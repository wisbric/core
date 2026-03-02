# core

Shared Go libraries used by NightOwl, BookOwl, and TicketOwl.

## What’s Inside

- `pkg/auth` — authentication (session cookies, OIDC JWT, PATs, API keys) and RBAC helpers
- `pkg/authadapter` — shared auth adapter base types for implementing `auth.Storage` in each service
- `pkg/config` — shared config parsing (env tags + defaults via `BaseConfig`)
- `pkg/httpserver` — common HTTP middleware, health/ready endpoints, metrics
- `pkg/platform` — Postgres/Redis helpers and migrations
- `pkg/tenant` — multi-tenant context + search_path handling
- `pkg/telemetry` — Prometheus + OpenTelemetry helpers
- `pkg/version` — build metadata

## Notes

- Dev header auth (`X-Tenant-Slug`) is **opt-in**: services must pass `allowDevHeader=true` (usually when `DEV_MODE=true`).
- Session cookies are signed with a shared secret; use the same secret across services for SSO.
- If `CORS_ALLOWED_ORIGINS` includes `*`, credentials are disabled to prevent unsafe configuration.

## Development

```bash
# From repo root
cd core

go test ./...
```
