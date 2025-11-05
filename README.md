<p align="center">
  <img src="web/static/brand/logo_dark.svg" alt="Shanraq Logo" width="64">
</p>

# Shanraq.org Framework

Shanraq is a batteries-included Go 1.23+ application framework focused on modern backend practices: typed configuration, PostgreSQL-first data access, structured logging, composable modules, and first-class observability.

## Highlights

- **Go 1.23+ ready**: Module targets Go `1.23` and pins toolchain `go1.24.1`, so you build on stable releases while staying future-proof.
- **PostgreSQL via pgxpool**: Sane defaults, lifecycle management, and trace hooks wired into Zap.
- **Declarative configuration**: Environment-aware configuration loader with `.env` support, file overrides, and typed structs.
- **Composable modules**: `shanraq.Module` contracts let teams build reusable features that can hook into init/start/routes independently.
- **Operational insights**: Built-in health, readiness, and Prometheus telemetry modules provide production-grade endpoints out of the box.
- **Domain modules included**: Authentication endpoints, background job queue/workers, and a Bootstrap 5.3 dashboard ship as ready-to-use modules.
- **Secure auth flows**: Refresh token rotation, password reset templates, and role-based access control (RBAC) tables with helper middleware.
- **SaaS ready**: Optional API key module and tenant-aware job queues let you scope workloads per customer out of the box.
- **Zero-downtime migrations**: Embedded Goose migrations run on startup, so schema upgrades stay versioned and automated.

## Project Layout

```
cmd/app/              # Reference application entry point
configs/              # Example configuration files
internal/config       # Configuration loader & validation
internal/db           # PostgreSQL connection factory & tracing
internal/httpserver   # Chi-powered HTTP server wrapper
internal/logging      # Zap logger builder
pkg/shanraq           # Public runtime & module contracts
pkg/modules/auth      # Email/password auth REST module
pkg/modules/apikeys   # Tenant API keys and middleware
pkg/modules/jobs      # Postgres-backed job queue & workers
pkg/modules/migrations  # Embedded Goose migrations
pkg/modules/notifier  # SMTP-backed email notifications
pkg/modules/webui     # Bootstrap 5.3 dashboard module
pkg/modules/...       # Additional modules (health, telemetry, etc.)
web/                  # html/template renderer, landing carousel, static assets
```

## Getting Started

Pick the lane that suits your workflow:

- **Native:** install Go ≥ 1.23 (the toolchain directive will fetch `go1.24.1`) and PostgreSQL, copy the sample config, then `go run ./cmd/app -config config.yaml`.
- **Docker:** `docker compose up --build` provisions a distroless image and local Postgres.

Quick steps:

1. `cp configs/config.example.yaml config.yaml`
2. `export SHANRAQ_AUTH_TOKEN_SECRET=$(openssl rand -hex 32)` (or set in `.env`)
3. Launch via Go or Docker as above.
4. Explore `/`, `/console`, `/docs`, `/jobs`, `/metrics`.

### Database Migrations

The migrations module runs automatically during startup. If you prefer to apply migrations before launching the HTTP server, run:

```bash
go run ./cmd/app -config config.yaml
```

The binary will create the target database if it does not yet exist, execute every embedded Goose migration under `pkg/modules/migrations/sql`, and then block serving requests.

### Docker Workflow

Prefer containers? The repo ships with a production-style `Dockerfile` and a `docker-compose.yml` suited for local development:

```bash
docker compose up --build
```

This command launches:

- `db`: PostgreSQL 16 with a persistent volume.
- `app`: the Shanraq binary built inside a minimal distroless image, configured via env vars to talk to the `db` service.

Expose ports `8080` (app) and `5432` (Postgres) on your host to interact with the stack just like a native run.

## Configuration Reference

