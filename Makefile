GO ?= go

.PHONY: run
run:
	$(GO) run ./cmd/app -config configs/config.example.yaml

# -p 1: the packages share one database. Run in parallel, as `go test` does by
# default, the articles harness creates an administrator while the auth package
# is counting them, and the last-admin guard tests see three where they staged
# two — sometimes refusing nothing, sometimes refusing a demotion that was
# legitimate. Both were reproduced locally at roughly one run in six, which is
# the rate CI was failing at. Serial packages cost a couple of seconds.
#
# Every database-backed test skips itself when SHANRAQ_TEST_DB is unset, and a
# skip is invisible behind a `| tail -3`. That is not a hypothetical: on
# 16 August 2026 six local runs in a row were reported as passing while every
# integration test had quietly not run, and the first honest execution happened
# in CI. So the default target refuses to pretend, and mirrors CI's unit job
# exactly — plain pass and race pass — so "make test is green" and "CI will be
# green" mean the same thing.
.PHONY: test
test: require-test-db
	$(GO) test -p 1 ./...
	$(GO) test -race -p 1 ./...

# The plain pass alone, for a tight edit-run loop. Still refuses to run blind.
.PHONY: test-quick
test-quick: require-test-db
	$(GO) test -p 1 ./...

# Deliberately without a database, framed so it cannot be mistaken for a full
# run. Useful when the change genuinely touches no storage.
.PHONY: test-nodb
test-nodb:
	@echo "── no database: integration tests will SKIP, this is not the full suite ──"
	@SHANRAQ_TEST_DB= $(GO) test -p 1 ./...
	@echo "── finished WITHOUT the integration tests ──"

.PHONY: require-test-db
require-test-db:
	@if [ -z "$$SHANRAQ_TEST_DB" ]; then \
	  echo ""; \
	  echo "SHANRAQ_TEST_DB is not set, so every database-backed test would skip"; \
	  echo "itself and the run would report a green that means nothing."; \
	  echo ""; \
	  echo "  createdb -T template0 -l en_US.UTF-8 -E UTF8 shanraq_test"; \
	  echo "      (the locale is not decoration: on a C-locale database"; \
	  echo "       lower() leaves Cyrillic alone and the name search fails)"; \
	  echo "  export DATABASE_URL=\"postgres://$$(whoami)@localhost:5432/shanraq_test?sslmode=disable\""; \
	  echo "  go run ./cmd/migrate"; \
	  echo "  export SHANRAQ_TEST_DB=\"\$$DATABASE_URL\""; \
	  echo ""; \
	  echo "The name must contain \"test\": the harness refuses any other database."; \
	  echo "To run only what needs no database:  make test-nodb"; \
	  echo ""; \
	  exit 1; \
	fi

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: smoke
smoke:
	./scripts/docker-smoke.sh

# A backup nobody has restored is a hope. Run this after every schema change,
# on the machine that holds the age private key.
#   make restore-test ARCHIVE=~/backups/shanraq-backup-20260819-023000Z.tar.gz.age
.PHONY: restore-test
restore-test:
	@if [ -z "$(ARCHIVE)" ]; then \
	  echo "usage: make restore-test ARCHIVE=<path to shanraq-backup-*.tar.gz[.age]>"; \
	  echo "       BACKUP_AGE_IDENTITY must point at the age private key for .age archives"; \
	  exit 1; \
	fi
	./scripts/backup-restore-test.sh "$(ARCHIVE)"

.PHONY: snapshots
snapshots:
	go generate ./web
