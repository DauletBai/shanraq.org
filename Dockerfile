# syntax=docker/dockerfile:1

FROM golang:1.26-bookworm AS builder
WORKDIR /src

ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/shanraq ./cmd/app

# Pre-create the media directory here (the distroless runtime has no shell to
# mkdir at start) so it can be COPYed in with the nonroot owner below.
RUN mkdir -p /app/data/media

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=builder /bin/shanraq /usr/local/bin/shanraq
COPY configs/config.example.yaml /app/config.yaml
# Writable media tree owned by the nonroot user (uid/gid 65532 in distroless).
COPY --from=builder --chown=65532:65532 /app/data /app/data

ENV SHANRAQ_LOGGING_MODE=production \
    SHANRAQ_CONFIG=/app/config.yaml

# Persist uploaded media across container restarts (mount a volume here).
VOLUME ["/app/data"]

ENTRYPOINT ["/usr/local/bin/shanraq"]
CMD ["-config", "/app/config.yaml"]