| Key | Description | Default |
| --- | ----------- | ------- |
| `environment` | Logical environment label | `development` |
| `server.address` | HTTP bind address | `:8080` |
| `server.read_timeout`/`write_timeout`/`idle_timeout` | Server timeouts | `15s/15s/60s` |
| `database.url` | PostgreSQL DSN (pgx format) | `postgres://postgres:postgres@localhost:5432/shanraq?sslmode=disable` |
| `database.max_conns` | Max pooled connections | `10` |
| `telemetry.enable_metrics` | Toggle Prometheus handler | `true` |
| `telemetry.metrics_path` | Metrics HTTP path | `/metrics` |
| `telemetry.tracing.enabled` | Enable OTLP tracing | `false` |
| `telemetry.tracing.endpoint` | OTLP collector endpoint | *(empty)* |
| `telemetry.tracing.sample_ratio` | Sampling fraction (0-1) | `0.1` |
| `telemetry.tracing.service_name` | OpenTelemetry `service.name` | `shanraq-app` |
| `logging.level` | Zap level (`debug/info/warn/error`) | `info` |
| `logging.mode` | `development` or `production` encoder | `production` |
| `auth.token_secret` | HMAC secret for JWT tokens | `replace-me-now` |
| `auth.token_ttl` | Token lifetime | `15m` |

Durations accept Go duration syntax (e.g. `30s`, `5m`).

### Environment Variables

| Variable | Purpose |
| -------- | ------- |
| `SHANRAQ_DATABASE_URL` | Overrides `database.url` with a full pgx DSN. |
| `SHANRAQ_AUTH_TOKEN_SECRET` | JWT signing key (32+ bytes recommended). |
| `SHANRAQ_AUTH_TOKEN_TTL` | Overrides access token lifetime (e.g. `45m`). |
| `SHANRAQ_SERVER_ADDRESS` | Bind address for the HTTP server. |
| `SHANRAQ_TELEMETRY_ENABLE_METRICS` | Toggle Prometheus endpoint (`true`/`false`). |
| `SHANRAQ_LOGGING_LEVEL` | Set log level (`debug`, `info`, etc.). |
| `SHANRAQ_LOGGING_MODE` | Overrides encoder (`development` or `production`). |
| `SHANRAQ_JOBS_WORKERS` | Optional override for worker count when wired in `cmd/app/main.go`. |
| `SHANRAQ_NOTIFICATIONS_SMTP_HOST` | SMTP host for outbound email (accompanied by `_PORT`, `_USERNAME`, `_PASSWORD`, `_FROM`). |
| `SHANRAQ_AUTH_MFA_TOTP_ENABLED` | Enables TOTP multi-factor authentication. |
| `SHANRAQ_AUTH_MFA_TOTP_ISSUER` | Overrides the issuer label used in authenticator apps. |

### Notifications

Configure outbound email via the `notifications.smtp` section:

```yaml
notifications:
  smtp:
    host: smtp.example.com
    port: 587
    username: api-key
    password: super-secret
    from: no-reply@example.com
```

Leaving `host` or `from` empty disables e-mail delivery; password reset links continue to be logged to the console for local development.

### UI Shortcuts (scan on mobile)

<p>
  <img src="https://api.qrserver.com/v1/create-qr-code/?size=128x128&data=http%3A%2F%2Flocalhost%3A8080%2F" alt="Home QR" style="margin-right: 1rem;"/>
  <img src="https://api.qrserver.com/v1/create-qr-code/?size=128x128&data=http%3A%2F%2Flocalhost%3A8080%2Fconsole" alt="Console QR" style="margin-right: 1rem;"/>
  <img src="https://api.qrserver.com/v1/create-qr-code/?size=128x128&data=http%3A%2F%2Flocalhost%3A8080%2Fdocs" alt="Docs QR"/>
</p>

## Building Modules

Create modules by implementing one or more of the optional interfaces:

```go
type WidgetModule struct{}

func (WidgetModule) Name() string { return "widget" }

func (WidgetModule) Init(ctx context.Context, rt *shanraq.Runtime) error {
    // One-time setup before HTTP server starts (migrations, caches, ...)
    return nil
}

func (WidgetModule) Routes(r chi.Router) {
    r.Get("/widgets", ...)
}

func (WidgetModule) Start(ctx context.Context, rt *shanraq.Runtime) error {
    // Long-running workers, schedulers, etc.
    <-ctx.Done()
    return ctx.Err()
}
```

Register modules in `cmd/app/main.go` (or your own binary):

```go
app := shanraq.New(cfg)
app.Register(widget.New())
app.Register(telemetry.New())
app.Register(health.New())
```

