#!/usr/bin/env bash
#
# data-restore.sh — restore a snapshot (from data-snapshot.sh) onto THIS server.
#
# Restores the Postgres database and extracts the uploaded media (incl. avatars).
# The target database should be EMPTY and freshly created — do NOT run
# `go run ./cmd/migrate` first; the dump already carries the schema at HEAD, so
# the app will see no pending migrations on boot.
#
# Usage:
#   SHANRAQ_DATABASE_URL=postgres://user:pass@host:5432/shanraq \
#   SHANRAQ_MEDIA_PARENT=./data \
#   scripts/data-restore.sh <snapshot-dir>
#
# SHANRAQ_MEDIA_PARENT is the PARENT of the media dir (the archive holds a
# top-level "media/"), and must match media.dir in the target's config
# (default media.dir = ./data/media → SHANRAQ_MEDIA_PARENT=./data).
#
set -euo pipefail

SNAP="${1:?usage: data-restore.sh <snapshot-dir>}"
DB_URL="${SHANRAQ_DATABASE_URL:?set SHANRAQ_DATABASE_URL to the TARGET database}"
MEDIA_PARENT="${SHANRAQ_MEDIA_PARENT:-./data}"

PGRESTORE="pg_restore"
if [[ -n "${PG_BIN:-}" ]]; then PGRESTORE="$PG_BIN/pg_restore"
elif [[ -x /opt/homebrew/opt/postgresql@16/bin/pg_restore ]]; then PGRESTORE=/opt/homebrew/opt/postgresql@16/bin/pg_restore
elif [[ -x /usr/local/opt/postgresql@16/bin/pg_restore ]]; then PGRESTORE=/usr/local/opt/postgresql@16/bin/pg_restore
fi

echo "→ restoring database ($("$PGRESTORE" --version))"
"$PGRESTORE" --no-owner --no-privileges --clean --if-exists -d "$DB_URL" "$SNAP/db.dump"

if [[ -f "$SNAP/media.tar.gz" ]]; then
  echo "→ extracting media into $MEDIA_PARENT"
  mkdir -p "$MEDIA_PARENT"
  tar -xzf "$SNAP/media.tar.gz" -C "$MEDIA_PARENT"
fi

echo "✓ restore complete. Start the app — it should report 'no migrations to run'."
echo "! Reminder: change the passwords of the seeded staff accounts on this server."
