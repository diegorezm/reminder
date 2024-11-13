# Variables
HOME_DIR=$(shell echo $$HOME)
DB_PATH=$(HOME_DIR)/.local/share/reminder/reminder.db
DB_DRIVER=sqlite3
DB_PATH=$(HOME_DIR)/.local/share/reminder/reminder.db
GOOSE_DBSTRING=file:$(DB_PATH)?cache=shared
GOOSE_MIGRATION_DIR=./internal/migrations

# Targets
.PHONY: build run clean up reset create

templ:
	@templ generate

build: templ
	@go build -o ./bin/reminder-cli ./cmd/reminder-cli/main.go
	@chmod +x ./bin/reminder-cli

run: build
	./bin/reminder-cli

server: build
	./bin/reminder-cli serve

clean:
	@rm -rf ./bin
	@rm -rf $(DB_PATH)

up:
	@GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) \
	GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose up

reset:
	@GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) \
	GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose reset

create:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) \
		goose create $(NAME) go