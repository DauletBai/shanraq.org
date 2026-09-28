# CloudLab

CloudLab is the running project for Shanraq's Cloud & DevOps course. It is a
small HTTP service with persistent notes, JSON access logs, health and readiness
endpoints, Prometheus-format metrics, graceful shutdown, container manifests,
CI, infrastructure examples, and backup/restore scripts.

Run locally:

```sh
go test ./...
go run ./cmd/cloudlab
```

Run the container stack:

```sh
docker compose up --build -d
./scripts/smoke.sh
docker compose logs -f app
```

Stop it without deleting the named data volume:

```sh
docker compose down
```

The course explains every command before it asks a learner to run it. Files in
this directory are reference checkpoints, not a replacement for the lessons.
