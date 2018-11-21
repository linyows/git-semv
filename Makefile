default: build

build:
	go build ./cmd/git-semv

test:
	go test ./...
