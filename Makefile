GO ?= go

.PHONY: run
run:
	$(GO) run ./cmd/app -config configs/config.example.yaml

.PHONY: test
test:
	$(GO) test ./...

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: smoke
smoke:
	./scripts/docker-smoke.sh

.PHONY: snapshots
snapshots:
	go generate ./web
