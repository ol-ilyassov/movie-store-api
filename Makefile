## [command]        # [definition]
## run              # runs application
## migrate/version  # prints migration tool version
## migrate/up       # runs up migrations
## migrate/down     # runs down migrations

DB_DSN := postgres://movies_api:pa55word@localhost/movies_api
GO_PATH := ~/go/bin/

.SILENT:

.PHONY: run
run:
	go run ./cmd/api
	@echo "- run finished"

.PHONY: migrate/version
migrate/version: 
	${GO_PATH}/migrate -version

.PHONY: migrate/up
migrate/up: 
	${GO_PATH}/migrate -path=./migrations -database=${DB_DSN} up
	@echo "- migrate/up finished"

.PHONY: migrate/down
migrate/down: 
	${GO_PATH}/migrate -path=./migrations -database=${DB_DSN} down
	@echo "- migrate/down finished"

.PHONY: migrate/create
migrate/create:
	${GO_PATH}/migrate create -seq -ext=.sql -dir=./migrations ${name}
	@echo "- migrate/create finished"