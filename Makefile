MIGRATIONS_BIN_PATH=./bin/migrations
APP_BIN_PATH=./bin/app

SRC_MIGRATIONS=./cmd/migrations
SRC_APP=./cmd/app

# Migrations
.PHONY: build_migrations
build_migrations:
	go build -o $(MIGRATIONS_BIN_PATH) $(SRC_MIGRATIONS)

.PHONY: up_migrations
up_migrations: build_migrations
	./$(MIGRATIONS_BIN_PATH) -action=up

.PHONY: down_migrations
down_migrations: build_migrations
	./$(MIGRATIONS_BIN_PATH) -action=down

# Application
.PHONY: build_app
build_app:
	go build -o $(APP_BIN_PATH) $(SRC_APP)

.PHONY: run_app
run_app: build_app
	./$(APP_BIN_PATH)

# Tests
.PHONY: run_tests
run_tests:
	go test ./...
