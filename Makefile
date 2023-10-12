.ONESHELL:
SHELL := /bin/bash

install:
	go mod vendor

test:
	go clean -testcache && go test -cover -v ./...

coverage:
	go clean -testcache
	go test -coverprofile=c.out ./...
	go tool cover -func=c.out
	go tool cover -html=c.out