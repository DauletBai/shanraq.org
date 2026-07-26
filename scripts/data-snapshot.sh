#!/usr/bin/env bash
#
# data-snapshot.sh — snapshot Shanraq data for moving between servers.
#
# Produces two files in <out-dir>:
#   db.dump        — the whole Postgres database (schema + rows + goose version)
#   media.tar.gz   — every uploaded file (AVATARS and images) from the media dir
#
# Both are needed: the database only stores the avatar's URL, the actual photo
# lives on disk in the media directory. Both files contain personal data
# (emails, password hashes, photos) — NEVER commit them to git.
#
# Usage:
#   SHANRAQ_DATABASE_URL=postgres://user:pass@host:5432/shanraq \
#   SHANRAQ_MEDIA_DIR=./data/media \
#   scripts/data-snapshot.sh [out-dir]
#
set -euo pipefail

DB_URL="${SHANRAQ_DATABASE_URL:-postgres://postgres:postgres@127.0.0.1:5432/shanraq?sslmode=disable}"
MEDIA_DIR="${SHANRAQ_MEDIA_DIR:-./data/media}"
OUT="${1:-shanraq-snapshot-$(date +%Y%m%d-%H%M%S)}"

# pg_dump MUST be at least the server's major version (the server is Postgres 16).
# Prefer an explicit PG_BIN, then a Homebrew postgresql@16, then PATH.
PGDUMP="pg_dump"
if [[ -n "${PG_BIN:-}" ]]; then PGDUMP="$PG_BIN/pg_dump"
elif [[ -x /opt/homebrew/opt/postgresql@16/bin/pg_dump ]]; then PGDUMP=/opt/homebrew/opt/postgresql@16/bin/pg_dump
elif [[ -x /usr/local/opt/postgresql@16/bin/pg_dump ]]; then PGDUMP=/usr/local/opt/postgresql@16/bin/pg_dump
fi

mkdir -p "$OUT"
echo "→ database → $OUT/db.dump   ($("$PGDUMP" --version))"
"$PGDUMP" -Fc -d "$DB_URL" -f "$OUT/db.dump"

if [[ -d "$MEDIA_DIR" ]]; then
  echo "→ media    → $OUT/media.tar.gz   ($MEDIA_DIR)"
  tar -czf "$OUT/media.tar.gz" -C "$(dirname "$MEDIA_DIR")" "$(basename "$MEDIA_DIR")"
else
  echo "! media dir $MEDIA_DIR not found — skipping media archive"
fi

echo "✓ snapshot ready in $OUT/ — copy BOTH files to the target server, then run data-restore.sh there."
