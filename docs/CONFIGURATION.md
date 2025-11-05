# Shanraq Configuration Guide

Shanraq loads configuration from `config.yaml` (or another path passed to `-config`) and transparently merges environment variables prefixed with `SHANRAQ_`. This guide documents every section, the effective defaults, and practical overrides.

## File Layout

```yaml
environment: development
server:
  address: ":8080"
  read_timeout: 15s
  write_timeout: 15s
  idle_timeout: 60s
database:
  url: postgres://postgres:postgres@127.0.0.1:5432/shanraq?sslmode=disable
  max_conns: 10
  min_conns: 1
  max_conn_lifetime: 30m
  max_conn_idle_time: 5m
  health_check_period: 30s
telemetry:
  enable_metrics: true
  metrics_path: /metrics
logging:
  level: info
  mode: production
auth:
  token_secret: replace-me-now
  token_ttl: 15m
```

### Environment

| Key | Description |
| --- | ----------- |
| `environment` | Displayed in the operator console and influences safe defaults. Production runs must override `auth.token_secret`. |

### Server

| Key | Description | Notes |
| --- | ----------- | ----- |
| `server.address` | Listen address/port. | Set to `0.0.0.0:8080` when running in containers. |
| `server.read_timeout` | Read deadline per request. | Go duration string. |
| `server.write_timeout` | Write deadline per request. | Adjust for long streaming responses. |
| `server.idle_timeout` | Keep-alive timeout. | Lower in front of load balancers. |

### Database

| Key | Description | Notes |
| --- | ----------- | ----- |
| `database.url` | pgx connection string. | Supports query parameters like `sslmode=require`. |
| `database.max_conns` / `min_conns` | Connection pool sizing. | Ensure `max_conns` < Postgres `max_connections`. |
| `database.max_conn_lifetime` | TTL per connection. | Helps recycle stale connections. |
| `database.max_conn_idle_time` | Idle timeout before release. | Lower for burst workloads. |
| `database.health_check_period` | Interval between `ping` checks. | `0` disables proactive checks. |

### Telemetry

| Key | Description |
| --- | ----------- |
| `telemetry.enable_metrics` | Enables `/metrics` (Prometheus). |
| `telemetry.metrics_path` | Custom path for the handler. |
| `telemetry.tracing.enabled` | Enables OpenTelemetry spans via OTLP exporter. |
| `telemetry.tracing.endpoint` | OTLP HTTP endpoint (e.g. `collector:4318`). |
| `telemetry.tracing.sample_ratio` | Sampling fraction between 0 and 1. |
| `telemetry.tracing.service_name` | Overrides the `service.name` resource attribute. |

### Logging

| Key | Description |
| --- | ----------- |
| `logging.level` | Log verbosity (`debug`, `info`, `warn`, `error`). |
| `logging.mode` | `development` enables console-friendly encoder; `production` emits JSON. |

### Authentication

| Key | Description | Notes |
| --- | ----------- | ----- |
| `auth.token_secret` | HMAC secret for JWT tokens. | Production mode panics when this remains the default. Use a 32+ byte random string. |
| `auth.token_ttl` | Access token lifetime. | Refresh tokens outlive this (30 days by default). |
| `auth.mfa.totp.enabled` | Enables TOTP-based MFA challenges during sign-in. | When enabled, users must verify a code from an authenticator app. |
| `auth.mfa.totp.issuer` | Issuer label displayed inside authenticator apps. | Defaults to `Shanraq`. |

## Environment Variable Overrides

Every key above may be overridden by prefixing `SHANRAQ_` and using upper snake case. Examples:

| Variable | Overrides |
| -------- | --------- |
| `SHANRAQ_SERVER_ADDRESS` | `server.address` |
| `SHANRAQ_DATABASE_URL` | `database.url` |
| `SHANRAQ_LOGGING_MODE` | `logging.mode` |
| `SHANRAQ_AUTH_MFA_TOTP_ENABLED` | `auth.mfa.totp.enabled` |
| `SHANRAQ_AUTH_MFA_TOTP_ISSUER` | `auth.mfa.totp.issuer` |
| `SHANRAQ_AUTH_TOKEN_SECRET` | `auth.token_secret` |
| `SHANRAQ_AUTH_TOKEN_TTL` | `auth.token_ttl` |

> **Tip:** Create an `.env` file alongside `config.yaml` for local development. `internal/config` auto-loads `.env` when present.

## Sample `.env` Snippet

```
SHANRAQ_DATABASE_URL=postgres://shanraq:shanraq@db:5432/shanraq?sslmode=disable
SHANRAQ_AUTH_TOKEN_SECRET=$(openssl rand -hex 32)
SHANRAQ_LOGGING_MODE=development
SHANRAQ_SERVER_ADDRESS=:8080
```

## Module-Specific Settings

- **Auth**: The new RBAC model stores roles in `auth_roles` and `auth_user_roles`. Use migrations or seed scripts to create additional roles, then assign them via SQL or bespoke handlers.
- **API Keys**: Customer credentials live in `auth_api_keys`. Keys are hashed at rest; expose creation endpoints only behind `auth.RequireRoles`. Demo seeds provision `sk_demo_operator_token` for the operator account—rotate it outside development.
- **Jobs**: Worker counts live in `cmd/app/main.go`. Expose an environment variable (e.g. `SHANRAQ_JOBS_WORKERS`) if you need runtime overrides.
- **Web UI**: Carousel and docs pull copy from `framework_about`. Update via SQL seeds or admin tooling.
- **Notifier**: Configure `notifications.smtp` to enable e-mail (host, port, username, password, from). Leaving host or from empty keeps delivery disabled while still logging reset links.

## Troubleshooting

- Cannot connect to Postgres? Set `database.url` to `postgres://...@127.0.0.1/...` instead of `localhost` if the resolver blocks IPv6.
- Production panic mentioning the token secret: set `SHANRAQ_AUTH_TOKEN_SECRET` or edit `config.yaml`.
- Metrics 404: ensure `telemetry.enable_metrics=true` and your reverse proxy allows the path defined in `telemetry.metrics_path`.
