.PHONY: docker-start docker-stop
.PHONY: golang-imports migrate-new migrate-install
.PHONY: test test-html-output

MIGRATE_VERSION=v4.18.1
SHELL := /bin/bash

include .env
export

docker-start:
	docker compose up -d

docker-stop:
	docker compost down --remove-orphans

golang-imports:
	goimports -w .

migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)

migrate-new:
	migrate create -ext sql -dir assets/files/migrations $(NAME)

test:
	go test -cover ./...

test-html-output:
	go test -coverprofile=c.out ./... && go tool cover -html=c.out && rm -f c.out
