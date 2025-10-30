# Shanraq.org Framework

Shanraq is a batteries-included Go 1.25.x application framework focused on modern backend practices: typed configuration, PostgreSQL-first data access, structured logging, composable modules, and first-class observability.

## Highlights

- **Go 1.25.x ready**: Module targets Go `1.25` with Go toolchain `1.22.3`, so you can build today and upgrade seamlessly.
- **PostgreSQL via pgxpool**: Sane defaults, lifecycle management, and trace hooks wired into Zap.
- **Declarative configuration**: Environment-aware configuration loader with `.env` support, file overrides, and typed structs.
- **Composable modules**: `shanraq.Module` contracts let teams build reusable features that can hook into init/start/routes independently.
- **Operational insights**: Built-in health, readiness, and Prometheus telemetry modules provide production-grade endpoints out of the box.
- **Domain modules included**: Authentication endpoints, background job queue/workers, and a Bootstrap 5.3 dashboard ship as ready-to-use modules.
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
pkg/modules/jobs      # Postgres-backed job queue & workers
pkg/modules/migrations  # Embedded Goose migrations
pkg/modules/webui     # Bootstrap 5.3 dashboard module
pkg/modules/...       # Additional modules (health, telemetry, etc.)
web/                  # html/template renderer, landing carousel, static assets
```

## Getting Started

1. Install Go >= 1.22.3 (toolchain) and PostgreSQL.
2. Copy and tweak the sample config (the binary refuses to start without it):
   ```bash
   cp configs/config.example.yaml config.yaml
   ```
3. (Optional) export env overrides. Every key supports the `SHANRAQ_` prefix (e.g. `SHANRAQ_SERVER_ADDRESS=:9090`).
4. Run the reference app:
   ```bash
   go run ./cmd/app -config config.yaml
   ```
5. Visit `http://localhost:8080/` for the landing carousel and queue explorer, `http://localhost:8080/console` for the operator dashboard, plus `http://localhost:8080/healthz`, `http://localhost:8080/readyz`, and `http://localhost:8080/metrics`.
6. Need to override secrets or DSNs? Use env vars such as `export SHANRAQ_AUTH_TOKEN_SECRET=$(openssl rand -hex 32)` before launching the binary.

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
| `logging.level` | Zap level (`debug/info/warn/error`) | `info` |
| `logging.mode` | `development` or `production` encoder | `production` |
| `auth.token_secret` | HMAC secret for JWT tokens | `replace-me-now` |
| `auth.token_ttl` | Token lifetime | `15m` |

Durations accept Go duration syntax (e.g. `30s`, `5m`).

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

## Testing

The repository ships with a focused configuration test. Run the full suite via:

```bash
go test ./...
```

Add module tests close to their packages to keep coverage meaningful.

### Quality & Verification

- `go test ./...` exercises configuration loading, authentication token helpers, transport utilities, and any additional unit tests you add. UI regressions are covered with renderer tests so the landing carousel, navigation collapse, and syntax-highlighted featurettes stay intact.
- `make test` is a convenience alias if you prefer the bundled Makefile.
- When touching templates, run `go test ./web`; the renderer compiles embedded HTML and regression tests confirm the layout renders home and dashboard views with the bundled brand assets.
- For integration flows (auth + jobs + web UI), start the stack locally (`go run ./cmd/app -config config.yaml`) and walk through: `curl http://localhost:8080/healthz`, create a demo user via `/auth/signup`, enqueue a job via `/jobs`, and confirm the console at `/console` reflects the change.
- Using Docker? `docker compose logs -f app` streams application logs (including Postgres connection retries) so you can observe startup sequencing.
## Auth Module

`pkg/modules/auth` exposes `/auth/signup`, `/auth/signin`, and `/auth/profile` JSON endpoints. Passwords are hashed with `bcrypt`, users live in the `auth_users` table, and stateless JWT tokens secure profile access. Configure the secret via `auth.token_secret` or the `SHANRAQ_AUTH_TOKEN_SECRET` env var.

## Background Jobs

- `POST /jobs` enqueues work with arbitrary JSON payloads and optional `run_at` timestamps.
- Workers (configured in `cmd/app/main.go`) poll Postgres, claim jobs with `FOR UPDATE SKIP LOCKED`, and retry automatically.
- Add business logic via `jobsModule.Handle("job-name", jobs.Handler)`; the example registers `send_welcome_email` using `jobs.LogHandler`.
- Explore and manage the queue from `/console` or via the JSON API: `GET /jobs?status=pending`, `POST /jobs/{id}/retry`, and `POST /jobs/{id}/cancel`.

The Web UI module queries the same queue to render status cards and recent jobs.

## Migrations

`pkg/modules/migrations` embeds Goose scripts under `pkg/modules/migrations/sql`. Every boot runs `goose.Up` inside the process, guaranteeing schema parity with binaries. Add migrations by dropping new `*.sql` files following Goose's timestamp naming and `-- +goose Up/Down` directives.

## Templates & Web UI

`web` hosts the renderer, landing carousel, and static bundle. `home.html` delivers a Bootstrap carousel whose copy pulls from the `framework_about` table, while `dashboard.html` powers the operator console. Shared partials such as `partials/queue.html` drive the queue explorer and modal form. Extend by adding new templates in `web/views` (with matching static assets under `web/static`).

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
