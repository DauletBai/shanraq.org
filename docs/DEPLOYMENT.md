# Shanraq Deployment Guide

This guide mirrors the in-app deployment page that ships at `/static/docs/deployment.html`, but it is versioned alongside the codebase for easy diffing. Read it together with the [Configuration Guide](CONFIGURATION.md) when tailoring `config.yaml` and environment variables.

## Docker Compose (production profile)

Create a dedicated compose file (for example `docker-compose.prod.yml`) so you can tune restart policies, secrets, and volumes separately from developer defaults:

```yaml
services:
  db:
    image: postgres:16
    restart: unless-stopped
    environment:
      POSTGRES_DB: shanraq
      POSTGRES_USER: shanraq
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - db-data:/var/lib/postgresql/data
  app:
    build: .
    restart: unless-stopped
    environment:
      SHANRAQ_DATABASE_URL: ${DATABASE_URL}
      SHANRAQ_AUTH_TOKEN_SECRET: ${AUTH_SECRET}
      SHANRAQ_AUTH_TOKEN_TTL: 45m
      SHANRAQ_SERVER_ADDRESS: 0.0.0.0:8080
    depends_on:
      - db
    command: ["/usr/local/bin/shanraq", "-config", "/app/config.yaml"]
volumes:
  db-data:
```

Run with:

```bash
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d
```

Recommendations:

- store credentials in your orchestrator or `.env.prod`
- point `SHANRAQ_DATABASE_URL` at a managed Postgres instance when leaving local development
- add `logging.mode: production` in `config.yaml`
- seed or migrate the database before scaling replicas: `docker compose run --rm app /usr/local/bin/shanraq -config /app/config.yaml` (the process runs migrations on startup and then exits when interrupted)
- rotate `SHANRAQ_AUTH_TOKEN_SECRET` when moving between environments; the runtime refuses to start with default secrets while `environment=production`.

## Role & Secret Management

The authentication module now persists roles in `auth_roles` and `auth_user_roles`. To add more roles in production:

```sql
INSERT INTO auth_roles (name, description) VALUES ('auditor', 'Read-only access to reports')
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description;

INSERT INTO auth_user_roles (user_id, role_id)
SELECT u.id, r.id FROM auth_users u, auth_roles r
WHERE u.email = 'auditor@example.com' AND r.name = 'auditor';
```

Refresh tokens are trimmed to the most recent five entries per user. When rotating secrets:

1. Deploy the new binary with the updated `SHANRAQ_AUTH_TOKEN_SECRET`.
2. Invalidate outstanding tokens via `UPDATE auth_refresh_tokens SET revoked_at = NOW();`.
3. Prompt operators to sign in again.

## Migrating existing data (database + media)

Two things hold state and BOTH must move together — the database stores an
avatar's URL, but the actual photo lives on disk in the media directory:

- **Postgres database** — accounts, articles, listings, author bios, the
  `avatar_url` column, etc.
- **Media directory** (`media.dir`, default `./data/media`) — uploaded files:
  **avatars** and article/listing images.

Use the helper scripts. `pg_dump`/`pg_restore` must match the server major
version (Postgres 16); the scripts auto-detect a Homebrew `postgresql@16`, or
set `PG_BIN=/path/to/pg16/bin`.

**On the source machine** — create a snapshot (two files: `db.dump` +
`media.tar.gz`):

```bash
SHANRAQ_DATABASE_URL=postgres://postgres:pass@127.0.0.1:5432/shanraq \
SHANRAQ_MEDIA_DIR=./data/media \
scripts/data-snapshot.sh ./snapshot
```

**Copy both files** to the target server (they contain personal data — never
commit them, transfer over SSH):

```bash
scp -r ./snapshot user@your-server:/opt/shanraq/snapshot
```

**On the target server** — create an EMPTY database (do NOT run
`cmd/migrate` first; the dump already carries the schema at HEAD), then restore:

