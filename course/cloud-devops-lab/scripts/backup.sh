#!/bin/sh
set -eu

destination=${1:-backups}
stamp=$(date -u +%Y%m%dT%H%M%SZ)
volume=${CLOUDLAB_VOLUME:-cloud-devops-lab_cloudlab-data}
mkdir -p "$destination"
destination=$(cd "$destination" && pwd)
archive=cloudlab-$stamp.tar.gz
docker run --rm \
  -v "$volume:/source:ro" \
  -v "$destination:/backup" \
  alpine:3.23 tar -czf "/backup/$archive" -C /source .
(
  cd "$destination"
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$archive" >"$archive.sha256"
  else
    shasum -a 256 "$archive" >"$archive.sha256"
  fi
)
printf 'Backup written: %s/%s\n' "$destination" "$archive"
