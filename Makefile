GO_PATH := ~/go/bin/
# better approach: $HOME/.profile
MOVIES_API_DB_DSN := postgres://movies_api:pa55word@localhost/movies_api
BINARY_NAME := movie-store

.SILENT:

# ------------
# Helpers:
# ------------

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ------------
# Application commands:
# ------------

## run: run the cmd/api application
.PHONY: run
run:
	go run ./cmd/api -db-dsn=${MOVIES_API_DB_DSN}
	@echo "- run finished"

## build: build the cmd/api application
.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/${BINARY_NAME} ./cmd/api
	@echo "- build finished"

# ------------
# Migration commands:
# ------------

## migrate/version: prints migration tool version
.PHONY: migrate/version
migrate/version: 
	${GO_PATH}/migrate -version

## migrate/up: runs up migrations
.PHONY: migrate/up
migrate/up: confirm
	${GO_PATH}/migrate -path=./migrations -database=${DB_DSN} up
	@echo "- migrate/up finished"

## migrate/down: runs down migrations
.PHONY: migrate/down
migrate/down: confirm
	${GO_PATH}/migrate -path=./migrations -database=${DB_DSN} down
	@echo "- migrate/down finished"

## migrate/create: creates new migration files
.PHONY: migrate/create
migrate/create:
	${GO_PATH}/migrate create -seq -ext=.sql -dir=./migrations ${name}
	@echo "- migrate/create finished"

# ------------
# Quality Control:
# ------------

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	${GO_PATH}/staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...
	@echo "- audit finished"

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor
	@echo "- vendor finished"

# ------------
# Additional Comments section:
# ------------

# - {-ldflags='-s} flag in build command, removes DWARF information and symbol table from the binary.
# - {-a} flag in build command, forces all packages to be rebuilt.
# - {go clean -cache} command, removes everything from the build cache.