The runtime injects a shared `*pgxpool.Pool`, `*zap.Logger`, and `chi.Router` instance, so modules stay cohesive yet decoupled.

### Enabling or Disabling Bundled Modules

The reference binary (`cmd/app/main.go`) wires the default stack:

```go
authModule := auth.New()
apiKeyModule := apikeys.New(
    apikeys.WithHTTPMiddleware(authModule.RequireRoles("admin")),
)

jobModule := jobs.New(
    jobs.WithWorkerCount(jobWorkers),
    jobs.WithPollInterval(jobPollSeconds),
    jobs.WithTenantResolver(jobs.AuthTenantResolver()),
    jobs.WithHTTPMiddleware(
        apiKeyModule.RequireAPIKey(),
        authModule.RequireRoles("user", "operator", "admin"),
    ),
)

app.Register(authModule)
app.Register(apiKeyModule) // remove if you do not need API keys
app.Register(jobModule)
app.Register(webui.New(jobWorkers, jobPollSeconds,
    webui.WithTenantResolver(jobs.AuthTenantResolver()),
))
```

To opt out of a capability, remove the corresponding `app.Register(...)` call and any middleware that depends on it. Conversely, drop in your own modules using the same pattern—it is safe to mix first-party and custom components.

## Testing

The repository ships with a focused configuration test. Run the full suite via:

```bash
go test ./...
```

Add module tests close to their packages to keep coverage meaningful.

### Quality & Verification

- `go test ./...` exercises configuration loading, authentication token helpers, transport utilities, and any additional unit tests you add. UI regressions are covered with renderer tests so the landing carousel, navigation collapse, and syntax-highlighted featurettes stay intact.
- `make test` is a convenience alias if you prefer the bundled Makefile.
- `make snapshots` refreshes renderer snapshots by running `UPDATE_SNAPSHOTS=1 go test ./web`.
- `make smoke` launches the Docker Compose stack, probes `/healthz`, `/readyz`, and `/metrics`, and tears the stack down afterwards.
- `make smoke` respects `SMOKE_APP_PORT` / `SMOKE_DB_PORT` so it will not clash with local Postgres.
- When touching templates, run `go test ./web`; the renderer compiles embedded HTML and regression tests confirm the layout renders home and dashboard views with the bundled brand assets.
- For integration flows (auth + jobs + web UI), start the stack locally (`go run ./cmd/app -config config.yaml`) and walk through: `curl http://localhost:8080/healthz`, create a demo user via `/auth/signup`, enqueue a job via `/jobs`, and confirm the console at `/console` reflects the change.
- Using Docker? `docker compose logs -f app` streams application logs (including Postgres connection retries) so you can observe startup sequencing.
## Auth Module

`pkg/modules/auth` exposes `/auth/signup`, `/auth/signin`, `/auth/refresh`, `/auth/signout`, and `/auth/profile` JSON endpoints. Passwords are hashed with `bcrypt`, users live in the `auth_users` table, and stateless JWT tokens secure profile access. Configure the secret via `auth.token_secret` or the `SHANRAQ_AUTH_TOKEN_SECRET` env var.

- Access tokens embed the user role; seeds create `admin@shanraq.org` / `operator@shanraq.org` demo accounts.
- RBAC lives in `auth_roles` and `auth_user_roles`, and helper middleware `auth.Module.RequireRoles(...)` makes it easy to guard routes.
- Every sign-in issues a refresh token persisted in Postgres; `/auth/refresh` rotates them, the module trims old tokens automatically, and `/auth/signout` revokes the provided token.
- Password reset pages live at `/auth/password/reset` and `/auth/password/confirm`, with tokens stored in `auth_password_resets`.
- Optional TOTP-based MFA can be enabled via `auth.mfa.totp`, prompting users to enrol an authenticator app before tokens are issued.
- Handlers can call `jobs.InfoFromContext(ctx)` to access worker metadata when enqueuing auth-related background jobs.
- Use the API key module (below) when you want to distribute scoped credentials to customers or service integrations.

## API Keys Module

`pkg/modules/apikeys` issues and validates tenant-scoped API keys. It exposes JSON endpoints under `/auth/apikeys` (list, create, revoke) and a middleware `apikeys.Module.RequireAPIKey()` that hydrates `auth.Claims` for downstream handlers. Keys are stored hashed in the `auth_api_keys` table, and every key is associated with an `auth_users` row.

