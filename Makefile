TEST ?= ./...

default: build

deps: export GO111MODULE=off
deps:
	go get golang.org/x/lint/golint
	go get github.com/pierrre/gotestcover
	go get github.com/goreleaser/goreleaser

build:
	go build ./cmd/git-semv

test:
	go test $(TEST) $(TESTARGS)
	go test -race $(TEST) $(TESTARGS) -coverprofile=coverage.txt -covermode=atomic

lint:
	golint -set_exit_status $(TEST)

ci: deps test lint
