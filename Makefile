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

## audit: run quality control checks
.PHONY: audit
audit: tidy test
	go mod verify
	test -z "$(shell gofmt -l .)" 
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## tidy: cleanup and format code and tidy modfile
.PHONY: tidy
tidy:
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

## install: install necessary dependencies
.PHONY: install
install:
	go install github.com/swaggo/swag/cmd/swag@latest
	go mod download
	go mod verify

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0

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
	npx tailwindcss -i ./cmd/web/static/input.css -o ./cmd/web/static/output.css --minify


# ==================================================================================== #
# COMMANDS
# ==================================================================================== #

## cmd/web/build: build the web application
.PHONY: cmd/web/build
cmd/web/build:
	go build -v -o=${MAIN_PACKAGE_PATH}/web ./cmd/web

## cmd/web/run: run the web application
.PHONY: cmd/web/run
cmd/web/run:
	${MAIN_PACKAGE_PATH}/web

