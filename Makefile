GO ?= go

all: test

.PHONY: test
test:
	$(GO) test -v ./... -timeout=10s

.PHONY: coverage
coverage:
	$(GO) test -v ./... -covermode=count -coverprofile=coverage.out -timeout=20s

.PHONY: build
build:
	goreleaser build --snapshot --rm-dist --single-target

.PHONY: plugin
plugin: build
	mv polr /usr/local/bin/kubectl-polr