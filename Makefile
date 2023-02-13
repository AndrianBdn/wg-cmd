VERSION=0.0
.PHONY: all fmt static precommit

all:
	go build -ldflags "-X main.version=0.10"

fmt:
	gofumpt -l -w .

static:
	go vet ./...
	staticcheck ./...

precommit: fmt static