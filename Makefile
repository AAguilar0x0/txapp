MAIN_PACKAGE_PATH := ./bin

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## git/hooks: setup git hooks
.PHONY: git/hooks
git/hooks:
	chmod +x .git-hooks/* && cp -R .git-hooks/* .git/hooks/

## audit: run quality control checks
.PHONY: audit
audit: check
	go mod verify
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## check: run code maintenance tasks (tidy dependencies, verify, clean, and format code)
.PHONY: check
check:
	go mod tidy -v
	go vet ./...
	go clean
	go fmt ./...

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## setup: install and configure necessary dependencies for development
.PHONY: setup
setup: install install/bin git/hooks

## install: install necessary dependencies
.PHONY: install
install:
	go mod download
	go mod verify

## install/bin: install binary dependency
.PHONY: install/bin
install/bin:
	go install github.com/air-verse/air@latest
	go install github.com/swaggo/swag/app/swag@latest
	go install github.com/sqlc-dev/sqlc/app/sqlc@latest
	go install github.com/a-h/templ/app/templ@latest

## docker: run the application with docker
.PHONY: docker
docker:
	docker compose up -d

## docker/build: rebuild the backend image then run the application with docker
.PHONY: docker/build
docker/build:
	docker compose up -d --build backend

## docker/db: run the db with docker
.PHONY: docker/db
docker/db:
	docker compose up -d db

## sqlc: sqlc generate
.PHONY: sqlc
sqlc:
	sqlc generate

## templ: templ generate
.PHONY: templ
templ:
	templ generate

## tailwindcss: generate css from tailwindcss
.PHONY: tailwindcss
tailwindcss:
	npx tailwindcss -i ./app/web/static/input.css -o ./app/web/static/output.css --minify

## swagger: generate swagger docs
.PHONY: swagger
swagger:
	(cd ./app/web && swag init --parseDependency)


# ==================================================================================== #
# APPLICATIONS
# ==================================================================================== #

## app/web/build: build the web application
.PHONY: app/web/build
app/web/build:
	go build -v -o=${MAIN_PACKAGE_PATH}/web ./app/web

## app/web/bin: execute the web application binary
.PHONY: app/web/bin
app/web/bin:
	${MAIN_PACKAGE_PATH}/web

## app/web/live: run the application with reloading on file changes
.PHONY: app/web/live
app/web/live:
	air

## app/migrate/run: run the migrate application
.PHONY: app/migrate/run
app/migrate/run:
	go run ./app/migrate/main.go

## app/sqlcore/run: run the sqlc to core codegen
.PHONY: app/sqlcore/run
app/sqlcore/run:
	go run ./app/sqlcore/main.go
	go fmt ./core/models
