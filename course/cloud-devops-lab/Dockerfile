# syntax=docker/dockerfile:1
FROM golang:1.26.6-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
ARG VERSION=dev
RUN CGO_ENABLED=0 go test ./... && \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o /out/cloudlab ./cmd/cloudlab && \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/cloudlab-healthcheck ./cmd/healthcheck && \
    mkdir -p /out/data && chown 65532:65532 /out/data

FROM scratch
COPY --from=build /out/cloudlab /cloudlab
COPY --from=build /out/cloudlab-healthcheck /cloudlab-healthcheck
COPY --chown=65532:65532 --from=build /out/data /data
USER 65532:65532
EXPOSE 8080
VOLUME ["/data"]
ENTRYPOINT ["/cloudlab"]
