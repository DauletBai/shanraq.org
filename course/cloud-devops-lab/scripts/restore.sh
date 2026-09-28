#!/bin/sh
set -eu

archive=${1:?usage: restore.sh backups/cloudlab-TIMESTAMP.tar.gz}
test -f "$archive"
test -f "$archive.sha256"
directory=$(cd "$(dirname "$archive")" && pwd)
name=$(basename "$archive")
volume=${CLOUDLAB_VOLUME:-cloud-devops-lab_cloudlab-data}
(
  cd "$directory"
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum -c "$name.sha256"
  else
    shasum -a 256 -c "$name.sha256"
  fi
)
docker compose stop app
docker run --rm \
  -v "$volume:/target" \
  -v "$directory:/backup:ro" \
  alpine:3.23 sh -c 'rm -rf /target/* && tar -xzf "/backup/'"$name"'" -C /target'
docker compose up -d app
printf 'Restore completed from %s\n' "$archive"
