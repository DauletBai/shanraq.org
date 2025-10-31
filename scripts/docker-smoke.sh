#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT_DIR}"

if ! command -v docker >/dev/null 2>&1; then
  echo "docker command not found. Install Docker Desktop or CLI before running the smoke check." >&2
  exit 1
fi

if ! docker compose version >/dev/null 2>&1; then
  echo "docker compose plugin is required (Docker CLI v20.10+)." >&2
  exit 1
fi

STACK_NAME="shanraq-smoke"
COMPOSE="docker compose"
APP_PORT="${SMOKE_APP_PORT:-18080}"
DB_PORT="${SMOKE_DB_PORT:-15432}"

echo "Using app port ${APP_PORT} and database port ${DB_PORT} for smoke test."

cleanup() {
  ${COMPOSE} -p "${STACK_NAME}" down --remove-orphans >/dev/null 2>&1 || true
}

trap cleanup EXIT

echo "🧪 Starting Shanraq Docker smoke test..."
APP_PORT="${APP_PORT}" DB_PORT="${DB_PORT}" ${COMPOSE} -p "${STACK_NAME}" up --build -d

deadline=$((SECONDS + 120))
until curl -fsS "http://localhost:${APP_PORT}/healthz" >/dev/null 2>&1; do
  if (( SECONDS >= deadline )); then
    echo "❌ Application failed to become healthy within timeout." >&2
    ${COMPOSE} -p "${STACK_NAME}" logs app
    exit 1
  fi
  sleep 2
done

echo "✅ /healthz passed"

if ! curl -fsS "http://localhost:${APP_PORT}/readyz" >/dev/null 2>&1; then
  echo "⚠️ /readyz returned non-success status" >&2
  ${COMPOSE} -p "${STACK_NAME}" logs app
  exit 1
fi
echo "✅ /readyz passed"

if ! curl -fsS "http://localhost:${APP_PORT}/metrics" | head -n 5; then
  echo "⚠️ Failed to read metrics endpoint" >&2
  ${COMPOSE} -p "${STACK_NAME}" logs app
  exit 1
fi

echo "🎉 Docker smoke check succeeded."