- Developers manage their keys via the REST API (protected with `auth.RequireRoles`).
- Integrations authenticate with `X-API-Key` or `Authorization: ApiKey <token>` headers—no JWT required.
- Combine `RequireAPIKey()` with `auth.RequireRoles(...)` to accept either JWTs or API keys on the same route set.
- Demo data seeds an operator key (`sk_demo_operator_token`) so the console and API explorer work immediately—rotate it outside local development.

## Background Jobs

- `POST /jobs` enqueues work with arbitrary JSON payloads and optional `run_at` timestamps. When a tenant resolver is configured, the module automatically persists the `user_id` alongside each job.
- Workers (configured in `cmd/app/main.go`) poll Postgres, claim jobs with `FOR UPDATE SKIP LOCKED`, and retry automatically.
- Add business logic via `jobsModule.Handle("job-name", jobs.Handler)`; the reference app wires a `send_welcome_email` handler that decodes payloads, logs attempts, and uses context metadata.
- Explore and manage the queue from `/console` or via the JSON API: `GET /jobs?status=pending`, `POST /jobs/{id}/retry`, and `POST /jobs/{id}/cancel`. The HTTP handlers enforce tenant scoping when a resolver is provided.

The Web UI module queries the same queue to render status cards and recent jobs.

## Migrations

`pkg/modules/migrations` embeds Goose scripts under `pkg/modules/migrations/sql`. Every boot runs `goose.Up` inside the process, guaranteeing schema parity with binaries. Add migrations by dropping new `*.sql` files following Goose's timestamp naming and `-- +goose Up/Down` directives.

The seed migration provisions demo content (`framework_about` marketing copy, sample jobs, and admin/operator accounts). Seeded users are marked with `password_reset_required = true`; initiate `/auth/password/reset` to generate a reset link before signing in.

## Templates & Web UI

`web` hosts the renderer, landing carousel, and static bundle. `home.html` delivers a Bootstrap carousel whose copy pulls from the `framework_about` table, while `dashboard.html` powers the operator console. Shared partials such as `partials/queue.html` drive the queue explorer and modal form. Extend by adding new templates in `web/views` (with matching static assets under `web/static`).

The operator dashboard surfaces throughput for the last hour, live failure ratios, and the next scheduled job so runbooks stay actionable without leaving the browser. When a tenant resolver runs, the dashboard automatically filters metrics and recent jobs to the active account.

## Deployment Guide

A production-minded playbook lives in [`docs/DEPLOYMENT.md`](docs/DEPLOYMENT.md) and is also served in-app at [`/static/docs/deployment.html`](http://localhost:8080/static/docs/deployment.html).

## Web Documentation

The framework ships with an embedded handbook at [`/docs`](http://localhost:8080/docs) that summarises everyday workflows:

- **Quick start commands** to run the stack locally or through Docker.
- **Module catalogue** covering authentication, background jobs, telemetry, and their primary endpoints.
- **Operator console overview** highlighting landing, console, and Jobs API entry points.

Update `pkg/modules/webui/webui.go` if you need to surface additional modules or custom commands—rendered content lives in `web/views/docs.html`.

### Customising the landing page copy

Edit the `framework_about` table to tailor the carousel headline, subheadline, and feature slides. The most recent row is loaded during startup:

```sql
INSERT INTO framework_about (headline, subheadline, feature_one, feature_two, feature_three)
VALUES (
  'Shanraq Console',
  'A Go-first framework for resilient backends.',
  'PostgreSQL-native data layer with migrations built-in.',
  'Composable module system for HTTP, workers, and observability.',
  'Cloud-ready tooling: Docker, Prometheus telemetry, structured logging.'
);
```

Restart the app (or apply the insert before boot) and the home carousel will reflect the update automatically.
## Documentation

- [Configuration Guide](docs/CONFIGURATION.md): Every `config.yaml` field, environment override, and module-specific knob.
- [Deployment Guide](docs/DEPLOYMENT.md): Production Compose profile, Helm notes, and operational checklist.
- In-app documentation at [`/docs`](http://localhost:8080/docs) mirrors these resources with QR shortcuts directly in the UI.
