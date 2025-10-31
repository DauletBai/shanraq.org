# Shanraq Deployment Guide

This guide mirrors the in-app deployment page that ships at `/static/docs/deployment.html`, but it is versioned alongside the codebase for easy diffing.

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

## Kubernetes / Helm basics

1. Build the provided `Dockerfile` (CGO disabled, uses distroless base).
2. Package configuration as a `ConfigMap`; sensitive overrides go into a `Secret` (database DSN, token secret).
3. Deploy the application and Postgres as separate `Deployments`/`StatefulSets`; only expose the app publicly.
4. Configure readiness (`/readyz`) and liveness (`/healthz`) probes.
5. Optionally add a `HorizontalPodAutoscaler` that reacts to CPU or queue-length metrics.
6. Terminate TLS at an ingress controller or service mesh.

> Tip: if you already operate Helm, wrap the instructions above into a chart with templated environment variables.

## Operational checklist

- Rotate refresh tokens when resetting passwords. `/auth/signout` and `/auth/password/confirm` already do this automatically.
- Back up the database — both `job_queue` and `auth_*` tables hold critical state.
- Scrape `/metrics`; the dashboard displays queue throughput using Prometheus counters.
- Set `SHANRAQ_AUTH_TOKEN_SECRET` to a 32+ byte random string in all non-local environments.
- Keep an eye on the `job_queue` table size; archive or prune old jobs if necessary.

## Environment variables

| Variable | Purpose |
| -------- | ------- |
| `SHANRAQ_DATABASE_URL` | Postgres DSN (pgx format). |
| `SHANRAQ_AUTH_TOKEN_SECRET` | JWT signing key. |
| `SHANRAQ_SERVER_ADDRESS` | HTTP bind address (default `:8080`). |
| `SHANRAQ_TELEMETRY_ENABLE_METRICS` | Toggle Prometheus handler (true/false). |
| `SHANRAQ_LOGGING_MODE` | `development` or `production`. |
| `SHANRAQ_JOBS_WORKERS` | Overrides worker count (if you expose it via config). |

## Further reading

- [Source repository](https://github.com/DauletBai/shanraq.org)
- In-app documentation at [`/docs`](http://localhost:8080/docs) (includes quick commands, schema overview, and module guide).
