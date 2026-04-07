# Load environment variables from .env if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

#-----------------------------------------#
###         Linting, formatting 		###
#-----------------------------------------#

.PHONY: lint-install
lint-install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.5

.PHONY: lint
lint:
	golangci-lint run --max-issues-per-linter=0 --max-same-issues=0 ./...

.PHONY: fmt
fmt:
	golangci-lint fmt ./...


#-----------------------------------------#
###         Database Migrations         ###
#-----------------------------------------#

.PHONY: migrate-create
migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir "./migrations" create $$name sql

.PHONY: migrate-up
migrate-up:
	go run ./cmd migrate up

.PHONY: migrate-down
migrate-down:
	go run ./cmd migrate down


#-----------------------------------------#
###               Test                  ###
#-----------------------------------------#

.PHONY: test
test:
	ENVIRONMENT=test go test -race ./pkg/...

.PHONY: test-system
test-system:
	./scripts/prepare-system-test.sh
	ENVIRONMENT=test go test -race -tags=system -p=1 ./tests/system/...


#-----------------------------------------#
###             Build, Run              ###
#-----------------------------------------#

.PHONY: infra-up
infra-up:
	docker-compose -f dev-infra.yaml --profile full up -d --build

.PHONY: infra-down
infra-down:
	docker-compose -f dev-infra.yaml --profile full down --remove-orphans

.PHONY: run
run: infra-up
	go run ./cmd run
