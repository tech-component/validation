.PHONY: docker-logs docker-start docker-stop
.PHONY: golang-imports migrate-new migrate-install
.PHONY: mock-gen mock-install test test-html-output

MIGRATE_VERSION=v4.18.1
SHELL := /bin/bash

include .env
export

docker-build:
	docker compose build

docker-logs:
	docker compose logs -f

docker-start:
	docker compose up -d

docker-stop:
	docker compose down --remove-orphans

golang-imports:
	goimports -w .

install:
	go install

migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)

migrate-new:
	migrate create -ext sql -dir assets/files/migrations $(NAME)

mock-gen:
	rm -rf mocks && mkdir -p mocks && go generate ./...

mock-install:
	go install github.com/matryer/moq@v0.5.1

test: mock-gen
	go test -cover ./...

test-html-output:
	go test -coverprofile=c.out ./... && go tool cover -html=c.out && rm -f c.out
