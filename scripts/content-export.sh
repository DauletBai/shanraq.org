#!/usr/bin/env bash
# Content-only export: the writing, and nothing that belongs to a person.
#
# The full backup cannot leave the country — it holds accounts, e-mail
# addresses, sellers' telephone numbers and avatars, and Article 12(2) of the
# Law on Personal Data keeps those in a database inside Kazakhstan. This archive
# holds published articles with their translations, the prediction ledger and
# the editable pages. It may live anywhere, which is the whole point: the one
# copy that survives losing the machine, the provider, or the country.
#
# Cover images: the uploaded ones are in here, the repository ones are not.
# Ninety of the hundred and eight articles point at /static/covers, which ships
# in the source tree and is therefore already off-site on GitHub. The other
# eighteen were uploaded and exist on this disk and nowhere else — those are
# copied. Only files a published article points at: the same volume holds
# avatars, and those belong to people.
#
# Environment (from /opt/shanraq/.env):
#   CONTENT_DIR          where archives are written            [/opt/shanraq/exports]
#   CONTENT_RETENTION    how many to keep locally              [7]
#   CONTENT_UPLOAD_CMD   off-site command, {file} is the path  [none = local only]
set -euo pipefail

log()  { printf '%s %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$*"; }
fail() { log "ERROR: $*"; exit 1; }

APP_DIR=/opt/shanraq
COMPOSE_FILE=docker-compose.prod.yml
cd "$APP_DIR"

# Read the file, never execute it. .env is written for docker compose and for
# systemd's EnvironmentFile, both of which parse it; a shell sources it, and one
# unquoted value with spaces in it — SHANRAQ_OPERATOR_LEGAL_NAME=«Казна
# Технолоджис» — is then a command that does not exist.
envval() {
  sed -n "s/^$1=//p" .env | head -1 | sed -e 's/^"\(.*\)"$/\1/' -e "s/^'\(.*\)'\$/\1/"
}
POSTGRES_PASSWORD="$(envval POSTGRES_PASSWORD)"
OUT_DIR="$(envval CONTENT_DIR)"; OUT_DIR="${OUT_DIR:-$APP_DIR/exports}"
RETENTION="$(envval CONTENT_RETENTION)"; RETENTION="${RETENTION:-7}"
CONTENT_UPLOAD_CMD="$(envval CONTENT_UPLOAD_CMD)"
STAMP="$(date -u +%Y%m%d-%H%M%S)"
WORK="$OUT_DIR/.work-$STAMP"

mkdir -p "$WORK"
trap 'rm -rf "$WORK"' EXIT

[[ -n "${POSTGRES_PASSWORD:-}" ]] || fail "POSTGRES_PASSWORD is not set in .env"
DSN="postgres://shanraq:${POSTGRES_PASSWORD}@db:5432/shanraq?sslmode=disable"

# The exporter ships inside the app image, so no toolchain is needed here. It
# joins the compose network, which is how db:5432 resolves.
log "exporting content"
docker compose -f "$COMPOSE_FILE" run --rm --no-deps \
  -v "$WORK:/out" \
  -v shanraq_media-data:/media-src:ro \
  -e DATABASE_URL="$DSN" -e EXPORT_DIR=/out -e MEDIA_ROOT=/media-src \
  --entrypoint /usr/local/bin/export app || fail "export failed"

[[ -s "$WORK/articles.json" ]] || fail "export produced no articles.json"

ARCHIVE="$OUT_DIR/shanraq-content-$STAMP.tar.gz"
tar -czf "$ARCHIVE" -C "$WORK" .
log "content ready: $ARCHIVE ($(du -h "$ARCHIVE" | cut -f1))"

# Not encrypted, deliberately: there is nothing in here that is not already
# published at shanraq.org, and an unencrypted copy is one a mirror can serve
# and a stranger can restore without holding our key.
if [[ -n "${CONTENT_UPLOAD_CMD:-}" ]]; then
  cmd="${CONTENT_UPLOAD_CMD//\{file\}/$ARCHIVE}"
  log "off-site upload: $cmd"
  eval "$cmd" || fail "off-site upload failed"
else
  log "CONTENT_UPLOAD_CMD is not set — the copy stays on this machine, which is the problem it exists to solve"
fi

ls -1t "$OUT_DIR"/shanraq-content-*.tar.gz 2>/dev/null | tail -n +"$((RETENTION + 1))" | while read -r old; do
  rm -f "$old" && log "pruned $(basename "$old")"
done
log "done."
