## [command]        # [definition]
## run              # runs application
## migrate/version  # prints migration tool version
## migrate/up       # runs up migrations
## migrate/down     # runs down migrations
## migrate/create   # creates new migration files

DB_DSN := postgres://movies_api:pa55word@localhost/movies_api
GO_PATH := ~/go/bin/

.SILENT:

# ------------
# Helper:

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ------------
# Application commands:

.PHONY: run
run:
	go run ./cmd/api
	@echo "- run finished"

# ------------
# Migration commands:

.PHONY: migrate/version
migrate/version: 
	${GO_PATH}/migrate -version

.PHONY: migrate/up
migrate/up: confirm
	${GO_PATH}/migrate -path=./migrations -database=${DB_DSN} up
	@echo "- migrate/up finished"

.PHONY: migrate/down
migrate/down: confirm
	${GO_PATH}/migrate -path=./migrations -database=${DB_DSN} down
	@echo "- migrate/down finished"

.PHONY: migrate/create
migrate/create:
	${GO_PATH}/migrate create -seq -ext=.sql -dir=./migrations ${name}
	@echo "- migrate/create finished"

# ------------