```bash
createdb -U postgres shanraq   # empty target DB
SHANRAQ_DATABASE_URL=postgres://postgres:pass@127.0.0.1:5432/shanraq \
SHANRAQ_MEDIA_PARENT=./data \
scripts/data-restore.sh ./snapshot
```

`SHANRAQ_MEDIA_PARENT` is the parent of the media dir (the archive holds a
top-level `media/`), so it must match the target's `media.dir` (default
`./data/media` → parent `./data`). Start the app; goose reports no pending
migrations.

> **After a restore onto a public server:** reset the passwords of the seeded
> staff accounts (`admin@`, `operator@`, …) — the dump carries their existing
> hashes. Use `adminctl` or a password-reset. Your own account, its trilingual
> bio, and your avatar come across intact.

## Kubernetes / Helm basics

1. Build the provided `Dockerfile` (CGO disabled, uses distroless base).
2. Package configuration as a `ConfigMap`; sensitive overrides go into a `Secret` (database DSN, token secret).
3. Deploy the application and Postgres as separate `Deployments`/`StatefulSets`; only expose the app publicly.
4. Configure readiness (`/readyz`) and liveness (`/healthz`) probes.
5. Optionally add a `HorizontalPodAutoscaler` that reacts to CPU or queue-length metrics.
6. Terminate TLS at an ingress controller or service mesh.

> Tip: if you already operate Helm, wrap the instructions above into a chart with templated environment variables.

### Helm Chart Outline

- `values.yaml`: expose `image`, `pullPolicy`, `service`, `resources`, `env`, and `config` sections.
- `templates/deployment.yaml`: mount the rendered `config.yaml`, inject env vars, and wire readiness/liveness probes.
- `templates/secret.yaml`: populate `SHANRAQ_AUTH_TOKEN_SECRET` using `lookup` for existing secrets to avoid regeneration.
- `templates/hpa.yaml`: optional HorizontalPodAutoscaler referencing CPU or custom Prometheus metrics (queue depth).
- `templates/service.yaml` and `templates/ingress.yaml`: expose HTTP under your chosen domain with TLS termination.

## Operational checklist

- Rotate refresh tokens when resetting passwords. `/auth/signout` and `/auth/password/confirm` already do this automatically.
- Keep `auth_roles` in sync with your business requirements; migrations can seed defaults, but production seeds belong in separate change scripts.
- Audit `auth_api_keys` when rotating credentials; revoke stale keys so per-tenant scopes remain enforced.
- Remove or rotate the demo operator key (`sk_demo_operator_token`) before exposing the API publicly.
- Configure SMTP credentials via `notifications.smtp.*` (or `SHANRAQ_NOTIFICATIONS_SMTP_*` env vars) so password resets reach end users.
- Back up the database — both `job_queue` and `auth_*` tables hold critical state.
- Scrape `/metrics`; the dashboard displays queue throughput using Prometheus counters.
- Set `SHANRAQ_AUTH_TOKEN_SECRET` to a 32+ byte random string in all non-local environments.
- Keep an eye on the `job_queue` table size; archive or prune old jobs if necessary.
- Watch for `"refresh token reuse attempt"` warnings in the logs — they indicate clients presenting revoked tokens.

## Environment variables

| Variable | Purpose |
| -------- | ------- |
| `SHANRAQ_DATABASE_URL` | Postgres DSN (pgx format). |
| `SHANRAQ_AUTH_TOKEN_SECRET` | JWT signing key. |
| `SHANRAQ_AUTH_TOKEN_TTL` | Overrides access token lifetime. |
| `SHANRAQ_SERVER_ADDRESS` | HTTP bind address (default `:8080`). |
| `SHANRAQ_TELEMETRY_ENABLE_METRICS` | Toggle Prometheus handler (true/false). |
| `SHANRAQ_LOGGING_MODE` | `development` or `production`. |
| `SHANRAQ_JOBS_WORKERS` | Overrides worker count (if you expose it via config). |

## Further reading

- [Source repository](https://github.com/DauletBai/shanraq.org)
- In-app documentation at [`/docs`](http://localhost:8080/docs) (includes quick commands, schema overview, and module guide).
