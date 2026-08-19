#!/usr/bin/env bash
#
# backup.sh — scheduled, encrypted, retained backup of Shanraq
# (the Postgres database AND the uploaded media: avatars, images).
#
# Zero-config: with just a reachable database it writes plaintext archives to
# ./backups. Add BACKUP_AGE_RECIPIENT to encrypt, and BACKUP_UPLOAD_CMD to ship
# off-site. Exits non-zero on ANY failure so cron/systemd can alert.
#
# Environment:
#   SHANRAQ_DATABASE_URL   Postgres DSN (direct mode)          [dev default]
#   SHANRAQ_MEDIA_DIR      media dir to archive (direct mode)  [./data/media]
#   BACKUP_COMPOSE_FILE    if set, dump via `docker compose -f <f> exec db`
#   BACKUP_MEDIA_VOLUME    if set, tar media from this docker volume (compose mode)
#   BACKUP_DIR             where archives are written          [./backups]
#   BACKUP_RETENTION       how many archives to keep locally   [7]
#   BACKUP_AGE_RECIPIENT   age public key → encrypt the archive  [none = plaintext]
#   BACKUP_UPLOAD_CMD      off-site command; {file} is replaced with the archive path
#   BACKUP_REQUIRE_OFFSITE set to 1 once a destination exists: a missing
#                          BACKUP_UPLOAD_CMD then fails the run instead of warning
#   PG_BIN                 dir holding a pg_dump matching the server (Postgres 16)
#
set -euo pipefail

log()  { printf '%s %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$*"; }
fail() { log "ERROR: $*"; exit 1; }

DB_URL="${SHANRAQ_DATABASE_URL:-postgres://postgres:postgres@127.0.0.1:5432/shanraq?sslmode=disable}"
MEDIA_DIR="${SHANRAQ_MEDIA_DIR:-./data/media}"
BACKUP_DIR="${BACKUP_DIR:-./backups}"
RETENTION="${BACKUP_RETENTION:-7}"
TS="$(date -u +%Y%m%d-%H%M%SZ)-$$"   # PID suffix avoids same-second overwrite
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT
mkdir -p "$BACKUP_DIR"

sha256() { command -v sha256sum >/dev/null && sha256sum "$@" || shasum -a 256 "$@"; }

# ---- database ----
log "dumping database…"
if [[ -n "${BACKUP_COMPOSE_FILE:-}" ]]; then
  docker compose -f "$BACKUP_COMPOSE_FILE" exec -T db pg_dump -U shanraq -Fc shanraq > "$WORK/db.dump" \
    || fail "pg_dump (compose) failed"
else
  PGDUMP="pg_dump"
  [[ -n "${PG_BIN:-}" ]] && PGDUMP="$PG_BIN/pg_dump"
  [[ -z "${PG_BIN:-}" && -x /opt/homebrew/opt/postgresql@16/bin/pg_dump ]] && PGDUMP=/opt/homebrew/opt/postgresql@16/bin/pg_dump
  [[ -z "${PG_BIN:-}" && -x /usr/local/opt/postgresql@16/bin/pg_dump ]] && PGDUMP=/usr/local/opt/postgresql@16/bin/pg_dump
  "$PGDUMP" -Fc -d "$DB_URL" -f "$WORK/db.dump" || fail "pg_dump failed"
fi

# ---- media ----
log "archiving media…"
if [[ -n "${BACKUP_MEDIA_VOLUME:-}" ]]; then
  docker run --rm -v "${BACKUP_MEDIA_VOLUME}:/data:ro" -v "$WORK:/out" alpine \
    tar czf /out/media.tar.gz -C /data . || fail "media archive (volume) failed"
elif [[ -d "$MEDIA_DIR" ]]; then
  tar czf "$WORK/media.tar.gz" -C "$(dirname "$MEDIA_DIR")" "$(basename "$MEDIA_DIR")" || fail "media archive failed"
else
  log "note: media dir '$MEDIA_DIR' not found — database only"
  tar czf "$WORK/media.tar.gz" -T /dev/null   # valid empty archive
fi

# ---- checksums + single archive ----
( cd "$WORK" && sha256 db.dump media.tar.gz > SHA256SUMS )
ARCHIVE="$BACKUP_DIR/shanraq-backup-$TS.tar.gz"
tar czf "$ARCHIVE" -C "$WORK" db.dump media.tar.gz SHA256SUMS

# ---- optional encryption (age) ----
if [[ -n "${BACKUP_AGE_RECIPIENT:-}" ]]; then
  command -v age >/dev/null || fail "BACKUP_AGE_RECIPIENT is set but 'age' is not installed"
  age -r "$BACKUP_AGE_RECIPIENT" -o "$ARCHIVE.age" "$ARCHIVE" || fail "age encryption failed"
  rm -f "$ARCHIVE"; ARCHIVE="$ARCHIVE.age"
  log "encrypted → $(basename "$ARCHIVE")"
else
  log "WARNING: BACKUP_AGE_RECIPIENT not set — the archive is NOT encrypted (do not ship it off-site as-is)"
fi
log "backup ready: $ARCHIVE ($(du -h "$ARCHIVE" | cut -f1))"

# ---- optional off-site upload ----
if [[ -n "${BACKUP_UPLOAD_CMD:-}" ]]; then
  cmd="${BACKUP_UPLOAD_CMD//\{file\}/$ARCHIVE}"
  log "off-site upload: $cmd"
  eval "$cmd" || fail "off-site upload failed"
elif [[ "${BACKUP_REQUIRE_OFFSITE:-0}" == "1" ]]; then
  # Once a destination exists, its disappearance must stop the run. A .env
  # rewritten by hand is exactly how a working off-site copy quietly becomes a
  # backup sitting on the machine it protects.
  fail "BACKUP_REQUIRE_OFFSITE=1 but BACKUP_UPLOAD_CMD is not set"
else
  log "WARNING: no BACKUP_UPLOAD_CMD — this archive stays on the machine it is"
  log "         protecting, so it survives a bad migration but not a lost host"
fi

# ---- retention: keep the newest N, prune older ----
ls -1t "$BACKUP_DIR"/shanraq-backup-*.tar.gz* 2>/dev/null | tail -n +"$((RETENTION + 1))" | while read -r old; do
  rm -f "$old" && log "pruned $(basename "$old")"
done

log "done."
