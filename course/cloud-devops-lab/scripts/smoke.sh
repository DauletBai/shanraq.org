#!/bin/sh
set -eu

base_url=${1:-http://localhost:8080}
curl --fail --silent --show-error "$base_url/healthz" >/dev/null
curl --fail --silent --show-error "$base_url/readyz" >/dev/null
curl --fail --silent --show-error "$base_url/metrics" | grep -q cloudlab_http_requests_total
printf 'CloudLab smoke check passed: %s\n' "$base_url"
