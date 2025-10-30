# syntax=docker/dockerfile:1

FROM golang:1.23-bullseye AS builder
WORKDIR /src

ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/shanraq ./cmd/app

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=builder /bin/shanraq /usr/local/bin/shanraq
COPY configs/config.example.yaml /app/config.yaml

ENV SHANRAQ_LOGGING_MODE=production \
    SHANRAQ_CONFIG=/app/config.yaml

ENTRYPOINT ["/usr/local/bin/shanraq"]
CMD ["-config", "/app/config.yaml"]
