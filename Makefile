export GOBIN := $(CURDIR)/out

.PHONY: build install-goose migrate-up migrate-down migrate-reset deploy

# GOOSE
GOOSE := $(GOBIN)/goose
GOOSE_VERSION := v3.27.2
GOOSE_DRIVER := sqlite3
GOOSE_DBSTRING := ./internal/db/sqlite.db
GOOSE_MIGRATION_DIR := ./internal/db/migrations
GOOSE_ARGS := -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING)

build:
	@GOOS=linux GOARCH=amd64 go build -o $(GOBIN)/currency_exchanger ./cmd

install-goose:
	@go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)

migrate-up:
	@$(GOOSE) $(GOOSE_ARGS) up

migrate-down:
	@$(GOOSE) $(GOOSE_ARGS) down

migrate-reset:
	@$(GOOSE) $(GOOSE_ARGS) reset

deploy: install-goose migrate-up build