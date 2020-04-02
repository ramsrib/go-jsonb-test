SHELL := /bin/bash
GOPATH = $(shell go env GOPATH)

.PHONY: build
build:
	docker build -t go-jsonb-test .

run-local:
	source local.env && go run .

run-db:
	docker-compose up -d db
	# wait for db to finish the initialization
	sleep 5;

run run-docker: run-db
	docker-compose up --build app

test:
	source local.env && go test .

clean clean-docker:
	docker-compose down -v

