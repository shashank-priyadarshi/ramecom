# ecommerce-microservices
#
# Usage:
#   make help
#   make test
#   make build
#   make run APP=gateway
#   make watch APP=auth
#   make up / make down
#
# APP must be one of: gateway auth user order payment notification

APP        ?= gateway
BIN_DIR    ?= bin
BIN        ?= $(BIN_DIR)/ecom
GO         ?= go
DOCKER     ?= docker
COMPOSE    ?= docker compose
AIR        ?= air
AIR_CFG    ?= build/.air.toml
AIR_PKG    ?= github.com/air-verse/air@latest
COVER_OUT  ?= coverage.out
# Avoid env var LDFLAGS (often set by Homebrew/Postgres) — use GO_LDFLAGS only.
GO_LDFLAGS ?= -s -w

APPS := gateway auth user order payment notification

.DEFAULT_GOAL := help

.PHONY: help \
	test test-race test-cover test-verbose \
	build build-docker \
	run run-all run-docker \
	up down restart logs ps \
	infra infra-down \
	watch \
	tidy fmt vet clean tools check-app

## help: Show this help
help:
	@echo "Targets:"
	@echo "  test            Run all tests"
	@echo "  test-race       Run tests with race detector"
	@echo "  test-cover      Run tests with coverage report"
	@echo "  test-verbose    Run tests verbose"
	@echo "  build           Build binary to $(BIN) (host)"
	@echo "  build-docker    Build app image via docker compose"
	@echo "  run             Run one app on host (APP=$(APP))"
	@echo "  run-all         Run all apps on host (background PIDs in tmp/run)"
	@echo "  run-docker      Start one app service in compose (APP=$(APP))"
	@echo "  up              docker compose up --build -d (full stack)"
	@echo "  down            docker compose down"
	@echo "  restart         Restart all compose services"
	@echo "  logs            Follow compose logs (optional SERVICE=...)"
	@echo "  ps              Compose process status"
	@echo "  infra           Start infra only (zookeeper kafka mariadb prometheus grafana)"
	@echo "  infra-down      Stop infra services"
	@echo "  watch           Live reload with air (APP=$(APP), config $(AIR_CFG))"
	@echo "  tidy            go mod tidy"
	@echo "  fmt             go fmt ./..."
	@echo "  vet             go vet ./..."
	@echo "  clean           Remove bin/, tmp/, coverage"
	@echo "  tools           Install air (live reload)"
	@echo ""
	@echo "Variables: APP=$(APP)  BIN=$(BIN)  COMPOSE=\"$(COMPOSE)\""
	@echo "Apps: $(APPS)"

# -----------------------------------------------------------------------------
# Tests
# -----------------------------------------------------------------------------

## test: Run all Go tests
test:
	$(GO) test ./...

## test-race: Run tests with -race
test-race:
	$(GO) test -race ./...

## test-cover: Run tests with coverage
test-cover:
	$(GO) test -coverprofile=$(COVER_OUT) ./...
	$(GO) tool cover -func=$(COVER_OUT)

## test-verbose: Verbose tests
test-verbose:
	$(GO) test -v ./...

# -----------------------------------------------------------------------------
# Build
# -----------------------------------------------------------------------------

## build: Build host binary
build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -ldflags="$(GO_LDFLAGS)" -o $(BIN) ./cmd
	@echo "built $(BIN)"

## build-docker: Build compose app image (ecom:local)
build-docker:
	$(COMPOSE) build

# -----------------------------------------------------------------------------
# Run (host)
# -----------------------------------------------------------------------------

check-app:
	@echo "$(APPS)" | tr ' ' '\n' | grep -qx "$(APP)" \
		|| (echo "error: APP must be one of: $(APPS) (got '$(APP)')" >&2; exit 1)

## run: Run a single app on the host (requires infra; uses build/.env/.<app>.env)
run: check-app
	$(GO) run ./cmd -app $(APP)

## run-all: Start all apps on host in background (logs under tmp/run/)
run-all:
	@mkdir -p tmp/run
	@for a in $(APPS); do \
		echo "starting $$a..."; \
		$(GO) run ./cmd -app $$a > tmp/run/$$a.log 2>&1 & echo $$! > tmp/run/$$a.pid; \
	done
	@echo "pids in tmp/run/*.pid  logs in tmp/run/*.log"
	@echo "stop with: make run-all-stop"

.PHONY: run-all-stop
run-all-stop:
	@if [ -d tmp/run ]; then \
		for f in tmp/run/*.pid; do \
			[ -f "$$f" ] || continue; \
			pid=$$(cat "$$f"); \
			kill "$$pid" 2>/dev/null || true; \
			rm -f "$$f"; \
		done; \
		echo "stopped host apps"; \
	else \
		echo "nothing to stop"; \
	fi

# -----------------------------------------------------------------------------
# Run (docker)
# -----------------------------------------------------------------------------

## run-docker: Start/recreate one app container (APP=...)
run-docker: check-app
	$(COMPOSE) up -d --build $(APP)

## up: Full stack (infra + all apps)
up:
	$(COMPOSE) up --build -d

## down: Stop full stack
down:
	$(COMPOSE) down

## restart: Restart all services
restart:
	$(COMPOSE) restart

## logs: Follow logs (SERVICE=gateway optional)
logs:
	$(COMPOSE) logs -f $(SERVICE)

## ps: Show compose status
ps:
	$(COMPOSE) ps

## infra: Start dependencies only
infra:
	$(COMPOSE) up -d zookeeper kafka mariadb prometheus grafana

## infra-down: Stop dependency services only
infra-down:
	$(COMPOSE) stop zookeeper kafka mariadb prometheus grafana

# -----------------------------------------------------------------------------
# Live reload (air)
# -----------------------------------------------------------------------------

## tools: Install air
tools:
	$(GO) install $(AIR_PKG)
	@echo "air installed (ensure $$(go env GOPATH)/bin is on PATH)"

## watch: Live-reload one app with air (APP=gateway|auth|...)
watch: check-app
	@command -v $(AIR) >/dev/null 2>&1 || { \
		echo "air not found; installing $(AIR_PKG)..."; \
		$(GO) install $(AIR_PKG); \
	}
	@mkdir -p tmp
	$(AIR) -c $(AIR_CFG) --build.full_bin "./tmp/main -app $(APP)"

# -----------------------------------------------------------------------------
# Hygiene
# -----------------------------------------------------------------------------

## tidy: go mod tidy
tidy:
	$(GO) mod tidy

## fmt: Format packages
fmt:
	$(GO) fmt ./...

## vet: Static checks
vet:
	$(GO) vet ./...

## clean: Remove build artifacts
clean:
	rm -rf $(BIN_DIR) tmp $(COVER_OUT)
	@echo "cleaned"
