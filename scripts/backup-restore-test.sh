#!/usr/bin/env bash
#
# backup-restore-test.sh — restore a Shanraq backup into a scratch database and
# check that what came back is a site, not a file that merely unpacked.
#
# A backup nobody has restored is a hope. The first one was verified by hand on
# 14 August 2026; this is that procedure written down so it can be repeated
# after every schema change, and so "we have backups" stops being a belief.
#
# It runs where the age private key lives — the owner's machine — because the
# key is deliberately not on the server it protects. Nothing here writes to
# production: it creates a scratch database, reads it, and drops it.
#
# Usage:
#   scripts/backup-restore-test.sh <archive.tar.gz.age | archive.tar.gz> [--keep]
#
# Environment:
#   BACKUP_AGE_IDENTITY   path to the age private key (required for .age archives)
#   RESTORE_TEST_DB       scratch database name        [shanraq_restore_test]
#   RESTORE_TEST_HOST     Postgres host                [localhost]
#   RESTORE_TEST_PORT     Postgres port                [5432]
#   RESTORE_TEST_USER     Postgres user                [$(whoami)]
#   PG_BIN                dir holding pg_restore/psql matching the server (16)
#
# Exits non-zero on the first thing that is not right, so a scheduler can treat
# a silent success as the only good news.
set -euo pipefail

