TEST ?= ./...

default: build

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

gitfetch:
	git fetch

major: build gitfetch
	eval `./git-semv major --bump`

minor: build gitfetch
	eval `./git-semv minor --bump`

patch: build gitfetch
	eval `./git-semv patch --bump`

dist:
	@test -z $(GITHUB_TOKEN) || goreleaser --rm-dist

clean:
	rm -rf coverage.txt
	git checkout go.*
