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

METRICS_FILE=""
cleanup() {
  # An `if`, not `[ … ] && rm`: with no metrics file the && list fails, and
  # under `set -e` a failing trap body stops before the stack comes down.
  if [ -n "${METRICS_FILE}" ]; then rm -f "${METRICS_FILE}"; fi
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

# Read the body to a file first, then show the top of it.
#
# `curl … | head -n 5` looks equivalent and is not: head exits after five lines
# and closes the pipe, curl takes a SIGPIPE writing the rest, and under
# `set -o pipefail` that fails the whole run. /metrics is about 90 KB, so
# whether curl finishes before head leaves is a race — one this passed locally
# and lost on every CI runner.
METRICS_FILE="$(mktemp)"
if ! curl -fsS "http://localhost:${APP_PORT}/metrics" -o "${METRICS_FILE}"; then
  echo "⚠️ Failed to read metrics endpoint" >&2
  ${COMPOSE} -p "${STACK_NAME}" logs app
  exit 1
fi
# A 200 is not the same as metrics: assert it is Prometheus exposition format
# rather than an error page served with the wrong status.
if ! grep -q '^# HELP ' "${METRICS_FILE}"; then
  echo "⚠️ /metrics answered but the body is not Prometheus exposition format" >&2
  head -n 5 "${METRICS_FILE}" >&2
  ${COMPOSE} -p "${STACK_NAME}" logs app
  exit 1
fi
head -n 5 "${METRICS_FILE}"
echo "✅ /metrics passed"

echo "🎉 Docker smoke check succeeded."