log()  { printf '%s %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$*"; }
fail() { printf '%s ERROR: %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$*" >&2; exit 1; }

ARCHIVE="${1:-}"
KEEP="${2:-}"
[[ -n "$ARCHIVE" ]] || fail "usage: $0 <archive.tar.gz.age|archive.tar.gz> [--keep]"
[[ -f "$ARCHIVE" ]] || fail "no such archive: $ARCHIVE"

DB="${RESTORE_TEST_DB:-shanraq_restore_test}"
HOST="${RESTORE_TEST_HOST:-localhost}"
PORT="${RESTORE_TEST_PORT:-5432}"
USER="${RESTORE_TEST_USER:-$(whoami)}"
[[ "$DB" == *test* ]] || fail "RESTORE_TEST_DB must name a test database; refusing '$DB'"

# The tools must match the server that wrote the dump, not whichever Postgres
# happens to be first in PATH. A machine with both 14 and 16 installed will
# hand you the older pg_restore and a file-header error that reads like a
# corrupt archive — the same fallback backup.sh uses, for the same reason.
if [[ -z "${PG_BIN:-}" ]]; then
  for candidate in /opt/homebrew/opt/postgresql@16/bin /usr/local/opt/postgresql@16/bin /usr/lib/postgresql/16/bin; do
    [[ -x "$candidate/pg_restore" ]] && { PG_BIN="$candidate"; break; }
  done
fi
PSQL="psql"; PGRESTORE="pg_restore"; CREATEDB="createdb"; DROPDB="dropdb"
if [[ -n "${PG_BIN:-}" ]]; then
  PSQL="$PG_BIN/psql"; PGRESTORE="$PG_BIN/pg_restore"; CREATEDB="$PG_BIN/createdb"; DROPDB="$PG_BIN/dropdb"
  log "using Postgres tools from $PG_BIN"
fi
conn=(-h "$HOST" -p "$PORT" -U "$USER")

# The decrypted copy holds e-mail addresses, telephone numbers and avatars. It
# lives in a temporary directory that goes away however this script ends —
# including the failure paths, which is when it is easiest to forget.
WORK="$(mktemp -d)"
cleanup() {
  rm -rf "$WORK"
  if [[ "$KEEP" != "--keep" ]]; then
    "$DROPDB" "${conn[@]}" --if-exists "$DB" 2>/dev/null || true
  else
    log "kept scratch database '$DB' — drop it yourself when done"
  fi
}
trap cleanup EXIT

# ---- decrypt ----
TARBALL="$WORK/backup.tar.gz"
if [[ "$ARCHIVE" == *.age ]]; then
  command -v age >/dev/null || fail "'age' is not installed"
  [[ -n "${BACKUP_AGE_IDENTITY:-}" ]] || fail "BACKUP_AGE_IDENTITY must point at the age private key"
  [[ -f "$BACKUP_AGE_IDENTITY" ]] || fail "no such identity file: $BACKUP_AGE_IDENTITY"
  log "decrypting…"
  age -d -i "$BACKUP_AGE_IDENTITY" -o "$TARBALL" "$ARCHIVE" || fail "decryption failed — wrong key?"
else
  log "archive is not encrypted; reading as-is"
  cp "$ARCHIVE" "$TARBALL"
fi

# ---- unpack and verify checksums ----
tar xzf "$TARBALL" -C "$WORK" || fail "archive did not unpack"
for member in db.dump media.tar.gz SHA256SUMS; do
  [[ -f "$WORK/$member" ]] || fail "archive is missing $member"
done
log "verifying checksums…"
if command -v sha256sum >/dev/null; then
  ( cd "$WORK" && sha256sum -c SHA256SUMS >/dev/null ) || fail "checksum mismatch — the archive is damaged"
else
  ( cd "$WORK" && shasum -a 256 -c SHA256SUMS >/dev/null ) || fail "checksum mismatch — the archive is damaged"
fi

# ---- media ----
MEDIA_FILES="$(tar tzf "$WORK/media.tar.gz" | grep -cv '/$' || true)"
log "media archive holds $MEDIA_FILES files"

# ---- restore ----
log "restoring into scratch database '$DB'…"
"$DROPDB" "${conn[@]}" --if-exists "$DB" 2>/dev/null || true
"$CREATEDB" "${conn[@]}" -T template0 -E UTF8 "$DB" || fail "could not create the scratch database"
# --no-owner because the roles on this machine are not the server's. Errors are
# not tolerated: a restore that "mostly worked" is the kind of backup that is
# discovered to be useless on the day it is needed.
"$PGRESTORE" "${conn[@]}" -d "$DB" --no-owner --exit-on-error "$WORK/db.dump" \
  || fail "pg_restore failed"

q() { "$PSQL" "${conn[@]}" -d "$DB" -At -c "$1"; }

TABLES="$(q "SELECT count(*) FROM information_schema.tables WHERE table_schema='public'")"
ARTICLES="$(q "SELECT count(*) FROM articles WHERE status='published'")"
TRANSLATIONS="$(q "SELECT count(*) FROM article_translations")"
USERS="$(q "SELECT count(*) FROM auth_users")"
ADMINS="$(q "SELECT count(*) FROM auth_users WHERE role IN ('admin','director')")"
LISTINGS="$(q "SELECT count(*) FROM listings")"
COMMENTS="$(q "SELECT count(*) FROM comments")"
GOOSE="$(q "SELECT max(version_id) FROM goose_db_version")"

cat <<REPORT

  restored from : $(basename "$ARCHIVE")
  tables        : $TABLES
  published     : $ARTICLES articles, $TRANSLATIONS translations
  accounts      : $USERS ($ADMINS administrator(s))
  listings      : $LISTINGS
  comments      : $COMMENTS
  media files   : $MEDIA_FILES
  schema version: $GOOSE

REPORT

# ---- what must be true for this to be a backup rather than a file ----
[[ "$TABLES"   -gt 20 ]] || fail "only $TABLES tables came back; the dump is not a whole schema"
[[ "$USERS"    -gt 0  ]] || fail "no accounts came back"
[[ "$ARTICLES" -gt 0  ]] || fail "no published articles came back"
[[ -n "$GOOSE"        ]] || fail "no migration history came back; the schema version is unknown"
# Restoring into a site nobody can administer is a restore that has to be
# repaired by hand before it is usable.
[[ "$ADMINS"   -gt 0  ]] || fail "the restored data holds no administrator"

log "restore test passed."
