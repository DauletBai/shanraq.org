# syntax=docker/dockerfile:1

FROM golang:1.27-bookworm AS builder
WORKDIR /src

ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
# The public module proxy occasionally resets a single HTTP/2 stream on a fresh
# GitHub runner. Everything already downloaded remains in the module cache, so
# retry the missing archive instead of turning a transient network reset into a
# failed release.
RUN for attempt in 1 2 3; do \
      go mod download && exit 0; \
      status=$?; \
      if [ "${attempt}" -eq 3 ]; then exit "${status}"; fi; \
      echo "go mod download failed (attempt ${attempt}/3); retrying..." >&2; \
      sleep $((attempt * 3)); \
    done

COPY . .
# TARGETOS/TARGETARCH are provided by BuildKit (default to the build host, or set
# via `docker build --platform linux/amd64`). This keeps the image arch and the
# binary arch consistent — a plain `docker build` on an ARM Mac no longer bakes an
# amd64 binary into an arm64 image.
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/shanraq ./cmd/app
# The content exporter ships in the same image so it can be run against the
# live database without a toolchain on the host: docker compose exec app export
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/export ./cmd/export

# Pre-create the media directory here (the distroless runtime has no shell to
# mkdir at start) so it can be COPYed in with the nonroot owner below.
RUN mkdir -p /app/data/media

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=builder /bin/shanraq /usr/local/bin/shanraq
COPY --from=builder /bin/export /usr/local/bin/export
COPY configs/config.example.yaml /app/config.yaml
# Writable media tree owned by the nonroot user (uid/gid 65532 in distroless).
COPY --from=builder --chown=65532:65532 /app/data /app/data

ENV SHANRAQ_LOGGING_MODE=production \
    SHANRAQ_CONFIG=/app/config.yaml

# Persist uploaded media across container restarts (mount a volume here).
VOLUME ["/app/data"]

ENTRYPOINT ["/usr/local/bin/shanraq"]
CMD ["-config", "/app/config.yaml"]
